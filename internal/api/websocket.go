package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type SafeWebSocketConn struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

func (s *SafeWebSocketConn) WriteMessage(messageType int, data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.conn.WriteMessage(messageType, data)
}

func (s *SafeWebSocketConn) ReadMessage() (messageType int, p []byte, err error) {
	return s.conn.ReadMessage()
}

func (s *SafeWebSocketConn) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.conn.Close()
}

type WebSocketHub struct {
	clients    map[*SafeWebSocketConn]bool
	broadcast  chan []byte
	register   chan *SafeWebSocketConn
	unregister chan *SafeWebSocketConn
	mu         sync.RWMutex
}

func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{
		clients:    make(map[*SafeWebSocketConn]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *SafeWebSocketConn),
		unregister: make(chan *SafeWebSocketConn),
	}
}

func (h *WebSocketHub) Run() {
	for {
		select {
		case conn := <-h.register:
			h.mu.Lock()
			h.clients[conn] = true
			h.mu.Unlock()

		case conn := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[conn]; ok {
				delete(h.clients, conn)
				conn.Close()
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for conn := range h.clients {
				if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
					go func(c *SafeWebSocketConn) {
						h.unregister <- c
					}(conn)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *WebSocketHub) SendJSON(v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		return
	}
	select {
	case h.broadcast <- data:
	default:
	}
}

func (s *Server) DashboardWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket upgrade error: %v", err)
		return
	}

	sConn := &SafeWebSocketConn{conn: conn}
	s.wsHub.register <- sConn

	// send stats immediately on connect
	stats := s.CollectStats()
	if data, err := json.Marshal(stats); err == nil {
		sConn.WriteMessage(websocket.TextMessage, data)
	}

	ticker := time.NewTicker(2 * time.Second)
	defer func() {
		ticker.Stop()
		s.wsHub.unregister <- sConn
	}()

	// read goroutine — detects client disconnect
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			if _, _, err := sConn.ReadMessage(); err != nil {
				return
			}
		}
	}()

	for {
		select {
		case <-ticker.C:
			stats := s.CollectStats()
			data, err := json.Marshal(stats)
			if err != nil {
				continue
			}
			if err := sConn.WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}
		case <-done:
			return
		}
	}
}
