package store

import (
	"errors"
)

type ReusableMockStore struct {
	MockGetAllCategories func() ([]Category, error)
	MockCreateCategory   func(category *Category) error
	MockGetCategoryByID  func(id int64) (*Category, error)
	MockUpdateCategory   func(category *Category) error
	MockDeleteCategory   func(id int64) error

	MockCreateBudgetLine        func(b *BudgetLine) (int64, error)
	MockGetBudgetLinesByMonthID func(monthID int) ([]BudgetLine, error)
	MockUpdateBudgetLine        func(b *BudgetLine) error
	MockDeleteBudgetLine        func(id int64) error
	MockUpdateActualLine        func(a *ActualLine) error
	MockGetActualLineByID       func(id int64) (*ActualLine, error)
	MockGetBudgetLineByID       func(id int64) (*BudgetLine, error)

	MockGetBoardData func(monthID int) (*BoardDataPayload, error)

	MockCanFinalizeMonth func(monthID int) (bool, string, error)
	MockFinalizeMonth    func(monthID int, snapJSON string) (int64, error)

	MockGetAnnualSnapshotsMetadataByYear func(year int) ([]AnnualSnapMeta, error)
	MockGetAnnualSnapshotJSONByID        func(snapID int64) (string, error)
}

func (m *ReusableMockStore) GetAllCategories() ([]Category, error) {
	if m.MockGetAllCategories != nil {
		return m.MockGetAllCategories()
	}
	return nil, errors.New("ReusableMockStore: MockGetAllCategories not implemented")
}

func (m *ReusableMockStore) CreateCategory(category *Category) error {
	if m.MockCreateCategory != nil {
		return m.MockCreateCategory(category)
	}
	return errors.New("ReusableMockStore: MockCreateCategory not implemented")
}

func (m *ReusableMockStore) GetCategoryByID(id int64) (*Category, error) {
	if m.MockGetCategoryByID != nil {
		return m.MockGetCategoryByID(id)
	}
	return nil, errors.New("ReusableMockStore: MockGetCategoryByID not implemented")
}

func (m *ReusableMockStore) UpdateCategory(category *Category) error {
	if m.MockUpdateCategory != nil {
		return m.MockUpdateCategory(category)
	}
	return errors.New("ReusableMockStore: MockUpdateCategory not implemented")
}

func (m *ReusableMockStore) DeleteCategory(id int64) error {
	if m.MockDeleteCategory != nil {
		return m.MockDeleteCategory(id)
	}
	return errors.New("ReusableMockStore: MockDeleteCategory not implemented")
}

func (m *ReusableMockStore) CreateBudgetLine(b *BudgetLine) (int64, error) {
	if m.MockCreateBudgetLine != nil {
		return m.MockCreateBudgetLine(b)
	}
	return 0, errors.New("ReusableMockStore: MockCreateBudgetLine not implemented")
}

func (m *ReusableMockStore) GetBudgetLinesByMonthID(monthID int) ([]BudgetLine, error) {
	if m.MockGetBudgetLinesByMonthID != nil {
		return m.MockGetBudgetLinesByMonthID(monthID)
	}
	return nil, errors.New("ReusableMockStore: MockGetBudgetLinesByMonthID not implemented")
}

func (m *ReusableMockStore) UpdateBudgetLine(b *BudgetLine) error {
	if m.MockUpdateBudgetLine != nil {
		return m.MockUpdateBudgetLine(b)
	}
	return errors.New("ReusableMockStore: MockUpdateBudgetLine not implemented")
}

func (m *ReusableMockStore) DeleteBudgetLine(id int64) error {
	if m.MockDeleteBudgetLine != nil {
		return m.MockDeleteBudgetLine(id)
	}
	return errors.New("ReusableMockStore: MockDeleteBudgetLine not implemented")
}

func (m *ReusableMockStore) UpdateActualLine(a *ActualLine) error {
	if m.MockUpdateActualLine != nil {
		return m.MockUpdateActualLine(a)
	}
	return errors.New("ReusableMockStore: MockUpdateActualLine not implemented")
}

func (m *ReusableMockStore) GetActualLineByID(id int64) (*ActualLine, error) {
	if m.MockGetActualLineByID != nil {
		return m.MockGetActualLineByID(id)
	}
	return nil, errors.New("ReusableMockStore: MockGetActualLineByID not implemented")
}

func (m *ReusableMockStore) GetBudgetLineByID(id int64) (*BudgetLine, error) {
	if m.MockGetBudgetLineByID != nil {
		return m.MockGetBudgetLineByID(id)
	}
	return nil, errors.New("ReusableMockStore: MockGetBudgetLineByID not implemented")
}

func (m *ReusableMockStore) GetBoardData(monthID int) (*BoardDataPayload, error) {
	if m.MockGetBoardData != nil {
		return m.MockGetBoardData(monthID)
	}
	return nil, errors.New("ReusableMockStore: MockGetBoardData not implemented")
}

func (m *ReusableMockStore) CanFinalizeMonth(monthID int) (bool, string, error) {
	if m.MockCanFinalizeMonth != nil {
		return m.MockCanFinalizeMonth(monthID)
	}
	return false, "", errors.New("ReusableMockStore: MockCanFinalizeMonth not implemented")
}

func (m *ReusableMockStore) FinalizeMonth(monthID int, snapJSON string) (int64, error) {
	if m.MockFinalizeMonth != nil {
		return m.MockFinalizeMonth(monthID, snapJSON)
	}
	return 0, errors.New("ReusableMockStore: MockFinalizeMonth not implemented")
}

func (m *ReusableMockStore) GetAnnualSnapshotsMetadataByYear(year int) ([]AnnualSnapMeta, error) {
	if m.MockGetAnnualSnapshotsMetadataByYear != nil {
		return m.MockGetAnnualSnapshotsMetadataByYear(year)
	}
	return nil, errors.New("ReusableMockStore: MockGetAnnualSnapshotsMetadataByYear not implemented")
}

func (m *ReusableMockStore) GetAnnualSnapshotJSONByID(snapID int64) (string, error) {
	if m.MockGetAnnualSnapshotJSONByID != nil {
		return m.MockGetAnnualSnapshotJSONByID(snapID)
	}
	return "", errors.New("ReusableMockStore: MockGetAnnualSnapshotJSONByID not implemented")
}
