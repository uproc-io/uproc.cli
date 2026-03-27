# AGENTS.md

This file is the CLI subproject guide for agentic coding in `cli/`.
The `cli/` folder is a standalone repository mounted via symlink in the mono-repo.

--------------------------------------------------------------------------------
Scope and precedence
--------------------------------------------------------------------------------

- Applies only to files under `cli/`.
- When this file conflicts with the mono-repo `AGENTS.md`, this file is authoritative for `cli/` changes.
- Keep changes minimal and aligned with existing CLI repository conventions.

--------------------------------------------------------------------------------
Build, lint, and test commands
--------------------------------------------------------------------------------

- Install deps: `go mod tidy`
- Build binary: `go build -o bizzmod`
- Run without build: `go run . --help`
- Format: `gofmt -w .`
- Lint (if available): `go vet ./...`
- Test: `go test ./...`
- Release local dry-run: `goreleaser release --snapshot --clean`
- Release by tag (CI): push tag `vX.Y.Z`

When adding new CLI commands:
- Prefer wrapping existing `/api/v1/external/*` endpoints directly.
- Keep a generic raw command (currently `request`) to avoid endpoint coverage gaps.
- Update `README.md` command list in the same change set.

Install command policy:
- `install` must consume `/api/v1/external/install` and render a full installation plan.
- `install` supports `--dry-run` to print every step without executing server changes.
- Default/expected operational usage is dry-run preview first; any future execution mode must be explicit and opt-in.

Distribution notes:
- Release automation is defined in `.github/workflows/release.yml` and `.goreleaser.yml`.
- Artifacts are produced for Linux/macOS/Windows on `amd64` + `arm64`.
- Packaging targets include GitHub Releases, Homebrew tap, and Scoop bucket.

Authentication UX notes:
- `login` supports args, env vars, and interactive prompt fallback.
- If `.env` is missing/incomplete and no args are provided, `login` prompts for required values and creates `.env`.

--------------------------------------------------------------------------------
Code style and safety
--------------------------------------------------------------------------------

- Prefer explicit, minimal changes over broad refactors.
- Follow the style already present in the touched files.
- Do not commit secrets or environment files.
- Add comments only when needed to explain non-obvious behavior.

--------------------------------------------------------------------------------
Changelog policy
--------------------------------------------------------------------------------

- If a CLI changelog is introduced, record behavior/contract changes in it within the same change set.
- Include verification commands once canonical test/build commands are documented.

Keep this file updated when CLI commands, conventions, or release workflows are formalized.
