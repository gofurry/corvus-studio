# Corvus Studio Phase 3 — Release + Checklist ExecPlan

This ExecPlan is a living document. Update its Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective sections whenever implementation evidence or assumptions change.

All commands run from `E:\Git\开源\agent\corvus` in PowerShell unless another working directory is stated.

## 1. Purpose and user-visible outcome

Phase 3 establishes a deterministic Steam Coming Soon preparation workflow. A user opens a Project, previews the built-in template, atomically creates one Steam Coming Soon Release Goal and its Checklist, filters tasks, adds user tasks, and transitions task and release states. SQLite remains the source of truth, so Core restarts and direct browser refreshes preserve the workspace.

Observable routes after local implementation:

- `/projects/:projectId` — Project overview with Release workspace entry and progress;
- `/projects/:projectId/release` — template preview, Release Goal creation, state and next actions;
- `/projects/:projectId/checklist` — task filters, selection, details, custom-task creation and state changes.

Phase 3 remains `In progress` until an authorized push produces successful Windows, macOS and Linux GitHub Actions evidence.

## 2. Current repository assessment

At Phase 3 start on 2026-07-31:

- branch `dev` has a clean worktree and is three local commits ahead of `origin/dev`;
- Phase 0, Phase 1 and Phase 2 are recorded as complete;
- the root Go workspace contains `apps/core`, `apps/launcher` and `agent`;
- Core provides configuration, logging, Echo v5, SQLite, goose, sqlc, OpenAPI-generated models, Project application/storage/API behavior and optional embedded Web assets;
- schema version 2 contains runtime metadata and Projects only;
- Web provides Project list/create/detail routes, native directory selection and a generated-client wrapper;
- no Release, Checklist, Template, Resource, Deliverable, Evidence, Agent or Steam integration exists;
- `go test ./apps/core/... ./apps/launcher/... ./agent/...` passed at plan start;
- `pnpm test` passed 12 existing Web tests at plan start;
- latest remote evidence is Actions run `30617549600` for Phase 2 baseline `80dca49`; local directory-picker commits are newer than that run.

## 3. Source-of-truth documents

Conflicts use this precedence:

1. `docs/development/Corvus_Studio_Development_Implementation_Plan_v0.1.md` defines Phase 3 as Release Goal, Steam Coming Soon Template and Checklist with status updates.
2. `docs/architecture/Corvus_Studio_Repository_Structure_Design_v0.1.md` defines Core domain/application/infrastructure/interface boundaries and Web feature placement.
3. `docs/architecture/Corvus_Studio_Technology_Stack_Decision_v0.1.md` fixes Go, Echo v5, SQLite, goose, sqlc, OpenAPI First, React, TypeScript and TanStack Query.
4. `docs/deployment/Corvus_Studio_Engineering_and_Deployment_Guide_v0.1.md` defines generated-code drift checks, CI and three-platform evidence.
5. `USAGE.md` defines developer commands and must contain only verified behavior.
6. `docs/architecture/Corvus_Studio_Data_Model_Design_v0.1.md` defines Release, Checklist and Template concepts and state names.
7. `docs/api/Corvus_Studio_Backend_API_Design_v0.1.md` defines create/get/transition Release and Checklist API conventions.
8. Functional, PRD, UX, frontend, system-design and product-design documents define template ownership, sources and the checklist user flow without moving Phase 4 relationships into Phase 3.
9. `docs/roadmap.md` records evidence-backed progress.

The initial template records current Steam sources without making runtime network calls: Coming Soon, Review Process, Graphical Asset Rules, Store Graphical Assets and Store Page Written Description.

## 4. Scope

Phase 3 will:

- implement Release, Checklist and immutable Template domain boundaries;
- add schema version 3, sqlc queries and transaction-backed repositories;
- embed and validate one `steam-coming-soon` template version `1.0.0` with eight Steam and four Corvus tasks;
- atomically create a Release Goal and independent Checklist copy;
- enforce one goal per Project and goal type;
- persist real Release and Checklist transition history;
- define OpenAPI-first template/release/checklist endpoints and regenerate clients;
- add Project Release and Checklist Web routes, filters and user-task creation;
- verify generation drift, transition rules, rollback, restart persistence and embedded route refresh;
- synchronize README, USAGE, roadmap, CI behavior and this ExecPlan.

