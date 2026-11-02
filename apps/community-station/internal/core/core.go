package core

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/pwojczyn/NieboGuard-Community/apps/community-station/internal/api"
	"github.com/pwojczyn/NieboGuard-Community/apps/community-station/internal/collectors"
	"github.com/pwojczyn/NieboGuard-Community/apps/community-station/internal/config"
	"github.com/pwojczyn/NieboGuard-Community/apps/community-station/internal/rtl"
	"github.com/pwojczyn/NieboGuard-Community/apps/community-station/internal/simulator"
	"github.com/pwojczyn/NieboGuard-Community/apps/community-station/internal/storage"
	"github.com/pwojczyn/NieboGuard-Community/apps/community-station/internal/websocket"
)

// Options configures the Community Station Core Engine.
type Options struct {
	ConfigPath string
	Port       int
	Host       string
	WorkDir    string
	Verbose    bool
	Logger     *slog.Logger
}

// Core represents the NieboGuard Community Station engine.
type Core struct {
	opts        Options
	config      *config.Configuration
	logger      *slog.Logger
	storage     *storage.Storage
	collector   *collectors.Manager
	apiServer   *api.Server
	wssServer   *websocket.Server
	rtlReceiver *rtl.Receiver
	simulator   *simulator.Engine
	shutdown    chan struct{}
	stopped     chan struct{}
}

// New creates a new Community Station Core engine instance.
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

	// Initialize storage
	storage, err := storage.New(cfg.WorkDir)
	if err != nil {
		return nil, fmt.Errorf("creating storage: %w", err)
	}

	// Initialize collector manager
	collector, err := collectors.New(cfg, storage, logger)
	if err != nil {
		return nil, fmt.Errorf("creating collector manager: %w", err)
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

	// Initialize RTL-SDR receiver if enabled
	var rtlReceiver *rtl.Receiver
	if cfg.RTLSDR.Enabled {
		rtlReceiver, err = rtl.NewReceiver(cfg.RTLSDR, logger)
		if err != nil {
			logger.Warn("RTL-SDR initialization failed, continuing without hardware", "error", err)
			// RTL-SDR not critical, continue without it
		}
	}

	// Initialize simulator if enabled
	var sim *simulator.Engine
	if cfg.Simulator.Enabled {
		simCfg := simulator.Config{
			Enabled:   cfg.Simulator.Enabled,
			Mode:      cfg.Simulator.Mode,
			Objects:   cfg.Simulator.Objects,
			StartLat:  cfg.Simulator.StartLat,
			StartLon:  cfg.Simulator.StartLon,
			SpeedMps:  cfg.Simulator.SpeedMps,
			AltM:      cfg.Simulator.AltM,
			RateHz:    cfg.Simulator.RateHz,
			Speed:     cfg.Simulator.Speed,
		}
		sim, err = simulator.NewEngine(cfg.Simulator.Mode, simCfg)
		if err != nil {
			return nil, fmt.Errorf("creating simulator: %w", err)
		}
	}

	core := &Core{
		opts:        opts,
		config:      cfg,
		logger:      logger,
		storage:     storage,
		collector:   collector,
		apiServer:   apiServer,
		wssServer:   wssServer,
		rtlReceiver: rtlReceiver,
		simulator:   sim,
		shutdown:    make(chan struct{}),
		stopped:     make(chan struct{}),
	}

	return core, nil
}

// Start starts the Community Station application.
func (c *Core) Start() error {
	c.logger.Info("NieboGuard Community Station started", "port", c.config.Port, "host", c.config.Host)

	// Initialize storage
	if err := c.storage.Initialize(); err != nil {
		return fmt.Errorf("storage initialization failed: %w", err)
	}

	// Start collectors
	if err := c.collector.StartAll(); err != nil {
		return fmt.Errorf("collector start failed: %w", err)
	}

	// Start API server
	if err := c.apiServer.Start(); err != nil {
		return fmt.Errorf("API server failed: %w", err)
	}

	// Start WebSocket server
	if err := c.wssServer.Start(c.config.Port+100, c.config.Host); err != nil {
		return fmt.Errorf("WebSocket server failed: %w", err)
	}

	// Start RTL-SDR receiver if enabled
	if c.rtlReceiver != nil {
		if err := c.rtlReceiver.Start(); err != nil {
			c.logger.Warn("RTL-SDR receiver start failed, continuing without hardware", "error", err)
		} else {
			c.logger.Info("RTL-SDR receiver started")
		}
	}

	// Start simulator if enabled
	if c.simulator != nil {
		c.logger.Info("Simulator started", "mode", c.simulator.Mode())
	}

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Handle signals in goroutine
	go func() {
		sig := <-sigChan
		c.logger.Info("received signal, shutting down...", "signal", sig)
		c.Stop()
	}()

	c.logger.Info("NieboGuard Community Station running", "port", c.config.Port)
	return nil
}

// Stop gracefully shuts down the Community Station application.
func (c *Core) Stop() error {
	c.logger.Info("NieboGuard Community Station stopping...")

	// Stop simulator
	if c.simulator != nil {
		c.simulator.Stop()
		c.logger.Info("Simulator stopped")
	}

	// Stop RTL-SDR receiver
	if c.rtlReceiver != nil {
		c.rtlReceiver.Stop()
		c.logger.Info("RTL-SDR receiver stopped")
	}

	// Stop WebSocket server
	c.wssServer.Stop()

	// Stop API server
	c.apiServer.Stop()

	// Stop collectors
	c.collector.StopAll()

	// Close storage
	if err := c.storage.Close(); err != nil {
		c.logger.Error("storage close error", "error", err)
	}

	close(c.stopped)
	c.logger.Info("NieboGuard Community Station stopped cleanly")
	return nil
}

// Health returns the health status of the application.
func (c *Core) Health() map[string]interface{} {
	status := map[string]interface{}{
		"status":   "ok",
		"service":  "nieboguard-community-station",
		"version":  "0.1.0",
		"host":     c.config.Host,
		"port":     c.config.Port,
		"sources":  len(c.collector.GetSources()),
	}

	// Add RTL-SDR status
	if c.rtlReceiver != nil {
		status["rtl_sdr"] = c.rtlReceiver.Status()
	} else {
		status["rtl_sdr"] = map[string]string{"status": "not_configured"}
	}

	// Add simulator status
	if c.simulator != nil {
		status["simulator"] = c.simulator.Status()
	} else {
		status["simulator"] = map[string]string{"status": "disabled"}
	}

	return status
}

// Sources returns information about data sources.
func (c *Core) Sources() []collectors.SourceInfo {
	return c.collector.GetSources()
}

// Objects returns aircraft events from storage.
func (c *Core) Objects(limit, offset int) ([]map[string]interface{}, error) {
	return c.storage.GetAircraft(limit, offset)
}

// SimulatorEngine returns the simulator engine if enabled.
func (c *Core) SimulatorEngine() *simulator.Engine {
	return c.simulator
}

// RTLReceiver returns the RTL-SDR receiver if enabled.
func (c *Core) RTLReceiver() *rtl.Receiver {
	return c.rtlReceiver
}