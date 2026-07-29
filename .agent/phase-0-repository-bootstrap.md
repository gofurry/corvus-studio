# Phase 0: Repository Bootstrap ExecPlan

> Status: Execution in progress — Milestones 2–3 complete; Milestone 4 started; Milestone 5 pending.
> Target repository root: `E:\Git\开源\agent\corvus`
> Prepared: 2026-07-29 (Asia/Shanghai)
> Maintenance standard: `.agent/PLANS.md`

## 1. Purpose and user-visible outcome

Phase 0 establishes the Corvus Studio repository and build foundation without implementing product behavior.

After completion, a developer can clone or open the repository root and:

- inspect a recognizable Monorepo containing Core, Web, Launcher, Agent, shared-package reservations, documentation, tests, tools, deployment reservations, and GitHub Actions;
- use one `go.work` file to discover and test the Core, Launcher, and Agent Go modules;
- build a minimal Core bootstrap binary;
- build and manually open a minimal Fyne Launcher window that contains no business UI and does not start Core;
- install one pnpm workspace lockfile, run the React development server, and observe a static “Corvus Studio — Repository bootstrap” page;
- run frontend lint, smoke tests, and production build from the repository root;
- compare the same local commands with the commands used by GitHub Actions;
- see accurate README and USAGE instructions that distinguish working Phase 0 commands from future commands.

Observable success is defined by commands and files in Sections 13 and 16. Phase 0 does not make `corvus serve`, an API, a database, an Agent, or any product workflow available.

## 2. Current repository assessment

Assessment date: 2026-07-29. The initial checks were read-only with respect to the repository. A later check observed that the approved root relocation had occurred between assistant turns; that external state change was not performed by this planning task.

### Repository location and files

The approved and current Git root is:

    E:\Git\开源\agent\corvus

It currently contains:

- `.git/`;
- `LICENSE`, containing the GNU Affero General Public License v3 text;
- an empty, untracked `README.md`;
- an untracked `USAGE.md`;
- an untracked `docs/` directory containing all 15 supplied design documents;
- the allowed `AGENTS.md` and `.agent/` planning artifacts.

The former nested `corvus-studio/` directory is absent. Because the relocation occurred outside the assistant's recorded milestone procedure, a pre-move SHA-256 manifest for all 15 documents is unavailable. Earlier evidence still proves that the four moved tracked documents are byte-identical to their HEAD blobs, and current inspection confirms all 15 expected document names remain present.

No application directories, Go modules, Go workspace, pnpm workspace, package manifests, CI workflows, or build configuration currently exist. `AGENTS.md` and `.agent/` were absent before this planning task and are the only planning artifacts added by it.

### Git state

- The approved outer root is a Git repository on `main`, tracking `origin/main` with `+0/-0` branch divergence.
- Remote: `https://github.com/gofurry/corvus-studio.git`.
- HEAD at assessment: `3bac4c69631db6bdbc96afe212b8608b2a0f36a3` (`add: system design docs`).
- The worktree is not clean and must not be cleaned or reset by Phase 0.
- Four documents tracked at the repository root appear deleted, while byte-identical blobs exist under the untracked `docs/` directory:
  - `Corvus_Studio_Launch_Functional_Specification_v0.1.md`;
  - `Corvus_Studio_Launch_PRD_v0.1.md`;
  - `Corvus_Studio_Launch_UX_IA_v0.1.md`;
  - `Corvus_Studio_Product_Design_Draft_v0.1.md`.
- `README.md`, `USAGE.md`, and `docs/` are untracked.
- Hash comparison confirmed the four moved document blobs match their HEAD versions exactly.

This is user-owned, uncommitted work. The bootstrap must preserve it and must not stage, commit, overwrite, or revert it.

### Local environment evidence

| Tool | Observed state | Phase 0 implication |
|---|---|---|
| Git | `2.51.0.windows.1` | Available |
| Go | `go1.26.5 windows/amd64` | Required toolchain is available |
| `GOTOOLCHAIN` | `auto` | Must not be relied on silently; record any automatic download |
| CGO | enabled | Suitable for native Fyne build when toolchain prerequisites exist |
| GCC | `14.2.0` | Present on current Windows host |
| Node.js | `v24.15.0` | Compatible with the selected Vite/test toolchain |
| npm | `9.6.7` | Not used for workspace installation |
| pnpm | `10.11.0` | Selected package-manager version |
| Corepack | `0.34.6` | Available |
| golangci-lint | `v1.64.8` | Too old; v2.12.2 is still required for Phase 0 lint evidence |
| goose | `v3.27.1` | Installed but out of scope and must not be run |
| sqlc | missing | Not a Phase 0 blocker; integration starts later |
| Fyne CLI | missing | Not required for `go build`; packaging is out of scope |
| Docker | missing | Not a Phase 0 blocker |
| make / clang | missing | Make is not required; Windows build uses GCC |

At the Milestone 2 start, branch `dev` pointed to `ad73355` with a clean worktree. No Phase 0 build, test, dependency installation, migration, or application command had been run at that point.

## 3. Source-of-truth documents

Conflicts use the following precedence. A lower source supplies context only when it does not expand or contradict a higher source.

1. `docs/Corvus_Studio_Development_Implementation_Plan_v0.1.md`
   - Defines Repository Bootstrap as Monorepo, Go module/workspace, pnpm workspace, GitHub Actions, README, and License.
   - Places Echo, Viper, Zap, SQLite, goose, sqlc, and Cobra in Phase 1.
2. `docs/Corvus_Studio_Repository_Structure_Design_v0.1.md`
   - Defines the target top-level tree, Core/Web/Launcher/Agent separation, package reservations, tests, tools, deployments, and `.github` layout.
3. `docs/Corvus_Studio_Technology_Stack_Decision_v0.1.md`
   - Fixes Go 1.26, Echo v5, SQLite/modernc, goose/sqlc, React 19, Vite 8, pnpm, Fyne, Google ADK Go, and quality tools.
4. `docs/Corvus_Studio_Engineering_and_Deployment_Guide_v0.1.md`
   - Defines GitHub Actions checks, React-to-Go production delivery direction, supported platforms, and absence of automatic updates in v0.1.
5. `USAGE.md`
   - Defines intended developer commands and experience; commands for future phases must be labelled rather than presented as currently working.

The following documents are complete background sources but do not authorize Phase 1+ implementation in this plan:

