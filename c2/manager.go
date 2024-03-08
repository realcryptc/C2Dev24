package c2

import (
	"fmt"
	"sync"
	"time"
)

type Agent struct {
	ID       string
	IP       string
	LastCall time.Time
	CmdQueue [][]string
}

type SafeAgentMap struct {
	mtx    sync.Mutex
	Agents map[string]*Agent // pointer to Agent Type ID: Agent Pointer
}

// AgentMap is implicitly started by being a global variable
var AgentMap SafeAgentMap = SafeAgentMap{Agents: make(map[string]*Agent)}

func (am *SafeAgentMap) Add(agent *Agent) {
	am.mtx.Lock()
	defer am.mtx.Unlock()

	if _, exists := am.Agents[agent.ID]; !exists {
		am.Agents[agent.ID] = agent
	}
	// Can add Logging here -- Ben says kinda extra though
}

// agent operations -- will manage commands that each agent is doing

func (am *SafeAgentMap) Get(agentID string) *Agent {
	// *Agent = pass by reference
	am.mtx.Lock()
	defer am.mtx.Unlock()
	if _, exists := am.Agents[agentID]; exists {
		return am.Agents[agentID]
	}
	return nil
}

func (am *SafeAgentMap) Enqueue(agentID string, cmd []string) error {
	agent := am.Get(agentID)
	if agent == nil {
		return fmt.Errorf("agent '%s' does not exist", agentID)
	}
	am.mtx.Lock()
	defer am.mtx.Unlock()

	agent.CmdQueue = append(agent.CmdQueue, cmd)
	return nil
}

func (am *SafeAgentMap) Dequeue(agentID string) ([]string, error) {
	agent := am.Get(agentID)
	if agent == nil {
		return nil, fmt.Errorf("agent '%s' does not exist", agentID)
	}

	am.mtx.Lock()
	defer am.mtx.Unlock()

	if len(agent.CmdQueue) < 1 {
		return nil, fmt.Errorf("agent '%s' has no queued commands", agentID)
	}
	// update queue to exclude first element
	cmd := agent.CmdQueue[0]
	agent.CmdQueue = agent.CmdQueue[1:]
	return cmd, nil
}
