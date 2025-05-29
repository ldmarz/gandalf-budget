package app

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

func SeedInitialMonth(db *sqlx.DB) error {
	log.Println("Checking if initial month seeding is required...")
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM months")
	if err != nil {
		log.Printf("Could not query count from months table (migrations might not have run yet): %v", err)
        if err != sql.ErrNoRows && !strings.Contains(err.Error(), "no such table") {
             return fmt.Errorf("failed to query count from months table: %w", err)
        }
        log.Println("Months table likely not yet created or empty; proceeding with seeding check logic.")
        count = 0
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