- `docs/Corvus_Studio_ADK_Implementation_Plan_v0.1.md`: future in-process ADK runtime and workflows.
- `docs/Corvus_Studio_Agent_Architecture_and_Tool_Specification_v0.1.md`: future Agent/tool boundaries and evidence model.
- `docs/Corvus_Studio_Backend_API_Design_v0.1.md`: future REST/SSE contract.
- `docs/Corvus_Studio_Data_Model_Design_v0.1.md`: future domains and persistence concepts.
- `docs/Corvus_Studio_Launch_Frontend_Specification_v0.1.md`: future pages and frontend libraries.
- `docs/Corvus_Studio_Launch_Functional_Specification_v0.1.md`: future business capabilities.
- `docs/Corvus_Studio_Launch_MVP_Roadmap_v0.1.md`: product roadmap; its Phase labels are subordinate to the Development Implementation Plan.
- `docs/Corvus_Studio_Launch_PRD_v0.1.md`: product requirements and v0.1 outcome.
- `docs/Corvus_Studio_Launch_System_Design_Document_v0.1.md`: future runtime/domain/storage boundaries.
- `docs/Corvus_Studio_Launch_UX_IA_v0.1.md`: future navigation and workflows.
- `docs/Corvus_Studio_Product_Design_Draft_v0.1.md`: product direction and licensing recommendation.

`LICENSE` is preserved as the current repository license. `README.md` is empty; `USAGE.md` is aspirational and may receive the minimal status/command corrections specified by this plan.

## 4. Scope

Phase 0 will:

- verify and preserve the already-normalized Git worktree at the approved outer repository root;
- create the target Monorepo skeleton and root engineering files;
- create exactly three Go modules for Core, Launcher, and Agent, joined by a root `go.work`;
- create a minimal buildable Core bootstrap command;
- create a minimal Fyne Launcher that proves native compilation and shows no business UI;
- create an empty Agent package boundary without ADK dependencies or behavior;
- create a pnpm workspace with a minimal React 19/Vite 8/TypeScript/SCSS application;
- add frontend lint, formatting check, one smoke test, and production build;
- reserve shared package and later infrastructure locations with short README placeholders only;
- add GitHub Actions for Go quality, frontend quality, and native Core/Launcher builds on Windows, macOS, and Linux;
- minimally synchronize README and USAGE with commands that actually exist;
- update this ExecPlan with observed progress and evidence during execution.

## 5. Non-goals

Phase 0 will not implement or wire:

- Echo, HTTP routes, SSE, Viper, Zap, lumberjack, Cobra commands, or a running Core server;
- SQLite, modernc.org/sqlite, goose execution, sqlc generation, database files, schema, or business tables;
- an OpenAPI contract or generated API client;
- Project, Release Goal, Checklist, Resource, Deliverable, Template, Project Codex, Review, Report, Asset Map, authentication, settings, backup, export, or Watch behavior;
- Steam templates, Steam rules, Steam submission, or Steam data analysis;
- Google ADK Go, Agent runtime, prompts, workflows, tools, providers, models, or model credentials;
- Ant Design, Zustand, TanStack Query, React Router, React Flow, react-markdown, product pages, product navigation, or business UI;
- Launcher Core process management, system tray behavior, browser launching, packaging, signing, notarization, or automatic updates;
- Docker, systemd, production deployment, formal installers, release workflows, or non-GitHub CI;
- Issue templates, Pull Request templates, `CONTRIBUTING.md`, enforced Conventional Commits, or Git commits.

OpenAPI source placement, sqlc layout, migrations, and production asset staging are documented as future destinations only. Their directories and functional files are not created in Phase 0 unless explicitly marked as a placeholder in Section 7.

## 6. Document conflicts and decisions

| Status | Conflict or ambiguity | Decision and rationale |
|---|---|---|
| Accepted | The initial working directory was outside the nested Git repository; the root was normalized between assistant turns. | The outer `corvus` directory is the permanent root. Do not repeat the move. Phase 0 begins by auditing the current root, exact document inventory, four tracked-document blob matches, and a new SHA-256 baseline for all 15 documents. |
| Accepted | Development Plan Phase 0 is infrastructure-only, while the lower-priority MVP Roadmap lists Go Core, React, Fyne, SQLite, and CI/CD under Phase 0. | Create minimal buildable Core/React/Fyne shells, but follow the higher-priority plan by deferring SQLite and all runtime/business wiring to Phase 1+. |
| Accepted | The Development Plan and USAGE use `go test ./...`, but three nested modules in a root `go.work` are not reliably covered by that pattern from a non-module root. | Use `go test ./apps/core/... ./apps/launcher/... ./agent/...`. Do not add a misleading root module merely to preserve the shorter command. User confirmed this choice. |
| Accepted | Repository design separates top-level `agent/`, while the runtime design says the Agent runs in the Core process. | `agent/` is an independent Go module and future source boundary; Core will import/run it later. Phase 0 adds no ADK dependency. |
| Accepted | Fyne creates native CGO/toolchain complexity during bootstrap. | Use a real Fyne v2.8.0 minimal window and validate native builds on all three operating systems, without implementing launcher responsibilities. User confirmed this choice. |
| Accepted | USAGE presents future commands as if they work now. | Preserve product guidance but add a clear phase-status distinction and list only verified Phase 0 commands as current. |
| Accepted | The repository structure names shared packages, but there are no contracts or shared components yet. | Create only explanatory README placeholders. Do not create package manifests, exports, generated code, or dependencies until a consuming phase exists. |
| Accepted | OpenAPI/sqlc/goose locations are needed, but their content is later-phase work. | Reserve future locations as `apps/core/openapi/`, `packages/api-client/`, `apps/core/sqlc.yaml`, Core-internal SQL query/generated packages, and `apps/core/migrations/`. Only `packages/api-client/README.md` exists in Phase 0. |
| Open Question | Product Design recommends AGPL-3.0-or-later for code and CC BY 4.0 for docs, while the repository has one AGPLv3 license and no separate docs license. | Preserve the existing `LICENSE` byte-for-byte. Decide separate documentation licensing outside Phase 0 with appropriate project/legal review. This does not block bootstrap. |

## 7. Proposed repository tree

Legend:

- `[existing]`: preserve existing content.
- `[plan]`: planning artifact created before implementation.
- `[P0]`: create functional Phase 0 content.
- `[reserve]`: create only the named explanatory placeholder; no implementation or package API.
- `[future]`: intended location, absent after Phase 0.
- `[generated]`: generated by an approved Phase 0 tool and committed if applicable.