## 5. Non-goals

Phase 3 will not implement:

- Resource, Deliverable, Evidence storage or Checklist relationships to those later domains;
- task deletion, content editing, reordering, ownership, due dates or assignment;
- template import, export, online refresh, upgrade diff or merge;
- Steam authentication, Steam API access, automated review submission or publishing;
- a Published state not present in the higher-priority data-model document;
- Home Dashboard, reports, SSE, React Flow, Zustand, Agent/ADK or model providers;
- Project editing, Launcher business UI, Docker, installers, signing or auto-update.

## 6. Document conflicts and decisions

| Status   | Conflict or ambiguity                                                                                       | Decision and rationale                                                                                                                        |
| -------- | ----------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------- |
| Accepted | Backend API design omits Release list and template-read routes needed to restore and preview the Web state. | Add read-only `GET /projects/{project_id}/releases` and `GET /release-templates/{template_key}` without rewriting versioned design documents. |
| Accepted | UX lists `Published`; the higher-priority data model stops at `Submitted`.                                  | Phase 3 uses Draft, Preparing, NeedsAttention, ReadyForReview, Ready and Submitted only.                                                      |
| Accepted | Product documents allow multiple template sources while the implementation plan names one Steam template.   | One versioned pack contains eight `platform_template` tasks and four `corvus_template` tasks.                                                 |
| Accepted | Template-generated tasks are user-owned, while Phase 3 only requires status editing.                        | Persist an independent copy and allow status changes plus new User tasks; defer content editing and deletion.                                 |
| Accepted | Future releases may target other goals, while v0.1 has only Steam Coming Soon.                              | Schema supports goal types, but `UNIQUE(project_id, goal_type)` permits only one Steam Coming Soon goal per Project.                          |
| Accepted | Goal creation and template generation could be separate operations.                                         | `POST /releases` creates the Goal and all template items in one transaction so partial workspaces cannot exist.                               |
| Accepted | Required tasks could be marked NotApplicable.                                                               | Required tasks must be Done before ReadyForReview and cannot become NotApplicable; recommended tasks may be NotApplicable.                    |
| Accepted | Remote policy was not expanded by the implementation request.                                               | Create local commits only. Keep Phase 3 In progress until a later authorized push and successful native CI run.                               |

## 7. Proposed repository tree

Relevant additions after Phase 3:

```text
.agent/phase-3-release-checklist.md
apps/core/
├── internal/
│   ├── domain/{release,checklist,template}/
│   ├── application/{release,checklist}/
│   ├── infrastructure/
│   │   ├── storage/{release_repository,checklist_repository}.go
│   │   └── templatecatalog/
│   │       └── templates/steam-coming-soon-v1.0.0.json
│   └── interfaces/httpserver/{releases,checklist,templates}.go
├── migrations/00003_release_checklist.sql
├── queries/{releases,checklist}.sql
└── openapi/openapi.yaml
apps/web/src/
├── api/{releases,checklist,templates}.ts
└── features/{release,checklist}/
packages/api-client/src/generated/
```

No Phase 4 directory or relationship table is created.

## 8. Milestones

### Milestone 1 — Start the living plan

Create this ExecPlan, mark Phase 3 In progress, validate Markdown/diff, and commit `docs: start Phase 3 release checklist`.

### Milestone 2 — Define the public contract

Extend OpenAPI with template, Release and Checklist schemas/routes; regenerate committed Go and TypeScript outputs twice; verify compilation and zero drift; commit `build(api): define release and checklist contracts`.

### Milestone 3 — Persist the workspace

Implement domains, transition rules, embedded template Catalog, migration 00003, sqlc queries, transaction repositories, histories and restart/rollback tests; commit `feat(core): persist release checklist workspace`.

### Milestone 4 — Expose the HTTP workflow

Compose services in Runtime, register routes before SPA fallback, add strict bounded request decoding, DTO mapping, filters, errors and OpenAPI contract tests; commit `feat(api): expose release checklist workflow`.

