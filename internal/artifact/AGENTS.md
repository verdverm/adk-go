# Package `internal/artifact`

This package provides the concrete, session-scoped implementation of the `agent.Artifacts` interface, bridging the high-level agent-facing API with the lower-level `artifact.Service`.

## Core Implementation (`artifacts.go`)

The primary component is the `Artifacts` struct, which serves as a secure proxy to the global artifact service.

### The `Artifacts` Struct

```go
type Artifacts struct {
	Service   artifact.Service
	AppName   string
	UserID    string
	SessionID string
}
```

### The Security Guarantee

The critical function of this package is to enforce **Session Isolation**.
1.  **Context Binding**: When an agent requests an artifact operation (e.g., `Save("image.png")`), it does not provide user or session identifiers.
2.  **Automatic Scoping**: This implementation automatically injects the `AppName`, `UserID`, and `SessionID` (stored in the struct) into the underlying `artifact.SaveRequest`.
3.  **Guarantee**: This guarantees that an agent executing within a specific session cannot accidentally read or write artifacts belonging to a different user or session.

### Method Signatures

The methods implement `agent.Artifacts` by transforming the call into a fully qualified service request.

```go
func (a *Artifacts) Save(ctx context.Context, name string, data *genai.Part) (*artifact.SaveResponse, error)
func (a *Artifacts) Load(ctx context.Context, name string) (*artifact.LoadResponse, error)
func (a *Artifacts) LoadVersion(ctx context.Context, name string, version int) (*artifact.LoadResponse, error)
func (a *Artifacts) List(ctx context.Context) (*artifact.ListResponse, error)
```
