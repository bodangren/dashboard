package ws

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestHub_RegisterAndBroadcast(t *testing.T) {
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