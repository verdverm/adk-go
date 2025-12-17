# Package `internal/agent`

This package encapsulates the unexported, foundational state and structural components necessary for defining and managing agents within the ADK's agent tree.

## Core State (`state.go`)

This file defines the base internal state that all public `agent.Agent` implementations embed and rely on.

*   `State` **struct**: An unexported struct that holds essential configuration fields, most notably the `AgentType`.
*   `Type` **constants**: Defines the concrete agent classification.
*   `Reveal` **function**: A utility function that provides access to the internal `*State`, allowing other internal ADK packages to inspect an agent's specific type and configuration.

```go
type State struct {
	AgentType Type
	Config    any
}

type Type string

const (
	TypeLLMAgent        Type = "LLMAgent"
	TypeLoopAgent       Type = "LoopAgent"
	TypeSequentialAgent Type = "SequentialAgent"
	TypeParallelAgent   Type = "ParallelAgent"
	TypeCustomAgent     Type = "CustomAgent"
)
```

**Key Functions:**
*   `Reveal(a Agent) *State`: Accesses the internal state of an agent.

## Agent Tree Structure (`parentmap/`)

This sub-package manages the hierarchy of agents, ensuring a valid agent tree structure.

```go
type Map map[string]agent.Agent
```

*   `Map` **type (in `map.go`)**: A map from an agent's name (`string`) to its parent (`agent.Agent`).
*   `New` **function**: Constructs the parent map by traversing the agent tree starting from the root, performing validation checks to ensure:
    *   Agents have at most one parent.
    *   All agent names are unique within the tree.
*   `RootAgent` **function**: Utility to traverse up the tree from any agent to find the root agent.
*   `ToContext` / `FromContext`: Functions for securely storing and retrieving the agent hierarchy map within a `context.Context`.

## Runtime Configuration (`runconfig/`)

This sub-package defines and manages invocation-specific runtime configurations that affect agent behavior.

```go
type RunConfig struct {
       StreamingMode StreamingMode
}

type StreamingMode string

const (
       StreamingModeNone StreamingMode = "none"
       StreamingModeSSE  StreamingMode = "sse"
       StreamingModeBidi StreamingMode = "bidi"
)
```

*   `RunConfig` **struct (in `run_config.go`)**: Holds configuration parameters for a single agent run.
*   `StreamingMode`: An enum-like type defining how content is streamed to the user, with constants:
    *   `StreamingModeNone`
    *   `StreamingModeSSE` (Server-Sent Events)
    *   `StreamingModeBidi` (Bidirectional)
*   `ToContext` / `FromContext`: Functions to serialize the `RunConfig` into and retrieve it from a `context.Context`, making it available throughout the agent's execution flow.
