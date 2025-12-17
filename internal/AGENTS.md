# ADK Internal Package Structure (`internal/`)

The `internal/` directory contains packages intended for exclusive use by the Agent Development Kit (ADK) codebase. Per Go's compilation rules, code outside of the ADK repository cannot import these packages. They encapsulate essential, unexported implementation details, utility logic, and testing infrastructure.

## Internal Architecture

The internal packages are organized into three layers:

1.  **State & Context**: The foundational data structures (`internal/agent`, `internal/context`) that carry state through the system.
2.  **Execution Engines**: The complex logic that drives agents (`internal/llminternal`) and tools (`internal/toolinternal`).
3.  **Service Adapters**: Concrete implementations of public interfaces (`internal/artifact`, `internal/memory`, `internal/sessioninternal`) that handle I/O and storage.

## Package Reference

### Agent Lifecycle and State Management

| Package | Description |
| :--- | :--- |
| [`internal/agent`](agent/AGENTS.md) | **Internal Agent State and Types.** This package holds the foundational state for agent abstractions. Key components include:<ul><li>`State`: The unexported base struct holding agent configuration details like `AgentType`.</li><li>`Type`: Defines various concrete agent types suchs as `TypeLLMAgent`, `TypeCustomAgent`, `TypeLoopAgent`, etc.</li><li>Sub-packages like `parentmap` (for agent tree traversal) and `runconfig` (for internal runtime configuration handling).</li></ul> |
| [`internal/context`](context/AGENTS.md) | **Invocation Context Implementation.** This package provides the concrete, unexported implementation of the primary context structures used during an agent's execution.<ul><li>`InvocationContext`: The concrete structure implementing the public `agent.InvocationContext` interface. It manages session, artifact, and memory service handles, as well as invocation-specific details like `InvocationID`, `Branch`, and `RunConfig`.</li><li>`NewInvocationContext`: The factory function used to initialize a new, unique context instance for each top-level agent run.</li></ul> |

### Execution Engines

| Package | Description |
| :--- | :--- |
| [`internal/llminternal`](llminternal/AGENTS.md) | **LLM Execution Flow and State.** This package defines the multi-turn conversational loop and the internal configurations for an LLM-based agent.<ul><li>`State`: Defines the internal configuration of an `LLMAgent`, including the `model.LLM`, lists of `tool.Tool` and `tool.Toolset`, instruction strings (`Instruction`, `GlobalInstruction`), and configuration for handling agent transfers (`DisallowTransferToParent`, `DisallowTransferToPeers`).</li><li>`Flow`: The central orchestration struct, holding callbacks (`BeforeModelCallbacks`, `AfterModelCallbacks`, etc.) and a pipeline of request/response processors.</li><li>`Flow.Run`: The main iterative function that defines the agent's core loop, continually calling `runOneStep` until a final response is achieved.</li><li>`runOneStep`: Executes a single turn of the loop, which involves: **Preprocessing** (running processors like `basicRequestProcessor`, `instructionsRequestProcessor`, `AgentTransferRequestProcessor`), **Calling the LLM** (`callLLM`), **Post-processing** (running processors like `nlPlanningResponseProcessor`), and **Handling Function Calls** (`handleFunctionCalls`).</li><li>`handleFunctionCalls`: Manages the execution of tools requested by the LLM, running tool callbacks (`BeforeToolCallbacks`, `AfterToolCallbacks`), executing the tool's `Run` method, and generating a `FunctionResponse` event for the next LLM turn.</li></ul> |
| [`internal/toolinternal`](toolinternal/AGENTS.md) | **Tool Processing Logic.** Contains internal interfaces (like `toolinternal.RequestProcessor`, `toolinternal.FunctionTool`) and logic used by the `Flow` to manage tool discovery, processing, and execution. |

### Service Implementations

| Package | Description |
| :--- | :--- |
| [`internal/sessioninternal`](sessioninternal/AGENTS.md) | **Session Service Implementation.** Provides the internal data structures and methods for manipulating the session state and event history, supporting the public `session.Service`. |
| [`internal/sessionutils`](sessionutils/AGENTS.md) | **Session Data Utilities.** Contains helper functions specifically designed for operations on session data, such as event list manipulation. |
| [`internal/artifact`](artifact/AGENTS.md) | **Artifact Service Implementation.** Houses the concrete logic and potentially the data models used by the public `artifact.Service` interface. |
| [`internal/memory`](memory/AGENTS.md) | **Memory Service Implementation.** Contains the concrete logic for the public `memory.Service` interface, including ingestion and search logic for long-term agent knowledge. |

### Utilities and Infrastructure

| Package | Description |
| :--- | :--- |
| [`internal/telemetry`](telemetry/AGENTS.md) | **Internal Telemetry Implementation.** Provides the concrete, unexported structures and integration points for ADK's observability system, linking agent events and tool calls to tracing backends (e.g., OpenTelemetry). |
| [`internal/httprr`](httprr/AGENTS.md) | **HTTP Record and Replay.** Implements `http.RoundTripper` to enable deterministic testing of network-dependent code (like LLM calls). It allows recording live HTTP interactions to a trace file and replaying them later, using request/response scrubbing functions to remove non-deterministic data and secrets. |
| [`internal/utils`](utils/AGENTS.md) | **General-Purpose Unexported Utilities.** A collection of helper functions used across multiple ADK packages for tasks that do not belong to a specific domain (e.g., general data manipulation, content parsing). |
| [`internal/typeutil`](typeutil/AGENTS.md) | **Type Safety.** JSON schema validation and safe type conversion. |
| [`internal/converters`](converters/AGENTS.md) | **Data Mapping.** Conversion between `genai` types and internal maps. |
| [`internal/cli`](cli/AGENTS.md) | **CLI Tools.** Helpers for building the `adk` command line interface. |
| [`internal/testutil`](testutil/AGENTS.md) | **Testing Helpers.** Contains general utility functions and fixtures exclusively intended for use within the ADK's comprehensive internal test suites. |
| [`internal/version`](version/AGENTS.md) | **ADK Version Tracking.** Code responsible for tracking and exposing the ADK library version string. |
