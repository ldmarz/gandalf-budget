package store

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// NewStore initializes and returns a new sqlx.DB connection.
// It also ensures the database file exists or creates it.
func NewStore(dataSourceName string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite3", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database %s: %w", dataSourceName, err)
	}
	log.Printf("Successfully connected to database: %s", dataSourceName)
	return db, nil
}

// Store defines the interface for database operations.
// It will include methods for categories, budget lines, and actuals.
type Store interface {
	// Category methods
	GetAllCategories() ([]Category, error)
	CreateCategory(category *Category) error
	GetCategoryByID(id int64) (*Category, error)
	UpdateCategory(category *Category) error
	DeleteCategory(id int64) error

	// BudgetLine and ActualLine methods
	CreateBudgetLine(b *BudgetLine) (int64, error)
	GetBudgetLinesByMonthID(monthID int) ([]BudgetLine, error)
	UpdateBudgetLine(b *BudgetLine) error
	DeleteBudgetLine(id int64) error
	UpdateActualLine(a *ActualLine) error
	GetActualLineByID(id int64) (*ActualLine, error)
	GetBudgetLineByID(id int64) (*BudgetLine, error)
}

// sqlStore provides a concrete implementation of the Store interface
// using an sqlx.DB database connection.
type sqlStore struct {
	*sqlx.DB
}

// NewSQLStore creates a new sqlStore with the given database connection.
func NewSQLStore(db *sqlx.DB) Store {
	return &sqlStore{DB: db}
}

// RunMigrations reads all *.sql files from the specified directory and executes them.
// It attempts to make migrations idempotent by checking for "table already exists" errors.
func RunMigrations(db *sqlx.DB, migrationsDir string) error {
	log.Printf("Looking for migrations in: %s", migrationsDir)
	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
	if err != nil {
		return fmt.Errorf("failed to glob migration files: %w", err)
	}
	if len(files) == 0 {
		log.Println("No migration files found.")
		return nil
	}

	log.Printf("Found %d migration files. Applying...", len(files))
	for _, file := range files {
		log.Printf("Applying migration: %s", file)
		queryBytes, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}
		query := string(queryBytes)
		_, err = db.Exec(query)
		if err != nil {
			// Basic idempotency check for table creation.
			// SQLite error message for "table already exists" can vary.
			// This checks for a common substring.
			if strings.Contains(err.Error(), "already exists") {
				log.Printf("Table in migration %s likely already exists, skipping: %v", file, err)
				continue // Skip this migration file
			}
			return fmt.Errorf("failed to execute migration file %s: %w", file, err)
		}
		log.Printf("Successfully applied migration: %s", file)
	}
	log.Println("All migrations applied successfully.")
	return nil
}
