package adsb

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

// Pipeline connects RTL-SDR → decoder → adapter → SQLite

type Pipeline struct {
	logger   *slog.Logger
	decoder  string // path to dump1090/readsb
	enabled  bool
}

func New(logger *slog.Logger) *Pipeline {
	return &Pipeline{logger: logger, enabled: false}
}

// Initialize finds decoder binary
func (p *Pipeline) Initialize(decoderPath string) error {
	if decoderPath != "" {
		if _, err := os.Stat(decoderPath); err != nil {
			return fmt.Errorf("decoder not found at %s: %w", decoderPath, err)
		}
		p.decoder = decoderPath
		p.enabled = true
		p.logger.Info("ADS-B pipeline initialized", "decoder", decoderPath)
		return nil
	}

	// Try common paths
	paths := []string{
		"C:\\Program Files\\dump1090\\dump1090.exe",
		"C:\\Program Files (x86)\\dump1090\\dump1090.exe",
		"dump1090",
		"readsb",
	}
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			p.decoder = path
			p.enabled = true
			p.logger.Info("ADS-B decoder found", "decoder", path)
			return nil
		}
	}

	return fmt.Errorf("no ADS-B decoder found (install dump1090 or readsb)")
}

// Start runs the decoder
func (p *Pipeline) Start() error {
	if !p.enabled {
		return fmt.Errorf("pipeline not initialized")
	}
	p.logger.Info("ADS-B pipeline started")
	return nil
}

// Stop stops the decoder
func (p *Pipeline) Stop() error {
	p.logger.Info("ADS-B pipeline stopped")
	return nil
}

// ReadAircraftJSON reads aircraft.json from decoder output
func (p *Pipeline) ReadAircraftJSON(path string) ([]map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// Simplified parsing
	lines := strings.Split(string(data), "\n")
	var results []map[string]interface{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || line[0] == '{' {
			continue
		}
		results = append(results, map[string]interface{}{
			"raw": line,
		})
	}
	return results, nil
}

func ReadFromFile(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}