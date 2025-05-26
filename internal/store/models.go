package store

// Category represents a budget category.
type Category struct {
	ID    int64  `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Color string `json:"color" db:"color"` // Tailwind color class
}

// Add other models here as needed, e.g.:
// type Month struct { ... }
// type BudgetLine struct { ... }
// type ActualLine struct { ... }
// type AnnualSnap struct { ... }
