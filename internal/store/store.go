package store

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func NewStore(dataSourceName string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite3", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database %s: %w", dataSourceName, err)
	}
	log.Printf("Successfully connected to database: %s", dataSourceName)
	return db, nil
}

type Store interface {
	GetAllCategories() ([]Category, error)
	CreateCategory(category *Category) error
	GetCategoryByID(id int64) (*Category, error)
	UpdateCategory(category *Category) error
	DeleteCategory(id int64) error

	CreateBudgetLine(b *BudgetLine) (int64, error)
	GetBudgetLinesByMonthID(monthID int) ([]BudgetLine, error)
	UpdateBudgetLine(b *BudgetLine) error
	DeleteBudgetLine(id int64) error
	UpdateActualLine(a *ActualLine) error
	GetActualLineByID(id int64) (*ActualLine, error)
	GetBudgetLineByID(id int64) (*BudgetLine, error)

	GetBoardData(monthID int) (*BoardDataPayload, error)

	CanFinalizeMonth(monthID int) (bool, string, error)
	FinalizeMonth(monthID int, snapJSON string) (int64, error)

	GetAnnualSnapshotsMetadataByYear(year int) ([]AnnualSnapMeta, error)
	GetAnnualSnapshotJSONByID(snapID int64) (string, error)
}

type sqlStore struct {
	DB *sqlx.DB
}

func (s *sqlStore) GetAnnualSnapshotsMetadataByYear(year int) ([]AnnualSnapMeta, error) {
	panic("implement me")
}

func NewSQLStore(db *sqlx.DB) Store {
	return &sqlStore{DB: db}
}

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
			if strings.Contains(err.Error(), "already exists") {
				log.Printf("Table in migration %s likely already exists, skipping: %v", file, err)
				continue
			}
			return fmt.Errorf("failed to execute migration file %s: %w", file, err)
		}
		log.Printf("Successfully applied migration: %s", file)
	}
	log.Println("All migrations applied successfully.")
	return nil
}

func (s *sqlStore) GetAnnualSnapshotJSONByID(snapID int64) (string, error) {
	var snapJSON string
	query := `SELECT snap_json FROM annual_snaps WHERE id = ?;`
	err := s.DB.Get(&snapJSON, query, snapID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", sql.ErrNoRows
		}
		return "", fmt.Errorf("error fetching annual snapshot JSON for ID %d: %w", snapID, err)
	}
	return snapJSON, nil
}
