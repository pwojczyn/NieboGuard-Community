package asterix

import (
	"math/rand"
	"testing"

	"github.com/pwojczyn/NieboGuard-Community/simulator/pkg/models"
)

func testEvent() *models.AircraftEvent {
	return &models.AircraftEvent{
		SourceID: "sim-station", SourceType: "radar", ObservationID: "obs-1",
		ObjectType: "aircraft", Hex: "abc123", Flight: "NGS001",
		Lat: 52.1, Lon: 21.05, AltM: 1000, Gs: 100, HeadingDeg: 90,
		Synthetic: true, SimulationNote: "educational simulation",
	}
}

func TestCAT048EncodeAndVerify(t *testing.T) {
	e := NewCAT048Encoder(52.0, 21.0, 100)
	e.Rng = rand.New(rand.NewSource(1))
	e.SetErrorPercentage(0)
	data, err := e.Encode(testEvent())
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	if len(data) < 22 {
		t.Fatalf("record too short: %d", len(data))
	}
	if err := VerifyCAT048Record(data); err != nil {
		t.Fatalf("verify: %v", err)
	}
	if !e.IsInRange(52.1, 21.05) {
		t.Fatal("expected in range")
	}
	rangeKm, azimuth, inRange := e.GetRangeInfo(52.1, 21.05)
	if !inRange || rangeKm <= 0 || azimuth < 0 || azimuth > 360 {
		t.Fatalf("range info: %v %v %v", rangeKm, azimuth, inRange)
	}
	if e.Stats()["simulated"] != true {
		t.Fatal("must be marked simulated")
	}
}

func TestCAT048RejectsOutOfRangeAndLoss(t *testing.T) {
	e := NewCAT048Encoder(0, 0, 0)
	e.MaxRangeKm = 10
	e.Rng = rand.New(rand.NewSource(1))
	if _, err := e.Encode(&models.AircraftEvent{Lat: 60, Lon: 60}); err == nil {
		t.Fatal("out-of-range should fail")
	}
	e2 := NewCAT048Encoder(0, 0, 0)
	e2.SetLossPercentage(100)
	if _, err := e2.Encode(&models.AircraftEvent{Lat: 0.01, Lon: 0.01}); err == nil {
		t.Fatal("simulated loss should fail")
	}
	if err := VerifyCAT048Record([]byte{0x01}); err == nil {
		t.Fatal("short record should fail")
	}
}
