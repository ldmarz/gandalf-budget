package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"gandalf-budget/internal/store" // For store types
)

func TestFinalizeMonthHandler(t *testing.T) {
	mockStore := &MockStore{}

	sampleBoardData := []store.BudgetLine{
		{ID: 1, MonthID: 1, CategoryID: 1, Label: "Item 1", Expected: 100, ActualAmount: ptrToFloat64(90)},
	}
	expectedSnapJSONBytes, _ := json.Marshal(sampleBoardData)
	expectedSnapJSON := string(expectedSnapJSONBytes)

	tests := []struct {
		name               string
		monthIDParam       string
		pathSuffix         string
		httpMethod         string
		setupMock          func(ms *MockStore)
		expectedStatusCode int
		expectedBody       interface{}
	}{
		{
			name:         "Successful month finalization",
			monthIDParam: "1",
			pathSuffix:   "/finalize",
			httpMethod:   http.MethodPut,
			setupMock: func(ms *MockStore) {
				ms.CanFinalizeMonthFunc = func(monthID int) (bool, string, error) {
					if monthID != 1 {
						return false, "mock error: unexpected monthID for CanFinalizeMonth", fmt.Errorf("unexpected monthID: %d", monthID)
					}
					return true, "", nil
				}
				ms.GetBoardDataFunc = func(monthID int) ([]store.BudgetLine, error) {
					if monthID != 1 {
						return nil, fmt.Errorf("mock error: unexpected monthID for GetBoardData: %d", monthID)
					}
					return sampleBoardData, nil
				}
				ms.FinalizeMonthFunc = func(monthID int, snapJSON string) (int64, error) {
					if monthID != 1 {
						return 0, fmt.Errorf("mock error: unexpected monthID for FinalizeMonth: %d", monthID)
					}
					if snapJSON != expectedSnapJSON {
						return 0, fmt.Errorf("mock error: snapJSON mismatch. Got %s, want %s", snapJSON, expectedSnapJSON)
					}
					return 2, nil
				}
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       map[string]interface{}{"message": "Month finalized successfully", "new_month_id": float64(2)},
		},
		{
			name:         "Cannot finalize month (CanFinalizeMonth returns false)",
			monthIDParam: "2",
			pathSuffix:   "/finalize",
			httpMethod:   http.MethodPut,
			setupMock: func(ms *MockStore) {
				ms.CanFinalizeMonthFunc = func(monthID int) (bool, string, error) {
					return false, "Actuals not set for all lines", nil
				}
				ms.GetBoardDataFunc = nil
				ms.FinalizeMonthFunc = nil
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       map[string]string{"error": "Actuals not set for all lines"},
		},
		{
			name:         "Error in CanFinalizeMonth",
			monthIDParam: "3",
			pathSuffix:   "/finalize",
			httpMethod:   http.MethodPut,
			setupMock: func(ms *MockStore) {
				ms.CanFinalizeMonthFunc = func(monthID int) (bool, string, error) {
					return false, "", errors.New("DB error checking finalization")
				}
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       map[string]string{"error": "Failed to check finalization status"},
		},
		{
			name:         "Error in GetBoardData for snapshot",
			monthIDParam: "4",
			pathSuffix:   "/finalize",
			httpMethod:   http.MethodPut,
			setupMock: func(ms *MockStore) {
				ms.CanFinalizeMonthFunc = func(monthID int) (bool, string, error) { return true, "", nil }
				ms.GetBoardDataFunc = func(monthID int) ([]store.BudgetLine, error) {
					return nil, errors.New("DB error fetching board data")
				}
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       map[string]string{"error": "Failed to generate snapshot data"},
		},
		{
			name:         "Error in FinalizeMonth store method",
			monthIDParam: "5",
			pathSuffix:   "/finalize",
			httpMethod:   http.MethodPut,
			setupMock: func(ms *MockStore) {
				ms.CanFinalizeMonthFunc = func(monthID int) (bool, string, error) { return true, "", nil }
				ms.GetBoardDataFunc = func(monthID int) ([]store.BudgetLine, error) { return sampleBoardData, nil }
				ms.FinalizeMonthFunc = func(monthID int, snapJSON string) (int64, error) {
					return 0, errors.New("DB error during finalization transaction")
				}
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       map[string]string{"error": "Failed to finalize month"},
		},
		{
			name:               "Invalid monthId in path (non-integer)",
			monthIDParam:       "abc",
			pathSuffix:         "/finalize",
			httpMethod:         http.MethodPut,
			setupMock:          func(ms *MockStore) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       map[string]string{"error": "Invalid Month ID format"},
		},
		{
			name:               "Malformed path (missing /finalize)",
			monthIDParam:       "1",
			pathSuffix:         "/somethingelse",
			httpMethod:         http.MethodPut,
			setupMock:          func(ms *MockStore) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       map[string]string{"error": "Invalid path structure for finalize endpoint"},
		},
		{
			name:               "Malformed path (too short)",
			monthIDParam:       "",
			pathSuffix:         "/finalize",
			httpMethod:         http.MethodPut,
			setupMock:          func(ms *MockStore) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       map[string]string{"error": "Invalid path structure for finalize endpoint"},
		},
        {
			name:               "Malformed path (no monthID)",
			monthIDParam:       "",
			pathSuffix:         "/finalize",
			httpMethod:         http.MethodPut,
			setupMock:          func(ms *MockStore) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       map[string]string{"error": "Invalid path structure for finalize endpoint"},
		},
		{
			name:         "Wrong HTTP method (GET)",
			monthIDParam: "1",
			pathSuffix:   "/finalize",
			httpMethod:   http.MethodGet,
			setupMock:    func(ms *MockStore) {},
			expectedStatusCode: http.StatusMethodNotAllowed,
			expectedBody:       "Method not allowed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockStore.CanFinalizeMonthFunc = nil
			mockStore.GetBoardDataFunc = nil
			mockStore.FinalizeMonthFunc = nil
			tc.setupMock(mockStore)

			path := fmt.Sprintf("/api/v1/months/%s%s", tc.monthIDParam, tc.pathSuffix)
			req, err := http.NewRequest(tc.httpMethod, path, nil)
			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}

			rr := httptest.NewRecorder()
			handler := FinalizeMonthHandler(mockStore)
			handler.ServeHTTP(rr, req)

			if rr.Code != tc.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d. Body: %s", tc.expectedStatusCode, rr.Code, rr.Body.String())
			}
			
			if contentType := rr.Header().Get("Content-Type"); strings.Contains(contentType, "application/json") {
				var actualBodyMap map[string]interface{}
				if err := json.Unmarshal(rr.Body.Bytes(), &actualBodyMap); err != nil {
					t.Fatalf("Could not unmarshal response body to map: %v. Body: %s", err, rr.Body.String())
				}

				expectedBodyBytes, _ := json.Marshal(tc.expectedBody)
				var expectedBodyMap map[string]interface{}
				_ = json.Unmarshal(expectedBodyBytes, &expectedBodyMap)
				
				if !reflect.DeepEqual(actualBodyMap, expectedBodyMap) {
					t.Errorf("Expected JSON body %+v, got %+v", expectedBodyMap, actualBodyMap)
				}
			} else {
				if expectedStr, ok := tc.expectedBody.(string); ok {
					if strings.TrimSpace(rr.Body.String()) != expectedStr {
						t.Errorf("Expected body string '%s', got '%s'", expectedStr, rr.Body.String())
					}
				} else if rr.Body.Len() > 0 {
					t.Errorf("Unexpected non-JSON body: %s", rr.Body.String())
				}
			}
		})
	}
}
