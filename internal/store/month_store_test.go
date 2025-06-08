package store

import (
	"database/sql"
	"encoding/json"
	"testing"
	// "github.com/jmoiron/sqlx" // Implicitly used
)

func TestCanFinalizeMonth(t *testing.T) {
	db := newTestDB(t)
	s := NewSQLStore(db).(*sqlStore)

	catID := createTestCategory(t, db, "Utilities", "bg-yellow-500")

	// Month 1: All actuals non-zero
	month1ID := createTestMonth(t, db, 2023, 1, false)
	bl1m1 := createTestBudgetLine(t, db, month1ID, catID, "Electricity", 100.0)
	createTestActualLine(t, db, bl1m1, 95.0)
	bl2m1 := createTestBudgetLine(t, db, month1ID, catID, "Water", 50.0)
	createTestActualLine(t, db, bl2m1, 50.0)

	// Month 2: Some actuals are zero
	month2ID := createTestMonth(t, db, 2023, 2, false)
	bl1m2 := createTestBudgetLine(t, db, month2ID, catID, "Internet", 70.0)
	createTestActualLine(t, db, bl1m2, 70.0)
	bl2m2 := createTestBudgetLine(t, db, month2ID, catID, "Gas", 80.0)
	createTestActualLine(t, db, bl2m2, 0.0) // Zero actual

	// Month 3: No budget lines
	month3ID := createTestMonth(t, db, 2023, 3, false)

	// Month 4: Budget lines, but one has no actual_lines record yet
	// This case depends on how `CanFinalizeMonth` query is written.
	// The current query `INNER JOIN actual_lines` would not count lines without an actual_line record.
	// The PRD says "If any budget lines have an actual amount of 0".
	// This implies an actual_line record exists. If it doesn't, it's not "0", it's "missing".
	// For the current implementation of CanFinalizeMonth, this month should be finalizable
	// if the other line with an actual IS NOT 0.
	month4ID := createTestMonth(t, db, 2023, 4, false)
	bl1m4 := createTestBudgetLine(t, db, month4ID, catID, "Phone", 60.0)
	createTestActualLine(t, db, bl1m4, 55.0)                       // Non-zero actual
	createTestBudgetLine(t, db, month4ID, catID, "Cable TV", 90.0) // No actual_line record for this one

	tests := []struct {
		name       string
		monthID    int
		wantCan    bool
		wantReason string // Expected non-empty if wantCan is false
		wantErr    bool
	}{
		{
			name:       "All actuals non-zero",
			monthID:    int(month1ID),
			wantCan:    true,
			wantReason: "",
			wantErr:    false,
		},
		{
			name:       "Some actuals are zero",
			monthID:    int(month2ID),
			wantCan:    false,
			wantReason: "1 budget lines still have zero actuals.", // Specific to current implementation
			wantErr:    false,
		},
		{
			name:       "No budget lines for the month", // Should be finalizable as no lines have zero actuals
			monthID:    int(month3ID),
			wantCan:    true,
			wantReason: "",
			wantErr:    false,
		},
		{
			name: "Line missing actual_lines record, other is non-zero",
			// This depends on interpretation: "actual amount of 0" vs "no actual record"
			// Current query: `INNER JOIN actual_lines ... WHERE al.actual = 0`
			// This means lines without an actual_lines record are NOT considered as having "actual amount of 0".
			// So, if all *existing* actuals are non-zero, it should be finalizable.
			monthID:    int(month4ID),
			wantCan:    true,
			wantReason: "",
			wantErr:    false,
		},
		{
			name:       "Invalid month ID",
			monthID:    999,
			wantCan:    true, // No lines with zero actuals for a non-existent month
			wantReason: "",
			wantErr:    false, // The function itself doesn't error for non-existent month
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCan, gotReason, err := s.CanFinalizeMonth(tt.monthID)
			if (err != nil) != tt.wantErr {
				t.Errorf("CanFinalizeMonth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCan != tt.wantCan {
				t.Errorf("CanFinalizeMonth() gotCan = %v, wantCan %v", gotCan, tt.wantCan)
			}
			if !tt.wantCan && gotReason == "" {
				t.Errorf("CanFinalizeMonth() expected a reason when returning false, but got empty string")
			}
			if tt.wantReason != "" && gotReason != tt.wantReason {
				t.Errorf("CanFinalizeMonth() gotReason = %q, wantReason %q", gotReason, tt.wantReason)
			}
		})
	}
}

