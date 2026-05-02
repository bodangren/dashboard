package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type ActivityEvent struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	Repo      string          `json:"repo"`
	Message   string          `json:"message"`
	Timestamp time.Time       `json:"timestamp"`
	Metadata  json.RawMessage `json:"metadata,omitempty"`
}

type ActivityHub struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan ActivityEvent
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.Mutex
	done       chan struct{}
}

func NewActivityHub() *ActivityHub {
	return &ActivityHub{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan ActivityEvent, 10),
		register:   make(chan *websocket.Conn, 10),
		unregister: make(chan *websocket.Conn, 10),
		done:       make(chan struct{}),
	}
}

func (h *ActivityHub) Start() {
	go h.run()
}

func (h *ActivityHub) run() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("activityHub: panic recovered: %v\n%s", r, string(debug.Stack()))
		}
	}()
	for {
		select {
		case <-h.done:
			return
		case conn := <-h.register:
			h.mu.Lock()
			h.clients[conn] = true
			h.mu.Unlock()
		case conn := <-h.unregister:
			h.mu.Lock()
			delete(h.clients, conn)
			h.mu.Unlock()
		case event := <-h.broadcast:
			h.mu.Lock()
			for conn := range h.clients {
				func() {
					defer func() {
						if r := recover(); r != nil {
							conn.Close()
							delete(h.clients, conn)
						}
					}()
					conn.SetWriteDeadline(time.Now().Add(time.Second))
					err := conn.WriteJSON(event)
					if err != nil {
						conn.Close()
						delete(h.clients, conn)
					}
				}()
			}
			h.mu.Unlock()
		}
	}
}

func (h *ActivityHub) Stop() {
	close(h.done)
	h.mu.Lock()
	for conn := range h.clients {
		conn.Close()
	}
	h.mu.Unlock()
}

func (h *ActivityHub) Broadcast(event ActivityEvent) {
	h.broadcast <- event
}

func (h *ActivityHub) Register(conn *websocket.Conn) {
	h.register <- conn
}

func (h *ActivityHub) Unregister(conn *websocket.Conn) {
	h.unregister <- conn
}

func (h *ActivityHub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	h.Register(conn)
}