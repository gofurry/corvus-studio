# Corvus Studio Usage & Development Guide

> 文档状态：Phase 3 Release + Checklist 已完成本地实施，等待远端三平台验证\
> 项目：Corvus Studio\
> 用途：Project、Release、Checklist、本地开发、验证命令和后续 Phase 边界

## 1. 当前实现状态

Corvus Studio 已完成 Phase 0、Phase 1、Phase 2，并完成 Phase 3 Release + Checklist 的本地实现：

- `corvus serve` 长期运行命令；
- Viper 配置文件、环境变量和 CLI 覆盖；
- Zap 控制台/文件结构化日志与 lumberjack 轮转；
- Echo v5 本地 HTTP runtime、`GET /healthz` 和优雅关闭；
- `modernc.org/sqlite`、嵌入式 goose migration 和 sqlc 查询边界；
- Project Domain、UUIDv7、目录校验、SQLite repository 和重启持久化；
- OpenAPI-first Project create/list/get 接口及生成的 Go/TypeScript 客户端边界；
- 由本地 Core 打开的原生目录选择器，以及不可用时的手动绝对路径回退；
- Release Goal、Checklist、状态历史、原子创建和重启持久化；
- 内嵌且经过启动校验的 Steam Coming Soon 模板 v1.0.0，包含 8 个 Steam 官方任务和 4 个 Corvus 推荐任务；
- OpenAPI-first Template/Release/Checklist 接口及同步生成的 Go/TypeScript 类型；
- React Router、TanStack Query 和 Ant Design 实现的项目、Release 状态及可恢复筛选条件的 Checklist 流程；
- Vite 开发代理与可选的生产 Web 嵌入构建；
- Core 配置、日志、存储、HTTP、OpenAPI 合约和生命周期测试。