func TestFinalizeMonth(t *testing.T) {
	db := newTestDB(t)
	s := NewSQLStore(db).(*sqlStore)

	// Setup: A month with budget lines and actuals (all non-zero)
	catFoodID := createTestCategory(t, db, "Food", "bg-red-500")
	catRentID := createTestCategory(t, db, "Rent", "bg-blue-500")

	currentYear := 2023
	currentMonthVal := 11 // November, so next is December
	originalMonthID := createTestMonth(t, db, currentYear, currentMonthVal, false)

	bl1 := createTestBudgetLine(t, db, originalMonthID, catFoodID, "Groceries", 300.50)
	createTestActualLine(t, db, bl1, 280.75)
	bl2 := createTestBudgetLine(t, db, originalMonthID, catRentID, "Apartment Rent", 1200.00)
	createTestActualLine(t, db, bl2, 1200.00)
	_ = createTestBudgetLine(t, db, originalMonthID, catFoodID, "Restaurants", 150.00) // No actual line, for testing robustness of cloning (should still clone)

	snapData := map[string]interface{}{"total_expected": 1650.50, "total_actual": 1480.75}
	snapJSONBytes, _ := json.Marshal(snapData)
	snapJSON := string(snapJSONBytes)

	// Execute FinalizeMonth
	newMonthID, err := s.FinalizeMonth(int(originalMonthID), snapJSON)
	if err != nil {
		t.Fatalf("FinalizeMonth() failed: %v", err)
	}

	// Verify:
	// 1. An annual_snaps record is created
	var snap AnnualSnap
	err = db.Get(&snap, "SELECT month_id, snap_json FROM annual_snaps WHERE month_id = ?", originalMonthID)
	if err != nil {
		t.Errorf("Failed to get annual_snap for original month %d: %v", originalMonthID, err)
	}
	if snap.SnapJSON != snapJSON {
		t.Errorf("Annual_snap JSON mismatch: got %s, want %s", snap.SnapJSON, snapJSON)
	}
	// CreatedAt check is harder without time mocking, so we'll skip precise check for now.

	// 2. The original month is marked finalized = 1
	var originalMonthFinalized bool
	err = db.Get(&originalMonthFinalized, "SELECT finalized FROM months WHERE id = ?", originalMonthID)
	if err != nil {
		t.Errorf("Failed to get finalized status for original month %d: %v", originalMonthID, err)
	}
	if !originalMonthFinalized {
		t.Errorf("Original month %d was not marked as finalized", originalMonthID)
	}

	// 3. A new month record is created for the next calendar month
	if newMonthID == 0 {
		t.Fatal("FinalizeMonth returned newMonthID as 0")
	}
	var newMonth Month
	err = db.Get(&newMonth, "SELECT year, month, finalized FROM months WHERE id = ?", newMonthID)
	if err != nil {
		t.Fatalf("Failed to get new month record for ID %d: %v", newMonthID, err)
	}

	expectedNextYear, expectedNextMonthVal := currentYear, currentMonthVal+1
	if expectedNextMonthVal > 12 {
		expectedNextMonthVal = 1
		expectedNextYear++
	}

	if newMonth.Year != expectedNextYear || newMonth.Month != expectedNextMonthVal {
		t.Errorf("New month record has incorrect year/month: got %d-%d, want %d-%d", newMonth.Year, newMonth.Month, expectedNextYear, expectedNextMonthVal)
	}
	if newMonth.Finalized {
		t.Errorf("New month record was created as finalized, but should not be")
	}

	// 4. Budget lines are cloned to the new month
	var clonedLines []BudgetLine
	err = db.Select(&clonedLines, "SELECT category_id, label, expected FROM budget_lines WHERE month_id = ? ORDER BY label", newMonthID)
	if err != nil {
		t.Fatalf("Failed to get cloned budget lines for new month %d: %v", newMonthID, err)
	}

	expectedClonedLines := []struct {
		CategoryID int64
		Label      string
		Expected   float64
	}{
		{catRentID, "Apartment Rent", 1200.00},
		{catFoodID, "Groceries", 300.50},
		{catFoodID, "Restaurants", 150.00}, // This line had no actual, should still be cloned
	}
	if len(clonedLines) != len(expectedClonedLines) {
		t.Fatalf("Number of cloned lines mismatch: got %d, want %d. Got: %+v", len(clonedLines), len(expectedClonedLines), clonedLines)
	}
	// Simple check for label and expected, assuming order by label
	for i, el := range expectedClonedLines {
		cl := clonedLines[i]
		if cl.CategoryID != int(el.CategoryID) || cl.Label != el.Label || cl.Expected != el.Expected {
			t.Errorf("Cloned line mismatch at index %d: got CatID %d, Label %s, Exp %.2f; want CatID %d, Label %s, Exp %.2f",
				i, cl.CategoryID, cl.Label, cl.Expected, el.CategoryID, el.Label, el.Expected)
		}
	}

	// 5. New actual_lines are created for the cloned budget lines, with actual = 0
	var actualsForNewMonth []ActualLine
	queryActuals := `
		SELECT al.actual 
		FROM actual_lines al
		JOIN budget_lines bl ON al.budget_line_id = bl.id
		WHERE bl.month_id = ?;
	`
	err = db.Select(&actualsForNewMonth, queryActuals, newMonthID)
	if err != nil {
		t.Fatalf("Failed to get actual lines for new month %d: %v", newMonthID, err)
	}
	if len(actualsForNewMonth) != len(expectedClonedLines) { // Should be one actual_line per cloned budget_line
		t.Fatalf("Number of actual lines for new month mismatch: got %d, want %d", len(actualsForNewMonth), len(expectedClonedLines))
	}
	for i, al := range actualsForNewMonth {
		if al.Actual != 0 {
			t.Errorf("Actual line for cloned budget line at index %d has non-zero actual: got %.2f, want 0.00", i, al.Actual)
		}
	}

	// Test case: Finalizing a month that wraps around the year (December -> January)
	t.Run("FinalizeMonth_YearWrap", func(t *testing.T) {
		yearWrapDB := newTestDB(t) // Fresh DB for this sub-test to avoid ID conflicts
		yearWrapStore := NewSQLStore(yearWrapDB).(*sqlStore)

		catID := createTestCategory(t, yearWrapDB, "TestCat", "bg-gray-500")
		decMonthID := createTestMonth(t, yearWrapDB, 2023, 12, false)
		blDec := createTestBudgetLine(t, yearWrapDB, decMonthID, catID, "Year End", 100.0)
		createTestActualLine(t, yearWrapDB, blDec, 100.0)

		nextNewMonthID, err := yearWrapStore.FinalizeMonth(int(decMonthID), "{}")
		if err != nil {
			t.Fatalf("FinalizeMonth for year wrap failed: %v", err)
		}

		var nextMonthData Month
		err = yearWrapDB.Get(&nextMonthData, "SELECT year, month FROM months WHERE id = ?", nextNewMonthID)
		if err != nil {
			t.Fatalf("Failed to get next month data for year wrap: %v", err)
		}
		if nextMonthData.Year != 2024 || nextMonthData.Month != 1 {
			t.Errorf("Year wrap incorrect: expected 2024-01, got %d-%d", nextMonthData.Year, nextMonthData.Month)
		}
	})

	// Test transaction rollback (conceptual - hard to perfectly simulate DB failure mid-tx without mocks)
	// We can test by trying to finalize a month that doesn't exist, or violating a constraint if possible.
	// For now, we rely on the fact that if any step fails, the `tx.Rollback()` should be called.
	// A more direct test would involve a mock DB or specific error injection.
	t.Run("FinalizeMonth_ErrorRollbackConceptual", func(t *testing.T) {
		errorDB := newTestDB(t)
		errorStore := NewSQLStore(errorDB).(*sqlStore)

		// Create a month and a budget line
		monthToFailID := createTestMonth(t, errorDB, 2025, 1, false)
		catToFailID := createTestCategory(t, errorDB, "FailCat", "col")
		createTestBudgetLine(t, errorDB, monthToFailID, catToFailID, "Line1", 100)
		// No actual line for this one, so CanFinalizeMonth (if called) might prevent this.
		// However, FinalizeMonth doesn't call CanFinalizeMonth internally.

		// To simulate an error, let's try to make one of the inserts fail.
		// One way without mocking the DB is to violate a constraint NOT NULL or UNIQUE if we can control it.
		// For example, if annual_snaps had a UNIQUE constraint on month_id (it doesn't by default in 001_init.sql).
		// Or, if we could make `time.Now()` fail (not possible here).

		// For this conceptual test, we'll call FinalizeMonth with a monthID that doesn't exist
		// which will cause the `tx.Get(&currentMonth, ...)` to fail.
		// The goal is to ensure no partial data is written.
		invalidMonthID := 9999
		_, err := errorStore.FinalizeMonth(invalidMonthID, "{}")
		if err == nil {
			t.Errorf("Expected FinalizeMonth to fail for invalid month ID, but it didn't")
		}

		// Check that no new month was created (as an indicator of rollback)
		var count int
		errCheck := errorDB.Get(&count, "SELECT COUNT(*) FROM months WHERE year = ? AND month = ?", 2025, 2) // Example next month
		if errCheck != nil && errCheck != sql.ErrNoRows {                                                    // COUNT should return 0, not ErrNoRows here.
			t.Logf("Error checking for next month: %v", errCheck)
		}
		if count > 0 {
			t.Errorf("A new month was created despite an error in FinalizeMonth, indicating potential rollback failure. Count: %d", count)
		}

		// Check that no annual_snap was created for the invalid month
		errCheck = errorDB.Get(&count, "SELECT COUNT(*) FROM annual_snaps WHERE month_id = ?", invalidMonthID)
		if errCheck != nil && errCheck != sql.ErrNoRows {
			t.Logf("Error checking for annual_snaps: %v", errCheck)
		}
		if count > 0 {
			t.Errorf("An annual_snap was created despite an error in FinalizeMonth. Count: %d", count)
		}
	})
}

