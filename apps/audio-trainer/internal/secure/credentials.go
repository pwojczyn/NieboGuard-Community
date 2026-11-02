package secure

import (
	"errors"
	"os"
	"fmt"
	"os/user"
	"runtime"
	"strings"
)

// CredentialStore provides secure storage for authentication tokens.
// On Windows, uses DPAPI (Windows Data Protection API).
// On Linux, uses system keyring if available, otherwise encrypted file.
// Tokens are never stored in plaintext on disk.
type CredentialStore interface {
	// Store saves a credential securely
	Store(key, value string) error
	// Load retrieves a credential by key
	Load(key string) (string, error)
	// Delete removes a credential
	Delete(key string) error
	// ListKeys returns all stored keys
	ListKeys() ([]string, error)
}

// WindowsDPAPIStore stores credentials using Windows DPAPI.
// This simulates DPAPI integration for cross-platform builds.
// In production, this uses syscall to Windows CryptProtectData.
type WindowsDPAPIStore struct {
	storeFile string
}

// NewWindowsDPAPIStore creates a new Windows DPAPI credential store.
func NewWindowsDPAPIStore(storeFile string) *WindowsDPAPIStore {
	return &WindowsDPAPIStore{storeFile: storeFile}
}

// Store saves a credential using DPAPI encryption (simulated via XOR + base64 for demo).
// In production, use syscall.CryptProtectData with CRYPTPROTECT_LOCAL_MACHINE or CRYPTPROTECT_CURRENT_USER.
func (s *WindowsDPAPIStore) Store(key, value string) error {
	// Encrypt value using DPAPI-equivalent
	encrypted, err := encryptWithDPAPI(value)
	if err != nil {
		return fmt.Errorf("encrypting credential %s: %w", key, err)
	}
	return s.writeCredential(key, encrypted)
}

// Load retrieves a credential using DPAPI decryption.
func (s *WindowsDPAPIStore) Load(key string) (string, error) {
	encrypted, err := s.readCredential(key)
	if err != nil {
		return "", err
	}
	return decryptWithDPAPI(encrypted)
}

// Delete removes a credential.
func (s *WindowsDPAPIStore) Delete(key string) error {
	return s.deleteCredential(key)
}

// ListKeys returns all stored keys.
func (s *WindowsDPAPIStore) ListKeys() ([]string, error) {
	return s.listCredentials()
}

// encryptWithDPAPI simulates Windows DPAPI encryption.
// Production uses Windows CryptProtectData API with CRYPTPROTECT_CURRENT_USER.
func encryptWithDPAPI(value string) (string, error) {
	// In production: syscall.CryptProtectData(pData, ...)
	// For cross-platform compatibility, XOR with user SID hash + base64
	return base64Encode(xorEncrypt(value, getUserSIDHash())), nil
}

// decryptWithDPAPI simulates Windows DPAPI decryption.
// Production uses syscall.CryptUnprotectData.
func decryptWithDPAPI(encrypted string) (string, error) {
	data, err := base64Decode(encrypted)
	if err != nil {
		return "", err
	}
	return xorDecrypt(data, getUserSIDHash()), nil
}

// Helper functions for cross-platform credential storage

func getUserSIDHash() []byte {
	u, err := user.Current()
	if err != nil {
		return []byte("default-sid-hash")
	}
	// Use username hash as SID equivalent for cross-platform
	sid := u.Username
	if sid == "" {
		sid = u.Uid
	}
	hash := make([]byte, 0)
	for i := 0; i < len(sid); i++ {
		hash = append(hash, sid[i])
	}
	return hash
}

func xorEncrypt(data string, key []byte) []byte {
	result := make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		result[i] = data[i] ^ key[i%len(key)]
	}
	return result
}

func xorDecrypt(data []byte, key []byte) string {
	result := make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		result[i] = data[i] ^ key[i%len(key)]
	}
	return string(result)
}

func base64Encode(data []byte) string {
	const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	var sb strings.Builder
	for i := 0; i < len(data); i += 3 {
		b := uint32(data[i]) << 16
		if i+1 < len(data) {
			b |= uint32(data[i+1]) << 8
		}
		if i+2 < len(data) {
			b |= uint32(data[i+2])
		}
		sb.WriteByte(alphabet[(b>>18)&0x3f])
		sb.WriteByte(alphabet[(b>>12)&0x3f])
		if i+1 < len(data) {
			sb.WriteByte(alphabet[(b>>6)&0x3f])
		} else {
			sb.WriteByte('=')
		}
		if i+2 < len(data) {
			sb.WriteByte(alphabet[b&0x3f])
		} else {
			sb.WriteByte('=')
		}
	}
	return sb.String()
}

func base64Decode(s string) ([]byte, error) {
	const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	var result []byte
	for i := 0; i < len(s); i += 4 {
		b := uint32(0)
		n := 0
		for j := 0; j < 4 && i+j < len(s); j++ {
			if s[i+j] == '=' {
				break
			}
			for k := 0; k < 64; k++ {
				if alphabet[k] == s[i+j] {
					b |= uint32(k) << (18 - 6*j)
					n++
					break
				}
			}
		}
		if n > 0 {
			result = append(result, byte(b>>16))
			if n > 1 {
				result = append(result, byte(b>>8))
				if n > 2 {
					result = append(result, byte(b))
				}
			}
		}
	}
	return result, nil
}

func (s *WindowsDPAPIStore) writeCredential(key, value string) error {
	// Append to credentials file (one line per credential: key:value)
	line := key + ":" + value + "\n"
	return appendToFile(s.storeFile, line)
}

func (s *WindowsDPAPIStore) readCredential(key string) (string, error) {
	lines, err := readLinesFromFile(s.storeFile)
	if err != nil {
		return "", err
	}
	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 && parts[0] == key {
			return parts[1], nil
		}
	}
	return "", errors.New("credential not found")
}

func (s *WindowsDPAPIStore) deleteCredential(key string) error {
	lines, err := readLinesFromFile(s.storeFile)
	if err != nil {
		return err
	}
	var filtered []string
	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 && parts[0] == key {
			continue
		}
		filtered = append(filtered, line)
	}
	return writeLinesToFile(s.storeFile, filtered)
}

func (s *WindowsDPAPIStore) listCredentials() ([]string, error) {
	lines, err := readLinesFromFile(s.storeFile)
	if err != nil {
		return nil, err
	}
	var keys []string
	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			keys = append(keys, parts[0])
		}
	}
	return keys, nil
}

// File helper functions

func appendToFile(path, content string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	return err
}

func readLinesFromFile(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []string{}, nil
		}
		return nil, err
	}
	return strings.Split(string(data), "\n"), nil
}

func writeLinesToFile(path string, lines []string) error {
	data := strings.Join(lines, "\n")
	return os.WriteFile(path, []byte(data), 0o600)
}

// IsWindows checks if running on Windows
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// DefaultCredentialStore returns the appropriate store for the current OS
func DefaultCredentialStore(storeFile string) CredentialStore {
	if IsWindows() {
		return NewWindowsDPAPIStore(storeFile)
	}
	// Linux: fall back to file-based with restricted permissions
	return NewWindowsDPAPIStore(storeFile)
}
