package http

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gandalf-budget/internal/app"
	"gandalf-budget/internal/store"
)

func TestGetDashboardData_Success(t *testing.T) {
	mockStore := &store.ReusableMockStore{
		MockGetBoardData: func(monthID int) (*store.BoardDataPayload, error) {
			if monthID == 1 {
				return &store.BoardDataPayload{
					MonthID:   1,
					Year:      2023,
					MonthName: "December",
					BudgetLines: []store.BudgetLineWithActual{
						{ID: 101, MonthID: 1, CategoryID: 1, CategoryName: "Food", CategoryColor: "blue", Label: "Groceries", ExpectedAmount: 500.00, ActualAmount: 480.50},
						{ID: 102, MonthID: 1, CategoryID: 2, CategoryName: "Housing", CategoryColor: "red", Label: "Rent", ExpectedAmount: 1500.00, ActualAmount: 1500.00},
						{ID: 103, MonthID: 1, CategoryID: 1, CategoryName: "Food", CategoryColor: "blue", Label: "Eating Out", ExpectedAmount: 150.00, ActualAmount: 180.75},
						{ID: 104, MonthID: 1, CategoryID: 3, CategoryName: "Utilities", CategoryColor: "green", Label: "Internet", ExpectedAmount: 60.00, ActualAmount: 60.00},
					},
				}, nil
			}
			return nil, errors.New("board data not found for month_id")
		},
		MockGetAllCategories: func() ([]store.Category, error) {
			return []store.Category{
				{ID: 1, Name: "Food", Color: "blue"},
				{ID: 2, Name: "Housing", Color: "red"},
				{ID: 3, Name: "Utilities", Color: "green"},
				{ID: 4, Name: "Entertainment", Color: "purple"},
			}, nil
		},
	}

	handler := GetDashboardData(mockStore)
	req := httptest.NewRequest("GET", "/api/v1/dashboard?month_id=1", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v. Body: %s", status, http.StatusOK, rr.Body.String())
		return
	}

	var payload app.DashboardPayload
	err := json.NewDecoder(rr.Body).Decode(&payload)
	if err != nil {
		t.Fatalf("could not unmarshal response body: %v. Body: %s", err, rr.Body.String())
	}

	expectedTotalExpected := 500.00 + 1500.00 + 150.00 + 60.00
	if payload.TotalExpected != expectedTotalExpected {
		t.Errorf("payload.TotalExpected = %f; want %f", payload.TotalExpected, expectedTotalExpected)
	}

	expectedTotalActual := 480.50 + 1500.00 + 180.75 + 60.00
	if payload.TotalActual != expectedTotalActual {
		t.Errorf("payload.TotalActual = %f; want %f", payload.TotalActual, expectedTotalActual)
	}

	expectedTotalDifference := expectedTotalExpected - expectedTotalActual
	if payload.TotalDifference != expectedTotalDifference {
		t.Errorf("payload.TotalDifference = %f; want %f", payload.TotalDifference, expectedTotalDifference)
	}

	if payload.MonthID != 1 {
		t.Errorf("payload.MonthID = %d; want %d", payload.MonthID, 1)
	}
	if payload.Year != 2023 {
		t.Errorf("payload.Year = %d; want %d", payload.Year, 2023)
	}
	if payload.Month != "December" {
		t.Errorf("payload.Month = %s; want %s", payload.Month, "December")
	}

	if len(payload.CategorySummaries) != 4 {
		t.Fatalf("len(payload.CategorySummaries) = %d; want %d", len(payload.CategorySummaries), 4)
	}

	var cat1Summary app.CategorySummary
	foundCat1 := false
	for _, summary := range payload.CategorySummaries {
		if summary.CategoryID == 1 {
			cat1Summary = summary
			foundCat1 = true
			break
		}
	}
	if !foundCat1 {
		t.Fatalf("CategorySummary for Food (ID 1) not found")
	}

	if cat1Summary.CategoryName != "Food" {
		t.Errorf("cat1Summary.CategoryName = %s; want Food", cat1Summary.CategoryName)
	}

	expectedCat1TotalExpected := 500.00 + 150.00
	if cat1Summary.TotalExpected != expectedCat1TotalExpected {
		t.Errorf("cat1Summary.TotalExpected = %f; want %f", cat1Summary.TotalExpected, expectedCat1TotalExpected)
	}
	expectedCat1TotalActual := 480.50 + 180.75
	if cat1Summary.TotalActual != expectedCat1TotalActual {
		t.Errorf("cat1Summary.TotalActual = %f; want %f", cat1Summary.TotalActual, expectedCat1TotalActual)
	}
	if cat1Summary.Difference != (expectedCat1TotalExpected - expectedCat1TotalActual) {
		t.Errorf("cat1Summary.Difference = %f; want %f", cat1Summary.Difference, (expectedCat1TotalExpected - expectedCat1TotalActual))
	}
	if len(cat1Summary.BudgetLines) != 2 {
		t.Errorf("len(cat1Summary.BudgetLines) = %d; want %d", len(cat1Summary.BudgetLines), 2)
	}

	var bl1_1 app.BudgetLineDetail
	foundBL1_1 := false
	for _, bl := range cat1Summary.BudgetLines {
		if bl.Label == "Groceries" {
			bl1_1 = bl
			foundBL1_1 = true
			break
		}
	}
	if !foundBL1_1 {
		t.Fatalf("BudgetLineDetail for 'Groceries' not found in Food category")
	}

	if bl1_1.ExpectedAmount != 500.00 {
		t.Errorf("Groceries ExpectedAmount = %f; want %f", bl1_1.ExpectedAmount, 500.00)
	}
	if bl1_1.ActualAmount != 480.50 {
		t.Errorf("Groceries ActualAmount = %f; want %f", bl1_1.ActualAmount, 480.50)
	}
	if bl1_1.Difference != (500.00 - 480.50) {
		t.Errorf("Groceries Difference = %f; want %f", bl1_1.Difference, (500.00 - 480.50))
	}

	var cat2Summary app.CategorySummary
	foundCat2 := false
	for _, summary := range payload.CategorySummaries {
		if summary.CategoryID == 2 {
			cat2Summary = summary
			foundCat2 = true
			break
		}
	}
	if !foundCat2 {
		t.Fatalf("CategorySummary for Housing (ID 2) not found")
	}
	if cat2Summary.CategoryName != "Housing" {
		t.Errorf("cat2Summary.CategoryName = %s; want Housing", cat2Summary.CategoryName)
	}

	expectedCat2TotalExpected := 1500.00
	if cat2Summary.TotalExpected != expectedCat2TotalExpected {
		t.Errorf("cat2Summary.TotalExpected = %f; want %f", cat2Summary.TotalExpected, expectedCat2TotalExpected)
	}
	expectedCat2TotalActual := 1500.00
	if cat2Summary.TotalActual != expectedCat2TotalActual {
		t.Errorf("cat2Summary.TotalActual = %f; want %f", cat2Summary.TotalActual, expectedCat2TotalActual)
	}
	if len(cat2Summary.BudgetLines) != 1 {
		t.Errorf("len(cat2Summary.BudgetLines) = %d; want %d", len(cat2Summary.BudgetLines), 1)
	}
}