### Milestone 5 — Build the Web workflow

Add Project workspace navigation, template preview, atomic creation, Release status controls, Checklist filtering/detail/status and custom tasks; add interaction tests; commit `feat(web): add release checklist workspace`.

### Milestone 6 — Audit local completion

Synchronize README, USAGE, roadmap, CI if required, and this living plan; run section 16; commit `docs: validate Phase 3 locally`. Leave remote Actions as the explicit remaining gate.

Before every commit inspect `git status --short`, unstaged and staged diffs, generated/local artifacts and obvious secrets; stage explicit paths only. A failed validation blocks its commit. Recover a completed milestone with `git revert <hash>`.

## 9. Go workspace strategy

All new code remains in existing `apps/core`; no Go module or workspace change is required. Launcher and Agent receive no Phase 3 behavior. Template JSON uses the standard library and `go:embed`, so no runtime dependency is added.

Canonical tests remain:

```powershell
go test ./apps/core/... ./apps/launcher/... ./agent/...
```

## 10. pnpm workspace strategy

The root remains the only lockfile owner. `packages/api-client` is regenerated from Core OpenAPI; `apps/web` consumes it through the workspace. No new frontend runtime dependency is planned. `packages/shared-types` and `packages/ui` remain placeholders; Zustand, React Flow, react-markdown and Playwright remain absent.

## 11. Build and embedding strategy

Development remains separate Core and Vite processes. Production continues staging `apps/web/dist` and compiling Core with `corvus_webui`. All new `/api/v1` routes register before the SPA fallback, and embedded tests cover direct refresh of Release and Checklist URLs. Launcher remains unchanged.

## 12. CI strategy

Existing Go/sqlc/OpenAPI/TypeScript drift, test, lint, frontend and native-platform jobs should cover Phase 3 after their test suites expand. Modify `.github/workflows/ci.yml` only if a new explicit command is required. Do not add Issue/PR templates or another CI system. Remote completion requires successful Ubuntu, Windows and macOS jobs for the Phase 3 commits.

## 13. Validation matrix

| Target      | Required evidence                                                                     |
| ----------- | ------------------------------------------------------------------------------------- |
| Domain      | UUIDv7, validation, full transition matrices, idempotent same-state behavior          |
| Template    | exact 12 items, stable keys/order/version, two sources, valid references              |
| SQLite      | migration 2→3, backup, atomic generation, conflict rollback, histories, reopen        |
| API         | 201/200/400/404/409/500, strict JSON, UUIDs, filters and OpenAPI compliance           |
| Web         | no-goal state, preview, create, refresh, filters, custom task, transitions and errors |
| Consistency | Required completion gates ReadyForReview; advanced goal states lock Required tasks    |
| Generation  | OpenAPI Go models, sqlc and TypeScript client regenerate without drift                |
| Embedding   | Release and Checklist browser refresh returns the SPA while APIs remain authoritative |
| Boundary    | no Resource, Deliverable, Evidence, Steam integration or Agent implementation         |

## 14. Idempotence and recovery

- Template generation uses exact key/version and a unique Project/goal-type constraint.
- Repeated Release creation returns conflict without additional items.
- Same-state transition requests return current state without a history row.
- Transaction failures roll back Goal, Checklist and history writes together.
- Generated files are changed only through pinned generators and must reproduce exactly.
- Existing Checklist copies are never overwritten by later Catalog versions.
- Migration 00003 is corrected in place only before its milestone commit; later fixes use a new migration.
- Tests use temporary data directories; never delete an unresolved or user Project path.
- Recovery uses focused fixes or `git revert`, never reset-hard, clean, amend or force push.

## 15. Risks and blockers

- Current official Steam rules may change; the template records version/review date/source links and never claims live synchronization.
- Phase 3 intentionally cannot validate images, trailers or external Steam state because Resource/Deliverable and Steam integration are later work.
- The local PATH may expose an obsolete golangci-lint; only pinned v2 evidence is acceptable.
- Existing API code uses manual Echo v5 handlers because the generator adapter targets Echo v4.
- A single large Web chunk already produces a non-failing Vite warning; Phase 3 should avoid unrelated bundling optimization.
- Latest local directory-picker commits have not yet received remote native-platform evidence and will be included in the next authorized push.

