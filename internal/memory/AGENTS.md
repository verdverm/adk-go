# ADK Internal Memory Components (`internal/memory/`)

This package provides the concrete, session-scoped implementation of the `agent.Memory` interface, bridging the high-level agent-facing API with the lower-level `memory.Service`.

## Memory Implementation (`memory.go`)

*   **`Memory` struct**: An unexported structure that embeds a public `memory.Service` and caches the identifying metadata: `AppName`, `UserID`, and `SessionID`.
*   **Method Implementations**: The methods (`AddSession`, `Search`) translate the simplified calls from an `agent.InvocationContext` into the structured request objects required by the underlying memory service.
    *   `AddSession` forwards the entire `session.Session` object to the service for ingestion.
    *   `Search` creates a `memory.SearchRequest` using the current `AppName` and `UserID`, scoped for long-term knowledge retrieval.