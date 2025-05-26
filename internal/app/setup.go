package app

import (
	"database/sql"
	"fmt"
	"log"
	"strings" // Required for strings.Contains
	"time"

	"github.com/jmoiron/sqlx"
)

// SeedInitialMonth checks if the months table is empty and, if so,
// seeds it with the current calendar month and year.
func SeedInitialMonth(db *sqlx.DB) error {
	log.Println("Checking if initial month seeding is required...")
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM months")
	if err != nil {
		// This might happen if the table doesn't exist yet (e.g., first run before migrations)
		// Or if there's some other DB error.
		// If migrations haven't run, this will fail. This function should run AFTER migrations.
		log.Printf("Could not query count from months table (migrations might not have run yet): %v", err)
		// Return nil because the migration process should handle table creation.
		// If migrations did run and this still fails, it's a more serious issue.
        // For now, assume migrations will create the table. If they did, and count fails, that's an error.
        // Let's refine: if the error is "no such table", it's fine here. Otherwise, it's an actual error.
        if err != sql.ErrNoRows && !strings.Contains(err.Error(), "no such table") {
             return fmt.Errorf("failed to query count from months table: %w", err)
        }
        // If "no such table" or ErrNoRows, means table is effectively empty or not yet created by migration.
        // We'll proceed assuming migration will create it, and then this logic will be fine on next startup.
        // Or, if it *was* created and is empty, count will be 0.
        log.Println("Months table likely not yet created or empty; proceeding with seeding check logic.")
        count = 0 // Treat as empty
	}

	if count == 0 {
		log.Println("Months table is empty. Seeding with current month and year.")
		currentYear := time.Now().Year()
		currentMonth := int(time.Now().Month())

		query := `INSERT INTO months (year, month, finalized) VALUES (?, ?, ?)`
		_, err = db.Exec(query, currentYear, currentMonth, 0)
		if err != nil {
			return fmt.Errorf("failed to insert initial month: %w", err)
		}
		log.Printf("Successfully seeded months table with %d-%02d.", currentYear, currentMonth)
	} else {
		log.Printf("Months table is not empty (count: %d). No seeding required.", count)
	}
	return nil
}
