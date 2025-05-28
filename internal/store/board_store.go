package store

import (
	"fmt" // For error formatting
	// Other necessary imports like "database/sql" if doing complex scanning,
	// but sqlx should handle it.
)

// GetBoardData retrieves all budget lines for a given monthID,
// joining with actual_lines to populate actual amounts.
func (s *sqlStore) GetBoardData(monthID int) ([]BudgetLine, error) {
	query := `
	SELECT
		bl.id,
		bl.month_id,
		bl.category_id,
		bl.label,
		bl.expected,
		al.id AS actual_id,
		al.actual AS actual_amount
	FROM budget_lines bl
	LEFT JOIN actual_lines al ON bl.id = al.budget_line_id
	WHERE bl.month_id = ?
	ORDER BY bl.id; -- Or some other meaningful order
	`
	var budgetLines []BudgetLine
	err := s.DB.Select(&budgetLines, query, monthID)
	if err != nil {
		return nil, fmt.Errorf("error fetching board data for month %d: %w", monthID, err)
	}
	return budgetLines, nil
}
