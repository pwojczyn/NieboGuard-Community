package r2

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	crypto "github.com/pwojczyn/NieboGuard-Community/apps/audio-trainer/internal/crypto"
)

// DownloadManager struct remains unchanged
type DownloadManager struct {
	logger    *slog.Logger
	downloads map[string]*DownloadStatus
	mu        sync.RWMutex
}

// DownloadStatus struct remains unchanged
type DownloadStatus struct {
	ObjectPath   string    `json:"object_path"`
	URL          string    `json:"url"`
	SizeBytes    int64     `json:"size_bytes"`
	Hash         string    `json:"sha256"`
	Downloaded   int64     `json:"downloaded"`
	Status       string    `json:"status"` // pending, downloading, completed, failed, verified
	StartedAt    time.Time `json:"started_at"`
	CompletedAt  time.Time `json:"completed_at"`
	Error        string    `json:"error,omitempty"`
}

func NewDownloadManager() *DownloadManager {
	return &DownloadManager{
		logger:    slog.Default().With("component", "r2-download"),
		downloads: make(map[string]*DownloadStatus),
	}
}

// DownloadObject downloads an object with manifest verification
func (dm *DownloadManager) DownloadObject(objectPath, url, expectedHash string, sizeBytes int64, destPath string) error {
	// 1. Verify manifest signature and file hashes before downloading
	manifestPath := filepath.Join(objectPath, "manifest.json")

	// Check if manifest exists and verify it
	manifestData, err := os.ReadFile(manifestPath)
	if err != nil {
		dm.updateStatus(objectPath, "failed", fmt.Sprintf("manifest not found: %v", err))
		return fmt.Errorf("manifest not found: %w", err)
	}

	var manifest crypto.Manifest
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		dm.updateStatus(objectPath, "failed", fmt.Sprintf("manifest JSON parse error: %v", err))
		return fmt.Errorf("manifest JSON parse error: %w", err)
	}

	// Verify manifest signature
	if err := crypto.VerifyManifest(&manifest); err != nil {
		dm.updateStatus(objectPath, "failed", fmt.Sprintf("manifest verification failed: %v", err))
		return fmt.Errorf("manifest verification failed: %w", err)
	}

	// Verify all files have matching hashes
	for _, file := range manifest.Files {
		filePath := filepath.Join(objectPath, file.Path)
		fileData, err := os.ReadFile(filePath)
		if err != nil {
			dm.updateStatus(objectPath, "failed", fmt.Sprintf("file not found: %s", filePath))
			return fmt.Errorf("file not found: %w", err)
		}

		computedHash := fmt.Sprintf("%x", sha256.Sum256(fileData))
		if computedHash != file.SHA256 {
			dm.updateStatus(objectPath, "failed", fmt.Sprintf("file %s hash mismatch", file.Path))
			return fmt.Errorf("file %s hash mismatch", file.Path)
		}
	}

	// 2. Proceed with download (rest of the function unchanged)
	dm.mu.Lock()
	dm.downloads[objectPath] = &DownloadStatus{
		ObjectPath: objectPath,
		URL:        url,
		SizeBytes:  sizeBytes,
		Hash:       expectedHash,
		Status:     "downloading",
		StartedAt:  time.Now(),
	}
	dm.mu.Unlock()

	// Ensure destination directory exists
	dir := filepath.Dir(destPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		dm.updateStatus(objectPath, "failed", fmt.Sprintf("creating directory: %v", err))
		return err
	}

	// Download to temporary file
	tempPath := destPath + ".tmp"
	outFile, err := os.Create(tempPath)
	if err != nil {
		dm.updateStatus(objectPath, "failed", fmt.Sprintf("creating temp file: %v", err))
		return err
	}
	defer outFile.Close()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		dm.updateStatus(objectPath, "failed", fmt.Sprintf("creating request: %v", err))
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		dm.updateStatus(objectPath, "failed", fmt.Sprintf("sending request: %v", err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		dm.updateStatus(objectPath, "failed", fmt.Sprintf("HTTP %d", resp.StatusCode))
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	// Download with SHA-256 verification
	hasher := sha256.New()
	multiWriter := io.MultiWriter(outFile, hasher)

	reader := &progressReader{reader: resp.Body, onProgress: func(n int64) {
		dm.mu.Lock()
		if ds, ok := dm.downloads[objectPath]; ok {
			ds.Downloaded = n
		}
		dm.mu.Unlock()
	}}

	if _, err := io.Copy(multiWriter, reader); err != nil {
		dm.updateStatus(objectPath, "failed", fmt.Sprintf("download error: %v", err))
		outFile.Close()
		os.Remove(tempPath)
		return err
	}

	outFile.Close()

	// Verify SHA-256
	computedHash := hex.EncodeToString(hasher.Sum(nil))
	if strings.ToLower(computedHash) != strings.ToLower(expectedHash) {
		dm.updateStatus(objectPath, "failed", fmt.Sprintf("SHA-256 mismatch: computed %s, expected %s", computedHash, expectedHash))
		os.Remove(tempPath)
		return fmt.Errorf("SHA-256 mismatch")
	}

	// Atomic rename
	if err := os.Rename(tempPath, destPath); err != nil {
		dm.updateStatus(objectPath, "failed", fmt.Sprintf("atomic rename: %v", err))
		return err
	}

	// Update status
	dm.mu.Lock()
	if ds, ok := dm.downloads[objectPath]; ok {
		ds.Status = "verified"
		ds.CompletedAt = time.Now()
		ds.Downloaded = sizeBytes
	}
	dm.mu.Unlock()

	return nil
}

func (dm *DownloadManager) updateStatus(objectPath, status, errMsg string) {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	if ds, ok := dm.downloads[objectPath]; ok {
		ds.Status = status
		ds.Error = errMsg
		if status == "failed" {
			ds.CompletedAt = time.Now()
		}
	}
}

// GetStatus returns download status for all objects
func (dm *DownloadManager) GetStatus() []DownloadStatus {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	statuses := make([]DownloadStatus, 0, len(dm.downloads))
	for _, ds := range dm.downloads {
		statuses = append(statuses, *ds)
	}
	return statuses
}

// progressReader wraps an io.Reader to track progress
type progressReader struct {
	reader   io.Reader
	onProgress func(n int64)
	total    int64
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.reader.Read(p)
	pr.total += int64(n)
	if pr.onProgress != nil {
		pr.onProgress(pr.total)
	}
	return n, err
}

func (dm *DownloadManager) GenerateSignedURL(objectPath string, expires time.Duration) (string, error) {
	hash := sha256.New()
	hash.Write([]byte(objectPath + expires.String()))
	signature := hex.EncodeToString(hash.Sum(nil)[:32])
	return fmt.Sprintf("https://signed-r2/%s?sig=%s&exp=%d", objectPath, signature, time.Now().Add(expires).Unix()), nil
}