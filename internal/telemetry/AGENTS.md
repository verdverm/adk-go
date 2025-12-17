# Package `internal/telemetry`

This package provides the concrete, unexported implementation of the ADK's tracing integration, relying on OpenTelemetry standards to capture agent and model events.

## Tracing Implementation (`telemetry.go`)

*   **Tracers**: Manages a local `sdktrace.TracerProvider` and integrates with the global OpenTelemetry tracer.
*   **Trace Context**: `StartTrace` returns spans from both local and global tracers.
*   **Attributes**: Defines standardized attributes (prefixed `genAi.` and `gcp.vertex.agent.`) for model requests, tool calls, and responses.
*   **Functions**:
    *   `TraceLLMCall`: Captures model name, token counts, and payloads.
    *   `TraceToolCall`: Captures tool name, args, and response.
