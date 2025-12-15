# ADK Internal CLI Components (`internal/cli/`)

The `internal/cli/` package is an internal utility layer that provides helper functions for building command-line applications related to the ADK. It is structured entirely within a single `util/` sub-package.

## Utilities (`util/`)

| File | Description |
| :--- | :--- |
| `flagset_helpers.go` | Contains `FormatFlagUsage` to capture and return the usage string for a standard Go `flag.FlagSet`. |
| `oscmd.go` | Provides utilities for executing operating system commands with structured, colored logging. It includes `LogStartStop` to track command lifecycle and `LogCommand` to run an `exec.Cmd` while routing its streams through a wrapper (`reprintableStream`) for pretty-printing. |
| `text_helpers.go` | Contains simple text formatting functions, such as `CenterString`, used for CLI output presentation. |