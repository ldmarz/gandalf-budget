package store

// Category represents a budget category.
type Category struct {
	ID    int64  `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Color string `json:"color" db:"color"` // Tailwind color class
}

// Add other models here as needed, e.g.:
// type Month struct { ... }

// BudgetLine represents a line item in a monthly budget.
type BudgetLine struct {
	ID           int     `json:"id" db:"id"`
	MonthID      int     `json:"month_id" db:"month_id"`
	CategoryID   int     `json:"category_id" db:"category_id"`
	Label        string  `json:"label" db:"label"`
	Expected     float64 `json:"expected" db:"expected"`
	ActualID     *int64  `json:"actual_id,omitempty" db:"actual_id"`         // New field
	ActualAmount *float64 `json:"actual_amount,omitempty" db:"actual_amount"` // New field
}

// ActualLine represents the actual spending for a budget line.
type ActualLine struct {
	ID           int64   `json:"id" db:"id"` // Changed to int64 to match ActualID in BudgetLine
	BudgetLineID int64   `json:"budget_line_id" db:"budget_line_id"` // Changed to int64
	Actual       float64 `json:"actual" db:"actual"`
}

// type AnnualSnap struct { ... }
