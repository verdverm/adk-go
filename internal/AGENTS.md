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
| [`internal/agent`](agent/AGENTS.md) | **Internal Agent State.** Holds the `State` struct embedded by all agents, the `Type` definitions, and the `parentmap` for tree traversal. |
| [`internal/context`](context/AGENTS.md) | **Context Implementation.** Concrete `InvocationContext` and `CallbackContext` that manage the lifecycle of a single run, including branching and configuration. |

### Execution Engines

| Package | Description |
| :--- | :--- |
| [`internal/llminternal`](llminternal/AGENTS.md) | **LLM Agent Engine.** The heart of the `LLMAgent`. Implements the `Flow` loop: Preprocessing -> LLM Call -> Tool Execution -> Postprocessing. |
| [`internal/toolinternal`](toolinternal/AGENTS.md) | **Tool Runtime.** Manages tool discovery, argument parsing, and the `toolContext` used during tool execution. |

### Service Implementations

| Package | Description |
| :--- | :--- |
| [`internal/sessioninternal`](sessioninternal/AGENTS.md) | **Mutable Session.** Wraps the public `Session` to allow internal state mutation (which is otherwise read-only to users). |
| [`internal/sessionutils`](sessionutils/AGENTS.md) | **Session Helpers.** Utilities for managing state deltas (`app:`, `user:`) and event filtering. |
| [`internal/artifact`](artifact/AGENTS.md) | **Artifact Adapter.** Implements `agent.Artifacts` with automatic session scoping. |
| [`internal/memory`](memory/AGENTS.md) | **Memory Adapter.** Implements `agent.Memory` with automatic user/app scoping. |

### Utilities and Infrastructure

| Package | Description |
| :--- | :--- |
| [`internal/telemetry`](telemetry/AGENTS.md) | **Tracing.** OpenTelemetry instrumentation for the entire ADK. |
| [`internal/httprr`](httprr/AGENTS.md) | **Testing Transport.** Record/Replay mechanism for deterministic LLM testing. |
| [`internal/utils`](utils/AGENTS.md) | **Helpers.** General utilities for `genai.Content` manipulation and function IDs. |
| [`internal/typeutil`](typeutil/AGENTS.md) | **Type Safety.** JSON schema validation and safe type conversion. |
| [`internal/converters`](converters/AGENTS.md) | **Data Mapping.** Conversion between `genai` types and internal maps. |
| [`internal/cli`](cli/AGENTS.md) | **CLI Tools.** Helpers for building the `adk` command line interface. |
| [`internal/testutil`](testutil/AGENTS.md) | **Test Fixtures.** Shared test runners and mock models. |
| [`internal/version`](version/AGENTS.md) | **Versioning.** Single source of truth for the ADK library version. |
