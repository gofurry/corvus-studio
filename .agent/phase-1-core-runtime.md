# Corvus Studio Phase 1 — Core Runtime ExecPlan

This ExecPlan is a living document. Update its Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective sections whenever implementation evidence or assumptions change.

All commands run from `E:\Git\开源\agent\corvus` in PowerShell unless a different working directory is stated.

## 1. Purpose and user-visible outcome

Phase 1 replaces the exiting Phase 0 Core bootstrap with a real local runtime. After completion, a developer can run:

```powershell
go run ./apps/core/cmd/corvus serve --data-dir .tmp/core-runtime
```

and observe a long-running Echo v5 server on `127.0.0.1:8765`, a migrated SQLite database under the selected data directory, structured console/file logs, and a successful `GET /healthz` response. Ctrl+C or context cancellation must stop the HTTP server and close storage cleanly.

The Vite dev server remains a separate process and proxies Core requests during development. A production validation path builds Vite, stages the generated assets inside the Core module, and builds a tagged Core binary that serves the embedded SPA. The default Go build remains independent of generated Web assets.

## 2. Current repository assessment

At Phase 1 start on 2026-07-29:

- branch `dev` points to `caf77e6` and is three commits ahead of `origin/dev`;
- the worktree is clean;
- the root Go workspace contains exactly `apps/core`, `apps/launcher`, and `agent`;
- `apps/core` has only a dependency-free bootstrap `main.go`, its smoke test, and `go.mod`;
- no Echo runtime, configuration loader, logger, database, migration, generated SQL, HTTP route, or production Web embedding exists;
- the pnpm workspace and minimal Web project build successfully according to Phase 0 evidence;
- GitHub Actions already runs Go/frontend quality and native Windows/macOS/Linux builds;
- Phase 0 is complete; no Phase 2+ product domain exists.

Observed toolchain:

| Tool          | Version/evidence                                                                    |
| ------------- | ----------------------------------------------------------------------------------- |
| Go            | `go1.26.5 windows/amd64`                                                            |
| GOWORK        | root `go.work`                                                                      |
| pnpm          | pinned `10.11.0`                                                                    |
| golangci-lint | project CI pins `v2.12.2`; local official v2 archive may be used if PATH remains v1 |

Read-only module queries on 2026-07-29 found these selected exact versions:

| Component      | Module                             | Version   |
| -------------- | ---------------------------------- | --------- |
| Echo           | `github.com/labstack/echo/v5`      | `v5.3.1`  |
| Viper          | `github.com/spf13/viper`           | `v1.21.0` |
| Zap            | `go.uber.org/zap`                  | `v1.28.0` |
| lumberjack     | `gopkg.in/natefinch/lumberjack.v2` | `v2.2.1`  |
| Cobra          | `github.com/spf13/cobra`           | `v1.10.2` |
| modernc SQLite | `modernc.org/sqlite`               | `v1.55.0` |
| goose          | `github.com/pressly/goose/v3`      | `v3.27.3` |
| sqlc tool      | `github.com/sqlc-dev/sqlc`         | `v1.31.1` |

## 3. Source-of-truth documents

Conflicts use this precedence:

1. `docs/development/Corvus_Studio_Development_Implementation_Plan_v0.1.md` defines Phase 1 as Echo, Viper, Zap, SQLite, goose, sqlc, and Cobra, accepted when `corvus serve` runs.
2. `docs/architecture/Corvus_Studio_Repository_Structure_Design_v0.1.md` places Core commands, internal layers, migrations, and configs under `apps/core`.
3. `docs/architecture/Corvus_Studio_Technology_Stack_Decision_v0.1.md` fixes the selected libraries, local-first model, configuration precedence, and pure-Go SQLite driver.
4. `docs/deployment/Corvus_Studio_Engineering_and_Deployment_Guide_v0.1.md` requires GitHub Actions, embedded production Web assets, migration-before-start, logging rotation, and three-platform support.
5. `USAGE.md` defines the developer experience and must be updated from Phase 0 claims to verified Phase 1 commands.
6. `docs/architecture/Corvus_Studio_Launch_System_Design_Document_v0.1.md` supplies local-first storage, single-writer, Runtime Config, and security boundaries.
7. `docs/api/Corvus_Studio_Backend_API_Design_v0.1.md` reserves `/api/v1` for later business APIs; Phase 1 adds no product endpoint.
8. `.agent/phase-0-repository-bootstrap.md` records the accepted development/prod Web delivery boundary and explicitly assigns staging/embedding to Phase 1.
9. `docs/roadmap.md` provides the maintained Phase 1 tasks and completion gates.

