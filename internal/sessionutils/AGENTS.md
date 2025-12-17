# Package `internal/sessionutils`

This package provides utility functions for manipulating and querying session data structures, specifically state maps and event lists.

## State Management (`utils.go`)

Handles the partitioning and merging of state keys based on prefixes.

**Prefixes:**
*   `app:`: Application-scoped state.
*   `user:`: User-scoped state.
*   `temp:`: Temporary state (not persisted).

**Key Functions:**
*   `ExtractStateDeltas(delta map[string]any)`: Splits a single delta map into `appState`, `userState`, and `sessionState` deltas based on prefixes. Ignores `temp:` keys.
*   `MergeStates(app, user, session)`: Combines the three state maps into a single map for client consumption, re-adding the prefixes.
