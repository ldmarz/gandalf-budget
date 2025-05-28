package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gandalf-budget/internal/store"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// MockStore is a mock implementation of the store.Store interface for testing.
type MockStore struct {
	// Category methods
	MockGetAllCategories func() ([]store.Category, error)
	MockCreateCategory   func(category *store.Category) error
	MockGetCategoryByID  func(id int64) (*store.Category, error)
	MockUpdateCategory   func(category *store.Category) error
	MockDeleteCategory   func(id int64) error

	// BudgetLine and ActualLine methods
	MockCreateBudgetLine        func(b *store.BudgetLine) (int64, error)
	MockGetBudgetLinesByMonthID func(monthID int) ([]store.BudgetLine, error)
	MockUpdateBudgetLine        func(b *store.BudgetLine) error
	MockDeleteBudgetLine        func(id int64) error
	MockUpdateActualLine        func(a *store.ActualLine) error
	MockGetActualLineByID       func(id int64) (*store.ActualLine, error)
	MockGetBudgetLineByID       func(id int64) (*store.BudgetLine, error)
}

// Implement store.Store interface for MockStore
func (m *MockStore) GetAllCategories() ([]store.Category, error) {
	if m.MockGetAllCategories != nil {
		return m.MockGetAllCategories()
	}
	return nil, fmt.Errorf("MockGetAllCategories not implemented")
}

func (m *MockStore) CreateCategory(category *store.Category) error {
	if m.MockCreateCategory != nil {
		return m.MockCreateCategory(category)
	}
	return fmt.Errorf("MockCreateCategory not implemented")
}

func (m *MockStore) GetCategoryByID(id int64) (*store.Category, error) {
	if m.MockGetCategoryByID != nil {
		return m.MockGetCategoryByID(id)
	}
	return nil, fmt.Errorf("MockGetCategoryByID not implemented")
}

func (m *MockStore) UpdateCategory(category *store.Category) error {
	if m.MockUpdateCategory != nil {
		return m.MockUpdateCategory(category)
	}
	return fmt.Errorf("MockUpdateCategory not implemented")
}

func (m *MockStore) DeleteCategory(id int64) error {
	if m.MockDeleteCategory != nil {
		return m.MockDeleteCategory(id)
	}
	return fmt.Errorf("MockDeleteCategory not implemented")
}

func (m *MockStore) CreateBudgetLine(b *store.BudgetLine) (int64, error) {
	if m.MockCreateBudgetLine != nil {
		return m.MockCreateBudgetLine(b)
	}
	return 0, fmt.Errorf("MockCreateBudgetLine not implemented")
}

func (m *MockStore) GetBudgetLinesByMonthID(monthID int) ([]store.BudgetLine, error) {
	if m.MockGetBudgetLinesByMonthID != nil {
		return m.MockGetBudgetLinesByMonthID(monthID)
	}
	return nil, fmt.Errorf("MockGetBudgetLinesByMonthID not implemented")
}

func (m *MockStore) UpdateBudgetLine(b *store.BudgetLine) error {
	if m.MockUpdateBudgetLine != nil {
		return m.MockUpdateBudgetLine(b)
	}
	return fmt.Errorf("MockUpdateBudgetLine not implemented")
}

func (m *MockStore) DeleteBudgetLine(id int64) error {
	if m.MockDeleteBudgetLine != nil {
		return m.MockDeleteBudgetLine(id)
	}
	return fmt.Errorf("MockDeleteBudgetLine not implemented")
}

func (m *MockStore) UpdateActualLine(a *store.ActualLine) error {
	if m.MockUpdateActualLine != nil {
		return m.MockUpdateActualLine(a)
	}
	return fmt.Errorf("MockUpdateActualLine not implemented")
}

func (m *MockStore) GetActualLineByID(id int64) (*store.ActualLine, error) {
	if m.MockGetActualLineByID != nil {
		return m.MockGetActualLineByID(id)
	}
	return nil, fmt.Errorf("MockGetActualLineByID not implemented")
}

