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
		`CREATE TABLE IF NOT EXISTS offline_packs (
			id TEXT NOT NULL,
			version TEXT NOT NULL,
			sha256 TEXT NOT NULL,
			bytes INTEGER NOT NULL DEFAULT 0,
			created_at TEXT NOT NULL,
			manifest_json TEXT NOT NULL,
			mode TEXT NOT NULL DEFAULT 'UNVERIFIED',
			sharing TEXT NOT NULL DEFAULT 'LOCAL_ONLY',
			status TEXT NOT NULL DEFAULT 'AVAILABLE',
			updated_at TEXT NOT NULL,
			PRIMARY KEY (id, version)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_offline_packs_created ON offline_packs(created_at DESC);`,

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
	ID         string    `json:"id"`
	Hex        string    `json:"hex"`
	Flight     string    `json:"flight"`
	Lat        float64   `json:"lat"`
	Lon        float64   `json:"lon"`
	AltBaro    int       `json:"alt_baro"`
	GS         float64   `json:"gs"`
	Track      float64   `json:"track"`
	Seen       float64   `json:"seen"`
	SourceID   string    `json:"source_id"`
	SourceType string    `json:"source_type"`
	ObservedAt time.Time `json:"observed_at"`
	Verified   int       `json:"verified"`
	Synthetic  int       `json:"synthetic"`
}

// OfflinePack is a verified offline data pack manifest record.
type OfflinePack struct {
	ID           string `json:"id"`
	Version      string `json:"version"`
	Sha256       string `json:"sha256"`
	Bytes        int64  `json:"bytes"`
	CreatedAt    string `json:"created_at"`
	ManifestJSON string `json:"-"`
	Mode         string `json:"mode"`
	Sharing      string `json:"sharing"`
	Status       string `json:"status"`
	UpdatedAt    string `json:"updated_at"`
}

// UpsertOfflinePack stores a manifest without accepting arbitrary file paths.
func (s *Storage) UpsertOfflinePack(p *OfflinePack) error {
	if p == nil || p.ID == "" || p.Version == "" || p.Sha256 == "" {
		return fmt.Errorf("invalid offline pack")
	}
	_, err := s.db.Exec(`INSERT INTO offline_packs (id,version,sha256,bytes,created_at,manifest_json,mode,sharing,status,updated_at)
		VALUES (?,?,?,?,?,?,?,?,?,?)
		ON CONFLICT(id,version) DO UPDATE SET sha256=excluded.sha256,bytes=excluded.bytes,
		manifest_json=excluded.manifest_json,mode=excluded.mode,sharing=excluded.sharing,
		status=excluded.status,updated_at=excluded.updated_at`,
		p.ID, p.Version, p.Sha256, p.Bytes, p.CreatedAt, p.ManifestJSON, p.Mode, p.Sharing, p.Status, p.UpdatedAt)
	return err
}

// ListOfflinePacks returns newest manifests, optionally filtered by pack id.
func (s *Storage) ListOfflinePacks(id string, limit int) ([]OfflinePack, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	var rows *sql.Rows
	var err error
	if id != "" {
		rows, err = s.db.Query(`SELECT id,version,sha256,bytes,created_at,manifest_json,mode,sharing,status,updated_at FROM offline_packs WHERE id=? ORDER BY created_at DESC LIMIT ?`, id, limit)
	} else {
		rows, err = s.db.Query(`SELECT id,version,sha256,bytes,created_at,manifest_json,mode,sharing,status,updated_at FROM offline_packs ORDER BY created_at DESC LIMIT ?`, limit)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []OfflinePack
	for rows.Next() {
		var p OfflinePack
		if err := rows.Scan(&p.ID, &p.Version, &p.Sha256, &p.Bytes, &p.CreatedAt, &p.ManifestJSON, &p.Mode, &p.Sharing, &p.Status, &p.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// CurrentOfflinePack returns the newest version of a pack.
func (s *Storage) CurrentOfflinePack(id string) (*OfflinePack, error) {
	packs, err := s.ListOfflinePacks(id, 1)
	if err != nil || len(packs) == 0 {
		return nil, err
	}
	return &packs[0], nil
}

// PruneOfflinePacks removes old manifest rows (not files).
func (s *Storage) PruneOfflinePacks(before time.Time) (int64, error) {
	r, err := s.db.Exec(`DELETE FROM offline_packs WHERE created_at < ?`, before.UTC().Format(time.RFC3339))
	if err != nil {
		return 0, err
	}
	n, _ := r.RowsAffected()
	return n, nil
}
