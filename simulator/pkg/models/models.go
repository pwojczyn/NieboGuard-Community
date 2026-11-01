package models

import "time"

// AircraftEvent is the unified observation model for all protocols.
type AircraftEvent struct {
	SourceID       string    `json:"source_id"`
	SourceType     string    `json:"source_type"`
	ObservationID  string    `json:"observation_id"`
	ObjectType     string    `json:"object_type"`
	Hex            string    `json:"hex,omitempty"`
	Flight         string    `json:"flight,omitempty"`
	Track          float64   `json:"track,omitempty"`
	Lat            float64   `json:"lat,omitempty"`
	Lon            float64   `json:"lon,omitempty"`
	AltM           float64   `json:"alt_m,omitempty"`
	Gs             float64   `json:"gs,omitempty"`
	HeadingDeg     float64   `json:"heading_deg,omitempty"`
	ObservedAt     time.Time `json:"observed_at"`
	ReceivedAt     time.Time `json:"received_at"`
	Quality        float64   `json:"quality,omitempty"`
	Verified       bool      `json:"verified"`
	Synthetic      bool      `json:"synthetic"`
	SimulationNote string    `json:"simulation_note"`
	Seen           float64   `json:"seen,omitempty"`
}

// SimulatedObservation represents a simulated/synthetic aircraft observation.
// All fields are marked synthetic=true to distinguish from real measurements.
type SimulatedObservation struct {
	ID             string    `json:"id"`
	SourceID       string    `json:"source_id"`
	SourceType     string    `json:"source_type"`
	ObservationID  string    `json:"observation_id"`
	ObjectType     string    `json:"object_type"`
	Hex            string    `json:"hex,omitempty"`
	Flight         string    `json:"flight,omitempty"`
	Lat            float64   `json:"lat,omitempty"`
	Lon            float64   `json:"lon,omitempty"`
	AltM           float64   `json:"alt_m,omitempty"`
	SpeedMps       float64   `json:"speed_mps,omitempty"`
	HeadingDeg     float64   `json:"heading_deg,omitempty"`
	ObservedAt     time.Time `json:"observed_at"`
	ReceivedAt     time.Time `json:"received_at"`
	Quality        float64   `json:"quality,omitempty"`
	Verified       bool      `json:"verified"`
	Synthetic      bool      `json:"synthetic"`
	SimulationNote string    `json:"simulation_note"`
	Seen           float64   `json:"seen,omitempty"`
}

// SimulatedTrack represents a tracked object trajectory
type SimulatedTrack struct {
	TrackID      string    `json:"track_id"`
	SourceID     string    `json:"source_id"`
	Observations []SimulatedObservation `json:"observations"`
	StartTime    time.Time `json:"start_time"`
	LastUpdate   time.Time `json:"last_update"`
	TTLSeconds   int       `json:"ttl_seconds"`
	Active       bool      `json:"active"`
}

// SimulationConfig is the configuration for simulation modes
type SimulationConfig struct {
	Mode           string  `json:"mode"`
	Speed          float64 `json:"speed"`
	Loop           bool    `json:"loop"`
	SyntheticCount int     `json:"synthetic_count"`
	StartLat       float64 `json:"start_lat"`
	StartLon       float64 `json:"start_lon"`
	DirDeg         float64 `json:"dir_deg"`
	SpeedMps       float64 `json:"speed_mps"`
	AltM           float64 `json:"alt_m"`
	RateHz         int     `json:"rate_hz"`
	Seed           int64   `json:"seed"`
}

// ASTERIXCat021Record is a simplified ASTERIX CAT021 record
type ASTERIXCat021Record struct {
	FSPEC         []byte `json:"fspec"`
	DataSourceID  uint16 `json:"data_source_id"`
	TrackNumber   uint16 `json:"track_number,omitempty"`
	Lat           float64 `json:"lat,omitempty"`
	Lon           float64 `json:"lon,omitempty"`
	Altitude      uint16 `json:"altitude,omitempty"`
	Speed         uint16 `json:"speed,omitempty"`
	Track         uint16 `json:"track,omitempty"`
	Timestamp     uint32 `json:"timestamp"`
}

// ASTERIXCat048Record is a simplified ASTERIX CAT048 record
type ASTERIXCat048Record struct {
	FSPEC           []byte `json:"fspec"`
	DataSourceID    uint16 `json:"data_source_id"`
	TargetAddress   uint32 `json:"target_address,omitempty"`
	SlantRange      uint16 `json:"slant_range,omitempty"`
	Azimuth         uint16 `json:"azimuth,omitempty"`
	TrackNumber     uint16 `json:"track_number,omitempty"`
	ModeACCode      uint16 `json:"mode_ac_code,omitempty"`
	ModeSCode       string `json:"mode_s_code,omitempty"`
	Position        [2]float64 `json:"position,omitempty"`
	TimeOfOccurrence uint32 `json:"time_of_occurrence"`
}

// ASTERIXCat062Record is a simplified ASTERIX CAT062 record
type ASTERIXCat062Record struct {
	FSPEC           []byte `json:"fspec"`
	DataSourceID    uint16 `json:"data_source_id"`
	TrackNumber     uint16 `json:"track_number"`
	Position        [2]float64 `json:"position"`
	Speed           uint16 `json:"speed,omitempty"`
	TrackAngle      uint16 `json:"track_angle,omitempty"`
	TrackStatus     uint8  `json:"track_status,omitempty"`
	TimeOfTrack     uint32 `json:"time_of_track"`
}
