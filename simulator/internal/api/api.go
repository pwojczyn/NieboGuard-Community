package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pwojczyn/NieboGuard-Community/simulator/pkg/models"
)

// Server provides simulator control API.

type Server struct {
	httpHandler *http.ServeMux
}

func New() *Server {
	s := &Server{httpHandler: http.NewServeMux()}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.httpHandler.HandleFunc("/sim/api/health", s.healthHandler)
	s.httpHandler.HandleFunc("/sim/api/status", s.statusHandler)
	s.httpHandler.HandleFunc("/sim/api/objects", s.objectsHandler)
	s.httpHandler.HandleFunc("/sim/api/run", s.runHandler)
	s.httpHandler.HandleFunc("/sim/api/stop", s.stopHandler)
	s.httpHandler.HandleFunc("/sim/api/scenarios", s.scenariosHandler)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.httpHandler.ServeHTTP(w, r)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"service": "nieboguard-simulator",
		"version": "0.1.0",
	})
}

func (s *Server) statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"mode":      "replay",
		"speed":     1.0,
		"running":   true,
		"objects":   0,
		"simulated": true,
	})
}

func (s *Server) objectsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode([]models.SimulatedObservation{
		{
			ID: "test-01", SourceID: "sim-01", SourceType: "synthetic",
			ObservationID: "sim-obs-01", ObjectType: "aircraft",
			Hex: "TEST01", Lat: 52.1, Lon: 19.2, AltM: 10000,
			Synthetic: true, SimulationNote: "TEST DATA",
			ObservedAt: time.Now(),
		},
	})
}

func (s *Server) runHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "started",
		"mode": r.URL.Query().Get("mode"),
		"simulated": true,
	})
}

func (s *Server) stopHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"status": "stopped"})
}

func (s *Server) scenariosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode([]string{"replay", "synthetic", "live_transform"})
}

func (s *Server) Health() string {
	return fmt.Sprintf("simulator OK %v", time.Now().Format("2006-01-02"))
}
