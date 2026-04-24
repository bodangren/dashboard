package ws

import (
	"log"
	"net/http"
	"runtime/debug"
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
	clients       map[*websocket.Conn]bool
	subscriptions map[string][]*websocket.Conn
	broadcast     chan LogEntry
	register      chan *websocket.Conn
	unregister    chan *websocket.Conn
	subscribe     chan subscribeMsg
	unsubscribe   chan subscribeMsg
	mu            sync.Mutex
	done          chan struct{}
}

type subscribeMsg struct {
	conn    *websocket.Conn
	agentID string
}

func NewHub() *Hub {
	return &Hub{
		clients:       make(map[*websocket.Conn]bool),
		subscriptions: make(map[string][]*websocket.Conn),
		broadcast:     make(chan LogEntry, 10),
		register:      make(chan *websocket.Conn, 10),
		unregister:    make(chan *websocket.Conn, 10),
		subscribe:     make(chan subscribeMsg, 10),
		unsubscribe:   make(chan subscribeMsg, 10),
		done:          make(chan struct{}),
	}
}

func (h *Hub) Start() {
	go h.run()
}

func (h *Hub) run() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("hub: panic recovered: %v\n%s", r, string(debug.Stack()))
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
			for agentID := range h.subscriptions {
				h.removeConnFromAgentUnlocked(conn, agentID)
			}
			h.mu.Unlock()
		case msg := <-h.subscribe:
			h.mu.Lock()
			h.clients[msg.conn] = true
			h.subscriptions[msg.agentID] = append(h.subscriptions[msg.agentID], msg.conn)
			h.mu.Unlock()
		case msg := <-h.unsubscribe:
			h.mu.Lock()
			delete(h.clients, msg.conn)
			h.removeConnFromAgentUnlocked(msg.conn, msg.agentID)
			h.mu.Unlock()
		case entry := <-h.broadcast:
			h.mu.Lock()
			if subs, ok := h.subscriptions[entry.AgentID]; ok && len(subs) > 0 {
				for _, conn := range subs {
					func() {
						defer func() {
							if r := recover(); r != nil {
								conn.Close()
								h.removeConnFromAgentUnlocked(conn, entry.AgentID)
								delete(h.clients, conn)
							}
						}()
						conn.SetWriteDeadline(time.Now().Add(time.Second))
						err := conn.WriteJSON(entry)
						if err != nil {
							conn.Close()
							h.removeConnFromAgentUnlocked(conn, entry.AgentID)
							delete(h.clients, conn)
						}
					}()
				}
			} else {
				for conn := range h.clients {
					func() {
						defer func() {
							if r := recover(); r != nil {
								conn.Close()
								delete(h.clients, conn)
							}
						}()
						conn.SetWriteDeadline(time.Now().Add(time.Second))
						err := conn.WriteJSON(entry)
						if err != nil {
							conn.Close()
							delete(h.clients, conn)
						}
					}()
				}
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) removeConnFromAgentUnlocked(conn *websocket.Conn, agentID string) {
	conns := h.subscriptions[agentID]
	for i, c := range conns {
		if c == conn {
			h.subscriptions[agentID] = append(conns[:i], conns[i+1:]...)
			break
		}
	}
	if len(h.subscriptions[agentID]) == 0 {
		delete(h.subscriptions, agentID)
	}
}

func (h *Hub) Stop() {
	close(h.done)
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

func (h *Hub) Subscribe(conn *websocket.Conn, agentID string) {
	h.subscribe <- subscribeMsg{conn: conn, agentID: agentID}
}

func (h *Hub) Unsubscribe(conn *websocket.Conn, agentID string) {
	h.unsubscribe <- subscribeMsg{conn: conn, agentID: agentID}
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
