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

		boardData, err := s.GetBoardData(monthID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "Month not found", http.StatusNotFound)
			} else {
				http.Error(w, "Failed to fetch board data: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}
		if boardData == nil {
			http.Error(w, "Failed to retrieve board data (nil payload)", http.StatusInternalServerError)
			return
		}

		allCategories, err := s.GetAllCategories()
		if err != nil {
			http.Error(w, "Failed to fetch categories: "+err.Error(), http.StatusInternalServerError)
			return
		}

		payload := app.DashboardPayload{
			MonthID: int(boardData.MonthID),
			Year:    boardData.Year,
			Month:   boardData.MonthName,
		}

		categorySummariesMap := make(map[int64]*app.CategorySummary)

		for _, cat := range allCategories {
			categorySummariesMap[cat.ID] = &app.CategorySummary{
				CategoryID:    int(cat.ID),
				CategoryName:  cat.Name,
				CategoryColor: cat.Color,
				BudgetLines:   []app.BudgetLineDetail{},
			}
		}

		for _, line := range boardData.BudgetLines {
			payload.TotalExpected += line.ExpectedAmount
			payload.TotalActual += line.ActualAmount

			summary, ok := categorySummariesMap[line.CategoryID]
			if !ok {
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
				BudgetLineID:   int(line.ID),
				Label:          line.Label,
				ExpectedAmount: line.ExpectedAmount,
				ActualAmount:   line.ActualAmount,
				Difference:     line.ExpectedAmount - line.ActualAmount,
			})
		}

		payload.CategorySummaries = []app.CategorySummary{}
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
