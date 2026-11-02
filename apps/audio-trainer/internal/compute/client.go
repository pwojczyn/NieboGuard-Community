package compute

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/pwojczyn/NieboGuard-Community/apps/audio-trainer/internal/secure"
)

// Client for NieboGuard Compute API (Cloudflare Worker backed)
type Client struct {
	baseURL      string
	stationID    string
	credStore    secure.CredentialStore
	tokenKey     string
	httpClient   *http.Client
	logger       *slog.Logger
}

// NewClient creates a new compute API client using secure credential storage.
// tokenKey is the key used to store/retrieve the token from the credential store.
func NewClient(baseURL, stationID, tokenKey string, credStore secure.CredentialStore) (*Client, error) {
	if credStore == nil {
		return nil, fmt.Errorf("credential store is required")
	}
	
	c := &Client{
		baseURL:    baseURL,
		stationID:  stationID,
		credStore:  credStore,
		tokenKey:   tokenKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		logger:     slog.Default().With("component", "compute-api"),
	}
	return c, nil
}

// getToken retrieves the current token from secure storage
func (c *Client) getToken() (string, error) {
	token, err := c.credStore.Load(c.tokenKey)
	if err != nil {
		return "", fmt.Errorf("loading token from secure storage: %w", err)
	}
	return token, nil
}

func (c *Client) do(method, path string, body interface{}) ([]byte, error) {
	url := c.baseURL + path

	var bodyReader io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBytes)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Load token securely from DPAPI/Credential Manager at request time
	token, err := c.getToken()
	if err == nil && token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request %s %s: %w", method, path, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// Enroll registers the station
func (c *Client) Enroll(request EnrollRequest) (*EnrollResponse, error) {
	respBody, err := c.do("POST", "/api/v2/compute/enroll", request)
	if err != nil {
		return nil, err
	}
	var resp EnrollResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("parsing enroll response: %w", err)
	}
	return &resp, nil
}

// TokenRequest requests an authorization token
func (c *Client) TokenRequest(request TokenRequest) (*TokenResponse, error) {
	respBody, err := c.do("POST", "/api/v2/compute/token", request)
	if err != nil {
		return nil, err
	}
	var resp TokenResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("parsing token response: %w", err)
	}
	return &resp, nil
}

// NextTask gets the next available task
func (c *Client) NextTask() (*Task, error) {
	respBody, err := c.do("GET", "/api/v2/compute/tasks/next", nil)
	if err != nil {
		return nil, err
	}
	var resp TaskResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("parsing next task response: %w", err)
	}
	return &resp.Task, nil
}

// Heartbeat sends task heartbeat
func (c *Client) Heartbeat(taskID string, progress float64) error {
	req := HeartbeatRequest{Progress: progress}
	_, err := c.do("POST", "/api/v2/compute/tasks/"+taskID+"/heartbeat", req)
	return err
}

// UploadURL gets signed upload URL
func (c *Client) UploadURL(taskID string) (*UploadURLResponse, error) {
	respBody, err := c.do("POST", "/api/v2/compute/tasks/"+taskID+"/upload-url", nil)
	if err != nil {
		return nil, err
	}
	var resp UploadURLResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("parsing upload url response: %w", err)
	}
	return &resp, nil
}

// CompleteTask marks task as complete
func (c *Client) CompleteTask(taskID string) error {
	_, err := c.do("POST", "/api/v2/compute/tasks/"+taskID+"/complete", nil)
	 return err
}

// FailTask marks task as failed
func (c *Client) FailTask(taskID string, reason string) error {
	req := FailTaskRequest{Reason: reason}
	_, err := c.do("POST", "/api/v2/compute/tasks/"+taskID+"/fail", req)
	return err
}

// ErrNoTasks is returned when no tasks are available
var ErrNoTasks = errors.New("no tasks available")

// Request/Response types
type EnrollRequest struct {
	StationID string `json:"station_id"`
	Version   string `json:"version"`
	Platform  string `json:"platform"`
	CPU       int    `json:"cpu"`
	RAM       int    `json:"ram"`
}

type EnrollResponse struct {
	StationID string `json:"station_id"`
	Token     string `json:"token"`
}

type TokenRequest struct {
	StationID string `json:"station_id"`
}

type TokenResponse struct {
	Token   string `json:"token"`
	Expires string `json:"expires"`
}

type HeartbeatRequest struct {
	Progress float64 `json:"progress"`
}

type UploadURLResponse struct {
	URL       string    `json:"url"`
	ExpiresAt time.Time `json:"expires_at"`
	Method    string    `json:"method"`
}

type FailTaskRequest struct {
	Reason string `json:"reason"`
}

// TaskResponse wraps a task
type TaskResponse struct {
	Task Task `json:"task"`
}

// Task represents a training task from the compute API
type Task struct {
	ID             string `json:"id"`
	StationID      string `json:"station_id"`
	Status         string `json:"status"` // pending, assigned, running, completed, failed
	DatasetID      string `json:"dataset_id"`
	DatasetVersion string `json:"dataset_version"`
	ModelID        string `json:"model_id"`
	ModelVersion   string `json:"model_version"`
	ModelSHA256    string `json:"base_model_sha256"`
	PackType       string `json:"pack_type"` // MINI, STANDARD, EXTENDED
	Epochs         int    `json:"epochs"`
	BatchSize      int    `json:"batch_size"`
	LearningRate   float64 `json:"learning_rate"`
	CpuLimit       int    `json:"cpu_limit_percent"`
	RAMLimit       int    `json:"ram_limit_gb"`
	TimeoutHours   int    `json:"timeout_hours"`
	SignedURLs     []SignedURL `json:"signed_urls"`
	CreatedAt      time.Time `json:"created_at"`
	ExpiresAt      time.Time `json:"expires_at"`
}

// SignedURL represents a signed download URL
type SignedURL struct {
	URL       string    `json:"url"`
	ObjectPath string `json:"object_path"`
	SHA256    string    `json:"sha256"`
	SizeBytes int64     `json:"size_bytes"`
	ExpiresAt time.Time `json:"expires_at"`
}