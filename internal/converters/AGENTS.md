# ADK Internal Converters Components (`internal/converters/`)

This package provides utility functions for converting internal data structures, primarily leveraging JSON marshalling to handle types where reflective libraries may be incompatible or incomplete.

## Conversion Utilities (`map_structure.go`)

*   **`ToMapStructure` function**: Converts an arbitrary Go value (`any`) into a standard `map[string]any`.
*   **Purpose**: It achieves the conversion by first marshalling the input to JSON bytes and then unmarshalling those bytes into the target map structure. This pattern is necessary to reliably convert types (like those from the `genai` package) which may not have the necessary struct field tags for libraries like `mapstructure`.