package main

import (
	"embed"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os" // Required for Go version < 1.20 for fs.Sub default behavior with embed.FS

	"gandalf-budget/internal/app" // Assuming module name is gandalf-budget
	"gandalf-budget/internal/store"

	_ "github.com/mattn/go-sqlite3" // SQLite driver, ensure this is in go.mod
)

//go:embed all:web/dist
var staticFiles embed.FS

func main() {
	log.Println("Starting Gandalf Budget application...")

	// Initialize database
	db, err := store.NewStore("budget.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Run database migrations
	if err := store.RunMigrations(db, "internal/store/migrations"); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	// Seed initial data (e.g., current month if db is empty)
	if err := app.SeedInitialMonth(db); err != nil {
		log.Fatalf("Failed to seed initial data: %v", err)
	}

	log.Println("Setting up embedded static file server...")

	distFS, err := fs.Sub(staticFiles, "web/dist")
	if err != nil {
		log.Fatal("Failed to get sub VFS for web/dist: ", err)
	}

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.FS(distFS))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Serve index.html for root or typical SPA direct paths
		if r.URL.Path == "/" || r.URL.Path == "/index.html" { // Add other SPA paths if needed
			f, err_open := distFS.Open("index.html")
			if err_open != nil {
				http.Error(w, "index.html not found in embedded FS", http.StatusInternalServerError)
				log.Printf("Error opening index.html from embed: %v", err_open)
				return
			}
			defer f.Close()

			fi, err_stat := f.Stat()
			if err_stat != nil {
				http.Error(w, "Failed to stat index.html", http.StatusInternalServerError)
				log.Printf("Error stating index.html from embed: %v", err_stat)
				return
			}
			// Ensure f is io.ReadSeeker, which it should be for embed.FS files
            rs, ok := f.(io.ReadSeeker)
            if !ok {
                http.Error(w, "Embedded file does not support seeking", http.StatusInternalServerError)
                log.Printf("Error: embedded index.html is not an io.ReadSeeker")
                return
            }
			http.ServeContent(w, r, "index.html", fi.ModTime(), rs)
			return
		}
		// For other paths (e.g., /assets/...), let the FileServer handle them.
		// It will serve files if they exist in distFS, or 404 if not.
		fileServer.ServeHTTP(w, r)
	})

	log.Println("Starting HTTP server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("ListenAndServe error: %v", err)
	}
}
