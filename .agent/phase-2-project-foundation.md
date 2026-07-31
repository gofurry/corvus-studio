# Corvus Studio Phase 2 — Project Foundation ExecPlan

This ExecPlan is a living document. Update its Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective sections whenever implementation evidence or assumptions change.

All commands run from `E:\Git\开源\agent\corvus` in PowerShell unless another working directory is stated.

## 1. Purpose and user-visible outcome

Phase 2 creates the first end-to-end Corvus Studio product workflow. After completion, a developer can start Core and Web, create a Project that references an existing local directory, view all Projects, and open a Project detail page. Closing and restarting Core or refreshing the browser must not lose the Project because SQLite is the source of truth.

The observable Web routes are:

- `/projects` — Project list and empty state;
- `/projects/new` — Project creation form;
- `/projects/:projectId` — Project detail loaded from Core.

“Open Project” means navigating to the in-app detail route. Phase 2 does not open the operating-system file manager and does not create, move, or write files in the referenced Project directory.

## 2. Current repository assessment

At Phase 2 start on 2026-07-31:

- branch `dev` points to `74f59de` and matches `origin/dev`;
- the worktree is clean;
- Phase 0 and Phase 1 are complete;
- the root Go workspace contains `apps/core`, `apps/launcher`, and `agent`;
- Core already provides Cobra/Viper configuration, Zap logging, Echo v5, SQLite through `modernc.org/sqlite`, embedded goose migrations, sqlc generation, `/healthz`, graceful shutdown, and optional embedded Web assets;
- the database has only schema version 1 and the infrastructure-only `corvus_runtime_metadata` table;
- no `internal/domain` Project code, Project migration, `/api/v1` business route, OpenAPI source, or generated API client exists;
- `packages/api-client`, `packages/shared-types`, and `packages/ui` are explanatory placeholders;
- Web is still the Phase 0 bootstrap page and depends only on React/ReactDOM plus its development toolchain;
- GitHub Actions already validates Go formatting/sqlc/test/lint, frontend format/lint/test/build, embedded Core build, and native Windows/macOS/Linux Core/Launcher builds.

Observed local tools:

| Tool | Version/evidence |
| --- | --- |
| Go | `go1.26.5 windows/amd64` |
| Go workspace | root `go.work` |
| Node | `v24.15.0` |
| pnpm | `10.11.0` |
| Git | `2.51.0.windows.1` |
| golangci-lint | local PATH is `v1.64.8`; CI pins `v2.12.2` |

The outdated local golangci-lint binary is not accepted as Phase 2 lint evidence. Use the same pinned v2 binary strategy recorded in the Phase 1 ExecPlan or rely on a later authorized CI run.

## 3. Source-of-truth documents

Conflicts use this precedence:

1. `docs/development/Corvus_Studio_Development_Implementation_Plan_v0.1.md` defines Phase 2 as Project Domain, Project API, Project page, and local storage.
2. `docs/architecture/Corvus_Studio_Repository_Structure_Design_v0.1.md` defines Core clean-architecture layers, Web feature placement, OpenAPI ownership, and generated client placement.
3. `docs/architecture/Corvus_Studio_Technology_Stack_Decision_v0.1.md` fixes OpenAPI First, UUIDv7, SQLite, React 19, TypeScript, Vite 8, pnpm, Ant Design, TanStack Query, and React Router.
4. `docs/deployment/Corvus_Studio_Engineering_and_Deployment_Guide_v0.1.md` defines CI, generated-code checks, embedded production assets, and three-platform support.
5. `USAGE.md` defines developer-facing commands and must be synchronized only with verified Phase 2 behavior.
6. `docs/architecture/Corvus_Studio_Data_Model_Design_v0.1.md` defines the conceptual Project fields and Active/Archived status.
7. `docs/api/Corvus_Studio_Backend_API_Design_v0.1.md` defines `/api/v1`, create/get Project routes, and the error envelope.
8. `docs/architecture/Corvus_Studio_Launch_System_Design_Document_v0.1.md` defines local-first storage, Project folders as first-class references, and Core as the single writer.
9. UX, frontend, functional, PRD, product, and MVP roadmap documents provide user-flow context without moving Release Goal or later domains into Phase 2.
10. `docs/roadmap.md` records implementation progress and requires create/list/open plus persistence evidence.

