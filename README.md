# Corvus Studio

Corvus Studio is a local-first workspace for preparing independent game releases. Phase 3 Release + Checklist is locally implemented: developers can create and reopen a Project, generate a versioned Steam Coming Soon Release Goal with its complete Checklist, filter or add tasks, and persist task and release progress through Core restarts and browser refreshes. [GitHub Actions run 30617549600](https://github.com/gofurry/corvus-studio/actions/runs/30617549600) remains the latest successful Windows, macOS, and Linux baseline for Phase 2; Phase 3 stays in progress until its local commits are pushed and the native matrix succeeds.

Project editing, deletion, and archiving are intentionally absent. Resource, Deliverable, Evidence, Asset Map, Agent/ADK, authentication, SSE, Steam network integration, template upgrades, and system file-manager integration remain later-phase work.

## What works

- A Go workspace with separate `apps/core`, `apps/launcher`, and `agent` modules.
- `corvus serve`, backed by Cobra, Viper, Zap/lumberjack, Echo v5, modernc SQLite, embedded goose migrations, and sqlc-generated queries.
- Loopback-only `GET /healthz` readiness on `127.0.0.1:8765` by default.
- OpenAPI-first `POST /api/v1/projects`, `GET /api/v1/projects`, and `GET /api/v1/projects/{project_id}` endpoints.
- An OpenAPI-first native directory-picker endpoint with manual absolute-path fallback.
- A Project domain and repository with UUIDv7 identifiers, validated existing-directory locations, and restart persistence.
- Release Goal and Checklist domains with explicit transitions, Required-task readiness gates, transaction history, and restart persistence.
- A validated, embedded `steam-coming-soon` template v1.0.0 with eight Steam tasks and four Corvus recommendations; runtime use requires no Steam network access.
- OpenAPI-first template, Release, and Checklist endpoints with one-transaction Goal/Checklist creation and generated Go/TypeScript contracts.
- A generated Go model boundary and generated TypeScript Fetch client sourced from the same OpenAPI document.
- A React 19, TypeScript, Vite 8, Ant Design, React Router, and TanStack Query flow for Projects, Release status, URL-restorable Checklist filters, custom tasks, and task details.
- A Vite development proxy and tagged production embedding path with SPA route fallback.
- A minimal Fyne launcher window with no product UI or Core process management.
- GitHub Actions definitions for Go/frontend quality and native Windows, macOS, and Linux Core tests/builds.

## Repository map

```text
apps/core/       Core runtime, Project/Release/Checklist domains and APIs, migrations, and Web embedding
apps/web/        React/Vite Project and release-preparation application
apps/launcher/   Minimal Fyne launcher shell
agent/           Reserved Agent Go module boundary
packages/        Generated API client plus reserved shared TypeScript boundaries
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
- `zenity` or `kdialog` for the native directory picker on Linux; manual path entry remains available without either tool

Do not silently substitute older project versions. Record a toolchain mismatch before proposing a version change.

## Setup and validation

Run commands from the repository root.

```powershell
pnpm install --frozen-lockfile

Push-Location .\apps\core
go tool oapi-codegen -config openapi/oapi-codegen.yaml openapi/openapi.yaml
go tool sqlc generate
go tool sqlc vet
Pop-Location

pnpm --filter @corvus-studio/api-client generate
git diff --exit-code -- apps/core/internal/interfaces/httpserver/api apps/core/internal/infrastructure/storage/sqlc packages/api-client/src/generated

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

Open `http://localhost:5173/projects`. On the create form, use **Browse** to open the operating system directory picker, or enter an absolute path manually. Corvus only records the selected directory’s normalized identity; it does not create, move, or write files inside that directory. Opening a project means navigating to its Corvus detail page.

From a Project, open **Release**, review template `steam-coming-soon` v1.0.0, and confirm **Generate release workspace**. Corvus inserts the Release Goal and all 12 tasks atomically. Use **Checklist** to filter by status, source, or category; inspect What, Why, Requirement, Source, and Status; add user-owned tasks; and move tasks forward. Required tasks must all be Done before the Release Goal can enter Ready for review.

The same flow is available through the API:

```powershell
$projectLocation = (Resolve-Path .).Path
$projectBody = @{
  name = 'Corvus Studio'
  description = 'Local Phase 3 smoke project'
  location = $projectLocation
  language = 'English'
  stage = 'concept'
} | ConvertTo-Json

$project = Invoke-RestMethod -Method Post `
  -Uri http://127.0.0.1:8765/api/v1/projects `
  -ContentType 'application/json' `
  -Body $projectBody
Invoke-RestMethod http://127.0.0.1:8765/api/v1/projects
Invoke-RestMethod "http://127.0.0.1:8765/api/v1/projects/$($project.id)"

$template = Invoke-RestMethod http://127.0.0.1:8765/api/v1/release-templates/steam-coming-soon
$releaseBody = @{
  project_id = $project.id
  goal_type = 'steam_coming_soon'
  template_key = $template.key
  template_version = $template.version
} | ConvertTo-Json
$release = Invoke-RestMethod -Method Post `
  -Uri http://127.0.0.1:8765/api/v1/releases `
  -ContentType 'application/json' `
  -Body $releaseBody
Invoke-RestMethod "http://127.0.0.1:8765/api/v1/releases/$($release.id)/checklist"
```

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

Core accepts loopback hosts only because authentication is not implemented. See [USAGE.md](USAGE.md) for flags, Project/Release/Checklist API behavior, storage, generation, validation, and troubleshooting.

## Documentation and license

- [USAGE.md](USAGE.md) contains current Phase 3 development commands and boundaries.
- [.agent/phase-3-release-checklist.md](.agent/phase-3-release-checklist.md) is the living Phase 3 execution and evidence log.
- [docs/roadmap.md](docs/roadmap.md) tracks implementation phases and completion gates.
- [docs/](docs/) contains the source design documents.

Code is licensed under the repository's [GNU Affero General Public License v3](LICENSE). A separate documentation-license decision remains pending.
