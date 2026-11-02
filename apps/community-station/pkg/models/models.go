package models

import "time"

// AircraftEvent represents an aircraft observation.
type AircraftEvent struct {
	ID         string    `json:"id"`
	Hex        string    `json:"hex"`
	Flight     string    `json:"flight"`
	Lat        float64   `json:"lat,omitempty"`
	Lon        float64   `json:"lon,omitempty"`
	AltBaro    int       `json:"alt_baro,omitempty"`
	GS         float64   `json:"gs,omitempty"`
	Track      float64   `json:"track,omitempty"`
	Seen       float64   `json:"seen,omitempty"`
	SourceID   string    `json:"source_id"`
	SourceType string    `json:"source_type"`
	ObservedAt time.Time `json:"observed_at"`
	Verified   bool      `json:"verified"`
	Synthetic  bool      `json:"synthetic"`
}