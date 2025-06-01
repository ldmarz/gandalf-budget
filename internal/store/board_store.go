package store

import (
	"fmt" // For error formatting
	// Other necessary imports like "database/sql" if doing complex scanning,
	// but sqlx should handle it.
	"database/sql"
)

func (s *sqlStore) GetBoardData(monthID int) (*BoardDataPayload, error) {
	var monthDetails struct {
		Year      int  `db:"year"`
		Month     int  `db:"month"`
		Finalized bool `db:"finalized"`
	}
	monthQuery := `SELECT year, month, finalized FROM months WHERE id = ?;`
	err := s.DB.Get(&monthDetails, monthQuery, monthID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
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
	FROM budget_lines bl
	JOIN categories c ON bl.category_id = c.id
	LEFT JOIN actual_lines al ON bl.id = al.budget_line_id
	WHERE bl.month_id = ?
	ORDER BY c.name, bl.label;
	`
	err = s.DB.Select(&budgetLinesWithActuals, query, monthID)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error fetching budget lines with actuals for month %d: %w", monthID, err)
	}

	payload := &BoardDataPayload{
		MonthID:     int64(monthID),
		Year:        monthDetails.Year,
		MonthName:   monthName,
		BudgetLines: budgetLinesWithActuals,
		IsFinalized: monthDetails.Finalized,
	}

	return payload, nil
}
