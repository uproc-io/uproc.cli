# Changelog

All notable changes in `bizzmod-cli` should be documented in this file.

## 2026-03-28

### Added
- Added `bizzmod update check <CUSTOMER_API_KEY>` as a dry-run-only verification command that fetches `/api/v1/external/install?dry_run=true` and runs local read-only checks for required binaries, services, env vars, and health endpoints.

### Changed
- Root command registration now includes the new `update` command group.
- CLI documentation now includes update-check usage and explicit no-apply behavior.

### Verification
- `go test ./...`
- `go vet ./...`
- `go build -o bizzmod`