No implementation blocker is known at plan start.

## 16. Completion criteria

Run from the repository root:

```powershell
Push-Location .\apps\core
go tool oapi-codegen -config openapi/oapi-codegen.yaml openapi/openapi.yaml
go tool sqlc generate
go tool sqlc vet
Pop-Location

pnpm --filter @corvus-studio/api-client generate
git diff --exit-code -- apps/core/internal/interfaces/httpserver/api apps/core/internal/infrastructure/storage/sqlc packages/api-client/src/generated

go test ./apps/core/... ./apps/launcher/... ./agent/...
go vet ./apps/core/... ./apps/launcher/... ./agent/...
golangci-lint run ./apps/core/... ./apps/launcher/... ./agent/...

pnpm install --frozen-lockfile
pnpm format:check
pnpm lint
pnpm test
pnpm build

go run ./apps/core/cmd/stage-web
go test -tags corvus_webui ./apps/core/...
go build -tags corvus_webui -o .tmp/corvus-phase3 ./apps/core/cmd/corvus
git diff --check
git status --short
```

Evidence must show atomic creation, duplicate rejection, persistence after Core reopen, direct embedded-route refresh, zero generated drift, no Phase 4/Agent behavior, a clean final worktree and one local commit per milestone. Local success leaves Phase 3 In progress until remote three-platform CI passes.

## 17. Progress, discoveries and decision log

### Progress

- [x] 2026-07-31: audited the clean `dev` baseline, relevant design documents, existing architecture and current official Steam documentation; obtained product decisions.
- [x] 2026-07-31: Milestone 1 created the living ExecPlan and marked Phase 3 In progress; Markdown formatting and `git diff --check` passed.
- [x] 2026-07-31: Milestone 2 added the template/Release/Checklist OpenAPI contract and regenerated Go/TypeScript clients; Core tests, API-client typecheck and consecutive-output hash comparison passed.
- [x] 2026-07-31: Milestone 3 implemented Release/Checklist/Template domains, the 12-item embedded Catalog, schema version 3, sqlc repositories, atomic workspace creation, transition histories and restart/rollback tests; sqlc generate/vet, Core tests and Go vet passed.
- [x] 2026-07-31: Milestone 4 composed the services into Runtime and exposed template/Release/Checklist Echo v5 routes with strict requests, filters, DTO/error mapping and OpenAPI exchange tests; focused uncached tests, full Core tests and Go vet passed, including a real restart-persistence HTTP test.
- [ ] Milestone 5 — Web workflow.
- [ ] Milestone 6 — full local completion audit.
- [ ] Authorized push and successful Windows/macOS/Linux Actions evidence.

### Surprises & Discoveries

- Backend API design has no Release-list or template-preview route, although a refreshable Web workflow requires both.
- UX includes Published while the higher-priority data model does not.
- The branch begins three local commits ahead of its remote baseline.
- On Windows, one immediate repeated TypeScript generation attempt hit a transient file lock while Prettier reopened `client/types.gen.ts`. A clean retry and two successful runs separated by two seconds produced the identical diff hash `848a1bfe7b7c864351ead23bb66061324084d999`; no source or generator change was required.
- Expanding the sqlc schema caused the pinned sqlc tool to record additional transitive module checksums in `go.work.sum`; no Go module version or direct dependency changed.

### Decision Log

- **2026-07-31 — Accepted:** one Steam Coming Soon goal per Project and goal type.
- **2026-07-31 — Accepted:** atomically create Goal plus Checklist from exact template version.
- **2026-07-31 — Accepted:** include Steam official and Corvus recommended tasks in one pack.
- **2026-07-31 — Accepted:** allow User tasks but no task deletion.
- **2026-07-31 — Accepted:** explicit Release transitions with Required-task readiness gates.
- **2026-07-31 — Accepted:** local milestone commits only; no push in this implementation run.

### Outcomes & Retrospective

Not started — implementation milestones are pending.
