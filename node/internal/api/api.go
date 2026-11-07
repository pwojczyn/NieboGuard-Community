package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/pwojczyn/NieboGuard-Community/node/internal/collectors"
	"github.com/pwojczyn/NieboGuard-Community/node/internal/storage"
	"github.com/pwojczyn/NieboGuard-Community/node/pkg/config"
)

// Server represents the HTTP API server.
type Server struct {
	cfg          *config.Configuration
	storage      *storage.Storage
	collectorMgr *collectors.Manager
	logger       *slog.Logger
	server       *http.Server
	security     *securityConfig
}

// New creates a new API server.
func New(cfg *config.Configuration, store *storage.Storage, collMgr *collectors.Manager, logger *slog.Logger) (*Server, error) {
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	if !loopbackHost(cfg.Host) && firstEnv("NNODE_AUTH_TOKEN", "NB_NODE_AUTH_TOKEN") == "" {
		return nil, fmt.Errorf("refusing non-loopback host %q without NNODE_AUTH_TOKEN", cfg.Host)
	}
	router := mux.NewRouter()
	s := &Server{
		cfg:          cfg,
		storage:      store,
		collectorMgr: collMgr,
		logger:       logger,
		security:     newSecurityConfig(),
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Handler: router,
		},
	}

	// Setup routes
	s.setupRoutes(router)
	router.Use(s.securityMiddleware)

	return s, nil
}

// Start starts the HTTP server.
func (s *Server) Start(port int, host string) error {
	addr := fmt.Sprintf("%s:%d", host, port)
	s.logger.Info("starting api server", "addr", addr)

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("api server failed", "error", err)
		}
	}()
	return nil
}

// Stop stops the HTTP server.
func (s *Server) Stop() {
	if s.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.server.Shutdown(ctx)
	}
}

// HealthHandler handles health checks.
func (s *Server) HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handleError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"status":  "ok",
		"service": "nieboguard-node",
		"version": "0.1.0",
	}
	json.NewEncoder(w).Encode(response)
}

// ReadyHandler handles readiness checks.
func (s *Server) ReadyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handleError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"status":   "ready",
		"database": "connected",
		"sources":  len(s.collectorMgr.GetSources()),
	}
	json.NewEncoder(w).Encode(response)
}

// ObjectsHandler handles aircraft objects.
func (s *Server) ObjectsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handleError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Parse query parameters
	query := r.URL.Query()
	limit := 100
	if l := query.Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	offset := 0
	if o := query.Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	events, err := s.storage.GetAircraft(limit, offset)
	if err != nil {
		s.logger.Error("failed to get aircraft", "error", err)
		handleError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"objects": events,
		"total":   len(events),
		"limit":   limit,
		"offset":  offset,
	})
}

// SourcesHandler returns information about data sources.
func (s *Server) SourcesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handleError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	sources := s.collectorMgr.GetSources()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sources)
}

// EventsHandler returns aggregated events.
func (s *Server) EventsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handleError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	query := r.URL.Query()
	limit := 50
	if l := query.Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	events, err := s.storage.GetAircraft(limit, 0)
	if err != nil {
		s.logger.Error("failed to get events", "error", err)
		handleError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"events": events,
		"total":  len(events),
		"limit":  limit,
	})
}

// NoticesHandler returns safety notices.
func (s *Server) NoticesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handleError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"notices": []interface{}{},
		"total":   0,
	})
}

// StatsHandler returns system statistics.
func (s *Server) StatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handleError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	stats := map[string]interface{}{
		"uptime":      "0h 0m",
		"sources":     len(s.collectorMgr.GetSources()),
		"objects":     0,
		"connections": 0,
		"version":     "0.1.0",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// AdminSourcesHandler handles source administration.
func (s *Server) AdminSourcesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handleError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "ok",
		"message":   "admin operation completed",
		"timestamp": time.Now().Unix(),
	})
}

var sha256Pattern = regexp.MustCompile(`^[a-f0-9]{64}$`)

