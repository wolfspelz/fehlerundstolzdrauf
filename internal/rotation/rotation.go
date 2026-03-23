package rotation

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/wolfspelz/fehlerundstolzdrauf/internal/db"
)

type Story struct {
	ID    int    `json:"id"`
	Year  string `json:"year"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

type Quote struct {
	ID          int    `json:"id"`
	Text        string `json:"text"`
	Attribution string `json:"attribution"`
}

type Historical struct {
	ID    int    `json:"id"`
	Year  string `json:"year"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

type NewEntry struct {
	Type        string `json:"type"`
	ID          int    `json:"id"`
	Year        string `json:"year,omitempty"`
	Title       string `json:"title,omitempty"`
	Text        string `json:"text,omitempty"`
	Attribution string `json:"attribution,omitempty"`
}

type Edition struct {
	Date            string       `json:"date"`
	Stories         []Story      `json:"stories"`
	New             []NewEntry   `json:"new"`
	Older           []NewEntry   `json:"older"`
	Quotes          []Quote      `json:"quotes"`
	Historical      []Historical `json:"historical"`
}

func dateToSeed(date string) int64 {
	var resetCount int
	db.DB.QueryRow("SELECT count FROM edition_resets WHERE date = ?", date).Scan(&resetCount)
	input := date
	if resetCount > 0 {
		input = date + "-" + strconv.Itoa(resetCount)
	}
	h := sha256.Sum256([]byte(input))
	return int64(binary.BigEndian.Uint64(h[:8]))
}

func ResetEdition(date string) error {
	_, err := db.DB.Exec("INSERT INTO edition_resets (date, count) VALUES (?, 1) ON CONFLICT(date) DO UPDATE SET count = count + 1", date)
	if err != nil {
		return err
	}
	_, err = db.DB.Exec("DELETE FROM edition_cache WHERE date = ?", date)
	return err
}

func GetEdition(date string) (*Edition, error) {
	// Check cache
	var cached string
	err := db.DB.QueryRow("SELECT content_json FROM edition_cache WHERE date = ?", date).Scan(&cached)
	if err == nil {
		var ed Edition
		if json.Unmarshal([]byte(cached), &ed) == nil {
			return &ed, nil
		}
	}

	// Generate new edition
	ed, err := generateEdition(date)
	if err != nil {
		return nil, err
	}

	// Cache it
	data, _ := json.Marshal(ed)
	db.DB.Exec("INSERT OR REPLACE INTO edition_cache (date, content_json) VALUES (?, ?)", date, string(data))

	return ed, nil
}

func generateEdition(date string) (*Edition, error) {
	rng := rand.New(rand.NewSource(dateToSeed(date)))

	stories, err := pickStories(rng, 5)
	if err != nil {
		return nil, fmt.Errorf("pick stories: %w", err)
	}

	quotes, err := pickQuotes(rng, 3)
	if err != nil {
		return nil, fmt.Errorf("pick quotes: %w", err)
	}

	historical, err := pickHistorical(rng, 2)
	if err != nil {
		return nil, fmt.Errorf("pick historical: %w", err)
	}

	// Collect IDs to exclude from "Neu" column
	excludeStoryIDs := make(map[int]bool)
	for _, s := range stories {
		excludeStoryIDs[s.ID] = true
	}
	excludeQuoteIDs := make(map[int]bool)
	for _, q := range quotes {
		excludeQuoteIDs[q.ID] = true
	}
	excludeHistoricalIDs := make(map[int]bool)
	for _, h := range historical {
		excludeHistoricalIDs[h.ID] = true
	}

	newEntries, err := pickNew(rng, excludeStoryIDs, excludeQuoteIDs, excludeHistoricalIDs, 4)
	if err != nil {
		return nil, fmt.Errorf("pick new: %w", err)
	}

	// Collect IDs from new entries to also exclude from older
	for _, n := range newEntries {
		switch n.Type {
		case "story":
			excludeStoryIDs[n.ID] = true
		case "quote":
			excludeQuoteIDs[n.ID] = true
		case "historical":
			excludeHistoricalIDs[n.ID] = true
		}
	}

	olderEntries, err := pickOlder(rng, excludeStoryIDs, excludeQuoteIDs, excludeHistoricalIDs, 4-len(newEntries))
	if err != nil {
		return nil, fmt.Errorf("pick older: %w", err)
	}

	// Update shown_count
	for _, s := range stories {
		db.DB.Exec("UPDATE stories SET shown_count = shown_count + 1 WHERE id = ?", s.ID)
	}
	for _, q := range quotes {
		db.DB.Exec("UPDATE quotes SET shown_count = shown_count + 1 WHERE id = ?", q.ID)
	}
	for _, h := range historical {
		db.DB.Exec("UPDATE historical SET shown_count = shown_count + 1 WHERE id = ?", h.ID)
	}

	return &Edition{
		Date:         date,
		Stories:      stories,
		New:          newEntries,
		Older:        olderEntries,
		Quotes:       quotes,
		Historical:   historical,
	}, nil
}

func pickNew(rng *rand.Rand, excludeStoryIDs, excludeQuoteIDs, excludeHistoricalIDs map[int]bool, max int) ([]NewEntry, error) {
	var entries []NewEntry

	// Recent stories
	rows, err := db.DB.Query("SELECT id, year, title, text FROM stories WHERE status IN ('unmoderated','approved') AND created_at >= date('now', '-1 month') ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var year, title, text string
		if err := rows.Scan(&id, &year, &title, &text); err != nil {
			return nil, err
		}
		if !excludeStoryIDs[id] {
			entries = append(entries, NewEntry{Type: "story", ID: id, Year: year, Title: title, Text: text})
		}
	}

	// Recent quotes
	rows2, err := db.DB.Query("SELECT id, text, attribution FROM quotes WHERE created_at >= date('now', '-1 month') ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows2.Close()
	for rows2.Next() {
		var id int
		var text, attribution string
		if err := rows2.Scan(&id, &text, &attribution); err != nil {
			return nil, err
		}
		if !excludeQuoteIDs[id] {
			entries = append(entries, NewEntry{Type: "quote", ID: id, Text: text, Attribution: attribution})
		}
	}

	// Recent historical
	rows3, err := db.DB.Query("SELECT id, year, title, text FROM historical WHERE created_at >= date('now', '-1 month') ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows3.Close()
	for rows3.Next() {
		var id int
		var year, title, text string
		if err := rows3.Scan(&id, &year, &title, &text); err != nil {
			return nil, err
		}
		if !excludeHistoricalIDs[id] {
			entries = append(entries, NewEntry{Type: "historical", ID: id, Year: year, Title: title, Text: text})
		}
	}

	// Shuffle and limit
	rng.Shuffle(len(entries), func(i, j int) { entries[i], entries[j] = entries[j], entries[i] })
	if len(entries) > max {
		entries = entries[:max]
	}

	return entries, nil
}

func pickOlder(rng *rand.Rand, excludeStoryIDs, excludeQuoteIDs, excludeHistoricalIDs map[int]bool, max int) ([]NewEntry, error) {
	if max <= 0 {
		return nil, nil
	}

	var entries []NewEntry

	// Older stories (> 1 month or NULL created_at)
	rows, err := db.DB.Query("SELECT id, year, title, text FROM stories WHERE status IN ('unmoderated','approved') AND (created_at < date('now', '-1 month') OR created_at IS NULL) ORDER BY shown_count ASC, id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var year, title, text string
		if err := rows.Scan(&id, &year, &title, &text); err != nil {
			return nil, err
		}
		if !excludeStoryIDs[id] {
			entries = append(entries, NewEntry{Type: "story", ID: id, Year: year, Title: title, Text: text})
		}
	}

	// Older quotes
	rows2, err := db.DB.Query("SELECT id, text, attribution FROM quotes WHERE created_at < date('now', '-1 month') OR created_at IS NULL ORDER BY shown_count ASC, id ASC")
	if err != nil {
		return nil, err
	}
	defer rows2.Close()
	for rows2.Next() {
		var id int
		var text, attribution string
		if err := rows2.Scan(&id, &text, &attribution); err != nil {
			return nil, err
		}
		if !excludeQuoteIDs[id] {
			entries = append(entries, NewEntry{Type: "quote", ID: id, Text: text, Attribution: attribution})
		}
	}

	// Older historical
	rows3, err := db.DB.Query("SELECT id, year, title, text FROM historical WHERE created_at < date('now', '-1 month') OR created_at IS NULL ORDER BY shown_count ASC, id ASC")
	if err != nil {
		return nil, err
	}
	defer rows3.Close()
	for rows3.Next() {
		var id int
		var year, title, text string
		if err := rows3.Scan(&id, &year, &title, &text); err != nil {
			return nil, err
		}
		if !excludeHistoricalIDs[id] {
			entries = append(entries, NewEntry{Type: "historical", ID: id, Year: year, Title: title, Text: text})
		}
	}

	rng.Shuffle(len(entries), func(i, j int) { entries[i], entries[j] = entries[j], entries[i] })
	if len(entries) > max {
		entries = entries[:max]
	}

	return entries, nil
}

func pickStories(rng *rand.Rand, n int) ([]Story, error) {
	rows, err := db.DB.Query("SELECT id, year, title, text FROM stories WHERE status IN ('unmoderated','approved') ORDER BY shown_count ASC, id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []Story
	for rows.Next() {
		var s Story
		if err := rows.Scan(&s.ID, &s.Year, &s.Title, &s.Text); err != nil {
			return nil, err
		}
		all = append(all, s)
	}

	if len(all) == 0 {
		return nil, nil
	}

	// Take from the lowest shown_count pool, then shuffle
	if len(all) > n*3 {
		all = all[:n*3]
	}
	rng.Shuffle(len(all), func(i, j int) { all[i], all[j] = all[j], all[i] })

	if len(all) > n {
		all = all[:n]
	}
	return all, nil
}

func pickQuotes(rng *rand.Rand, n int) ([]Quote, error) {
	rows, err := db.DB.Query("SELECT id, text, attribution FROM quotes ORDER BY shown_count ASC, id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []Quote
	for rows.Next() {
		var q Quote
		if err := rows.Scan(&q.ID, &q.Text, &q.Attribution); err != nil {
			return nil, err
		}
		all = append(all, q)
	}

	if len(all) == 0 {
		return nil, nil
	}

	if len(all) > n*3 {
		all = all[:n*3]
	}
	rng.Shuffle(len(all), func(i, j int) { all[i], all[j] = all[j], all[i] })

	if len(all) > n {
		all = all[:n]
	}
	return all, nil
}

func pickHistorical(rng *rand.Rand, n int) ([]Historical, error) {
	rows, err := db.DB.Query("SELECT id, year, title, text FROM historical ORDER BY shown_count ASC, id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []Historical
	for rows.Next() {
		var h Historical
		if err := rows.Scan(&h.ID, &h.Year, &h.Title, &h.Text); err != nil {
			return nil, err
		}
		all = append(all, h)
	}

	if len(all) == 0 {
		return nil, nil
	}

	if len(all) > n*3 {
		all = all[:n*3]
	}
	rng.Shuffle(len(all), func(i, j int) { all[i], all[j] = all[j], all[i] })

	if len(all) > n {
		all = all[:n]
	}
	return all, nil
}

func Today() string {
	return time.Now().UTC().Format("2006-01-02")
}

// tableToEntryType maps DB table names to NewEntry.Type values
func tableToEntryType(table string) string {
	switch table {
	case "stories":
		return "story"
	case "quotes":
		return "quote"
	default:
		return table // "historical" stays "historical"
	}
}

// UpdateCachedEntry updates a specific entry in the cached edition JSON without regenerating.
func UpdateCachedEntry(date string, table string, id int, fields map[string]interface{}) {
	var cached string
	if db.DB.QueryRow("SELECT content_json FROM edition_cache WHERE date = ?", date).Scan(&cached) != nil {
		return
	}

	var ed Edition
	if json.Unmarshal([]byte(cached), &ed) != nil {
		return
	}

	updated := false

	switch table {
	case "stories":
		for i := range ed.Stories {
			if ed.Stories[i].ID == id {
				if v, ok := fields["year"].(string); ok {
					ed.Stories[i].Year = v
				}
				if v, ok := fields["title"].(string); ok {
					ed.Stories[i].Title = v
				}
				if v, ok := fields["text"].(string); ok {
					ed.Stories[i].Text = v
				}
				updated = true
			}
		}
	case "quotes":
		for i := range ed.Quotes {
			if ed.Quotes[i].ID == id {
				if v, ok := fields["text"].(string); ok {
					ed.Quotes[i].Text = v
				}
				if v, ok := fields["attribution"].(string); ok {
					ed.Quotes[i].Attribution = v
				}
				updated = true
			}
		}
	case "historical":
		for i := range ed.Historical {
			if ed.Historical[i].ID == id {
				if v, ok := fields["year"].(string); ok {
					ed.Historical[i].Year = v
				}
				if v, ok := fields["title"].(string); ok {
					ed.Historical[i].Title = v
				}
				if v, ok := fields["text"].(string); ok {
					ed.Historical[i].Text = v
				}
				updated = true
			}
		}
	}

	// Also update in New and Older sections
	entryType := tableToEntryType(table)
	for i := range ed.New {
		if ed.New[i].Type == entryType && ed.New[i].ID == id {
			if v, ok := fields["year"].(string); ok {
				ed.New[i].Year = v
			}
			if v, ok := fields["title"].(string); ok {
				ed.New[i].Title = v
			}
			if v, ok := fields["text"].(string); ok {
				ed.New[i].Text = v
			}
			if v, ok := fields["attribution"].(string); ok {
				ed.New[i].Attribution = v
			}
			updated = true
		}
	}
	for i := range ed.Older {
		if ed.Older[i].Type == entryType && ed.Older[i].ID == id {
			if v, ok := fields["year"].(string); ok {
				ed.Older[i].Year = v
			}
			if v, ok := fields["title"].(string); ok {
				ed.Older[i].Title = v
			}
			if v, ok := fields["text"].(string); ok {
				ed.Older[i].Text = v
			}
			if v, ok := fields["attribution"].(string); ok {
				ed.Older[i].Attribution = v
			}
			updated = true
		}
	}

	if updated {
		data, _ := json.Marshal(ed)
		db.DB.Exec("UPDATE edition_cache SET content_json = ? WHERE date = ?", string(data), date)
	}
}

// RemoveCachedEntry removes a specific entry from the cached edition JSON.
func RemoveCachedEntry(date string, table string, id int) {
	var cached string
	if db.DB.QueryRow("SELECT content_json FROM edition_cache WHERE date = ?", date).Scan(&cached) != nil {
		return
	}

	var ed Edition
	if json.Unmarshal([]byte(cached), &ed) != nil {
		return
	}

	switch table {
	case "stories":
		for i := range ed.Stories {
			if ed.Stories[i].ID == id {
				ed.Stories = append(ed.Stories[:i], ed.Stories[i+1:]...)
				break
			}
		}
	case "quotes":
		for i := range ed.Quotes {
			if ed.Quotes[i].ID == id {
				ed.Quotes = append(ed.Quotes[:i], ed.Quotes[i+1:]...)
				break
			}
		}
	case "historical":
		for i := range ed.Historical {
			if ed.Historical[i].ID == id {
				ed.Historical = append(ed.Historical[:i], ed.Historical[i+1:]...)
				break
			}
		}
	}

	entryType := tableToEntryType(table)
	for i := range ed.New {
		if ed.New[i].Type == entryType && ed.New[i].ID == id {
			ed.New = append(ed.New[:i], ed.New[i+1:]...)
			break
		}
	}
	for i := range ed.Older {
		if ed.Older[i].Type == entryType && ed.Older[i].ID == id {
			ed.Older = append(ed.Older[:i], ed.Older[i+1:]...)
			break
		}
	}

	data, _ := json.Marshal(ed)
	db.DB.Exec("UPDATE edition_cache SET content_json = ? WHERE date = ?", string(data), date)
}
