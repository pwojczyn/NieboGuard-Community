package tests

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/pwojczyn/NieboGuard-Community/simulator/internal/military"
	"github.com/pwojczyn/NieboGuard-Community/simulator/pkg/models"
)

// TestLink16Simulator_GenerateTrack verifies that GenerateTrack produces valid data.
func TestLink16Simulator_GenerateTrack(t *testing.T) {
	sim := military.NewLink16Simulator("test-link16-01")
	event, err := sim.GenerateTrack()

	if err != nil {
		t.Fatalf("GenerateTrack returned error: %v", err)
	}

	if event == nil {
		t.Fatal("GenerateTrack returned nil event")
	}

	// Verify all required fields
	if event.SourceID != "test-link16-01" {
		t.Errorf("Expected SourceID 'test-link16-01', got '%s'", event.SourceID)
	}

	if !event.Synthetic {
		t.Error("Expected Synthetic=true")
	}

	if event.SourceType != "sim-link16" {
		t.Errorf("Expected SourceType 'sim-link16', got '%s'", event.SourceType)
	}

	if event.Hex == "" {
		t.Error("Expected Hex to be set")
	}

	if event.Flight == "" {
		t.Error("Expected Flight to be set")
	}

	if event.Lat == 0 && event.Lon == 0 {
		t.Error("Expected Lat/Lon to be set")
	}

	if event.AltM <= 0 {
		t.Error("Expected AltM > 0")
	}

	if event.Gs <= 0 {
		t.Error("Expected Gs > 0")
	}

	if event.Track < 0 || event.Track > 360 {
		t.Error("Expected Track in range 0-360")
	}

	if event.Synthetic != true {
		t.Error("Expected Synthetic=true")
	}

	if event.SimulationNote == "" {
		t.Error("Expected SimulationNote to be set")
	}

	// Verify simulation note contains educational marker
	if !contains(event.SimulationNote, "EDUCATIONAL") {
		t.Errorf("Expected SimulationNote to contain 'EDUCATIONAL', got '%s'", event.SimulationNote)
	}

	if !contains(event.SimulationNote, "SIMULATION") {
		t.Errorf("Expected SimulationNote to contain 'SIMULATION', got '%s'", event.SimulationNote)
	}
}

// TestLink16Simulator_GenerateMultipleTracks verifies GenerateMultipleTracks works.
func TestLink16Simulator_GenerateMultipleTracks(t *testing.T) {
	sim := military.NewLink16Simulator("test-link16-02")
	tracks, err := sim.GenerateMultipleTracks(5)

	if err != nil {
		t.Fatalf("GenerateMultipleTracks returned error: %v", err)
	}

	if len(tracks) != 5 {
		t.Errorf("Expected 5 tracks, got %d", len(tracks))
	}

	for i, track := range tracks {
		if track.SourceID != "test-link16-02" {
			t.Errorf("Track %d: Expected SourceID 'test-link16-02', got '%s'", i, track.SourceID)
		}

		if !track.Synthetic {
			t.Errorf("Track %d: Expected Synthetic=true", i)
		}

		if track.SourceType != "sim-link16" {
			t.Errorf("Track %d: Expected SourceType 'sim-link16', got '%s'", i, track.SourceType)
		}

		if track.Lat == 0 && track.Lon == 0 {
			t.Errorf("Track %d: Expected Lat/Lon to be set", i)
		}

		if track.AltM <= 0 {
			t.Errorf("Track %d: Expected AltM > 0", i)
		}

		if track.Gs <= 0 {
			t.Errorf("Track %d: Expected Gs > 0", i)
		}

		if track.Track < 0 || track.Track > 360 {
			t.Errorf("Track %d: Expected Track in range 0-360", i)
		}

		if track.SimulationNote == "" {
			t.Errorf("Track %d: Expected SimulationNote to be set", i)
		}
	}
}