// OfflinePacksHandler returns verified offline manifest metadata.
func (s *Server) OfflinePacksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handleError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	packs, err := s.storage.ListOfflinePacks(id, limit)
	if err != nil {
		s.logger.Error("listing offline packs", "error", err)
		handleError(w, http.StatusInternalServerError, "internal error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"packs": packs, "total": len(packs)})
}

// CurrentOfflinePackHandler returns the newest version of one pack.
func (s *Server) CurrentOfflinePackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handleError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	id := strings.TrimSpace(mux.Vars(r)["id"])
	if id == "" {
		handleError(w, http.StatusBadRequest, "pack id required")
		return
	}
	pack, err := s.storage.CurrentOfflinePack(id)
	if err != nil {
		s.logger.Error("reading offline pack", "error", err)
		handleError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if pack == nil {
		handleError(w, http.StatusNotFound, "pack not found")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pack)
}

// AdminOfflinePackHandler registers manifest metadata. It never accepts a file path.
func (s *Server) AdminOfflinePackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handleError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var body struct {
		ID           string `json:"id"`
		Version      string `json:"version"`
		Sha256       string `json:"sha256"`
		Bytes        int64  `json:"bytes"`
		CreatedAt    string `json:"created_at"`
		Mode         string `json:"mode"`
		Sharing      string `json:"sharing"`
		Status       string `json:"status"`
		ManifestJSON string `json:"manifest_json"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64*1024)).Decode(&body); err != nil {
		handleError(w, http.StatusBadRequest, "invalid json")
		return
	}
	body.ID = strings.TrimSpace(body.ID)
	body.Version = strings.TrimSpace(body.Version)
	body.Sha256 = strings.ToLower(strings.TrimSpace(body.Sha256))
	if body.ID == "" || body.Version == "" || !sha256Pattern.MatchString(body.Sha256) || body.Bytes < 0 {
		handleError(w, http.StatusBadRequest, "invalid manifest fields")
		return
	}
	if body.Mode == "" {
		body.Mode = "UNVERIFIED"
	}
	if body.Sharing == "" {
		body.Sharing = "LOCAL_ONLY"
	}
	if body.Mode == "SIMULATED" && body.Sharing != "LOCAL_ONLY" {
		handleError(w, http.StatusBadRequest, "simulated pack must be LOCAL_ONLY")
		return
	}
	if body.Status == "" {
		body.Status = "AVAILABLE"
	}
	if body.CreatedAt == "" {
		body.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}
	if body.ManifestJSON == "" {
		body.ManifestJSON = "{}"
	}
	pack := &storage.OfflinePack{ID: body.ID, Version: body.Version, Sha256: body.Sha256, Bytes: body.Bytes, CreatedAt: body.CreatedAt, ManifestJSON: body.ManifestJSON, Mode: body.Mode, Sharing: body.Sharing, Status: body.Status, UpdatedAt: time.Now().UTC().Format(time.RFC3339)}
	if err := s.storage.UpsertOfflinePack(pack); err != nil {
		s.logger.Error("upserting offline pack", "error", err)
		handleError(w, http.StatusInternalServerError, "internal error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"ok": true, "pack": pack})
}

// setupRoutes configures the HTTP routes.
func (s *Server) setupRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/health", s.HealthHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/ready", s.ReadyHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"version": "0.1.0",
		})
	}).Methods(http.MethodGet)

	router.HandleFunc("/api/v1/sources", s.SourcesHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/objects", s.ObjectsHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/events", s.EventsHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/notices", s.NoticesHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/stats", s.StatsHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/offline/packs", s.OfflinePacksHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/offline/packs/{id}", s.CurrentOfflinePackHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/admin/offline/packs", s.AdminOfflinePackHandler).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/admin/sources", s.AdminSourcesHandler).Methods(http.MethodPost)
}

// handleError sends an error response.
func handleError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]interface{}{
			"code":    http.StatusText(status),
			"message": message,
		},
	})
}
