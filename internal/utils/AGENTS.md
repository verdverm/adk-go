# Package `internal/utils`

This package provides a collection of general-purpose, unexported helper functions used across the ADK for tasks ranging from LLM content manipulation to structured data validation.

## LLM Content and Struct Utilities (`utils.go`)

*   **Function Call ID Management**:
    *   `PopulateClientFunctionCallID`: Inserts an ADK-prefixed UUID into empty `FunctionCall.ID` fields within `genai.Content`. This allows the ADK runtime to accurately map calls to responses in the event stream.
    *   `RemoveClientFunctionCallID`: Clears these ADK-prefixed IDs from content before it is sent to the LLM.
*   **Content Extraction**: Provides simple helpers (`FunctionCalls`, `FunctionResponses`, `TextParts`) to safely extract specific part types from a `genai.Content` object.
*   **Instructions**: `AppendInstructions` utility to safely add or append system instructions to an `model.LLMRequest`.
*   **Safety**: `Must` helper function for panicking on non-nil errors in configuration steps.

## Schema Validation Utilities (`schema_utils.go`)

*   **`ValidateMapOnSchema`**: The core validation function that checks if a `map[string]any` (e.g., function arguments or results) conforms to a given `genai.Schema`, ensuring all required properties are present and their types match the schema definition.
