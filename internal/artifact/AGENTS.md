# ADK Internal Artifact Components (`internal/artifact/`)

This package provides the concrete, session-scoped implementation of the `agent.Artifacts` interface, bridging the high-level agent-facing API with the lower-level `artifact.Service`.

## Artifacts Implementation (`artifacts.go`)

*   **`Artifacts` struct**: An unexported structure that embeds a public `artifact.Service` and caches the identifying metadata: `AppName`, `UserID`, and `SessionID`.
*   **Method Implementations**: The methods (`Save`, `Load`, `LoadVersion`, `List`) translate the simplified calls from an `agent.InvocationContext` (which only provides the artifact name and data) into the full, structured request objects (`artifact.SaveRequest`, `artifact.LoadRequest`, etc.) required by the underlying artifact storage service. This ensures that artifact operations are correctly scoped to the current user's session.
*   **Interface Guarantee**: The `Artifacts` struct explicitly verifies it implements the `agent.Artifacts` interface: `var _ agent.Artifacts = (*Artifacts)(nil)`.