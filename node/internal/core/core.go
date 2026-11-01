package core

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/pwojczyn/NieboGuard-Community/node/pkg/config"
)

// Options configures the Core Engine.
type Options struct {
	ConfigPath string
	Port       int
	Host       string
	WorkDir    string
	Logger     *slog.Logger
}

// Core represents the NieboGuard Node engine.
type Core struct {
	opts        Options
	config      *config.Configuration
	logger      *slog.Logger
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

	core := &Core{
		opts:   opts,
		config: cfg,
		logger: logger,
	}

	return core, nil
}

func (c *Core) Start() error {
	c.logger.Info("NieboGuard Node started", "port", c.config.Port, "host", c.config.Host)
	return nil
}

func (c *Core) Stop() error {
	c.logger.Info("NieboGuard Node stopped")
	return nil
}