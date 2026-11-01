package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Configuration represents the application configuration.
type Configuration struct {
	Port     int    `yaml:"port"`
	Host     string `yaml:"host"`
	WorkDir  string `yaml:"workdir"`
	Verbose  bool   `yaml:"verbose"`
}

// Load loads configuration from a YAML file.
func Load(path string) (*Configuration, error) {
	if path == "" {
		// Use default configuration
		return &Configuration{
			Port:    8787,
			Host:    "127.0.0.1",
			WorkDir: "",
			Verbose: false,
		}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Configuration
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing YAML: %w", err)
	}

	// Set defaults if needed
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

	return &cfg, nil
}