GitHub Actions [run 30617549600](https://github.com/gofurry/corvus-studio/actions/runs/30617549600) 已针对 Phase 2 基线提交 `80dca49` 通过 Go quality、Frontend quality、Ubuntu、macOS 和 Windows 原生任务。Phase 3 只有当前 Windows 本地证据；尚未推送，不能把该运行描述成覆盖 Phase 3。

仍未实现：Project 编辑/删除/归档、Dashboard、Checklist 任务编辑/删除/排序、模板升级合并、Resource、Deliverable、Evidence、Steam API/在线模板更新、Asset Map、Watch、SSE、ADK Agent、模型 Provider、登录鉴权、系统目录打开、Docker、systemd、安装包、签名、公证和自动更新。

## 2. Developer Environment

项目固定使用：

- Go 1.26.5；
- Node.js 24.15.x；
- pnpm 10.11.0；
- golangci-lint 2.12.2；
- 构建 Launcher 所需的本机 C 编译器与 Fyne 图形开发依赖。

在仓库根目录检查：

```powershell
go version
go env GOWORK
node --version
pnpm --version
golangci-lint --version
```

Windows、macOS 和 Linux 都是正式目标平台。Docker 不是当前开发环境前提。

## 3. Repository Setup

```powershell
git clone <repository>
cd corvus-studio
pnpm install --frozen-lockfile
```

确认根 `go.work` 只聚合 Core、Launcher 和 Agent：

```powershell
go env GOWORK
go work edit -json
go list -m
```

仓库只使用根目录的 `pnpm-lock.yaml`。不要在 `apps/web` 或 `packages/*` 中生成独立锁文件。

## 4. 运行 Core

推荐在开发时显式指定仓库内已忽略的数据目录：

```powershell
go run ./apps/core/cmd/corvus serve --data-dir .tmp/core-runtime
```

启动成功后可以观察到：

- Core 监听 `127.0.0.1:8765`；
- `.tmp/core-runtime/corvus.db` 存在；
- `.tmp/core-runtime/logs/corvus.log` 包含 JSON 日志；
- 首次启动自动执行嵌入式 migration；
- Ctrl+C 或 SIGTERM 触发 HTTP 优雅关闭和数据库关闭。

检查 readiness：

```powershell
Invoke-RestMethod http://127.0.0.1:8765/healthz
```

预期字段：

```json
{ "status": "ok", "database": "ok", "schema_version": 3 }
```

当前 schema version 是 `3`。`/healthz` 是运行状态接口，不是产品 API；公开产品 API 位于 `/api/v1`。

### CLI 参数

```powershell
go run ./apps/core/cmd/corvus serve --help
```

当前参数：

- `--config`：显式 YAML/YML/JSON 配置文件；
- `--data-dir`：数据库、日志和 migration 备份的基础目录；
- `--host`：监听主机，Phase 1 仅接受 loopback；
- `--port`：监听端口；
- `--database`：覆盖 SQLite 文件路径；
- `--log-file`：覆盖轮转日志路径；
- `--log-level`：`debug`、`info`、`warn`、`error` 等 Zap 级别。

## 5. 配置规则

优先级固定为：

```text
CLI > CORVUS_* 环境变量 > 配置文件 > 默认值
```

显式配置示例位于 `apps/core/configs/corvus.example.yaml`：

```powershell
go run ./apps/core/cmd/corvus serve --config .\apps\core\configs\corvus.example.yaml
```

未传 `--config` 时，Core 会在操作系统用户配置目录下的 `corvus-studio/` 中依次查找 `config.yaml`、`config.yml`、`config.json`；不存在属于正常情况，不会自动写入文件。

常用环境变量：

```powershell
$env:CORVUS_RUNTIME_DATA_DIR = '.tmp/core-runtime'
$env:CORVUS_SERVER_HOST = '127.0.0.1'
$env:CORVUS_SERVER_PORT = '8765'
$env:CORVUS_STORAGE_PATH = '.tmp/core-runtime/corvus.db'
$env:CORVUS_LOGGING_PATH = '.tmp/core-runtime/logs/corvus.log'
$env:CORVUS_LOGGING_LEVEL = 'debug'
go run ./apps/core/cmd/corvus serve
```

默认数据目录由 Go 的 `os.UserConfigDir()` 决定，并追加 `corvus-studio`。默认数据库是 `<data_dir>/corvus.db`，默认日志是 `<data_dir>/logs/corvus.log`。未实现鉴权前，非 loopback 地址会被拒绝。

## 6. Project API

公开合约的唯一来源是 `apps/core/openapi/openapi.yaml`：

| Method | Path                              | 成功结果                             |
| ------ | --------------------------------- | ------------------------------------ |
| `GET`  | `/api/v1/projects`                | `200`，按创建时间倒序返回 `items`    |
| `POST` | `/api/v1/projects`                | `201`、Project 和 `Location` 响应头  |
| `GET`  | `/api/v1/projects/{project_id}`   | `200`，返回指定 Project 的持久化详情 |
| `POST` | `/api/v1/system/select-directory` | `200`，返回选择的绝对目录或取消状态  |

创建示例：

```powershell
$projectLocation = (Resolve-Path .).Path
$projectBody = @{
  name = 'Corvus Studio'
  description = 'Phase 3 local project'
  location = $projectLocation
  steam_app_id = 480
  language = 'English'
  stage = 'concept'
} | ConvertTo-Json

$project = Invoke-RestMethod -Method Post `
  -Uri http://127.0.0.1:8765/api/v1/projects `
  -ContentType 'application/json' `
  -Body $projectBody
Invoke-RestMethod http://127.0.0.1:8765/api/v1/projects
Invoke-RestMethod "http://127.0.0.1:8765/api/v1/projects/$($project.id)"
```

`stage` 只接受 `concept`、`development`、`release_preparation`、`released`。创建状态固定为 `active`。错误使用统一 `error` 信封，当前错误码是 `validation_failed`、`project_not_found`、`project_location_conflict`、`directory_picker_unavailable` 和 `internal_error`。请求体上限是 1 MiB，未知字段、尾随 JSON 和非 UUIDv7 ID 会被拒绝。

重复的规范目录返回 `409 project_location_conflict`，不会新增记录或写入目标目录。“打开项目”仅指进入 Corvus 详情页，不会调用系统文件管理器。

目录选择接口只接受 JSON `{"purpose":"project_location"}`，避免普通跨站请求触发本机对话框。用户取消时返回 `selected: false`，不属于错误。Windows 使用系统 FolderBrowserDialog，macOS 使用系统目录选择器；Linux 优先使用 `zenity`，其次使用 `kdialog`。两者都不存在时 Web 会提示用户继续手动输入。

## 7. Release 与 Checklist API

Phase 3 的公开接口同样以 `apps/core/openapi/openapi.yaml` 为唯一来源：

| Method | Path                                       | 行为                                     |
| ------ | ------------------------------------------ | ---------------------------------------- |
| `GET`  | `/api/v1/release-templates/{template_key}` | 读取内置模板 metadata 和任务预览         |
| `GET`  | `/api/v1/projects/{project_id}/releases`   | 列出 Project 的 Release Goals            |
| `POST` | `/api/v1/releases`                         | 原子创建 Release Goal 和完整 Checklist   |
| `GET`  | `/api/v1/releases/{release_id}`            | 读取 Goal 和 Checklist 汇总              |
| `POST` | `/api/v1/releases/{release_id}/transition` | 显式推进或退回 Release 状态              |
| `GET`  | `/api/v1/releases/{release_id}/checklist`  | 按 status、source、category 筛选任务     |
| `POST` | `/api/v1/checklist`                        | 添加服务端标记为 `user` 来源的自定义任务 |
| `GET`  | `/api/v1/checklist/{item_id}`              | 读取任务详情                             |
| `POST` | `/api/v1/checklist/{item_id}/transition`   | 修改任务状态                             |

通过 API 创建工作区：

```powershell
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
$checklist = Invoke-RestMethod "http://127.0.0.1:8765/api/v1/releases/$($release.id)/checklist"
```

一次创建会在同一事务内写入 Goal 和全部 12 个模板任务；重复 Goal 返回 `409 release_goal_conflict`，不会增加记录。模板是项目内独立副本，运行时不会访问 Steam 网络，也不会在模板版本变化时自动合并。

Checklist 状态为 `not_started`、`in_progress`、`needs_review`、`done`、`blocked`、`not_applicable`。Required 任务不能标记为 `not_applicable`；所有 Required 任务都为 Done 后，Release 才能从 `preparing` 进入 `ready_for_review`。进入 `ready_for_review`、`ready` 或 `submitted` 后，Required 任务被锁定，必须先把 Release 退回准备状态。相同目标状态是幂等操作，不重复写历史。

工作流错误继续使用统一信封，错误码为 `template_not_found`、`release_not_found`、`release_goal_conflict`、`release_not_ready`、`release_state_conflict`、`checklist_item_not_found`、`invalid_transition`、`validation_failed` 和 `internal_error`。冲突与 readiness gate 返回 `409`；失败操作不留下半完成数据。

## 8. Web 开发与生产嵌入

开发模式使用两个独立进程。在第一个终端运行 Core：

```powershell
go run ./apps/core/cmd/corvus serve --data-dir .tmp/core-runtime
```

在第二个终端运行 Vite：

```powershell
pnpm dev:web
```

Vite 默认把 `/api` 和 `/healthz` 代理到 `http://127.0.0.1:8765`。Core 使用其他端口时，在启动 Vite 前设置：

```powershell
$env:CORVUS_CORE_URL = 'http://127.0.0.1:18765'
pnpm dev:web
```

打开 `http://localhost:5173/projects` 可以查看项目列表、进入创建表单和打开应用内详情页。创建时点击 **Browse** 选择已有目录，或手动输入绝对路径；Core 只读取和规范化目录元数据，不会创建、移动或修改项目文件。

进入 Project 后：

1. 打开 **Release**，查看模板版本、来源和 Steam/Corvus 任务数量；
2. 确认 **Generate release workspace**，一次性创建 Goal 和 Checklist；
3. 打开 **Checklist**，使用 status/source/category 筛选任务；
4. 选择任务查看 What、Why、Requirement、Source 和 Status，或添加自定义任务；
5. 在 Release 页面根据 Required 完成度显式推进发布状态。

筛选条件和选中任务写入 URL query，直接刷新 `/projects/:projectId/checklist` 可以恢复上下文。窄屏使用任务详情 Drawer。当前不包含 Evidence、Resource、Deliverable、文件验证或 Steam 自动提交。

生产嵌入验证：

```powershell
pnpm build
go run ./apps/core/cmd/stage-web
New-Item -ItemType Directory -Force .tmp | Out-Null
go build -tags corvus_webui -o .tmp/corvus-embedded.exe ./apps/core/cmd/corvus
```

`stage-web` 只允许写入 Core module 内已忽略的 `apps/core/internal/interfaces/webui/dist/`。带 `corvus_webui` 标签的二进制使用 `go:embed` 提供静态资源和 SPA fallback；普通 Go build 不依赖 Node 或已有 `dist/`。

## 9. SQLite、goose、sqlc 与 OpenAPI 生成

Core 使用纯 Go 的 `modernc.org/sqlite`，启动时启用 WAL、foreign keys、busy timeout 和 normal synchronous。goose 在 HTTP 监听前执行嵌入式 migration。对于已有且非空、存在待执行 migration 的数据库，Core 会先在数据库同级 `backups/` 中创建 SQLite 快照；快照失败会阻止 migration 和启动。

数据库包含基础设施表、`projects`，以及 Phase 3 的 `release_goals`、`checklist_items` 和两类状态历史 STRICT 表。Goal 与模板任务在一个事务内创建；实体状态变化与历史记录也在同一事务内写入。`projects.location_key` 只用于平台一致的唯一性判断，不通过公开 API 返回。

验证 sqlc 生成结果：

```powershell
Push-Location .\apps\core
go tool oapi-codegen -config openapi/oapi-codegen.yaml openapi/openapi.yaml
go tool sqlc generate
go tool sqlc vet
Pop-Location

pnpm --filter @corvus-studio/api-client generate
git diff --exit-code -- apps/core/internal/interfaces/httpserver/api apps/core/internal/infrastructure/storage/sqlc packages/api-client/src/generated
```

没有 `corvus migrate`、`corvus backup` 或 `corvus doctor` 命令。不要用这些未来命令替代自动启动 migration。

## 10. Validation

### Go workspace

根目录没有 `go.mod`，因此不要把根级 `go test ./...` 当作三 module 全量验证。规范命令是：

```powershell
go test ./apps/core/... ./apps/launcher/... ./agent/...
golangci-lint run ./apps/core/... ./apps/launcher/... ./agent/...
```

检查格式：

```powershell
$phase3GoFiles = Get-ChildItem -Recurse -File -Include *.go -Path '.\apps\core','.\apps\launcher','.\agent'
$phase3Unformatted = $phase3GoFiles | ForEach-Object { gofmt -l $_.FullName }
if ($phase3Unformatted) { throw "Unformatted Go files: $phase3Unformatted" }
```

### Frontend

```powershell
pnpm install --frozen-lockfile
pnpm format:check
pnpm lint
pnpm test
pnpm build
```

成功后 `apps/web/dist/index.html` 存在；`dist/` 不提交。

### Native builds

```powershell
New-Item -ItemType Directory -Force .tmp | Out-Null
go build -o .tmp/corvus.exe ./apps/core/cmd/corvus
go build -o .tmp/corvus-launcher.exe ./apps/launcher
```

GitHub Actions 在 Windows、macOS、Linux 原生 runner 上执行全部 Core 测试，其中包含 Project、Release 和 Checklist 的 migration、原子性、状态历史及重启持久化场景，并构建 Core/Launcher。现有生成漂移步骤自动覆盖扩展后的 OpenAPI、sqlc 和 TypeScript client，因此 Phase 3 不需要新增平行工作流。必须等当前提交实际推送并运行后才能声明三平台通过。交叉编译、打包、签名、公证和发布不属于 Phase 3。

## 11. Launcher 与 Agent 边界

```powershell
go run ./apps/launcher
```

Launcher 仍只显示最小 Fyne 窗口，不启停 Core、不提供托盘功能，也不承载业务 UI。`agent/` 仍是独立 Go module 边界；Phase 3 没有 ADK、模型 Provider、workflow、tool 或 prompt 实现。

## 12. Contribution Rules

- 使用 Trunk Based Development；
- 不强制 Conventional Commits；
- 不创建 Issue Template 或 Pull Request Template；
- 当前不创建 `CONTRIBUTING.md`；
- 每个逻辑变更都创建本地提交；
- 每次变更都必须提供可运行验证证据；
- 不得在没有证据时声称测试、构建或远端 CI 通过。

## 13. Troubleshooting

### Core 拒绝 host

当前没有鉴权，只允许 `localhost` 或 loopback IP。不要通过绑定 `0.0.0.0` 绕过这一安全边界。

### 端口已占用

```powershell
go run ./apps/core/cmd/corvus serve --port 18765 --data-dir .tmp/core-runtime
```

同时启动 Vite 时，将 `CORVUS_CORE_URL` 指向相同端口。

### SQLite 启动失败

先检查配置中数据库路径是否误指向目录，以及父目录是否可写。不要删除现有数据库或 migration 备份来掩盖错误。

### Project location 被拒绝

确认路径是当前平台上的已有绝对目录，而不是相对路径、文件或不存在的位置。符号链接会解析到真实目录；同一真实目录只能创建一个 Project。不要为了通过校验而让 Corvus 自动创建用户目录。

### Browse 无法打开目录选择器

Windows 需要系统 PowerShell，macOS 需要 `osascript`。Linux 安装 `zenity` 或 `kdialog` 后重试。即使原生选择器不可用，Project location 输入框仍可手动填写已有目录的绝对路径。

### Launcher 无法构建

确认已安装当前平台的 [Fyne prerequisites](https://docs.fyne.io/started/)，并检查 `go env CGO_ENABLED` 与本机 C 编译器。

### pnpm 安装不一致

不要在子目录单独安装或生成锁文件。回到仓库根目录执行：

```powershell
pnpm install --frozen-lockfile
```

如果 lockfile 与 manifest 不一致，不要删除 lockfile 或关闭严格检查；先核对依赖变更。

## 14. 文档导航

1. `README.md`：当前能力和快速验证。
2. `USAGE.md`：Core、Web、存储、验证和排错。
3. `.agent/phase-3-release-checklist.md`：Phase 3 执行计划、进度和证据。
4. `docs/roadmap.md`：阶段状态和验收门槛。
5. `docs/`：产品、架构、API、数据模型、Agent 和工程设计的信息源。
