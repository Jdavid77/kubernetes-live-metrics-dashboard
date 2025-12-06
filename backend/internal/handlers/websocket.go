package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Jdavid77/kubernetes-dashboard/internal/kubernetes"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins (configure in production)
	},
}

// WebSocketHandler handles WebSocket connections for real-time metrics
type WebSocketHandler struct {
	aggregator      *kubernetes.Aggregator
	refreshInterval time.Duration
	clients         map[*websocket.Conn]bool
	broadcast       chan []byte
	register        chan *websocket.Conn
	unregister      chan *websocket.Conn
	mu              sync.Mutex
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(aggregator *kubernetes.Aggregator, refreshInterval time.Duration) *WebSocketHandler {
	handler := &WebSocketHandler{
		aggregator:      aggregator,
		refreshInterval: refreshInterval,
		clients:         make(map[*websocket.Conn]bool),
		broadcast:       make(chan []byte),
		register:        make(chan *websocket.Conn),
		unregister:      make(chan *websocket.Conn),
	}

	// Start the broadcast loop
	go handler.handleConnections()
	go handler.broadcastMetrics()

	return handler
}

// ServeHTTP handles WebSocket upgrade requests
func (h *WebSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return
	}

	// Register new client
	h.register <- conn

	// Keep connection alive and handle disconnection
	go h.handleClient(conn)
}

// handleClient handles individual client connections
func (h *WebSocketHandler) handleClient(conn *websocket.Conn) {
	defer func() {
		h.unregister <- conn
		conn.Close()
	}()

	// Read messages from client (for keep-alive)
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

// handleConnections manages client registrations and broadcasts
func (h *WebSocketHandler) handleConnections() {
	for {
		select {
		case conn := <-h.register:
			h.mu.Lock()
			h.clients[conn] = true
			h.mu.Unlock()
			log.Printf("New WebSocket client connected. Total clients: %d", len(h.clients))

		case conn := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[conn]; ok {
				delete(h.clients, conn)
				conn.Close()
			}
			h.mu.Unlock()
			log.Printf("WebSocket client disconnected. Total clients: %d", len(h.clients))

		case message := <-h.broadcast:
			h.mu.Lock()
			for conn := range h.clients {
				err := conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("Error broadcasting to client: %v", err)
					conn.Close()
					delete(h.clients, conn)
				}
			}
			h.mu.Unlock()
		}
	}
}

// broadcastMetrics periodically fetches and broadcasts metrics
func (h *WebSocketHandler) broadcastMetrics() {
	ticker := time.NewTicker(h.refreshInterval)
	defer ticker.Stop()

	for range ticker.C {
		h.mu.Lock()
		clientCount := len(h.clients)
		h.mu.Unlock()

		// Only fetch metrics if there are connected clients
		if clientCount == 0 {
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		metrics, err := h.aggregator.GetClusterMetrics(ctx)
		cancel()

		if err != nil {
			log.Printf("Error fetching metrics for broadcast: %v", err)
			continue
		}

		// Create message
		message := map[string]interface{}{
			"type": "metrics_update",
			"data": metrics,
		}

		messageBytes, err := json.Marshal(message)
		if err != nil {
			log.Printf("Error marshaling metrics: %v", err)
			continue
		}

		// Broadcast to all clients
		h.broadcast <- messageBytes
	}
}
