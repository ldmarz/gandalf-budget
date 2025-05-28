package store

import (
	"fmt" // For error formatting
	// Other necessary imports like "database/sql" if doing complex scanning,
	// but sqlx should handle it.
	"database/sql" // For sql.ErrNoRows
)

// GetBoardData retrieves comprehensive data for the monthly budget board.
// It includes month details and all budget lines with their actuals and category info.
func (s *sqlStore) GetBoardData(monthID int) (*BoardDataPayload, error) {
	// 1. Fetch month details
	var monthDetails struct {
		Year  int `db:"year"`
		Month int `db:"month"`
	}
	monthQuery := `SELECT year, month FROM months WHERE id = ?;`
	err := s.DB.Get(&monthDetails, monthQuery, monthID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows // Month not found
		}
		return nil, fmt.Errorf("error fetching month details for month ID %d: %w", monthID, err)
	}

	monthName := ""
	switch monthDetails.Month {
	case 1: monthName = "January"
	case 2: monthName = "February"
	case 3: monthName = "March"
	case 4: monthName = "April"
	case 5: monthName = "May"
	case 6: monthName = "June"
	case 7: monthName = "July"
	case 8: monthName = "August"
	case 9: monthName = "September"
	case 10: monthName = "October"
	case 11: monthName = "November"
	case 12: monthName = "December"
	default: monthName = "Unknown"
	}

	// 2. Fetch budget lines with actuals and category info
	var budgetLinesWithActuals []BudgetLineWithActual
	query := `
	SELECT
		bl.id,
		bl.month_id,
		bl.category_id,
		c.name AS category_name,
		c.color AS category_color,
		bl.label,
		bl.expected AS expected_amount,
		COALESCE(al.actual, 0) AS actual_amount 
		-- Using COALESCE to ensure actual_amount is 0 if no matching actual_line,
		-- assuming actual_lines might be missing for newly created budget_lines
		-- or if ActualID on BudgetLine is nullable and not set.
		-- The PRD implies actual_lines are auto-created, so LEFT JOIN might be sufficient
		-- if actual_lines.actual can be NULL. If actual_lines.actual is NOT NULL and defaults to 0,
		-- then COALESCE is still safe.
	FROM budget_lines bl
	JOIN categories c ON bl.category_id = c.id
	LEFT JOIN actual_lines al ON bl.id = al.budget_line_id
	WHERE bl.month_id = ?
	ORDER BY c.name, bl.label; -- Meaningful order
	`
	err = s.DB.Select(&budgetLinesWithActuals, query, monthID)
	if err != nil && err != sql.ErrNoRows { // sql.ErrNoRows is ok here, means no budget lines
		return nil, fmt.Errorf("error fetching budget lines with actuals for month %d: %w", monthID, err)
	}
	// If err is sql.ErrNoRows, budgetLinesWithActuals will be an empty slice, which is correct.

	payload := &BoardDataPayload{
		MonthID:     int64(monthID), // Convert int to int64 if your struct field is int64
		Year:        monthDetails.Year,
		MonthName:   monthName,
		BudgetLines: budgetLinesWithActuals,
	}

	return payload, nil
}
