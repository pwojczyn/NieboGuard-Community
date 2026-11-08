package collectors

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pwojczyn/NieboGuard-Community/node/internal/storage"
	"github.com/pwojczyn/NieboGuard-Community/node/pkg/config"
)

func testManager(t *testing.T, source config.SourceConfig) (*Manager, *storage.Storage) {
	t.Helper()
	store, err := storage.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Initialize(); err != nil {
		t.Fatal(err)
	}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	m, err := NewManager(&config.Configuration{Sources: []config.SourceConfig{source}}, store, logger)
	if err != nil {
		t.Fatal(err)
	}
	return m, store
}

func waitForAircraft(t *testing.T, store *storage.Storage) []storage.AircraftEvent {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		events, err := store.GetAircraft(10, 0)
		if err != nil {
			t.Fatal(err)
		}
		if len(events) > 0 {
			return events
		}
		time.Sleep(25 * time.Millisecond)
	}
	t.Fatal("timed out waiting for aircraft")
	return nil
}

func TestConfiguredADSBJSONCollectorStoresAircraft(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"now":1,"aircraft":[{"hex":"abc123","flight":"NG123","lat":52.1,"lon":21.0,"alt_baro":12000,"gs":420,"track":90,"seen":0.2}]}`)
	}))
	defer server.Close()

	m, store := testManager(t, config.SourceConfig{ID: "local_adsb", Type: "adsb_json", Name: "Test ADS-B", URL: server.URL, IntervalSeconds: 1, Enabled: true})
	if err := m.StartAll(); err != nil {
		t.Fatal(err)
	}
	defer m.StopAll()
	events := waitForAircraft(t, store)
	if events[0].Hex != "abc123" || events[0].Lat != 52.1 || events[0].SourceID != "local_adsb" {
		t.Fatalf("unexpected event: %+v", events[0])
	}
	source, ok := m.GetSource("local_adsb")
	if !ok || source.Status != "active" || source.Count == 0 {
		t.Fatalf("unexpected source: %+v ok=%v", source, ok)
	}
}

func TestConfiguredHTTPJSONCollectorStoresAircraft(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `[{"hex":"def456","flight":"NG456","lat":50.2,"lon":19.1,"alt_baro":8000,"gs":300,"track":180,"seen":0.4}]`)
	}))
	defer server.Close()

	m, store := testManager(t, config.SourceConfig{ID: "http_json", Type: "http_json", Name: "Test HTTP", URL: server.URL, IntervalSeconds: 1, Enabled: true})
	if err := m.StartAll(); err != nil {
		t.Fatal(err)
	}
	defer m.StopAll()
	events := waitForAircraft(t, store)
	if events[0].Hex != "def456" || events[0].Lon != 19.1 || events[0].SourceType != "http_json" {
		t.Fatalf("unexpected event: %+v", events[0])
	}
}

func TestConfiguredSourceRejectsNonHTTPURL(t *testing.T) {
	m, _ := testManager(t, config.SourceConfig{ID: "local_adsb", Type: "adsb_json", URL: "file:///tmp/aircraft.json", Enabled: true})
	if err := m.StartAll(); err != nil {
		t.Fatal(err)
	}
	defer m.StopAll()
	source, ok := m.GetSource("local_adsb")
	if !ok || source.Status != "inactive" {
		t.Fatalf("expected inactive source, got %+v ok=%v", source, ok)
	}
}
