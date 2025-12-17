# ADK Go Project Structure

The Agent Development Kit (ADK) for Go is a framework for building stateful, interactive AI agents. It is structured to separate the public API (interfaces, models) from the complex internal execution logic.

## System Overview

The ADK operates on a **Session-based Event Loop** model:

1.  **Session**: Represents a conversation with a user. It holds the persistent `State` (key-value) and `Events` (history).
2.  **Runner**: The runtime that loads a session, determines which `Agent` should handle the next turn, and executes it.
3.  **Agent**: A logical unit (e.g., `LLMAgent`, `LoopAgent`) that receives the current context and produces `Events`.
4.  **Events**: Immutable records of interactions (User messages, Model responses, Tool calls).

## Package Directory

### Core Agent Components

These packages contain the essential logic for defining and running agents.

| Package | Description |
| :--- | :--- |
| [`agent/`](agent/AGENTS.md) | **The Agent Interface.** Defines `Agent`, `InvocationContext`, and the `llmagent`, `workflowagents` implementations. |
| [`model/`](model/AGENTS.md) | **LLM Abstraction.** Defines `LLM`, `LLMRequest`, `LLMResponse`. Adapters for Gemini/Vertex AI. |
| [`tool/`](tool/AGENTS.md) | **Tool Use.** Defines `Tool`, `Context`, and standard tools (`geminitool`, `functiontool`). |
| [`runner/`](runner/AGENTS.md) | **Execution Runtime.** Orchestrates the Agent-Session loop. |
| [`session/`](session/AGENTS.md) | **State Management.** Defines `Session`, `Event`, `State`. |

### Auxiliary Components

These packages support the core agent functionality.

| Package | Description |
| :--- | :--- |
| [`memory/`](memory/AGENTS.md) | **Long-term Memory.** RAG (Retrieval Augmented Generation) service interface. |
| [`artifact/`](artifact/AGENTS.md) | **File Storage.** Managed file storage for sessions. |
| [`server/`](server/AGENTS.md) | **Serving.** REST API and Agent-to-Agent (A2A) protocols. |
| [`telemetry/`](telemetry/AGENTS.md) | **Observability.** OpenTelemetry integration. |
| [`util/`](util/AGENTS.md) | **Utilities.** Helpers for instructions and data. |

## Examples and Entry Points

| Package | Description |
| :--- | :--- |
| `cmd/` | Contains the main entry points (executables) for the project, often used for example applications or CLI tools. |
| `examples/` | Standalone example code demonstrating various features and usage patterns of the ADK. |

## Internal Structure Summary

The `internal/` directory is composed of unexported packages that provide the concrete, runtime implementation details for the public interfaces.

[`internal/`](internal/AGENTS.md) **Private Implementation.** Contains the heavy lifting for LLM execution (`llminternal`), context management, and service adapters. **DO NOT IMPORT.**

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
