# Repository Guidelines

## Project Structure & Modules
- `main.go`: HTTP server, routes, and app wiring (port `3000`).
- `visitorstore/`: SQLite-backed persistence (via `modernc.org/sqlite`).
- `morse/`: Generates WAV audio for callsigns.
- `templates/` and `static/`: Embedded UI templates and assets.
- `cmd/`: Utility CLIs (e.g., `importcsv`, `exportbolt`).
- `deploy/`: Systemd unit, startup scripts, SELinux note.
- `docs/`, `images/`, `results/`: Documentation, screenshots, output data.
- Tests live next to code as `*_test.go`.

## Build, Test, and Run
- Build: `make build` → produces `./fieldday`.
- ARM build: `make build-raspi` → ARMv7 binary.
- Run locally: `./fieldday test.db` then open `http://localhost:3000`.
- Tests: `go test ./...` (optionally `go test -cover ./...`).
- Release (local snapshot): `goreleaser release --snapshot --clean` → artifacts in `dist/`.

## Coding Style & Conventions
- Language: Go. Use `gofmt` defaults and idiomatic Go.
- Format/lint: `gofmt -s -w .`, `go vet ./...` before PRs.
- Packages: lowercase short names (`visitorstore`, `morse`).
- Exports: PascalCase for exported, lowerCamel for unexported.
- Files: `snake_case.go`; tests as `name_test.go` with table-driven tests when useful.

## Testing Guidelines
- Framework: standard `testing` package.
- Scope: prefer unit tests near implementation (`visitorstore`, `morse`, handlers in `main_test.go`).
- DB tests: use throwaway SQLite files (see tests) and `defer os.Remove(...)`.
- Run all: `go test ./...` before pushing.

## Commit & Pull Requests
- Commits: short, imperative subject (e.g., "Add goreleaser config"), optionally include scope and PR/issue refs ("(#10)").
- PRs: include purpose, implementation notes, how to test (commands/URLs), and screenshots for UI changes (`images/`). Link related issues.

## Security & Configuration
- Data path: first CLI arg is the SQLite file; it is created if missing.
- Secrets: none required. Avoid committing data—`.gitignore` excludes `*.db`, `*.csv`, `results/`, and `dist/`.
- Service: see `deploy/fieldday.service`; SELinux note in `Makefile` (httpd_can_network_connect).

## Architecture Overview
- Minimal `net/http` server, embedded assets via `go:embed`, persistence in `visitorstore`, and browser-played Morse via `/morse-audio`.