Official upstream evidence is used only to verify selected APIs and compatibility: Echo v5 graceful shutdown, Viper precedence, goose embedded/context-aware providers, sqlc v2 configuration, and modernc SQLite DSN pragmas.

## 4. Scope

Phase 1 will:

- replace the bootstrap command with a Cobra root and working `serve` subcommand;
- load defaults, optional YAML/JSON config, `CORVUS_*` environment variables, and CLI overrides with `CLI > ENV > config file > defaults` precedence;
- default to loopback-only `127.0.0.1:8765` and reject non-loopback hosts while authentication is absent;
- place runtime data under the OS user config directory by default and support an explicit `--data-dir` for development/tests;
- create structured Zap console/file logging with lumberjack rotation;
- open SQLite through `modernc.org/sqlite` with foreign keys, WAL, normal synchronous mode, and a busy timeout;
- embed goose migrations in the Core binary and apply pending migrations before HTTP startup;
- establish a minimal infrastructure-only runtime metadata schema and sqlc-generated query boundary;
- expose only operational `GET /healthz`, including database and schema readiness;
- shut down cleanly on Ctrl+C/SIGTERM or cancelled context;
- proxy `/api` and `/healthz` from Vite to the configured Core development address;
- provide a portable Web staging command and tagged embedded-SPA Core build;
- add tests and CI checks for configuration, logging, migrations, generated code drift, health behavior, lifecycle, and native SQLite runtime.

## 5. Non-goals

Phase 1 will not implement:

- Project, Release Goal, Checklist, Resource, Deliverable, Template, Asset Map, Codex, Review, Report, Export, or Watch behavior;
- any `/api/v1` business route or OpenAPI contract;
- login, authentication, authorization, API keys, secure credential storage, or remote network exposure;
- SSE, Agent/ADK, model providers, tools, prompts, or workflows;
- Launcher Core lifecycle, tray, browser opening, or business UI;
- user-facing backup/restore, `corvus backup`, `corvus doctor`, or deployment artifacts;
- Docker, systemd, installers, signing, notarization, automatic updates, or release packaging.

## 6. Document conflicts and decisions

| Status   | Conflict or ambiguity                                                                                                              | Decision and rationale                                                                                                                                                               |
| -------- | ---------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| Accepted | No default Core port is specified.                                                                                                 | Use configurable `127.0.0.1:8765`. Loopback is safe without authentication and gives `corvus serve` a zero-argument path.                                                            |
| Accepted | Config/data/log directories are unspecified.                                                                                       | Search optional config in the OS user config directory under `corvus-studio`; derive database/log defaults from `runtime.data_dir`; explicit `--config` missing or invalid is fatal. |
| Accepted | API design contains only future business routes.                                                                                   | Add `/healthz` outside `/api/v1` as an operational endpoint; do not add placeholder business APIs.                                                                                   |
| Accepted | sqlc normally needs schema and queries, but business tables begin in Phase 2.                                                      | Create only `corvus_runtime_metadata` plus infrastructure health/metadata queries. This proves SQLite/goose/sqlc without anticipating product domains.                               |
| Accepted | Engineering migration flow requires backup, while user-facing Backup is Phase 7.                                                   | Phase 1 may create an automatic pre-migration SQLite snapshot for an existing non-empty database. It exposes no user backup workflow or command.                                     |
| Accepted | Repository design lists `migrate`, `doctor`, and `backup` commands, but the Phase 1 acceptance is only `corvus serve`.             | Implement only `serve`. Pending migrations run automatically before serving; additional commands remain in their owning later decisions.                                             |
| Accepted | Production Web embedding needs generated files inside the Core module, but native Core builds must not require a prior pnpm build. | Default builds exclude Web assets. `corvus_webui` build tag enables an embedded filesystem after an explicit portable staging command. Generated staged files remain ignored.        |
| Accepted | Headless/Docker modes may later need non-loopback binding, but Phase 1 has no authentication.                                      | Reject non-loopback hosts now. A later security/deployment decision may deliberately expand this boundary.                                                                           |

