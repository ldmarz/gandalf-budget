package http

import (
	"database/sql" // Added
	"encoding/json"
	"errors" // Added
	"log"    // Added
	"net/http"
	"strconv"
	"strings" // Added
	"time"

	"gandalf-budget/internal/store"
)

func GetAnnualReport(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		yearStr := r.URL.Query().Get("year")
		if yearStr == "" {
			http.Error(w, "year query parameter is required", http.StatusBadRequest)
			return
		}

		year, err := strconv.Atoi(yearStr)
		if err != nil {
			http.Error(w, "Invalid year format: must be an integer", http.StatusBadRequest)
			return
		}
		currentYear := time.Now().Year()
		if year < 2000 || year > currentYear+5 {
			http.Error(w, "Year out of reasonable range", http.StatusBadRequest)
			return
		}

		snapshotsMeta, err := s.GetAnnualSnapshotsMetadataByYear(year)
		if err != nil {
			http.Error(w, "Failed to retrieve annual report data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if snapshotsMeta == nil {
			snapshotsMeta = []store.AnnualSnapMeta{}
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(snapshotsMeta); err != nil {
			http.Error(w, "Failed to marshal JSON response: "+err.Error(), http.StatusInternalServerError)
		}
	}
}

func GetSnapshotDetail(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pathPrefix := "/api/v1/reports/snapshots/"
		idStr := strings.TrimPrefix(r.URL.Path, pathPrefix)

		if idStr == "" {
			http.Error(w, "Snapshot ID missing in URL path", http.StatusBadRequest)
			return
		}

		snapID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid Snapshot ID format: must be an integer", http.StatusBadRequest)
			return
		}

		snapJSON, err := s.GetAnnualSnapshotJSONByID(snapID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "Snapshot not found", http.StatusNotFound)
			} else {
				http.Error(w, "Failed to retrieve snapshot data: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write([]byte(snapJSON))
		if err != nil {
			log.Printf("Error writing snapshot JSON to response for ID %d: %v", snapID, err)
		}
	}
}
