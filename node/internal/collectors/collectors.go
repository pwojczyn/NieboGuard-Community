package collectors

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/pwojczyn/NieboGuard-Community/node/internal/storage"
	"github.com/pwojczyn/NieboGuard-Community/node/pkg/config"
)

// Collector is the interface for all data source adapters.
type Collector interface {
	Start() error
	Stop() error
	Status() SourceStatus
	GetSource() SourceInfo
}

// SourceInfo provides metadata about a source.
type SourceInfo struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	URL     string `json:"url"`
	Status  string `json:"status"`
	License string `json:"license"`
	Count   int    `json:"count"`
}

// SourceStatus provides runtime status.
type SourceStatus struct {
	Active    bool     `json:"active"`
	LastFetch string   `json:"last_fetch"`
	Records   int      `json:"records"`
	Errors    []string `json:"errors"`
}

// Manager manages data source collectors.
type Manager struct {
	cfg     *config.Configuration
	store   *storage.Storage
	logger  *slog.Logger
	sources map[string]Collector
}

// NewManager creates a new collector manager.
func NewManager(cfg *config.Configuration, store *storage.Storage, logger *slog.Logger) (*Manager, error) {
	if cfg == nil {
		cfg = &config.Configuration{}
	}
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	m := &Manager{
		cfg:     cfg,
		store:   store,
		logger:  logger,
		sources: make(map[string]Collector),
	}

	// Configured sources are authoritative when present. Empty configuration
	// preserves the original built-in source list.
	if len(cfg.Sources) > 0 {
		m.sources = make(map[string]Collector)
		for _, source := range cfg.Sources {
			if !source.Enabled || strings.TrimSpace(source.ID) == "" {
				continue
			}
			if source.Type == "" {
				source.Type = "adsb_json"
			}
			m.sources[source.ID] = newPollingSource(source, store, logger)
		}
	} else {
		m.initDefaultSources()
	}

	return m, nil
}

// initDefaultSources initializes built-in collectors.
func (m *Manager) initDefaultSources() {
	m.sources["local_adsb"] = &localADSbSource{
		name:   "Local ADS-B",
		logger: m.logger,
		store:  m.store,
	}

	m.sources["http_json"] = &httpJSONSource{
		name:   "HTTP JSON Polling",
		logger: m.logger,
		store:  m.store,
	}

	m.sources["rss_atom"] = &rssAtomSource{
		name:   "RSS/Atom Feed",
		logger: m.logger,
	}
}

// StartAll starts all registered collectors.
func (m *Manager) StartAll() error {
	m.logger.Info("starting collectors", "count", len(m.sources))
	for id, collector := range m.sources {
		if err := collector.Start(); err != nil {
			m.logger.Error("collector start failed", "id", id, "error", err)
			// Continue with other collectors, don't fail completely
		} else {
			m.logger.Info("collector started", "id", id)
		}
	}
	return nil
}

// StopAll stops all registered collectors.
func (m *Manager) StopAll() error {
	m.logger.Info("stopping collectors", "count", len(m.sources))
	for id, collector := range m.sources {
		if err := collector.Stop(); err != nil {
			m.logger.Error("collector stop failed", "id", id, "error", err)
		} else {
			m.logger.Info("collector stopped", "id", id)
		}
	}
	return nil
}

// GetSources returns all registered sources.
func (m *Manager) GetSources() []SourceInfo {
	var sources []SourceInfo
	for _, collector := range m.sources {
		sources = append(sources, collector.GetSource())
	}
	return sources
}

// GetSource returns a specific source by ID.
func (m *Manager) GetSource(id string) (SourceInfo, bool) {
	if collector, ok := m.sources[id]; ok {
		return collector.GetSource(), true
	}
	return SourceInfo{}, false
}

// localADSbSource represents a local ADS-B adapter (stub for P0).
type localADSbSource struct {
	name   string
	logger *slog.Logger
	store  *storage.Storage
	ctx    context.Context
}

func (l *localADSbSource) Start() error {
	l.logger.Info("local ADS-B collector started")
	return nil
}

func (l *localADSbSource) Stop() error {
	l.logger.Info("local ADS-B collector stopped")
	return nil
}

func (l *localADSbSource) Status() SourceStatus {
	return SourceStatus{
		Active:    true,
		LastFetch: "now",
		Records:   0,
	}
}

func (l *localADSbSource) GetSource() SourceInfo {
	return SourceInfo{
		ID:      "local_adsb",
		Type:    "local_adsb",
		Name:    l.name,
		Status:  "active",
		License: "public_domain",
	}
}

// httpJSONSource represents an HTTP JSON adapter (stub for P0).
type httpJSONSource struct {
	name   string
	logger *slog.Logger
	store  *storage.Storage
}

func (h *httpJSONSource) Start() error {
	h.logger.Info("HTTP JSON collector started")
	return nil
}

func (h *httpJSONSource) Stop() error {
	h.logger.Info("HTTP JSON collector stopped")
	return nil
}

func (h *httpJSONSource) Status() SourceStatus {
	return SourceStatus{
		Active:    true,
		LastFetch: "now",
		Records:   0,
	}
}

func (h *httpJSONSource) GetSource() SourceInfo {
	return SourceInfo{
		ID:      "http_json",
		Type:    "http_json",
		Name:    h.name,
		Status:  "inactive",
		License: "custom",
	}
}

// rssAtomSource represents an RSS adapter (stub for P0).
type rssAtomSource struct {
	name   string
	logger *slog.Logger
}

func (r *rssAtomSource) Start() error {
	r.logger.Info("RSS/Atom collector started")
	return nil
}

func (r *rssAtomSource) Stop() error {
	r.logger.Info("RSS/Atom collector stopped")
	return nil
}

func (r *rssAtomSource) Status() SourceStatus {
	return SourceStatus{
		Active:    true,
		LastFetch: "now",
		Records:   0,
	}
}

func (r *rssAtomSource) GetSource() SourceInfo {
	return SourceInfo{
		ID:      "rss_atom",
		Type:    "rss_atom",
		Name:    r.name,
		Status:  "inactive",
		License: "public_domain",
	}
}