```text
corvus/
├── .agent/                                      [plan]
│   ├── PLANS.md
│   └── phase-0-repository-bootstrap.md
├── .github/
│   ├── workflows/
│   │   └── ci.yml                               [P0]
│   ├── SECURITY.md                              [future: contact policy required]
│   ├── ISSUE_TEMPLATE/                          [must not create]
│   └── pull_request_template.md                  [must not create]
├── agent/                                       [P0: Go module boundary]
│   ├── doc.go                                   [P0]
│   ├── go.mod                                   [P0]
│   ├── runtime/                                 [future]
│   ├── workflows/                               [future]
│   ├── tools/                                   [future]
│   ├── context/                                 [future]
│   ├── prompts/                                 [future]
│   └── providers/                               [future]
├── apps/
│   ├── core/                                    [P0: Go module]
│   │   ├── cmd/
│   │   │   └── corvus/
│   │   │       ├── main.go                      [P0]
│   │   │       └── main_test.go                 [P0]
│   │   ├── go.mod                               [P0]
│   │   ├── internal/                            [future]
│   │   ├── openapi/                             [future]
│   │   ├── migrations/                          [future]
│   │   ├── configs/                             [future]
│   │   └── sqlc.yaml                            [future]
│   ├── launcher/                                [P0: Go/Fyne module]
│   │   ├── main.go                              [P0]
│   │   ├── go.mod                               [P0]
│   │   └── go.sum                               [generated]
│   └── web/                                     [P0: pnpm package]
│       ├── src/
│       │   ├── styles/
│       │   │   └── index.scss                   [P0]
│       │   ├── App.test.tsx                     [P0]
│       │   ├── App.tsx                          [P0]
│       │   ├── main.tsx                         [P0]
│       │   ├── test-setup.ts                    [P0]
│       │   ├── pages/                           [future]
│       │   ├── features/                        [future]
│       │   ├── components/                      [future]
│       │   ├── api/                             [future]
│       │   ├── hooks/                           [future]
│       │   └── stores/                          [future]
│       ├── eslint.config.js                     [P0]
│       ├── index.html                           [P0]
│       ├── package.json                         [P0]
│       ├── tsconfig.app.json                    [P0]
│       ├── tsconfig.json                        [P0]
│       ├── tsconfig.node.json                   [P0]
│       └── vite.config.ts                       [P0]
├── packages/
│   ├── api-client/
│   │   └── README.md                            [reserve]
│   ├── shared-types/
│   │   └── README.md                            [reserve]
│   └── ui/
│       └── README.md                            [reserve]
├── configs/
│   └── README.md                                [reserve]
├── deployments/
│   └── README.md                                [reserve; no deployment files]
├── docs/                                        [existing; preserve bytes]
│   ├── Corvus_Studio_ADK_Implementation_Plan_v0.1.md
│   ├── Corvus_Studio_Agent_Architecture_and_Tool_Specification_v0.1.md
│   ├── Corvus_Studio_Backend_API_Design_v0.1.md
│   ├── Corvus_Studio_Data_Model_Design_v0.1.md
│   ├── Corvus_Studio_Development_Implementation_Plan_v0.1.md
│   ├── Corvus_Studio_Engineering_and_Deployment_Guide_v0.1.md
│   ├── Corvus_Studio_Launch_Frontend_Specification_v0.1.md
│   ├── Corvus_Studio_Launch_Functional_Specification_v0.1.md
│   ├── Corvus_Studio_Launch_MVP_Roadmap_v0.1.md
│   ├── Corvus_Studio_Launch_PRD_v0.1.md
│   ├── Corvus_Studio_Launch_System_Design_Document_v0.1.md
│   ├── Corvus_Studio_Launch_UX_IA_v0.1.md
│   ├── Corvus_Studio_Product_Design_Draft_v0.1.md
│   ├── Corvus_Studio_Repository_Structure_Design_v0.1.md
│   └── Corvus_Studio_Technology_Stack_Decision_v0.1.md
├── scripts/
│   └── README.md                                [reserve]
├── tests/
│   ├── integration/README.md                    [reserve]
│   ├── scenario/README.md                       [reserve]
│   └── fixtures/README.md                       [reserve]
├── tools/
│   └── README.md                                [reserve]
├── .editorconfig                                [P0]
├── .gitattributes                               [P0]
├── .gitignore                                   [P0]
├── .go-version                                  [P0: 1.26.5]
├── .golangci.yml                                [P0: v2 format]
├── .node-version                                [P0: 24.15.0]
├── .npmrc                                       [P0]
├── .prettierignore                              [P0]
├── .prettierrc.json                             [P0]
├── AGENTS.md                                    [plan]
├── LICENSE                                      [existing; unchanged]
├── README.md                                    [existing empty file; P0 update]
├── USAGE.md                                     [existing; minimal P0 update]
├── go.work                                      [P0]
├── go.work.sum                                  [generated only if Go creates it]
├── package.json                                 [P0]
├── pnpm-lock.yaml                               [generated with pnpm 10.11.0]
└── pnpm-workspace.yaml                          [P0]
```

## 8. Milestones

All commands in this section run from `E:\Git\开源\agent\corvus` in PowerShell unless stated otherwise. Update Progress immediately after each acceptance check.

### Milestone 1 — Audit and protect the normalized repository root

**Goal:** verify the already-normalized outer Git root, preserve the pre-existing dirty state, and establish current document evidence without repeating a filesystem move.

**Paths:** read-only inspection of `.git/`, `docs/`, `LICENSE`, `README.md`, and `USAGE.md`. The obsolete `corvus-studio/` path must remain absent.

**Operations and evidence:**

```powershell
$phase0Evidence = Join-Path ([IO.Path]::GetTempPath()) 'corvus-phase0-preflight'
$phase0ExpectedDocs = @(
    'Corvus_Studio_ADK_Implementation_Plan_v0.1.md',
    'Corvus_Studio_Agent_Architecture_and_Tool_Specification_v0.1.md',
    'Corvus_Studio_Backend_API_Design_v0.1.md',
    'Corvus_Studio_Data_Model_Design_v0.1.md',
    'Corvus_Studio_Development_Implementation_Plan_v0.1.md',
    'Corvus_Studio_Engineering_and_Deployment_Guide_v0.1.md',
    'Corvus_Studio_Launch_Frontend_Specification_v0.1.md',
    'Corvus_Studio_Launch_Functional_Specification_v0.1.md',
    'Corvus_Studio_Launch_MVP_Roadmap_v0.1.md',
    'Corvus_Studio_Launch_PRD_v0.1.md',
    'Corvus_Studio_Launch_System_Design_Document_v0.1.md',
    'Corvus_Studio_Launch_UX_IA_v0.1.md',
    'Corvus_Studio_Product_Design_Draft_v0.1.md',
    'Corvus_Studio_Repository_Structure_Design_v0.1.md',
    'Corvus_Studio_Technology_Stack_Decision_v0.1.md'
)

if ((git rev-parse --show-toplevel) -ne 'E:/Git/开源/agent/corvus') {
    throw 'Unexpected Git root'
}
if (Test-Path -LiteralPath '.\corvus-studio') {
    throw 'Obsolete nested repository path reappeared'
}

$phase0ActualDocs = Get-ChildItem -LiteralPath '.\docs' -File |
    Select-Object -ExpandProperty Name |
    Sort-Object
$phase0NameDiff = Compare-Object ($phase0ExpectedDocs | Sort-Object) $phase0ActualDocs
if ($phase0NameDiff) {
    throw "Document inventory mismatch: $phase0NameDiff"
}

foreach ($phase0Name in @(
    'Corvus_Studio_Launch_Functional_Specification_v0.1.md',
    'Corvus_Studio_Launch_PRD_v0.1.md',
    'Corvus_Studio_Launch_UX_IA_v0.1.md',
    'Corvus_Studio_Product_Design_Draft_v0.1.md'
)) {
    $phase0HeadHash = git rev-parse ('HEAD:' + $phase0Name)
    $phase0CurrentHash = git hash-object ('docs/' + $phase0Name)
    if ($phase0HeadHash -ne $phase0CurrentHash) {
        throw "Tracked document content changed: $phase0Name"
    }
}

New-Item -ItemType Directory -Force -Path $phase0Evidence | Out-Null
git status --porcelain=v2 --branch |
    Set-Content -LiteralPath (Join-Path $phase0Evidence 'git-status-current.txt')
Get-ChildItem -LiteralPath '.\docs' -File |
    Sort-Object Name |
    Get-FileHash -Algorithm SHA256 |
    Select-Object Path, Hash |
    Export-Csv -NoTypeInformation -LiteralPath (Join-Path $phase0Evidence 'docs-current.csv')
```

