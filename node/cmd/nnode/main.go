package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/pwojczyn/NieboGuard-Community/node/internal/core"
)

func main() {
	configPath := flag.String("config", "", "path to configuration file")
	port := flag.Int("port", 0, "HTTP server port")
	host := flag.String("host", "", "server listen address")
	workdir := flag.String("workdir", "", "working directory for data storage")
	verbose := flag.Bool("verbose", false, "verbose logging")
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

	// Initialize core engine
	opts := core.Options{
		ConfigPath: *configPath,
		Port:       *port,
		Host:       *host,
		WorkDir:    *workdir,
		Logger:     logger,
	}

	app, err := core.New(opts)
	if err != nil {
		slog.Error("failed to initialize core", "error", err)
		os.Exit(1)
	}

	// Start application
	if err := app.Start(); err != nil {
		slog.Error("application failed", "error", err)
		os.Exit(1)
	}

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Graceful shutdown
	slog.Info("shutting down...")
	if err := app.Stop(); err != nil {
		slog.Error("shutdown error", "error", err)
		os.Exit(1)
	}
	slog.Info("stopped cleanly")
	fmt.Println("NieboGuard Node stopped")
}
