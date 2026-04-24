package ws

import (
	"net/http"
	"net/url"
)

type LogStreamHandler struct {
	hub    *Hub
	getLog func(agentID string) (string, error)
}

func NewLogStreamHandler(hub *Hub, getLogFunc func(agentID string) (string, error)) *LogStreamHandler {
	return &LogStreamHandler{
		hub:    hub,
		getLog: getLogFunc,
	}
}

func (h *LogStreamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	agentID := r.URL.Query().Get("agent")
	if agentID == "" {
		http.Error(w, "missing agent parameter", http.StatusBadRequest)
		return
	}

	decodedAgentID, err := url.PathUnescape(agentID)
	if err != nil {
		http.Error(w, "invalid agent ID", http.StatusBadRequest)
		return
	}

	if _, err := h.getLog(decodedAgentID); err != nil {
		http.Error(w, "agent not found", http.StatusNotFound)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	h.hub.Subscribe(conn, decodedAgentID)
	defer h.hub.Unsubscribe(conn, decodedAgentID)

	<-make(chan struct{})
}