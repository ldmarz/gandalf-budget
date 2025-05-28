package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings" // For TrimPrefix

	"gandalf-budget/internal/store"
)

// GetBoardDataHandler fetches and returns budget lines with actuals for a given month.
func GetBoardDataHandler(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract monthId from path, e.g., /api/v1/board-data/123
		pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/board-data/"), "/")
		if len(pathParts) == 0 || pathParts[0] == "" {
			http.Error(w, "Month ID is required", http.StatusBadRequest)
			return
		}
		monthIDStr := pathParts[0]

		monthID, err := strconv.Atoi(monthIDStr)
		if err != nil {
			http.Error(w, "Invalid Month ID format", http.StatusBadRequest)
			return
		}

		boardData, err := s.GetBoardData(monthID)
		if err != nil {
			// Log the error server-side
			// log.Printf("Error fetching board data: %v", err)
			http.Error(w, "Failed to fetch board data", http.StatusInternalServerError)
			return
		}

		if boardData == nil {
			boardData = &store.BoardDataPayload{}
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(boardData); err != nil {
			// Log the error server-side
			// log.Printf("Error encoding board data to JSON: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}
