package http

import (
	"encoding/json"
	"log" // For server-side logging
	"net/http"
	"strconv"
	"strings" // For TrimPrefix

	"gandalf-budget/internal/store"
)

// FinalizeMonthHandler handles the logic for finalizing a month.
func FinalizeMonthHandler(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/months/"), "/")
		if len(pathParts) < 2 || pathParts[1] != "finalize" {
			http.Error(w, "Invalid path structure for finalize endpoint", http.StatusBadRequest)
			return
		}
		monthIDStr := pathParts[0]

		monthID, err := strconv.Atoi(monthIDStr)
		if err != nil {
			http.Error(w, "Invalid Month ID format", http.StatusBadRequest)
			return
		}

		// 1. Check if month can be finalized
		canFinalize, reason, err := s.CanFinalizeMonth(monthID)
		if err != nil {
			log.Printf("Error checking if month %d can be finalized: %v", monthID, err)
			http.Error(w, "Failed to check finalization status", http.StatusInternalServerError)
			return
		}
		if !canFinalize {
			http.Error(w, reason, http.StatusBadRequest) // Send the reason to the client
			return
		}

		// 2. Define dashboard payload for snapshot (using GetBoardData for now)
		// This part assumes GetBoardData is suitable for the snapshot.
		// It might need adjustment based on actual dashboard requirements in Milestone 5.
		boardData, err := s.GetBoardData(monthID)
		if err != nil {
			log.Printf("Error fetching board data for snapshot (month %d): %v", monthID, err)
			http.Error(w, "Failed to generate snapshot data", http.StatusInternalServerError)
			return
		}
		snapJSONBytes, err := json.Marshal(boardData)
		if err != nil {
			log.Printf("Error marshalling snapshot data for month %d: %v", monthID, err)
			http.Error(w, "Failed to prepare snapshot data", http.StatusInternalServerError)
			return
		}
		snapJSON := string(snapJSONBytes)

		// 3. Finalize the month
		newMonthID, err := s.FinalizeMonth(monthID, snapJSON)
		if err != nil {
			log.Printf("Error finalizing month %d: %v", monthID, err)
			http.Error(w, "Failed to finalize month", http.StatusInternalServerError)
			return
		}

		// 4. Respond with success
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":      "Month finalized successfully",
			"new_month_id": newMonthID,
		})
	}
}
