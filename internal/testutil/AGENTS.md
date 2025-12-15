# ADK Internal Test Utilities (`internal/testutil/`)

This package provides comprehensive fixtures and utilities necessary for deterministic and reliable testing of ADK agents, models, and execution components.

## Agent Testing Infrastructure (`test_agent_runner.go`)

*   **`TestAgentRunner`**: A high-level helper struct designed to simplify end-to-end testing of `agent.Agent` implementations. It manages the creation of sessions, wraps the public `runner.Runner`, and provides entry points (`Run`, `RunContent`, `RunContentWithConfig`) to execute the agent.
*   **`MockModel`**: A core test fixture that implements the `model.LLM` interface.
    *   It allows tests to pre-program responses (`Responses`) and captures all incoming requests (`Requests`) for assertion, eliminating external network dependencies.
    *   It supports both blocking (`Generate`) and streaming (`GenerateStream`) LLM modes.
*   **Collectors**: Helper functions (`CollectEvents`, `CollectParts`, `CollectTextParts`) for cleanly consuming and aggregating the `iter.Seq2` stream of results returned by the runner.

## LLM Testing Utilities (`genai.go`)

*   **`NewGeminiTransport`**: A utility function that integrates the LLM client configuration with the `internal/httprr` record/replay mechanism. This ensures that tests involving LLM calls can be recorded once and replayed indefinitely for determinism.
*   **`scrubGeminiRequest`**: An `httprr` scrubbing function specifically for Gemini/LLM API requests. It removes variable headers (`X-Goog-Api-Key`, `User-Agent`) and canonicalizes the JSON request body (using `json.Compact`) to prevent testing instability due to non-deterministic serialization or sensitive data logging.