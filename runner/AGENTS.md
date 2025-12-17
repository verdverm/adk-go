# Package `runner`

This package provides a runtime for ADK agents.

## The Runner (`runner.go`)

The `Runner` manages the execution of the agent within a session, handling message processing, event generation, and interaction with various services like artifact storage, session management, and memory.

```go
type Runner struct {
	// ... unexported fields ...
}
```

**Key Functions:**
*   `New(cfg Config) (*Runner, error)`: Creates a new runner with the provided services and root agent.
*   `Run(ctx, userID, sessionID, msg, cfg)`: The main entry point.
    1.  Loads the session state.
    2.  Determines which agent in the hierarchy should handle the request (`findAgentToRun`).
    3.  Creates the `InvocationContext`.
    4.  Executes the agent's `Run` loop.
    5.  Persists the resulting events to the session service.

### Agent Resolution Strategy
The runner inspects the session history to find the last message. If it was a function call or a transfer, it might route execution to a sub-agent. Otherwise, it defaults to the root agent.
