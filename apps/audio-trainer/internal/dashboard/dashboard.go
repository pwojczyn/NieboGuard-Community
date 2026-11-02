package dashboard

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/pwojczyn/NieboGuard-Community/apps/audio-trainer/internal/config"
)

type Dashboard struct {
	cfg      *config.Configuration
	upTime   time.Time
	wsConns  map[interface{}]bool
	mu       sync.RWMutex
}

func New(cfg *config.Configuration) (*Dashboard, error) {
	d := &Dashboard{
		cfg:      cfg,
		upTime:   time.Now(),
		wsConns:  make(map[interface{}]bool),
	}
	return d, nil
}

func (d *Dashboard) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/", d.handleIndex).Methods("GET")
	r.HandleFunc("/api/v1/health", d.handleHealth).Methods("GET")
	r.HandleFunc("/api/v1/ready", d.handleReady).Methods("GET")
}

func (d *Dashboard) handleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Audio AI Trainer Dashboard\n")
}

func (d *Dashboard) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "{\"status\": \"healthy\"}")
}

func (d *Dashboard) handleReady(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "{\"status\": \"ready\"}")
}