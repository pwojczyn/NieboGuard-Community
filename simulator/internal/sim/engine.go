package sim

import (
	"fmt"
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
	Mode          Mode
	Speed         float64 // 0.25, 1, 2, 5, 10
	Loop          bool
	SyntheticCount int
	StartLat     float64
	StartLon     float64
	DirDeg       float64 // direction
	SpeedMps     float64
	AltM         float64
	RateHz       int // updates per second
	Seed         int64
}

func NewEngine(mode Mode) *Engine {
	return &Engine{
		Mode:    mode,
		Speed:   1.0,
		Loop:    false,
		SyntheticCount: 10,
		StartLat: 52.2297,
		StartLon: 21.0122,
		DirDeg:   90,
		SpeedMps: 250,
		AltM:     10000,
		RateHz:   1,
		Seed:     42,
	}
}

func (e *Engine) NextFrame() ([]models.SimulatedObservation, error) {
	switch e.Mode {
	case ModeSynthetic:
		return e.syntheticFrame()
	case ModeReplay:
		return e.replayFrame()
	case ModeLiveTransform:
		return e.liveFrame()
	default:
		return nil, fmt.Errorf("unknown mode")
	}
}

func (e *Engine) syntheticFrame() ([]models.SimulatedObservation, error) {
	obs := make([]models.SimulatedObservation, 0, e.SyntheticCount)
	rng := randFromSeed(e.Seed)
	now := time.Now()

	for i := 0; i < e.SyntheticCount; i++ {
		t := now.Add(-time.Duration(i) * time.Second)
		lat := e.StartLat + (float64(i) * 0.001) + (rng.Float64() - 0.5)*0.01
		lon := e.StartLon + (float64(i) * 0.001) + (rng.Float64() - 0.5)*0.01

		obs = append(obs, models.SimulatedObservation{
			ID:           fmt.Sprintf("SIM-%04d", i),
			SourceID:     "sim-sensor-01",
			SourceType:   "synthetic",
			ObservationID: fmt.Sprintf("sim-obs-%d-%d", now.Unix(), i),
			ObjectType:   "aircraft",
			Hex:          fmt.Sprintf("SIM%04X", i),
			Flight:       fmt.Sprintf("TEST%03d", i),
			Lat:          lat,
			Lon:          lon,
			AltM:         10000 + float64(i)*50 + rng.Float64()*500,
			SpeedMps:     250 + rng.Float64()*50,
			HeadingDeg:   90 + float64(i)*10,
			ObservedAt:   t,
			ReceivedAt:   t,
			Quality:      0.95,
			Verified:     false,
			Synthetic:    true,
			SimulationNote: "TEST DATA — NOT REAL",
		})
	}
	return obs, nil
}

func (e *Engine) replayFrame() ([]models.SimulatedObservation, error) {
	return nil, fmt.Errorf("replay requires recorded file")
}

func (e *Engine) liveFrame() ([]models.SimulatedObservation, error) {
	return nil, fmt.Errorf("live transform requires live source")
}

func randFromSeed(s int64) *randSource {
	return &randSource{seed: s}
}

type randSource struct {
	seed int64
}

func (r *randSource) Float64() float64 {
	r.seed = (r.seed*1103515245 + 12345) & 0x7fffffff
	return float64(r.seed) / float64(0x7fffffff)
}
