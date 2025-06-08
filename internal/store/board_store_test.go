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
		name            string
		monthID         int
		expectedPayload *BoardDataPayload
		expectError     bool
	}{
		{
			name:    "Month with budget lines and actual lines",
			monthID: int(month1ID),
			expectedPayload: &BoardDataPayload{
				MonthID:   month1ID,
				Year:      2023,
				MonthName: "January",
				BudgetLines: []BudgetLineWithActual{
					{ID: bl1m1, MonthID: month1ID, CategoryID: cat1ID, CategoryName: "Food", CategoryColor: "bg-red-500", Label: "Groceries", ExpectedAmount: 500.0, ActualAmount: 480.0},
					{ID: bl2m1, MonthID: month1ID, CategoryID: cat2ID, CategoryName: "Travel", CategoryColor: "bg-blue-500", Label: "Gas", ExpectedAmount: 100.0, ActualAmount: 110.0},
				},
				IsFinalized: false,
			},
			expectError: false,
		},
		{
			name:    "Month with no budget lines",
			monthID: int(month2ID),
			expectedPayload: &BoardDataPayload{
				MonthID:     month2ID,
				Year:        2023,
				MonthName:   "February",
				BudgetLines: []BudgetLineWithActual{},
				IsFinalized: false,
			},
			expectError: false,
		},
		{
			name:    "Month with lines but varied actuals",
			monthID: int(month3ID),
			expectedPayload: &BoardDataPayload{
				MonthID:   month3ID,
				Year:      2023,
				MonthName: "March",
				BudgetLines: []BudgetLineWithActual{
					{ID: bl1m3, MonthID: month3ID, CategoryID: cat1ID, CategoryName: "Food", CategoryColor: "bg-red-500", Label: "Restaurant", ExpectedAmount: 150.0, ActualAmount: 0.0},
					{ID: bl2m3, MonthID: month3ID, CategoryID: cat2ID, CategoryName: "Travel", CategoryColor: "bg-blue-500", Label: "Bus Pass", ExpectedAmount: 50.0, ActualAmount: 0.0},
				},
				IsFinalized: false,
			},
			expectError: false,
		},
		{
			name:        "Invalid monthID (non-existent)",
			monthID:     999,
			expectError: true,
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

			if !reflect.DeepEqual(payload, tc.expectedPayload) {
				t.Errorf("Expected payload %+v, got %+v", tc.expectedPayload, payload)
			}
		})
	}
}