## 4. Scope

Phase 2 will:

- create a Project domain with UUIDv7 identity, four lifecycle stages, Active/Archived status, validation, and explicit domain errors;
- add Project application service ports for repository, clock, and ID generation;
- add a schema version 2 Project migration and sqlc create/get/list queries;
- store normalized references to existing local directories without modifying those directories;
- define OpenAPI 3.0.3 create/list/get Project endpoints and shared error schemas;
- generate committed Go API models and a committed TypeScript Fetch client from the same specification;
- expose manual Echo v5 handlers for create/list/get;
- build the minimum React routes for listing, creating, and opening Projects;
- add domain, storage, HTTP contract, Web interaction, persistence, generation-drift, and native-platform validation;
- update README, USAGE, CI, roadmap, and this ExecPlan with current evidence.

## 5. Non-goals

Phase 2 will not implement:

- Project edit, delete, archive, restore, search, pagination, import, export, or operating-system folder opening;
- Project dashboard data;
- Release Goal, Checklist, Steam templates, Resource, Deliverable, Asset Map, Watch, reports, or review workflows;
- Agent/ADK, model providers, tools, prompts, SSE, or AI behavior;
- login, authentication, remote network exposure, or credential storage;
- Launcher lifecycle integration or desktop business UI;
- Docker, systemd, installers, signing, notarization, automatic update, or formal release packaging;
- a shared UI component library or global Zustand store without a Phase 2 use case.

## 6. Document conflicts and decisions

| Status | Conflict or ambiguity | Decision and rationale |
| --- | --- | --- |
| Accepted | The Backend API design specifies create/get but omits Project list; the maintained roadmap requires create/list/open. | Add `GET /api/v1/projects`. The higher-priority implementation plan and roadmap acceptance require a list to reopen an existing Project. Record this reconciliation here; do not rewrite the versioned design document. |
| Accepted | Project `stage` exists but allowed values are undefined. | Use `concept`, `development`, `release_preparation`, and `released` in OpenAPI, domain, SQLite, and UI. |
| Accepted | Project location may be a Corvus-managed folder or external workspace. | Phase 2 accepts only an existing absolute directory, resolves and stores its canonical reference, and never writes to it. Managed directories are deferred. |
| Accepted | “Open Project” could mean in-app selection or OS file-manager launch. | Navigate to `/projects/:projectId` and reload from Core. OS integration is deferred. |
| Accepted | The API design includes `/projects/{id}/dashboard`, whose data depends on later domains. | Do not add the dashboard route in Phase 2. |
| Accepted | `oapi-codegen` can generate an Echo server adapter, but its documented adapter imports Echo v4 while the repository uses Echo v5. | Pin `oapi-codegen v2.8.0`, generate models only, and register manual Echo v5 routes. |
| Accepted | The repository policy requires actual three-platform evidence, but this task explicitly selects local commits only. | Complete all local checks, keep Phase 2 In progress, and leave the remote CI acceptance item open until push is separately authorized. |

## 7. Proposed repository tree

Relevant additions and modifications after Phase 2:

```text
.agent/
└── phase-2-project-foundation.md
apps/
├── core/
│   ├── internal/
│   │   ├── application/
│   │   │   └── project/
│   │   ├── domain/
│   │   │   └── project/
│   │   ├── infrastructure/
│   │   │   └── storage/
│   │   │       ├── project_repository.go
│   │   │       └── sqlc/                  # generated, committed
│   │   └── interfaces/
│   │       └── httpserver/
│   │           ├── api/
│   │           │   └── models.gen.go      # generated, committed
│   │           └── projects.go
│   ├── migrations/
│   │   └── 00002_projects.sql
│   ├── openapi/
│   │   ├── openapi.yaml
│   │   └── oapi-codegen.yaml
│   └── queries/
│       └── projects.sql
└── web/
    └── src/
        ├── api/
        ├── app/
        └── features/
            └── projects/
packages/
└── api-client/
    ├── src/
    │   ├── generated/                     # generated, committed
    │   └── index.ts
    ├── openapi-ts.config.ts
    ├── package.json
    └── tsconfig.json
```