func (m *MockStore) GetBudgetLineByID(id int64) (*store.BudgetLine, error) {
	if m.MockGetBudgetLineByID != nil {
		return m.MockGetBudgetLineByID(id)
	}
	return nil, fmt.Errorf("MockGetBudgetLineByID not implemented")
}

// TestCreateBudgetLineHandler tests the CreateBudgetLineHandler.
func TestCreateBudgetLineHandler(t *testing.T) {
	mockStore := &MockStore{}
	handler := CreateBudgetLineHandler(mockStore) // Assuming this is the correct handler name

	t.Run("successful creation", func(t *testing.T) {
		expectedID := int64(1)
		inputBudgetLine := store.BudgetLine{
			MonthID:    1,
			CategoryID: 1,
			Label:      "Test Groceries",
			Expected:   150.00,
		}
		createdBudgetLine := store.BudgetLine{
			ID:         int(expectedID), // Note: handler sets this
			MonthID:    inputBudgetLine.MonthID,
			CategoryID: inputBudgetLine.CategoryID,
			Label:      inputBudgetLine.Label,
			Expected:   inputBudgetLine.Expected,
		}

		mockStore.MockCreateBudgetLine = func(bl *store.BudgetLine) (int64, error) {
			if bl.Label != inputBudgetLine.Label || bl.Expected != inputBudgetLine.Expected || bl.MonthID != inputBudgetLine.MonthID || bl.CategoryID != inputBudgetLine.CategoryID {
				t.Errorf("MockCreateBudgetLine called with unexpected data: got %+v, want label %s", bl, inputBudgetLine.Label)
			}
			// The handler will update the ID of the passed-in budget line, so we don't check it here.
			return expectedID, nil
		}

		budgetLineJSON, _ := json.Marshal(inputBudgetLine)
		req := httptest.NewRequest("POST", "/api/v1/budget-lines", bytes.NewReader(budgetLineJSON))
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d. Body: %s", http.StatusCreated, rr.Code, rr.Body.String())
		}

		var respBody store.BudgetLine
		if err := json.Unmarshal(rr.Body.Bytes(), &respBody); err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		if respBody.ID != int(expectedID) {
			t.Errorf("expected budget line ID %d, got %d", expectedID, respBody.ID)
		}
		if respBody.Label != createdBudgetLine.Label {
			t.Errorf("expected label %s, got %s", createdBudgetLine.Label, respBody.Label)
		}
		// Add more assertions for other fields if necessary
	})

	t.Run("invalid request body - bad JSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/budget-lines", strings.NewReader("{not_json"))
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("missing required fields", func(t *testing.T) {
		budgetLineJSON := `{"label":"Test"}` // Missing month_id, category_id
		req := httptest.NewRequest("POST", "/api/v1/budget-lines", strings.NewReader(budgetLineJSON))
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("store error on creation", func(t *testing.T) {
		mockStore.MockCreateBudgetLine = func(bl *store.BudgetLine) (int64, error) {
			return 0, fmt.Errorf("database error")
		}
		budgetLineJSON := `{"label":"Test","expected":100,"month_id":1,"category_id":1}`
		req := httptest.NewRequest("POST", "/api/v1/budget-lines", strings.NewReader(budgetLineJSON))
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
		}
	})
}