func TestGetDashboardData_InvalidMonthID(t *testing.T) {
	mockStore := &store.ReusableMockStore{}
	handler := GetDashboardData(mockStore)

	req := httptest.NewRequest("GET", "/api/v1/dashboard?month_id=abc", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expectedErrorMsg := "Invalid month_id: must be an integer"
	if body := rr.Body.String(); !strings.Contains(body, expectedErrorMsg) {
		t.Errorf("handler returned unexpected body: got %s want to contain %s", body, expectedErrorMsg)
	}
}

func TestGetDashboardData_MonthIDRequired(t *testing.T) {
	mockStore := &store.ReusableMockStore{}
	handler := GetDashboardData(mockStore)

	req := httptest.NewRequest("GET", "/api/v1/dashboard", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expectedErrorMsg := "month_id query parameter is required"
	if body := rr.Body.String(); !strings.Contains(body, expectedErrorMsg) {
		t.Errorf("handler returned unexpected body: got %s want to contain %s", body, expectedErrorMsg)
	}
}

func TestGetDashboardData_MonthNotFound(t *testing.T) {
	mockStore := &store.ReusableMockStore{
		MockGetMonthByID: func(id int) (*store.Month, error) {
			return nil, sql.ErrNoRows
		},
	}
	handler := GetDashboardData(mockStore)

	req := httptest.NewRequest("GET", "/api/v1/dashboard?month_id=999", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
	expectedErrorMsg := "Month not found"
	if body := rr.Body.String(); !strings.Contains(body, expectedErrorMsg) {
		t.Errorf("handler returned unexpected body: got %s want to contain %s", body, expectedErrorMsg)
	}
}

func TestGetDashboardData_ErrorGetMonthByID_Other(t *testing.T) {
	mockStore := &store.ReusableMockStore{
		MockGetMonthByID: func(id int) (*store.Month, error) {
			return nil, errors.New("some other database error")
		},
	}
	handler := GetDashboardData(mockStore)
	req := httptest.NewRequest("GET", "/api/v1/dashboard?month_id=1", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
	expectedErrorMsg := "Failed to fetch month details"
	if body := rr.Body.String(); !strings.Contains(body, expectedErrorMsg) {
		t.Errorf("handler returned unexpected body: got %s want to contain %s", body, expectedErrorMsg)
	}
}

func TestGetDashboardData_ErrorGetBoardData(t *testing.T) {
	mockStore := &store.ReusableMockStore{
		MockGetMonthByID: func(id int) (*store.Month, error) {
			return &store.Month{ID: 1, Year: 2023, Month: 12}, nil
		},
		MockGetBoardData: func(monthID int) ([]store.BudgetLine, error) {
			return nil, errors.New("failed to fetch board data")
		},
	}
	handler := GetDashboardData(mockStore)

	req := httptest.NewRequest("GET", "/api/v1/dashboard?month_id=1", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
	expectedErrorMsg := "Failed to fetch board data"
	if body := rr.Body.String(); !strings.Contains(body, expectedErrorMsg) {
		t.Errorf("handler returned unexpected body: got '%s' want to contain '%s'", body, expectedErrorMsg)
	}
}

func TestGetDashboardData_ErrorGetAllCategories(t *testing.T) {
	mockStore := &store.ReusableMockStore{
		MockGetMonthByID: func(id int) (*store.Month, error) {
			return &store.Month{ID: 1, Year: 2023, Month: 12}, nil
		},
		MockGetBoardData: func(monthID int) ([]store.BudgetLine, error) {
			return []store.BudgetLine{}, nil
		},
		MockGetAllCategories: func() ([]store.Category, error) {
			return nil, errors.New("failed to fetch categories")
		},
	}
	handler := GetDashboardData(mockStore)

	req := httptest.NewRequest("GET", "/api/v1/dashboard?month_id=1", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
	expectedErrorMsg := "Failed to fetch categories"
	if body := rr.Body.String(); !strings.Contains(body, expectedErrorMsg) {
		t.Errorf("handler returned unexpected body: got '%s' want to contain '%s'", body, expectedErrorMsg)
	}
}

func floatEquals(a, b, tolerance float64) bool {
	if (a-b) < tolerance && (b-a) < tolerance {
		return true
	}
	return false
}

func TestMain(m *testing.M) {
	m.Run()
}
