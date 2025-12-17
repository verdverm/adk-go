# Package `internal/sessioninternal`

This package provides internal implementation details for the public `session.Service` and related state structures.

## Mutable Session Implementation (`mutablesession.go`)

*   **`MutableSession` struct**: An unexported wrapper around an underlying `session.Session` that implements the mutable methods required during an agent's run.
*   **Purpose**: This struct implements `session.Session` and its embedded `session.State` interface, delegating read operations to the underlying session and write operations (`State().Set()`) to the internal `MutableState` interface of the stored session. This allows agents to safely modify the session state within the runner.

## Mutable State Contract (`session.go`)

*   **`MutableState` interface**: A small, unexported interface defining only the `Set(string, any) error` method. This acts as the internal contract that must be implemented by any concrete state store to allow session state modification.

```go
type MutableState interface {
	Set(string, any) error
}
```
