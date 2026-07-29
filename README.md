# Corvus Studio

Corvus Studio is a local-first workspace for preparing independent game releases. **Phase 0: Repository Bootstrap is complete**, and Phase 1 has not started: the monorepo, buildable application shells, local validation commands, and CI definition exist, but product features do not.

## What works in Phase 0

- A Go workspace with separate `apps/core`, `apps/launcher`, and `agent` modules.
- A minimal Core executable that prints `Corvus Studio core bootstrap` and exits.
- A minimal Fyne launcher window with no product UI or Core process management.
- A React 19, TypeScript, and Vite 8 Web bootstrap page with format, lint, test, and build checks.
- GitHub Actions definitions for Go/frontend quality and native Windows, macOS, and Linux builds.

Echo, SQLite, migrations, OpenAPI generation, ADK Agent behavior, authentication, Steam workflows, and all product domains begin in later phases.

## Repository map

```text
apps/core/       Go Core executable shell
apps/web/        React/Vite Web application shell
apps/launcher/   Minimal Fyne launcher shell
agent/           Reserved Agent Go module boundary
packages/        Reserved shared TypeScript package boundaries
configs/         Repository-wide configuration boundary
deployments/     Future deployment boundary; no deployment implementation
scripts/         Future portable automation boundary
tests/           Future integration, scenario, and fixture boundaries
tools/           Future repository-owned tool boundary
docs/            Product and engineering design documents
```

### Why some frontend files are at the repository root

This is a pnpm monorepo, not an isolated frontend repository. `pnpm-workspace.yaml`, the single `pnpm-lock.yaml`, the root `package.json`, `.npmrc`, and `.node-version` define consistent installation and commands for `apps/web` and future packages under `packages/`. The shared Prettier files also belong at the root so those packages use one formatting policy. Application-specific source, Vite, TypeScript, ESLint, tests, and styles remain under `apps/web/`.

The Go equivalents—`go.work`, `go.work.sum`, `.go-version`, and `.golangci.yml`—coordinate the three Go modules from the repository root.

## Prerequisites

- Go 1.26.5
- Node.js 24.15.x
- pnpm 10.11.0
- A native compiler and the [Fyne prerequisites](https://docs.fyne.io/started/) for your operating system when building or running the Launcher
- golangci-lint 2.12.2 for the local Go lint command

Do not substitute older project versions silently. If a required version is unavailable, record the mismatch before changing the repository.

## Setup and validation

Run commands from the repository root.

```powershell
pnpm install --frozen-lockfile

go test ./apps/core/... ./apps/launcher/... ./agent/...
golangci-lint run ./apps/core/... ./apps/launcher/... ./agent/...

pnpm format:check
pnpm lint
pnpm test
pnpm build
```

The frontend build succeeds when `apps/web/dist/index.html` is created. GitHub Actions runs the same quality commands and additionally builds Core and Launcher natively on Windows, macOS, and Linux.

## Run the bootstrap shells

```powershell
# Prints the bootstrap marker and exits.
go run ./apps/core/cmd/corvus

# Starts the React development server.
pnpm dev:web

# Opens the minimal desktop launcher window.
go run ./apps/launcher
```

There is no `corvus serve` command in Phase 0, and the Web dev server is not yet connected to Core. Production embedding is also reserved for a later phase.

## Documentation and license

- [USAGE.md](USAGE.md) separates current Phase 0 commands from future product workflows.
- [.agent/phase-0-repository-bootstrap.md](.agent/phase-0-repository-bootstrap.md) is the executable Phase 0 plan and evidence log.
- [docs/roadmap.md](docs/roadmap.md) tracks implementation phases and completion evidence.
- [docs/](docs/) contains the source design documents.

Code is licensed under the repository's [GNU Affero General Public License v3](LICENSE). A separate documentation-license decision remains pending.
