# Package `internal/version`

This package is solely responsible for exposing the current version string of the ADK Go library.

## Version Constant (`version.go`)

*   `Version` **constant**: A string constant (e.g., `"0.3.0"`) that holds the current version of the ADK.
*   **Purpose**: This constant is imported by other internal packages (like `internal/llminternal` and `internal/telemetry`) and is included in metadata sent to the LLM (e.g., in the User-Agent header or request fields) for identification and telemetry.
