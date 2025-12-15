# ADK Internal Telemetry Components (`internal/telemetry/`)

This package provides the concrete, unexported implementation of the ADK's tracing integration, relying on OpenTelemetry standards to capture agent and model events.

## Tracing Implementation (`telemetry.go`)

*   **Tracers and Processors**: Manages a local `sdktrace.TracerProvider` and integrates with the global OpenTelemetry tracer. `AddSpanProcessor` allows external components to register custom span processors for exporting trace data.
*   **Trace Context Management**: `StartTrace` returns a pair of spans (from the local and global tracers) to ensure all trace consumers receive the necessary context.
*   **Standardized Attributes**: Defines a comprehensive set of constant string attributes (prefixed `genAi.` and `gcp.vertex.agent.`) for model requests, tool calls, and responses, ensuring consistency across the ADK.
*   **`TraceLLMCall`**: Creates spans for calls to the Large Language Model, capturing details like model name, token counts (`UsageMetadata`), and request/response payloads (with filtering for large file parts).
*   **`TraceToolCall` / `TraceMergedToolCalls`**: Creates spans for individual and batched tool execution, capturing tool names, descriptions, and the tool's response payload.
*   **`safeSerialize`**: Utility to reliably convert Go objects to JSON strings for trace attributes, handling serialization errors gracefully.