## 7. Proposed repository tree

Only relevant Phase 1 additions are shown:

```text
apps/core/
├── cmd/
│   ├── corvus/
│   │   ├── main.go
│   │   └── main_test.go
│   └── stage-web/
│       └── main.go
├── configs/
│   └── corvus.example.yaml
├── internal/
│   ├── application/
│   │   ├── runtime.go
│   │   └── runtime_test.go
│   ├── infrastructure/
│   │   ├── config/
│   │   ├── logging/
│   │   └── storage/
│   │       └── sqlc/                 # generated, committed
│   └── interfaces/
│       ├── cli/
│       ├── httpserver/
│       └── webui/
│           └── dist/                 # generated, ignored
├── migrations/
│   ├── 00001_core_runtime.sql
│   └── embed.go
├── queries/
│   └── runtime.sql
├── go.mod
├── go.sum
└── sqlc.yaml
```

No `internal/domain` directory is created until Phase 2 has a real Project domain.

## 8. Milestones

### Milestone 1 — Configuration and logging foundation

**Goal:** pin Viper/Zap/lumberjack and validate deterministic configuration precedence and rotating file logging.

**Paths:** `apps/core/go.mod`, `apps/core/go.sum`, `apps/core/configs/`, `apps/core/internal/infrastructure/config/`, `apps/core/internal/infrastructure/logging/`.

**Operations:** add exact dependencies; implement default path resolution, YAML/JSON/env/CLI loading and validation; implement Zap console + lumberjack file cores; add tests and example config.

**Validation:**

```powershell
go test ./apps/core/internal/infrastructure/config/... ./apps/core/internal/infrastructure/logging/...
$files = Get-ChildItem -Recurse -File -Include *.go -Path '.\apps\core'
$unformatted = $files | ForEach-Object { gofmt -l $_.FullName }
if ($unformatted) { throw "Unformatted Go files: $unformatted" }
```

Expected: tests exit 0; no unformatted paths; config tests prove default < file < env < CLI; logging test observes JSON in the configured file.

**Recovery:** fix only the failing package or revert the milestone commit with `git revert`; do not delete unrelated user configuration.

### Milestone 2 — SQLite, goose, and sqlc foundation

**Goal:** create and reopen a migrated pure-Go SQLite database with committed generated query code.

**Paths:** Core module dependency files, `migrations/`, `queries/`, `sqlc.yaml`, `internal/infrastructure/storage/`, generated `storage/sqlc/`.

**Operations:** pin modernc/goose/sqlc; add the sqlc Go tool directive; create the infrastructure migration and queries; generate code; implement open/ping/version/migration/close; add idempotence and pre-migration snapshot tests.

**Validation:**

```powershell
Push-Location .\apps\core
go tool sqlc generate
go tool sqlc vet
Pop-Location
git diff --exit-code -- apps/core/internal/infrastructure/storage/sqlc
go test ./apps/core/internal/infrastructure/storage/...
```

Expected: generation/vet/tests exit 0; generated code is unchanged after regeneration; a temporary DB reaches schema version 1 twice without duplicate work; the sqlc health query returns 1.

