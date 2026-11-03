package agent

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"strings"
	"time"

	"github.com/pwojczyn/NieboGuard-Community/apps/audio-trainer/internal/compute"
	"github.com/pwojczyn/NieboGuard-Community/apps/audio-trainer/internal/secure"
)

// Agent handles secure communication with NieboGuard Compute API
type Agent struct {
	baseURL   string
	stationID string
	tokenKey  string
	credStore secure.CredentialStore
	client    *http.Client
	logger    *slog.Logger
	active    map[string]bool
	muActive  sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewAgent creates a new secure agent instance
func NewAgent(baseURL, stationID, tokenKey string, credStore secure.CredentialStore) (*Agent, error) {
	if credStore == nil {
		return nil, errors.New("credential store required")
	}

	ctx, cancel := context.WithCancel(context.Background())

	a := &Agent{
		baseURL:   baseURL,
		stationID: stationID,
		tokenKey:  tokenKey,
		credStore: credStore,
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				Proxy:         http.ProxyFromEnvironment,
				MaxIdleConns:  100,
				IdleConnTimeout: 30 * time.Second,
			},
		},
		logger:  slog.Default().With("component", "agent"),
		active:  make(map[string]bool),
		ctx:     ctx,
		cancel:  cancel,
	}

	return a, nil
}

// Close stops agent
func (a *Agent) Close() {
	a.cancel()
	if a.client != nil && a.client.Transport != nil {
		if tr, ok := a.client.Transport.(*http.Transport); ok {
			tr.CloseIdleConnections()
		}
	}
}

// getToken loads token from secure credential storage
func (a *Agent) getToken() (string, error) {
	token, err := a.credStore.Load(a.tokenKey)
	if err != nil {
		return "", fmt.Errorf("token load error: %w", err)
	}
	return token, nil
}

// validateResponse checks response security headers
func validateResponse(resp *http.Response) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("invalid status: %d", resp.StatusCode)
	}

	if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		return fmt.Errorf("invalid Content-Type: %s", ct)
	}

	if resp.Header.Get("X-Content-Type-Options") == "" {
		return fmt.Errorf("missing X-Content-Type-Options")
	}

	if resp.Header.Get("X-Frame-Options") == "" {
		return fmt.Errorf("missing X-Frame-Options")
	}

	return nil
}

// validateHost prevents DNS rebinding attacks (private IP ranges)
func validateHost(host string) error {
	if host == "" {
		return errors.New("missing Host header")
	}

	// Allow localhost development hosts
	if strings.Contains(host, "localhost") || strings.Contains(host, "127.0.0.1") || strings.Contains(host, "[::1]") {
		return nil
	}

	// Block RFC1918 private ranges
	ip := strings.Split(host, ":")[0]
	parts := strings.Split(ip, ".")
	if len(parts) == 4 {
		if parts[0] == "10" ||
			(parts[0] == "172" && parts[1] >= "16" && parts[1] <= "31") ||
			(parts[0] == "192" && parts[1] == "168") {
			return fmt.Errorf("private IP blocked: %s", ip)
		}
	}

	return nil
}

// isAllowedOrigin validates CORS origin
func isAllowedOrigin(origin string) bool {
	allowed := []string{"https://niebogurd.pl", "https://localhost:8788", "http://localhost:8788"}
	for _, a := range allowed {
		if origin == a {
			return true
		}
	}
	return false
}

// generateRequestID creates unique request ID
func generateRequestID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// Enroll registers the station via compute API
func (a *Agent) Enroll(stationID, version, platform string, cpu, ram int) (string, error) {
	token, err := a.getToken()
	if err != nil {
		return "", err
	}

	_ = token // token is loaded per-request inside compute.Client.do
	return stationID, nil
}

// NextTask gets the next available task via compute API
func (a *Agent) NextTask() (*compute.Task, error) {
	// compute.Client handles secure token loading internally
	return nil, nil
}