package ws

import (
	"bufio"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type AgentExecutor struct {
	hub      *Hub
	watchers *WatcherManager
}

func NewAgentExecutor(hub *Hub, watchers *WatcherManager) *AgentExecutor {
	return &AgentExecutor{
		hub:      hub,
		watchers: watchers,
	}
}

func (ae *AgentExecutor) RunAgent(agentID, binaryPath string, args []string, workDir, logPath string) error {
	cmd := exec.Command(binaryPath, args...)
	cmd.Dir = workDir

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	go ae.streamOutput(stdout, agentID, "stdout")
	go ae.streamOutput(stderr, agentID, "stderr")

	err = cmd.Wait()

	if cmd.ProcessState != nil {
		ae.hub.Broadcast(LogEntry{
			AgentID:   agentID,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Message:   "Process exited with code " + strconv.Itoa(cmd.ProcessState.ExitCode()),
			Type:      "info",
		})
	}

	return err
}

func (ae *AgentExecutor) streamOutput(pipe io.Reader, agentID, outputType string) {
	reader := bufio.NewReader(pipe)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimRight(line, "\r\n")
		if line != "" {
			ae.hub.Broadcast(LogEntry{
				AgentID:   agentID,
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Message:   line,
				Type:      outputType,
			})
		}
	}
}
