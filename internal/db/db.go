package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

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

CREATE TABLE IF NOT EXISTS featured (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    year_range TEXT NOT NULL,
    title TEXT NOT NULL,
    intro TEXT NOT NULL,
    quote TEXT,
    outro TEXT,
    created_at TEXT DEFAULT (datetime('now')),
    shown_count INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS quotes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    text TEXT NOT NULL,
    attribution TEXT NOT NULL,
    shown_count INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS historical (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    year TEXT NOT NULL,
    title TEXT NOT NULL,
    text TEXT NOT NULL,
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
	dir := filepath.Dir(dbPath)
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

	// Seed if empty
	var count int
	DB.QueryRow("SELECT COUNT(*) FROM stories").Scan(&count)
	if count == 0 {
		if err := runSeedFile("/data/seed.sql"); err != nil {
			// Try relative path for local dev
			if err2 := runSeedFile("data/seed.sql"); err2 != nil {
				return fmt.Errorf("seed db: %w (also tried relative: %v)", err, err2)
			}
		}
	}

	return nil
}

func runSeedFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	_, err = DB.Exec(string(data))
	return err
}
