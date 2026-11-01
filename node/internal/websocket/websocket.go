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
	"github.com/pwojczyn/NieboGuard-Community/node/internal/collectors"
	"github.com/pwojczyn/NieboGuard-Community/node/internal/storage"
	"github.com/pwojczyn/NieboGuard-Community/node/pkg/config"
)

// Connection represents a WebSocket client connection.
type Connection struct {
	conn         *websocket.Conn
	clientID     string
	mu           sync.Mutex
	subscriptions []string
}

// Server represents the WebSocket server component.
type Server struct {
	cfg          *config.Configuration
	storage      *storage.Storage
	collectorMgr *collectors.Manager
	logger       *slog.Logger
	server       *http.Server
	clients      map[string]*Connection
	mu           sync.Mutex
	shutdown     chan struct{}
	stopped      chan struct{}
}

// New creates a new WebSocket server instance.
func New(cfg *config.Configuration, store *storage.Storage, collMgr *collectors.Manager, logger *slog.Logger) (*Server, error) {
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	s := &Server{
		cfg:          cfg,
		storage:      store,
		collectorMgr: collMgr,
		logger:       logger,
		clients:      make(map[string]*Connection),
		shutdown:     make(chan struct{}),
		stopped:      make(chan struct{}),
	}

	s.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Host, cfg.Port+100),
		Handler: http.HandlerFunc(s.handleStream),
	}

	return s, nil
}

// Start begins accepting WebSocket connections.
func (s *Server) Start(port int, host string) error {
	addr := fmt.Sprintf("%s:%d", host, port+100)
	s.logger.Info("starting websocket server", "addr", addr)

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("websocket server failed", "error", err)
		}
	}()
	return nil
}

// Stop gracefully shuts down the WebSocket server and closes all client connections.
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
	s.clients = make(map[string]*Connection)
	s.mu.Unlock()

	close(s.shutdown)
	close(s.stopped)
}

// handleStream handles WebSocket upgrade and connection management.
func (s *Server) handleStream(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Error("websocket upgrade error", "error", err)
		http.Error(w, "upgrade failed", http.StatusBadRequest)
		return
	}

	clientID := fmt.Sprintf("client_%d", time.Now().UnixNano())
	wsConn := &Connection{
		conn: conn,
		clientID: clientID,
	}

	s.mu.Lock()
	s.clients[clientID] = wsConn
	s.mu.Unlock()

	// Send initial snapshot of recent data
	wsConn.sendSnapshot()

	// Start message handler goroutine
	go wsConn.readLoop(s.logger)

	// Start keepalive pings
	go wsConn.sendPings()

	s.logger.Info("websocket client connected", "client_id", clientID)
}

// upgrader upgrades HTTP connections to WebSocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return r.Host == "localhost"
	},
}

// readLoop reads messages from the WebSocket connection.
func (conn *Connection) readLoop(logger *slog.Logger) {
	defer func() {
		conn.conn.Close()
	}()

	for {
		conn.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		mt, msg, err := conn.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				logger.Error("websocket read error", "client_id", conn.clientID)
			}
			return
		}

		if mt == websocket.TextMessage {
			conn.handleIncomingMessage(msg)
		}
	}
}

// handleIncomingMessage processes incoming client messages.
func (conn *Connection) handleIncomingMessage(msg []byte) {
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
						conn.subscriptions = append(conn.subscriptions, s)
					}
				}
			}
		case "unsubscribe":
			if topics, ok := data["topics"].([]interface{}); ok {
				for _, t := range topics {
					if s, ok := t.(string); ok {
						for i, sub := range conn.subscriptions {
							if sub == s {
								conn.subscriptions = append(conn.subscriptions[:i], conn.subscriptions[i+1:]...)
								break
							}
						}
					}
				}
			}
		}
	}
}

// sendSnapshot sends an initial data snapshot to the connected client.
func (conn *Connection) sendSnapshot() {
	// In a real implementation, this would get data from storage
	// For now, send an empty snapshot
	message := map[string]interface{}{
		"type": "snapshot",
		"data": []interface{}{},
		"seq":  1,
	}
	conn.sendJSON(message)
}

// sendPings periodically sends ping frames to keep the connection alive.
func (conn *Connection) sendPings() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if err := conn.conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(5*time.Second)); err != nil {
			return
		}
	}
}

// sendJSON serializes and sends a JSON message to the client.
func (conn *Connection) sendJSON(msg map[string]interface{}) error {
	conn.mu.Lock()
	defer conn.mu.Unlock()

	if err := conn.conn.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
		return err
	}

	return conn.conn.WriteJSON(msg)
}