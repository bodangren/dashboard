package ws

import (
	"sync"
)

type WatcherManager struct {
	watchers map[string]*LogWatcher
	mu       sync.Mutex
	hub      *Hub
}

func NewWatcherManager(hub *Hub) *WatcherManager {
	return &WatcherManager{
		watchers: make(map[string]*LogWatcher),
		hub:      hub,
	}
}

func (wm *WatcherManager) StartWatching(agentID, logPath string) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	if _, exists := wm.watchers[agentID]; exists {
		return
	}

	watcher := NewLogWatcher(wm.hub, agentID, logPath)
	wm.watchers[agentID] = watcher
	watcher.Start()
}

func (wm *WatcherManager) StopWatching(agentID string) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	if watcher, exists := wm.watchers[agentID]; exists {
		watcher.Stop()
		delete(wm.watchers, agentID)
	}
}

func (wm *WatcherManager) StopAll() {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	for agentID, watcher := range wm.watchers {
		watcher.Stop()
		delete(wm.watchers, agentID)
	}
}
