# Package `internal/httprr`

This package provides a mechanism for recording and replaying HTTP interactions, primarily used to ensure deterministic and fast unit testing of components that interact with external services, such as the LLM.

## RecordReplay Implementation (`rr.go`)

*   **`RecordReplay` struct**: Implements the `http.RoundTripper` interface, enabling it to be used as an `http.Client` transport.
*   **Modes of Operation**:
    *   **Record Mode**: Sends the request to the real network and writes the (scrubbed) request and response to a trace file.
    *   **Replay Mode**: Reads the trace file and returns the corresponding response for a given request, completely bypassing the network.
*   **Scrubbing**: The `ScrubReq` and `ScrubResp` methods allow registration of functions to canonicalize non-deterministic data (like timestamps or IDs) and remove secrets from the log before saving.
