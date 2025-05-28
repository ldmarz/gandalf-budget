package http

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"gandalf-budget/internal/app"
	"gandalf-budget/internal/store"
)

// GetDashboardData handles the request for dashboard data.
// It expects a "month_id" query parameter.
func GetDashboardData(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		monthIDStr := r.URL.Query().Get("month_id")
		if monthIDStr == "" {
			http.Error(w, "month_id query parameter is required", http.StatusBadRequest)
			return
		}

		monthID, err := strconv.Atoi(monthIDStr)
		if err != nil {
			http.Error(w, "Invalid month_id: must be an integer", http.StatusBadRequest)
			return
		}

		// Fetch board data which now includes month details and budget lines with actuals
		boardData, err := s.GetBoardData(monthID) // Expects *store.BoardDataPayload
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) { // Assuming store returns sql.ErrNoRows if month not found
				http.Error(w, "Month not found", http.StatusNotFound)
			} else {
				http.Error(w, "Failed to fetch board data: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}
		// If boardData is nil for any other reason (should ideally not happen if err is nil)
		if boardData == nil {
			http.Error(w, "Failed to retrieve board data (nil payload)", http.StatusInternalServerError)
			return
		}

		// Fetch all categories to ensure every category is listed in the dashboard,
		// even if it has no budget lines for the selected month.
		allCategories, err := s.GetAllCategories() // Changed from GetCategories
		if err != nil {
			http.Error(w, "Failed to fetch categories: "+err.Error(), http.StatusInternalServerError)
			return
		}

		payload := app.DashboardPayload{
			MonthID: int(boardData.MonthID), // Convert int64 to int
			Year:    boardData.Year,
			Month:   boardData.MonthName, // Directly use MonthName
			// Totals and summaries will be calculated below
		}

		categorySummariesMap := make(map[int64]*app.CategorySummary) // Keyed by CategoryID (int64)

		for _, cat := range allCategories {
			// Ensure cat.ID is int64 if that's the key type
			// The store.Category has ID int64. app.CategorySummary.CategoryID is int.
			// This might require a cast or consistent types. Assuming app.CategorySummary.CategoryID is int64 for now.
			// Let's check app.DashboardPayload definition for CategorySummary.CategoryID type.
			// app.CategorySummary.CategoryID is int. store.Category.ID is int64.
			// This is an inconsistency. For now, I'll cast cat.ID to int for map key if necessary,
			// but the struct fields should align. Assuming app.CategorySummary.CategoryID can be int64.
			// If app.CategorySummary.CategoryID must be int, then cat.ID needs conversion.
			// The map key should be int64 to match store.Category.ID.
			// app.CategorySummary.CategoryID should be int64.
			// Let's assume app.CategorySummary.CategoryID is changed to int64. (This is outside current scope to change app struct)
			// For now, I'll use cat.ID (int64) as key and assign to app.CategorySummary.CategoryID (int) with conversion. This is not ideal.
			// The `app.CategorySummary` has `CategoryID int`.
			// `store.Category` has `ID int64`.
			// `store.BudgetLineWithActual` has `CategoryID int64`.
			// This requires careful casting.

			// Correct approach: categorySummariesMap key should be int64 (category ID type from store.BudgetLineWithActual)
			// app.CategorySummary.CategoryID should be int64.
			// For now, I'll assume app.CategorySummary.CategoryID is int as per current app struct.

			categorySummariesMap[cat.ID] = &app.CategorySummary{
				CategoryID:    int(cat.ID), // store.Category.ID is int64, app.CategorySummary.CategoryID is int
				CategoryName:  cat.Name,
				CategoryColor: cat.Color,
				BudgetLines:   []app.BudgetLineDetail{},
			}
		}

		for _, line := range boardData.BudgetLines { // line is store.BudgetLineWithActual
			payload.TotalExpected += line.ExpectedAmount
			payload.TotalActual += line.ActualAmount

			summary, ok := categorySummariesMap[line.CategoryID]
			if !ok {
				// This means a budget line exists for a category not in the allCategories list.
				// This should ideally not happen if data is consistent.
				// Or, it means a category was perhaps deleted after budget lines were created.
				// For robustness, we can create a summary on the fly or log an error.
				// For now, let's assume if it's not in allCategories, we might skip it or log.
				// However, BudgetLineWithActual now has CategoryName and CategoryColor.
				// So, we can create the summary if it doesn't exist.
				summary = &app.CategorySummary{
					CategoryID:    int(line.CategoryID),
					CategoryName:  line.CategoryName,
					CategoryColor: line.CategoryColor,
					BudgetLines:   []app.BudgetLineDetail{},
				}
				categorySummariesMap[line.CategoryID] = summary
			}

			summary.TotalExpected += line.ExpectedAmount
			summary.TotalActual += line.ActualAmount

			summary.BudgetLines = append(summary.BudgetLines, app.BudgetLineDetail{
				BudgetLineID:   int(line.ID), // store.BudgetLineWithActual.ID is int64
				Label:          line.Label,
				ExpectedAmount: line.ExpectedAmount,
				ActualAmount:   line.ActualAmount,
				Difference:     line.ExpectedAmount - line.ActualAmount,
			})
		}

		payload.CategorySummaries = []app.CategorySummary{} // Initialize
		for _, summary := range categorySummariesMap {
			summary.Difference = summary.TotalExpected - summary.TotalActual
			payload.CategorySummaries = append(payload.CategorySummaries, *summary)
		}

		payload.TotalDifference = payload.TotalExpected - payload.TotalActual

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			http.Error(w, "Failed to marshal JSON response: "+err.Error(), http.StatusInternalServerError)
		}
	}
}
