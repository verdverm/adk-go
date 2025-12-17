# Package `internal/toolinternal`

This package provides the core internal logic, interfaces, and context implementations required for managing and executing tools during the LLM agent flow.

## Core Tool Interfaces and Types (`tool.go`)

*   `FunctionTool` **interface**: Extends the public `tool.Tool` interface by adding:
    *   **`Declaration() genai.FunctionDeclaration`**: Provides the tool's structured declaration for the LLM.
    *   `Run(ctx tool.Context, args any) (result map[string]any, err error)`: The concrete, internal method used by the LLM flow to execute the tool and return a structured result map.
*   `RequestProcessor` **interface**: Defines the internal contract for any processor that needs to modify a `model.LLMRequest` (e.g., adding contents, instructions, or dynamic tools) before it is sent to the LLM.

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

*   `toolContext` **struct**: The concrete implementation of the public `tool.Context` interface.
*   **Artifacts Change Tracking**: It embeds a wrapped `internalArtifacts` struct. This wrapper's `Save` method automatically records a `session.EventActions.ArtifactDelta` entry containing the new artifact `Version`, ensuring that artifact changes made by a tool are committed to the session history.
*   **Creation**: `NewToolContext` is the factory method, responsible for initializing the context, assigning a unique `FunctionCallID`, and ensuring it has an `EventActions` structure for recording state and artifact changes.

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

*   `PackTool` **function**: A utility function critical for preparing the `model.LLMRequest` for API compatibility. It iterates through a list of tools and consolidates all `genai.FunctionDeclaration` objects into a single `genai.Tool` object within `req.Config.Tools`. This is necessary because some LLM endpoints expect all function declarations to be batched together. It also prevents duplicate tool names.
