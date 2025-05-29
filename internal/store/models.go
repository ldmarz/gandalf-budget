package store

import "time"

type Category struct {
	ID    int64  `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Color string `json:"color" db:"color"`
}

type Month struct {
	ID        int64 `json:"id" db:"id"`
	Year      int   `json:"year" db:"year"`
	Month     int   `json:"month" db:"month"`
	Finalized bool  `json:"finalized" db:"finalized"`
}

type BudgetLine struct {
	ID           int     `json:"id" db:"id"`
	MonthID      int     `json:"month_id" db:"month_id"`
	CategoryID   int     `json:"category_id" db:"category_id"`
	Label        string  `json:"label" db:"label"`
	Expected     float64 `json:"expected" db:"expected"`
	ActualID     *int64  `json:"actual_id,omitempty" db:"actual_id"`
	ActualAmount *float64 `json:"actual_amount,omitempty" db:"actual_amount"`
}

type ActualLine struct {
	ID           int64   `json:"id" db:"id"`
	BudgetLineID int64   `json:"budget_line_id" db:"budget_line_id"`
	Actual       float64 `json:"actual" db:"actual"`
}

type AnnualSnap struct {
	ID        int64  `json:"id" db:"id"`
	MonthID   int64  `json:"month_id" db:"month_id"`
	SnapJSON  string `json:"snap_json" db:"snap_json"`
	CreatedAt string `json:"created_at" db:"created_at"`
}

type AnnualSnapMeta struct {
	ID            int64     `json:"id" db:"id"`
	MonthID       int64     `json:"month_id" db:"month_id"`
	Year          int       `json:"year" db:"year"`
	Month         string    `json:"month" db:"month_name"`
	SnapCreatedAt time.Time `json:"snap_created_at" db:"created_at"`
}

type BudgetLineWithActual struct {
	ID             int64   `json:"id" db:"id"`
	MonthID        int64   `json:"month_id" db:"month_id"`
	CategoryID     int64   `json:"category_id" db:"category_id"`
	CategoryName   string  `json:"category_name" db:"category_name"`
	CategoryColor  string  `json:"category_color" db:"category_color"`
	Label          string  `json:"label" db:"label"`
	ExpectedAmount float64 `json:"expected_amount" db:"expected_amount"`
	ActualAmount   float64 `json:"actual_amount" db:"actual_amount"`
}

type BoardDataPayload struct {
	MonthID     int64                  `json:"month_id"`
	Year        int                    `json:"year"`
	MonthName   string                 `json:"month_name"`
	BudgetLines []BudgetLineWithActual `json:"budget_lines"`
}