Expected output: no error; exactly 15 expected documents; all four tracked-document comparisons match; current Git status and a 15-document SHA-256 baseline exist under the temporary evidence directory.

**Acceptance observed on 2026-07-29:** Git root is `E:/Git/开源/agent/corvus`, the nested path is absent, all 15 expected names are present, and the four tracked documents match their HEAD blobs. Full pre-move SHA-256 evidence for the other 11 originally untracked documents cannot be reconstructed and is retained as a documented limitation.

**Recovery:** this revised milestone is read-only except for temporary evidence. If any assertion fails, stop before modifying repository content and investigate the mismatch. Do not recreate the nested directory, repeat the move, or use Git reset/checkout.

### Milestone 2 — Create the root skeleton and three Go module shells

**Goal:** establish repository boundaries and prove Core, Agent, and real Fyne Launcher compilation without runtime or business dependencies.

**Paths:** root engineering files, reserved README files, `apps/core`, `apps/launcher`, `agent`, and `go.work` as listed in Section 7.

**Tool precondition:**

```powershell
go version
node --version
pnpm --version
gcc --version
```

Required observations:

- Go reports `go1.26.x`; target patch for `.go-version` and CI is `1.26.5`.
- Node reports `v24.15.0` or another approved Node 24 version satisfying project engines.
- pnpm reports exactly `10.11.0` before creating or changing `pnpm-lock.yaml`.
- Windows reports a usable GCC.

If Go remains 1.25.8, stop before running Go initialization or validation. Do not lower `go` directives. Toolchain installation is an external prerequisite, not an operation in this ExecPlan.

**Operations:**

- Create root ignore/editor/version/lint files using explicit content, not an interactive generator.
- Create explanatory placeholders for `configs`, `deployments`, `scripts`, `tests`, `tools`, and the three shared package directories.
- Create Core `main.go` with a small testable `run(io.Writer)` function that prints `Corvus Studio core bootstrap` and exits successfully; add a test for that exact output.
- Create `agent/doc.go` with package documentation only.
- Create Launcher `main.go` using `fyne.io/fyne/v2/app` and a simple label in one window. It must not spawn Core, add a tray, open a browser, or contain product navigation.
- Initialize only the three approved module paths and one workspace.

Exact initialization commands, guarded so they run only when the corresponding manifest is absent:

```powershell
go -C .\apps\core mod init github.com/gofurry/corvus-studio/apps/core
go -C .\apps\launcher mod init github.com/gofurry/corvus-studio/apps/launcher
go -C .\agent mod init github.com/gofurry/corvus-studio/agent
go work init .\apps\core .\apps\launcher .\agent
go -C .\apps\launcher get fyne.io/fyne/v2@v2.8.0
go -C .\apps\core mod tidy
go -C .\apps\launcher mod tidy
go -C .\agent mod tidy
go work sync
gofmt -w .\apps\core\cmd\corvus\main.go .\apps\core\cmd\corvus\main_test.go .\apps\launcher\main.go .\agent\doc.go
```

If manifests were created explicitly before these commands, skip `mod init`/`work init`, verify their module/use values, and run only dependency/tidy/sync commands.

**Acceptance:**

```powershell
go env GOWORK
go work edit -json
go list -m
go test ./apps/core/... ./apps/launcher/... ./agent/...

$phase0Bin = Join-Path ([IO.Path]::GetTempPath()) 'corvus-phase0-bin'
New-Item -ItemType Directory -Force -Path $phase0Bin | Out-Null
go build -o (Join-Path $phase0Bin 'corvus.exe') ./apps/core/cmd/corvus
go build -o (Join-Path $phase0Bin 'corvus-launcher.exe') ./apps/launcher
& (Join-Path $phase0Bin 'corvus.exe')
```

Expected output:

- `GOWORK` is the root `go.work` absolute path;
- `go work edit -json` contains exactly `./apps/core`, `./apps/launcher`, and `./agent` in `Use`;
- `go list -m` lists the three approved module paths;
- tests exit zero, including the Core output test;
- both temporary executables exist;
- running Core prints `Corvus Studio core bootstrap` and exits zero.

Optional manual observation on a graphical Windows host:

```powershell
go run ./apps/launcher
```

A minimal Corvus Studio Launcher window appears and closes normally. This manual check is not a substitute for native CI builds.

**Recovery:** keep pre-existing-file backups under the temporary evidence directory. If module initialization fails, restore only manifests that existed before the milestone; otherwise move the known new module/workspace files to the temporary evidence directory for inspection. Do not delete or edit `docs`, `LICENSE`, `README.md`, `USAGE.md`, `AGENTS.md`, or `.agent/`. Rerun `go mod tidy` only with Go 1.26 and the pinned Fyne version.

### Milestone 3 — Create the pnpm and minimal Web workspace

**Goal:** provide one deterministic pnpm workspace with a minimal React page that lints, tests, and builds.

**Paths:** root `package.json`, `pnpm-workspace.yaml`, `.npmrc`, Prettier files, `pnpm-lock.yaml`, and `apps/web` files listed in Section 7.

**Operations:**

- Write manifests and source explicitly; do not use `create-vite` against a possibly populated path.
- Root package name is `corvus-studio`, is private, and declares `packageManager: pnpm@10.11.0`.
- Root scripts delegate `dev:web`, `lint`, `test`, `build`, and `format:check` to workspace packages.
- Web package name is `@corvus-studio/web`, is private, and uses `vite`, `tsc -b && vite build`, `eslint`, `vitest run`, and Prettier check scripts.
- The Web source renders only a bootstrap heading/description and imports one SCSS file.
- The smoke test renders the app and asserts the Corvus Studio heading.
- Vite build output remains `apps/web/dist`. No proxy, API call, or embed copy is added because Core has no server.

