package ws

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

const (
	inMask    = syscall.IN_MODIFY | syscall.IN_CREATE | syscall.IN_DELETE | syscall.IN_MOVE
	inBufSize = 1024 * (syscall.SizeofInotifyEvent + 256)
)

type LogWatcher struct {
	hub       *Hub
	agentID   string
	logPath   string
	lastSize  int64
	stop      chan struct{}
	done      chan struct{}
	inotifyFd int
	watchDir  string
	watchName string
}

func NewLogWatcher(hub *Hub, agentID, logPath string) *LogWatcher {
	dir := filepath.Dir(logPath)
	name := filepath.Base(logPath)
	return &LogWatcher{
		hub:       hub,
		agentID:   agentID,
		logPath:   logPath,
		stop:      make(chan struct{}),
		done:      make(chan struct{}),
		inotifyFd: -1,
		watchDir:  dir,
		watchName: name,
	}
}

func (lw *LogWatcher) Start() {
	go lw.watch()
}

func (lw *LogWatcher) Stop() {
	close(lw.stop)
	if lw.inotifyFd != -1 {
		syscall.Close(lw.inotifyFd)
	}
	<-lw.done
}

func (lw *LogWatcher) watch() {
	defer close(lw.done)

	if err := lw.initInotify(); err != nil {
		lw.watchFallback()
		return
	}

	buf := make([]byte, inBufSize)
	for {
		select {
		case <-lw.stop:
			return
		default:
			n, err := syscall.Read(lw.inotifyFd, buf)
			if err != nil {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			lw.processEvents(buf[:n])
		}
	}
}

func (lw *LogWatcher) initInotify() error {
	fd, err := syscall.InotifyInit()
	if err != nil {
		return err
	}
	lw.inotifyFd = fd

	wd, err := syscall.InotifyAddWatch(fd, lw.watchDir, inMask)
	if err != nil {
		syscall.Close(fd)
		return err
	}
	_ = wd

	lw.check()
	return nil
}

func (lw *LogWatcher) processEvents(buf []byte) {
	offset := 0
	for offset+syscall.SizeofInotifyEvent <= len(buf) {
		header := (*syscall.InotifyEvent)(unsafe.Pointer(&buf[offset]))
		if header.Len > 0 {
			nameEnd := offset + syscall.SizeofInotifyEvent + int(header.Len)
			if nameEnd > len(buf) {
				break
			}
			name := string(buf[offset+syscall.SizeofInotifyEvent : nameEnd])
			if name == lw.watchName || (name == lw.watchName+".tmp") {
				lw.check()
			}
		}
		offset += syscall.SizeofInotifyEvent + int(header.Len)
	}
}

func (lw *LogWatcher) watchFallback() {
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
