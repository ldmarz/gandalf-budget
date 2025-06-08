package store

import (
	"reflect"
	"testing"
)

func TestGetBoardData(t *testing.T) {
	db := newTestDB(t)
	s := NewSQLStore(db).(*sqlStore)

	cat1ID := createTestCategory(t, db, "Food", "bg-red-500")
	cat2ID := createTestCategory(t, db, "Travel", "bg-blue-500")

	month1ID := createTestMonth(t, db, 2023, 1, false)
	month2ID := createTestMonth(t, db, 2023, 2, false)
	month3ID := createTestMonth(t, db, 2023, 3, false)

	bl1m1 := createTestBudgetLine(t, db, month1ID, cat1ID, "Groceries", 500.0)
	bl2m1 := createTestBudgetLine(t, db, month1ID, cat2ID, "Gas", 100.0)
	createTestActualLine(t, db, bl1m1, 480.0)
	createTestActualLine(t, db, bl2m1, 110.0)

	bl1m3 := createTestBudgetLine(t, db, month3ID, cat1ID, "Restaurant", 150.0)
	createTestActualLine(t, db, bl1m3, 0)

	bl2m3 := createTestBudgetLine(t, db, month3ID, cat2ID, "Bus Pass", 50.0)

	tests := []struct {
		name          string
		monthID       int
		expectedLines []BudgetLineWithActual
		expectError   bool
	}{
		{
			name:    "Month with budget lines and actual lines",
			monthID: int(month1ID),
			expectedLines: []BudgetLineWithActual{
				{ID: bl1m1, MonthID: month1ID, CategoryID: cat1ID, CategoryName: "Food", CategoryColor: "bg-red-500", Label: "Groceries", ExpectedAmount: 500.0, ActualAmount: 480.0},
				{ID: bl2m1, MonthID: month1ID, CategoryID: cat2ID, CategoryName: "Travel", CategoryColor: "bg-blue-500", Label: "Gas", ExpectedAmount: 100.0, ActualAmount: 110.0},
			},
			expectError: false,
		},
		{
			name:          "Month with no budget lines",
			monthID:       int(month2ID),
			expectedLines: []BudgetLineWithActual{},
			expectError:   true,
		},
		{
			name:    "Month with lines but varied actuals",
			monthID: int(month3ID),
			expectedLines: []BudgetLineWithActual{
				{ID: bl1m3, MonthID: month3ID, CategoryID: cat1ID, CategoryName: "Food", CategoryColor: "bg-red-500", Label: "Restaurant", ExpectedAmount: 150.0, ActualAmount: 0},
				{ID: bl2m3, MonthID: month3ID, CategoryID: cat2ID, CategoryName: "Travel", CategoryColor: "bg-blue-500", Label: "Bus Pass", ExpectedAmount: 50.0, ActualAmount: 0},
			},
			expectError: false,
		},
		{
			name:          "Invalid monthID (non-existent)",
			monthID:       999,
			expectedLines: []BudgetLineWithActual{},
			expectError:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			payload, err := s.GetBoardData(tc.monthID)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error, but got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("Did not expect an error, but got: %v", err)
			}

			if len(payload.BudgetLines) == 0 && len(tc.expectedLines) == 0 {
			} else if !reflect.DeepEqual(payload.BudgetLines, tc.expectedLines) {
				t.Errorf("Expected lines %+v, but got %+v", tc.expectedLines, payload.BudgetLines)
				for i := 0; i < len(payload.BudgetLines) || i < len(tc.expectedLines); i++ {
					if i < len(payload.BudgetLines) && i < len(tc.expectedLines) {
						if !reflect.DeepEqual(payload.BudgetLines[i], tc.expectedLines[i]) {
							t.Logf("Difference at index %d: Expected %+v, Got %+v", i, tc.expectedLines[i], payload.BudgetLines[i])
						}
					} else if i < len(payload.BudgetLines) {
						t.Logf("Extra line at index %d (got): %+v", i, payload.BudgetLines[i])
					} else {
						t.Logf("Missing line at index %d (expected): %+v", i, tc.expectedLines[i])
					}
				}
			}
		})
	}
}

func ptrToInt64(v int64) *int64 {
	return &v
}

func ptrToFloat64(v float64) *float64 {
	return &v
}
