Implement a `memory/filesystem.go` that uses a configurable base directory and allows organizing `<basedir>/memory/<path>/.../<name>.md` where the path is determined by `<app>/<user>/<session-id>/memory/<path-arg>.md`

Ensure you implement the `memory.Service` interface so that we can use this instead of the `memory/inmemory.go`

1. Read all files in `memory/` to understand the current code
2. Understand the semantics of the Service interface and inmemory implementation
3. Write a filesystem implementation that strictly adheres to the `Service` interface
4. Do not duplicate implementation that already exists

The filesystem implementation should

1. Use the inmemory implementation internally, we don't want to read all the files in every search
2. When a new session is added, it gets added to the inmemory and writen to disk as simple json.
3. Add a load function that will walk a directory of session json files and fill the inmemory again