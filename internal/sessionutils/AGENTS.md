# ADK Internal Session Utilities (`internal/sessionutils/`)

This package provides utility functions for manipulating and querying session data structures, specifically the list of conversation events.

## Event Utilities (`utils.go`)

*   **`FilterEventsByBranch` function**: Takes a slice of `session.Event` and a branch identifier (`string`).
*   **Purpose**: It returns a new slice of events, containing only those that belong to the specified branch. This is crucial for isolating the conversation history relevant to a particular agent or workflow branch during execution.