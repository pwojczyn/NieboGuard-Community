package asterix

import (
	"testing"
	"time"

	"github.com/pwojczyn/NieboGuard-Community/simulator/pkg/models"
)

func cat062Event() *models.AircraftEvent {
	return &models.AircraftEvent{
		SourceID: "sim-station", SourceType: "radar", ObservationID: "obs-1",
		ObjectType: "aircraft", Hex: "abc123", Flight: "NGS001",
		Lat: 52.1, Lon: 21.05, AltM: 1000, Gs: 100, HeadingDeg: 90,
		Synthetic: true, SimulationNote: "educational simulation",
	}
}

func TestCAT062TrackLifecycle(t *testing.T) {
	e := &CAT062Encoder{DataSourceID: 7}
	first, err := e.Encode(cat062Event())
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	if len(first) < 16 {
		t.Fatalf("record too short: %d", len(first))
	}
	if len(e.ActiveTracks) != 1 {
		t.Fatalf("expected one track, got %d", len(e.ActiveTracks))
	}
	second, err := e.Encode(cat062Event())
	if err != nil {
		t.Fatalf("second encode: %v", err)
	}
	if len(e.ActiveTracks) != 1 {
		t.Fatalf("same object must reuse track, got %d", len(e.ActiveTracks))
	}
	if string(first[3:5]) != string(second[3:5]) {
		t.Fatal("track number should be stable")
	}
	other := cat062Event()
	other.Flight = "NGS002"
	if _, err := e.Encode(other); err != nil {
		t.Fatalf("other encode: %v", err)
	}
	if len(e.ActiveTracks) != 2 {
		t.Fatalf("expected two tracks, got %d", len(e.ActiveTracks))
	}
}

func TestCAT062CleanupExpired(t *testing.T) {
	e := &CAT062Encoder{DataSourceID: 7}
	event := cat062Event()
	if _, err := e.Encode(event); err != nil {
		t.Fatal(err)
	}
	for _, track := range e.ActiveTracks {
		track.LastUpdate = time.Now().Add(-2 * time.Hour)
	}
	e.CleanupExpired(time.Hour)
	if len(e.ActiveTracks) != 0 {
		t.Fatalf("expired tracks should be removed, got %d", len(e.ActiveTracks))
	}
}
