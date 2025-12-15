# ADK Internal Package Structure (`internal/`) - Detailed Reference

The `internal/` directory contains packages intended for exclusive use by the Agent Development Kit (ADK) codebase. Per Go's compilation rules, code outside of the ADK repository cannot import these packages. They encapsulate essential, unexported implementation details, utility logic, and testing infrastructure.

---

## Agent Lifecycle and State Management

| Package | Description |
| :--- | :--- |
| `internal/agent/` | **Internal Agent State and Types.** This package holds the foundational state for agent abstractions. Key components include:<ul><li>`State`: The unexported base struct holding agent configuration details like `AgentType`.</li><li>`Type`: Defines various concrete agent types suchs as `TypeLLMAgent`, `TypeCustomAgent`, `TypeLoopAgent`, etc.</li><li>Sub-packages like `parentmap` (for agent tree traversal) and `runconfig` (for internal runtime configuration handling).</li></ul> |
| `internal/context/` | **Invocation Context Implementation.** This package provides the concrete, unexported implementation of the primary context structures used during an agent's execution.<ul><li>`InvocationContext`: The concrete structure implementing the public `agent.InvocationContext` interface. It manages session, artifact, and memory service handles, as well as invocation-specific details like `InvocationID`, `Branch`, and `RunConfig`.</li><li>`NewInvocationContext`: The factory function used to initialize a new, unique context instance for each top-level agent run.</li></ul> |

## LLM Agent Core Execution Flow

The `internal/llminternal/` package is the most complex, housing the entire business logic for `LLMAgent` execution.

| Package | Description |
| :--- | :--- |
| `internal/llminternal/` | **LLM Execution Flow and State.** This package defines the multi-turn conversational loop and the internal configurations for an LLM-based agent.<ul><li>`State`: Defines the internal configuration of an `LLMAgent`, including the `model.LLM`, lists of `tool.Tool` and `tool.Toolset`, instruction strings (`Instruction`, `GlobalInstruction`), and configuration for handling agent transfers (`DisallowTransferToParent`, `DisallowTransferToPeers`).</li><li>`Flow`: The central orchestration struct, holding callbacks (`BeforeModelCallbacks`, `AfterModelCallbacks`, etc.) and a pipeline of request/response processors.</li><li>`Flow.Run`: The main iterative function that defines the agent's core loop, continually calling `runOneStep` until a final response is achieved.</li><li>`runOneStep`: Executes a single turn of the loop, which involves: **Preprocessing** (running processors like `basicRequestProcessor`, `instructionsRequestProcessor`, `AgentTransferRequestProcessor`), **Calling the LLM** (`callLLM`), **Post-processing** (running processors like `nlPlanningResponseProcessor`), and **Handling Function Calls** (`handleFunctionCalls`).</li><li>`handleFunctionCalls`: Manages the execution of tools requested by the LLM, running tool callbacks (`BeforeToolCallbacks`, `AfterToolCallbacks`), executing the tool's `Run` method, and generating a `FunctionResponse` event for the next LLM turn.</li></ul>|

## Service Implementation Details and Utilities

| Package | Description |
| :--- | :--- |
| `internal/artifact/` | **Artifact Service Implementation.** Houses the concrete logic and potentially the data models used by the public `artifact.Service` interface. |
| `internal/memory/` | **Memory Service Implementation.** Contains the concrete logic for the public `memory.Service` interface, including ingestion and search logic for long-term agent knowledge. |
| `internal/sessioninternal/` | **Session Service Implementation.** Provides the internal data structures and methods for manipulating the session state and event history, supporting the public `session.Service`. |
| `internal/sessionutils/` | **Session Data Utilities.** Contains helper functions specifically designed for operations on session data, such as event list manipulation. |
| `internal/toolinternal/` | **Tool Processing Logic.** Contains internal interfaces (like `toolinternal.RequestProcessor`, `toolinternal.FunctionTool`) and logic used by the `Flow` to manage tool discovery, processing, and execution. |
| `internal/utils/` | **General-Purpose Unexported Utilities.** A collection of helper functions used across multiple ADK packages for tasks that do not belong to a specific domain (e.g., general data manipulation, content parsing). |
| `internal/version/` | **ADK Version Tracking.** Code responsible for tracking and exposing the ADK library version string. |

## Testing and Observability

| Package | Description |
| :--- | :--- |
| `internal/httprr/` | **HTTP Record and Replay.** Implements `http.RoundTripper` to enable deterministic testing of network-dependent code (like LLM calls). It allows recording live HTTP interactions to a trace file and replaying them later, using request/response scrubbing functions to remove non-deterministic data and secrets. |
| `internal/telemetry/` | **Internal Telemetry Implementation.** Provides the concrete, unexported structures and integration points for ADK's observability system, linking agent events and tool calls to tracing backends (e.g., OpenTelemetry). |
| `internal/testutil/` | **Testing Helpers.** Contains general utility functions and fixtures exclusively intended for use within the ADK's comprehensive internal test suites. |

*Note: Directories like `internal/cli/`, `internal/converters/`, and `internal/typeutil/` are also present, containing related unexported logic.*