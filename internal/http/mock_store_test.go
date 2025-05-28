package http

import (
	"gandalf-budget/internal/store" // Import the actual store package
)

// MockStore is a mock implementation of the store.Store interface for testing HTTP handlers.
type MockStore struct {
	// --- Category methods ---
	GetAllCategoriesFunc    func() ([]store.Category, error)
	CreateCategoryFunc      func(category *store.Category) error
	GetCategoryByIDFunc     func(id int64) (*store.Category, error)
	UpdateCategoryFunc      func(category *store.Category) error
	DeleteCategoryFunc      func(id int64) error

	// --- BudgetLine and ActualLine methods ---
	CreateBudgetLineFunc        func(b *store.BudgetLine) (int64, error)
	GetBudgetLinesByMonthIDFunc func(monthID int) ([]store.BudgetLine, error)
	UpdateBudgetLineFunc        func(b *store.BudgetLine) error
	DeleteBudgetLineFunc        func(id int64) error
	UpdateActualLineFunc        func(a *store.ActualLine) error
	GetActualLineByIDFunc       func(id int64) (*store.ActualLine, error)
	GetBudgetLineByIDFunc       func(id int64) (*store.BudgetLine, error)

	// --- Board data methods ---
	GetBoardDataFunc func(monthID int) ([]store.BudgetLine, error)

	// --- Month finalization methods ---
	CanFinalizeMonthFunc func(monthID int) (bool, string, error)
	FinalizeMonthFunc    func(monthID int, snapJSON string) (int64, error)
}

// --- Category methods implementation ---
func (m *MockStore) GetAllCategories() ([]store.Category, error) {
	if m.GetAllCategoriesFunc != nil {
		return m.GetAllCategoriesFunc()
	}
	return nil, nil // Default behavior
}
func (m *MockStore) CreateCategory(category *store.Category) error {
	if m.CreateCategoryFunc != nil {
		return m.CreateCategoryFunc(category)
	}
	return nil
}
func (m *MockStore) GetCategoryByID(id int64) (*store.Category, error) {
	if m.GetCategoryByIDFunc != nil {
		return m.GetCategoryByIDFunc(id)
	}
	return nil, nil
}
func (m *MockStore) UpdateCategory(category *store.Category) error {
	if m.UpdateCategoryFunc != nil {
		return m.UpdateCategoryFunc(category)
	}
	return nil
}
func (m *MockStore) DeleteCategory(id int64) error {
	if m.DeleteCategoryFunc != nil {
		return m.DeleteCategoryFunc(id)
	}
	return nil
}

// --- BudgetLine and ActualLine methods implementation ---
func (m *MockStore) CreateBudgetLine(b *store.BudgetLine) (int64, error) {
	if m.CreateBudgetLineFunc != nil {
		return m.CreateBudgetLineFunc(b)
	}
	return 0, nil
}
func (m *MockStore) GetBudgetLinesByMonthID(monthID int) ([]store.BudgetLine, error) {
	if m.GetBudgetLinesByMonthIDFunc != nil {
		return m.GetBudgetLinesByMonthIDFunc(monthID)
	}
	return nil, nil
}
func (m *MockStore) UpdateBudgetLine(b *store.BudgetLine) error {
	if m.UpdateBudgetLineFunc != nil {
		return m.UpdateBudgetLineFunc(b)
	}
	return nil
}
func (m *MockStore) DeleteBudgetLine(id int64) error {
	if m.DeleteBudgetLineFunc != nil {
		return m.DeleteBudgetLineFunc(id)
	}
	return nil
}
func (m *MockStore) UpdateActualLine(a *store.ActualLine) error {
	if m.UpdateActualLineFunc != nil {
		return m.UpdateActualLineFunc(a)
	}
	return nil
}
func (m *MockStore) GetActualLineByID(id int64) (*store.ActualLine, error) {
	if m.GetActualLineByIDFunc != nil {
		return m.GetActualLineByIDFunc(id)
	}
	return nil, nil
}
func (m *MockStore) GetBudgetLineByID(id int64) (*store.BudgetLine, error) {
	if m.GetBudgetLineByIDFunc != nil {
		return m.GetBudgetLineByIDFunc(id)
	}
	return nil, nil
}

// --- Board data methods implementation ---
func (m *MockStore) GetBoardData(monthID int) ([]store.BudgetLine, error) {
	if m.GetBoardDataFunc != nil {
		return m.GetBoardDataFunc(monthID)
	}
	return nil, nil // Default behavior
}

// --- Month finalization methods implementation ---
func (m *MockStore) CanFinalizeMonth(monthID int) (bool, string, error) {
	if m.CanFinalizeMonthFunc != nil {
		return m.CanFinalizeMonthFunc(monthID)
	}
	return false, "", nil // Default behavior
}
func (m *MockStore) FinalizeMonth(monthID int, snapJSON string) (int64, error) {
	if m.FinalizeMonthFunc != nil {
		return m.FinalizeMonthFunc(monthID, snapJSON)
	}
	return 0, nil // Default behavior
}
