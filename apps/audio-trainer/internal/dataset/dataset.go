package dataset

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/pwojczyn/NieboGuard-Community/apps/audio-trainer/internal/config"
)

type DatasetEntry struct {
	Path         string    `json:"path"`
	Duration     float64   `json:"duration"`
	SampleRate   int       `json:"sample_rate"`
	Channels     int       `json:"channels"`
	BitDepth     int       `json:"bit_depth"`
	Labels       []string  `json:"labels"`
	Verified     bool      `json:"verified"`
	ImportedAt   time.Time `json:"imported_at"`
	Status       string    `json:"status"`
}

type Event struct {
	Type      string `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

type DatasetManager struct {
	logger  *slog.Logger
	cfg     *config.Configuration
	history []DatasetEntry
	events  chan Event
	mu      sync.RWMutex
}

func New(cfg *config.Configuration) (*DatasetManager, error) {
	logger := slog.Default().With("component", "dataset")
	dm := &DatasetManager{
		logger:   logger,
		cfg:      cfg,
		events:   make(chan Event, 100),
		history:  []DatasetEntry{},
	}

	if err := os.MkdirAll(filepath.Join(cfg.WorkDir, cfg.Dataset.ImportDir), 0o755); err != nil {
		return nil, fmt.Errorf("creating import directory: %w", err)
	}
	return dm, nil
}

func (dm *DatasetManager) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			time.Sleep(30 * time.Second)
		}
	}()
}

func (dm *DatasetManager) ImportWAV(path string) error {
	info := DatasetEntry{
		Path:         path,
		Duration:     10.0,
		SampleRate:   16000,
		Channels:     1,
		BitDepth:     16,
		Labels:       []string{"unknown"},
		Verified:     false,
		Status:       "processing",
		ImportedAt:   time.Now(),
	}

	dm.mu.Lock()
	defer dm.mu.Unlock()

	dm.history = append(dm.history, info)
	dm.events <- Event{
		Type:      "wav_imported",
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"path": path,
		},
	}
	return nil
}

func (dm *DatasetManager) VerifyWAV(path string) error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	for i := range dm.history {
		if dm.history[i].Path == path {
			dm.history[i].Verified = true
			dm.history[i].Status = "verified"
			dm.history[i].Labels = []string{"verified"}
			break
		}
	}
	dm.events <- Event{
		Type:      "wav_verified",
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"path": path,
		},
	}
	return nil
}

func (dm *DatasetManager) GetHistory() []DatasetEntry {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	return dm.history
}

func (dm *DatasetManager) GetEntry(path string) (*DatasetEntry, bool) {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	for _, entry := range dm.history {
		if entry.Path == path {
			return &entry, true
		}
	}
	return nil, false
}
