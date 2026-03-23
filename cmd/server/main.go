package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/wolfspelz/fehlerundstolzdrauf/internal/api"
	"github.com/wolfspelz/fehlerundstolzdrauf/internal/db"
)

func main() {
	dbPath := "/data/fehlerundstolzdrauf.db"
	if v := os.Getenv("DB_PATH"); v != "" {
		dbPath = v
	}

	api.AdminToken = os.Getenv("ADMIN_TOKEN")
	if api.AdminToken == "" {
		log.Fatal("ADMIN_TOKEN environment variable is required")
	}

	if err := db.Init(dbPath); err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}
	defer db.DB.Close()

	// Public API
	http.HandleFunc("/api/edition", api.HandleEdition)
	http.HandleFunc("/api/submit", api.HandleSubmit)
	http.HandleFunc("/api/stats", api.HandlePublicStats)

	// Admin API
	http.HandleFunc("/admin/submissions", api.AdminAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			api.HandleSubmissions(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))
	http.HandleFunc("/admin/submissions/", api.AdminAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			api.HandleUpdateSubmission(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))
	http.HandleFunc("/admin/stories", api.AdminAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			api.HandleListStories(w, r)
		} else if r.Method == http.MethodPost {
			api.HandleCreateStory(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))
	http.HandleFunc("/admin/quotes", api.AdminAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			api.HandleListQuotes(w, r)
		} else if r.Method == http.MethodPost {
			api.HandleCreateQuote(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))
	http.HandleFunc("/admin/historical", api.AdminAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			api.HandleListHistorical(w, r)
		} else if r.Method == http.MethodPost {
			api.HandleCreateHistorical(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))
	http.HandleFunc("/admin/new-edition", api.AdminAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			api.HandleNewEdition(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))
	http.HandleFunc("/admin/backup", api.AdminAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			api.HandleBackup(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))
	http.HandleFunc("/admin/stats", api.AdminAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			api.HandleStats(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	// Admin delete: /admin/{type}/{id}
	http.HandleFunc("/admin/", api.AdminAuth(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/admin/")
		parts := strings.Split(path, "/")

		// Handle PUT and DELETE for /{type}/{id} pattern
		if len(parts) == 2 {
			if r.Method == http.MethodPut {
				api.HandleUpdate(w, r)
				return
			}
			if r.Method == http.MethodDelete {
				api.HandleDelete(w, r)
				return
			}
		}

		// Don't interfere with other /admin/ routes already registered
		http.NotFound(w, r)
	}))

	// Static files
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)

	port := "80"
	if v := os.Getenv("PORT"); v != "" {
		port = v
	}

	fmt.Printf("Server starting on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
