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

// GetAnnualReport handles requests for annual report data.
// It expects a "year" query parameter.
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
		// Basic validation for a reasonable year range
		currentYear := time.Now().Year()
		if year < 2000 || year > currentYear+5 { // Allow a bit into the future for planning
			http.Error(w, "Year out of reasonable range", http.StatusBadRequest)
			return
		}

		snapshotsMeta, err := s.GetAnnualSnapshotsMetadataByYear(year)
		if err != nil {
			// Log the error for server-side visibility
			// Consider a more specific error for "not found" if that's distinct from a general store error
			// For now, any error from the store is treated as a 500.
			http.Error(w, "Failed to retrieve annual report data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// If snapshotsMeta is nil (which can happen if the store returns nil, nil),
		// convert it to an empty slice for JSON marshalling to ensure `[]` instead of `null`.
		if snapshotsMeta == nil {
			snapshotsMeta = []store.AnnualSnapMeta{}
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(snapshotsMeta); err != nil {
			// Log this error as well
			http.Error(w, "Failed to marshal JSON response: "+err.Error(), http.StatusInternalServerError)
		}
	}
}

// GetSnapshotDetail handles requests for a specific annual snapshot's JSON data.
// It expects a "snapId" as part of the URL path.
func GetSnapshotDetail(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Path: /api/v1/reports/snapshots/{snapId}
		// Since we are using http.ServeMux, we need to parse the ID from the path manually.
		// Example: /api/v1/reports/snapshots/123
		// We expect the router to register this for a path prefix like "/api/v1/reports/snapshots/"
		// and then we extract the part after the prefix.
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
			if errors.Is(err, sql.ErrNoRows) { // Make sure to import "database/sql"
				http.Error(w, "Snapshot not found", http.StatusNotFound)
			} else {
				// Log the error for server-side visibility
				http.Error(w, "Failed to retrieve snapshot data: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		// Write the JSON string directly. No need to re-encode if it's already JSON.
		// Ensure snapJSON is indeed a JSON string.
		_, err = w.Write([]byte(snapJSON))
		if err != nil {
			// This error usually happens if the connection is closed or there's an issue writing.
			// Log it, but the header might have already been sent.
			// Consider logging this, but it's hard to send a different HTTP error at this point.
			log.Printf("Error writing snapshot JSON to response for ID %d: %v", snapID, err)
		}
	}
}
