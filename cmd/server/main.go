package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"gandalf-budget/internal/app"
	httpinternal "gandalf-budget/internal/http"
	"gandalf-budget/internal/store"

	// "github.com/jmoiron/sqlx" // Not directly used here, but db is of this type
	_ "github.com/mattn/go-sqlite3"
)

//go:embed all:web/dist
var staticFiles embed.FS

func main() {
	log.Println("Starting Gandalf Budget application...")

	db, err := store.NewStore("budget.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	if err := store.RunMigrations(db, "internal/store/migrations"); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	if err := app.SeedInitialMonth(db); err != nil {
		log.Fatalf("Failed to seed initial data: %v", err)
	}


	log.Println("Setting up router...")
	distFS, err := fs.Sub(staticFiles, "web/dist")
	if err != nil {
		log.Fatalf("Failed to create sub VFS for web/dist: %v", err)
	}

	router := httpinternal.NewRouter(distFS, db) // Pass db connection

	log.Println("Starting HTTP server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("ListenAndServe error: %v", err)
	}
}
