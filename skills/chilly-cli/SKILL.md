---
name: chilly-cli
description: Use `chilly` to operate chill.institute from the terminal. Trigger this skill for requests about chilly, chill institute, or putio/put.io workflows when an agent needs to authenticate, inspect CLI state, query user data, search indexers, manage transfers or downloads, change user or local settings, or inspect the current command contract. Prefer safe canonical command patterns and JSON output for automation.
---

# Chilly CLI

Use `chilly` as the local CLI entrypoint for chill.institute. For agent workflows, default to `--output json` and parse `stdout` only.

## Defaults

- Use the installed `chilly` binary directly.
- If `chilly` is not on `PATH`, stop and help the user install the CLI before continuing.
- Check `chilly settings get api-base-url --output json` before assuming which hosted environment is active.
- Use `--profile <name>` or `--config <path>` when you need isolated local state.
- Use `schema` or `--describe` when you need the current local contract before running a command.
- Use `doctor` when auth, config path, profile, or environment state looks inconsistent.
- Prefer top-level canonical commands like `search`, `whoami`, `list-top-movies`, and `add-transfer` over nested aliases.
- Use `--fields` when a read command supports it and you only need a stable subset of the payload.
- Use `--dry-run` on supported mutating commands when you need a safe preview.
- Omit `--output json` only when a human explicitly wants the built-in terminal summary.

## Auth

- Interactive login: `chilly auth login`
- Print the login URL instead of auto-opening a browser: `chilly auth login --no-browser`
- Store an existing token directly: `chilly auth login --token <token>`
- Verify the current login: `chilly whoami --output json`

