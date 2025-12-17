# Package `session`

This package defines the conversation structure, including the `Session` interface, turn-based `Event` history, and `State` for session-scoped key-value storage.

## Core Interfaces (`session.go`)

### Session and State

```go
type Session interface {
	ID() string
	AppName() string
	UserID() string
	State() State
	Events() Events
	LastUpdateTime() time.Time
}

type State interface {
	Get(string) (any, error)
	Set(string, any) error
	All() iter.Seq2[string, any]
}
```

### Event Model

An `Event` represents a single interaction or step in the conversation.

```go
type Event struct {
	model.LLMResponse

	ID        string
	Timestamp time.Time

	InvocationID string
	Branch       string
	Author       string

	Actions EventActions
	LongRunningToolIDs []string
}
```

## The Session Service (`service.go`)

The `Service` interface manages the lifecycle of sessions.

```go
type Service interface {
	Create(context.Context, *CreateRequest) (*CreateResponse, error)
	Get(context.Context, *GetRequest) (*GetResponse, error)
	List(context.Context, *ListRequest) (*ListResponse, error)
	Delete(context.Context, *DeleteRequest) error
	AppendEvent(context.Context, Session, *Event) error
}
```

## Implementations

*   **In-Memory**: A simple map-based implementation for testing and transient sessions.
*   **Database (`database/`)**: A GORM-based implementation for persistent storage.