**Recovery:** retain failed DB fixtures only outside the repository for diagnosis; repair migration code with a new migration after it is committed—never silently rewrite an applied migration.

### Milestone 3 — Echo, Cobra, lifecycle, and Web delivery

**Goal:** make `corvus serve` observable and stoppable, with dev proxy and optional production SPA embedding.

**Paths:** Core application/interfaces/cmd packages, Web staging command, `apps/web/vite.config.ts`, `.gitignore`.

**Operations:** pin Echo/Cobra; compose config/logger/storage/server; add signal context; add `/healthz`; implement graceful shutdown; implement loopback guard; add staging and tagged embed packages; configure Vite proxy; add tests.

**Validation:**

```powershell
go test ./apps/core/...
go run ./apps/core/cmd/corvus --help

$runtimeDir = Join-Path ([IO.Path]::GetTempPath()) 'corvus-phase1-smoke'
$process = Start-Process -FilePath 'go' -ArgumentList @('run','./apps/core/cmd/corvus','serve','--data-dir',$runtimeDir,'--port','18765') -PassThru -WindowStyle Hidden
try {
    $response = Invoke-RestMethod -Uri 'http://127.0.0.1:18765/healthz'
    if ($response.status -ne 'ok') { throw 'Core health check failed' }
} finally {
    Stop-Process -Id $process.Id -ErrorAction SilentlyContinue
}

pnpm build
go run ./apps/core/cmd/stage-web
go build -tags corvus_webui -o .tmp/corvus-embedded ./apps/core/cmd/corvus
```

Expected: tests and builds exit 0; help contains `serve`; health returns `status=ok`; database/log files exist; the embedded build succeeds after staging.

**Recovery:** terminate only the recorded smoke-process PID, preserve temporary runtime files for diagnosis, and never recursively remove a computed repository path.

### Milestone 4 — CI, documentation, and final audit

**Goal:** align CI and developer docs with the exact implemented commands and close Phase 1 with current evidence.

**Paths:** `.github/workflows/ci.yml`, `README.md`, `USAGE.md`, `docs/roadmap.md`, this ExecPlan.

**Operations:** add sqlc drift and embedded build checks; run Core SQLite tests on all native runners; document config/env/serve/health/storage/logs/build commands; mark roadmap tasks only after evidence; audit forbidden Phase 2+ content.

**Validation:**

```powershell
go test ./apps/core/... ./apps/launcher/... ./agent/...
golangci-lint run ./apps/core/... ./apps/launcher/... ./agent/...
pnpm install --frozen-lockfile
pnpm format:check
pnpm lint
pnpm test
pnpm build
git diff --check
```

Expected: all local commands exit 0; CI YAML parses; docs match implementation; no Phase 2+ symbol/schema/route is introduced. Remote three-platform completion remains pending until an actual Actions run exists.

**Recovery:** correct only the failing implementation/docs/CI slice, rerun its validation, and create a follow-up local commit. Do not amend, force-push, or weaken checks.

## 9. Go workspace strategy

The root `go.work` remains unchanged with exactly three modules. All Phase 1 dependencies and sqlc tool pin belong to `apps/core/go.mod`; they must not leak into Launcher or Agent. The canonical cross-module test remains:

```powershell
go test ./apps/core/... ./apps/launcher/... ./agent/...
```

`go test ./...` from the non-module root remains invalid as full workspace coverage.

## 10. pnpm workspace strategy

No new frontend library is needed. Phase 1 modifies only Vite dev proxy configuration and uses the existing frozen workspace. Ant Design, Zustand, TanStack Query, React Router, React Flow, and react-markdown remain deferred. The single root `pnpm-lock.yaml` must not change unless a real frontend dependency is added.

## 11. Build and embedding strategy

Development mode runs:

```powershell
go run ./apps/core/cmd/corvus serve
pnpm dev:web
```

Vite proxies `/api` and `/healthz` to `http://127.0.0.1:8765`; HMR remains on the Vite port.

