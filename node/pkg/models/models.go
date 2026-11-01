package models

import (
	"time"
)

// Aircraft represents an aircraft observation.
type Aircraft struct {
	ID           string    `json:"id"`
	Hex          string    `json:"hex"`
	Flight       string    `json:"flight"`
	Lat          float64   `json:"lat,omitempty"`
	Lon          float64   `json:"lon,omitempty"`
	AltBaro      int       `json:"alt_baro,omitempty"`
	GS           float64   `json:"gs,omitempty"`
	Track        float64   `json:"track,omitempty"`
	Seen         float64   `json:"seen,omitempty"`
	SourceID     string    `json:"source_id"`
	SourceType   string    `json:"source_type"`
	ObservedAt   time.Time `json:"observed_at"`
	Verified     bool      `json:"verified"`
	Synthetic    bool      `json:"synthetic"`
}

// Notice represents a safety notice or alert.
type Notice struct {
	ID           string    `json:"id"`
	Source       string    `json:"source"`
	Category     string    `json:"category"`
	Title        string    `json:"title"`
	Region       string    `json:"region"`
	Verification string    `json:"verification"`
	PublishedAt  time.Time `json:"published_at"`
	ReceivedAt   time.Time `json:"received_at"`
}

// Source represents a data source.
type Source struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Name        string    `json:"name"`
	URL         string    `json:"url"`
	Status      string    `json:"status"`
	LastUpdate  time.Time `json:"last_update"`
	License     string    `json:"license"`
	RecordCount int       `json:"record_count"`
}

// HealthInfo represents node health information.
type HealthInfo struct {
	Status      string    `json:"status"`
	Uptime      string    `json:"uptime"`
	Version     string    `json:"version"`
	StartTime   time.Time `json:"start_time"`
	Objects     int       `json:"objects"`
	Sources     int       `json:"sources"`
	Connections int       `json:"connections"`
}

// APIResponse is a generic API response wrapper.
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

// APIError represents an API error.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}