# cli-editor

Experimental terminal text editor written in Go.

The project started as a single-file prototype and later evolved into a modular
codebase with terminal UI, key bindings, dialogs, syntax highlighting, file
utilities, configuration, Unicode experiments, and an AI-assistant prototype.

The original single-file implementation is preserved as
`manual_tests/test_editor.go`.

## Build

```sh
go build ./cmd/editor
```

The main command builds successfully with the toolchain used when this archive
was created.

## Historical state

This repository preserves an authentic development experiment rather than a
polished release. `go test ./...` is not currently clean:

- `manual_tests` contains multiple standalone programs with their own `main`;
- `pkg/utils/file.go` contains an unfinished `ioutil.ReadAll` call.

Generated ARM binaries, large scratch data, editor metadata, and macOS metadata
were intentionally excluded from the archive.

Originally developed during 2025 and archived in September 2026.
