package ws

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type LogEntry struct {
	AgentID   string `json:"agent_id"`
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
	Type      string `json:"type"`
}

type Hub struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan LogEntry
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan LogEntry, 10),
		register:   make(chan *websocket.Conn, 10),
		unregister: make(chan *websocket.Conn, 10),
	}
}

func (h *Hub) Start() {
	go h.run()
}

func (h *Hub) run() {
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	for {
		select {
		case conn := <-h.register:
			h.mu.Lock()
			h.clients[conn] = true
			h.mu.Unlock()
		case conn := <-h.unregister:
			h.mu.Lock()
			delete(h.clients, conn)
			h.mu.Unlock()
		case entry := <-h.broadcast:
			h.mu.Lock()
			for conn := range h.clients {
				conn.SetWriteDeadline(time.Now().Add(time.Second))
				err := conn.WriteJSON(entry)
				if err != nil {
					delete(h.clients, conn)
				}
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) Stop() {
	close(h.broadcast)
	close(h.register)
	close(h.unregister)
	h.mu.Lock()
	for conn := range h.clients {
		conn.Close()
	}
	h.mu.Unlock()
}

func (h *Hub) Broadcast(entry LogEntry) {
	h.broadcast <- entry
}

func (h *Hub) Register(conn *websocket.Conn) {
	h.register <- conn
}

func (h *Hub) Unregister(conn *websocket.Conn) {
	h.unregister <- conn
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	h.Register(conn)
}