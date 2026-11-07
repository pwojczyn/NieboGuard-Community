package api

import (
	"crypto/subtle"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const localAuthHeader = "X-NG-Local-Auth"

type rateBucket struct {
	start time.Time
	count int
}

type securityConfig struct {
	token     string
	require   bool
	limit     int
	window    time.Duration
	mu        sync.Mutex
	buckets   map[string]rateBucket
	lastPrune time.Time
}

func firstEnv(names ...string) string {
	for _, name := range names {
		if v := strings.TrimSpace(os.Getenv(name)); v != "" {
			return v
		}
	}
	return ""
}

func loopbackHost(host string) bool {
	h := strings.Trim(strings.TrimSpace(host), "[]")
	if h == "localhost" {
		return true
	}
	ip := net.ParseIP(h)
	return ip != nil && ip.IsLoopback()
}

func newSecurityConfig() *securityConfig {
	token := firstEnv("NNODE_AUTH_TOKEN", "NB_NODE_AUTH_TOKEN")
	require := token != "" || firstEnv("NNODE_REQUIRE_AUTH", "NB_NODE_REQUIRE_AUTH") == "1"
	return &securityConfig{
		token: token, require: require, limit: 120, window: time.Minute,
		buckets: make(map[string]rateBucket),
	}
}

func loopbackAddress(remote string) bool {
	host, _, err := net.SplitHostPort(remote)
	if err != nil {
		host = remote
	}
	ip := net.ParseIP(strings.Trim(host, "[]"))
	return ip != nil && ip.IsLoopback()
}

func (s *securityConfig) authorized(r *http.Request) bool {
	if !s.require {
		return true
	}
	got := strings.TrimSpace(r.Header.Get(localAuthHeader))
	if got == "" {
		auth := strings.TrimSpace(r.Header.Get("Authorization"))
		if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
			got = strings.TrimSpace(auth[7:])
		}
	}
	if len(got) != len(s.token) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(got), []byte(s.token)) == 1
}

func (s *securityConfig) allow(remote string, now time.Time) (bool, int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if now.Sub(s.lastPrune) > 10*time.Minute {
		for k, b := range s.buckets {
			if now.Sub(b.start) > 2*s.window {
				delete(s.buckets, k)
			}
		}
		s.lastPrune = now
	}
	b := s.buckets[remote]
	if b.start.IsZero() || now.Sub(b.start) >= s.window {
		b = rateBucket{start: now, count: 0}
	}
	if b.count >= s.limit {
		s.buckets[remote] = b
		return false, int(s.window.Seconds())
	}
	b.count++
	s.buckets[remote] = b
	return true, 0
}

func (s *Server) securityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Cache-Control", "no-store")
		remote := r.RemoteAddr
		ok, retry := s.security.allow(remote, time.Now())
		if !ok {
			w.Header().Set("Retry-After", strconvItoa(retry))
			http.Error(w, `{"error":{"code":"rate_limited","message":"Too many requests"}}`, http.StatusTooManyRequests)
			return
		}
		path := r.URL.Path
		// Health/readiness are safe local probes; all data and admin paths are protected.
		publicProbe := path == "/api/v1/health" || path == "/api/v1/ready" || path == "/api/v1/version"
		if !publicProbe {
			if !loopbackAddress(remote) && s.security.token == "" {
				http.Error(w, `{"error":{"code":"loopback_required","message":"Local API is loopback-only"}}`, http.StatusForbidden)
				return
			}
			if strings.HasPrefix(path, "/api/v1/admin/") && s.security.token == "" {
				http.Error(w, `{"error":{"code":"auth_not_configured","message":"Admin API requires NNODE_AUTH_TOKEN"}}`, http.StatusForbidden)
				return
			}
			if !s.security.authorized(r) {
				w.Header().Set("WWW-Authenticate", "Bearer")
				http.Error(w, `{"error":{"code":"unauthorized","message":"Invalid local auth token"}}`, http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func strconvItoa(v int) string {
	if v <= 0 {
		return "60"
	}
	const digits = "0123456789"
	b := make([]byte, 0, 4)
	for v > 0 {
		b = append([]byte{digits[v%10]}, b...)
		v /= 10
	}
	return string(b)
}
