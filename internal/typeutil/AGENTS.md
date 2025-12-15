# ADK Internal Type Utilities (`internal/typeutil/`)

This package provides utility functions for converting between Go types, with an emphasis on using JSON marshalling and unmarshalling to handle complex structural conversion and optional JSON schema validation.

## Type Conversion with Schema Validation (`convert.go`)

*   **`ConvertToWithJSONSchema[From, To any]` function**: A generic function designed to convert a value of type `From` to a value of type `To`.
*   **Conversion Mechanism**: The function uses `json.Marshal` and subsequent `json.Unmarshal` to perform the conversion. This method is preferred when direct type assertion or reflection is inadequate for handling complex struct tags or data formats.
*   **Validation**: It accepts an optional `jsonschema.Resolved` object. If provided, the function will validate the marshaled content (after unmarshalling into a generic map structure) against the schema to ensure data integrity before finally unmarshalling it into the target Go type.