Default Core builds contain no Web assets. Production validation runs Vite build, stages `apps/web/dist` to ignored `apps/core/internal/interfaces/webui/dist`, then builds with `-tags corvus_webui`. The Echo catch-all serves embedded assets and SPA fallback only in that tagged binary. API/health routes retain precedence.

The Launcher is unchanged and does not start or stop Core in this phase.

## 12. CI strategy

Retain GitHub Actions as the only CI. Required Phase 1 additions:

- Ubuntu Go quality regenerates sqlc and fails on generated-code drift;
- Ubuntu frontend quality stages the built Web output and compiles the tagged embedded Core;
- Windows/macOS/Linux native jobs run Core tests so modernc SQLite migrations execute natively, then build Core and Launcher as before;
- caches continue to key on all `go.sum`, `go.work.sum`, and `pnpm-lock.yaml` files.

Issue/PR templates and `CONTRIBUTING.md` remain absent.

## 13. Validation matrix

| Target                 | Command/observation                        | Required evidence                          |
| ---------------------- | ------------------------------------------ | ------------------------------------------ |
| Config defaults        | config package tests                       | loopback, port 8765, OS user path          |
| Precedence             | config package tests                       | CLI > ENV > YAML/JSON > defaults           |
| Logging                | logging package tests                      | JSON file entry and valid level            |
| SQLite                 | storage tests                              | pure-Go open/ping/close on temp path       |
| goose                  | storage tests                              | schema version 1 and idempotent reopen     |
| sqlc                   | `go tool sqlc generate`; `vet`; Git diff   | deterministic committed output             |
| HTTP                   | httptest/server tests                      | `/healthz` 200/503 behavior                |
| Lifecycle              | application/CLI tests and process smoke    | start, cancellation, graceful exit         |
| Dev proxy              | Vite config test/build inspection          | `/api` and `/healthz` target loopback Core |
| Embedded Web           | stage + tagged Go build/test               | SPA index served by tagged handler         |
| Windows local          | full Go tests, build, process health smoke | exit 0 and observable files/JSON           |
| Linux/macOS/Windows CI | native Core tests/build                    | actual successful Actions jobs             |
| Phase boundary         | dependency/path/symbol audit               | no Phase 2+ business implementation        |

## 14. Idempotence and recovery

- Config loading never writes a config file and treats an absent implicit file as normal.
- Runtime directories use `MkdirAll`; existing files are never overwritten except append/rotation-managed logs and SQLite-managed files.
- Goose owns schema versioning; repeated startup applies no duplicate migration.
- Committed migrations are immutable after use; corrections require a higher-numbered migration.
- sqlc generation is deterministic and CI rejects drift.
- Web staging writes only to the resolved ignored `apps/core/internal/interfaces/webui/dist` target after validating source/destination boundaries.
- Smoke tests use unique OS temporary directories and recorded PIDs/listeners.
- Every logical milestone is a local commit. Safe rollback uses `git revert <commit>`; no hard reset, clean, or destructive checkout.

## 15. Risks and blockers

- Echo v5.3.1 and modernc/sqlite v1.55.0 are recent releases; tests must exercise their actual current APIs and all native platforms.
- sqlc as a pinned Go tool materially expands `go.sum` and download/cache size; it must remain Core-local and must not become a fourth workspace module.
- Windows process smoke needs careful PID cleanup because `go run` may spawn a child binary; package-level lifecycle tests are the primary deterministic evidence.
- SQLite migration backup semantics must be tested on a real existing database; a failed snapshot must stop migration.
- WAL creates `-wal` and `-shm` sidecars during use; shutdown tests must verify clean closure without committing runtime files.
- Loopback-only binding deliberately limits headless remote access until authentication/deployment decisions exist.
- Default native CI can be slow because Fyne and modernc both compile substantial platform code.
- Remote GitHub Actions evidence cannot be claimed until the user pushes or otherwise authorizes a run.

