package agents

import (
	"sync"
	"testing"
	"time"
)

func TestAgentState(t *testing.T) {
	state := &AgentState{ExitCode: 1, LastError: "command not found"}
	if state.ExitCode != 1 {
		t.Errorf("expected ExitCode 1, got %d", state.ExitCode)
	}
	if state.LastError != "command not found" {
		t.Errorf("expected LastError 'command not found', got %q", state.LastError)
	}
}

func TestAgentStateMap_New(t *testing.T) {
	m := NewAgentStateMap()
	if m == nil {
		t.Fatal("expected non-nil AgentStateMap")
	}
	if m.states == nil {
		t.Error("expected states map to be initialized")
	}
}

func TestAgentStateMap_Get_NotFound(t *testing.T) {
	m := NewAgentStateMap()
	state := m.Get("nonexistent-id")
	if state != nil {
		t.Error("expected nil for nonexistent agent")
	}
}

func TestAgentStateMap_SetAndGet(t *testing.T) {
	m := NewAgentStateMap()
	agentID := "0 */4 * * *:/home/user/proj:gpt-4o"
	state := &AgentState{ExitCode: 1, LastError: "test error"}

	m.Set(agentID, state)

	got := m.Get(agentID)
	if got == nil {
		t.Fatal("expected state after Set")
	}
	if got.ExitCode != 1 {
		t.Errorf("expected ExitCode 1, got %d", got.ExitCode)
	}
	if got.LastError != "test error" {
		t.Errorf("expected LastError 'test error', got %q", got.LastError)
	}
}

func TestAgentStateMap_Clear(t *testing.T) {
	m := NewAgentStateMap()
	agentID := "0 */4 * * *:/home/user/proj:gpt-4o"
	state := &AgentState{ExitCode: 0}

	m.Set(agentID, state)
	m.Clear(agentID)

	got := m.Get(agentID)
	if got != nil {
		t.Error("expected nil after Clear")
	}
}

func TestAgentStateMap_Overwrite(t *testing.T) {
	m := NewAgentStateMap()
	agentID := "0 */4 * * *:/home/user/proj:gpt-4o"

	m.Set(agentID, &AgentState{ExitCode: 1, LastError: "first error"})
	m.Set(agentID, &AgentState{ExitCode: 2, LastError: "second error"})

	got := m.Get(agentID)
	if got.ExitCode != 2 {
		t.Errorf("expected ExitCode 2, got %d", got.ExitCode)
	}
	if got.LastError != "second error" {
		t.Errorf("expected LastError 'second error', got %q", got.LastError)
	}
}

func TestAgentStateMap_ConcurrentAccess(t *testing.T) {
	m := NewAgentStateMap()
	agentIDs := []string{
		"0 */4 * * *:/home/user/proj:gpt-4o",
		"0 */4 * * *:/home/user/proj:gpt-3.5",
		"0 */4 * * *:/home/user/proj:claude",
	}
	done := make(chan struct{})
	wg := sync.WaitGroup{}

	for _, agentID := range agentIDs {
		wg.Add(3)
		go func(id string) {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				m.Set(id, &AgentState{ExitCode: i})
			}
		}(agentID)
		go func(id string) {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				m.Get(id)
			}
		}(agentID)
		go func(id string) {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				m.Clear(id)
			}
		}(agentID)
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("concurrent access test timed out - possible deadlock or race condition")
	}
}
