package core

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/pwojczyn/NieboGuard-Community/node/internal/api"
	"github.com/pwojczyn/NieboGuard-Community/node/internal/collectors"
	"github.com/pwojczyn/NieboGuard-Community/node/internal/storage"
	"github.com/pwojczyn/NieboGuard-Community/node/internal/websocket"
	"github.com/pwojczyn/NieboGuard-Community/node/pkg/config"
)

// Options configures the Core Engine.
type Options struct {
	ConfigPath string
	Port       int
	Host       string
	WorkDir    string
	Verbose    bool
	Logger     *slog.Logger
}

// Core represents the NieboGuard Node engine.
type Core struct {
	opts        Options
	config      *config.Configuration
	logger      *slog.Logger
	storage     *storage.Storage
	collector   *collectors.Manager
	apiServer   *api.Server
	wssServer   *websocket.Server
	shutdown    chan struct{}
	stopped     chan struct{}
}

// New creates a new Core engine instance.
func New(opts Options) (*Core, error) {
	logger := opts.Logger
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	cfg, err := config.Load(opts.ConfigPath)
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}

	if cfg.Port == 0 {
		cfg.Port = 8787
	}
	if cfg.Host == "" {
		cfg.Host = "127.0.0.1"
	}
	if cfg.WorkDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("getting home directory: %w", err)
		}
		cfg.WorkDir = filepath.Join(home, ".nieboguard")
	}
	if cfg.WorkDir != "" {
		if err := os.MkdirAll(cfg.WorkDir, 0755); err != nil {
			return nil, fmt.Errorf("creating work directory: %w", err)
		}
	}

	// Initialize storage
	storage, err := storage.New(cfg.WorkDir)
	if err != nil {
		return nil, fmt.Errorf("creating storage: %w", err)
	}

	// Initialize collector manager
	collector, err := collectors.NewManager(cfg, storage, logger)
	if err != nil {
		return nil, fmt.Errorf("creating collector: %w", err)
	}

	// Initialize API server
	apiServer, err := api.New(cfg, storage, collector, logger)
	if err != nil {
		return nil, fmt.Errorf("creating API server: %w", err)
	}

	// Initialize WebSocket server
	wssServer, err := websocket.New(cfg, storage, collector, logger)
	if err != nil {
		return nil, fmt.Errorf("creating WebSocket server: %w", err)
	}

	core := &Core{
		opts:        opts,
		config:      cfg,
		logger:      logger,
		storage:     storage,
		collector:   collector,
		apiServer:   apiServer,
		wssServer:   wssServer,
		shutdown:    make(chan struct{}),
		stopped:     make(chan struct{}),
	}

	return core, nil
}

func (c *Core) Start() error {
	c.logger.Info("NieboGuard Node started", "port", c.config.Port, "host", c.config.Host)

	// Initialize storage
	if err := c.storage.Initialize(); err != nil {
		return fmt.Errorf("storage initialization failed: %w", err)
	}

	// Start collectors
	if err := c.collector.StartAll(); err != nil {
		return fmt.Errorf("collector start failed: %w", err)
	}

	// Start API server
	if err := c.apiServer.Start(c.config.Port, c.config.Host); err != nil {
		return fmt.Errorf("API server failed: %w", err)
	}

	// Start WebSocket server
	if err := c.wssServer.Start(c.config.Port+100, c.config.Host); err != nil {
		return fmt.Errorf("WebSocket server failed: %w", err)
	}

	return nil
}

func (c *Core) Stop() error {
	c.logger.Info("NieboGuard Node stopped")
	return nil
}
