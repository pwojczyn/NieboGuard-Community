package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/pwojczyn/NieboGuard-Community/simulator/internal/api"
	"github.com/pwojczyn/NieboGuard-Community/simulator/internal/sim"
)

func main() {
	port := flag.Int("port", 0, "HTTP server port")
	host := flag.String("host", "", "server listen address")
	verbose := flag.Bool("verbose", false, "verbose logging")
	flag.Parse()

	logLevel := slog.LevelInfo
	if *verbose {
		logLevel = slog.LevelDebug
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(logger)

	if *port == 0 {
		*port = 8789
	}
	if *host == "" {
		*host = "127.0.0.1"
	}

	simEngine := sim.NewEngine(sim.ModeSynthetic)
	simEngine.SyntheticCount = 10
	simEngine.StartLat = 52.2297
	simEngine.StartLon = 21.0122
	simEngine.DirDeg = 90
	simEngine.SpeedMps = 250
	simEngine.AltM = 10000
	simEngine.RateHz = 1
	simEngine.Seed = 42

	simAPI := api.New()

	mux := http.NewServeMux()
	mux.HandleFunc("/sim/api/", simAPI.ServeHTTP)
	mux.HandleFunc("/sim/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"simulated": true,
			"mode":      simEngine.Mode,
			"speed":     simEngine.Speed,
			"objects":   simEngine.SyntheticCount,
		})
	})

	addr := fmt.Sprintf("%s:%d", *host, *port)
	logger.Info("NieboGuard Simulator started", "addr", addr)

	srv := &http.Server{Addr: addr, Handler: mux}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("simulator server failed", "error", err)
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Info("shutting down simulator...")
	srv.Close()
	logger.Info("simulator stopped")
}

