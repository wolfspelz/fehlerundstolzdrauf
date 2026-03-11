package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var dbDir string

var DB *sql.DB

const schema = `
CREATE TABLE IF NOT EXISTS stories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    year TEXT NOT NULL,
    title TEXT NOT NULL,
    text TEXT NOT NULL,
    created_at TEXT DEFAULT (datetime('now')),
    status TEXT DEFAULT 'unmoderated',
    shown_count INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS quotes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    text TEXT NOT NULL,
    attribution TEXT NOT NULL,
    created_at TEXT DEFAULT (datetime('now')),
    shown_count INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS historical (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    year TEXT NOT NULL,
    title TEXT NOT NULL,
    text TEXT NOT NULL,
    created_at TEXT DEFAULT (datetime('now')),
    shown_count INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS edition_cache (
    date TEXT PRIMARY KEY,
    content_json TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS edition_resets (
    date TEXT PRIMARY KEY,
    count INTEGER DEFAULT 1
);
`

func Init(dbPath string) error {
	dbDir = filepath.Dir(dbPath)
	dir := dbDir
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create db dir: %w", err)
	}

	var err error
	DB, err = sql.Open("sqlite3", dbPath+"?_journal_mode=WAL")
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}

	if _, err := DB.Exec(schema); err != nil {
		return fmt.Errorf("create schema: %w", err)
	}

	migrate()

	// Seed if empty
	var count int
	DB.QueryRow("SELECT COUNT(*) FROM stories").Scan(&count)
	if count == 0 {
		seedPaths := []string{"/seed.sql", "internal/db/seed.sql"}
		var seedErr error
		for _, p := range seedPaths {
			if seedErr = runSeedFile(p); seedErr == nil {
				break
			}
		}
		if seedErr != nil {
			return fmt.Errorf("seed db: %w", seedErr)
		}
	}

	return nil
}

func needsMigration() bool {
	if !columnExists("quotes", "created_at") || !columnExists("historical", "created_at") {
		return true
	}
	if tableExists("featured") {
		return true
	}
	return false
}

func migrate() {
	if !needsMigration() {
		return
	}

	if _, err := backupData(dbDir); err != nil {
		log.Printf("Warning: pre-migration backup failed: %v", err)
	}

	// Add created_at to quotes and historical (added when featured was removed)
	for _, table := range []string{"quotes", "historical"} {
		if !columnExists(table, "created_at") {
			DB.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN created_at TEXT DEFAULT (datetime('now'))", table))
		}
	}
	// Drop unused featured table
	DB.Exec("DROP TABLE IF EXISTS featured")
}

func tableExists(table string) bool {
	var name string
	err := DB.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&name)
	return err == nil
}

type BackupResult struct {
	File       string         `json:"file"`
	FileSize   int64          `json:"file_size"`
	TableRows  map[string]int `json:"table_rows"`
	TotalRows  int            `json:"total_rows"`
}

func Backup() (*BackupResult, error) {
	return backupData(dbDir)
}

func backupData(dir string) (*BackupResult, error) {
	ts := time.Now().Format("2006-01-02_150405")
	path := filepath.Join(dir, fmt.Sprintf("backup_%s.sql", ts))

	var b strings.Builder
	b.WriteString("-- Backup " + ts + "\n\n")

	tables := []string{"stories", "quotes", "historical", "featured", "edition_cache", "edition_resets"}
	result := &BackupResult{
		File:      path,
		TableRows: make(map[string]int),
	}

	for _, table := range tables {
		if !tableExists(table) {
			continue
		}
		rows, err := DB.Query(fmt.Sprintf("SELECT * FROM %s", table))
		if err != nil {
			continue
		}
		cols, _ := rows.Columns()
		vals := make([]interface{}, len(cols))
		ptrs := make([]interface{}, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		rowCount := 0
		for rows.Next() {
			rows.Scan(ptrs...)
			b.WriteString(fmt.Sprintf("INSERT INTO %s (%s) VALUES (", table, strings.Join(cols, ", ")))
			for i, v := range vals {
				if i > 0 {
					b.WriteString(", ")
				}
				if v == nil {
					b.WriteString("NULL")
				} else {
					switch val := v.(type) {
					case int64:
						b.WriteString(fmt.Sprintf("%d", val))
					case float64:
						b.WriteString(fmt.Sprintf("%g", val))
					default:
						s := fmt.Sprintf("%s", val)
						s = strings.ReplaceAll(s, "'", "''")
						b.WriteString("'" + s + "'")
					}
				}
			}
			b.WriteString(");\n")
			rowCount++
		}
		rows.Close()
		if rowCount > 0 {
			result.TableRows[table] = rowCount
			result.TotalRows += rowCount
		}
	}

	if err := os.WriteFile(path, []byte(b.String()), 0644); err != nil {
		return nil, fmt.Errorf("write backup: %w", err)
	}

	info, err := os.Stat(path)
	if err == nil {
		result.FileSize = info.Size()
	}

	log.Printf("Database backup written to %s", path)
	return result, nil
}

func columnExists(table, column string) bool {
	rows, err := DB.Query(fmt.Sprintf("PRAGMA table_info(%s)", table))
	if err != nil {
		return false
	}
	defer rows.Close()
	for rows.Next() {
		var cid int
		var name, typ string
		var notnull int
		var dflt sql.NullString
		var pk int
		if rows.Scan(&cid, &name, &typ, &notnull, &dflt, &pk) == nil && name == column {
			return true
		}
	}
	return false
}

func runSeedFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	_, err = DB.Exec(string(data))
	return err
}
