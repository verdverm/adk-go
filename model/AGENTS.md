# Package `model`

This package defines the interfaces and data structures for interacting with Large Language Models (LLMs) in a model-agnostic way.

## Core Abstractions (`llm.go`)

### The LLM Interface

```go
type LLM interface {
	Name() string
	GenerateContent(ctx context.Context, req *LLMRequest, stream bool) iter.Seq2[*LLMResponse, error]
}
```

### Request and Response

The `LLMRequest` and `LLMResponse` structures wrap the underlying provider's types (currently heavily relying on `genai` types) to provide a standard exchange format.

```go
type LLMRequest struct {
	Model    string
	Contents []*genai.Content
	Config   *genai.GenerateContentConfig
	Tools    map[string]any // Internal tool representation
}

type LLMResponse struct {
	Content           *genai.Content
	UsageMetadata     *genai.GenerateContentResponseUsageMetadata
	// ... metadata ...
	Partial           bool // True if this is a chunk of a stream
	TurnComplete      bool // True if the model turn is finished
}
```

## Implementations

### Gemini (`gemini/`)
The primary implementation backed by the Google GenAI SDK.

*   **Constructor**: `gemini.NewModel(ctx, modelName, cfg)`
*   **Features**: Supports streaming, function calling, and standard content generation.
