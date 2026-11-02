package simulator

import (
	"fmt"
	"sync"
	"time"
)

// Engine manages simulation modes for the Community Station.
type Engine struct {
	mode            Mode
	speed           float64
	loop            bool
	syntheticCount  int
	startLat        float64
	startLon        float64
	dirDeg          float64
	speedMps        float64
	altM            float64
	rateHz          int
	seed            int64

	mu              sync.RWMutex
	running         bool
	observations    []Observation
}

// Mode represents the simulation mode.
type Mode string

const (
	ModeReplay         Mode = "replay"
	ModeSynthetic      Mode = "synthetic"
	ModeLiveTransform  Mode = "live_transform"
)

// Observation represents a simulated aircraft observation.
type Observation struct {
	ID            string    `json:"id"`
	SourceID      string    `json:"source_id"`
	SourceType    string    `json:"source_type"`
	Hex           string    `json:"hex"`
	Flight        string    `json:"flight"`
	Lat           float64   `json:"lat"`
	Lon           float64   `json:"lon"`
	AltM          float64   `json:"alt_m"`
	SpeedMps      float64   `json:"speed_mps"`
	HeadingDeg    float64   `json:"heading_deg"`
	ObservedAt    time.Time `json:"observed_at"`
	Synthetic     bool      `json:"synthetic"`
	SimulationNote string   `json:"simulation_note"`
}

// Config provides simulation parameters.
type Config struct {
	Enabled         bool    `json:"enabled"`
	Mode            string  `json:"mode"`
	Objects         int     `json:"objects"`
	StartLat        float64 `json:"start_lat"`
	StartLon        float64 `json:"start_lon"`
	DirDeg          float64 `json:"dir_deg"`
	SpeedMps        float64 `json:"speed_mps"`
	AltM            float64 `json:"alt_m"`
	RateHz          int     `json:"rate_hz"`
	Speed           float64 `json:"speed"`
}

// NewEngine creates a new simulator engine.
func NewEngine(modeStr string, cfg Config) (*Engine, error) {
	mode := ModeSynthetic
	switch modeStr {
	case "replay":
		mode = ModeReplay
	case "live_transform":
		mode = ModeLiveTransform
	case "synthetic":
		mode = ModeSynthetic
	}

	engine := &Engine{
		mode:           mode,
		speed:          cfg.Speed,
		syntheticCount: cfg.Objects,
		startLat:       cfg.StartLat,
		startLon:       cfg.StartLon,
		dirDeg:         cfg.DirDeg,
		speedMps:       cfg.SpeedMps,
		altM:           cfg.AltM,
		rateHz:         cfg.RateHz,
		seed:           42,
		running:        false,
		observations:   make([]Observation, 0, cfg.Objects),
	}

	if engine.syntheticCount < 1 {
		engine.syntheticCount = 10
	}
	if engine.rateHz < 1 {
		engine.rateHz = 1
	}

	return engine, nil
}

// Start begins the simulation.
func (e *Engine) Start() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.running {
		return
	}

	e.running = true
	e.observations = make([]Observation, 0, e.syntheticCount)

	go e.simulateLoop()
}

// Stop stops the simulation.
func (e *Engine) Stop() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.running = false
	e.observations = e.observations[:0]
}

// Status returns the current simulation status.
func (e *Engine) Status() map[string]interface{} {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return map[string]interface{}{
		"running":   e.running,
		"mode":      string(e.mode),
		"speed":     e.speed,
		"objects":   len(e.observations),
		"synthetic": true,
		"simulated": true,
	}
}

// Mode returns the current simulation mode.
func (e *Engine) Mode() Mode {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.mode
}

// GetObservations returns current simulated observations.
func (e *Engine) GetObservations() []Observation {
	e.mu.RLock()
	defer e.mu.RUnlock()

	result := make([]Observation, len(e.observations))
	copy(result, e.observations)
	return result
}

// SetMode changes the simulation mode.
func (e *Engine) SetMode(modeStr string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	switch modeStr {
	case "replay", "synthetic", "live_transform":
		e.mode = Mode(modeStr)
	}
}

// SetCount updates the synthetic count.
func (e *Engine) SetCount(count int) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if count > 0 {
		e.syntheticCount = count
	}
}

// SetSpeed updates the simulation speed.
func (e *Engine) SetSpeed(speed float64) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if speed > 0 {
		e.speed = speed
	}
}

// simulateLoop runs the simulation in the background.
func (e *Engine) simulateLoop() {
	ticker := time.NewTicker(time.Duration(1000/e.rateHz) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			e.mu.Lock()
			if !e.running {
				e.mu.Unlock()
				return
			}

			e.generateObservations()
			e.mu.Unlock()
		default:
			e.mu.Lock()
			running := e.running
			e.mu.Unlock()

			if !running {
				return
			}

			time.Sleep(10 * time.Millisecond)
		}
	}
}

// generateObservations creates synthetic observations.
func (e *Engine) generateObservations() {
	now := time.Now()

	e.observations = e.observations[:0]

	for i := 0; i < e.syntheticCount; i++ {
		lat := e.startLat + float64(i)*0.001
		lon := e.startLon + float64(i)*0.001

		obs := Observation{
			ID:             fmt.Sprintf("SIM-%02d", i),
			SourceID:       "sim-sensor-01",
			SourceType:     "synthetic",
			Hex:            fmt.Sprintf("SIM%02X", i),
			Flight:         fmt.Sprintf("TEST%03d", i+1),
			Lat:            lat,
			Lon:            lon,
			AltM:           e.altM + float64(i)*50,
			SpeedMps:       e.speedMps + float64(i)*2,
			HeadingDeg:     e.dirDeg,
			ObservedAt:     now,
			Synthetic:      true,
			SimulationNote: "EDUCATIONAL SIMULATION - SYNTHETIC DATA",
		}

		e.observations = append(e.observations, obs)
	}
}