`packages/shared-types` and `packages/ui` remain placeholders. Project API types come from OpenAPI and must not be copied into `shared-types`.

## 8. Milestones

### Milestone 1 — Start the living plan

**Goal:** record the decision-complete scope and show Phase 2 as active without claiming implementation.

**Paths:** `.agent/phase-2-project-foundation.md`, `docs/roadmap.md`.

**Operations:** create this ExecPlan; update the roadmap summary/table/Phase 2 section to In progress and link this plan.

**Validation:**

```powershell
git diff --check
git status --short
```

Expected: only the ExecPlan and roadmap are changed; no product source exists.

**Recovery:** correct the two planning files or revert the milestone commit. Do not reset or modify Phase 1 history.

### Milestone 2 — OpenAPI contract and deterministic generators

**Goal:** make one OpenAPI source generate compile-checked Go models and a reusable TypeScript Fetch client.

**Paths:** Core OpenAPI files and module files, `packages/api-client/`, root/frontend manifests and lockfile.

**Operations:** define create/list/get and error schemas; pin `oapi-codegen v2.8.0` as a Core Go tool; generate models only; initialize the API client package with exact `openapi-ts` tooling; add root `generate:api`; add exact Web dependencies.

**Validation:**

```powershell
Push-Location .\apps\core
go tool oapi-codegen -config openapi/oapi-codegen.yaml openapi/openapi.yaml
Pop-Location
pnpm --filter @corvus-studio/api-client generate
pnpm --filter @corvus-studio/api-client build
go test ./apps/core/...
pnpm install --frozen-lockfile
```

Expected: Go models and TypeScript client compile; a second generation produces no tracked diff.

**Recovery:** repair the source spec/config and regenerate. Never hand-edit generated files or replace confirmed versions silently.

### Milestone 3 — Project domain and persistence

**Goal:** create, list, fetch, close, and reopen Projects through domain/application/storage layers.

**Paths:** Core domain/application/storage packages, migration 00002, Project queries, sqlc output.

**Operations:** add domain types/errors/validation; add repository/clock/ID ports and service; create STRICT schema with a unique normalized location key; generate sqlc; implement the SQLite adapter; add restart and failure tests.

**Validation:**

```powershell
Push-Location .\apps\core
go tool sqlc generate
go tool sqlc vet
Pop-Location
go test ./apps/core/internal/domain/project/... ./apps/core/internal/application/project/... ./apps/core/internal/infrastructure/storage/...
```

Expected: migration reaches version 2; tests prove UUIDv7, validation, duplicate rejection, stable ordering, and persistence after reopen.

**Recovery:** fix an uncommitted migration in place. After the migration is committed, correct schema changes with a new migration rather than rewriting history.

### Milestone 4 — Echo v5 Project API

**Goal:** expose the OpenAPI behavior through the existing loopback Core runtime.

**Paths:** Core HTTP interface and runtime composition packages plus tests.

**Operations:** add dependency-based server construction; register routes before SPA fallback; use strict JSON decoding and a 1 MiB limit; map DTO/domain/storage errors; add OpenAPI-backed request/response tests.

**Validation:**

```powershell
go test ./apps/core/internal/interfaces/httpserver/... ./apps/core/internal/application/...
go test ./apps/core/...
```

Expected: create returns 201 and Location; list/get return 200; invalid/not-found/conflict/internal cases use the required envelope; all tested traffic conforms to the specification.

