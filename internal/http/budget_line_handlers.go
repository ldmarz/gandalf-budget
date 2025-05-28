package http

import (
	"encoding/json"
	"fmt"
	"gandalf-budget/internal/store"
	"log"
	"net/http"
	"strconv" // For Atoi
	"strings" // For TrimSuffix and Split
)

// CreateBudgetLineHandler handles the creation of a new budget line.
// It expects a JSON payload with month_id, category_id, label, and expected amount.
// It creates the budget line and an associated actual line with 0 actual amount.
func CreateBudgetLineHandler(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var bl store.BudgetLine
		if err := json.NewDecoder(r.Body).Decode(&bl); err != nil {
			log.Printf("Error decoding request body for create budget line: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Basic validation
		if bl.MonthID == 0 || bl.CategoryID == 0 || bl.Label == "" {
			http.Error(w, "Missing required fields: month_id, category_id, label", http.StatusBadRequest)
			return
		}

		budgetLineID, err := s.CreateBudgetLine(&bl)
		if err != nil {
			log.Printf("Error creating budget line: %v", err)
			http.Error(w, fmt.Sprintf("Failed to create budget line: %v", err), http.StatusInternalServerError)
			return
		}

		// Set the ID on the object that was passed in, as it's now populated
		bl.ID = int(budgetLineID) // Assuming budgetLineID is int64, and bl.ID is int

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(bl); err != nil {
			log.Printf("Error encoding created budget line to JSON: %v", err)
		}
	}
}

// UpdateActualLineHandler handles updates to an existing actual line.
// It expects the actual line ID in the URL path (e.g., /actual-lines/{id})
// and a JSON payload with the 'actual' amount.
func UpdateActualLineHandler(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		pathParts := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")
		idStr := pathParts[len(pathParts)-1]
		actualLineID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid actual line ID in path", http.StatusBadRequest)
			return
		}

		var reqBody struct {
			Actual *float64 `json:"actual"`
		}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			log.Printf("Error decoding request body for update actual line ID %d: %v", actualLineID, err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if reqBody.Actual == nil {
			http.Error(w, "Missing 'actual' field in request body", http.StatusBadRequest)
			return
		}
		if *reqBody.Actual < 0 {
			http.Error(w, "Invalid 'actual' amount: must be non-negative", http.StatusBadRequest)
			return
		}

		// Fetch existing actual line
		al, err := s.GetActualLineByID(actualLineID)
		if err != nil {
			log.Printf("Error fetching actual line ID %d for update: %v", actualLineID, err)
			http.Error(w, "Failed to retrieve actual line for update", http.StatusInternalServerError)
			return
		}
		if al == nil {
			http.Error(w, "Actual line not found", http.StatusNotFound)
			return
		}

		al.Actual = *reqBody.Actual

		if err := s.UpdateActualLine(al); err != nil {
			log.Printf("Error updating actual line ID %d: %v", actualLineID, err)
			http.Error(w, fmt.Sprintf("Failed to update actual line: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(al); err != nil {
			log.Printf("Error encoding updated actual line to JSON for ID %d: %v", actualLineID, err)
		}
	}
}

// UpdateBudgetLineHandler handles updates to an existing budget line.
// It expects the budget line ID in the URL path (e.g., /budget-lines/{id})
// and a JSON payload with fields to update (label, expected).
func UpdateBudgetLineHandler(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Assuming ID is extracted from path, e.g. using a router like Gorilla Mux.
		// For standard library, it's more complex. Let's assume ID is the last path segment.
		// e.g. /api/v1/budget-lines/123
		pathParts := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")
		idStr := pathParts[len(pathParts)-1]
		budgetLineID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid budget line ID in path", http.StatusBadRequest)
			return
		}

		var reqBody struct {
			Label    *string  `json:"label"`
			Expected *float64 `json:"expected"`
		}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			log.Printf("Error decoding request body for update budget line ID %d: %v", budgetLineID, err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Fetch existing budget line
		bl, err := s.GetBudgetLineByID(budgetLineID)
		if err != nil {
			log.Printf("Error fetching budget line ID %d for update: %v", budgetLineID, err)
			http.Error(w, "Failed to retrieve budget line for update", http.StatusInternalServerError)
			return
		}
		if bl == nil {
			http.Error(w, "Budget line not found", http.StatusNotFound)
			return
		}

		// Update fields if provided in the request
		if reqBody.Label != nil {
			bl.Label = *reqBody.Label
		}
		if reqBody.Expected != nil {
			bl.Expected = *reqBody.Expected
		}
		// MonthID and CategoryID are not updatable through this handler for now.

		if err := s.UpdateBudgetLine(bl); err != nil {
			log.Printf("Error updating budget line ID %d: %v", budgetLineID, err)
			http.Error(w, fmt.Sprintf("Failed to update budget line: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(bl); err != nil {
			log.Printf("Error encoding updated budget line to JSON for ID %d: %v", budgetLineID, err)
		}
	}
}

// DeleteBudgetLineHandler handles the deletion of a budget line.
// It expects the budget line ID in the URL path (e.g., /budget-lines/{id}).
func DeleteBudgetLineHandler(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		pathParts := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")
		idStr := pathParts[len(pathParts)-1]
		budgetLineID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid budget line ID in path", http.StatusBadRequest)
			return
		}

		err = s.DeleteBudgetLine(budgetLineID)
		if err != nil {
			// Consider specific error types, e.g., sql.ErrNoRows for not found
			log.Printf("Error deleting budget line ID %d: %v", budgetLineID, err)
			http.Error(w, fmt.Sprintf("Failed to delete budget line: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent) // 204 No Content
	}
}

// GetBudgetLinesByMonthIDHandler handles fetching budget lines for a specific month.
// It expects a 'month_id' query parameter.
func GetBudgetLinesByMonthIDHandler(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		monthIDStr := r.URL.Query().Get("month_id")
		if monthIDStr == "" {
			http.Error(w, "Missing 'month_id' query parameter", http.StatusBadRequest)
			return
		}

		monthID, err := strconv.Atoi(monthIDStr)
		if err != nil {
			http.Error(w, "Invalid 'month_id' query parameter: must be an integer", http.StatusBadRequest)
			return
		}

		budgetLines, err := s.GetBudgetLinesByMonthID(monthID)
		if err != nil {
			log.Printf("Error getting budget lines for month ID %d: %v", monthID, err)
			http.Error(w, fmt.Sprintf("Failed to get budget lines: %v", err), http.StatusInternalServerError)
			return
		}

		if budgetLines == nil {
			budgetLines = []store.BudgetLine{} // Return empty list instead of null
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(budgetLines); err != nil {
			log.Printf("Error encoding budget lines to JSON: %v", err)
		}
	}
}