Initial dependency set and versions:

- Runtime: `react@19.2.8`, `react-dom@19.2.8`.
- Build/language: `vite@8.1.5`, `@vitejs/plugin-react@6.0.4`, `typescript@6.0.2`, `sass@1.102.0`.
- Vite/Rolldown optional dependency policy: exclude only `@rolldown/binding-wasm32-wasi`; Corvus targets native Windows, macOS, and Linux bindings, and this optional fallback contains an internally unsatisfiable Emnapi peer graph under strict pnpm validation.
- Lint/format: `eslint@10.8.0`, `@eslint/js@10.0.1`, `typescript-eslint@8.65.0`, `eslint-plugin-react-hooks@7.1.1`, `eslint-plugin-react-refresh@0.5.3`, `globals@17.8.0`, `prettier@3.9.6`.
- Test: `vitest@4.1.10`, `jsdom@30.0.1`, `@testing-library/react@16.3.2`, `@testing-library/dom@10.4.1`, `@testing-library/jest-dom@7.0.0`.
- Types: `@types/node@24.13.3`, `@types/react@19.2.17`, `@types/react-dom@19.2.3`.

Do not add Ant Design, Zustand, TanStack Query, React Router, React Flow, react-markdown, or a shared workspace package manifest.

**Lockfile creation and acceptance:**

```powershell
pnpm --version
pnpm install
pnpm list -r --depth -1
pnpm format:check
pnpm lint
pnpm test
pnpm build
pnpm install --frozen-lockfile
Test-Path -LiteralPath '.\apps\web\dist\index.html'
```

Expected output:

- pnpm is exactly 10.11.0 before both install commands;
- the recursive list contains the root and `@corvus-studio/web`, but no placeholder package;
- format, lint, and test exit zero; Vitest reports the bootstrap smoke test passed;
- Vite reports a successful production build;
- frozen reinstall reports an up-to-date lockfile and exits zero;
- the final path check prints `True`.

**Recovery:** never run another package manager and never regenerate the lockfile under a different pnpm version. If installation fails, retain `pnpm-lock.yaml` for inspection, remove only generated `node_modules`/`dist` paths or move them outside the repository, correct manifests, and rerun `pnpm install` with 10.11.0. Restore any pre-existing manifest from its temporary backup rather than using Git cleanup.

### Milestone 4 — Add CI and align developer documentation

**Goal:** make local quality commands visible in GitHub Actions and document only capabilities that exist at the end of Phase 0.

**Paths:** `.github/workflows/ci.yml`, `.golangci.yml`, `README.md`, `USAGE.md`, and this ExecPlan's living sections.

**CI operations:**

- Trigger on pull requests and pushes to `main`; grant `contents: read` only.
- Use current action majors verified during planning: `actions/checkout@v7`, `actions/setup-go@v7`, `actions/setup-node@v7`, `pnpm/action-setup@v6`, and `golangci/golangci-lint-action@v9`.
- Go quality job on Ubuntu:
  - set up Go from `.go-version` and cache `**/go.sum`;
  - fail if `gofmt -l` reports files under `apps/core`, `apps/launcher`, or `agent`;
  - run `go test ./apps/core/... ./apps/launcher/... ./agent/...`;
  - run golangci-lint v2.12.2 with the same explicit package patterns.
- Frontend job on Ubuntu:
  - set up pnpm 10.11.0 and Node from `.node-version`;
  - cache the pnpm store using `pnpm-lock.yaml`;
  - run `pnpm install --frozen-lockfile`, `pnpm format:check`, `pnpm lint`, `pnpm test`, and `pnpm build`.
- Native build matrix on `ubuntu-latest`, `windows-latest`, and `macos-latest`:
  - Linux installs `gcc libgl1-mesa-dev xorg-dev libwayland-dev libxkbcommon-dev`;
  - Windows uses `msys2/setup-msys2@v2` to provide `mingw-w64-x86_64-gcc` and verifies `gcc --version`;
  - macOS verifies `xcode-select -p` and `clang --version`;
  - every runner builds `./apps/core/cmd/corvus` and `./apps/launcher` natively;
  - no packaging, signing, artifact publishing, or cross-compilation is performed.

**Documentation operations:**

- Populate the empty README with project identity, Phase 0 status, repository map, prerequisites, exact root validation commands, links to USAGE/design docs, and AGPL license notice.
- Preserve USAGE's product overview but insert a prominent implementation-status section.
- Replace current-development instructions with the working root commands and mark `corvus serve`, migration/sqlc, Docker, Agent, product workflows, and formal desktop build as future-phase material.
- Do not edit any file under `docs/` or change `LICENSE`.

**Local acceptance:**

```powershell
$phase0GoFiles = Get-ChildItem -Recurse -File -Include *.go -Path '.\apps\core','.\apps\launcher','.\agent'
$phase0Unformatted = $phase0GoFiles | ForEach-Object { gofmt -l $_.FullName }
if ($phase0Unformatted) { throw "Unformatted Go files: $phase0Unformatted" }

golangci-lint --version
golangci-lint run ./apps/core/... ./apps/launcher/... ./agent/...
go test ./apps/core/... ./apps/launcher/... ./agent/...
pnpm install --frozen-lockfile
pnpm format:check
pnpm lint
pnpm test
pnpm build
git diff --check
```

Required: golangci-lint reports v2.12.2. If the local tool remains v1.64.8, record local lint as blocked and rely only on an actually observed CI v2 result later; do not claim a local pass.

**Boundary audit:**

```powershell
$phase0Forbidden = 'labstack/echo|modernc\.org/sqlite|pressly/goose|sqlc|google.*adk|ReleaseGoal|Checklist|Deliverable|AssetMap|Steam'
rg -n --glob '*.go' --glob '*.ts' --glob '*.tsx' $phase0Forbidden apps agent packages
Get-ChildItem -Recurse -File -Filter '*.sql'
Test-Path '.\.github\ISSUE_TEMPLATE'
Test-Path '.\.github\pull_request_template.md'
Test-Path '.\CONTRIBUTING.md'
```

Expected output: ripgrep and SQL search return no application matches; all three forbidden repository-policy path checks print `False`.

**Acceptance and remote evidence:** locally validate workflow command parity and `git diff --check`. Because this plan does not commit or push, a GitHub Actions run is not automatically available. Record remote jobs as pending until a user-authorized push/PR causes a run; never state that GitHub Actions passed based only on YAML inspection.

**Recovery:** back up current untracked `README.md` and `USAGE.md` to the temporary evidence directory before editing. On failure, restore those exact copies and move newly created CI/config files to the evidence directory. Do not touch design documents or the license.

### Milestone 5 — Final evidence audit and handoff

**Goal:** prove that Phase 0 is complete, reproducible, and contains no later-phase implementation.

