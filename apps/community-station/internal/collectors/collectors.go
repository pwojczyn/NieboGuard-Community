package collectors

import (
	"log/slog"

	"github.com/pwojczyn/NieboGuard-Community/apps/community-station/internal/config"
	"github.com/pwojczyn/NieboGuard-Community/apps/community-station/internal/storage"
	"github.com/pwojczyn/NieboGuard-Community/apps/community-station/pkg/models"
)

// SourceInfo provides metadata about a source.
type SourceInfo struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	Count     int    `json:"count"`
	Synthetic bool   `json:"synthetic"`
	Hardware  bool   `json:"hardware"`
}

// SourceStatus provides runtime status.
type SourceStatus struct {
	Active    bool     `json:"active"`
	Records   int      `json:"records"`
	LastFetch string   `json:"last_fetch"`
}

// Collector interface for data sources.
type Collector interface {
	Start() error
	Stop() error
	Status() SourceStatus
	SourceInfo() SourceInfo
}

// Manager manages all collectors.
type Manager struct {
	cfg     *config.Configuration
	store   *storage.Storage
	logger  *slog.Logger
	sources map[string]Collector
}

// New creates a new manager.
func New(cfg *config.Configuration, store *storage.Storage, logger *slog.Logger) (*Manager, error) {
	m := &Manager{
		cfg:     cfg,
		store:   store,
		logger:  logger,
		sources: make(map[string]Collector),
	}
	m.initSources()
	return m, nil
}

func (m *Manager) initSources() {
	m.sources["simulator"] = &simSource{
		name:   "Simulator",
		logger: m.logger,
	}
}

// StartAll starts collectors.
func (m *Manager) StartAll() error {
	m.logger.Info("starting collectors")
	return nil
}

// StopAll stops collectors.
func (m *Manager) StopAll() error {
	m.logger.Info("stopping collectors")
	return nil
}

// GetSources returns all sources.
func (m *Manager) GetSources() []SourceInfo {
	var sources []SourceInfo
	for _, c := range m.sources {
		sources = append(sources, c.SourceInfo())
	}
	return sources
}

// InsertAircraft inserts aircraft events into storage.
func (m *Manager) InsertAircraft(events []models.AircraftEvent) error {
	return nil
}

// simSource is a simulated source for demonstration.
type simSource struct {
	name   string
	logger *slog.Logger
}

func (s *simSource) Start() error  { s.logger.Info("simulator source started"); return nil }
func (s *simSource) Stop() error   { return nil }
func (s *simSource) Status() SourceStatus {
	return SourceStatus{Active: true, LastFetch: "now", Records: 0}
}
func (s *simSource) SourceInfo() SourceInfo {
	return SourceInfo{ID: "simulator", Type: "simulator", Name: s.name, Status: "active", Synthetic: true, Count: 0}
}

// GetSource returns a specific collector.
func (m *Manager) GetSource(id string) (SourceInfo, bool) {
	if c, ok := m.sources[id]; ok {
		return c.SourceInfo(), true
	}
	return SourceInfo{}, false
}
