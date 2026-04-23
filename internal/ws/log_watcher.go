package ws

import (
	"bufio"
	"os"
	"strings"
	"time"
)

type LogWatcher struct {
	hub       *Hub
	agentID   string
	logPath   string
	lastSize  int64
	stop      chan struct{}
	done      chan struct{}
}

func NewLogWatcher(hub *Hub, agentID, logPath string) *LogWatcher {
	return &LogWatcher{
		hub:     hub,
		agentID: agentID,
		logPath: logPath,
		stop:    make(chan struct{}),
		done:    make(chan struct{}),
	}
}

func (lw *LogWatcher) Start() {
	go lw.watch()
}

func (lw *LogWatcher) Stop() {
	close(lw.stop)
	<-lw.done
}

func (lw *LogWatcher) watch() {
	defer close(lw.done)

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-lw.stop:
			return
		case <-ticker.C:
			lw.check()
		}
	}
}

func (lw *LogWatcher) check() {
	info, err := os.Stat(lw.logPath)
	if err != nil {
		return
	}

	size := info.Size()
	if size <= lw.lastSize {
		if size < lw.lastSize {
			lw.lastSize = 0
		}
		return
	}

	f, err := os.Open(lw.logPath)
	if err != nil {
		return
	}
	defer f.Close()

	if _, err := f.Seek(lw.lastSize, 0); err != nil {
		return
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			entry := LogEntry{
				AgentID:   lw.agentID,
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Message:   line,
				Type:      inferLineType(line),
			}
			lw.hub.Broadcast(entry)
		}
	}
	lw.lastSize = size
}

func inferLineType(line string) string {
	lower := strings.ToLower(line)
	if strings.Contains(lower, "error") || strings.Contains(lower, "fatal") || strings.Contains(lower, "panic") {
		return "stderr"
	}
	if strings.Contains(lower, "warn") {
		return "warn"
	}
	return "stdout"
}