Run every applicable command in Section 13, record platform/exit code/summary in Progress, compare current document hashes with Milestone 1's current baseline, and inspect `git status --short` without staging or committing.

Expected result: completion criteria in Section 16 are either backed by current evidence or explicitly marked pending/blocked. Do not declare Phase 0 complete while a required local or CI-independent criterion lacks evidence. Remote CI remains a named follow-up if no authorized GitHub run exists.

Recovery is milestone-specific: fix only the failing Phase 0 artifact, rerun its smallest validation set, then rerun the final matrix. Never repair by implementing a later-phase dependency or by discarding existing user work.

## 9. Go workspace strategy

`go.work` is at the repository root. There is no root `go.mod`.

Module boundaries are:

| Path | Module path | Phase 0 responsibility | Reason for boundary |
|---|---|---|---|
| `apps/core` | `github.com/gofurry/corvus-studio/apps/core` | Minimal Core command only | Core owns future HTTP/domain/storage runtime and produces the main binary. |
| `apps/launcher` | `github.com/gofurry/corvus-studio/apps/launcher` | Minimal Fyne window | Native UI dependencies and build constraints remain isolated from headless Core. |
| `agent` | `github.com/gofurry/corvus-studio/agent` | Empty package boundary | Agent source is top-level by design and can later be imported into Core without becoming Core-internal code. |

No module is created for the repository root, `packages`, `tests`, or `tools`. A new Go module requires a later explicit architectural reason and an update to `go.work` and this decision record.

All module manifests use `go 1.26.0`; `.go-version` and CI select Go 1.26.5. Fyne is the only non-standard-library Phase 0 Go dependency. Echo, modernc.org/sqlite, goose, sqlc, zap, lumberjack, cobra, viper, easyhash, UUID libraries, and Google ADK are absent.

Workspace discovery and cross-module validation use:

```powershell
go env GOWORK
go work edit -json
go list -m
go test ./apps/core/... ./apps/launcher/... ./agent/...
```

Do not use root `go test ./...` as the acceptance command: from a workspace root that is not itself a module, the pattern does not provide the intended three-module coverage.

Future Core-to-Agent imports must use the Agent module path and be resolved by released/pinned module versions outside a workspace; Phase 0 does not add a temporary `replace` directive to any `go.mod`. The local `go.work` provides development resolution when an import is introduced later.

## 10. pnpm workspace strategy

The root `pnpm-workspace.yaml` includes:

```yaml
packages:
  - apps/web
  - packages/*
```

The root `package.json` is private, pins `pnpm@10.11.0`, declares Node `>=24.15.0 <25`, and provides the public developer command surface:

- `pnpm dev:web`
- `pnpm format:check`
- `pnpm lint`
- `pnpm test`
- `pnpm build`

`apps/web` is the only Phase 0 workspace package with a `package.json`. It contains the minimal React/Vite toolchain and no product feature dependencies.

Reservations:

- `packages/api-client`: README only. Future generated OpenAPI client output; no handwritten substitute.
- `packages/shared-types`: README only. Future TypeScript-only cross-package types that are not already generated from OpenAPI.
- `packages/ui`: README only. Future reusable components after a real second consumer exists.

The wildcard intentionally discovers these packages when future manifests are added. It does not make README-only directories into pnpm packages.

There is exactly one root `pnpm-lock.yaml`. Developers never run `pnpm install` from a package in a way that creates another lockfile. Any manifest change requires pnpm 10.11.0, a reviewed root lockfile diff, and a final frozen install.

## 11. Build and embedding strategy

### Development mode

Long-term development uses two independent processes:

- Go Core serves the local API.
- Vite serves the React application with HMR and forwards API requests to Core.

Phase 0 does not yet have an HTTP server or API, so it validates the processes independently:

- Core bootstrap builds/runs and exits after printing its marker.
- `pnpm dev:web` starts Vite and displays only the bootstrap page.
- No proxy target or port is committed until Phase 1 defines the Core server configuration.

### Production embedding

The future production sequence is:

1. build `apps/web` to `apps/web/dist`;
2. stage those assets into a generated directory inside the Core module, such as `apps/core/internal/webui/dist`;
3. embed from that in-module directory using `go:embed`;
4. serve the embedded SPA from Echo with API routes kept separate.

Go cannot embed arbitrary files outside the owning module/package tree. Therefore the future build must use an explicit, portable staging step rather than `//go:embed` against `../web/dist`. The staging utility, embed package, SPA fallback, cache headers, and Echo handler belong to Phase 1 and are not created or validated in Phase 0.

### Launcher

The Phase 0 Launcher validates only:

- Fyne can be resolved and compiled;
- a native window with Corvus Studio identity can be opened manually;
- Core and Agent dependencies do not leak into the Launcher module.

Core process lifecycle, system tray, browser launch, single-instance rules, port discovery, shutdown handling, packaging, and signing are future work.

## 12. CI strategy

Phase 0 must implement one GitHub Actions workflow with three concerns:

1. **Go quality on Ubuntu:** format check, explicit workspace test, and golangci-lint v2.12.2.
2. **Frontend quality on Ubuntu:** pinned pnpm/Node setup, frozen install, formatting check, lint, test, and build.
3. **Native build matrix:** Core and Fyne Launcher build on Windows, macOS, and Linux with platform-native toolchains.

Caching:

- `actions/setup-go` caches build/module data using discovered `**/go.sum` files.
- `actions/setup-node` caches pnpm store data with `pnpm-lock.yaml` as the dependency path.
- Cache contents are optimization only; a cache miss must still produce the same result.

Phase 0 CI must not:

- run database migrations, sqlc, OpenAPI generation, Agent tests, scenario tests, Docker builds, packaging, signing, releases, or deployment;
- upload release artifacts;
- introduce another CI provider;
- add Issue or PR templates.

Future additions are made when their owning phase has real content: OpenAPI drift/generation checks, migration/sqlc checks, integration tests, scenario/Agent regressions, end-to-end browser tests, database platform smoke tests, packaging, signing, and release jobs.

## 13. Validation matrix

