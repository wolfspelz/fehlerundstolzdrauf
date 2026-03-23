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

func HandleAdminSubmissions(w http.ResponseWriter, r *http.Request) {
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

func HandleAdminUpdateSubmission(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func HandleAdminListStories(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, year, title, text, created_at, status, shown_count FROM stories ORDER BY id")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Story struct {
		ID         int    `json:"id"`
		Year       string `json:"year"`
		Title      string `json:"title"`
		Text       string `json:"text"`
		CreatedAt  string `json:"created_at"`
		Status     string `json:"status"`
		ShownCount int    `json:"shown_count"`
	}

	var results []Story
	for rows.Next() {
		var s Story
		rows.Scan(&s.ID, &s.Year, &s.Title, &s.Text, &s.CreatedAt, &s.Status, &s.ShownCount)
		results = append(results, s)
	}
	if results == nil {
		results = []Story{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func HandleAdminListQuotes(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, text, attribution, shown_count FROM quotes ORDER BY id")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Quote struct {
		ID          int    `json:"id"`
		Text        string `json:"text"`
		Attribution string `json:"attribution"`
		ShownCount  int    `json:"shown_count"`
	}

	var results []Quote
	for rows.Next() {
		var q Quote
		rows.Scan(&q.ID, &q.Text, &q.Attribution, &q.ShownCount)
		results = append(results, q)
	}
	if results == nil {
		results = []Quote{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func HandleAdminListHistorical(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, year, title, text, shown_count FROM historical ORDER BY id")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Historical struct {
		ID         int    `json:"id"`
		Year       string `json:"year"`
		Title      string `json:"title"`
		Text       string `json:"text"`
		ShownCount int    `json:"shown_count"`
	}

	var results []Historical
	for rows.Next() {
		var h Historical
		rows.Scan(&h.ID, &h.Year, &h.Title, &h.Text, &h.ShownCount)
		results = append(results, h)
	}
	if results == nil {
		results = []Historical{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func HandleAdminCreateStory(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id, "status": "ok"})
}

func HandleAdminCreateQuote(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id, "status": "ok"})
}

func HandleAdminCreateHistorical(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id, "status": "ok"})
}

func HandleAdminNewEdition(w http.ResponseWriter, r *http.Request) {
	if err := rotation.ResetEdition(rotation.Today()); err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func HandleAdminStats(w http.ResponseWriter, r *http.Request) {
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

func HandleAdminBackup(w http.ResponseWriter, r *http.Request) {
	result, err := db.Backup()
	if err != nil {
		http.Error(w, "Backup failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func HandleAdminUpdate(w http.ResponseWriter, r *http.Request) {
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

	// Whitelist tables and their updatable fields
	type tableConfig struct {
		fields []string
	}
	validTables := map[string]tableConfig{
		"stories":    {fields: []string{"year", "title", "text"}},
		"quotes":     {fields: []string{"text", "attribution"}},
		"historical": {fields: []string{"year", "title", "text"}},
	}

	config, ok := validTables[table]
	if !ok {
		http.Error(w, "Invalid type", http.StatusBadRequest)
		return
	}

	var raw map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Build SET clause from provided fields
	var setClauses []string
	var args []interface{}
	for _, field := range config.fields {
		if val, exists := raw[field]; exists {
			setClauses = append(setClauses, field+" = ?")
			args = append(args, val)
		}
	}

	if len(setClauses) == 0 {
		http.Error(w, "No valid fields to update", http.StatusBadRequest)
		return
	}

	args = append(args, id)
	query := "UPDATE " + table + " SET " + strings.Join(setClauses, ", ") + " WHERE id = ?"
	result, err := db.DB.Exec(query, args...)
	if err != nil {
		log.Printf("Update error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	rotation.UpdateCachedEntry(rotation.Today(), table, id, raw)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func HandleAdminDelete(w http.ResponseWriter, r *http.Request) {
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

	rotation.RemoveCachedEntry(rotation.Today(), table, id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
