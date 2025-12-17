# Package `internal/context`

This package provides the concrete, unexported implementations of the public context interfaces (`agent.InvocationContext` and `agent.CallbackContext`), managing the environment and state propagation during an agent's execution.

## Invocation Context (`invocation_context.go`)

The `InvocationContext` is the runtime backbone of an agent. It carries all the state required for a single execution run.

```go
type InvocationContext struct {
	context.Context

	params       InvocationContextParams
	invocationID string
}

type InvocationContextParams struct {
	Artifacts agent.Artifacts
	Memory    agent.Memory
	Session   session.Session

	Branch string
	Agent  agent.Agent

	UserContent   *genai.Content
	RunConfig     *agent.RunConfig
	EndInvocation bool
}
```

**Key Features:**
*   **Immutable Core**: Most fields (`Agent`, `Session`, `Artifacts`) are fixed for the duration of the invocation.
*   **Run Configuration**: Carries `RunConfig` which determines behaviors like streaming modes (`StreamingModeSSE`, `StreamingModeBidi`).
*   **Branching**: Tracks the `Branch` ID, which is essential for filtering conversation history in non-linear workflows.
*   **Purpose**: To provide all necessary dependencies and state (read-only and mutable) to the currently executing agent and its sub-components.

## Readonly Context (`readonly_context.go`)

*   **`ReadonlyContext` struct**: A wrapper around `InvocationContext` that implements `agent.ReadonlyContext`.
*   **Purpose**: Used in contexts where state mutation is disallowed, such as when an `tool.Toolset` is determining which tools to expose, ensuring access only to read-only information like `AppName`, `SessionID`, `UserID`, and `ReadonlyState`.

## Callback Context (`callback_context.go`)

The `callbackContext` is a specialized view provided to tools and callbacks.

*   **`callbackContext` struct**: Implements the `agent.CallbackContext` interface, used by callbacks and tools.
*   **State Delta Tracking**: It is initialized with a `stateDelta` map. When `State().Set()` is called, the change is recorded here rather than immediately in the session, allowing atomic updates via `session.EventActions`.
*   **Artifact Delta Tracking**: Wraps `agent.Artifacts` to record new artifact versions in `session.EventActions.ArtifactDelta`.
*   **State Access**: Its implementation of the `session.State` interface prioritizes reading from the local `StateDelta` map before falling back to the session's global state, ensuring that edits made earlier in the same turn are visible immediately.
