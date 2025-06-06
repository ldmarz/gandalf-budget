package http

import (
	"database/sql" // For sql.ErrNoRows
	"encoding/json"
	"log"
	"net/http"
	"strconv" // For parsing ID from path
	"strings" // For TrimPrefix

	"gandalf-budget/internal/store"
)

func HandleGetCategories(storage store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		categories, err := storage.GetAllCategories()
		if err != nil {
			log.Printf("Error in HandleGetCategories calling store.GetAllCategories: %v", err)
			http.Error(w, "Failed to retrieve categories", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(categories)
	}
}

func HandleCreateCategory(storage store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var newCategory store.Category
		if err := json.NewDecoder(r.Body).Decode(&newCategory); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		if newCategory.Name == "" || newCategory.Color == "" {
			http.Error(w, "Category name and color are required", http.StatusBadRequest)
			return
		}
		err := storage.CreateCategory(&newCategory)
		if err != nil {
			http.Error(w, "Failed to create category", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newCategory)
	}
}

func HandleUpdateCategory(storage store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		pathParts := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")
		idStr := pathParts[len(pathParts)-1]

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Printf("Error parsing category ID from path '%s': %v", idStr, err)
			http.Error(w, "Invalid category ID in path", http.StatusBadRequest)
			return
		}

		var categoryToUpdate store.Category
		if err := json.NewDecoder(r.Body).Decode(&categoryToUpdate); err != nil {
			log.Printf("Error decoding request body for update category ID %d: %v", id, err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		categoryToUpdate.ID = id

		if categoryToUpdate.Name == "" || categoryToUpdate.Color == "" {
			http.Error(w, "Category name and color are required for update", http.StatusBadRequest)
			return
		}

		err = storage.UpdateCategory(&categoryToUpdate)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Category not found or no changes needed", http.StatusNotFound)
			} else {
				log.Printf("Error in HandleUpdateCategory calling store.UpdateCategory for ID %d: %v", id, err)
				http.Error(w, "Failed to update category", http.StatusInternalServerError)
			}
			return
		}

		updatedCategory, err := storage.GetCategoryByID(id)
		if err != nil || updatedCategory == nil {
			log.Printf("Error fetching updated category ID %d after update: %v", id, err)
			http.Error(w, "Failed to retrieve category after update", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(updatedCategory); err != nil {
			log.Printf("Error encoding updated category to JSON for ID %d: %v", id, err)
		}
	}
}

func HandleDeleteCategory(storage store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		pathParts := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")
		idStr := pathParts[len(pathParts)-1]

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Printf("Error parsing category ID from path '%s' for delete: %v", idStr, err)
			http.Error(w, "Invalid category ID in path", http.StatusBadRequest)
			return
		}

		err = storage.DeleteCategory(id)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Category not found", http.StatusNotFound)
			} else {
				log.Printf("Error in HandleDeleteCategory calling store.DeleteCategory for ID %d: %v", id, err)
				http.Error(w, "Failed to delete category", http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
