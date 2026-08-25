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

// Hub manages active WebSocket connections and broadcasts messages to them.
type Hub struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.Mutex
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

// run processes registrations and broadcasts until ctx is cancelled.
func (h *Hub) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			h.mu.Lock()
			for conn := range h.clients {
				_ = conn.Close()
			}
			h.mu.Unlock()
			return

		case conn := <-h.register:
			h.mu.Lock()
			h.clients[conn] = true
			h.mu.Unlock()
			log.Printf("WebSocket client connected. Total: %d", len(h.clients))

		case conn := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[conn]; ok {
				delete(h.clients, conn)
				_ = conn.Close()
			}
			h.mu.Unlock()
			log.Printf("WebSocket client disconnected. Total: %d", len(h.clients))

		case message := <-h.broadcast:
			h.mu.Lock()
			for conn := range h.clients {
				if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
					log.Printf("WebSocket write error: %v", err)
					_ = conn.Close()
					delete(h.clients, conn)
				}
			}
			h.mu.Unlock()
		}
	}
}

// broadcastMetrics fetches metrics on each tick and sends to the hub.
func broadcastMetrics(ctx context.Context, hub *Hub, src kubernetes.MetricsSource, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			hub.mu.Lock()
			clientCount := len(hub.clients)
			hub.mu.Unlock()
			if clientCount == 0 {
				continue
			}

			fetchCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			metrics, err := src.GetClusterMetrics(fetchCtx)
			cancel()
			if err != nil {
				log.Printf("Error fetching metrics for broadcast: %v", err)
				continue
			}

			msg, err := json.Marshal(map[string]any{
				"type": "metrics_update",
				"data": metrics,
			})
			if err != nil {
				log.Printf("Error marshaling metrics: %v", err)
				continue
			}

			select {
			case hub.broadcast <- msg:
			default:
				// drop if the hub's buffer is full (slow consumer)
			}
		}
	}
}

// WebSocketHandler upgrades HTTP connections and hands them to the hub.
type WebSocketHandler struct {
	hub *Hub
}

// NewWebSocketHandler creates a handler and starts the hub and ticker goroutines.
// Both goroutines stop when ctx is cancelled.
func NewWebSocketHandler(ctx context.Context, aggregator kubernetes.MetricsSource, refreshInterval time.Duration) *WebSocketHandler {
	hub := newHub()
	go hub.run(ctx)
	go broadcastMetrics(ctx, hub, aggregator, refreshInterval)
	return &WebSocketHandler{hub: hub}
}

// ServeHTTP upgrades the connection and keeps it alive until the client disconnects.
func (h *WebSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	h.hub.register <- conn

	// Drain incoming frames; unregister on disconnect.
	go func() {
		defer func() { h.hub.unregister <- conn }()
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}
	}()
}
