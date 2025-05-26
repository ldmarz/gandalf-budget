package http

import (
	"bytes"
	"database/sql" // Used by store, and sql.ErrNoRows is relevant for handler logic
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"sort" // For sorting slices in tests, if needed for robust comparison

	"gandalf-budget/internal/store" // Adjust to your module path
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // In-memory SQLite for testing
)

// setupInMemoryDB creates a new in-memory SQLite database for testing
// and runs migrations.
func setupInMemoryDB(t *testing.T) *sqlx.DB {
	// Using ":memory:" for a clean DB for each test function call if setup is per-test.
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}

	// Simplified schema, matching the one in 001_init.sql for categories.
	schema := `
	CREATE TABLE categories (
	  id INTEGER PRIMARY KEY AUTOINCREMENT,
	  name TEXT UNIQUE NOT NULL,
	  color TEXT NOT NULL
	);`
	_, err = db.Exec(schema)
	if err != nil {
		db.Close()
		t.Fatalf("Failed to create schema: %v", err)
	}
	return db
}

func TestHandleGetCategories(t *testing.T) {
	db := setupInMemoryDB(t)
	defer db.Close()

	// Seed some data
	initialCategories := []store.Category{
		{Name: "Food", Color: "bg-red-500"}, // Will get ID 1
		{Name: "Travel", Color: "bg-blue-500"}, // Will get ID 2
		{Name: "Entertainment", Color: "bg-green-500"}, // Will get ID 3
	}
	for _, cat := range initialCategories {
		// Insert and ignore ID for initial seed, as it's auto-generated
		_, err := db.Exec(`INSERT INTO categories (name, color) VALUES (?, ?)`, cat.Name, cat.Color)
		if err != nil {
			t.Fatalf("Failed to seed category '%s': %v", cat.Name, err)
		}
	}
	
	// Update initialCategories with expected IDs and sort them by name for comparison
	// This assumes store.GetAllCategories sorts by name
	expectedCategories := []store.Category{
		{ID: 3, Name: "Entertainment", Color: "bg-green-500"},
		{ID: 1, Name: "Food", Color: "bg-red-500"},
		{ID: 2, Name: "Travel", Color: "bg-blue-500"},
	}
	sort.Slice(expectedCategories, func(i, j int) bool {
        return expectedCategories[i].Name < expectedCategories[j].Name
    })


	req, err := http.NewRequest("GET", "/api/v1/categories", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := HandleGetCategories(db)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v. Body: %s", status, http.StatusOK, rr.Body.String())
	}

	var actual []store.Category
	if err := json.Unmarshal(rr.Body.Bytes(), &actual); err != nil {
		t.Fatalf("Could not unmarshal response: %v. Body: %s", err, rr.Body.String())
	}

	if len(actual) != len(expectedCategories) {
		t.Fatalf("handler returned unexpected body size: got %d want %d categories. Body: %s", len(actual), len(expectedCategories), rr.Body.String())
	}
	
	// The actual response from store.GetAllCategories is sorted by name.
	// So, we compare `actual` directly with our `expectedCategories` which is also sorted by name.
	for i := range actual {
		if actual[i].Name != expectedCategories[i].Name || actual[i].Color != expectedCategories[i].Color {
			t.Errorf("Mismatch at index %d. Got Name: %s, Color: %s. Expected Name: %s, Color: %s",
				i, actual[i].Name, actual[i].Color, expectedCategories[i].Name, expectedCategories[i].Color)
		}
		// We don't check ID here because the seed data's IDs are assigned by the DB.
		// The important part is that the names and colors match in the correct order.
		// If ID checking is critical, fetch the seeded data first to get their assigned IDs.
		// For this version, the expectedCategories are manually assigned IDs based on insertion order,
		// then re-sorted by name. This should match the output of GetAllCategories.
		// The actual IDs are 1 for Food, 2 for Travel, 3 for Entertainment.
		// After sorting expectedCategories by name:
		// Entertainment (ID 3), Food (ID 1), Travel (ID 2)
		// The `actual` slice from the handler will also be sorted by name.
		// So, `actual[0]` should be Entertainment, `actual[1]` Food, `actual[2]` Travel.
		// Their IDs in the `actual` slice should reflect their DB IDs.
		if actual[i].ID != expectedCategories[i].ID {
             t.Errorf("Mismatch ID at index %d. Got ID: %d. Expected ID: %d for Name: %s",
                i, actual[i].ID, expectedCategories[i].ID, actual[i].Name)
        }
	}
}


func TestHandleCreateCategory(t *testing.T) {
	db := setupInMemoryDB(t)
	defer db.Close()

	categoryPayload := store.Category{Name: "Electronics", Color: "bg-gray-500"}
	payloadBytes, _ := json.Marshal(categoryPayload)

	req, err := http.NewRequest("POST", "/api/v1/categories", bytes.NewBuffer(payloadBytes))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := HandleCreateCategory(db)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v. Body: %s", status, http.StatusCreated, rr.Body.String())
	}

	var createdCategory store.Category
	if err := json.Unmarshal(rr.Body.Bytes(), &createdCategory); err != nil {
		t.Fatalf("Could not unmarshal response: %v", err)
	}

	if createdCategory.Name != categoryPayload.Name {
		t.Errorf("handler returned unexpected category name: got '%s' want '%s'", createdCategory.Name, categoryPayload.Name)
	}
	if createdCategory.Color != categoryPayload.Color {
		t.Errorf("handler returned unexpected category color: got '%s' want '%s'", createdCategory.Color, categoryPayload.Color)
	}
	if createdCategory.ID == 0 {
		t.Errorf("handler returned category with zero ID")
	}

	// Verify it's in the database
	var dbCategory store.Category
	err = db.Get(&dbCategory, "SELECT id, name, color FROM categories WHERE id = ?", createdCategory.ID)
	if err != nil {
		t.Fatalf("Could not fetch created category from DB: %v", err)
	}
	if dbCategory.Name != categoryPayload.Name {
		t.Errorf("DB category name mismatch: got '%s' want '%s'", dbCategory.Name, categoryPayload.Name)
	}
}

func TestHandleCreateCategory_Validation(t *testing.T) {
	db := setupInMemoryDB(t) 
	defer db.Close()

	tests := []struct {
		name           string
		payload        map[string]string 
		expectedStatus int
		expectedBody   string 
	}{
		{"MissingName", map[string]string{"color": "bg-red-500"}, http.StatusBadRequest, "Category name and color are required"},
		{"MissingColor", map[string]string{"name": "Test"}, http.StatusBadRequest, "Category name and color are required"},
		{"EmptyPayload", map[string]string{}, http.StatusBadRequest, "Category name and color are required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payloadBytes, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/api/v1/categories", bytes.NewBuffer(payloadBytes))
			req.Header.Set("Content-Type", "application/json")
			
			rr := httptest.NewRecorder()
			handler := HandleCreateCategory(db) 
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v. Body: %s", status, tt.expectedStatus, rr.Body.String())
			}
			if !bytes.Contains(rr.Body.Bytes(), []byte(tt.expectedBody)) {
				t.Errorf("handler returned unexpected body: got '%s' want to contain '%s'", rr.Body.String(), tt.expectedBody)
			}
		})
	}
}
// Add more tests for Update and Delete, including not found cases.
// func TestHandleUpdateCategory(t *testing.T) { /* ... */ }
// func TestHandleDeleteCategory(t *testing.T) { /* ... */ }
// func TestHandleUpdateCategory_NotFound(t *testing.T) { /* ... */ }
// func TestHandleDeleteCategory_NotFound(t *testing.T) { /* ... */ }
