package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"gandalf-budget/internal/store" // For store.BudgetLine
)

func TestGetBoardDataHandler(t *testing.T) {
	mockStore := &MockStore{} // Uses the mock from mock_store_test.go

	tests := []struct {
		name               string
		monthIDParam       string // Simulates the path parameter
		setupMock          func(ms *MockStore)
		expectedStatusCode int
		expectedBody       interface{} // Can be a slice of BudgetLine or an error map
	}{
		{
			name:         "Successful fetch",
			monthIDParam: "1",
			setupMock: func(ms *MockStore) {
				ms.GetBoardDataFunc = func(monthID int) ([]store.BudgetLine, error) {
					if monthID != 1 {
						return nil, fmt.Errorf("unexpected monthID: %d", monthID)
					}
					// ActualID and ActualAmount are pointers
					id1, amount1 := int64(101), float64(50.0)
					id2, amount2 := int64(102), float64(75.0)
					return []store.BudgetLine{
						{ID: 1, MonthID: 1, CategoryID: 1, Label: "Line 1", Expected: 100.0, ActualID: &id1, ActualAmount: &amount1},
						{ID: 2, MonthID: 1, CategoryID: 2, Label: "Line 2", Expected: 150.0, ActualID: &id2, ActualAmount: &amount2},
					}, nil
				}
			},
			expectedStatusCode: http.StatusOK,
			expectedBody: []store.BudgetLine{
				{ID: 1, MonthID: 1, CategoryID: 1, Label: "Line 1", Expected: 100.0, ActualID: ptrToInt64(101), ActualAmount: ptrToFloat64(50.0)},
				{ID: 2, MonthID: 1, CategoryID: 2, Label: "Line 2", Expected: 150.0, ActualID: ptrToInt64(102), ActualAmount: ptrToFloat64(75.0)},
			},
		},
		{
			name:         "Empty board data",
			monthIDParam: "2",
			setupMock: func(ms *MockStore) {
				ms.GetBoardDataFunc = func(monthID int) ([]store.BudgetLine, error) {
					if monthID != 2 {
						return nil, fmt.Errorf("unexpected monthID: %d", monthID)
					}
					return []store.BudgetLine{}, nil
				}
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       []store.BudgetLine{},
		},
		{
			name:         "Store error on GetBoardData",
			monthIDParam: "3",
			setupMock: func(ms *MockStore) {
				ms.GetBoardDataFunc = func(monthID int) ([]store.BudgetLine, error) {
					return nil, errors.New("database is down")
				}
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       map[string]string{"error": "Failed to fetch board data"},
		},
		{
			name:               "Invalid monthId in path (non-integer)",
			monthIDParam:       "abc",
			setupMock:          func(ms *MockStore) { /* No store call expected */ },
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       map[string]string{"error": "Invalid Month ID format"},
		},
		{
			name:               "Missing monthId in path",
			monthIDParam:       "", // Handled by router or handler if path is just /api/v1/board-data/
			setupMock:          func(ms *MockStore) { /* No store call expected */ },
			expectedStatusCode: http.StatusBadRequest, // Based on current handler logic for empty path part
			expectedBody:       map[string]string{"error": "Month ID is required"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock(mockStore) // Apply the mock setup for this test case

			// Create a request. The path needs to match how the handler extracts the ID.
			// The handler uses strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/board-data/"), "/")
			// So, we need a path like "/api/v1/board-data/1"
			path := fmt.Sprintf("/api/v1/board-data/%s", tc.monthIDParam)
			req, err := http.NewRequest("GET", path, nil)
			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}

			rr := httptest.NewRecorder()
			handler := GetBoardDataHandler(mockStore)
			handler.ServeHTTP(rr, req)

			if rr.Code != tc.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d. Body: %s", tc.expectedStatusCode, rr.Code, rr.Body.String())
			}

			// Check body for non-error cases or specific error messages
			if tc.expectedStatusCode == http.StatusOK {
				var actualBody []store.BudgetLine
				if err := json.Unmarshal(rr.Body.Bytes(), &actualBody); err != nil {
					// If it's expected to be an empty array, but it's something else non-JSON
					if reflect.DeepEqual(tc.expectedBody, []store.BudgetLine{}) && rr.Body.String() == "[]" {
						// This is fine
					} else {
						t.Fatalf("Could not unmarshal response body for OK status: %v. Body: %s", err, rr.Body.String())
					}
				}
				if !reflect.DeepEqual(actualBody, tc.expectedBody) {
					expectedJSON, _ := json.Marshal(tc.expectedBody)
					t.Errorf("Expected body %s, got %s", string(expectedJSON), rr.Body.String())
				}
			} else { // Error cases, expect a map[string]string error object
				var actualErrorBody map[string]string
				if err := json.Unmarshal(rr.Body.Bytes(), &actualErrorBody); err != nil {
					// If the raw string matches (e.g. "Method not allowed" which is not JSON)
					if expectedStr, ok := tc.expectedBody.(string); ok && rr.Body.String() == expectedStr {
						// This is fine
					} else {
						t.Logf("Raw error body: %s", rr.Body.String())
						t.Fatalf("Could not unmarshal error response body: %v. Body: %s", err, rr.Body.String())
					}
				} else if !reflect.DeepEqual(actualErrorBody, tc.expectedBody) {
					expectedJSON, _ := json.Marshal(tc.expectedBody)
					t.Errorf("Expected error body %s, got %s", string(expectedJSON), rr.Body.String())
				}
			}
		})
	}

	// Test wrong HTTP method
	t.Run("Wrong HTTP method", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/api/v1/board-data/1", nil) // Using POST
		if err != nil {
			t.Fatalf("Could not create request: %v", err)
		}
		rr := httptest.NewRecorder()
		handler := GetBoardDataHandler(mockStore) // mockStore can be empty for this
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status code %d for wrong method, got %d", http.StatusMethodNotAllowed, rr.Code)
		}
		// The default http.Error for StatusMethodNotAllowed writes "Method Not Allowed\n"
		// which is not JSON. So we check the string directly.
		// If we wanted JSON, the handler would need to write it.
		// For now, the current handler writes "Method not allowed" (without \n)
		if rr.Body.String() != "Method not allowed" { // Adjusted to match actual handler output
			t.Errorf("Expected body 'Method not allowed', got '%s'", rr.Body.String())
		}
	})
}

// Helper function to get a pointer to an int64 value for tests.
func ptrToInt64(v int64) *int64 {
	return &v
}

// Helper function to get a pointer to a float64 value for tests.
func ptrToFloat64(v float64) *float64 {
	return &v
}
