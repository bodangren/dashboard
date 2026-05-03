package api

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"dashboard/internal/agents"
	"dashboard/internal/ai"
	"dashboard/internal/ws"
)

type ActivityHandler struct {
	repos             []string
	getCommits        GetCommitsFunc
	readCrontab       agents.ReadFunc
	stateMap          *agents.AgentStateMap
	pullMu            *sync.RWMutex
	lastPullTime      map[string]time.Time
	lastPullErr       map[string]string
	recentAgentEvents []ActivityEvent
	eventMu           sync.Mutex
	activityHub       *ws.ActivityHub
	enhancer          *ai.ActivityEnhancer
}

type ActivityHandlerOption func(*ActivityHandler)

func WithActivityRepos(repos []string) ActivityHandlerOption {
	return func(h *ActivityHandler) { h.repos = repos }
}

func WithActivityGetCommits(fn GetCommitsFunc) ActivityHandlerOption {
	return func(h *ActivityHandler) { h.getCommits = fn }
}

func WithActivityReadCrontab(fn agents.ReadFunc) ActivityHandlerOption {
	return func(h *ActivityHandler) { h.readCrontab = fn }
}

func WithActivityStateMap(sm *agents.AgentStateMap) ActivityHandlerOption {
	return func(h *ActivityHandler) { h.stateMap = sm }
}

func WithActivityPullState(pullMu *sync.RWMutex, lastPullTime map[string]time.Time, lastPullErr map[string]string) ActivityHandlerOption {
	return func(h *ActivityHandler) {
		h.pullMu = pullMu
		h.lastPullTime = lastPullTime
		h.lastPullErr = lastPullErr
	}
}

func WithActivityHub(hub *ws.ActivityHub) ActivityHandlerOption {
	return func(h *ActivityHandler) { h.activityHub = hub }
}

func WithActivityEnhancer(e *ai.ActivityEnhancer) ActivityHandlerOption {
	return func(h *ActivityHandler) { h.enhancer = e }
}

func NewActivityHandler(opts ...ActivityHandlerOption) *ActivityHandler {
	h := &ActivityHandler{
		recentAgentEvents: make([]ActivityEvent, 0),
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

func (ah *ActivityHandler) SetRepos(repos []string) {
	ah.repos = repos
}

func (ah *ActivityHandler) HandleActivity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	typesParam := r.URL.Query().Get("types")
	sinceParam := r.URL.Query().Get("since")
	limitParam := r.URL.Query().Get("limit")

	limit := 50
	if limitParam != "" {
		if l, err := parsePositiveInt(limitParam); err == nil && l > 0 && l <= 200 {
			limit = l
		}
	}

	var since time.Time
	if sinceParam != "" {
		if t, err := time.Parse(time.RFC3339, sinceParam); err == nil {
			since = t
		}
	}

	typeFilters := make(map[EventType]bool)
	if typesParam != "" {
		for _, t := range strings.Split(typesParam, ",") {
			t = strings.TrimSpace(t)
			switch EventType(t) {
			case EventTypeCommit, EventTypeAgent, EventTypePull:
				typeFilters[EventType(t)] = true
			}
		}
	} else {
		typeFilters[EventTypeCommit] = true
		typeFilters[EventTypeAgent] = true
		typeFilters[EventTypePull] = true
	}

	events := ah.gatherEvents(typeFilters, since, limit)

	if len(events) > limit {
		events = events[:limit]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"events": events})
}

func (ah *ActivityHandler) gatherEvents(typeFilters map[EventType]bool, since time.Time, limit int) []ActivityEvent {
	type result struct {
		events []ActivityEvent
	}

	var wg sync.WaitGroup
	results := make([]result, 0, 3)
	var mu sync.Mutex

	if typeFilters[EventTypeCommit] && ah.getCommits != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			events := ah.gatherCommitEvents(since, limit)
			mu.Lock()
			results = append(results, result{events: events})
			mu.Unlock()
		}()
	}

	if typeFilters[EventTypeAgent] && ah.readCrontab != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			events := ah.gatherAgentEvents(since, limit)
			mu.Lock()
			results = append(results, result{events: events})
			mu.Unlock()
		}()
	}

	if typeFilters[EventTypePull] && ah.pullMu != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			events := ah.gatherPullEvents(since, limit)
			mu.Lock()
			results = append(results, result{events: events})
			mu.Unlock()
		}()
	}

	wg.Wait()

	var all []ActivityEvent
	for _, r := range results {
		all = append(all, r.events...)
	}

	sort.Slice(all, func(i, j int) bool {
		return all[i].Timestamp.After(all[j].Timestamp)
	})

	return all
}

