package sim

import (
	"sync"
	"time"

	"github.com/pwojczyn/NieboGuard-Community/simulator/pkg/models"
)

// Engine manages simulation modes

type Mode int

const (
	ModeReplay Mode = iota
	ModeSynthetic
	ModeLiveTransform
)

type Engine struct {
	Mode            Mode
	Speed           float64 // 0.25, 1, 2, 5, 10
	Loop            bool
	SyntheticCount  int
	StartLat        float64
	StartLon        float64
	DirDeg          float64 // direction
	SpeedMps        float64
	AltM            float64
	RateHz          int // updates per second
	Seed            int64
	
	// Optimization: object pool for observation reuse
	obsPool sync.Pool
	
	// Pre-allocated buffers
	ids       []string
	observations []models.SimulatedObservation
	
	// Time tracking
	lastFrame time.Time
	frameCount uint64
}

func NewEngine(mode Mode) *Engine {
	e := &Engine{
		Mode:            mode,
		Speed:           1.0,
		Loop:            false,
		SyntheticCount:  10,
		StartLat:        52.2297,
		StartLon:        21.0122,
		DirDeg:          90,
		SpeedMps:        250,
		AltM:            10000,
		RateHz:          1,
		Seed:            42,
		lastFrame:       time.Now(),
	}
	
	// Initialize object pool
	e.obsPool = sync.Pool{
		New: func() interface{} {
			return make([]models.SimulatedObservation, 0, e.SyntheticCount)
		},
	}
	
	// Pre-allocate IDs
	e.ids = make([]string, e.SyntheticCount)
	for i := 0; i < e.SyntheticCount; i++ {
		e.ids[i] = "SIM-" + string(rune('0'+i%10)) // simplified
	}
	
	// Pre-allocate observations slice
	e.observations = make([]models.SimulatedObservation, 0, e.SyntheticCount)
	
	return e
}

// SetCount updates the synthetic count and reinitializes pools
func (e *Engine) SetCount(count int) {
	e.SyntheticCount = count
	e.ids = make([]string, count)
	for i := 0; i < count; i++ {
		e.ids[i] = "SIM-" + string(rune('0'+i%10))
	}
	e.observations = make([]models.SimulatedObservation, 0, count)
	e.obsPool = sync.Pool{
		New: func() interface{} {
			return make([]models.SimulatedObservation, 0, count)
		},
	}
}

func (e *Engine) NextFrame() ([]models.SimulatedObservation, error) {
	// Rate limiting
	if e.RateHz > 0 {
		minInterval := time.Second / time.Duration(e.RateHz)
		sinceLast := time.Since(e.lastFrame)
		if sinceLast < minInterval {
			time.Sleep(minInterval - sinceLast)
		}
	}
	e.lastFrame = time.Now()
	e.frameCount++
	
	switch e.Mode {
	case ModeSynthetic:
		return e.syntheticFrame()
	case ModeReplay:
		return e.replayFrame()
	case ModeLiveTransform:
		return e.liveFrame()
	default:
		return nil, nil // Return empty instead of error for unknown mode
	}
}

func (e *Engine) syntheticFrame() ([]models.SimulatedObservation, error) {
	// Get pooled slice
	obs := e.obsPool.Get().([]models.SimulatedObservation)
	obs = obs[:0] // Reset length, keep capacity
	
	rng := randFromSeed(e.Seed + int64(e.frameCount))
	now := time.Now()
	
	// Pre-calculate constants
	latStep := 0.001
	lonStep := 0.001
	latNoise := 0.005
	lonNoise := 0.005
	altBase := 10000.0
	altStep := 50.0
	altNoise := 500.0
	speedBase := 250.0
	speedNoise := 50.0
	headingBase := 90.0
	headingStep := 10.0
	
	for i := 0; i < e.SyntheticCount; i++ {
		t := now.Add(-time.Duration(i) * time.Second)
		lat := e.StartLat + (float64(i) * latStep) + (rng.Float64() - 0.5) * latNoise * 2
		lon := e.StartLon + (float64(i) * lonStep) + (rng.Float64() - 0.5) * lonNoise * 2
		
		// Use math.Floor for deterministic IDs
		id := e.ids[i]
		
		obs = append(obs, models.SimulatedObservation{
			ID:             id,
			SourceID:       "sim-sensor-01",
			SourceType:     "synthetic",
			ObservationID:  "sim-obs-" + string(rune(now.Unix()%1e6)) + "-" + string(rune(i)),
			ObjectType:     "aircraft",
			Hex:            "SIM" + string(rune(i%16)) + string(rune(i%16)) + string(rune(i%16)) + string(rune(i%16)),
			Flight:         "TEST" + string(rune(i/100%10+'0')) + string(rune(i/10%10+'0')) + string(rune(i%10+'0')),
			Lat:            lat,
			Lon:            lon,
			AltM:           altBase + float64(i)*altStep + rng.Float64()*altNoise,
			SpeedMps:       speedBase + rng.Float64()*speedNoise,
			HeadingDeg:     headingBase + float64(i)*headingStep,
			ObservedAt:     t,
			ReceivedAt:     t,
			Quality:        0.95,
			Verified:       false,
			Synthetic:      true,
			SimulationNote: "EDUCATIONAL SIMULATION - SYNTHETIC DATA",
			Seen:           time.Since(t).Seconds(),
		})
	}
	
	// Return slice (caller should return to pool after use)
	// For now, create a copy to avoid pool issues
	result := make([]models.SimulatedObservation, len(obs))
	copy(result, obs)
	
	// Return to pool
	e.obsPool.Put(obs)
	
	return result, nil
}

func (e *Engine) replayFrame() ([]models.SimulatedObservation, error) {
	return nil, nil // Return empty instead of error for graceful handling
}

func (e *Engine) liveFrame() ([]models.SimulatedObservation, error) {
	return nil, nil // Return empty instead of error for graceful handling
}

// Fast random source using xorshift algorithm (faster than LCG)
type randSource struct {
	state uint64
}

func randFromSeed(s int64) *randSource {
	return &randSource{state: uint64(s) & 0xFFFFFFFFFFFFFFFF}
}

func (r *randSource) Float64() float64 {
	r.state ^= r.state << 13
	r.state ^= r.state >> 7
	r.state ^= r.state << 17
	return float64(r.state) / float64(0xFFFFFFFFFFFFFFFF)
}

// Stats returns engine performance statistics
func (e *Engine) Stats() map[string]interface{} {
	return map[string]interface{}{
		"mode":            e.Mode,
		"synthetic_count": e.SyntheticCount,
		"rate_hz":         e.RateHz,
		"frame_count":     e.frameCount,
		"uptime_sec":      time.Since(e.lastFrame).Seconds(),
		"seed":            e.Seed,
	}
}