// TestFinalizeMonth_NoBudgetLines: Test finalizing a month that has no budget lines.
func TestFinalizeMonth_NoBudgetLines(t *testing.T) {
	db := newTestDB(t)
	s := NewSQLStore(db).(*sqlStore)

	currentYear, currentMonthVal := 2024, 1
	originalMonthID := createTestMonth(t, db, currentYear, currentMonthVal, false)
	// No budget lines created for this month.

	snapJSON := `{"message": "No budget lines to snapshot"}`

	newMonthID, err := s.FinalizeMonth(int(originalMonthID), snapJSON)
	if err != nil {
		t.Fatalf("FinalizeMonth() for month with no budget lines failed: %v", err)
	}

	// Verify:
	// 1. Annual snap created
	var snapCount int
	err = db.Get(&snapCount, "SELECT COUNT(*) FROM annual_snaps WHERE month_id = ?", originalMonthID)
	if err != nil || snapCount != 1 {
		t.Errorf("Expected 1 annual_snap, got count %d, err: %v", snapCount, err)
	}

	// 2. Original month finalized
	var originalMonthFinalized bool
	err = db.Get(&originalMonthFinalized, "SELECT finalized FROM months WHERE id = ?", originalMonthID)
	if err != nil || !originalMonthFinalized {
		t.Errorf("Original month not finalized. Got finalized=%v, err: %v", originalMonthFinalized, err)
	}

	// 3. New month created
	var newMonth Month
	err = db.Get(&newMonth, "SELECT year, month FROM months WHERE id = ?", newMonthID)
	if err != nil {
		t.Fatalf("Failed to get new month record: %v", err)
	}
	expectedNextYear, expectedNextMonthVal := currentYear, currentMonthVal+1
	if newMonth.Year != expectedNextYear || newMonth.Month != expectedNextMonthVal {
		t.Errorf("New month has incorrect year/month: got %d-%d, want %d-%d", newMonth.Year, newMonth.Month, expectedNextYear, expectedNextMonthVal)
	}

	// 4. No budget lines cloned (since there were none)
	var clonedLinesCount int
	err = db.Get(&clonedLinesCount, "SELECT COUNT(*) FROM budget_lines WHERE month_id = ?", newMonthID)
	if err != nil || clonedLinesCount != 0 {
		t.Errorf("Expected 0 cloned budget lines, got %d, err: %v", clonedLinesCount, err)
	}

	// 5. No actual lines created for new month (since no budget lines)
	var actualsCount int
	queryActuals := `SELECT COUNT(al.id) FROM actual_lines al JOIN budget_lines bl ON al.budget_line_id = bl.id WHERE bl.month_id = ?`
	err = db.Get(&actualsCount, queryActuals, newMonthID)
	// This query will naturally return 0 if no budget_lines exist for newMonthID.
	// If there were budget lines but no actuals, that would be an issue. Here, 0 budget lines means 0 actuals.
	if err != nil && err != sql.ErrNoRows { // ErrNoRows could happen if Get is used on COUNT that returns nothing (unlikely for COUNT)
		t.Errorf("Error counting actual lines for new month: %v", err)
	}
	if actualsCount != 0 {
		t.Errorf("Expected 0 actual lines for new month, got %d", actualsCount)
	}
}
