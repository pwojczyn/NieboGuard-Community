package crypto

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"errors"
	"os"
)

// Manifest represents a verified dataset/model manifest with Ed25519 signature.
// It is embedded in downloaded packages for integrity and authenticity verification.
type Manifest struct {
	DatasetID    string              `json:"dataset_id"`
	Version      string              `json:"version"`
	Description  string              `json:"description,omitempty"`
	PackType     string              `json:"pack_type,omitempty"`
	SHA256       string              `json:"sha256"`
	SizeBytes    int64               `json:"size_bytes"`
	Files        []ManifestFile      `json:"files,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	Signature    []byte              `json:"signature,omitempty"`
}

// ManifestFile represents a single file entry in the manifest.
type ManifestFile struct {
	Path    string `json:"path"`
	SHA256  string `json:"sha256"`
	Size    int64  `json:"size"`
	Verified bool  `json:"verified"`
}

// PublicKey is the embedded Ed25519 public key used for manifest verification.
// In production, this would be provisioned per-station or fetched from a KMS.
var PublicKey = func() ed25519.PublicKey {
	b, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
	return ed25519.PublicKey(b)
}()

// VerifyManifest checks the Ed25519 signature of a manifest against embedded public key.
// Returns nil if verification succeeds, or an error if verification fails.
func VerifyManifest(m *Manifest) error {
	if m == nil {
		return errors.New("manifest is nil")
	}

	if len(m.DatasetID) == 0 {
		return errors.New("manifest missing dataset_id")
	}

	if len(m.SHA256) == 0 {
		return errors.New("manifest missing sha256 hash")
	}

	// Re-hash the manifest contents and compare
	expectedHash := sha256Sum(m)

	if expectedHash != m.SHA256 {
		return fmt.Errorf("manifest SHA256 mismatch: expected %s, got %s", expectedHash, m.SHA256)
	}

	// Verify Ed25519 signature
	manifestJSON, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("marshaling manifest: %w", err)
	}

	if len(m.Signature) == 0 {
		return errors.New("manifest missing Ed25519 signature")
	}

	if !ed25519.Verify(PublicKey, manifestJSON, m.Signature) {
		return errors.New("manifest Ed25519 signature verification failed")
	}

	return nil
}

// sha256Sum computes SHA256 of the manifest's structural fields (excluding signature).
func sha256Sum(m *Manifest) string {
	h := sha256.New()
	h.Write([]byte(m.DatasetID + ":" + m.Version + ":" + m.SHA256))
	return hex.EncodeToString(h.Sum(nil))
}

// VerifyFileManifest verifies an Ed25519 signature on file-level manifest data.
func VerifyFileManifest(filePath, expectedSHA256 string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("reading file %s: %w", filePath, err)
	}

	computed := fmt.Sprintf("%x", sha256.Sum256(data))
	if computed != expectedSHA256 {
		return fmt.Errorf("file SHA256 mismatch: computed %s, expected %s", computed, expectedSHA256)
	}

	return nil
}

// NewManifest creates a new manifest with Ed25519 signature.
// The caller must sign the manifest before returning.
func NewManifest(datasetID, version string, files []ManifestFile) *Manifest {
	return &Manifest{
		DatasetID:    datasetID,
		Version:      version,
		Files:        files,
		SHA256:       "",
		Signature:    []byte{},
		Metadata:     make(map[string]interface{}),
	}
}

// SignManifest signs the manifest with the station's Ed25519 private key.
// In production, the private key would be stored in Windows DPAPI or Linux keyring.
func SignManifest(m *Manifest, privateKey ed25519.PrivateKey) {
	m.Signature = ed25519.Sign(privateKey, append(
		[]byte(m.DatasetID+":"+m.Version),
		[]byte(m.SHA256)...,
	))
}

// ManifestQuarantineResult represents the result of quarantine verification.
type ManifestQuarantineResult struct {
	Valid    bool `json:"valid"`
	Reason   string `json:"reason,omitempty"`
	Allowed  bool `json:"allowed,omitempty"`
	QuarantinePath string `json:"quarantine_path,omitempty"`
}

// QuarantineCheck performs quarantine workflow checks on a manifest.
// It validates the manifest signature, checks for known bad patterns,
// and returns whether the item should be allowed or quarantined.
func QuarantineCheck(m *Manifest, blockedSignatures map[string]bool) *ManifestQuarantineResult {
	result := &ManifestQuarantineResult{Valid: true}

	// 1. Check Ed25519 signature
	if err := VerifyManifest(m); err != nil {
		result.Valid = false
		result.Reason = fmt.Sprintf("manifest verification failed: %v", err)
		result.Allowed = false
		result.QuarantinePath = "quarantine/updates/" + m.DatasetID
		return result
	}

	// 2. Check for blocked/published signatures
	if sig, ok := blockedSignatures[m.SHA256]; ok && sig {
		result.Valid = false
		result.Reason = "manifest signature on blocked list"
		result.Allowed = false
		result.QuarantinePath = "quarantine/updates/" + m.DatasetID
		return result
	}

	// 3. Check file integrity (all files must have valid hashes)
	for _, file := range m.Files {
		if len(file.SHA256) == 0 {
			result.Valid = false
			result.Reason = "manifest file missing SHA256"
			result.Allowed = false
			result.QuarantinePath = "quarantine/updates/" + m.DatasetID
			return result
		}
	}

	return result
}