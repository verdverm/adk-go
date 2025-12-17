# Package `artifact`

This package defines the service interface for managing file artifacts associated with a session.

## Core Interfaces (`service.go`)

### Service

The `Service` interface provides storage operations for artifacts, including versioning support.

```go
type Service interface {
	Save(ctx context.Context, req *SaveRequest) (*SaveResponse, error)
	Load(ctx context.Context, req *LoadRequest) (*LoadResponse, error)
	Delete(ctx context.Context, req *DeleteRequest) error
	List(ctx context.Context, req *ListRequest) (*ListResponse, error)
	Versions(ctx context.Context, req *VersionsRequest) (*VersionsResponse, error)
}
```

### Data Models

Artifacts are identified by a tuple of `(AppName, UserID, SessionID, FileName)`.

```go
type SaveRequest struct {
	AppName, UserID, SessionID, FileName string
	Part *genai.Part
	Version int64 // Optional, for optimistic concurrency or explicit versioning
}
```