**Recovery:** retain the OpenAPI contract and correct only the handler/mapper/service layer. Do not loosen the contract to hide implementation failures.

### Milestone 5 — Minimum Project Web flow

**Goal:** let a user list, create, and open Projects through accessible React routes.

**Paths:** `apps/web/src/` and Web manifest.

**Operations:** add Router, QueryClient, Ant Design shell, API wrapper, Project pages/components/styles, and interaction tests.

**Validation:**

```powershell
pnpm format:check
pnpm lint
pnpm test
pnpm build
```

Expected: all checks pass; tests cover loading/empty/error/list/create/detail behavior; build emits `apps/web/dist/index.html`.

**Recovery:** fix the failing feature slice without changing generated client files. Restore only through a new fix or `git revert` of the milestone commit.

### Milestone 6 — CI, docs, and local completion audit

**Goal:** align automation and developer documentation, collect full local evidence, and leave remote CI clearly pending.

**Paths:** CI workflow, README, USAGE, roadmap, and this ExecPlan.

**Operations:** add OpenAPI/client drift checks; retain native Project tests; document commands and user flow; update all living sections and check only evidence-backed roadmap items.

**Validation:** run the complete command set in section 16.

Expected: local checks pass, worktree is clean after the final commit, and the roadmap remains In progress solely because remote CI was not authorized.

**Recovery:** create a follow-up fix commit for the failing subsystem; never amend or weaken validation to manufacture completion.

## 9. Go workspace strategy

The root `go.work` remains unchanged. All Project code and its dependencies belong to the existing `apps/core` module. Launcher and Agent receive no new dependency or product behavior.

Core adds direct dependencies/tool pins only when imported:

- `github.com/google/uuid v1.6.0`;
- `github.com/getkin/kin-openapi v0.145.0` for contract tests;
- `github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen v2.8.0` as a Go tool.
- `github.com/oapi-codegen/runtime v1.6.0` for generated UUID transport models.

The canonical workspace test remains:

```powershell
go test ./apps/core/... ./apps/launcher/... ./agent/...
```

## 10. pnpm workspace strategy

The root remains the workspace coordinator and owns the only `pnpm-lock.yaml`.

- `packages/api-client` becomes `@corvus-studio/api-client` and contains committed generated Fetch SDK/types. `openapi-ts 0.99.0` emits the Fetch implementation into the generated package, so the deprecated standalone `@hey-api/client-fetch` package is not retained.
- `apps/web` consumes the package through the pnpm workspace and owns Ant Design, TanStack Query, and React Router runtime dependencies.
- root `package.json` gains only orchestration such as `generate:api`; frontend runtime dependencies do not move to root.
- `packages/shared-types` and `packages/ui` remain documentation-only placeholders.
- Zustand, React Flow, and react-markdown remain uninstalled until their owning phases need them.

No subdirectory may run an independent install or create another lockfile.

## 11. Build and embedding strategy

Development remains two processes:

```powershell
go run ./apps/core/cmd/corvus serve
pnpm dev:web
```

The generated Fetch client uses same-origin `/api/v1`; Vite continues proxying `/api` to loopback Core. Production validation builds Web, stages `apps/web/dist`, and compiles Core with `corvus_webui`. API and health routes must take precedence over the SPA fallback, including direct browser refresh on `/projects/:id`.

Launcher remains unchanged and does not start Core or host Project UI.

## 12. CI strategy

GitHub Actions remains the only CI.

Required Phase 2 additions:

- Go quality regenerates OpenAPI Go models and sqlc output, then fails on tracked drift;
- frontend quality regenerates the TypeScript client and fails on tracked drift before format/lint/test/build;
- native Windows/macOS/Linux jobs continue `go test ./apps/core/...`, which now includes Project migration/repository/API persistence tests;
- existing dependency caches continue to key on all Go sums and the single pnpm lockfile;
- no Issue/PR template or `CONTRIBUTING.md` is added.

Because this task is local-only, the workflow definition may be validated locally, but actual native job success remains pending.

## 13. Validation matrix

