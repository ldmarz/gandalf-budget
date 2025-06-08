package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"gandalf-budget/internal/store" // For store types
)

func TestGetBoardDataHandler(t *testing.T) {
	mockStore := &store.ReusableMockStore{}

	tests := []struct {
		name               string
		monthIDParam       string
		setupMock          func(ms *store.ReusableMockStore)
		expectedStatusCode int
		expectedBody       interface{}
	}{
		{
			name:         "Successful fetch",
			monthIDParam: "1",
			setupMock: func(ms *store.ReusableMockStore) {
				ms.MockGetBoardData = func(monthID int) (*store.BoardDataPayload, error) {
					if monthID != 1 {
						return nil, fmt.Errorf("unexpected monthID: %d", monthID)
					}
					return &store.BoardDataPayload{
						MonthID:   1,
						Year:      2024,
						MonthName: "January",
						BudgetLines: []store.BudgetLineWithActual{
							{ID: 1, MonthID: 1, CategoryID: 1, CategoryName: "Food", CategoryColor: "", Label: "Line 1", ExpectedAmount: 100.0, ActualAmount: 50.0},
							{ID: 2, MonthID: 1, CategoryID: 2, CategoryName: "Rent", CategoryColor: "", Label: "Line 2", ExpectedAmount: 150.0, ActualAmount: 75.0},
						},
					}, nil
				}
			},
			expectedStatusCode: http.StatusOK,
			expectedBody: &store.BoardDataPayload{
				MonthID:   1,
				Year:      2024,
				MonthName: "January",
				BudgetLines: []store.BudgetLineWithActual{
					{ID: 1, MonthID: 1, CategoryID: 1, CategoryName: "Food", CategoryColor: "", Label: "Line 1", ExpectedAmount: 100.0, ActualAmount: 50.0},
					{ID: 2, MonthID: 1, CategoryID: 2, CategoryName: "Rent", CategoryColor: "", Label: "Line 2", ExpectedAmount: 150.0, ActualAmount: 75.0},
				},
			},
		},
		{
			name:         "Empty board data",
			monthIDParam: "2",
			setupMock: func(ms *store.ReusableMockStore) {
				ms.MockGetBoardData = func(monthID int) (*store.BoardDataPayload, error) {
					if monthID != 2 {
						return nil, fmt.Errorf("unexpected monthID: %d", monthID)
					}
					return &store.BoardDataPayload{
						MonthID:     2,
						Year:        2024,
						MonthName:   "January",
						BudgetLines: []store.BudgetLineWithActual{},
					}, nil
				}
			},
			expectedStatusCode: http.StatusOK,
			expectedBody: &store.BoardDataPayload{
				MonthID:     2,
				Year:        2024,
				MonthName:   "January",
				BudgetLines: []store.BudgetLineWithActual{},
			},
		},
		{
			name:         "Store error on GetBoardData",
			monthIDParam: "3",
			setupMock: func(ms *store.ReusableMockStore) {
				ms.MockGetBoardData = func(monthID int) (*store.BoardDataPayload, error) {
					return nil, errors.New("database is down")
				}
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       map[string]string{"error": "Failed to fetch board data"},
		},
		{
			name:               "Invalid monthId in path (non-integer)",
			monthIDParam:       "abc",
			setupMock:          func(ms *store.ReusableMockStore) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       map[string]string{"error": "Invalid Month ID format"},
		},
		{
			name:               "Missing monthId in path",
			monthIDParam:       "",
			setupMock:          func(ms *store.ReusableMockStore) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       map[string]string{"error": "Month ID is required"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock(mockStore)

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

			if tc.expectedStatusCode == http.StatusOK {
				var actualBody store.BoardDataPayload
				if err := json.Unmarshal(rr.Body.Bytes(), &actualBody); err != nil {
					t.Fatalf("Could not unmarshal response body for OK status: %v. Body: %s", err, rr.Body.String())
				}
				if !reflect.DeepEqual(&actualBody, tc.expectedBody) {
					expectedJSON, _ := json.Marshal(tc.expectedBody)
					t.Errorf("Expected body %s, got %s", string(expectedJSON), rr.Body.String())
				}
			} else {
				var actualErrorBody map[string]string
				if err := json.Unmarshal(rr.Body.Bytes(), &actualErrorBody); err != nil {
					if expectedStr, ok := tc.expectedBody.(string); ok && rr.Body.String() == expectedStr {
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

	t.Run("Wrong HTTP method", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/api/v1/board-data/1", nil)
		if err != nil {
			t.Fatalf("Could not create request: %v", err)
		}
		rr := httptest.NewRecorder()
		handler := GetBoardDataHandler(mockStore)
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status code %d for wrong method, got %d", http.StatusMethodNotAllowed, rr.Code)
		}
		if rr.Body.String() != "Method not allowed" {
			t.Errorf("Expected body 'Method not allowed', got '%s'", rr.Body.String())
		}
	})
}

func ptrToInt64(v int64) *int64 {
	return &v
}

func ptrToFloat64(v float64) *float64 {
	return &v
}
