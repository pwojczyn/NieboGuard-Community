package military

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/pwojczyn/NieboGuard-Community/simulator/pkg/models"
)

// Link16Simulator generates synthetic Link 16-style data for educational purposes.
type Link16Simulator struct {
	SourceID      string
	Simulated     bool
	MinAlt        int
	MaxAlt        int
	MinSpeed      int
	MaxSpeed      int
	MinLat        float64
	MinLon        float64
	Seed          int64
	RateHz        int
}

// NewLink16Simulator creates a new Link 16 simulator instance.
func NewLink16Simulator(sourceID string) *Link16Simulator {
	return &Link16Simulator{
		SourceID:    sourceID,
		Simulated:   true,
		MinAlt:      5000,
		MaxAlt:      40000,
		MinSpeed:    100,
		MaxSpeed:    800,
		MinLat:      51.5, // Poland region minimum
		MinLon:      18.0, // Poland region minimum
		Seed:        time.Now().UnixNano(),
		RateHz:      1,
	}
}

// GenerateTrack returns a synthetic Link 16 track object.
func (s *Link16Simulator) GenerateTrack() (*models.AircraftEvent, error) {
	// Generate random values within specified ranges
	lat := randFloat64(-0.5, 0.5) + 52.0 // Poland region
	lon := randFloat64(-0.5, 0.5) + 19.0  // Poland region
	alt := rand.Intn(s.MaxAlt-s.MinAlt) + s.MinAlt
	speed := rand.Intn(s.MaxSpeed-s.MinSpeed) + s.MinSpeed
	heading := rand.Intn(360)

	// Generate synthetic track ID (SIM-XXXXX format)
	trackID := fmt.Sprintf("SIM-%05d", rand.Intn(100000))

	// Generate timestamp
	timestamp := time.Now().UTC().Format(time.RFC3339)

	// Create AircraftEvent with all required fields
	event := &models.AircraftEvent{
		SourceID:      s.SourceID,
		SourceType:    "sim-link16",
		ObservationID: fmt.Sprintf("link16-obs-%d", time.Now().Unix()),
		ObjectType:    "aircraft",
		Hex:           fmt.Sprintf("SIM%04X", rand.Intn(10000)),
		Flight:        fmt.Sprintf("LINK16-SIM-%03d", rand.Intn(100)),
		Lat:           lat,
		Lon:           lon,
		AltM:          float64(alt),
		Gs:            float64(speed),
		Track:         float64(heading),
		Seen:          1.0,
		Synthetic:     true,
		SimulationNote: "EDUCATIONAL SIMULATION - PUBLIC SPECIFICATIONS ONLY",
		ObservedAt:    time.Now(),
		ReceivedAt:    time.Now(),
		Quality:       0.98,
		Verified:      false,
	}

	return event, nil
}

// GenerateMultipleTracks generates multiple synthetic tracks.
func (s *Link16Simulator) GenerateMultipleTracks(count int) ([]models.AircraftEvent, error) {
	if count <= 0 {
		return nil, fmt.Errorf("count must be positive")
	}

	tracks := make([]models.AircraftEvent, 0, count)
	rng := rand.New(rand.NewSource(s.Seed))

	for i := 0; i < count; i++ {
		lat := s.MinLat + (rng.Float64()*1.0) // ±1 degree
		lon := s.MinLon + (rng.Float64()*1.0)  // ±1 degree
		alt := rand.Intn(s.MaxAlt-s.MinAlt) + s.MinAlt
		speed := rand.Intn(s.MaxSpeed-s.MinSpeed) + s.MinSpeed
		heading := rand.Intn(360)

		trackID := fmt.Sprintf("SIM-%05d", rng.Intn(100000))

		event := &models.AircraftEvent{
			SourceID:      s.SourceID,
			SourceType:    "sim-link16",
			ObservationID: fmt.Sprintf("link16-obs-%d", time.Now().Unix() + i),
			ObjectType:    "aircraft",
			Hex:           fmt.Sprintf("SIM%04X", rng.Intn(10000)),
			Flight:        fmt.Sprintf("SIM-LINK16-%03d", i),
			Lat:           lat,
			Lon:           lon,
			AltM:          float64(alt),
			Gs:            float64(speed),
			Track:         float64(heading),
			Seen:          1.0,
			Synthetic:     true,
			SimulationNote: "EDUCATIONAL SIMULATION - PUBLIC SPECIFICATIONS ONLY",
			ObservedAt:    time.Now(),
			ReceivedAt:    time.Now(),
			Quality:       0.98,
			Verified:      false,
		}

		tracks = append(tracks, *event)
	}

	return tracks, nil
}

