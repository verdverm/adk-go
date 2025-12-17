# Package `internal/testutil`

This package provides comprehensive fixtures and utilities necessary for deterministic and reliable testing of ADK agents, models, and execution components.

## LLM Testing Utilities (`genai.go`)

### Gemini Transport

The core utility is `NewGeminiTransport`, which integrates the `internal/httprr` record/replay mechanism with the Google GenAI SDK.

```go
// NewGeminiTransport returns the genai.ClientConfig configured for record and replay.
func NewGeminiTransport(rrfile string) (http.RoundTripper, error)
```

### Determinism & Scrubbing

To ensure that tests are stable (deterministic) and safe (no secrets logged), this package implements a rigorous request scrubbing logic (`scrubGeminiRequest`):

1.  **Header Removal**: Strips variable headers like `X-Goog-Api-Key`, `X-Goog-Api-Client`, and `User-Agent`.
2.  **JSON Canonicalization**: The Gemini API uses protobuf-generated JSON, which can have non-deterministic whitespace. This utility parses and re-serializes (compacts) the JSON body to ensure byte-for-byte equality across test runs.

### Test Agent Runner (`test_agent_runner.go`)

Provides helper functions to execute agents in a simplified test harness, mocking necessary session context.
