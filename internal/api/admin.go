package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/wolfspelz/fehlerundstolzdrauf/internal/db"
	"github.com/wolfspelz/fehlerundstolzdrauf/internal/rotation"
)

var AdminToken string

func AdminAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") || strings.TrimPrefix(auth, "Bearer ") != AdminToken {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func HandleSubmissions(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	if status == "" {
		status = "unmoderated"
	}

	rows, err := db.DB.Query("SELECT id, year, title, text, created_at, status, shown_count FROM stories WHERE status = ? ORDER BY created_at DESC", status)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type StoryAdmin struct {
		ID         int    `json:"id"`
		Year       string `json:"year"`
		Title      string `json:"title"`
		Text       string `json:"text"`
		CreatedAt  string `json:"created_at"`
		Status     string `json:"status"`
		ShownCount int    `json:"shown_count"`
	}

	var results []StoryAdmin
	for rows.Next() {
		var s StoryAdmin
		rows.Scan(&s.ID, &s.Year, &s.Title, &s.Text, &s.CreatedAt, &s.Status, &s.ShownCount)
		results = append(results, s)
	}

	if results == nil {
		results = []StoryAdmin{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func HandleUpdateSubmission(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path: /admin/submissions/123
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/admin/submissions/"), "/")
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Status != "approved" && req.Status != "hidden" && req.Status != "unmoderated" {
		http.Error(w, "Status must be 'approved', 'hidden', or 'unmoderated'", http.StatusBadRequest)
		return
	}

	result, err := db.DB.Exec("UPDATE stories SET status = ? WHERE id = ?", req.Status, id)
	if err != nil {
		log.Printf("UpdateSubmission error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// Invalidate cache
	db.DB.Exec("DELETE FROM edition_cache WHERE date = ?", rotation.Today())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func HandleCreateStory(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Year  string `json:"year"`
		Title string `json:"title"`
		Text  string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	result, err := db.DB.Exec("INSERT INTO stories (year, title, text, status) VALUES (?, ?, ?, 'approved')",
		req.Year, req.Title, req.Text)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	db.DB.Exec("DELETE FROM edition_cache WHERE date = ?", rotation.Today())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id, "status": "ok"})
}

func HandleCreateQuote(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Text        string `json:"text"`
		Attribution string `json:"attribution"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	result, err := db.DB.Exec("INSERT INTO quotes (text, attribution) VALUES (?, ?)", req.Text, req.Attribution)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	db.DB.Exec("DELETE FROM edition_cache WHERE date = ?", rotation.Today())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id, "status": "ok"})
}

func HandleCreateHistorical(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Year  string `json:"year"`
		Title string `json:"title"`
		Text  string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	result, err := db.DB.Exec("INSERT INTO historical (year, title, text) VALUES (?, ?, ?)",
		req.Year, req.Title, req.Text)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	db.DB.Exec("DELETE FROM edition_cache WHERE date = ?", rotation.Today())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id, "status": "ok"})
}

func HandleResetEdition(w http.ResponseWriter, r *http.Request) {
	if err := rotation.ResetEdition(rotation.Today()); err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func HandleStats(w http.ResponseWriter, r *http.Request) {
	stats := make(map[string]interface{})

	var total, approved, unmoderated, hidden int
	db.DB.QueryRow("SELECT COUNT(*) FROM stories").Scan(&total)
	db.DB.QueryRow("SELECT COUNT(*) FROM stories WHERE status='approved'").Scan(&approved)
	db.DB.QueryRow("SELECT COUNT(*) FROM stories WHERE status='unmoderated'").Scan(&unmoderated)
	db.DB.QueryRow("SELECT COUNT(*) FROM stories WHERE status='hidden'").Scan(&hidden)

	var quotesCount, historicalCount int
	db.DB.QueryRow("SELECT COUNT(*) FROM quotes").Scan(&quotesCount)
	db.DB.QueryRow("SELECT COUNT(*) FROM historical").Scan(&historicalCount)

	stats["stories_total"] = total
	stats["stories_approved"] = approved
	stats["stories_unmoderated"] = unmoderated
	stats["stories_hidden"] = hidden
	stats["quotes"] = quotesCount
	stats["historical"] = historicalCount

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func HandleDelete(w http.ResponseWriter, r *http.Request) {
	// Path: /admin/{type}/{id}
	path := strings.TrimPrefix(r.URL.Path, "/admin/")
	parts := strings.Split(path, "/")
	if len(parts) != 2 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	table := parts[0]
	id, err := strconv.Atoi(parts[1])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Whitelist tables
	validTables := map[string]bool{"stories": true, "quotes": true, "historical": true}
	if !validTables[table] {
		http.Error(w, "Invalid type", http.StatusBadRequest)
		return
	}

	result, err := db.DB.Exec("DELETE FROM "+table+" WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	db.DB.Exec("DELETE FROM edition_cache WHERE date = ?", rotation.Today())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