| Target | Command/observation | Required evidence |
| --- | --- | --- |
| Domain | domain package tests | stages, lengths, IDs, UTC, errors |
| Project path | domain/service tests | absolute existing directory, no writes, canonical uniqueness |
| SQLite migration | storage tests | version 2, idempotent reopen, backup path for upgrade |
| Repository | storage tests | create/get/list/conflict and deterministic order |
| OpenAPI Go models | regenerate plus Git diff | committed deterministic output |
| TypeScript client | regenerate/build plus Git diff | committed deterministic output |
| HTTP API | handler and contract tests | 201/200/400/404/409/500 envelopes |
| Persistence | close/reopen integration test | created Project remains readable |
| Web | Vitest interaction tests | list, empty, create, navigate, detail, errors |
| Web build | `pnpm build` | `apps/web/dist/index.html` exists |
| Embedded Core | stage plus tagged test/build | direct Project routes and SPA refresh coexist |
| Go workspace | explicit three-module test | exit 0 |
| Lint | pinned golangci-lint v2 | exit 0 |
| Native CI | actual Actions run | pending until push is authorized |
| Phase boundary | source/dependency scan | no Phase 3+ implementation |

## 14. Idempotence and recovery

- Inspect `git status`, unstaged diff, and staged diff before every change and commit.
- Generate Go/TypeScript/sqlc output only from pinned tools and committed inputs.
- Re-running generation must produce no diff; never edit generated files manually.
- `pnpm install --frozen-lockfile` is the verification path after the intentional dependency/lockfile commit.
- Migration 00002 uses guarded create/drop statements and is tested through repeated open. Once committed, use a new migration for corrections.
- Project creation never changes the referenced directory, so rollback only removes the database row created by the test fixture.
- Use temp directories outside tracked source for runtime tests and never recursively delete unresolved paths.
- Stage explicit paths and revert completed milestones with `git revert <commit>`. Do not use reset-hard, clean, amend, rebase, or force push.

## 15. Risks and blockers

- The local golangci-lint is v1 and cannot validate the repository’s v2 configuration; use the previously recorded pinned v2 binary path or leave that single check pending.
- `oapi-codegen` Echo server output targets Echo v4, so generating handlers would create an incompatible framework dependency; models-only generation avoids this.
- Hey API is pre-1.0 and may change generated output; exact pinning and committed drift checks are mandatory.
- Filesystem canonicalization differs by platform. Windows normalizes the uniqueness key case-insensitively; Unix-like systems retain case.
- A referenced directory may disappear after Project creation. Phase 2 preserves the reference and does not add watch/recovery UI.
- Core is still loopback-only without authentication. Do not expand the bind boundary.
- Actual Windows/macOS/Linux evidence cannot be renewed without an authorized push/PR.
- macOS signing, notarization, and installers are unrelated to Phase 2.

No blocking issue is known at plan start.

## 16. Completion criteria

Local implementation is complete only when all applicable commands pass:

```powershell
Push-Location .\apps\core
go tool oapi-codegen -config openapi/oapi-codegen.yaml openapi/openapi.yaml
go tool sqlc generate
go tool sqlc vet
Pop-Location

pnpm --filter @corvus-studio/api-client generate
git diff --exit-code -- apps/core/internal/interfaces/httpserver/api apps/core/internal/infrastructure/storage/sqlc packages/api-client/src/generated

go test ./apps/core/... ./apps/launcher/... ./agent/...
golangci-lint run ./apps/core/... ./apps/launcher/... ./agent/...

pnpm install --frozen-lockfile
pnpm format:check
pnpm lint
pnpm test
pnpm build

go run ./apps/core/cmd/stage-web
go test -tags corvus_webui ./apps/core/...
go build -tags corvus_webui -o .tmp/corvus-phase2 ./apps/core/cmd/corvus
git diff --check
git status --short
```

Evidence must also show:

