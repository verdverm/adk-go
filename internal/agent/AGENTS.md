# ADK Internal Agent Components (`internal/agent/`)

The `internal/agent/` package encapsulates the unexported, foundational state and structural components necessary for defining and managing agents within the ADK's agent tree.

## Core State (`state.go`)

This file defines the base internal state that all public `agent.Agent` implementations embed and rely on:

*   **`State` struct**: An unexported struct that holds essential configuration fields, most notably the `AgentType`.
*   **`Type` constants**: Defines the concrete agent classification, such as:
    *   `TypeLLMAgent`
    *   `TypeLoopAgent`
    *   `TypeSequentialAgent`
    *   `TypeParallelAgent`
    *   `TypeCustomAgent`
*   **`Reveal` function**: A utility function that provides access to the internal `*State`, allowing other internal ADK packages to inspect an agent's specific type and configuration.

## Agent Tree Structure (`parentmap/`)

This sub-package manages the hierarchy of agents, ensuring a valid agent tree structure:

*   **`Map` type (in `map.go`)**: A map from an agent's name (`string`) to its parent `agent.Agent`.
*   **`New` function**: Constructs the parent map by traversing the agent tree starting from the root, performing validation checks to ensure:
    *   Agents have at most one parent.
    *   All agent names are unique within the tree.
*   **`RootAgent` function**: Utility to traverse up the tree from any agent to find the root agent.
*   **`ToContext` / `FromContext`**: Functions for securely storing and retrieving the agent hierarchy map within a `context.Context`.

## Runtime Configuration (`runconfig/`)

This sub-package defines and manages invocation-specific runtime configurations that affect agent behavior:

*   **`RunConfig` struct (in `run_config.go`)**: Holds configuration parameters for a single agent run.
*   **`StreamingMode`**: An enum-like type defining how content is streamed to the user, with constants:
    *   `StreamingModeNone`
    *   `StreamingModeSSE` (Server-Sent Events)
    *   `StreamingModeBidi` (Bidirectional)
*   **`ToContext` / `FromContext`**: Functions to serialize the `RunConfig` into and retrieve it from a `context.Context`, making it available throughout the agent's execution flow.