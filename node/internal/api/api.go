package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
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
