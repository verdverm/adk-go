# ADK Internal Context Components (`internal/context/`)

This package provides the concrete, unexported implementations of the public context interfaces (`agent.InvocationContext` and `agent.CallbackContext`), managing the environment and state propagation during an agent's execution.

## Invocation Context (`invocation_context.go`)

*   **`InvocationContext` struct**: The core context object that is created at the start of an agent run. It embeds a standard `context.Context` and holds all run-specific parameters: `Agent`, `Session`, `Artifacts`, `Memory`, `UserContent`, `RunConfig`, and a unique `InvocationID`.
*   **Purpose**: To provide all necessary dependencies and state (read-only and mutable) to the currently executing agent and its sub-components. It also tracks if the invocation has been manually ended (`EndInvocation`).

## Readonly Context (`readonly_context.go`)

*   **`ReadonlyContext` struct**: A wrapper around `InvocationContext` that implements `agent.ReadonlyContext`.
*   **Purpose**: Used in contexts where state mutation is disallowed, such as when an `tool.Toolset` is determining which tools to expose, ensuring access only to read-only information like `AppName`, `SessionID`, `UserID`, and `ReadonlyState`.

## Callback Context (`callback_context.go`)

*   **`callbackContext` struct**: Implements the `agent.CallbackContext` interface, primarily used by `BeforeAgentCallback`, `AfterAgentCallback`, and `tool.Tool` implementations.
*   **State Delta Tracking**: It is initialized with a `stateDelta` map. When a callback or tool calls its `State().Set()` method, the change is recorded in this map, allowing the resulting session event to carry the exact state changes (`session.EventActions.StateDelta`).
*   **Artifact Delta Tracking**: It wraps the `agent.Artifacts` interface in `internalArtifacts` to ensure that any call to `Save` also records the new artifact version in `session.EventActions.ArtifactDelta`.
*   **State Access**: Its implementation of the `session.State` interface prioritizes reading from the local `StateDelta` map before falling back to the session's global state, ensuring that edits made earlier in the same turn are visible immediately.