func (ah *ActivityHandler) gatherCommitEvents(since time.Time, limit int) []ActivityEvent {
	var events []ActivityEvent
	for _, repoPath := range ah.repos {
		commits, err := ah.getCommits(repoPath, limit)
		if err != nil {
			continue
		}
		for _, c := range commits {
			if !since.IsZero() && c.Timestamp.Before(since) {
				continue
			}
			hashLen := 8
			if len(c.Hash) < hashLen {
				hashLen = len(c.Hash)
			}
			meta, err := json.Marshal(CommitEventMetadata{
				Hash:   c.Hash,
				Author: c.Author,
			})
			if err != nil {
				continue
			}
			event := ActivityEvent{
				ID:        "commit-" + c.Hash[:hashLen],
				Type:      EventTypeCommit,
				Repo:      repoPath,
				Message:   c.Message,
				Timestamp: c.Timestamp,
				Metadata:  meta,
			}
			if ah.enhancer != nil {
				commitInfo := ai.CommitInfo{
					Hash:    c.Hash,
					Author:  c.Author,
					Message: c.Message,
					Body:    c.Body,
				}
				summary, flags, err := ah.enhancer.EnhanceEvent(repoPath, c.Hash, []ai.CommitInfo{commitInfo})
				if err == nil && summary != "" {
					event.Summary = summary
					event.Flags = flags
				}
			}
			events = append(events, event)
		}
	}
	return events
}

func (ah *ActivityHandler) gatherAgentEvents(since time.Time, limit int) []ActivityEvent {
	ah.eventMu.Lock()
	defer ah.eventMu.Unlock()

	var events []ActivityEvent
	for _, e := range ah.recentAgentEvents {
		if !since.IsZero() && e.Timestamp.Before(since) {
			continue
		}
		events = append(events, e)
		if len(events) >= limit {
			break
		}
	}
	return events
}

func (ah *ActivityHandler) gatherPullEvents(since time.Time, limit int) []ActivityEvent {
	ah.pullMu.RLock()
	defer ah.pullMu.RUnlock()

	var events []ActivityEvent
	for repo, t := range ah.lastPullTime {
		if !since.IsZero() && t.Before(since) {
			continue
		}
		meta := PullEventMetadata{Success: true}
		if err, ok := ah.lastPullErr[repo]; ok && err != "" {
			meta.Success = false
			meta.Error = err
		}
		metaJSON, err := json.Marshal(meta)
		if err != nil {
			continue
		}
		events = append(events, ActivityEvent{
			ID:        "pull-" + repo,
			Type:      EventTypePull,
			Repo:      repo,
			Message:   "Pull " + repo,
			Timestamp: t,
			Metadata:  metaJSON,
		})
		if len(events) >= limit {
			break
		}
	}
	return events
}

func (ah *ActivityHandler) RecordAgentEvent(event ActivityEvent) {
	ah.eventMu.Lock()
	ah.recentAgentEvents = append([]ActivityEvent{event}, ah.recentAgentEvents...)
	if len(ah.recentAgentEvents) > 100 {
		ah.recentAgentEvents = ah.recentAgentEvents[:100]
	}
	ah.eventMu.Unlock()
	if ah.activityHub != nil {
		go ah.activityHub.Broadcast(ws.ActivityEvent{
			ID:        event.ID,
			Type:      string(event.Type),
			Repo:      event.Repo,
			Message:   event.Message,
			Timestamp: event.Timestamp,
			Summary:   event.Summary,
			Flags:     event.Flags,
			Metadata:  event.Metadata,
		})
	}
}

func parsePositiveInt(s string) (int, error) {
	var n int
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, nil
		}
		n = n*10 + int(c-'0')
	}
	return n, nil
}
