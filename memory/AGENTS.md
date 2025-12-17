# Package `memory`

This package defines the entities to interact with agent memory (long-term knowledge).

## Core Interfaces (`service.go`)

### Service

The `Service` interface allows agents to ingest session history and search for relevant information from past interactions.

```go
type Service interface {
	AddSession(ctx context.Context, s session.Session) error
	Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error)
}
```

### Data Models

```go
type SearchRequest struct {
	Query   string
	UserID  string
	AppName string
}

type Entry struct {
	Content *genai.Content
	Author string
	Timestamp time.Time
}
```