// Placeholder for other handler tests
func TestGetBudgetLinesByMonthIDHandler(t *testing.T) {
	mockStore := &MockStore{}
	handler := GetBudgetLinesByMonthIDHandler(mockStore)

	t.Run("successful retrieval", func(t *testing.T) {
		monthID := 1
		actualID1 := int64(101)
		actualAmount1 := 150.0
		actualID2 := int64(102)
		actualAmount2 := 25.0

		expectedBudgetLines := []store.BudgetLine{
			{ID: 1, MonthID: monthID, CategoryID: 1, Label: "Food", Expected: 200, ActualID: &actualID1, ActualAmount: &actualAmount1},
			{ID: 2, MonthID: monthID, CategoryID: 2, Label: "Gas", Expected: 50, ActualID: &actualID2, ActualAmount: &actualAmount2},
			{ID: 3, MonthID: monthID, CategoryID: 3, Label: "Rent", Expected: 1000, ActualID: nil, ActualAmount: nil}, // Case with no actual line
		}
		mockStore.MockGetBudgetLinesByMonthID = func(mID int) ([]store.BudgetLine, error) {
			if mID != monthID {
				t.Errorf("expected month ID %d, got %d", monthID, mID)
			}
			return expectedBudgetLines, nil
		}

		req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/budget-lines?month_id=%d", monthID), nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d. Body: %s", http.StatusOK, rr.Code, rr.Body.String())
		}

		var respBody []store.BudgetLine
		if err := json.Unmarshal(rr.Body.Bytes(), &respBody); err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}
		if len(respBody) != len(expectedBudgetLines) {
			t.Fatalf("expected %d budget lines, got %d", len(expectedBudgetLines), len(respBody))
		}

		for i, respLine := range respBody {
			expectedLine := expectedBudgetLines[i]
			if respLine.ID != expectedLine.ID || respLine.Label != expectedLine.Label || respLine.Expected != expectedLine.Expected {
				t.Errorf("mismatch in line %d basic data: expected %+v, got %+v", i, expectedLine, respLine)
			}
			// Compare ActualID
			if expectedLine.ActualID == nil && respLine.ActualID != nil {
				t.Errorf("mismatch in line %d ActualID: expected nil, got %v", i, *respLine.ActualID)
			} else if expectedLine.ActualID != nil && (respLine.ActualID == nil || *respLine.ActualID != *expectedLine.ActualID) {
				t.Errorf("mismatch in line %d ActualID: expected %v, got %v", i, *expectedLine.ActualID, respLine.ActualID)
			}

			// Compare ActualAmount
			if expectedLine.ActualAmount == nil && respLine.ActualAmount != nil {
				t.Errorf("mismatch in line %d ActualAmount: expected nil, got %v", i, *respLine.ActualAmount)
			} else if expectedLine.ActualAmount != nil && (respLine.ActualAmount == nil || *respLine.ActualAmount != *expectedLine.ActualAmount) {
				t.Errorf("mismatch in line %d ActualAmount: expected %v, got %v", i, *expectedLine.ActualAmount, respLine.ActualAmount)
			}
		}
	})

	t.Run("missing month_id query param", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/budget-lines", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("invalid month_id query param", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/budget-lines?month_id=abc", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("store error on retrieval", func(t *testing.T) {
		monthID := 2
		mockStore.MockGetBudgetLinesByMonthID = func(mID int) ([]store.BudgetLine, error) {
			return nil, fmt.Errorf("database error")
		}
		req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/budget-lines?month_id=%d", monthID), nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
		}
	})

	t.Run("no budget lines found", func(t *testing.T) {
		monthID := 3
		mockStore.MockGetBudgetLinesByMonthID = func(mID int) ([]store.BudgetLine, error) {
			return []store.BudgetLine{}, nil // Empty slice
		}
		req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/budget-lines?month_id=%d", monthID), nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
		}
		var respBody []store.BudgetLine
		if err := json.Unmarshal(rr.Body.Bytes(), &respBody); err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}
		if len(respBody) != 0 {
			t.Errorf("expected empty list, got %d items", len(respBody))
		}
	})
}

