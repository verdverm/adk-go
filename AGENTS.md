# ADK Go Project Structure

The Agent Development Kit (ADK) for Go is structured into several top-level packages, reflecting a modular and idiomatic Go project layout.

## Core Agent Components

These packages contain the essential logic for defining and running agents.

| Package | Description |
| :--- | :--- |
| `agent/` | Defines the `Agent` interface, core abstractions (`InvocationContext`), and callback types (`BeforeAgentCallback`, `AfterAgentCallback`) for defining agent behavior. |
| `model/` | Defines the `LLM` interface and structures (`LLMRequest`, `LLMResponse`) for model-agnostic interaction with Large Language Models (LLMs), such as Gemini. |
| `tool/` | Defines the `Tool` interface and `Context` for creating, managing, and calling external functions that agents can utilize. |
| `runner/` | Manages the execution environment and lifecycle of an agent within a session, handling message processing and service orchestration. |
| `session/` | Defines the conversation structure, including the `Session` interface, turn-based `Event` history, and `State` for session-scoped key-value storage. |

## Auxiliary Components

These packages support the core agent functionality, managing state, data, and serving.

| Package | Description |
| :--- | :--- |
| `memory/` | Defines the `Service` interface for long-term agent knowledge, enabling storage and retrieval of session content across different user sessions. |
| `artifact/` | Defines the `Service` interface for managing file artifacts (`Save`, `Load`, `Delete`, `List`) associated with a session, including versioning support. |
| `server/` | Hosts protocol implementations (e.g., ADK REST, ADK A2A) to expose and serve ADK agents as services. |
| `telemetry/` | Provides mechanisms, such as `RegisterSpanProcessor`, for observability, tracing, and metrics for agent events. |
| `util/` | A collection of general-purpose utility functions and helpers used throughout the codebase. |
| `internal/` | Packages intended for internal use only by the ADK. Go's compiler prevents external packages from importing these. |

## Examples and Entry Points

| Package | Description |
| :--- | :--- |
| `cmd/` | Contains the main entry points (executables) for the project, often used for example applications or CLI tools. |
| `examples/` | Standalone example code demonstrating various features and usage patterns of the ADK. |

## Internal Structure Summary

The `internal/` directory is composed of unexported packages that provide the concrete, runtime implementation details for the public interfaces.

| Package | Description |
| :--- | :--- |
| `internal/agent/` | Unexported foundational state and agent tree structures (`State`, `Type`, `parentmap`, `runconfig`). |
| `internal/context/` | Concrete implementations of `InvocationContext` and `CallbackContext` for agent runtime environment and state management. |
| `internal/llminternal/` | Core business logic for LLM agents, defining the multi-turn execution flow (`Flow`) and request/response processors. |
| `internal/memory/` | Session-scoped implementation of the public `memory.Service`. |
| `internal/artifact/` | Session-scoped implementation of the public `artifact.Service` with version tracking. |
| `internal/sessioninternal/` | Internal structures enabling mutable session state and event history manipulation. |
| `internal/toolinternal/` | Internal interfaces (`FunctionTool`, `RequestProcessor`) and logic for tool discovery and execution. |
| **...and others...** | Includes `cli/`, `converters/`, `utils/`, `telemetry/`, `testutil/`, and `version/` helpers. |

### LLM Agent Execution Flow (Condensed Summary)

The agent's iterative reasoning and action loop is defined in `internal/llminternal/base_flow.go`.

| Phase/Component | Role and Related Components |
| :--- | :--- |
| **`Flow` Struct** | Central engine defined in `base_flow.go`, holding extensible pipelines: `RequestProcessors`, `ResponseProcessors`, and execution Callbacks (`BeforeTool`, `AfterModel`, etc.). |
| **`Flow.Run` Loop** | The core iterative entry point. It calls `runOneStep` repeatedly until `lastEvent.IsFinalResponse()` is true, enabling multi-turn reasoning and tool-use cycles. |
| **Preprocessing** | Prepares `model.LLMRequest` (see `basic_processor.go`, `instruction_processor.go`, `contents_processor.go`). This phase injects session history, substitutes instructions with state variables, and adds `FunctionDeclaration` schemas from tools. |
| **LLM Call & Post-Processing** | `f.callLLM` executes the model. On response, `f.postprocess` runs `ResponseProcessors` to handle model outputs, and the event is finalized (assigning client IDs to function calls). |
| **Tool & Function Handling**| `f.handleFunctionCalls` executes requested tools via `toolinternal.FunctionTool`. The tool output is formatted as a new `FunctionResponse` event (Role: "user") to be sent back to the LLM in the next loop iteration. |
| **Agent Transfer** | Checks the yielded event for `ev.Actions.TransferToAgent` (via `AgentTransferRequestProcessor`). If present, control is immediately handed off to the target agent's `Run` method. |