- a Project survives closing and reopening Core storage;
- browser refresh on the detail route reloads the Project from Core;
- invalid input and duplicate location produce no partial row and no directory modification;
- OpenAPI, Go models, HTTP behavior, and TypeScript client agree;
- README and USAGE commands match the implementation;
- no later-phase domain or endpoint has been added;
- every logical milestone has a local commit and rollback path.

Phase 2 remains `In progress` until an actual authorized GitHub Actions run passes Go quality, frontend quality, and native Windows/macOS/Linux jobs.

## 17. Progress, discoveries and decision log

### Progress

- [x] 2026-07-31: audited the clean `dev` baseline at `74f59de`, read Phase 2 source documents, and obtained the three product decisions plus local-only remote policy.
- [x] 2026-07-31: Milestone 1 prepared the living plan and roadmap In progress state; `git diff --check` is the pre-commit validation.
- [x] 2026-07-31: Milestone 2 defined OpenAPI create/list/get, generated deterministic Go and TypeScript clients, and passed `pnpm generate:api`, frozen install, API-client typecheck, root format/build, and `go test ./apps/core/...`.
- [x] 2026-07-31: Milestone 3 implemented Project domain/application/storage, migration 00002, canonical existing-directory references, and restart persistence; sqlc generate/vet and focused/full Core tests passed.
- [x] 2026-07-31: Milestone 4 implemented manual Echo v5 create/list/get handlers, strict bounded JSON, error envelopes, runtime composition, OpenAPI request/response tests, and a real HTTP restart persistence test; focused/full Core tests and Go vet passed.
- [x] 2026-07-31: Milestone 5 implemented Project list/create/detail routes, generated-client wrapper, TanStack Query cache behavior, Ant Design UI, and interaction tests; format/lint, 7 Vitest cases, API-client typecheck, and Vite build passed.
- [ ] Milestone 6: complete local audit, docs, CI definition, and final local commit.
- [ ] Remote completion: push only after explicit authorization and obtain a green native CI run.

### Surprises & Discoveries

- The Backend API design lacks the Project list route required by the implementation roadmap.
- The selected Go generator’s Echo adapter imports Echo v4, while Corvus Studio is already on Echo v5.
- The local PATH still exposes golangci-lint v1.64.8 even though CI and the repository configuration use v2.12.2.
- `@hey-api/client-fetch 0.13.1` installed with an upstream deprecation warning, and the `openapi-ts 0.99.0` output imports only its own generated Fetch implementation. The unnecessary standalone dependency was removed rather than committing deprecated code.
- Generated Go UUID fields require `github.com/oapi-codegen/runtime/types`; the compatible current runtime `v1.6.0` was added explicitly.
- The first Ant Design production bundle is about 893 kB minified (287 kB gzip) and triggers Vite’s non-failing 500 kB chunk warning. Phase 2 keeps the simple route structure; route-level optimization is a later performance task, not a correctness blocker.

### Decision Log

- **2026-07-31 — Accepted:** add `GET /api/v1/projects` to complete create/list/open, without editing the versioned API design document.
- **2026-07-31 — Accepted:** use `concept`, `development`, `release_preparation`, and `released` as the Phase 2 stage enum.
- **2026-07-31 — Accepted:** Project creation references an existing absolute directory and never writes to it.
- **2026-07-31 — Accepted:** opening a Project means navigating to an application detail route.
- **2026-07-31 — Accepted:** generate Go models only and manually adapt Echo v5.
- **2026-07-31 — Accepted deviation:** keep the pinned `openapi-ts 0.99.0` Fetch plugin but omit deprecated `@hey-api/client-fetch 0.13.1`; generated output is self-contained and compile-checked.
- **2026-07-31 — Accepted:** create local milestone commits but do not push; retain remote CI as the final open gate.

### Outcomes & Retrospective

Milestones 1–5 established the living plan, deterministic OpenAPI boundary, Project persistence/API, and the minimum user-facing Web flow. Automated evidence now covers create/list/detail, direct detail-route reload, errors, and create-to-detail navigation.
