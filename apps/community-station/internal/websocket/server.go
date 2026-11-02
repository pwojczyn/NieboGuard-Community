package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pwojczyn/NieboGuard-Community/apps/community-station/internal/collectors"
	"github.com/pwojczyn/NieboGuard-Community/apps/community-station/internal/config"
	"github.com/pwojczyn/NieboGuard-Community/apps/community-station/internal/storage"
)

// Server represents the WebSocket server component.
type Server struct {
	cfg          *config.Configuration
	storage      *storage.Storage
	collector    *collectors.Manager
	logger       *slog.Logger
	server       *http.Server
	clients      map[string]*Connection
	mu           sync.RWMutex
	shutdown     chan struct{}
	stopped      chan struct{}
}

// New creates a new WebSocket server instance.
func New(cfg *config.Configuration, store *storage.Storage, collector *collectors.Manager, logger *slog.Logger) (*Server, error) {
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	s := &Server{
		cfg:         cfg,
		storage:     store,
		collector:   collector,
		logger:      logger,
		clients:     make(map[string]*Connection),
		shutdown:    make(chan struct{}),
		stopped:     make(chan struct{}),
	}

	return s, nil
}

// Start begins accepting WebSocket connections.
func (s *Server) Start(port int, host string) error {
	s.logger.Info("starting WebSocket server", "addr", fmt.Sprintf("%s:%d", host, port+100))

	go func() {
		s.server = &http.Server{
			Addr:    fmt.Sprintf("%s:%d", host, port+100),
			Handler: http.HandlerFunc(s.websocketHandle),
		}
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("WebSocket server failed", "error", err)
		}
	}()

	return nil
}

// Stop gracefully shuts down the WebSocket server.
func (s *Server) Stop() {
	if s.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.server.Shutdown(ctx)
	}

	s.mu.Lock()
	for _, conn := range s.clients {
		conn.conn.Close()
	}
	s.mu.Unlock()

	close(s.shutdown)
	close(s.stopped)
}

// websocketHandle handles WebSocket connections.
func (s *Server) websocketHandle(w http.ResponseWriter, r *http.Request) {
	// Upgrade to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Error("WebSocket upgrade error", "error", err)
		http.Error(w, "upgrade failed", http.StatusBadRequest)
		return
	}

	// Create connection
	clientID := fmt.Sprintf("client_%d", time.Now().UnixNano())
	clientConn := &Connection{
		conn:     conn,
		clientID: clientID,
	}
	s.mu.Lock()
	s.clients[clientID] = clientConn
	s.mu.Unlock()

	// Send initial snapshot
	clientConn.sendSnapshot()

	// Start message handling
	go clientConn.readLoop(s.logger)
}

// upgrader upgrades HTTP connections to WebSocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return r.Host == "localhost"
	},
}

// Connection represents a WebSocket client connection.
type Connection struct {
	conn         *websocket.Conn
	clientID     string
	mu           sync.Mutex
	subscriptions []string
}

// sendSnapshot sends an initial data snapshot to the connected client.
func (c *Connection) sendSnapshot() {
	message := map[string]interface{}{
		"type": "snapshot",
		"data": []interface{}{},
		"seq":  1,
	}
	c.sendJSON(message)
}

// sendPings periodically sends ping frames to keep the connection alive.
func (c *Connection) sendPings() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if err := c.conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(5*time.Second)); err != nil {
			return
		}
	}
}

// sendJSON serializes and sends a JSON message to the client.
func (c *Connection) sendJSON(msg map[string]interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.conn.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
		return err
	}
	return c.conn.WriteJSON(msg)
}

// readLoop reads messages from the WebSocket connection.
func (c *Connection) readLoop(logger *slog.Logger) {
	defer func() {
		c.conn.Close()
	}()

	for {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		mt, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				logger.Error("websocket read error", "client_id", c.clientID)
			}
			return
		}

		if mt == websocket.TextMessage {
			c.handleIncomingMessage(msg)
		}
	}
}

// handleIncomingMessage processes incoming client messages.
func (c *Connection) handleIncomingMessage(msg []byte) {
	var data map[string]interface{}
	if err := json.Unmarshal(msg, &data); err != nil {
		return
	}

	if op, ok := data["op"].(string); ok {
		switch op {
		case "subscribe":
			if topics, ok := data["topics"].([]interface{}); ok {
				for _, t := range topics {
					if s, ok := t.(string); ok {
						c.mu.Lock()
						c.subscriptions = append(c.subscriptions, s)
						c.mu.Unlock()
					}
				}
			}
		case "unsubscribe":
			if topics, ok := data["topics"].([]interface{}); ok {
				for _, t := range topics {
					if s, ok := t.(string); ok {
						c.mu.Lock()
						for i, sub := range c.subscriptions {
							if sub == s {
								c.subscriptions = append(c.subscriptions[:i], c.subscriptions[i+1:]...)
								c.mu.Unlock()
								break
							}
						}
						c.mu.Unlock()
					}
				}
			}
		}
	}
}
