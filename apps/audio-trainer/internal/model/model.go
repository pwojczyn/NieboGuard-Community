package model

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/pwojczyn/NieboGuard-Community/apps/audio-trainer/internal/config"
)

type ModelManager struct {
	logger   *slog.Logger
	cfg      *config.Configuration
	status   Status
	mu       sync.RWMutex
	events    chan Event
	history  []ModelEntry
}

type Status string

const (
	StatusIdle       Status = "idle"
	StatusLoading    Status = "loading"
	StatusReady      Status = "ready"
	StatusTraining   Status = "training"
	StatusEvaluating Status = "evaluating"
	StatusError      Status = "error"
	StatusPaused     Status = "paused"
)

type ModelEntry struct {
	Path       string    `json:"path"`
	Format     string    `json:"format"`
	Size       int64     `json:"size"`
	Loaded     bool      `json:"loaded"`
	Accuracy   float64   `json:"accuracy"`
	Precision  float64   `json:"precision"`
	Recall     float64   `json:"recall"`
	F1Score    float64   `json:"f1_score"`
	ConfMatrix []int     `json:"conf_matrix"`
	Classes    []string  `json:"classes"`
	Version    string    `json:"version"`
	CreatedAt  time.Time `json:"created_at"`
}

type Event struct {
	Type      string `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

func New(cfg *config.Configuration) (*ModelManager, error) {
	logger := slog.Default().With("component", "model")
	mm := &ModelManager{
		logger:   logger,
		cfg:      cfg,
		status:   StatusIdle,
		events:   make(chan Event, 100),
		history:  []ModelEntry{},
	}
	return mm, nil
}

func (mm *ModelManager) Start(ctx context.Context) {
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

func (mm *ModelManager) LoadModel(path string) error {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	if mm.status == StatusError {
		return fmt.Errorf("model manager in error state")
	}

	mm.status = StatusLoading
	mm.events <- Event{
		Type:      "model_loading",
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"path": path,
		},
	}

	// Check file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		mm.status = StatusError
		mm.events <- Event{
			Type:      "model_load_error",
			Timestamp: time.Now(),
			Data: map[string]interface{}{
				"error": err.Error(),
			},
		}
		return err
	}

	mm.status = StatusReady
	mm.events <- Event{
		Type:      "model_loaded",
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"path": path,
		},
	}
	return nil
}

func (mm *ModelManager) Status() Status {
	mm.mu.RLock()
	defer mm.mu.RUnlock()
	return mm.status
}

func (mm *ModelManager) SetStatus(s Status) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	mm.status = s
}

func (mm *ModelManager) Pause() error {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	if mm.status != StatusTraining {
		return fmt.Errorf("cannot pause from status %s", mm.status)
	}

	mm.status = StatusPaused
	return nil
}

func (mm *ModelManager) GetInfo() interface{} {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	return map[string]interface{}{
		"status":   mm.status,
		"model":    mm.history,
	}
}