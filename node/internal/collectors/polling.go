package collectors

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pwojczyn/NieboGuard-Community/node/internal/storage"
	"github.com/pwojczyn/NieboGuard-Community/node/pkg/config"
)

// pollingSource fetches a configured JSON source and stores aircraft events.
// It is intentionally loopback-oriented; the URL is validated before every poll.
type pollingSource struct {
	cfg    config.SourceConfig
	store  *storage.Storage
	logger *slog.Logger

	mu        sync.RWMutex
	active    bool
	cancel    context.CancelFunc
	count     int
	lastFetch time.Time
	lastError string
}

func newPollingSource(cfg config.SourceConfig, store *storage.Storage, logger *slog.Logger) *pollingSource {
	if cfg.IntervalSeconds <= 0 {
		cfg.IntervalSeconds = 5
	}
	if cfg.Name == "" {
		cfg.Name = cfg.ID
	}
	return &pollingSource{cfg: cfg, store: store, logger: logger}
}

func (p *pollingSource) Start() error {
	if p.cfg.URL == "" {
		return fmt.Errorf("source URL required")
	}
	if err := validSourceURL(p.cfg.URL); err != nil {
		return err
	}
	p.mu.Lock()
	if p.active {
		p.mu.Unlock()
		return nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel
	p.active = true
	p.lastError = ""
	p.mu.Unlock()
	go p.loop(ctx)
	return nil
}

func (p *pollingSource) Stop() error {
	p.mu.Lock()
	cancel := p.cancel
	p.cancel = nil
	p.active = false
	p.mu.Unlock()
	if cancel != nil {
		cancel()
	}
	return nil
}

func (p *pollingSource) loop(ctx context.Context) {
	p.fetch(ctx)
	ticker := time.NewTicker(time.Duration(p.cfg.IntervalSeconds) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.fetch(ctx)
		}
	}
}

func (p *pollingSource) fetch(ctx context.Context) {
	fetchCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	records, err := p.read(fetchCtx)
	if err != nil {
		p.setError(err)
		return
	}
	now := time.Now().UTC()
	stored := 0
	for _, record := range records {
		if record.Hex == "" || record.Lat == 0 || record.Lon == 0 {
			continue
		}
		event := &storage.AircraftEvent{
			ID:         fmt.Sprintf("%s-%d-%s", p.cfg.ID, now.UnixNano(), record.Hex),
			Hex:        record.Hex,
			Flight:     record.Flight,
			Lat:        record.Lat,
			Lon:        record.Lon,
			AltBaro:    record.AltBaro,
			GS:         record.GS,
			Track:      record.Track,
			Seen:       record.Seen,
			SourceID:   p.cfg.ID,
			SourceType: p.cfg.Type,
			ObservedAt: now,
			Verified:   0,
			Synthetic:  0,
		}
		if err := p.store.InsertAircraft(event); err != nil {
			p.logger.Error("storing aircraft", "source", p.cfg.ID, "error", err)
			continue
		}
		stored++
	}
	p.mu.Lock()
	p.count = stored
	p.lastFetch = now
	p.lastError = ""
	p.mu.Unlock()
}

func (p *pollingSource) read(ctx context.Context) ([]AircraftRecord, error) {
	if err := validSourceURL(p.cfg.URL); err != nil {
		return nil, err
	}
	typ := strings.ToLower(strings.TrimSpace(p.cfg.Type))
	if typ == "http_json" {
		adapter := &HTTPPollAdapter{URL: p.cfg.URL, Timeout: 8 * time.Second}
		items, err := adapter.Fetch()
		if err != nil {
			return nil, err
		}
		return recordsFromMaps(items), nil
	}
	adapter := &ADSBJSONAdapter{URL: p.cfg.URL, Timeout: 8 * time.Second}
	return adapter.Fetch()
}

func validSourceURL(raw string) error {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || u.Host == "" {
		return fmt.Errorf("invalid source URL")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("source URL must use http or https")
	}
	if u.User != nil || u.Fragment != "" {
		return fmt.Errorf("source URL must not contain credentials or fragment")
	}
	return nil
}

func recordsFromMaps(items []map[string]interface{}) []AircraftRecord {
	records := make([]AircraftRecord, 0, len(items))
	for _, item := range items {
		record := AircraftRecord{
			Hex:     strings.TrimSpace(stringValue(item["hex"])),
			Flight:  strings.TrimSpace(stringValue(item["flight"])),
			Lat:     floatValue(item["lat"]),
			Lon:     floatValue(item["lon"]),
			AltBaro: int(floatValue(item["alt_baro"])),
			GS:      floatValue(item["gs"]),
			Track:   floatValue(item["track"]),
			Seen:    floatValue(item["seen"]),
		}
		if record.Hex != "" && record.Lat != 0 && record.Lon != 0 {
			records = append(records, record)
		}
	}
	return records
}

func stringValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	case nil:
		return ""
	default:
		return fmt.Sprint(v)
	}
}

func floatValue(value interface{}) float64 {
	switch v := value.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case json.Number:
		f, _ := v.Float64()
		return f
	case string:
		f, _ := strconv.ParseFloat(strings.TrimSpace(v), 64)
		return f
	default:
		return 0
	}
}

func (p *pollingSource) setError(err error) {
	p.mu.Lock()
	p.lastError = err.Error()
	p.mu.Unlock()
	p.logger.Warn("source fetch failed", "source", p.cfg.ID, "error", err)
}

func (p *pollingSource) Status() SourceStatus {
	p.mu.RLock()
	defer p.mu.RUnlock()
	status := SourceStatus{Active: p.active, Records: p.count}
	if !p.lastFetch.IsZero() {
		status.LastFetch = p.lastFetch.Format(time.RFC3339)
	}
	if p.lastError != "" {
		status.Errors = []string{p.lastError}
	}
	return status
}

func (p *pollingSource) GetSource() SourceInfo {
	p.mu.RLock()
	defer p.mu.RUnlock()
	status := "inactive"
	if p.active {
		status = "active"
	}
	if p.lastError != "" {
		status = "error"
	}
	return SourceInfo{
		ID:      p.cfg.ID,
		Type:    p.cfg.Type,
		Name:    p.cfg.Name,
		URL:     p.cfg.URL,
		Status:  status,
		License: "custom",
		Count:   p.count,
	}
}
