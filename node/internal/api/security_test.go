package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func testSecurity(limit int) *securityConfig {
	return &securityConfig{token: "secret-token", require: true, limit: limit, window: time.Minute, buckets: make(map[string]rateBucket)}
}

func hit(s *Server, method, path, remote, token string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, path, nil)
	r.RemoteAddr = remote
	if token != "" {
		r.Header.Set(localAuthHeader, token)
	}
	w := httptest.NewRecorder()
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) })
	s.securityMiddleware(next).ServeHTTP(w, r)
	return w
}

func TestSecurityPublicProbeAndLoopback(t *testing.T) {
	s := &Server{security: testSecurity(20)}
	if got := hit(s, http.MethodGet, "/api/v1/health", "127.0.0.1:4000", "").Code; got != http.StatusNoContent {
		t.Fatalf("health = %d", got)
	}
	if got := hit(s, http.MethodGet, "/api/v1/objects", "127.0.0.1:4000", "").Code; got != http.StatusUnauthorized {
		t.Fatalf("objects without token = %d", got)
	}
	if got := hit(s, http.MethodGet, "/api/v1/objects", "127.0.0.1:4000", "wrong").Code; got != http.StatusUnauthorized {
		t.Fatalf("wrong token = %d", got)
	}
	if got := hit(s, http.MethodGet, "/api/v1/objects", "127.0.0.1:4000", "secret-token").Code; got != http.StatusNoContent {
		t.Fatalf("objects token = %d", got)
	}
	if got := hit(s, http.MethodGet, "/api/v1/objects", "192.0.2.10:4000", "").Code; got != http.StatusUnauthorized {
		t.Fatalf("remote without token = %d", got)
	}
}

func TestSecurityAdminRequiresConfiguredToken(t *testing.T) {
	s := &Server{security: &securityConfig{require: false, limit: 20, window: time.Minute, buckets: make(map[string]rateBucket)}}
	if got := hit(s, http.MethodPost, "/api/v1/admin/sources", "127.0.0.1:4000", "").Code; got != http.StatusForbidden {
		t.Fatalf("admin without config = %d", got)
	}
}

func TestSecurityRateLimit(t *testing.T) {
	s := &Server{security: testSecurity(2)}
	if got := hit(s, http.MethodGet, "/api/v1/health", "127.0.0.1:4000", "").Code; got != http.StatusNoContent {
		t.Fatalf("first = %d", got)
	}
	if got := hit(s, http.MethodGet, "/api/v1/health", "127.0.0.1:4000", "").Code; got != http.StatusNoContent {
		t.Fatalf("second = %d", got)
	}
	if got := hit(s, http.MethodGet, "/api/v1/health", "127.0.0.1:4000", "").Code; got != http.StatusTooManyRequests {
		t.Fatalf("third = %d", got)
	}
}

func TestLoopbackAddress(t *testing.T) {
	if !loopbackAddress("127.0.0.1:1") || !loopbackAddress("[::1]:1") || loopbackAddress("192.0.2.1:1") {
		t.Fatal("loopback parse")
	}
}

func TestLoopbackHost(t *testing.T) {
	if !loopbackHost("127.0.0.1") || !loopbackHost("localhost") || !loopbackHost("::1") || loopbackHost("0.0.0.0") || loopbackHost("example.com") {
		t.Fatal("bind host policy")
	}
}
