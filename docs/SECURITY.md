# Security

`chilly` ships with checksum verification for install and update paths, plus optional provenance verification for released archives.

## What Is Verified

- Published archives include `checksums.txt`.
- `scripts/install.sh` verifies archive checksums before installation.
- `chilly self-update` verifies archive checksums before replacing the current executable.
- GitHub Actions publishes release artifact attestations for released archives.

## Verify A Release

```bash
VERSION="$(gh release view --repo chill-institute/cli --json tagName -q .tagName)"
ARCHIVE="chilly_${VERSION#v}_darwin_arm64.tar.gz"

gh release download "$VERSION" --repo chill-institute/cli --pattern "$ARCHIVE"
gh attestation verify "$ARCHIVE" --repo chill-institute/cli
```

Adjust the archive name for your platform when needed.

## Operational Notes

- Fresh installs and upgrades should use the current release line so the latest updater fixes are present.
- Attestation verification is optional for normal use, but recommended when you want provenance from the GitHub release workflow.
