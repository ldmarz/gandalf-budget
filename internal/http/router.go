package http

import (
	"io"
	"io/fs"
	"log"
	"net/http"
	"path"
	"strings" // Required for path manipulation

	"github.com/jmoiron/sqlx"

	"gandalf-budget/internal/store"
)

func NewRouter(staticFS fs.FS, db *sqlx.DB) *http.ServeMux {
	mux := http.NewServeMux()
	appStore := store.NewSQLStore(db)

	mux.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"status": "ok"}`)
		log.Println("HIT: /api/v1/health")
	})

	mux.HandleFunc("/api/v1/dashboard", GetDashboardData(appStore))
	mux.HandleFunc("/api/v1/export/json", ExportJSONHandler(appStore)) // New route
	mux.HandleFunc("/api/v1/reports/annual", GetAnnualReport(appStore))

	mux.HandleFunc("/api/v1/reports/snapshots/", GetSnapshotDetail(appStore))
	
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

	mux.HandleFunc("/api/v1/categories/", func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/categories/")
		idStr = strings.TrimSuffix(idStr, "/")

		if idStr == "" {
			http.Error(w, "Category ID missing", http.StatusBadRequest)
			return
		}
		
		switch r.Method {
		case http.MethodPut:
			HandleUpdateCategory(appStore)(w, r)
		case http.MethodDelete:
			HandleDeleteCategory(appStore)(w, r)
		default:
			http.Error(w, "Method not allowed for specific category item (/api/v1/categories/:id)", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/v1/budget-lines", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			CreateBudgetLineHandler(appStore)(w, r)
		case http.MethodGet:
			GetBudgetLinesByMonthIDHandler(appStore)(w, r)
		default:
			http.Error(w, "Method not allowed for /api/v1/budget-lines collection", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/v1/budget-lines/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			UpdateBudgetLineHandler(appStore)(w, r)
		case http.MethodDelete:
			DeleteBudgetLineHandler(appStore)(w, r)
		default:
			http.Error(w, "Method not allowed for specific budget line item (/api/v1/budget-lines/:id)", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/v1/actual-lines/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			UpdateActualLineHandler(appStore)(w, r)
		default:
			http.Error(w, "Method not allowed for specific actual line item (/api/v1/actual-lines/:id)", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/v1/months/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/months/"), "/")
			if len(pathParts) >= 2 && pathParts[1] == "finalize" {
				FinalizeMonthHandler(appStore)(w, r)
			} else {
				http.NotFound(w,r)
			}
		} else {
			http.Error(w, "Method not allowed for /api/v1/months/", http.StatusMethodNotAllowed)
		}
	})

	fileServer := http.FileServer(http.FS(staticFS))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			if !isKnownAPIPath(r.URL.Path) {
				log.Printf("Unknown or improperly routed API path received by root fallback: %s", r.URL.Path)
				http.NotFound(w, r)
				return
			}
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

func serveIndexHTML(w http.ResponseWriter, r *http.Request, staticFS fs.FS) {
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

func isKnownAPIPath(path string) bool {
	knownAPIPrefixes := []string{
		"/api/v1/health",
		"/api/v1/dashboard",
		"/api/v1/reports/annual",
		"/api/v1/reports/snapshots/",
		"/api/v1/categories",
		"/api/v1/budget-lines",
		"/api/v1/actual-lines/",
		"/api/v1/months/",
		"/api/v1/export/json", // Add this line
	}
	for _, prefix := range knownAPIPrefixes {
		if strings.HasPrefix(path, prefix) {
			if (prefix == "/api/v1/categories" && path != "/api/v1/categories" && !strings.HasPrefix(path, "/api/v1/categories/")) ||
			   (prefix == "/api/v1/budget-lines" && path != "/api/v1/budget-lines" && !strings.HasPrefix(path, "/api/v1/budget-lines/")) {
				continue
			}
			return true
		}
	}
	return false
}
