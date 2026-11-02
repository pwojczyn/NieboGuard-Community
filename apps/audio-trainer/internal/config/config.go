package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Configuration describes Audio AI Trainer settings.
type Configuration struct {
	Port     int    `yaml:"port"`
	Host     string `yaml:"host"`
	WorkDir  string `yaml:"workdir"`
	Language string `yaml:"language"`
	Theme    string `yaml:"theme"`

	Training TrainingConfig `yaml:"training"`
	Dataset  DatasetConfig  `yaml:"dataset"`
	Compute  ComputeConfig  `yaml:"compute"`
	Consent  ConsentConfig  `yaml:"consent"`
}

type TrainingConfig struct {
	Enabled             bool    `yaml:"enabled"`
	Mode                string  `yaml:"mode"`
	Backbone            string  `yaml:"backbone"`
	Epochs              int     `yaml:"epochs"`
	BatchSize           int     `yaml:"batch_size"`
	LearningRate        float64 `yaml:"learning_rate"`
	MaxThreads          int     `yaml:"max_threads"`
	CpuLimitPercent     int     `yaml:"cpu_limit_percent"`
	CheckpointDir       string  `yaml:"checkpoint_dir"`
	ResumeFromCheckpoint string `yaml:"resume_from_checkpoint"`
	YAMNetModelPath     string  `yaml:"yamnet_model_path"`
	TensorFlowEnabled   bool    `yaml:"tensorflow_enabled"`
	TensorFlowPath      string  `yaml:"tensorflow_path"`
}

type DatasetConfig struct {
	ImportDir     string `yaml:"import_dir"`
	ManifestPath  string `yaml:"manifest_path"`
	SplitTrain    int    `yaml:"split_train"`
	SplitVal      int    `yaml:"split_val"`
	SplitTest     int    `yaml:"split_test"`
	ResampleHz    int    `yaml:"resample_hz"`
	Mono          bool   `yaml:"mono"`
	Normalize     bool   `yaml:"normalize"`
}

type ComputeConfig struct {
	Enabled   bool   `yaml:"enabled"`
	ServerURL string `yaml:"server_url"`
	// Token jest ładowany z bezpiecznego magazynu (Windows DPAPI / Linux keyring)
	StoredAtRuntime bool `yaml:"stored_at_runtime"`
	StationID string `yaml:"station_id"`
}

type ConsentConfig struct {
	Acknowledged bool   `yaml:"acknowledged"`
	Timestamp    string `yaml:"timestamp"`
	UploadResults bool  `yaml:"upload_results"`
	ShareData    bool   `yaml:"share_data"`
}

func DefaultConfigPath() string {
	return filepath.Join(os.Getenv("USERPROFILE"), "NieboGuard", "AudioTrainer", "config.yaml")
}

func Default() *Configuration {
	return &Configuration{
		Port:      8788,
		Host:      "127.0.0.1",
		WorkDir:   "%LOCALAPPDATA%\\NieboGuard\\AudioTrainer",
		Language:  "pl",
		Theme:     "dark",
		Training: TrainingConfig{
			Enabled:            false,
			Mode:               "cpu-lite",
			Backbone:           "yamnet",
			Epochs:             50,
			BatchSize:          32,
			LearningRate:       0.001,
			MaxThreads:         4,
			CpuLimitPercent:    70,
			CheckpointDir:      "checkpoints",
			YAMNetModelPath:    "",
			TensorFlowEnabled:  false,
			TensorFlowPath:     "",
		},
		Dataset: DatasetConfig{
			ImportDir:     "imports",
			ManifestPath:  "manifest.yaml",
			SplitTrain:    70,
			SplitVal:      20,
			SplitTest:     10,
			ResampleHz:    16000,
			Mono:          true,
			Normalize:     true,
		},
		Compute: ComputeConfig{
			Enabled:   false,
			ServerURL: "",
			StoredAtRuntime: false,
			StationID: "",
		},
		Consent: ConsentConfig{
			Acknowledged: false,
			UploadResults: false,
			ShareData:    false,
		},
	}
}

func Load(path string) (*Configuration, error) {
	cfg := Default()
	if path == "" {
		return cfg, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, Save(cfg, path)
		}
		return nil, fmt.Errorf("reading config %s: %w", path, err)
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config %s: %w", path, err)
	}
	return cfg, nil
}

func Save(cfg *Configuration, path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("creating config dir %s: %w", dir, err)
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("writing config %s: %w", path, err)
	}
	return nil
}

// GetToken retrieves the compute token from secure storage (Windows DPAPI / Linux keyring).
// Returns empty string if not stored yet.
func GetToken(storeFile string) (string, error) {
	// Will be implemented using secure.CredentialStore
	return "", nil
}

// SetToken stores the compute token in secure storage (Windows DPAPI / Linux keyring).
// Never saves token in plaintext config file.
func SetToken(storeFile, token string) error {
	// Will be implemented using secure.CredentialStore
	return nil
}