func TestUpdateBudgetLineHandler(t *testing.T) {
	mockStore := &MockStore{}
	handler := UpdateBudgetLineHandler(mockStore)

	t.Run("successful update", func(t *testing.T) {
		budgetLineID := int64(1)
		originalBudgetLine := &store.BudgetLine{
			ID:         int(budgetLineID),
			MonthID:    1,
			CategoryID: 1,
			Label:      "Old Label",
			Expected:   100.00,
		}
		updatePayload := struct {
			Label    *string  `json:"label"`
			Expected *float64 `json:"expected"`
		}{
			Label:    pointy.String("New Label"),
			Expected: pointy.Float64(200.00),
		}
		updatedBudgetLine := store.BudgetLine{
			ID:         int(budgetLineID),
			MonthID:    originalBudgetLine.MonthID,
			CategoryID: originalBudgetLine.CategoryID,
			Label:      *updatePayload.Label,
			Expected:   *updatePayload.Expected,
		}

		mockStore.MockGetBudgetLineByID = func(id int64) (*store.BudgetLine, error) {
			if id != budgetLineID {
				t.Errorf("expected GetBudgetLineByID with ID %d, got %d", budgetLineID, id)
			}
			return originalBudgetLine, nil
		}
		mockStore.MockUpdateBudgetLine = func(bl *store.BudgetLine) error {
			if bl.ID != int(budgetLineID) || bl.Label != *updatePayload.Label || bl.Expected != *updatePayload.Expected {
				t.Errorf("UpdateBudgetLine called with unexpected data: got %+v", bl)
			}
			return nil
		}

		payloadBytes, _ := json.Marshal(updatePayload)
		req := httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/budget-lines/%d", budgetLineID), bytes.NewReader(payloadBytes))
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d. Body: %s", http.StatusOK, rr.Code, rr.Body.String())
		}

		var respBody store.BudgetLine
		if err := json.Unmarshal(rr.Body.Bytes(), &respBody); err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}
		if respBody.Label != updatedBudgetLine.Label || respBody.Expected != updatedBudgetLine.Expected {
			t.Errorf("expected updated data %+v, got %+v", updatedBudgetLine, respBody)
		}
	})

	t.Run("budget line not found", func(t *testing.T) {
		budgetLineID := int64(2)
		mockStore.MockGetBudgetLineByID = func(id int64) (*store.BudgetLine, error) {
			return nil, nil // Simulate not found
		}
		payload := `{"label":"Any Label"}`
		req := httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/budget-lines/%d", budgetLineID), strings.NewReader(payload))
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
		}
	})

	t.Run("invalid budget line ID in path", func(t *testing.T) {
		payload := `{"label":"Any Label"}`
		req := httptest.NewRequest("PUT", "/api/v1/budget-lines/abc", strings.NewReader(payload))
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})
	
	t.Run("invalid request body - bad JSON", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/api/v1/budget-lines/1", strings.NewReader("{not_json"))
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d. Body: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
		}
	})

	t.Run("store error on GetBudgetLineByID", func(t *testing.T) {
		budgetLineID := int64(3)
		mockStore.MockGetBudgetLineByID = func(id int64) (*store.BudgetLine, error) {
			return nil, fmt.Errorf("database error on get")
		}
		payload := `{"label":"Any Label"}`
		req := httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/budget-lines/%d", budgetLineID), strings.NewReader(payload))
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
		}
	})

	t.Run("store error on UpdateBudgetLine", func(t *testing.T) {
		budgetLineID := int64(4)
		originalBudgetLine := &store.BudgetLine{ID: int(budgetLineID), Label: "Old"}
		mockStore.MockGetBudgetLineByID = func(id int64) (*store.BudgetLine, error) {
			return originalBudgetLine, nil
		}
		mockStore.MockUpdateBudgetLine = func(bl *store.BudgetLine) error {
			return fmt.Errorf("database error on update")
		}
		payload := `{"label":"New Label"}`
		req := httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/budget-lines/%d", budgetLineID), strings.NewReader(payload))
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
		}
	})
}

