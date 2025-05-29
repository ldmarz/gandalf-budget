package store

import (
	"database/sql" // Used by sqlx, good to have explicitly if direct use is ever needed.
	"fmt"
	"time"
)

func (s *sqlStore) CanFinalizeMonth(monthID int) (bool, string, error) {
	var count int
	query := `
	SELECT COUNT(bl.id)
	FROM budget_lines bl
	INNER JOIN actual_lines al ON bl.id = al.budget_line_id 
	WHERE bl.month_id = ? AND al.actual = 0;
	`

	err := s.DB.Get(&count, query, monthID)
	if err != nil {
		return false, "", fmt.Errorf("error checking finalization status for month %d: %w", monthID, err)
	}

	if count > 0 {
		return false, fmt.Sprintf("%d budget lines still have zero actuals.", count), nil
	}
	return true, "", nil
}

func (s *sqlStore) FinalizeMonth(monthID int, snapJSON string) (int64, error) {
	tx, err := s.DB.Beginx()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	var committed bool = false
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	createdAt := time.Now().Format("2006-01-02 15:04:05")
	_, err = tx.Exec(`
		INSERT INTO annual_snaps (month_id, snap_json, created_at)
		VALUES (?, ?, ?);`, monthID, snapJSON, createdAt)
	if err != nil {
		return 0, fmt.Errorf("failed to create annual snap for month %d: %w", monthID, err)
	}

	_, err = tx.Exec(`UPDATE months SET finalized = 1 WHERE id = ?;`, monthID)
	if err != nil {
		return 0, fmt.Errorf("failed to mark month %d as finalized: %w", monthID, err)
	}

	var currentMonth Month
	err = tx.Get(&currentMonth, `SELECT id, year, month, finalized FROM months WHERE id = ?;`, monthID)
	if err != nil {
		return 0, fmt.Errorf("failed to get current month details for month %d: %w", monthID, err)
	}

	nextYear, nextMonthVal := currentMonth.Year, currentMonth.Month+1
	if nextMonthVal > 12 {
		nextMonthVal = 1
		nextYear++
	}

	res, err := tx.Exec(`
		INSERT INTO months (year, month, finalized)
		VALUES (?, ?, 0);`, nextYear, nextMonthVal)
	if err != nil {
		return 0, fmt.Errorf("failed to create next month record for %d-%d: %w", nextYear, nextMonthVal, err)
	}
	newMonthID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get ID of new month (%d-%d): %w", nextYear, nextMonthVal, err)
	}

	var budgetLines []BudgetLine
	err = tx.Select(&budgetLines, `
		SELECT category_id, label, expected 
		FROM budget_lines WHERE month_id = ?;`, monthID)
	if err != nil {
		if err == sql.ErrNoRows {
		} else {
			return 0, fmt.Errorf("failed to fetch budget lines for month %d: %w", monthID, err)
		}
	}
	
	for _, bl := range budgetLines {
		clonedLineRes, err := tx.Exec(`
			INSERT INTO budget_lines (month_id, category_id, label, expected)
			VALUES (?, ?, ?, ?);`, newMonthID, bl.CategoryID, bl.Label, bl.Expected)
		if err != nil {
			return 0, fmt.Errorf("failed to clone budget line (label: %s) for new month %d: %w", bl.Label, newMonthID, err)
		}
		newBudgetLineID, err := clonedLineRes.LastInsertId()
		if err != nil {
			return 0, fmt.Errorf("failed to get ID of cloned budget line (label: %s) for new month %d: %w", bl.Label, newMonthID, err)
		}

		_, err = tx.Exec(`
			INSERT INTO actual_lines (budget_line_id, actual)
			VALUES (?, 0);`, newBudgetLineID)
		if err != nil {
			return 0, fmt.Errorf("failed to create actual line for cloned budget line ID %d (label: %s): %w", newBudgetLineID, bl.Label, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction for finalizing month %d and creating month %d: %w", monthID, newMonthID, err)
	}
	committed = true
	return newMonthID, nil
}
