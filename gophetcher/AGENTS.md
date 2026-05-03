# AGENTS

## Purpose

`gophetcher` is a small Go CLI that audits targets listed in a text file.
It supports `http://`, `https://`, and `file://` targets and prints per-scheme result tables.

## Working Commands

- Run the CLI with `go run ./cmd/main.go -file assets/targets.txt`.

## Entry Point

- The executable starts in [cmd/main.go](cmd/main.go).
- The main orchestration lives in [internal/analysis/analysis_handler.go](internal/analysis/analysis_handler.go).
- The CLI requires the `-file` flag and only accepts `.txt` files; see [internal/file/file_service.go](internal/file/file_service.go).

## Code Layout

- [internal/analysis](internal/analysis): routing, service interfaces, concrete analyzers, and audit result models.
- [internal/file](internal/file): command-line file selection and target file reading.
- [pkg/url_validator.go](pkg/url_validator.go): shared URL validation logic.
- [assets/targets.txt](assets/targets.txt): example input file.

## Conventions

- Keep new behavior behind the existing `AnalysisService` interface when adding another target scheme.
- Preserve the current split between orchestration in handlers and per-scheme work in services.
- Follow the existing user-facing error style: explicit, direct `error: ...` messages.
- Prefer small result-shaping helpers over formatting logic inline in loops.

## Known Repository Constraints

- Web checks use a hard-coded 12 second HTTP timeout in [internal/analysis/web_analysis_service.go](internal/analysis/web_analysis_service.go).
- File checks expect valid `file://` URLs and operate on decoded local filesystem paths in [internal/analysis/file_analysis_service.go](internal/analysis/file_analysis_service.go).