// Helper for tests needing pointers to basic types
var pointy = struct {
    String  func(s string) *string
    Float64 func(f float64) *float64
    Int     func(i int) *int
}{
    String:  func(s string) *string { return &s },
    Float64: func(f float64) *float64 { return &f },
    Int:     func(i int) *int { return &i },
}


func TestDeleteBudgetLineHandler(t *testing.T) {
	mockStore := &MockStore{}
	handler := DeleteBudgetLineHandler(mockStore)

	t.Run("successful deletion", func(t *testing.T) {
		budgetLineID := int64(1)
		deleteCalled := false
		mockStore.MockDeleteBudgetLine = func(id int64) error {
			if id != budgetLineID {
				t.Errorf("expected DeleteBudgetLine with ID %d, got %d", budgetLineID, id)
			}
			deleteCalled = true
			return nil
		}

		req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/budget-lines/%d", budgetLineID), nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusNoContent {
			t.Errorf("expected status %d, got %d. Body: %s", http.StatusNoContent, rr.Code, rr.Body.String())
		}
		if !deleteCalled {
			t.Error("expected MockDeleteBudgetLine to be called")
		}
	})

	t.Run("invalid budget line ID in path", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v1/budget-lines/abc", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("store error on deletion", func(t *testing.T) {
		budgetLineID := int64(2)
		mockStore.MockDeleteBudgetLine = func(id int64) error {
			return fmt.Errorf("database error")
		}
		req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/budget-lines/%d", budgetLineID), nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
		}
	})

	// Note: store.DeleteBudgetLine is expected to handle "not found" by returning an error.
	// If it were to return sql.ErrNoRows specifically and the handler distinguished this,
	// a separate test for http.StatusNotFound would be appropriate.
	// Current handler returns InternalServerError for any store error.
}

