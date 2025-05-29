package store

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func newTestDB(t *testing.T) *sqlx.DB {
	t.Helper()

	db, err := sqlx.Connect("sqlite3", ":memory:?_foreign_keys=on")
	if err != nil {
		t.Fatalf("Failed to connect to in-memory sqlite3: %v", err)
	}

	migrationPath := filepath.Join("..", "..", "internal", "store", "migrations", "001_init.sql")
	
	queryBytes, err := os.ReadFile(migrationPath)
	if err != nil {
		altMigrationPath := filepath.Join("migrations", "001_init.sql")
		queryBytes, err = os.ReadFile(altMigrationPath)
		if err != nil {
			t.Fatalf("Failed to read migration file from %s or %s: %v", migrationPath, altMigrationPath, err)
		}
	}

	query := string(queryBytes)
	_, err = db.Exec(query)
	if err != nil {
		t.Fatalf("Failed to execute migration: %v", err)
	}

	t.Cleanup(func() {
		err := db.Close()
		if err != nil {
			t.Errorf("Failed to close test database: %v", err)
		}
	})

	return db
}

func createTestCategory(t *testing.T, db *sqlx.DB, name string, color string) int64 {
	t.Helper()
	res, err := db.Exec("INSERT INTO categories (name, color) VALUES (?, ?)", name, color)
	if err != nil {
		t.Fatalf("Failed to create test category %s: %v", name, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		t.Fatalf("Failed to get last insert ID for category %s: %v", name, err)
	}
	return id
}

func createTestMonth(t *testing.T, db *sqlx.DB, year int, month int, finalized bool) int64 {
	t.Helper()
	finalizedInt := 0
	if finalized {
		finalizedInt = 1
	}
	res, err := db.Exec("INSERT INTO months (year, month, finalized) VALUES (?, ?, ?)", year, month, finalizedInt)
	if err != nil {
		t.Fatalf("Failed to create test month %d-%d: %v", year, month, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		t.Fatalf("Failed to get last insert ID for month %d-%d: %v", year, month, err)
	}
	return id
}

func createTestBudgetLine(t *testing.T, db *sqlx.DB, monthID int64, categoryID int64, label string, expected float64) int64 {
	t.Helper()
	res, err := db.Exec("INSERT INTO budget_lines (month_id, category_id, label, expected) VALUES (?, ?, ?, ?)",
		monthID, categoryID, label, expected)
	if err != nil {
		t.Fatalf("Failed to create test budget line '%s': %v", label, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		t.Fatalf("Failed to get last insert ID for budget line '%s': %v", label, err)
	}
	return id
}

func createTestActualLine(t *testing.T, db *sqlx.DB, budgetLineID int64, actual float64) int64 {
	t.Helper()
	res, err := db.Exec("INSERT INTO actual_lines (budget_line_id, actual) VALUES (?, ?)", budgetLineID, actual)
	if err != nil {
		t.Fatalf("Failed to create test actual line for budget_line_id %d: %v", budgetLineID, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		t.Fatalf("Failed to get last insert ID for actual line (budget_line_id %d): %v", budgetLineID, err)
	}
	return id
}
