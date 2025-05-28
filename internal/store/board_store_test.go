package store

import (
	"reflect"
	"testing"
	// "github.com/jmoiron/sqlx" // Implicitly used via sqlStore
)

func TestGetBoardData(t *testing.T) {
	db := newTestDB(t) // Uses the helper from store_test_setup_test.go
	s := NewSQLStore(db).(*sqlStore) // Assuming NewSQLStore returns the interface, cast to access DB directly or use methods

	// Setup initial data
	cat1ID := createTestCategory(t, db, "Food", "bg-red-500")
	cat2ID := createTestCategory(t, db, "Travel", "bg-blue-500")

	month1ID := createTestMonth(t, db, 2023, 1, false)
	month2ID := createTestMonth(t, db, 2023, 2, false) // Month with no budget lines
	month3ID := createTestMonth(t, db, 2023, 3, false) // Month with lines but no actuals

	// Budget lines for month 1
	bl1m1 := createTestBudgetLine(t, db, month1ID, cat1ID, "Groceries", 500.0)
	bl2m1 := createTestBudgetLine(t, db, month1ID, cat2ID, "Gas", 100.0)
	createTestActualLine(t, db, bl1m1, 480.0)
	createTestActualLine(t, db, bl2m1, 110.0)

	// Budget lines for month 3 (no actuals, or actuals are 0)
	bl1m3 := createTestBudgetLine(t, db, month3ID, cat1ID, "Restaurant", 150.0)
	createTestActualLine(t, db, bl1m3, 0) // Actual line with 0 value
	
	bl2m3 := createTestBudgetLine(t, db, month3ID, cat2ID, "Bus Pass", 50.0)
	// No actual line for bl2m3, so actual_id and actual_amount should be nil/zero depending on join

	// Define test cases
	tests := []struct {
		name          string
		monthID       int
		expectedLines []BudgetLine
		expectError   bool
	}{
		{
			name:    "Month with budget lines and actual lines",
			monthID: int(month1ID),
			expectedLines: []BudgetLine{
				{ID: int(bl1m1), MonthID: int(month1ID), CategoryID: int(cat1ID), Label: "Groceries", Expected: 500.0, ActualID: ptrToInt64(1), ActualAmount: ptrToFloat64(480.0)}, // Assuming actual_id is 1
				{ID: int(bl2m1), MonthID: int(month1ID), CategoryID: int(cat2ID), Label: "Gas", Expected: 100.0, ActualID: ptrToInt64(2), ActualAmount: ptrToFloat64(110.0)}, // Assuming actual_id is 2
			},
			expectError: false,
		},
		{
			name:          "Month with no budget lines",
			monthID:       int(month2ID),
			expectedLines: []BudgetLine{}, // Expect an empty slice, not nil
			expectError:   false,
		},
		{
			name:    "Month with lines but varied actuals",
			monthID: int(month3ID),
			expectedLines: []BudgetLine{
				// For bl1m3, actual_lines.id will be 3, actual_lines.actual is 0
				{ID: int(bl1m3), MonthID: int(month3ID), CategoryID: int(cat1ID), Label: "Restaurant", Expected: 150.0, ActualID: ptrToInt64(3), ActualAmount: ptrToFloat64(0.0)},
				// For bl2m3, there's no corresponding actual_lines record, so ActualID and ActualAmount should be nil (or zero if coalesce is used in query)
				// The current GetBoardData query uses LEFT JOIN, so ActualID and ActualAmount will be nil for bl2m3.
				// The BudgetLine struct uses pointers for ActualID and ActualAmount to handle this.
				{ID: int(bl2m3), MonthID: int(month3ID), CategoryID: int(cat2ID), Label: "Bus Pass", Expected: 50.0, ActualID: nil, ActualAmount: nil},
			},
			expectError: false,
		},
		{
			name:          "Invalid monthID (non-existent)",
			monthID:       999,
			expectedLines: []BudgetLine{}, // Expect an empty slice for non-existent month
			expectError:   false,          // The function itself doesn't error, just returns no lines
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lines, err := s.GetBoardData(tc.monthID)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error, but got nil")
				}
				return // Don't check lines if error was expected
			}
			if err != nil {
				t.Fatalf("Did not expect an error, but got: %v", err)
			}

			// Normalize actual_id if it's 0 due to no record vs. actual record with ID 0 (not possible with autoincrement)
			// For this test, we assume actual_id starts from 1.
			// The important part is comparing the structure and relevant values.
			// reflect.DeepEqual is strict with nil vs empty.
			if len(lines) == 0 && len(tc.expectedLines) == 0 {
				// Both are empty, so they are equal in this context.
			} else if !reflect.DeepEqual(lines, tc.expectedLines) {
				t.Errorf("Expected lines %+v, but got %+v", tc.expectedLines, lines)
				// For detailed comparison:
				for i := 0; i < len(lines) || i < len(tc.expectedLines); i++ {
					if i < len(lines) && i < len(tc.expectedLines) {
						if !reflect.DeepEqual(lines[i], tc.expectedLines[i]) {
							t.Logf("Difference at index %d: Expected %+v, Got %+v", i, tc.expectedLines[i], lines[i])
							// Further breakdown by field if needed
							if !reflect.DeepEqual(lines[i].ActualID, tc.expectedLines[i].ActualID) {
								t.Logf("  ActualID diff: Expected %v, Got %v", tc.expectedLines[i].ActualID, lines[i].ActualID)
							}
							if !reflect.DeepEqual(lines[i].ActualAmount, tc.expectedLines[i].ActualAmount) {
								t.Logf("  ActualAmount diff: Expected %v, Got %v", tc.expectedLines[i].ActualAmount, lines[i].ActualAmount)
							}
						}
					} else if i < len(lines) {
						t.Logf("Extra line at index %d (got): %+v", i, lines[i])
					} else {
						t.Logf("Missing line at index %d (expected): %+v", i, tc.expectedLines[i])
					}
				}
			}
		})
	}
}

// Helper function to get a pointer to an int64 value.
func ptrToInt64(v int64) *int64 {
	return &v
}

// Helper function to get a pointer to a float64 value.
func ptrToFloat64(v float64) *float64 {
	return &v
}
