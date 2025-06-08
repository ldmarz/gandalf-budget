package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"database/sql" // Required for sql.ErrNoRows
	"gandalf-budget/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestGetAnnualReport_Success(t *testing.T) {
	mockTime := time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC)
	expectedSnapshots := []store.AnnualSnapMeta{
		{ID: 1, MonthID: 10, Year: 2023, Month: "January", SnapCreatedAt: mockTime},
		{ID: 2, MonthID: 11, Year: 2023, Month: "February", SnapCreatedAt: mockTime.Add(time.Hour * 24 * 30)},
	}

	mockStore := &store.ReusableMockStore{
		MockGetAnnualSnapshotsMetadataByYear: func(year int) ([]store.AnnualSnapMeta, error) {
			if year == 2023 {
				return expectedSnapshots, nil
			}
			return nil, errors.New("unexpected year for mock")
		},
	}

	handler := GetAnnualReport(mockStore)
	req := httptest.NewRequest("GET", "/api/v1/reports/annual?year=2023", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")

	var actualSnapshots []store.AnnualSnapMeta
	err := json.NewDecoder(rr.Body).Decode(&actualSnapshots)
	assert.NoError(t, err, "could not unmarshal response body")
	assert.Equal(t, expectedSnapshots, actualSnapshots, "returned data does not match mock data")
}

func TestGetAnnualReport_NoData(t *testing.T) {
	mockStore := &store.ReusableMockStore{
		MockGetAnnualSnapshotsMetadataByYear: func(year int) ([]store.AnnualSnapMeta, error) {
			if year == 2024 {
				return []store.AnnualSnapMeta{}, nil
			}
			return nil, errors.New("unexpected year for mock")
		},
	}

	handler := GetAnnualReport(mockStore)
	req := httptest.NewRequest("GET", "/api/v1/reports/annual?year=2024", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
	assert.JSONEq(t, `[]`, rr.Body.String(), "response body should be an empty JSON array")
}

func TestGetAnnualReport_InvalidYearParameter(t *testing.T) {
	mockStore := &store.ReusableMockStore{}
	handler := GetAnnualReport(mockStore)

	req1 := httptest.NewRequest("GET", "/api/v1/reports/annual?year=abc", nil)
	rr1 := httptest.NewRecorder()
	handler.ServeHTTP(rr1, req1)

	assert.Equal(t, http.StatusBadRequest, rr1.Code, "handler returned wrong status code for non-integer year")
	assert.Contains(t, rr1.Body.String(), "Invalid year format: must be an integer", "incorrect error message for non-integer year")

	req2 := httptest.NewRequest("GET", "/api/v1/reports/annual", nil)
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req2)

	assert.Equal(t, http.StatusBadRequest, rr2.Code, "handler returned wrong status code for missing year")
	assert.Contains(t, rr2.Body.String(), "year query parameter is required", "incorrect error message for missing year")

	req3 := httptest.NewRequest("GET", "/api/v1/reports/annual?year=100", nil)
	rr3 := httptest.NewRecorder()
	handler.ServeHTTP(rr3, req3)
	assert.Equal(t, http.StatusBadRequest, rr3.Code)
	assert.Contains(t, rr3.Body.String(), "Year out of reasonable range")

}

func TestGetAnnualReport_StoreError(t *testing.T) {
	mockStore := &store.ReusableMockStore{
		MockGetAnnualSnapshotsMetadataByYear: func(year int) ([]store.AnnualSnapMeta, error) {
			return nil, errors.New("database connection failed")
		},
	}

	handler := GetAnnualReport(mockStore)
	req := httptest.NewRequest("GET", "/api/v1/reports/annual?year=2023", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "handler returned wrong status code for store error")
	expectedErrorMsg := "Failed to retrieve annual report data: database connection failed"
	assert.Equal(t, expectedErrorMsg, strings.TrimSpace(rr.Body.String()), "handler returned unexpected error message for store error")
}

func TestGetSnapshotDetail_Success(t *testing.T) {
	expectedJSON := `{"month_id":1,"year":2023,"month":"January","total_expected":1000,"total_actual":950,"categories":[]}`
	mockStore := &store.ReusableMockStore{
		MockGetAnnualSnapshotJSONByID: func(snapID int64) (string, error) {
			if snapID == 42 {
				return expectedJSON, nil
			}
			return "", errors.New("unexpected snapID for mock")
		},
	}

	handler := GetSnapshotDetail(mockStore)
	req := httptest.NewRequest("GET", "/api/v1/reports/snapshots/42", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"), "incorrect Content-Type header")
	assert.JSONEq(t, expectedJSON, rr.Body.String(), "response body does not match expected JSON")
}

func TestGetSnapshotDetail_NotFound(t *testing.T) {
	mockStore := &store.ReusableMockStore{
		MockGetAnnualSnapshotJSONByID: func(snapID int64) (string, error) {
			if snapID == 404 {
				return "", sql.ErrNoRows
			}
			return "", errors.New("unexpected snapID for mock")
		},
	}

	handler := GetSnapshotDetail(mockStore)
	req := httptest.NewRequest("GET", "/api/v1/reports/snapshots/404", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code, "handler returned wrong status code for not found")
	assert.Contains(t, rr.Body.String(), "Snapshot not found", "incorrect error message for not found")
}

func TestGetSnapshotDetail_InvalidID(t *testing.T) {
	mockStore := &store.ReusableMockStore{}
	handler := GetSnapshotDetail(mockStore)

	req := httptest.NewRequest("GET", "/api/v1/reports/snapshots/abc", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "handler returned wrong status code for invalid ID")
	assert.Contains(t, rr.Body.String(), "Invalid Snapshot ID format: must be an integer", "incorrect error message for invalid ID")

	reqEmpty := httptest.NewRequest("GET", "/api/v1/reports/snapshots/", nil)
	rrEmpty := httptest.NewRecorder()
	handler.ServeHTTP(rrEmpty, reqEmpty)
	assert.Equal(t, http.StatusBadRequest, rrEmpty.Code, "handler returned wrong status code for empty ID")
	assert.Contains(t, rrEmpty.Body.String(), "Snapshot ID missing in URL path", "incorrect error message for empty ID")
}

func TestGetSnapshotDetail_StoreError(t *testing.T) {
	mockStore := &store.ReusableMockStore{
		MockGetAnnualSnapshotJSONByID: func(snapID int64) (string, error) {
			return "", errors.New("internal database error")
		},
	}

	handler := GetSnapshotDetail(mockStore)
	req := httptest.NewRequest("GET", "/api/v1/reports/snapshots/77", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "handler returned wrong status code for store error")
	expectedMsg := "Failed to retrieve snapshot data: internal database error"
	assert.Equal(t, expectedMsg, strings.TrimSpace(rr.Body.String()), "incorrect error message for store error")
}

func TestGetSnapshotDetail_NotFound_WithSqlErrNoRows(t *testing.T) {
	mockStore := &store.ReusableMockStore{
		MockGetAnnualSnapshotJSONByID: func(snapID int64) (string, error) {
			if snapID == 404 {
				return "", sql.ErrNoRows
			}
			return "", errors.New("unexpected snapID for mock")
		},
	}

	handler := GetSnapshotDetail(mockStore)
	req := httptest.NewRequest("GET", "/api/v1/reports/snapshots/404", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code, "handler returned wrong status code for not found")
	assert.Contains(t, rr.Body.String(), "Snapshot not found", "incorrect error message for not found")
}
