package api

import (
	"encoding/json"
	"time"
)

type EventType string

const (
	EventTypeCommit EventType = "commit"
	EventTypeAgent  EventType = "agent"
	EventTypePull   EventType = "pull"
)

type ActivityEvent struct {
	ID        string          `json:"id"`
	Type      EventType       `json:"type"`
	Repo      string          `json:"repo"`
	Message   string          `json:"message"`
	Timestamp time.Time       `json:"timestamp"`
	Summary   string          `json:"summary,omitempty"`
	Flags     []string        `json:"flags,omitempty"`
	Metadata  json.RawMessage `json:"metadata,omitempty"`
}

type CommitEventMetadata struct {
	Hash    string `json:"hash"`
	Author  string `json:"author"`
	Count   int    `json:"count,omitempty"`
}

type AgentEventMetadata struct {
	AgentID   string `json:"agent_id"`
	AgentName string `json:"agent_name"`
	Status    string `json:"status"`
	ExitCode  int    `json:"exit_code,omitempty"`
	Error     string `json:"error,omitempty"`
}

type PullEventMetadata struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}