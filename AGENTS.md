# CLI

`cli` is the `chill.institute` command-line client. Treat it as an agent-first SDK surface with a human CLI on top.

## Stack

- Go with Cobra for commands
- manual HTTP RPC client for the hosted API
- local config store for auth token and API base URL
- `mise` for shared repo tasks

## Commands

- `go build ./cmd/chilly`
- `go run ./cmd/chilly version --output json`
- `go test ./...`
- `mise run verify`
- `mise run fmt`

For command-surface changes, also run:

- `go run ./cmd/chilly <command> --help`

## Conventions

- Prefer machine-readable contracts first. New behavior should have a stable JSON story before nicer human formatting.
- Keep `stdout` for command results and `stderr` for prompts, progress, warnings, and recovery hints.
- Keep auth requirements, flags, and schema/describe surfaces explicit.
- If a command surface, auth flow, default, or output contract changes, update the user-facing `chilly-cli` skill in [skills/](./skills/) in the same pass.

## Read More

- command and transport architecture: [docs/ARCHITECTURE.md](./docs/ARCHITECTURE.md)
- security posture: [SECURITY.md](./SECURITY.md)
- setup, release flow, and validation: [CONTRIBUTING.md](./CONTRIBUTING.md)