## 16. Completion criteria

Phase 1 is complete only when evidence proves all of the following:

- exact selected dependencies are pinned in Core and no unrelated module changes occur;
- `corvus serve` starts without required arguments, exposes healthy loopback HTTP, and stops gracefully;
- config precedence and validation work for defaults, YAML/JSON, environment, and CLI;
- Zap/lumberjack writes structured rotating file logs;
- modernc SQLite opens cross-platform, goose applies the embedded initial migration idempotently, and sqlc output is reproducible;
- an existing database is protected before any pending migration is applied;
- development proxy and tagged production embedding both build as designed;
- local full Go/frontend format, lint, test, and build commands pass;
- actual Windows/macOS/Linux CI runs pass before claiming remote completion;
- README, USAGE, roadmap, and this ExecPlan match observed commands/results;
- no Phase 2+ domain, business schema, API, UI, Agent, deployment, or Launcher lifecycle is implemented.

## 17. Progress, discoveries and decision log

### Progress

- [x] 2026-07-29 — Read Phase 1 source documents, audited the clean Phase 0 repository, verified the local toolchain, and queried exact selected module versions without mutating the repository.
- [x] 2026-07-29 18:53 +08:00 — Milestone 1: pinned Viper v1.21.0, Zap v1.28.0, and lumberjack v2.2.1; added YAML/JSON, `CORVUS_*`, and explicit override loading; enforced loopback/port/log validation; added OS-derived database/log paths and example YAML; added rotating structured console/file logging. Package tests and full Core tests passed, gofmt reported no files, and official golangci-lint v2.12.2 reported zero issues.
- [x] Milestone 1 — Configuration and logging foundation.
- [x] 2026-07-29 19:05 +08:00 — Milestone 2: pinned modernc SQLite v1.55.0 and goose v3.27.3, pinned sqlc v1.31.1 with a Core-local Go tool directive, added an embedded infrastructure-only schema migration and deterministic generated queries, and implemented WAL/foreign-key/busy-timeout storage open, health, version, close, idempotent migration, and `VACUUM INTO` pre-migration snapshot behavior. sqlc generate/vet and drift check passed; storage and full Core tests passed; official golangci-lint v2.12.2 reported zero issues.
- [x] Milestone 2 — SQLite, goose, and sqlc foundation.
- [x] 2026-07-29 19:16 +08:00 — Milestone 3: pinned Echo v5.3.1 and Cobra v1.10.2; replaced the bootstrap with `corvus serve`; composed config, rotating structured logs, migrated SQLite, readiness HTTP, and graceful context shutdown; added `/healthz`, Vite development proxies, atomic Web staging, tagged `go:embed` SPA delivery, cache rules, and package/integration tests. Default and tagged Core tests passed, official golangci-lint v2.12.2 reported zero issues, frontend format/lint/test/build passed, the embedded binary built, and a real Windows process returned `{"status":"ok","database":"ok","schema_version":1}` while creating its database and log.
- [x] Milestone 3 — Echo, Cobra, lifecycle, and Web delivery.
- [x] 2026-07-29 — Milestone 4 local work: CI now verifies sqlc drift, the tagged embedded Core, and native Core tests; README/USAGE/roadmap match the implemented runtime and Phase boundary. Windows local evidence passed gofmt over 27 files, sqlc generate/vet/drift, all three Go modules, official golangci-lint v2.12.2, frozen pnpm install, frontend format/lint/test/build, tagged Core tests, default/tagged Core builds, Launcher build, live readiness, live embedded index delivery, cache behavior, policy-file audit, and Phase 2 symbol/schema/route audit.
- [x] 2026-07-31 — Remote evidence: GitHub Actions run `30611808221` passed Go quality, Frontend quality, Ubuntu, macOS, and Windows against commit `59a478c`.
- [x] Milestone 4 — CI, documentation, and final audit.

### Surprises & Discoveries

