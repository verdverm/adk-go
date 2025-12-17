# Package `internal/toolinternal`

This package provides the core internal logic, interfaces, and context implementations required for managing and executing tools during the LLM agent flow.

## Core Tool Interfaces and Types (`tool.go`)

```go
type FunctionTool interface {
	tool.Tool
	Declaration() *genai.FunctionDeclaration
	Run(ctx tool.Context, args any) (result map[string]any, err error)
}

type RequestProcessor interface {
	ProcessRequest(ctx tool.Context, req *model.LLMRequest) error
}
```

## Tool Context Implementation (`context.go`)

The `toolContext` implements `tool.Context`. It is responsible for tracking changes (state, artifacts) during tool execution.

```go
type toolContext struct {
	agent.CallbackContext
	invocationContext agent.InvocationContext
	functionCallID    string
	eventActions      *session.EventActions
	artifacts         *internalArtifacts
}
```

**Key Features:**
*   **Artifact Tracking**: `internalArtifacts` wraps `agent.Artifacts`. Its `Save` method automatically records a `session.EventActions.ArtifactDelta` with the new artifact version.
*   **Factory**: `NewToolContext` initializes the context with a unique `functionCallID`.

## Tool Packaging Utility (`toolutils/`)

Helper for packing tools into LLM requests.

**Key Functions:**
*   `PackTool(req *model.LLMRequest, tool Tool) error`: Consolidates function declarations. If `genai.Tool` with declarations exists, it appends; otherwise, it creates a new one. This ensures compatibility with APIs that expect batched declarations.