// TestLink16Simulator_GenerateScenario verifies scenario generation.
func TestLink16Simulator_GenerateScenario(t *testing.T) {
	sim := military.NewLink16Simulator("test-link16-03")

	// Test formation scenario
	tracks, err := sim.GenerateScenario("formation")
	if err != nil {
		t.Fatalf("GenerateScenario 'formation' returned error: %v", err)
	}

	if len(tracks) != 3 {
		t.Errorf("Expected 3 tracks for 'formation', got %d", len(tracks))
	}

	for i, track := range tracks {
		if track.SimulationNote == "" {
			t.Errorf("Track %d: Expected SimulationNote", i)
		}
	}

	// Test patrol scenario
	tracks, err = sim.GenerateScenario("patrol")
	if err != nil {
		t.Fatalf("GenerateScenario 'patrol' returned error: %v", err)
	}

	if len(tracks) != 1 {
		t.Errorf("Expected 1 track for 'patrol', got %d", len(tracks))
	}

	if tracks[0].Flight != "PATROL-001" {
		t.Errorf("Expected Flight 'PATROL-001', got '%s'", tracks[0].Flight)
	}

	// Test intercept scenario
	tracks, err = sim.GenerateScenario("intercept")
	if err != nil {
		t.Fatalf("GenerateScenario 'intercept' returned error: %v", err)
	}

	if len(tracks) != 2 {
		t.Errorf("Expected 2 tracks for 'intercept', got %d", len(tracks))
	}

	// Default scenario (5 tracks)
	tracks, err = sim.GenerateScenario("")
	if err != nil {
		t.Fatalf("GenerateScenario '' returned error: %v", err)
	}

	if len(tracks) != 5 {
		t.Errorf("Expected 5 tracks for default, got %d", len(tracks))
	}
}

// TestLink16Simulator_GenerateTrackPacket verifies packet generation.
func TestLink16Simulator_GenerateTrackPacket(t *testing.T) {
	sim := military.NewLink16Simulator("test-link16-04")
	packet, err := sim.GenerateTrackPacket()

	if err != nil {
		t.Fatalf("GenerateTrackPacket returned error: %v", err)
	}

	if packet == nil {
		t.Fatal("GenerateTrackPacket returned nil packet")
	}

	// Parse the JSON packet
	var event models.AircraftEvent
	if err := json.Unmarshal(packet, &event); err != nil {
		t.Fatalf("Failed to parse packet JSON: %v", err)
	}

	// Verify fields
	if event.SourceID != "test-link16-04" {
		t.Errorf("Expected SourceID 'test-link16-04', got '%s'", event.SourceID)
	}

	if !event.Synthetic {
		t.Error("Expected Synthetic=true")
	}

	if event.SimulationNote == "" {
		t.Error("Expected SimulationNote to be set")
	}

	if !contains(event.SimulationNote, "EDUCATIONAL") {
		t.Errorf("Expected SimulationNote to contain 'EDUCATIONAL', got '%s'", event.SimulationNote)
	}
}

// TestLink16Simulator_Stats verifies Stats output.
func TestLink16Simulator_Stats(t *testing.T) {
	sim := military.NewLink16Simulator("test-link16-05")
	stats := sim.Stats()

	if stats == nil {
		t.Fatal("Stats returned nil")
	}

	if _, ok := stats["source_id"]; !ok {
		t.Error("Expected stats to contain 'source_id'")
	}

	if _, ok := stats["simulated"]; !ok {
		t.Error("Expected stats to contain 'simulated'")
	}

	if stats["source_id"] != "test-link16-05" {
		t.Errorf("Expected source_id 'test-link16-05', got '%v'", stats["source_id"])
	}

	if stats["simulated"] != true {
		t.Error("Expected simulated=true")
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestLink16Simulator_UnmarshalJSON verifies JSON unmarshaling
func TestLink16Simulator_UnmarshalJSON(t *testing.T) {
	// This test verifies the JSON parsing structure
	// It's based on the GenerateTrackPacket output format
	
	testJSON := `{
		"source_id": "test-link16-01",
		"source_type": "sim-link16",
		"observation_id": "link16-obs-001",
		"object_type": "aircraft",
		"hex": "SIM0123",
		"flight": "LINK16-SIM-001",
		"lat": 52.1,
		"lon": 19.2,
		"alt_m": 10000.0,
		"gs": 450.0,
		"track": 90.0,
		"seen": 1.0,
		"synthetic": true,
		"simulation_note": "EDUCATIONAL SIMULATION - PUBLIC SPECIFICATIONS ONLY",
		"observed_at": "2026-01-01T12:00:00Z",
		"received_at": "2026-01-01T12:00:00Z",
		"quality": 0.98,
		"verified": false
	}`
	
	var event models.AircraftEvent
	if err := json.Unmarshal([]byte(testJSON), &event); err != nil {
		t.Fatalf("Failed to unmarshal test JSON: %v", err)
	}
	
	if event.SourceID != "test-link16-01" {
		t.Errorf("Expected SourceID 'test-link16-01', got '%s'", event.SourceID)
	}
	
	if !event.Synthetic {
		t.Error("Expected Synthetic=true")
	}
	
	if event.SimulationNote != "EDUCATIONAL SIMULATION - PUBLIC SPECIFICATIONS ONLY" {
		t.Errorf("Expected SimulationNote 'EDUCATIONAL SIMULATION - PUBLIC SPECIFICATIONS ONLY', got '%s'", event.SimulationNote)
	}
}
