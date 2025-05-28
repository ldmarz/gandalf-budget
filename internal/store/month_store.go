package store

import (
	"database/sql" // Used by sqlx, good to have explicitly if direct use is ever needed.
	"fmt"
	"time"
	// sqlx is implicitly used via s.DB, which is *sqlx.DB
	// "github.com/jmoiron/sqlx"
)

// CanFinalizeMonth checks if a month can be finalized.
// A month can be finalized if all its budget lines have a non-zero actual amount.
// Note: The PRD says "all budget lines have an actual amount" which implies actual_id IS NOT NULL.
// The provided query checks `al.actual = 0`. This means an actual_line record exists, but its value is 0.
// If the intent is that *every* budget line *must* have an actual that is *not* zero, the query is correct.
// If the intent is that *every* budget line *must* have an associated actual_line (even if its value is 0),
// then the query might need adjustment (e.g., checking for budget_lines without a corresponding actual_line).
// Given the PRD's "If any budget lines have an actual amount of 0, the month cannot be finalized", the current query seems correct.
func (s *sqlStore) CanFinalizeMonth(monthID int) (bool, string, error) {
	var count int
	// This query identifies budget lines for the given month that are linked to an actual_line
	// where the actual amount is zero.
	query := `
	SELECT COUNT(bl.id)
	FROM budget_lines bl
	INNER JOIN actual_lines al ON bl.id = al.budget_line_id 
	WHERE bl.month_id = ? AND al.actual = 0;
	`
	// If there are budget lines that *do not have* an actual_lines record yet,
	// they would not be caught by this query. The PRD implies an actual_line record
	// is created when a budget_line is created, possibly with a default of 0.
	// Assuming actual_lines are always present for budget_lines in a month being considered for finalization.

	err := s.DB.Get(&count, query, monthID)
	if err != nil {
		// If no rows are found (e.g. monthID doesn't exist, or no budget lines with actual=0),
		// sql.ErrNoRows might be returned by s.DB.Get if it expects exactly one row.
		// However, COUNT should always return one row (with count 0 if no matches).
		// So, any error here is likely a real issue.
		return false, "", fmt.Errorf("error checking finalization status for month %d: %w", monthID, err)
	}

	if count > 0 {
		return false, fmt.Sprintf("%d budget lines still have zero actuals.", count), nil
	}
	return true, "", nil
}

// FinalizeMonth finalizes a given month and prepares the next month.
// It creates an annual snapshot, marks the month as finalized,
// creates the next month's record, and clones budget lines.
func (s *sqlStore) FinalizeMonth(monthID int, snapJSON string) (int64, error) {
	tx, err := s.DB.Beginx()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	// Using a labeled defer for explicit rollback or commit checking
	var committed bool = false
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	// 1. Create Annual Snap
	// Assuming 'Month' struct and 'BudgetLine' struct are defined in models.go in the same package.
	createdAt := time.Now().Format("2006-01-02 15:04:05") // Standard SQL datetime format
	_, err = tx.Exec(`
		INSERT INTO annual_snaps (month_id, snap_json, created_at)
		VALUES (?, ?, ?);`, monthID, snapJSON, createdAt)
	if err != nil {
		return 0, fmt.Errorf("failed to create annual snap for month %d: %w", monthID, err)
	}

	// 2. Mark current month as finalized
	// Assuming 'months' table has 'id' and 'finalized' columns.
	_, err = tx.Exec(`UPDATE months SET finalized = 1 WHERE id = ?;`, monthID)
	if err != nil {
		return 0, fmt.Errorf("failed to mark month %d as finalized: %w", monthID, err)
	}

	// 3. Determine next month's year and month value
	var currentMonth Month // This is store.Month from models.go
	err = tx.Get(&currentMonth, `SELECT id, year, month, finalized FROM months WHERE id = ?;`, monthID)
	if err != nil {
		return 0, fmt.Errorf("failed to get current month details for month %d: %w", monthID, err)
	}

	nextYear, nextMonthVal := currentMonth.Year, currentMonth.Month+1
	if nextMonthVal > 12 {
		nextMonthVal = 1
		nextYear++
	}

	// 4. Create new month record for next month
	// Assuming 'months' table has 'year', 'month', 'finalized' columns.
	res, err := tx.Exec(`
		INSERT INTO months (year, month, finalized)
		VALUES (?, ?, 0);`, nextYear, nextMonthVal)
	if err != nil {
		// This could fail due to UNIQUE constraint on (year, month) if it exists and next month already exists.
		return 0, fmt.Errorf("failed to create next month record for %d-%d: %w", nextYear, nextMonthVal, err)
	}
	newMonthID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get ID of new month (%d-%d): %w", nextYear, nextMonthVal, err)
	}

	// 5. Get budget lines from current month
	var budgetLines []BudgetLine // This is store.BudgetLine from models.go
	// Selecting only fields needed for cloning.
	// The BudgetLine struct in models.go has ID, MonthID, CategoryID, Label, Expected, ActualID, ActualAmount.
	// We only need CategoryID, Label, Expected for cloning the line itself.
	err = tx.Select(&budgetLines, `
		SELECT category_id, label, expected 
		FROM budget_lines WHERE month_id = ?;`, monthID) // Removed bl.id as it's not used for insertion
	if err != nil {
		if err == sql.ErrNoRows {
			// No budget lines to clone, which is a valid scenario.
			// The rest of the loop will be skipped.
		} else {
			return 0, fmt.Errorf("failed to fetch budget lines for month %d: %w", monthID, err)
		}
	}
	
	// 6. Clone budget lines and create new actual lines for the new month
	for _, bl := range budgetLines {
		// The BudgetLine struct has CategoryID, Label, Expected.
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

		// Create a corresponding actual_lines record with actual = 0
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
	committed = true // Mark as committed
	return newMonthID, nil
}
