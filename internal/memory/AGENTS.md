# Package `internal/memory`

This package provides the concrete, session-scoped implementation of the `agent.Memory` interface, bridging the high-level agent-facing API with the lower-level `memory.Service`.

## Core Implementation (`memory.go`)

Similar to the `artifact` package, this package enforces session isolation for memory operations.

*   `Memory` **struct**: An unexported structure that embeds a public `memory.Service` and caches the identifying metadata: `AppName`, `UserID`, and `SessionID`.
*   **Method Implementations**: The methods (`AddSession`, `Search`) translate the simplified calls from an `agent.InvocationContext` into the structured request objects required by the underlying memory service.
    *   `AddSession` forwards the entire `session.Session` object to the service for ingestion.
    *   `Search` creates a `memory.SearchRequest` using the current `AppName` and `UserID`, scoped for long-term knowledge retrieval.

### The `Memory` Struct

```go
type Memory struct {
	Service   memory.Service
	SessionID string
	UserID    string
	AppName   string
}
```

### Scoping & Isolation

The `Memory` struct automatically injects the `AppName` and `UserID` into all memory search and storage requests. This ensures that agents only access long-term memory relevant to the current user and application context.

### Method Signatures

```go
// AddSession ingests the current session into the memory store.
func (a *Memory) AddSession(ctx context.Context, session session.Session) error

// Search retrieves relevant memory entries based on a query.
func (a *Memory) Search(ctx context.Context, query string) (*memory.SearchResponse, error)
```
