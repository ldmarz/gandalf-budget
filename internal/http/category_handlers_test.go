package http

import (
	"bytes"
	"database/sql" // Used by store, and sql.ErrNoRows is relevant for handler logic
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv" // For converting int64 to string for URL paths
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

// seedCategory is a helper to insert a category and return it with its ID.
func seedCategory(t *testing.T, db *sqlx.DB, category store.Category) store.Category {
	t.Helper()
	_, err := db.Exec(`INSERT INTO categories (name, color) VALUES (?, ?)`, category.Name, category.Color)
	if err != nil {
		t.Fatalf("Failed to seed category '%s': %v", category.Name, err)
	}
	// Retrieve the category to get its ID
	var seededCategory store.Category
	err = db.Get(&seededCategory, "SELECT id, name, color FROM categories WHERE name = ?", category.Name)
	if err != nil {
		t.Fatalf("Failed to retrieve seeded category '%s': %v", category.Name, err)
	}
	return seededCategory
}

func TestHandleUpdateCategory_Success(t *testing.T) {
	db := setupInMemoryDB(t)
	defer db.Close()

	initialCategory := seedCategory(t, db, store.Category{Name: "Initial Name", Color: "bg-initial-500"})

	updatePayload := store.Category{Name: "Updated Name", Color: "bg-updated-500"}
	payloadBytes, _ := json.Marshal(updatePayload)

	req, err := http.NewRequest("PUT", "/api/v1/categories/"+strconv.FormatInt(initialCategory.ID, 10), bytes.NewBuffer(payloadBytes))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := HandleUpdateCategory(db)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v. Body: %s", status, http.StatusOK, rr.Body.String())
	}

	var updatedCategory store.Category
	if err := json.Unmarshal(rr.Body.Bytes(), &updatedCategory); err != nil {
		t.Fatalf("Could not unmarshal response: %v. Body: %s", err, rr.Body.String())
	}

	if updatedCategory.Name != updatePayload.Name {
		t.Errorf("handler returned unexpected category name: got '%s' want '%s'", updatedCategory.Name, updatePayload.Name)
	}
	if updatedCategory.Color != updatePayload.Color {
		t.Errorf("handler returned unexpected category color: got '%s' want '%s'", updatedCategory.Color, updatePayload.Color)
	}
	if updatedCategory.ID != initialCategory.ID {
		t.Errorf("handler returned unexpected category ID: got %d want %d", updatedCategory.ID, initialCategory.ID)
	}

	// Verify in DB
	dbCategory, err := store.GetCategoryByID(db, initialCategory.ID)
	if err != nil {
		t.Fatalf("Could not fetch category from DB after update: %v", err)
	}
	if dbCategory == nil {
		t.Fatalf("Category with ID %d not found in DB after update", initialCategory.ID)
	}
	if dbCategory.Name != updatePayload.Name {
		t.Errorf("DB category name mismatch after update: got '%s' want '%s'", dbCategory.Name, updatePayload.Name)
	}
	if dbCategory.Color != updatePayload.Color {
		t.Errorf("DB category color mismatch after update: got '%s' want '%s'", dbCategory.Color, updatePayload.Color)
	}
}

func TestHandleUpdateCategory_NotFound(t *testing.T) {
	db := setupInMemoryDB(t)
	defer db.Close()

	updatePayload := store.Category{Name: "Updated Name", Color: "bg-updated-500"}
	payloadBytes, _ := json.Marshal(updatePayload)

	nonExistentID := int64(999)
	req, err := http.NewRequest("PUT", "/api/v1/categories/"+strconv.FormatInt(nonExistentID, 10), bytes.NewBuffer(payloadBytes))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := HandleUpdateCategory(db)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v. Body: %s", status, http.StatusNotFound, rr.Body.String())
	}
}

func TestHandleUpdateCategory_Validation(t *testing.T) {
	db := setupInMemoryDB(t)
	defer db.Close()

	initialCategory := seedCategory(t, db, store.Category{Name: "Initial Name", Color: "bg-initial-500"})

	tests := []struct {
		name           string
		payload        map[string]string
		expectedStatus int
		expectedBody   string
	}{
		{"MissingName", map[string]string{"color": "bg-updated-500"}, http.StatusBadRequest, "Category name and color are required for update"},
		{"MissingColor", map[string]string{"name": "Updated Name"}, http.StatusBadRequest, "Category name and color are required for update"},
		{"EmptyName", map[string]string{"name": "", "color": "bg-updated-500"}, http.StatusBadRequest, "Category name and color are required for update"},
		{"EmptyColor", map[string]string{"name": "Updated Name", "color": ""}, http.StatusBadRequest, "Category name and color are required for update"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payloadBytes, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("PUT", "/api/v1/categories/"+strconv.FormatInt(initialCategory.ID, 10), bytes.NewBuffer(payloadBytes))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := HandleUpdateCategory(db)
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

func TestHandleDeleteCategory_Success(t *testing.T) {
	db := setupInMemoryDB(t)
	defer db.Close()

	categoryToDelete := seedCategory(t, db, store.Category{Name: "To Delete", Color: "bg-delete-500"})

	req, err := http.NewRequest("DELETE", "/api/v1/categories/"+strconv.FormatInt(categoryToDelete.ID, 10), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := HandleDeleteCategory(db)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v. Body: %s", status, http.StatusNoContent, rr.Body.String())
	}

	// Verify in DB
	deletedCategory, err := store.GetCategoryByID(db, categoryToDelete.ID)
	if err != nil && err != sql.ErrNoRows { // store.GetCategoryByID returns nil, nil for not found, not sql.ErrNoRows directly
		t.Fatalf("Error fetching category from DB after delete: %v", err)
	}
	if deletedCategory != nil {
		t.Errorf("Category with ID %d was found in DB after delete, but should have been deleted. Found: %+v", categoryToDelete.ID, deletedCategory)
	}
}

func TestHandleDeleteCategory_NotFound(t *testing.T) {
	db := setupInMemoryDB(t)
	defer db.Close()

	nonExistentID := int64(999)
	req, err := http.NewRequest("DELETE", "/api/v1/categories/"+strconv.FormatInt(nonExistentID, 10), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := HandleDeleteCategory(db)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v. Body: %s", status, http.StatusNotFound, rr.Body.String())
	}
}
