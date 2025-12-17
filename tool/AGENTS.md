# Package `tool`

This package defines the interfaces for tools that can be called by an agent.

## Core Interfaces (`tool.go`)

### The Tool Interface

```go
type Tool interface {
	Name() string
	Description() string
	IsLongRunning() bool
}
```

### Context

The `Context` is passed to the tool execution, providing access to the current state and actions.

```go
type Context interface {
	agent.CallbackContext
	FunctionCallID() string
	Actions() *session.EventActions
	SearchMemory(context.Context, string) (*memory.SearchResponse, error)
}
```

### Toolset

A `Toolset` is a collection of tools that can be dynamically exposed based on the current context.

```go
type Toolset interface {
	Name() string
	Tools(ctx agent.ReadonlyContext) ([]Tool, error)
}
```

## Built-in Tools

*   **`geminitool`**: Wraps GenAI tools (like Google Search) or standard `genai.Tool` structs.
*   **`agenttool`**: Allows an agent to call another agent as a function.
*   **`functiontool`**: Wraps a Go function as a tool.
*   **`loadartifactstool`**: Allows the LLM to read content from stored artifacts.
*   **`exitlooptool`**: A control flow tool for breaking out of loop agents.
*   **`mcptoolset`**: Supports the Model Context Protocol (MCP) for external tool integration.
