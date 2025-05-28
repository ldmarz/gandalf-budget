// Package app contains the core application logic.
package app

// DashboardPayload represents the data structure for the dashboard.
type DashboardPayload struct {
	MonthID         int               `json:"month_id"`         // The ID of the month.
	Year            int               `json:"year"`             // The year.
	Month           string            `json:"month"`            // The name of the month (e.g., "January").
	TotalExpected   float64           `json:"total_expected"`   // The total expected amount for the month.
	TotalActual     float64           `json:"total_actual"`     // The total actual amount for the month.
	TotalDifference float64           `json:"total_difference"` // The difference between total expected and total actual.
	CategorySummaries []CategorySummary `json:"category_summaries"` // A slice of category summaries.
}

// CategorySummary represents the summary of a specific category.
type CategorySummary struct {
	CategoryID    int                `json:"category_id"`    // The ID of the category.
	CategoryName  string             `json:"category_name"`  // The name of the category.
	CategoryColor string             `json:"category_color"` // The color associated with the category.
	TotalExpected float64            `json:"total_expected"` // The total expected amount for this category.
	TotalActual   float64            `json:"total_actual"`   // The total actual amount for this category.
	Difference    float64            `json:"difference"`     // The difference between expected and actual for this category.
	BudgetLines   []BudgetLineDetail `json:"budget_lines"`   // A slice of budget line details for this category.
}

// BudgetLineDetail represents the details of a specific budget line item.
type BudgetLineDetail struct {
	BudgetLineID   int     `json:"budget_line_id"`   // The ID of the budget line.
	Label          string  `json:"label"`            // The label or name of the budget line item.
	ExpectedAmount float64 `json:"expected_amount"`  // The expected amount for this budget line item.
	ActualAmount   float64 `json:"actual_amount"`    // The actual amount for this budget line item.
	Difference     float64 `json:"difference"`       // The difference between expected and actual for this budget line item.
}
