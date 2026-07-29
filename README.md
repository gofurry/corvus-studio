# Corvus Studio

Corvus Studio is a local-first workspace for preparing independent game releases. Phase 1 now has a working local Core runtime: configuration, structured logs, SQLite migrations, a readiness endpoint, graceful shutdown, and an optional embedded React build are implemented. Remote Windows/macOS/Linux CI evidence for this Phase 1 revision remains pending until the local commits are pushed.

Product domains are intentionally absent. Project, Release Goal, Checklist, Resource, Deliverable, Steam templates, Asset Map, Agent/ADK, authentication, SSE, and OpenAPI business APIs begin in later phases.

## What works

- A Go workspace with separate `apps/core`, `apps/launcher`, and `agent` modules.
- `corvus serve`, backed by Cobra, Viper, Zap/lumberjack, Echo v5, modernc SQLite, embedded goose migrations, and sqlc-generated queries.
- Loopback-only `GET /healthz` readiness on `127.0.0.1:8765` by default.
- A React 19, TypeScript, and Vite 8 Web shell with a development proxy and a tagged production embedding path.
- A minimal Fyne launcher window with no product UI or Core process management.
- GitHub Actions definitions for Go/frontend quality and native Windows, macOS, and Linux Core tests/builds.

## Repository map

```text
apps/core/       Go Core runtime, migrations, queries, and Web embedding boundary
apps/web/        React/Vite Web application shell
apps/launcher/   Minimal Fyne launcher shell
agent/           Reserved Agent Go module boundary
packages/        Reserved shared TypeScript package boundaries
configs/         Repository-wide configuration boundary
deployments/     Future deployment boundary; no deployment implementation
scripts/         Future portable automation boundary
tests/           Future integration, scenario, and fixture boundaries
tools/           Future repository-owned tool boundary
docs/            Product and engineering source documents plus the roadmap
```

### Why workspace files are at the repository root

This is a pnpm monorepo, not an isolated frontend repository. `pnpm-workspace.yaml`, the single `pnpm-lock.yaml`, the root `package.json`, `.npmrc`, and `.node-version` coordinate `apps/web` and future packages under `packages/`. Shared Prettier files define one formatting policy. Application-specific source, Vite, TypeScript, ESLint, tests, and styles remain under `apps/web/`.

The Go equivalents—`go.work`, `go.work.sum`, `.go-version`, and `.golangci.yml`—coordinate the three Go modules. These workspace files do not add frontend code to the backend or become part of the Core binary.

## Prerequisites

- Go 1.26.5
- Node.js 24.15.x
- pnpm 10.11.0
- golangci-lint 2.12.2 for local linting
- A native compiler and the [Fyne prerequisites](https://docs.fyne.io/started/) when building or running the Launcher

Do not silently substitute older project versions. Record a toolchain mismatch before proposing a version change.

## Setup and validation

Run commands from the repository root.

```powershell
pnpm install --frozen-lockfile

Push-Location .\apps\core
go tool sqlc generate
go tool sqlc vet
Pop-Location
git diff --exit-code -- apps/core/internal/infrastructure/storage/sqlc

go test ./apps/core/... ./apps/launcher/... ./agent/...
golangci-lint run ./apps/core/... ./apps/launcher/... ./agent/...

pnpm format:check
pnpm lint
pnpm test
pnpm build
```

## Run Core and Web locally

Start Core with an explicit development data directory:

```powershell
go run ./apps/core/cmd/corvus serve --data-dir .tmp/core-runtime
```

Core creates `.tmp/core-runtime/corvus.db` and `.tmp/core-runtime/logs/corvus.log`, applies pending embedded migrations before listening, and serves readiness at `http://127.0.0.1:8765/healthz`. Stop it with Ctrl+C.

In another terminal, start Vite:

```powershell
pnpm dev:web
```

Vite proxies `/api` and `/healthz` to Core. Set `CORVUS_CORE_URL` before starting Vite when Core uses a different loopback port.

To validate production Web embedding:

```powershell
pnpm build
go run ./apps/core/cmd/stage-web
New-Item -ItemType Directory -Force .tmp | Out-Null
go build -tags corvus_webui -o .tmp/corvus-embedded.exe ./apps/core/cmd/corvus
```

The staging directory and build artifacts are generated and ignored. A normal Core build remains independent of Node and frontend assets.

## Configuration

Configuration precedence is `CLI > CORVUS_* environment variables > YAML/JSON file > defaults`. Copy `apps/core/configs/corvus.example.yaml` when an explicit file is useful:

```powershell
go run ./apps/core/cmd/corvus serve --config .\apps\core\configs\corvus.example.yaml
```

Phase 1 accepts loopback hosts only because authentication is not implemented. See [USAGE.md](USAGE.md) for flags, environment variables, storage behavior, and troubleshooting.

## Documentation and license

- [USAGE.md](USAGE.md) contains current Phase 1 development commands and boundaries.
- [.agent/phase-1-core-runtime.md](.agent/phase-1-core-runtime.md) is the living Phase 1 execution and evidence log.
- [docs/roadmap.md](docs/roadmap.md) tracks implementation phases and completion gates.
- [docs/](docs/) contains the source design documents.

Code is licensed under the repository's [GNU Affero General Public License v3](LICENSE). A separate documentation-license decision remains pending.
