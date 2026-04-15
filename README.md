# Bizzmod CLI (Go)

Minimal CLI to authenticate and call Bizzmod External API endpoints (`/api/v1/external/*`).

## Requirements

- Go 1.22+

## Setup

1. Copy env file:

```bash
cp .env.sample .env
```

2. Fill required values:

- `BIZZMOD_API_URL`
- `CUSTOMER_DOMAIN`
- `CUSTOMER_API_KEY`
- `CUSTOMER_USER_EMAIL`

You can also store credentials with `login` (recommended for local usage).

## Build and run

```bash
go mod tidy
go build -o bizzmod
./bizzmod --help
```

Or run directly:

```bash
go run . --help
```

## Distribution

This CLI is configured for multi-platform binary distribution using GoReleaser.

Targets:
- Linux: `amd64`, `arm64`
- macOS: `amd64`, `arm64`
- Windows: `amd64`, `arm64`

Packaging:
- GitHub Releases artifacts + checksums
- Homebrew tap formula update
- Scoop manifest update

### Release process

1. Run local checks:

```bash
gofmt -w .
go vet ./...
go test ./...
```

2. Optional local release dry-run:

```bash
goreleaser release --snapshot --clean
```

3. Create and push a version tag:

```bash
git tag v0.1.0
git push origin v0.1.0
```

GitHub Actions (`.github/workflows/release.yml`) will publish the release.

## Commands

### Auth

```bash
bizzmod login
```

Stores credentials in OS user config path (`bizzmod-cli/config.json`).

`login` reads credentials in this order:
- command arguments (optional, still supported)
- environment variables (`.env` or shell env)
- interactive prompt for missing values

`CUSTOMER_DOMAIN` must be the customer domain identifier (not a URL).

Example using environment values:

```bash
bizzmod login
```

If `.env` is missing or incomplete, `login` prompts for required values and writes a new `.env` file.

### Raw external request

```bash
bizzmod request <METHOD> <PATH> [JSON_BODY]
```

Example:

```bash
bizzmod request GET /api/v1/external/modules
```

Output is always rendered as readable tables/lists (never raw JSON).
When backend response includes `{ success, message, data }`, CLI prints only `data`.

### Module commands

```bash
bizzmod module list
bizzmod module get <module_slug>
bizzmod module kpis <module_slug>
bizzmod module collections <module_slug>
bizzmod module collection <module_slug> <collection_name> [--page 1 --sort-field key --sort-order asc --filter-field key --filter-value val]
bizzmod module data <module_slug> <collection_name> [--page 1 --sort-field key --sort-order asc --filter-field key --filter-value val]
bizzmod module upload <module_slug> <collection_name> <file_path>
bizzmod module webhook <module_slug> <collection_name> <payload_json>
```

### Admin commands

```bash
bizzmod admin users list [--customer-id 1]
bizzmod admin users get <user_id>

bizzmod admin customers list
bizzmod admin customers get <customer_id>

bizzmod admin credentials list [--customer-id 1 --category ai --type api_key]
bizzmod admin credentials get <credential_id>

bizzmod admin modules list
bizzmod admin modules get <module_slug>

bizzmod admin tickets list
bizzmod admin tickets get <ticket_id>

bizzmod admin logs --module-slug <module_slug> [--level all --page 1]
bizzmod admin ai-requests [--customer-id 1 --module-slug financial-reconciliation --page 1 --limit 25]
bizzmod admin changelog
```

Admin create/update subcommands are currently hidden from help output.
Admin create/update commands run interactive contract mode (contracts fetched from API):

```bash
bizzmod admin users create
bizzmod admin users update
bizzmod admin customers create
bizzmod admin customers update
bizzmod admin credentials create
bizzmod admin credentials update
bizzmod admin tickets create
bizzmod admin tickets update
```

All admin commands use external API endpoints under `/api/v1/external/admin/*`, except ticket commands that use `/api/v1/external/tickets/*`.

### Interactive mode

```bash
bizzmod interactive
```

Inside interactive mode, run commands without the binary name:

```text
bizzmod> module list
bizzmod> module get order-track
bizzmod> request GET /api/v1/external/modules
bizzmod> help
bizzmod> exit
```

### Install plan (dry-run)

```bash
bizzmod install <CUSTOMER_API_KEY> --dry-run
```

This command fetches `/api/v1/external/install` and shows the full installation plan (release versions, required services, and ordered steps) without executing changes on the server.

### Update check (dry-run only)

```bash
bizzmod update check <CUSTOMER_API_KEY>
```

This command validates update readiness using `/api/v1/external/install?dry_run=true` plus local read-only checks (docker, dokploy, required services, required env vars, and health endpoints). It never executes deployment/apply actions.

## Notes

- All calls send headers required by backend external auth:
  - `x-api-key`
  - `x-customer-domain`
  - `x-user-email`
- `request` allows calling any current/future external endpoint without waiting for a dedicated subcommand.
- CLI output is always displayed in list/table format (never JSON output).
