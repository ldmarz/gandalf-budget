package store

import (
	"fmt"
	"math"
)

// CreateBudgetLine inserts a new BudgetLine into the budget_lines table
// and creates an associated ActualLine with an actual amount of 0.
// It returns the ID of the newly created BudgetLine.
func (s *sqlStore) CreateBudgetLine(b *BudgetLine) (int64, error) {
	tx, err := s.DB.Beginx()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback in case of error

	// Insert BudgetLine
	stmt, err := tx.PrepareNamed(`
		INSERT INTO budget_lines (month_id, category_id, label, expected)
		VALUES (:month_id, :category_id, :label, :expected)
		RETURNING id`)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare budget_lines insert statement: %w", err)
	}
	defer stmt.Close()

	var budgetLineID int64
	if err := stmt.Get(&budgetLineID, b); err != nil {
		return 0, fmt.Errorf("failed to execute budget_lines insert statement: %w", err)
	}

	// Create associated ActualLine
	actualLine := &ActualLine{
		BudgetLineID: budgetLineID,
		Actual:       0,
	}
	stmtActual, err := tx.PrepareNamed(`
		INSERT INTO actual_lines (budget_line_id, actual)
		VALUES (:budget_line_id, :actual)`)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare actual_lines insert statement: %w", err)
	}
	defer stmtActual.Close()

	if _, err := stmtActual.Exec(actualLine); err != nil {
		return 0, fmt.Errorf("failed to execute actual_lines insert statement: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return budgetLineID, nil
}

// GetBudgetLinesByMonthID retrieves all BudgetLines for a given month_ID,
// including associated actual line data.
func (s *sqlStore) GetBudgetLinesByMonthID(monthID int) ([]BudgetLine, error) {
	var budgetLines []BudgetLine
	query := `
		SELECT
			bl.id, bl.month_id, bl.category_id, bl.label, bl.expected,
			al.id AS actual_id, al.actual AS actual_amount
		FROM budget_lines bl
		LEFT JOIN actual_lines al ON bl.id = al.budget_line_id
		WHERE bl.month_id = $1
		ORDER BY bl.id`
	err := s.DB.Select(&budgetLines, query, monthID)
	if err != nil {
		return nil, fmt.Errorf("failed to get budget lines by month ID %d: %w", monthID, err)
	}
	if budgetLines == nil {
		return []BudgetLine{}, nil // Return empty slice instead of nil
	}
	return budgetLines, nil
}

// UpdateBudgetLine updates the label and expected amount of an existing BudgetLine.
func (s *sqlStore) UpdateBudgetLine(b *BudgetLine) error {
	_, err := s.DB.NamedExec(`
		UPDATE budget_lines
		SET label = :label, expected = :expected
		WHERE id = :id`, b)
	if err != nil {
		return fmt.Errorf("failed to update budget line with ID %d: %w", b.ID, err)
	}
	return nil
}

// UpdateActualLine updates the actual amount of an existing ActualLine.
// It also validates that the amount is non-negative and rounds it to 2 decimal places.
func (s *sqlStore) UpdateActualLine(a *ActualLine) error {
	if a.Actual < 0 {
		return fmt.Errorf("actual amount must be non-negative, got %f", a.Actual)
	}
	a.Actual = math.Round(a.Actual*100) / 100

	_, err := s.DB.NamedExec(`
		UPDATE actual_lines
		SET actual = :actual
		WHERE id = :id`, a)
	if err != nil {
		return fmt.Errorf("failed to update actual line with ID %d: %w", a.ID, err)
	}
	return nil
}

// GetActualLineByID retrieves an ActualLine by its ID.
func (s *sqlStore) GetActualLineByID(id int64) (*ActualLine, error) {
	var actualLine ActualLine
	err := s.DB.Get(&actualLine, "SELECT id, budget_line_id, actual FROM actual_lines WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("failed to get actual line with ID %d: %w", id, err)
	}
	return &actualLine, nil
}

// GetBudgetLineByID retrieves a BudgetLine by its ID.
func (s *sqlStore) GetBudgetLineByID(id int64) (*BudgetLine, error) {
	var budgetLine BudgetLine
	err := s.DB.Get(&budgetLine, "SELECT id, month_id, category_id, label, expected FROM budget_lines WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("failed to get budget line with ID %d: %w", id, err)
	}
	return &budgetLine, nil
}

// DeleteBudgetLine deletes a BudgetLine and its associated ActualLine.
func (s *sqlStore) DeleteBudgetLine(id int64) error {
	tx, err := s.DB.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete ActualLine first to maintain referential integrity if not using CASCADE DELETE
	_, err = tx.Exec("DELETE FROM actual_lines WHERE budget_line_id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete actual line for budget line ID %d: %w", id, err)
	}

	// Delete BudgetLine
	res, err := tx.Exec("DELETE FROM budget_lines WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete budget line with ID %d: %w", id, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected after deleting budget line ID %d: %w", id, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no budget line found with ID %d to delete", id)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction for deleting budget line ID %d: %w", id, err)
	}

	return nil
}
