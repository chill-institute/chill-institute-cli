# Mutate Reference

Use this reference for side-effecting commands.

## Rules

- Use `--dry-run --output json` before executing a mutation.
- Prefer `--json` or `--json @-` when you already know the exact request body.
- Keep `stdout` for the contract and ignore `stderr` for parsing.
- After a hosted mutation, rerun a read command to observe real server state instead of inferring it from the mutation response alone.

## Transfer Mutations

- Preview add-transfer: `chilly add-transfer --url "magnet:?xt=..." --dry-run --output json`
- Exact request body: `printf '{"url":"magnet:?xt=..."}' | chilly add-transfer --json @- --dry-run --output json`
- Nested alias: `printf '{"url":"magnet:?xt=..."}' | chilly user transfer add --json @- --dry-run --output json`

## Hosted Settings Mutations

- Exact request body: `printf '{"settings":{"showTopMovies":true}}' | chilly user settings set --json @- --dry-run --output json`
- Bare settings object shorthand: `printf '{"showTopMovies":true}' | chilly user settings set --json @- --dry-run --output json`
- One-field patch: `chilly user settings set sort-by title --dry-run --output json`

## Download Folder Mutations

- Patch by ID: `chilly user download-folder set 42 --dry-run --output json`
- Exact request body: `printf '{"downloadFolderId":42}' | chilly user download-folder set --json @- --dry-run --output json`
- Clear by patch: `chilly user download-folder clear --dry-run --output json`
- Clear by request body: `printf '{"settings":{"downloadFolderId":null}}' | chilly user download-folder clear --json @- --dry-run --output json`

## Local Mutations

- Local API host preview: `chilly settings set api-base-url https://api.chill.institute --dry-run --output json`
- Local API host from stdin JSON: `printf '{"key":"api-base-url","value":"https://api.chill.institute"}' | chilly settings set --json @- --dry-run --output json`
- Auth login preview: `printf '{"token":"token-from-setup","skip_verify":true}' | chilly auth login --json @- --dry-run --output json`
- Auth logout preview: `chilly auth logout --dry-run --output json`

## Update Mutations

- Check only: `chilly self-update --check --output json`
- Check with exact request body: `chilly self-update --json '{"check":true}' --output json`
- Preview install resolution: `chilly self-update --version v0.1.0 --dry-run --output json`
