package training

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/pwojczyn/NieboGuard-Community/apps/audio-trainer/internal/config"
)

type Trainer struct {
	logger   *slog.Logger
	cfg      *config.Configuration
	status   Status
	mu       sync.RWMutex
	ctx      context.Context
	stop     context.CancelFunc
	history  []HistoryEntry
	events   chan Event
}

type Status string

const (
	StatusIdle       Status = "idle"
	StatusReady      Status = "ready"
	StatusTraining   Status = "training"
	StatusPaused     Status = "paused"
	StatusError      Status = "error"
)

type HistoryEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Event     Event `json:"event"`
}

type Event struct {
	Type      string `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

func New(cfg *config.Configuration) (*Trainer, error) {
	logger := slog.Default().With("component", "trainer")
	t := &Trainer{
		logger:  logger,
		cfg:     cfg,
		status:  StatusIdle,
		events:  make(chan Event, 100),
		history: []HistoryEntry{},
	}
	return t, nil
}

func (t *Trainer) Start(ctx context.Context) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.ctx != nil {
		return // already started
	}
	t.ctx, t.stop = context.WithCancel(ctx)

	go func() {
		for {
			select {
			case <-t.ctx.Done():
				return
			case ev := <-t.events:
				t.mu.Lock()
				t.history = append(t.history, HistoryEntry{Timestamp: ev.Timestamp, Event: ev})
				// Keep only last 1000 entries
				if len(t.history) > 1000 {
					t.history = t.history[1:]
				}
				t.mu.Unlock()
			}
		}
	}()
}

func (t *Trainer) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.stop != nil {
		t.stop()
		t.stop = nil
		t.ctx = nil
		t.logger.Info("Trainer stopped")
	}
}

func (t *Trainer) Status() Status {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.status
}

func (t *Trainer) SetStatus(s Status) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.status = s
	t.logger.Info("Status changed", "status", s)
}

func (t *Trainer) StartTraining() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.status != StatusIdle && t.status != StatusPaused {
		return fmt.Errorf("cannot start training from status %s", t.status)
	}

	// Record start event
	t.events <- Event{
		Type:      "training_start",
		Timestamp: time.Now(),
		Data:      map[string]interface{}{},
	}

	// Simulate training progress for demonstration
	go func() {
		t.SetStatus(StatusTraining)
		times := []time.Duration{1 * time.Second, 2 * time.Second, 3 * time.Second}
		for i, _ := range times {
			select {
			case <-t.ctx.Done():
				return
			default:
			}
			t.events <- Event{
				Type:      "epoch_complete",
				Timestamp: time.Now(),
				Data: map[string]interface{}{
					"epoch":     i + 1,
					"loss":      0.5 - float64(i)*0.1,
					"accuracy":  0.6 + float64(i)*0.05,
				},
			}
		}

		t.SetStatus(StatusIdle)
		t.events <- Event{
			Type:      "training_complete",
			Timestamp: time.Now(),
			Data: map[string]interface{}{
				"epochs_completed": len(times),
			},
		}
	}()

	return nil
}

func (t *Trainer) Pause() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.status != StatusTraining {
		return fmt.Errorf("cannot pause from status %s", t.status)
	}

	// Pause the simulated training
	t.events <- Event{
		Type:      "training_paused",
		Timestamp: time.Now(),
		Data:      map[string]interface{}{},
	}
	t.SetStatus(StatusPaused)
	return nil
}

func (t *Trainer) GetHistory(limit int) []HistoryEntry {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if limit <= 0 || limit > len(t.history) {
		limit = len(t.history)
	}

	result := make([]HistoryEntry, limit)
	copy(result, t.history[len(t.history)-limit:])
	return result
}

func (t *Trainer) ExportHistory() ([]HistoryEntry, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.history, nil
}