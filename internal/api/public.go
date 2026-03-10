package api

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/wolfspelz/fehlerundstolzdrauf/internal/db"
	"github.com/wolfspelz/fehlerundstolzdrauf/internal/rotation"
)

// Rate limiter
var (
	rateMu    sync.Mutex
	rateMap   = make(map[string]time.Time)
)

func HandleEdition(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	edition, err := rotation.GetEdition(rotation.Today())
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=300")
	json.NewEncoder(w).Encode(edition)
}

type SubmitRequest struct {
	Year  string `json:"year"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

func HandleSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Rate limiting
	ip := clientIP(r)
	rateMu.Lock()
	if last, ok := rateMap[ip]; ok && time.Since(last) < time.Hour {
		rateMu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(map[string]string{"error": "Nur eine Einreichung pro Stunde. Bitte später nochmal versuchen."})
		return
	}
	rateMap[ip] = time.Now()
	rateMu.Unlock()

	var req SubmitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Ungültige Anfrage."})
		return
	}

	// Validate
	req.Year = strings.TrimSpace(req.Year)
	req.Title = strings.TrimSpace(req.Title)
	req.Text = strings.TrimSpace(req.Text)

	if req.Year == "" || req.Title == "" || req.Text == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Jahr, Titel und Text sind Pflichtfelder."})
		return
	}

	if utf8.RuneCountInString(req.Text) > 500 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Text darf maximal 500 Zeichen lang sein."})
		return
	}

	// Strip HTML tags
	req.Text = stripHTML(req.Text)
	req.Title = stripHTML(req.Title)
	req.Year = stripHTML(req.Year)

	_, err := db.DB.Exec("INSERT INTO stories (year, title, text, status) VALUES (?, ?, ?, 'unmoderated')",
		req.Year, req.Title, req.Text)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Fehler beim Speichern."})
		return
	}

	// Invalidate today's cache so new story can appear
	db.DB.Exec("DELETE FROM edition_cache WHERE date = ?", rotation.Today())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Danke für deine Geschichte. Sie ist ab sofort sichtbar."})
}

func stripHTML(s string) string {
	var b strings.Builder
	inTag := false
	for _, r := range s {
		if r == '<' {
			inTag = true
			continue
		}
		if r == '>' {
			inTag = false
			continue
		}
		if !inTag {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func clientIP(r *http.Request) string {
	// Check X-Forwarded-For (behind reverse proxy)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	if xri := r.Header.Get("X-Real-Ip"); xri != "" {
		return xri
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
