package storage

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Storage manages local SQLite database for Community Station.
type Storage struct {
	db  *sql.DB
	dir string
}

// New creates a new Storage instance.
func New(dir string) (*Storage, error) {
	if dir == "" {
		return nil, fmt.Errorf("directory required")
	}

	dbPath := filepath.Join(dir, "community_station.db")

	db, err := sql.Open("sqlite3", dbPath+"?_journal=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	s := &Storage{
		db:  db,
		dir: dir,
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
	}

	for _, sqlStr := range schemas {
		if _, err := s.db.Exec(sqlStr); err != nil {
			return fmt.Errorf("creating table/index: %w", err)
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
func (s *Storage) InsertAircraft(event interface{}) error {
	// This would be fully implemented with proper type conversion
	// For now, return nil (stub for integration)
	return nil
}

// GetAircraft gets aircraft events with pagination.
func (s *Storage) GetAircraft(limit int, offset int) ([]map[string]interface{}, error) {
	rows, err := s.db.Query(`SELECT id, hex, flight, lat, lon, alt_baro, gs, track, seen, source_id, source_type, observed_at, verified, synthetic FROM aircraft_events ORDER BY observed_at DESC LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []map[string]interface{}
	for rows.Next() {
		var id, hex, flight, sourceID, sourceType string
		var lat, lon, gs, track, seen float64
		var altBaro int
		var observedAt time.Time
		var verified, synthetic int

		err := rows.Scan(&id, &hex, &flight, &lat, &lon, &altBaro, &gs, &track, &seen, &sourceID, &sourceType, &observedAt, &verified, &synthetic)
		if err != nil {
			continue
		}

		events = append(events, map[string]interface{}{
			"id": id,
			"hex": hex,
			"flight": flight,
			"lat": lat,
			"lon": lon,
			"alt_baro": altBaro,
			"gs": gs,
			"track": track,
			"seen": seen,
			"source_id": sourceID,
			"source_type": sourceType,
			"observed_at": observedAt,
			"verified": verified == 1,
			"synthetic": synthetic == 1,
		})
	}
	return events, rows.Err()
}

// AircraftEvent represents an aircraft observation.
type AircraftEvent struct {
	ID            string     `json:"id"`
	Hex           string     `json:"hex"`
	Flight        string     `json:"flight"`
	Lat           float64    `json:"lat"`
	Lon           float64    `json:"lon"`
	AltBaro       int        `json:"alt_baro"`
	GS            float64    `json:"gs"`
	Track         float64    `json:"track"`
	Seen          float64    `json:"seen"`
	SourceID      string     `json:"source_id"`
	SourceType    string     `json:"source_type"`
	ObservedAt    time.Time `json:"observed_at"`
	Verified      int        `json:"verified"`
	Synthetic     int        `json:"synthetic"`
}