| Target | Command or observation | Expected evidence | Phase 0 requirement |
|---|---|---|---|
| Windows local root | `git rev-parse --show-toplevel` | Approved outer absolute root | Required |
| Document preservation | Compare against Milestone 1 current SHA-256 baseline; verify four HEAD blobs | No post-baseline differences; four tracked blobs match | Required |
| Go toolchain | `go version` | Go 1.26.x; planned patch 1.26.5 | Required |
| Go workspace | `go env GOWORK`; `go work edit -json`; `go list -m` | Root workspace and exactly three modules | Required |
| Go tests | `go test ./apps/core/... ./apps/launcher/... ./agent/...` | Exit 0; Core bootstrap test passes | Required |
| Core build/run | Build to temporary path, execute binary | Exit 0 and exact bootstrap marker | Required |
| Launcher Windows build | `go build` to temporary path | Exit 0 and executable exists | Required |
| Launcher manual smoke | `go run ./apps/launcher` | Minimal window appears/closes | Recommended local observation |
| Go formatting | `gofmt -l` audit | No paths printed | Required |
| Go lint | golangci-lint v2.12.2 explicit patterns | Exit 0 | Required; CI may supply evidence if local version blocked |
| pnpm version | `pnpm --version` | `10.11.0` | Required before lockfile mutation |
| pnpm workspace | `pnpm list -r --depth -1` | Root and Web only | Required |
| Frozen install | `pnpm install --frozen-lockfile` | Exit 0, unchanged lockfile | Required |
| Frontend formatting/lint | `pnpm format:check`; `pnpm lint` | Exit 0 | Required |
| Frontend tests | `pnpm test` | Bootstrap render test passes | Required |
| Frontend build | `pnpm build`; check `apps/web/dist/index.html` | Exit 0 and file exists | Required |
| Linux CI | Native matrix job | Core and Fyne builds exit 0 | Required before claiming remote CI complete |
| macOS CI | Native matrix job | Core and Fyne builds exit 0 | Required before claiming remote CI complete |
| Windows CI | Native matrix job | Core and Fyne builds exit 0 | Required before claiming remote CI complete |
| Phase boundary | Forbidden dependency/symbol and SQL searches | No application matches | Required |
| Repository policy | Test forbidden template/contribution paths | All false | Required |
| Worktree review | `git status --short`; `git diff --check` | Only intended new/updated paths; no whitespace errors | Required |

If GitHub Actions cannot run because no commit/push is authorized, record the three remote jobs as pending. Local checks do not constitute remote CI evidence.

## 14. Idempotence and recovery

- **Existing files:** inspect every target before creation. If an unexpected file exists, stop and compare; never use force overwrite. Planning files are expected at the outer root and are preserved.
- **Normalized root:** do not repeat the completed external relocation. Assert the outer Git root and absence of `corvus-studio/` before other milestones.
- **Document baseline:** retain the current 15-document SHA-256 inventory as the only available baseline for the 11 originally untracked documents; continue verifying the four tracked documents against HEAD blobs.
- **Go initialization:** run `mod init` and `work init` only when manifests are absent. Otherwise verify module/use declarations and use `go mod tidy`/`go work sync` with the pinned toolchain.
- **Go module count:** do not create extra modules to work around commands. The three approved boundaries are an invariant.
- **Generated Go sums:** let Go update only module/workspace sum files. Review the paths and never hand-edit checksums.
- **pnpm lockfile:** generate once from the root using pnpm 10.11.0. All later verification uses `--frozen-lockfile`. Never run npm/yarn installation or a second pnpm version.
- **Partial frontend install:** generated `node_modules` and `dist` may be safely moved out of the repository and regenerated; manifests and lockfile must be preserved for diagnosis.
- **Documentation:** capture README/USAGE copies before editing and hash all files under `docs/`. Do not format or rewrite design documents.
- **CI/config recovery:** newly created files can be moved to the temporary evidence directory. Do not use destructive Git commands because the worktree was dirty before Phase 0.
- **Repeated validation:** validation commands may update build/module/package caches and ignored build outputs, but must not alter tracked or user-owned source unexpectedly. A second frozen install and second build must succeed without manifest changes.

## 15. Risks and blockers

### Go 1.26

Go 1.26 is available from the official distribution service; 1.26.5 was the current stable patch observed during planning. The local host has 1.25.8. This is a local blocker for implementation evidence, not permission to lower the project version. Install/activate Go 1.26 outside this plan, then record `go version` before continuing.

### Echo v5

The module proxy exposed `github.com/labstack/echo/v5` through v5.3.1, so the selected major exists. It currently requires Go 1.25 or later. Echo is still a Phase 1 dependency and must be absent from Phase 0 module graphs.

### React 19, Vite 8, Node, and pnpm

Vite 8.1.5 requires Node `^20.19.0 || >=22.12.0`; Node 24.15.0 satisfies it. React 19 and the selected Vite React plugin are available. `jsdom@30.0.1` requires Node `^24.15.0` on the Node 24 line, which is why the minimum is 24.15.0. pnpm 10.11.0 supports the selected Node version and is already installed. The lockfile remains the compatibility evidence; no silent package downgrade or substitution is allowed.

### Fyne native dependencies

Fyne v2.8.0 is available and its module requires Go 1.22 or later. Desktop compilation requires CGO, a C compiler, and platform graphics headers. Windows uses MinGW/MSYS2, macOS uses Xcode command-line tools and system graphics frameworks, and Ubuntu needs X11/OpenGL/Wayland development packages. Native CI avoids unsupported assumptions from cross-compilation. Packaging is not validated.

### modernc.org/sqlite

modernc.org/sqlite v1.54.0 was available and requires Go 1.25 or later. It preserves the no-CGO database strategy, but actual runtime and cross-platform database smoke tests cannot be claimed until Phase 1 imports and exercises it. Phase 0 records placement only.

### Lint tool mismatch

The local golangci-lint v1.64.8 is not the approved v2 tool and may not understand v2 configuration or Go 1.26. Local lint remains blocked until v2.12.2 is available. The CI action pins v2.12.2; only an observed action run is CI lint evidence.

### Git and document state

The dirty worktree makes destructive cleanup unacceptable. The repository is now at the approved outer root, and the obsolete nested directory is absent. README and USAGE are untracked, so Git alone cannot restore their current content; temporary copies are mandatory. The four moved design documents are byte-identical to HEAD and must remain in `docs/` according to user intent. The other 11 untracked documents have a current baseline but no reconstructable pre-move hash manifest.

### Deferred platform/release concerns

macOS signing/notarization, Windows/Linux packaging, Docker, systemd, release artifacts, automatic update, and deployment verification are deliberately excluded. Their absence is not a Phase 0 failure.

## 16. Completion criteria

Phase 0 is complete only when evidence shows all of the following:

- the approved outer directory is the Git root and the obsolete inner container is absent;
- all 15 design documents retain the Milestone 1 current SHA-256 baseline, the four tracked documents continue to match their HEAD blobs, and `LICENSE` is unchanged;
- the target Phase 0 tree exists, while paths marked `[future]` remain absent;
- `go.work` is recognized and contains exactly Core, Launcher, and Agent modules;
- each minimal Go module can be tested or built using the explicit workspace commands;
- the Core bootstrap output test passes and the binary prints the documented marker;
- Fyne Launcher builds natively on Windows, Linux, and macOS; its Phase 0 UI contains no business behavior;
- the root pnpm workspace installs under pnpm 10.11.0 from the single frozen lockfile;
- the Web package passes formatting check, lint, smoke test, and production build;
- GitHub Actions contains the same commands as local documentation and has no deployment/release behavior;
- README and USAGE describe only genuinely available Phase 0 commands as current;
- no application dependency, import, SQL file, type, route, page, or workflow implements Phase 1+ capabilities;
- no Issue/PR templates, `CONTRIBUTING.md`, Docker/systemd files, installers, signing files, or automatic-update configuration exist;
- `git diff --check` passes, the worktree review identifies only intended Phase 0/planning paths plus the preserved pre-existing user changes, and no Git commit was created;
- all locally runnable required checks have recorded exit codes and results; any remote CI run not authorized or available is explicitly pending rather than reported as passed.