If the browser is on another machine, open [https://chill.institute/auth/cli-token](https://chill.institute/auth/cli-token) in a signed-in browser and copy the token.

## Core Workflows

### Inspect The CLI Contract

- List all command and procedure metadata: `chilly schema --output json`
- Inspect one command: `chilly schema command search --output json`
- Inspect one procedure: `chilly schema procedure chill.v4.UserService/Search --output json`
- Describe a command without executing it: `chilly search --describe --output json`

### Inspect Environment And Auth

- Show current API host: `chilly settings get api-base-url --output json`
- Show current config path and profile: `chilly settings path --output json`
- Run full local diagnostics: `chilly doctor --output json`
- Verify the active user: `chilly whoami --output json`

### Search Safely

- Start with `chilly user indexers --output json`.
- Choose a suitable indexer using `tags` and any available health or status signals.
- Run one scoped search at a time, for example: `chilly search --query "dune" --indexer-id yts --output json`
- Use `--fields` when you only need a subset, for example: `chilly search --query "dune" --fields results.title --output json`

### Mutate Safely

- Preview supported mutations first, for example: `chilly add-transfer --url "magnet:?xt=..." --dry-run --output json`
- Execute the mutation with the same canonical command.
- After `add-transfer`, use `chilly get-transfer <id> --output json` to read real hosted transfer state.
- Prefer `chilly user settings set <field> <value> --output json` for routine settings changes.
- Use `chilly user settings set --json '<payload>' --output json` only for full settings payload updates.

## Command Patterns

### Environment Discovery

- Full local diagnostics: `chilly doctor --output json`
- Narrow diagnostics to specific fields: `chilly doctor --fields auth.status,config.profile --output json`
- Current API host: `chilly settings get api-base-url --output json`
- Full local config: `chilly settings show --output json`
- Config file path: `chilly settings path --output json`
- Config file path for an isolated profile: `chilly settings path --profile dev --output json`

### Schema And Describe

- List all known command and procedure metadata: `chilly schema --output json`
- Inspect one command: `chilly schema command search --output json`
- Inspect one nested command: `chilly schema command "settings get" --output json`
- Inspect one procedure: `chilly schema procedure chill.v4.UserService/Search --output json`
- Describe a command without executing it: `chilly search --describe --output json`

### Version And Update

- Show the installed version: `chilly version`
- Show build metadata: `chilly version --output json`
- Check whether a newer release exists: `chilly self-update --check --output json`
- Install the latest release over the current binary: `chilly self-update`
- Install a specific release: `chilly self-update --version v0.1.0`

### Shell Completion

- Generate zsh completions: `chilly completion zsh`
- Generate bash completions: `chilly completion bash`
- Generate fish completions: `chilly completion fish`

### Authentication

- Interactive browser-assisted login: `chilly auth login`
- Manual browser opening: `chilly auth login --no-browser`
- Existing token: `chilly auth login --token <token>`
- Logout: `chilly auth logout --output json`
- Preview logout without clearing the saved token: `chilly auth logout --dry-run --output json`
- Verify current auth: `chilly whoami --output json`

### Read Commands

- Search with the built-in terminal summary: `chilly search --query "blade runner"`
- Search as JSON: `chilly search --query "blade runner" --output json`
- Search with field selection: `chilly search --query "blade runner" --fields results.title,results.magnetLink --output json`
- Search with a specific indexer: `chilly search --query "blade runner" --indexer-id <id> --output json`
- User profile summary: `chilly whoami`
- User profile as JSON: `chilly whoami --output json`
- User profile with selected fields: `chilly whoami --fields username,email --output json`
- User indexers: `chilly user indexers --output json`
- Top movies summary: `chilly list-top-movies`
- Top movies as JSON: `chilly list-top-movies --output json`
- Top movies with selected fields: `chilly list-top-movies --fields movies.title --output json`
- User settings: `chilly user settings get --output json`
- User settings with selected fields: `chilly user settings get --fields showTopMovies,sortBy --output json`
- Current download folder: `chilly user download-folder --output json`
- Preview setting the download folder: `chilly user download-folder set 42 --dry-run --output json`
- Preview clearing the download folder: `chilly user download-folder clear --dry-run --output json`
- Inspect one folder: `chilly user folder get 0 --output json`

### Mutating Commands

- Add transfer: `chilly add-transfer --url "magnet:?xt=..." --output json`
- Get one transfer: `chilly get-transfer 42 --output json`
- Get transfer status and file fields: `chilly get-transfer 42 --fields transfer.status,transfer.statusMessage,transfer.percentDone,transfer.fileId,transfer.fileUrl --output json`
- Preview add-transfer without executing it: `chilly add-transfer --url "magnet:?xt=..." --dry-run --output json`
- Replace user settings with a full JSON payload: `chilly user settings set --json '{"showTopMovies":true}' --output json`
- Patch one setting: `chilly user settings set show-top-movies true --output json`
- Patch enum settings with friendly values: `chilly user settings set sort-by title --output json`
- Preview a one-field settings patch: `chilly user settings set sort-by title --dry-run --output json`
- Preview the full settings payload: `chilly user settings set --json '{"showTopMovies":true}' --dry-run --output json`
- Preview a local CLI settings change: `chilly settings set api-base-url https://api.chill.institute --dry-run --output json`

### Aliases

- Nested transfer add alias: `chilly user transfer add --url "magnet:?xt=..." --output json`
- Nested transfer get alias: `chilly user transfer get 42 --output json`
- Nested transfer dry-run alias: `chilly user transfer add --url "magnet:?xt=..." --dry-run --output json`

## Output And Safety

- Prefer `--output json` for automation.
- Prefer this flow for agent search work: `chilly user indexers --output json`, choose a suitable indexer using `tags` and any available health or status signals, then run one scoped `chilly search --indexer-id <id>` request at a time.
- After `add-transfer`, prefer `get-transfer` for real hosted transfer state instead of inferring it from the add response alone.
- Prefer top-level canonical commands when both top-level and `user ...` forms exist.
- Expect prompts and browser-login hints on `stderr`; parse only `stdout`.
- Expect transient loading indicators on `stderr` in pretty mode for network-backed commands.
- Human-facing notices may appear on `stderr`.
- Command data is intended to be parsed from `stdout`.
- Expect failures in `--output json` mode to appear as a single JSON envelope on `stderr`.
- Use `whoami` after auth changes when you need positive confirmation that the token works.
