package app

type DashboardPayload struct {
	MonthID         int               `json:"month_id"`
	Year            int               `json:"year"`
	Month           string            `json:"month"`
	TotalExpected   float64           `json:"total_expected"`
	TotalActual     float64           `json:"total_actual"`
	TotalDifference float64           `json:"total_difference"`
	CategorySummaries []CategorySummary `json:"category_summaries"`
}

type CategorySummary struct {
	CategoryID    int                `json:"category_id"`
	CategoryName  string             `json:"category_name"`
	CategoryColor string             `json:"category_color"`
	TotalExpected float64            `json:"total_expected"`
	TotalActual   float64            `json:"total_actual"`
	Difference    float64            `json:"difference"`
	BudgetLines   []BudgetLineDetail `json:"budget_lines"`
}

type BudgetLineDetail struct {
	BudgetLineID   int     `json:"budget_line_id"`
	Label          string  `json:"label"`
	ExpectedAmount float64 `json:"expected_amount"`
	ActualAmount   float64 `json:"actual_amount"`
	Difference     float64 `json:"difference"`
}
