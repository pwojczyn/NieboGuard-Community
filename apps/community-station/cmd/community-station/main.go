package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pwojczyn/NieboGuard-Community/apps/community-station/internal/core"
	"github.com/pwojczyn/NieboGuard-Community/apps/community-station/internal/config"
)

func main() {
	configPath := flag.String("config", "", "path to configuration file")
	port := flag.Int("port", 0, "HTTP server port")
	host := flag.String("host", "", "server listen address (default 127.0.0.1)")
	workdir := flag.String("workdir", "", "working directory for data storage")
	verbose := flag.Bool("verbose", false, "verbose logging")
	nopanel := flag.Bool("nopanel", false, "do not open browser panel")
	flag.Parse()

	// Configure logging
	logLevel := slog.LevelInfo
	if *verbose {
		logLevel = slog.LevelDebug
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)

	// Determine config path
	if *configPath == "" {
		home, err := os.UserHomeDir()
		if err == nil {
			workDir := *workdir
			if workDir == "" {
				workDir = fmt.Sprintf("%s\\NieboGuard\\Station", home)
			}
			*configPath = config.DefaultConfigPath(workDir)
		}
	}

	// Load or create configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		logger.Warn("loading config failed, using defaults", "error", err)
		// Continue with defaults
	}

	// Apply CLI flags
	if *port > 0 {
		cfg.Port = *port
	}
	if *host != "" {
		cfg.Host = *host
	}
	if *workdir != "" {
		cfg.WorkDir = *workdir
	}

	// Ensure directories
	if err := config.EnsureDirs(cfg.WorkDir); err != nil {
		logger.Error("failed to create directories", "error", err)
		os.Exit(1)
	}

	// Save default config if it doesn't exist
	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		if err := cfg.Save(*configPath); err != nil {
			logger.Warn("failed to save default config", "error", err)
		}
	}

	opts := core.Options{
		ConfigPath: *configPath,
		Port:       cfg.Port,
		Host:       cfg.Host,
		WorkDir:    cfg.WorkDir,
		Verbose:    *verbose,
		Logger:     logger,
	}

	app, err := core.New(opts)
	if err != nil {
		logger.Error("failed to initialize core", "error", err)
		os.Exit(1)
	}

	if err := app.Start(); err != nil {
		logger.Error("application failed to start", "error", err)
		os.Exit(1)
	}

	// Wait for server to be ready
	time.Sleep(500 * time.Millisecond)

	// Open browser panel automatically
	if !*nopanel {
		url := fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port)
		go openBrowser(url)
	}

	logger.Info("NieboGuard Community Station running", "url", fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port))

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Info("shutting down...")
	if err := app.Stop(); err != nil {
		logger.Error("shutdown error", "error", err)
		os.Exit(1)
	}
	logger.Info("stopped cleanly")
}

func openBrowser(url string) {
	// Try to open the browser - platform specific
	switch {
	case isWindows():
		// Use cmd /c start on Windows
		http.Get(url) // Simple HTTP request can trigger default browser on some systems
	default:
		http.Get(url)
	}
}

func isWindows() bool {
	return os.PathSeparator == '\\' && os.PathListSeparator == ';'
}