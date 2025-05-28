package http

import (
	"io"
	"io/fs"
	"log"
	"net/http"
	"path"
	"strings" // Required for path manipulation

	"github.com/jmoiron/sqlx"

	"gandalf-budget/internal/store" // Added for store.Store
)

func NewRouter(staticFS fs.FS, db *sqlx.DB) *http.ServeMux {
	mux := http.NewServeMux()
	appStore := store.NewSQLStore(db) // Create store instance

	mux.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"status": "ok"}`)
		log.Println("HIT: /api/v1/health")
	})
	
	// Handler for /api/v1/categories (collections)
	mux.HandleFunc("/api/v1/categories", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			HandleGetCategories(appStore)(w, r)
		case http.MethodPost:
			HandleCreateCategory(appStore)(w, r)
		default:
			http.Error(w, "Method not allowed for /api/v1/categories collection", http.StatusMethodNotAllowed)
		}
	})

	// Handler for /api/v1/categories/:id (specific item)
	mux.HandleFunc("/api/v1/categories/", func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/categories/")
		idStr = strings.TrimSuffix(idStr, "/")

		if idStr == "" {
			http.Error(w, "Category ID missing", http.StatusBadRequest)
			return
		}
		
		// ID parsing will be done in handlers, as they might need the raw string or convert to int64.
		// Handlers are responsible for ensuring the ID is valid for their use case.

		switch r.Method {
		case http.MethodPut:
			HandleUpdateCategory(appStore)(w, r)
		case http.MethodDelete:
			HandleDeleteCategory(appStore)(w, r) // Added
		// case http.MethodGet: // For getting a single category by ID
			// HandleGetCategoryByID(appStore)(w,r) 
		default:
			http.Error(w, "Method not allowed for specific category item (/api/v1/categories/:id)", http.StatusMethodNotAllowed)
		}
	})

	// Handler for /api/v1/budget-lines (collection)
	mux.HandleFunc("/api/v1/budget-lines", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			CreateBudgetLineHandler(appStore)(w, r)
		case http.MethodGet:
			GetBudgetLinesByMonthIDHandler(appStore)(w, r) // Uses query param, not path param
		default:
			http.Error(w, "Method not allowed for /api/v1/budget-lines collection", http.StatusMethodNotAllowed)
		}
	})

	// Handler for /api/v1/budget-lines/:id (specific item)
	mux.HandleFunc("/api/v1/budget-lines/", func(w http.ResponseWriter, r *http.Request) {
		// ID parsing will be done in handlers
		switch r.Method {
		case http.MethodPut:
			UpdateBudgetLineHandler(appStore)(w, r)
		case http.MethodDelete:
			DeleteBudgetLineHandler(appStore)(w, r)
		default:
			http.Error(w, "Method not allowed for specific budget line item (/api/v1/budget-lines/:id)", http.StatusMethodNotAllowed)
		}
	})

	// Handler for /api/v1/actual-lines/:id (specific item - only update for now)
	mux.HandleFunc("/api/v1/actual-lines/", func(w http.ResponseWriter, r *http.Request) {
		// ID parsing will be done in handlers
		switch r.Method {
		case http.MethodPut:
			UpdateActualLineHandler(appStore)(w, r)
		default:
			http.Error(w, "Method not allowed for specific actual line item (/api/v1/actual-lines/:id)", http.StatusMethodNotAllowed)
		}
	})

	// Static file serving (same as before)
	fileServer := http.FileServer(http.FS(staticFS))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// ... (static file serving logic as before)
		if strings.HasPrefix(r.URL.Path, "/api/") {
			// This is a fallback for API routes not caught by more specific handlers.
			log.Printf("Unknown or improperly routed API path received by root fallback: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		p := path.Clean(r.URL.Path)
		if p == "/" || p == "/index.html" {
			serveIndexHTML(w, r, staticFS)
			return
		}
		f, err := staticFS.Open(strings.TrimPrefix(p, "/"))
		if err != nil {
			log.Printf("Static file %s not found, serving index.html (SPA fallback)", p)
			serveIndexHTML(w, r, staticFS)
			return
		}
		f.Close()
		fileServer.ServeHTTP(w, r)
	})
	return mux
}

// serveIndexHTML function (same as before)
func serveIndexHTML(w http.ResponseWriter, r *http.Request, staticFS fs.FS) {
	// ... (implementation as before)
	log.Printf("Serving index.html for path: %s", r.URL.Path)
	f, err := staticFS.Open("index.html")
	if err != nil { http.Error(w, "index.html not found", http.StatusInternalServerError); return }
	defer f.Close()
	fi, err_stat := f.Stat(); 
	if err_stat != nil { http.Error(w, "stat failed for index.html", http.StatusInternalServerError); return }
	rs, ok := f.(io.ReadSeeker);
	if !ok { http.Error(w, "seek failed for index.html", http.StatusInternalServerError); return }
	http.ServeContent(w, r, "index.html", fi.ModTime(), rs)
}
