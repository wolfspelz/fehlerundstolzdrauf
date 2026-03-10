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

type Featured struct {
	ID        int    `json:"id"`
	YearRange string `json:"year_range"`
	Title     string `json:"title"`
	Intro     string `json:"intro"`
	Quote     string `json:"quote,omitempty"`
	Outro     string `json:"outro,omitempty"`
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

type Edition struct {
	Date        string       `json:"date"`
	Stories     []Story      `json:"stories"`
	Featured    *Featured    `json:"featured"`
	Quotes      []Quote      `json:"quotes"`
	Historical  []Historical `json:"historical"`
	TotalStories int         `json:"total_stories"`
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

	featured, err := pickFeatured(rng)
	if err != nil {
		return nil, fmt.Errorf("pick featured: %w", err)
	}

	quotes, err := pickQuotes(rng, 3)
	if err != nil {
		return nil, fmt.Errorf("pick quotes: %w", err)
	}

	historical, err := pickHistorical(rng, 2)
	if err != nil {
		return nil, fmt.Errorf("pick historical: %w", err)
	}

	var totalStories int
	db.DB.QueryRow("SELECT COUNT(*) FROM stories WHERE status IN ('unmoderated','approved')").Scan(&totalStories)

	// Update shown_count
	for _, s := range stories {
		db.DB.Exec("UPDATE stories SET shown_count = shown_count + 1 WHERE id = ?", s.ID)
	}
	if featured != nil {
		db.DB.Exec("UPDATE featured SET shown_count = shown_count + 1 WHERE id = ?", featured.ID)
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
		Featured:     featured,
		Quotes:       quotes,
		Historical:   historical,
		TotalStories: totalStories,
	}, nil
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

func pickFeatured(rng *rand.Rand) (*Featured, error) {
	rows, err := db.DB.Query("SELECT id, year_range, title, intro, COALESCE(quote,''), COALESCE(outro,'') FROM featured ORDER BY shown_count ASC, id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []Featured
	for rows.Next() {
		var f Featured
		if err := rows.Scan(&f.ID, &f.YearRange, &f.Title, &f.Intro, &f.Quote, &f.Outro); err != nil {
			return nil, err
		}
		all = append(all, f)
	}

	if len(all) == 0 {
		return nil, nil
	}

	// Pick from lowest shown_count
	pool := all
	if len(pool) > 5 {
		pool = pool[:5]
	}
	pick := pool[rng.Intn(len(pool))]
	return &pick, nil
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
