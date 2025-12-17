# Package `agent`

This package is the heart of the ADK, defining the universal `Agent` interface and providing its primary implementations.

## Core Interfaces (`agent.go`)

### The Agent Interface

All agents, regardless of their internal logic (LLM, code, workflow), must implement this interface.

```go
type Agent interface {
	Name() string
	Description() string
	Run(InvocationContext) iter.Seq2[*session.Event, error]
	SubAgents() []Agent

	internal() *agent
}
```

### Context (`context.go`)

The `InvocationContext` provides the runtime environment for an agent's execution.

```go
type InvocationContext interface {
	context.Context
	Agent() Agent
	Artifacts() Artifacts
	Memory() Memory
	Session() session.Session
	InvocationID() string
	Branch() string
	UserContent() *genai.Content
	RunConfig() *RunConfig
	EndInvocation()
	Ended() bool
}
```

## Agent Implementations

### LLM Agent (`llmagent/`)
The most common agent type, driven by a Large Language Model (e.g., Gemini). It manages the conversational loop, tool execution, and prompt engineering.

*   **Constructor**: `llmagent.New(cfg llmagent.Config)`
*   **Features**: Supports `Instructions`, `Tools`, `FunctionCalling`, and `AgentTransfer`.

### Remote Agent (`remoteagent/`)
Enables the "Agent-to-Agent" (A2A) protocol, allowing an agent to communicate with another agent running in a different process or service.

*   **Constructor**: `remoteagent.NewA2A(cfg remoteagent.A2AConfig)`

### Workflow Agents (`workflowagents/`)
Control flow agents that manage the execution of sub-agents.

*   **Loop Agent**: Repeatedly runs sub-agents.
    *   `loopagent.New(cfg loopagent.Config)`
*   **Sequential Agent**: Runs sub-agents one after another.
*   **Parallel Agent**: Runs sub-agents concurrently.