- Phase 0's accepted embedding strategy explicitly assigns the staging utility, embedded filesystem, SPA fallback, cache headers, and Echo handler to Phase 1, even though the short Development Implementation Plan only names backend libraries.
- The maintained latest module versions are newer than the Phase 0 risk snapshot: modernc SQLite is v1.55.0 and goose is v3.27.3; both still require Go versions satisfied by Go 1.26.5.
- Echo v5.3.1 uses the new `*echo.Context` handler API and offers context-driven graceful shutdown through `StartConfig`.
- Initial downloads from the configured `goproxy.cn` failed twice for the large modernc archive (HTTP/2 internal error and unexpected EOF), and the official Go proxy was unreachable over the host's IPv6 route. A direct upstream module download succeeded; smaller transitive modules then resolved through the configured proxy. No proxy setting was committed.
- modernc SQLite v1.55.0 declares `modernc.org/libc` v1.74.3, but Go reports that version retracted by its author for a name-resolution lock leak/deadlock fixed in v1.74.4. The minimal fixed indirect patch was pinned. A later full-graph retraction audit was inconclusive because the configured checksum proxy returned HTTP 504 while checking an unrelated module; selected runtime/tool modules remain explicitly verified.
- Windows process smoke used the built tagged binary rather than `go run`, so the recorded PID was the actual Core process. Package-level lifecycle tests supplied deterministic graceful-cancellation evidence; the smoke process was force-stopped only after its health/database/log observations were captured.
- The first Phase 1 remote run (`30610942309`) exposed a macOS-only `EBADF` from calling `Sync` on `/dev/stderr`; all other jobs passed. Commit `59a478c` made the console sink write-only with a no-op Zap sync while retaining file close errors, and added a regression test that fails if the console's underlying `Sync` is called.

### Decision Log

- **2026-07-29 — Accepted:** bind only to `127.0.0.1:8765` by default and reject non-loopback Phase 1 configuration because no authentication exists.
- **2026-07-29 — Accepted:** reserve `/api/v1` for product APIs and use only `/healthz` for operational readiness.
- **2026-07-29 — Accepted:** use an infrastructure-only runtime metadata migration/query to prove goose/sqlc without anticipating Phase 2 models.
- **2026-07-29 — Accepted:** use a `corvus_webui` build tag so ordinary Go builds do not depend on generated Web assets.
- **2026-07-29 — Accepted:** override only the retracted indirect `modernc.org/libc` v1.74.3 with its author-designated v1.74.4 fix while retaining modernc SQLite v1.55.0.
- **2026-07-29 — Accepted:** create local commits after each validated milestone and do not push without explicit user authorization.
- **2026-07-29 — Accepted:** allow Vite's Core proxy target to be overridden with `CORVUS_CORE_URL`, while retaining `http://127.0.0.1:8765` as the zero-configuration development default.
- **2026-07-29 — Accepted:** keep roadmap and ExecPlan incomplete until this revision has actual three-platform GitHub Actions evidence; Phase 0's earlier successful run is not reusable as Phase 1 evidence.
- **2026-07-31 — Accepted:** do not broadly ignore macOS `EBADF`; prevent terminal syncing at the sink boundary so genuine logger/file close errors remain observable.

### Outcomes & Retrospective

Phase 1 is complete. Core now has a real loopback-only runtime, deterministic configuration, structured rotating logs, pure-Go SQLite with guarded embedded migrations, a sqlc boundary, readiness HTTP, graceful lifecycle handling, and both development-proxy and tagged embedded-Web delivery paths. No Phase 2+ domain or product route was introduced.

Local Windows validation passed the full Go/frontend generation, format, lint, test, build, process-smoke, embedded-Web, policy, and phase-boundary matrix. GitHub Actions run `30611808221` then passed Go quality, Frontend quality, Ubuntu, macOS, and Windows against the macOS-safe logging revision `59a478c`. The repository is ready to plan Phase 2 without carrying an open Phase 1 gate.
