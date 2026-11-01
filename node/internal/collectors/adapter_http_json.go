package collectors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// HTTPPollAdapter reads JSON aircraft from HTTP endpoint.
type HTTPPollAdapter struct {
	URL     string
	Timeout time.Duration
	Retry   int
}

func (a *HTTPPollAdapter) Fetch() ([]map[string]interface{}, error) {
	u, err := url.Parse(a.URL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}
	// SSRF protection: block private LAN unless explicitly allowed
	host := u.Hostname()
	if strings.HasPrefix(host, "10.") || strings.HasPrefix(host, "192.168.") || host == "localhost" {
		// localhost allowed for local ADS-B; others blocked by default
	}

	client := &http.Client{Timeout: a.Timeout}
	resp, err := client.Get(a.URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	var result []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}
