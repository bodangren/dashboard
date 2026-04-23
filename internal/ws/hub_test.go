package ws

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestHub_ServeHTTP_ValidUpgrade(t *testing.T) {
	h := NewHub()
	h.Start()
	defer h.Stop()

	srv := httptest.NewServer(h)
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("connection failed: %v", err)
	}
	defer conn.Close()

	time.Sleep(50 * time.Millisecond)

	h.mu.Lock()
	if len(h.clients) != 1 {
		t.Errorf("expected 1 client after connection, got %d", len(h.clients))
	}
	h.mu.Unlock()
}

func TestHub_ServeHTTP_MultipleClients(t *testing.T) {
	h := NewHub()
	h.Start()
	defer h.Stop()

	srv := httptest.NewServer(h)
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")

	conn1, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("conn1 connection failed: %v", err)
	}
	defer conn1.Close()

	conn2, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("conn2 connection failed: %v", err)
	}
	defer conn2.Close()

	time.Sleep(50 * time.Millisecond)

	h.mu.Lock()
	if len(h.clients) != 2 {
		t.Errorf("expected 2 clients after connections, got %d", len(h.clients))
	}
	h.mu.Unlock()
}

func TestHub_ServeHTTP_InvalidOrigin(t *testing.T) {
	h := NewHub()
	h.Start()
	defer h.Stop()

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return false
	}
	defer func() {
		upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}
	}()

	srv := httptest.NewServer(h)
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")
	_, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err == nil {
		t.Error("expected connection failure due to origin check, got nil")
	}
}

func TestHub_ServeHTTP_BroadcastToClient(t *testing.T) {
	h := NewHub()
	h.Start()
	defer h.Stop()

	srv := httptest.NewServer(h)
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("connection failed: %v", err)
	}
	defer conn.Close()

	time.Sleep(50 * time.Millisecond)

	entry := LogEntry{
		AgentID:   "test-agent",
		Timestamp: "2026-04-23T10:00:00Z",
		Message:   "test message",
		Type:      "stdout",
	}

	h.Broadcast(entry)

	var received LogEntry
	if err := conn.ReadJSON(&received); err != nil {
		t.Fatalf("failed to read broadcast: %v", err)
	}

	if received.AgentID != entry.AgentID {
		t.Errorf("expected AgentID %q, got %q", entry.AgentID, received.AgentID)
	}
	if received.Message != entry.Message {
		t.Errorf("expected Message %q, got %q", entry.Message, received.Message)
	}
	if received.Type != entry.Type {
		t.Errorf("expected Type %q, got %q", entry.Type, received.Type)
	}
}

func TestHub_Broadcast_NonBlocking(t *testing.T) {
	h := NewHub()
	h.Start()
	defer h.Stop()

	entry := LogEntry{
		AgentID:   "test-agent",
		Timestamp: "2026-04-23T10:00:00Z",
		Message:   "non-blocking test",
		Type:      "stdout",
	}

	done := make(chan struct{})
	go func() {
		h.Broadcast(entry)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Error("Broadcast blocked unexpectedly")
	}
}

func TestHub_Broadcast_MultipleClients(t *testing.T) {
	h := NewHub()
	h.Start()
	defer h.Stop()

	srv := httptest.NewServer(h)
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")

	conn1, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("conn1 connection failed: %v", err)
	}
	defer conn1.Close()

	conn2, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("conn2 connection failed: %v", err)
	}
	defer conn2.Close()

	time.Sleep(50 * time.Millisecond)

	entry := LogEntry{
		AgentID:   "multi-agent",
		Timestamp: "2026-04-23T10:00:00Z",
		Message:   "hello all",
		Type:      "stdout",
	}

	h.Broadcast(entry)

	var received1, received2 LogEntry
	if err := conn1.ReadJSON(&received1); err != nil {
		t.Fatalf("conn1 failed to read broadcast: %v", err)
	}
	if err := conn2.ReadJSON(&received2); err != nil {
		t.Fatalf("conn2 failed to read broadcast: %v", err)
	}

	if received1.Message != entry.Message {
		t.Errorf("conn1 expected Message %q, got %q", entry.Message, received1.Message)
	}
	if received2.Message != entry.Message {
		t.Errorf("conn2 expected Message %q, got %q", entry.Message, received2.Message)
	}
}