## 17. Progress, discoveries and decision log

### Progress

- [x] 2026-07-29 — Inspected the actual outer and nested directories without modifying repository content.
- [x] 2026-07-29 — Read all 15 design documents, README, USAGE, and LICENSE.
- [x] 2026-07-29 — Recorded Git state, document move equivalence, local tool versions, and dependency availability evidence.
- [x] 2026-07-29 — Confirmed target root, three-module strategy, explicit Go test command, and real Fyne skeleton with the user.
- [x] 2026-07-29 — Created the three allowed planning artifacts only.
- [x] 2026-07-29 16:53 +08:00 — Milestone 1 audit: outer Git root confirmed, nested path absent, all 15 expected document names present, and four tracked documents match HEAD blobs. Limitation: no pre-move SHA-256 baseline exists for the 11 originally untracked documents.
- [x] 2026-07-29 17:08 +08:00 — Execution preflight: branch `dev`, HEAD `ad73355`, clean worktree, Go 1.26.5, Node 24.15.0, pnpm 10.11.0, GCC 14.2.0; golangci-lint remains v1.64.8.
- [x] 2026-07-29 17:20 +08:00 — Milestone 2 workspace/format validation: `go env GOWORK`, `go work edit -json`, and `go list -m` identified exactly the root workspace and three approved modules; `gofmt -l` returned no files.
- [x] 2026-07-29 17:20 +08:00 — Milestone 2 test/build validation: explicit three-module `go test` exited 0; Core and Windows Fyne Launcher built to the temporary directory; Core printed the exact bootstrap marker; forbidden Go module scan returned no Phase 1+ dependencies; `git diff --check` exited 0.
- [x] Milestone 2 — Root skeleton and Go module shells created and locally validated on Windows.
- [x] 2026-07-29 17:28 +08:00 — Milestone 3 dependency validation: pnpm 10.11.0 generated one root lockfile; strict peer validation remained enabled; frozen install exited 0 after excluding only Rolldown's broken optional WASM fallback.
- [x] 2026-07-29 17:28 +08:00 — Milestone 3 frontend validation: workspace list contained root and Web only; format, ESLint, one Vitest/Testing Library smoke test, TypeScript/Vite build, and frozen reinstall all exited 0; `apps/web/dist/index.html` existed.
- [x] Milestone 3 — pnpm and minimal Web workspace created and locally validated on Windows.
- [ ] Milestone 4 — Add CI and align developer documentation.
- [ ] Milestone 5 — Complete final evidence audit and handoff.

### Surprises & Discoveries

- The provided working directory was not the Git root; the complete worktree was nested one level below it.
- The nested worktree was already dirty. Four tracked root documents had been moved byte-for-byte into an untracked `docs/` directory, and most supplied documentation was untracked.
- The local Go version is 1.25.8 even though Go 1.26 is available and mandated; Node and pnpm already satisfy the selected frontend baseline.
- Echo v5 is now a published stable module line, so its risk is Phase 1 integration rather than availability.
- Fyne remains the only Phase 0 component requiring native C/graphics build prerequisites.
- A root `go test ./...` command would not express the intended coverage for a non-module workspace root, so explicit module patterns are necessary.
- The repository-root relocation occurred between assistant turns rather than through Milestone 1's recorded command sequence. The root and document inventory are present, but the complete pre-move hash manifest cannot be reconstructed.
- Before execution, the repository had advanced externally to commit `ad73355`, both `main` and `dev` referenced that commit, the active branch was `dev`, and the worktree was clean.
- Go 1.26.5 became available between planning and execution, removing the Milestone 2 toolchain blocker without changing the project version.
- The first two full workspace test attempts timed out during the initial Windows Fyne/GLFW native compilation at 124 and 304 seconds. An isolated `go test -x ./apps/launcher/...` completed successfully after the build cache was populated, and the required full three-module test then completed successfully in 6.4 seconds.
- The first pnpm install failed under the required strict peer policy: Vite 8.1.5 resolved Rolldown 1.1.5, whose optional WASM chain needs Emnapi 2.x while the optional binding itself pins Emnapi 1.11.1. Adding top-level Emnapi 2.0.0-alpha.3 packages did not change that nested peer context and was reverted. The broken optional WASM fallback was excluded while native Windows/macOS/Linux bindings remain enabled and strict peer checks remain active.
- The first frontend format check found two source/config files that required the pinned Prettier rewrite. After the build generated `apps/web/dist`, a second check also revealed that a filtered workspace script does not automatically discover the root `.prettierignore`; the script now names the root ignore file explicitly and uses workspace-wide output globs.

### Decision Log

- **2026-07-29 — Accepted:** treat `E:\Git\开源\agent\corvus` as the permanent repository root. The relocation is already present; do not repeat it. Use the revised read-only Milestone 1 audit as the evidence boundary.
- **2026-07-29 — Accepted:** use three Go modules—Core, Launcher, and Agent—with no root module. This preserves dependency and architectural boundaries.
- **2026-07-29 — Accepted:** use `go test ./apps/core/... ./apps/launcher/... ./agent/...` as the canonical workspace test command.
- **2026-07-29 — Accepted:** use a real minimal Fyne v2.8.0 Launcher and native three-platform build matrix during Phase 0.
- **2026-07-29 — Accepted:** defer Echo, SQLite, goose, sqlc, OpenAPI generation, Agent/ADK, product UI libraries, and all business behavior to their owning later phases.
- **2026-07-29 — Accepted:** preserve the existing AGPLv3 license unchanged; separate documentation licensing remains an open, non-blocking later decision.
- **2026-07-29 — Accepted:** commit each completed milestone locally after validation; do not push unless the user explicitly requests it.
- **2026-07-29 — Accepted:** retain `strict-peer-dependencies=true` and exclude only Rolldown's broken optional WASM binding via `ignoredOptionalDependencies`; do not weaken peer validation or remove native bindings for the three supported platforms.

### Outcomes & Retrospective

Milestones 1–3 are complete. The three approved Go modules and root skeleton build/test locally on Windows, including a real Fyne v2.8.0 Launcher shell. The pnpm workspace installs from one frozen lockfile and the minimal Web package passes formatting, lint, smoke test, and production build. Milestone 4 is in progress; Milestone 5 is pending. No Phase 1 capability, database migration, or design-document edit has occurred.