// GenerateTrackPacket returns a synthetic Link 16 message packet.
func (s *Link16Simulator) GenerateTrackPacket() ([]byte, error) {
	// This is a simplified representation - actual Link 16 is binary
	track := &models.AircraftEvent{
		SourceID:      s.SourceID,
		SourceType:    "sim-link16",
		ObservationID: "sim-obs-001",
		ObjectType:    "aircraft",
		Hex:           "SIM0123",
		Flight:        "LINK16-SIM-001",
		Lat:           52.1,
		Lon:           19.2,
		AltM:          10000.0,
		Gs:            450.0,
		Track:         90.0,
		Seen:          1.0,
		Synthetic:     true,
		SimulationNote: "EDUCATIONAL SIMULATION - PUBLIC SPECIFICATIONS ONLY",
		ObservedAt:    time.Now(),
		ReceivedAt:    time.Now(),
		Quality:       0.98,
		Verified:      false,
	}

	// Return as JSON for educational purposes (actual Link 16 is binary)
	return json.Marshal(track)
}

// GenerateScenario generates a specific simulation scenario.
func (s *Link16Simulator) GenerateScenario(scenario string) ([]models.AircraftEvent, error) {
	switch scenario {
	case "formation":
		return s.generateFormationScenario()
	case "patrol":
		return s.generatePatrolScenario()
	case "intercept":
		return s.generateInterceptScenario()
	default:
		return s.GenerateMultipleTracks(5)
	}
}

func (s *Link16Simulator) generateFormationScenario() ([]models.AircraftEvent, error) {
	// Simple formation pattern
	tracks := make([]models.AircraftEvent, 3)
	rng := rand.New(rand.NewSource(s.Seed + 1))

	for i := 0; i < 3; i++ {
		tracks[i] = &models.AircraftEvent{
			SourceID:      s.SourceID,
			SourceType:    "sim-link16",
			ObservationID: fmt.Sprintf("link16-obs-formation-%d", i),
			ObjectType:    "aircraft",
			Hex:           fmt.Sprintf("SIM%04X", i+1),
			Flight:        fmt.Sprintf("FORM-%d", i),
			Lat:           52.1 + float64(i*0.001),
			Lon:           19.2 + float64(i*0.001),
			AltM:          10000.0 + float64(i*100),
			Gs:            400.0 + float64(i*20),
			Track:         90.0 + float64(i*10),
			Seen:          1.0,
			Synthetic:     true,
			SimulationNote: "EDUCATIONAL SIMULATION - PUBLIC SPECIFICATIONS ONLY",
			ObservedAt:    time.Now(),
			ReceivedAt:    time.Now(),
			Quality:       0.98,
			Verified:      false,
		}
	}
	return tracks, nil
}

func (s *Link16Simulator) generatePatrolScenario() ([]models.AircraftEvent, error) {
	tracks := make([]models.AircraftEvent, 1)
	rng := rand.New(rand.NewSource(s.Seed + 2))

	tracks[0] = &models.AircraftEvent{
		SourceID:      s.SourceID,
		SourceType:    "sim-link16",
		ObservationID: "link16-obs-patrol-001",
		ObjectType:    "aircraft",
		Hex:           "SIM0123",
		Flight:        "PATROL-001",
		Lat:           52.1,
		Lon:           19.2,
		AltM:          10000.0,
		Gs:            300.0,
		Track:         0.0,
		Seen:          1.0,
		Synthetic:     true,
		SimulationNote: "EDUCATIONAL SIMULATION - PUBLIC SPECIFICATIONS ONLY",
		ObservedAt:    time.Now(),
		ReceivedAt:    time.Now(),
		Quality:       0.98,
		Verified:      false,
	}

	return tracks, nil
}

func (s *Link16Simulator) generateInterceptScenario() ([]models.AircraftEvent, error) {
	tracks := make([]models.AircraftEvent, 2)
	rng := rand.New(rand.NewSource(s.Seed + 3))

	for i := 0; i < 2; i++ {
		tracks[i] = &models.AircraftEvent{
			SourceID:      s.SourceID,
			SourceType:    "sim-link16",
			ObservationID: fmt.Sprintf("link16-obs-intercept-%d", i),
			ObjectType:    "aircraft",
			Hex:           fmt.Sprintf("SIM%04X", i+1),
			Flight:        fmt.Sprintf("INTERCEPT-%d", i),
			Lat:           52.0 + float64(i*0.001),
			Lon:           19.1 + float64(i*0.001),
			AltM:          8000.0 + float64(i*500),
			Gs:            500.0 + float64(i*100),
			Track:         270.0 + float64(i*45),
			Seen:          1.0,
			Synthetic:     true,
			SimulationNote: "EDUCATIONAL SIMULATION - PUBLIC SPECIFICATIONS ONLY",
			ObservedAt:    time.Now(),
			ReceivedAt:    time.Now(),
			Quality:       0.98,
			Verified:      false,
		}
	}
	return tracks, nil
}

// Stats returns simulation statistics
func (s *Link16Simulator) Stats() map[string]interface{} {
	return map[string]interface{}{
		"source_id":     s.SourceID,
		"simulated":     s.Simulated,
		"total_tracks":  0, // Will be updated externally
		"min_alt":       s.MinAlt,
		"max_alt":       s.MaxAlt,
		"min_speed":     s.MinSpeed,
		"max_speed":     s.MaxSpeed,
		"seed":          s.Seed,
		"rate_hz":       s.RateHz,
	}
}