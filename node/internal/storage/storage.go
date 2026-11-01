package storage

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Storage manages local SQLite database.
type Storage struct {
	db   *sql.DB
	dir  string
	path string
}

// New creates a new Storage instance.
func New(dir string) (*Storage, error) {
	if dir == "" {
		return nil, fmt.Errorf("directory required")
	}

	dbPath := filepath.Join(dir, "nieboguard.db")

	db, err := sql.Open("sqlite3", dbPath+"?_journal=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	s := &Storage{
		db:   db,
		dir:  dir,
		path: dbPath,
	}

	return s, nil
}

// Initialize creates required tables and indexes.
func (s *Storage) Initialize() error {
	schemas := []string{
		`CREATE TABLE IF NOT EXISTS aircraft_events (
			id TEXT PRIMARY KEY,
			hex TEXT NOT NULL,
			flight TEXT,
			lat REAL,
			lon REAL,
			alt_baro INTEGER,
			gs REAL,
			track REAL,
			seen REAL,
			source_id TEXT,
			source_type TEXT,
			observed_at DATETIME,
			received_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			verified INTEGER DEFAULT 0,
			synthetic INTEGER DEFAULT 0
		);`,
		`CREATE INDEX IF NOT EXISTS idx_aircraft_hex ON aircraft_events(hex);`,
		`CREATE INDEX IF NOT EXISTS idx_aircraft_source ON aircraft_events(source_id);`,
		`CREATE INDEX IF NOT EXISTS idx_aircraft_time ON aircraft_events(observed_at);`,

		`CREATE TABLE IF NOT EXISTS notices (
			id TEXT PRIMARY KEY,
			source TEXT,
			category TEXT,
			title TEXT,
			region TEXT,
			verification TEXT,
			published_at DATETIME,
			received_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,

		`CREATE TABLE IF NOT EXISTS sources (
			id TEXT PRIMARY KEY,
			type TEXT,
			name TEXT,
			url TEXT,
			status TEXT DEFAULT 'inactive',
			last_update DATETIME,
			license TEXT,
			count_records INTEGER DEFAULT 0
		);`,

		`CREATE TABLE IF NOT EXISTS configuration (
			key TEXT PRIMARY KEY,
			value BLOB
		);`,

		`CREATE INDEX IF NOT EXISTS idx_notices_time ON notices(published_at);`,
	}

	for _, sql := range schemas {
		if _, err := s.db.Exec(sql); err != nil {
			return fmt.Errorf("creating table: %w", err)
		}
	}

	// Enable WAL mode explicitly
	_, err := s.db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		return fmt.Errorf("enabling WAL: %w", err)
	}

	return nil
}

// Close closes the database connection.
func (s *Storage) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// InsertAircraft inserts an aircraft event.
func (s *Storage) InsertAircraft(event *AircraftEvent) error {
	_, err := s.db.Exec(`INSERT OR IGNORE INTO aircraft_events 
		(id, hex, flight, lat, lon, alt_baro, gs, track, seen, source_id, source_type, observed_at, verified, synthetic) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		event.ID,
		event.Hex,
		event.Flight,
		event.Lat,
		event.Lon,
		event.AltBaro,
		event.GS,
		event.Track,
		event.Seen,
		event.SourceID,
		event.SourceType,
		event.ObservedAt,
		event.Verified,
		event.Synthetic,
	)
	return err
}

// GetAircraft gets aircraft events with pagination.
func (s *Storage) GetAircraft(limit int, offset int) ([]AircraftEvent, error) {
	rows, err := s.db.Query(`SELECT id, hex, flight, lat, lon, alt_baro, gs, track, seen, source_id, source_type, observed_at, verified, synthetic FROM aircraft_events ORDER BY observed_at DESC LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []AircraftEvent
	for rows.Next() {
		var e AircraftEvent
		err := rows.Scan(&e.ID, &e.Hex, &e.Flight, &e.Lat, &e.Lon, &e.AltBaro, &e.GS, &e.Track, &e.Seen, &e.SourceID, &e.SourceType, &e.ObservedAt, &e.Verified, &e.Synthetic)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, rows.Err()
}

// AircraftEvent represents an aircraft observation.
type AircraftEvent struct {
	ID        string     `json:"id"`
	Hex       string     `json:"hex"`
	Flight    string     `json:"flight"`
	Lat       float64    `json:"lat"`
	Lon       float64    `json:"lon"`
	AltBaro   int        `json:"alt_baro"`
	GS        float64    `json:"gs"`
	Track     float64    `json:"track"`
	Seen      float64    `json:"seen"`
	SourceID  string     `json:"source_id"`
	SourceType string    `json:"source_type"`
	ObservedAt time.Time `json:"observed_at"`
	Verified  int        `json:"verified"`
	Synthetic int        `json:"synthetic"`
}
