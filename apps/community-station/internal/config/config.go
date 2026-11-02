package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Configuration represents the Community Station configuration.
type Configuration struct {
	Port     int    `yaml:"port"`
	Host     string `yaml:"host"`
	WorkDir  string `yaml:"workdir"`
	Verbose  bool   `yaml:"verbose"`
	Language string `yaml:"language"`
	Theme    string `yaml:"theme"`

	// RTL-SDR configuration
	RTLSDR RTLSDRConfig `yaml:"rtl_sdr"`

	// Simulator configuration
	Simulator SimulatorConfig `yaml:"simulator"`

	// Data sharing configuration
	DataSharing DataSharingConfig `yaml:"data_sharing"`
}

// RTLSDRConfig contains RTL-SDR receiver settings.
type RTLSDRConfig struct {
	DeviceIndex  int    `yaml:"device_index"`
	DeviceSerial string `yaml:"device_serial"`
	Frequency    int    `yaml:"frequency"` // Hz
	SampleRate   int    `yaml:"sample_rate"`
	GainMode     string `yaml:"gain_mode"` // auto/manual
	Gain         float64 `yaml:"gain"`
	PPM          int    `yaml:"ppm"`
	Enabled      bool   `yaml:"enabled"`
	DecoderPath  string `yaml:"decoder_path"` // path to external decoder if used
}

// SimulatorConfig contains simulator settings.
type SimulatorConfig struct {
	Enabled     bool   `yaml:"enabled"`
	Mode        string `yaml:"mode"` // synthetic, replay, live_transform
	Objects     int    `yaml:"objects"`
	StartLat    float64 `yaml:"start_lat"`
	StartLon    float64 `yaml:"start_lon"`
	SpeedMps    float64 `yaml:"speed_mps"`
	AltM        float64 `yaml:"alt_m"`
	RateHz      int    `yaml:"rate_hz"`
	Speed       float64 `yaml:"speed"`
}

// DataSharingConfig contains optional data sharing settings.
type DataSharingConfig struct {
	Enabled     bool   `yaml:"enabled"`
	URL         string `yaml:"url"`
	Token       string `yaml:"token"`
	StationID   string `yaml:"station_id"`
	IntervalSec int    `yaml:"interval_sec"`
	ShareSimulated bool `yaml:"share_simulated"`
}

// Load loads configuration from a YAML file.
func Load(path string) (*Configuration, error) {
	if path == "" {
		return &Configuration{
			Port:     8787,
			Host:     "127.0.0.1",
			WorkDir:  "",
			Verbose:  false,
			Language: "pl",
			Theme:    "dark",
			RTLSDR: RTLSDRConfig{
				DeviceIndex:  0,
				Frequency:    1090000000,
				SampleRate:   2400000,
				GainMode:     "auto",
				Gain:         0,
				PPM:          0,
				Enabled:      false,
			},
			Simulator: SimulatorConfig{
				Enabled:    false,
				Mode:       "synthetic",
				Objects:    10,
				StartLat:   52.2297,
				StartLon:   21.0122,
				SpeedMps:   250,
				AltM:       10000,
				RateHz:     1,
				Speed:      1.0,
			},
			DataSharing: DataSharingConfig{
				Enabled:         false,
				URL:             "",
				Token:           "",
				StationID:       "",
				IntervalSec:     30,
				ShareSimulated:  false,
			},
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

	// Set defaults
	if cfg.Port == 0 {
		cfg.Port = 8787
	}
	if cfg.Host == "" {
		cfg.Host = "127.0.0.1"
	}
	if cfg.Language == "" {
		cfg.Language = "pl"
	}
	if cfg.Theme == "" {
		cfg.Theme = "dark"
	}
	if cfg.WorkDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("getting home directory: %w", err)
		}
		cfg.WorkDir = filepath.Join(home, "NieboGuard", "Station")
	}
	if err := os.MkdirAll(cfg.WorkDir, 0755); err != nil {
		return nil, fmt.Errorf("creating work directory: %w", err)
	}

	return &cfg, nil
}

// Save writes configuration to a YAML file.
func (c *Configuration) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}
	return nil
}

// DefaultConfigPath returns the default configuration file path.
func DefaultConfigPath(workDir string) string {
	return filepath.Join(workDir, "config.yaml")
}

// DataDir returns the data directory.
func DataDir(workDir string) string {
	return filepath.Join(workDir, "data")
}

// LogsDir returns the logs directory.
func LogsDir(workDir string) string {
	return filepath.Join(workDir, "logs")
}

// DecodersDir returns the decoders directory.
func DecodersDir(workDir string) string {
	return filepath.Join(workDir, "decoders")
}

// EnsureDirs creates all required directories.
func EnsureDirs(workDir string) error {
	dirs := []string{
		workDir,
		DataDir(workDir),
		LogsDir(workDir),
		DecodersDir(workDir),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return err
		}
	}
	return nil
}