func TestUpdateActualLineHandler(t *testing.T) {
	mockStore := &MockStore{}
	handler := UpdateActualLineHandler(mockStore)

	t.Run("successful update", func(t *testing.T) {
		actualLineID := int64(1)
		originalActualLine := &store.ActualLine{
			ID:           actualLineID, // Corrected: ID type is int64
			BudgetLineID: 10,
			Actual:       50.00,
		}
		updatePayload := struct {
			Actual *float64 `json:"actual"`
		}{
			Actual: pointy.Float64(75.50), // Valid amount
		}
		
		var capturedActualLine store.ActualLine
		mockStore.MockGetActualLineByID = func(id int64) (*store.ActualLine, error) {
			if id != actualLineID {
				t.Errorf("expected GetActualLineByID with ID %d, got %d", actualLineID, id)
			}
			return originalActualLine, nil
		}
		mockStore.MockUpdateActualLine = func(al *store.ActualLine) error {
			capturedActualLine = *al // Capture the line passed to the mock
			if al.ID != actualLineID || al.Actual != *updatePayload.Actual {
				t.Errorf("UpdateActualLine called with unexpected data: got %+v, want actual %f", al, *updatePayload.Actual)
			}
			return nil
		}

		payloadBytes, _ := json.Marshal(updatePayload)
		req := httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/actual-lines/%d", actualLineID), bytes.NewReader(payloadBytes))
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d. Body: %s", http.StatusOK, rr.Code, rr.Body.String())
		}

		var respBody store.ActualLine
		if err := json.Unmarshal(rr.Body.Bytes(), &respBody); err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}
		if respBody.Actual != *updatePayload.Actual {
			t.Errorf("expected updated actual %.2f, got %.2f", *updatePayload.Actual, respBody.Actual)
		}
		if capturedActualLine.Actual != 75.50 {
             t.Errorf("expected actual amount %f in mock store, got %f", 75.50, capturedActualLine.Actual)
        }
	})
	
	t.Run("update with amount needing rounding", func(t *testing.T) {
		actualLineID := int64(5)
		originalActualLine := &store.ActualLine{ID: actualLineID, BudgetLineID: 11, Actual: 10.0}
		updatePayload := struct{ Actual *float64 `json:"actual"` }{Actual: pointy.Float64(123.456)}
		expectedRounded := 123.46
		
		var capturedActualLine store.ActualLine
		mockStore.MockGetActualLineByID = func(id int64) (*store.ActualLine, error) { return originalActualLine, nil }
		mockStore.MockUpdateActualLine = func(al *store.ActualLine) error {
			capturedActualLine = *al
			return nil
		}

		payloadBytes, _ := json.Marshal(updatePayload)
		req := httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/actual-lines/%d", actualLineID), bytes.NewReader(payloadBytes))
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d. Body: %s", http.StatusOK, rr.Code, rr.Body.String())
		}
		if capturedActualLine.Actual != expectedRounded {
			t.Errorf("expected rounded actual amount %f in mock store, got %f", expectedRounded, capturedActualLine.Actual)
		}
		var respBody store.ActualLine
		if err := json.Unmarshal(rr.Body.Bytes(), &respBody); err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}
		if respBody.Actual != expectedRounded { // Handler should return the store-modified (rounded) value
			t.Errorf("expected response actual %.2f, got %.2f", expectedRounded, respBody.Actual)
		}
	})

	t.Run("update with negative amount - handler validation", func(t *testing.T) {
		actualLineID := int64(6)
		// No need to mock store calls as handler should reject first
		updatePayload := struct{ Actual *float64 `json:"actual"` }{Actual: pointy.Float64(-10.50)}
		
		payloadBytes, _ := json.Marshal(updatePayload)
		req := httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/actual-lines/%d", actualLineID), bytes.NewReader(payloadBytes))
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d for negative amount, got %d. Body: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
		}
		if !strings.Contains(rr.Body.String(), "must be non-negative") {
			t.Errorf("expected error message about non-negative amount, got: %s", rr.Body.String())
		}
	})

	t.Run("actual line not found", func(t *testing.T) {
		actualLineID := int64(2) // Corrected: ID type is int64
		mockStore.MockGetActualLineByID = func(id int64) (*store.ActualLine, error) {
			return nil, nil // Simulate not found
		}
		payload := `{"actual":100.0}`
		req := httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/actual-lines/%d", actualLineID), strings.NewReader(payload))
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
		}
	})

	t.Run("invalid actual line ID in path", func(t *testing.T) {
		payload := `{"actual":100.0}`
		req := httptest.NewRequest("PUT", "/api/v1/actual-lines/abc", strings.NewReader(payload))
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})
	
	t.Run("invalid request body - bad JSON", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/api/v1/actual-lines/1", strings.NewReader("{not_json"))
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d. Body: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
		}
	})

	t.Run("missing actual field in body", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/api/v1/actual-lines/1", strings.NewReader("{}")) // Empty JSON
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d. Body: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
		}
	})

	t.Run("store error on GetActualLineByID", func(t *testing.T) {
		actualLineID := int64(3)
		mockStore.MockGetActualLineByID = func(id int64) (*store.ActualLine, error) {
			return nil, fmt.Errorf("database error on get")
		}
		payload := `{"actual":100.0}`
		req := httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/actual-lines/%d", actualLineID), strings.NewReader(payload))
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
		}
	})

	t.Run("store error on UpdateActualLine", func(t *testing.T) {
		actualLineID := int64(4) // Corrected: ID type is int64
		originalActualLine := &store.ActualLine{ID: actualLineID, Actual: 50.0}
		mockStore.MockGetActualLineByID = func(id int64) (*store.ActualLine, error) {
			return originalActualLine, nil
		}
		mockStore.MockUpdateActualLine = func(al *store.ActualLine) error {
			return fmt.Errorf("database error on update")
		}
		payload := `{"actual":100.0}`
		req := httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/actual-lines/%d", actualLineID), strings.NewReader(payload))
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
		}
	})
}
