package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"github.com/pwojczyn/NieboGuard-Community/apps/audio-trainer/internal/config"
	"github.com/pwojczyn/NieboGuard-Community/apps/audio-trainer/internal/dashboard"
	"github.com/pwojczyn/NieboGuard-Community/apps/audio-trainer/internal/training"
	"github.com/pwojczyn/NieboGuard-Community/apps/audio-trainer/internal/dataset"
	"github.com/pwojczyn/NieboGuard-Community/apps/audio-trainer/internal/model"
)

func main() {
	// Parse command line arguments
	configPath := ""
	port := 0
	host := ""
	
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--config":
			if i+1 < len(os.Args) {
				configPath = os.Args[i+1]
				i++
			}
		case "--port":
			if i+1 < len(os.Args) {
				if p, err := strconv.Atoi(os.Args[i+1]); err == nil {
					port = p
				}
				i++
			}
		case "--host":
			if i+1 < len(os.Args) {
				host = os.Args[i+1]
				i++
			}
		}
	}

	// Load configuration
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Override config with command line arguments
	if port > 0 {
		cfg.Port = port
	}
	if host != "" {
		cfg.Host = host
	}

	// Ensure workdir exists
	cfg.WorkDir = filepath.Clean(os.ExpandEnv(cfg.WorkDir))
	if err := os.MkdirAll(cfg.WorkDir, 0o755); err != nil {
		log.Fatalf("Failed to create workdir %s: %v", cfg.WorkDir, err)
	}

	// Create components
	trainer, err := training.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create trainer: %v", err)
	}

	datasetManager, err := dataset.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create dataset manager: %v", err)
	}

	modelManager, err := model.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create model manager: %v", err)
	}

	api, err := dashboard.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create dashboard API: %v", err)
	}

	// Start background processes
	ctx, cancel := context.WithCancel(context.Background())
	go trainer.Start(ctx)
	go datasetManager.Start(ctx)
	go modelManager.Start(ctx)

	// Start API server
	go func() {
		r := mux.NewRouter()
		api.RegisterRoutes(r)
		addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
		log.Printf("Starting Audio AI Trainer on %s", addr)
		if err := http.ListenAndServe(addr, r); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutting down...")
	cancel()

	// Graceful shutdown with timeout
	timer := time.NewTimer(30 * time.Second)
	select {
	case <-timer.C:
		log.Println("Shutdown timeout reached")
	case <-ctx.Done():
		timer.Stop()
	}

	log.Println("Audio AI Trainer stopped")
}