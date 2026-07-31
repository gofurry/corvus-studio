# Corvus Studio Usage & Development Guide

> 文档状态：Phase 2 本地实现已完成，远端三平台验证待执行\
> 项目：Corvus Studio\
> 用途：Project Foundation、本地开发、验证命令和后续 Phase 边界

## 1. 当前实现状态

Corvus Studio 已完成 Phase 0、Phase 1，并完成 Phase 2 Project Foundation 的本地实现：

- `corvus serve` 长期运行命令；
- Viper 配置文件、环境变量和 CLI 覆盖；
- Zap 控制台/文件结构化日志与 lumberjack 轮转；
- Echo v5 本地 HTTP runtime、`GET /healthz` 和优雅关闭；
- `modernc.org/sqlite`、嵌入式 goose migration 和 sqlc 查询边界；
- Project Domain、UUIDv7、目录校验、SQLite repository 和重启持久化；
- OpenAPI-first Project create/list/get 接口及生成的 Go/TypeScript 客户端边界；
- React Router、TanStack Query 和 Ant Design 实现的项目列表、创建和详情流程；
- Vite 开发代理与可选的生产 Web 嵌入构建；
- Core 配置、日志、存储、HTTP、OpenAPI 合约和生命周期测试。

GitHub Actions [run 30611808221](https://github.com/gofurry/corvus-studio/actions/runs/30611808221) 已针对 Phase 1 修正提交 `59a478c` 通过 Go quality、Frontend quality、Ubuntu、macOS 和 Windows 原生任务。Phase 2 本轮不推送，因此不能把该历史运行当作 Phase 2 的远端证据。

仍未实现：Project 编辑/删除/归档、Dashboard、Release Goal、Checklist、Resource、Deliverable、Steam 模板、Asset Map、Watch、SSE、ADK Agent、模型 Provider、登录鉴权、系统目录打开、Docker、systemd、安装包、签名、公证和自动更新。

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
{ "status": "ok", "database": "ok", "schema_version": 2 }
```

当前 schema version 是 `2`。`/healthz` 是运行状态接口，不是产品 API；Project API 位于 `/api/v1/projects`。

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

| Method | Path                            | 成功结果                             |
| ------ | ------------------------------- | ------------------------------------ |
| `GET`  | `/api/v1/projects`              | `200`，按创建时间倒序返回 `items`    |
| `POST` | `/api/v1/projects`              | `201`、Project 和 `Location` 响应头  |
| `GET`  | `/api/v1/projects/{project_id}` | `200`，返回指定 Project 的持久化详情 |

创建示例：

```powershell
$projectLocation = (Resolve-Path .).Path
$projectBody = @{
  name = 'Corvus Studio'
  description = 'Phase 2 local project'
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

`stage` 只接受 `concept`、`development`、`release_preparation`、`released`。创建状态固定为 `active`。错误使用统一 `error` 信封，当前错误码是 `validation_failed`、`project_not_found`、`project_location_conflict` 和 `internal_error`。请求体上限是 1 MiB，未知字段、尾随 JSON 和非 UUIDv7 ID 会被拒绝。

重复的规范目录返回 `409 project_location_conflict`，不会新增记录或写入目标目录。“打开项目”仅指进入 Corvus 详情页，不会调用系统文件管理器。

## 7. Web 开发与生产嵌入

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

打开 `http://localhost:5173/projects` 可以查看项目列表、进入创建表单和打开应用内详情页。创建时 `location` 必须是已有绝对目录；Core 只读取和规范化目录元数据，不会创建、移动或修改项目文件。

生产嵌入验证：

```powershell
pnpm build
go run ./apps/core/cmd/stage-web
New-Item -ItemType Directory -Force .tmp | Out-Null
go build -tags corvus_webui -o .tmp/corvus-embedded.exe ./apps/core/cmd/corvus
```

`stage-web` 只允许写入 Core module 内已忽略的 `apps/core/internal/interfaces/webui/dist/`。带 `corvus_webui` 标签的二进制使用 `go:embed` 提供静态资源和 SPA fallback；普通 Go build 不依赖 Node 或已有 `dist/`。

## 8. SQLite、goose、sqlc 与 OpenAPI 生成

Core 使用纯 Go 的 `modernc.org/sqlite`，启动时启用 WAL、foreign keys、busy timeout 和 normal synchronous。goose 在 HTTP 监听前执行嵌入式 migration。对于已有且非空、存在待执行 migration 的数据库，Core 会先在数据库同级 `backups/` 中创建 SQLite 快照；快照失败会阻止 migration 和启动。

数据库包含基础设施表 `corvus_runtime_metadata` 和 Phase 2 `projects` STRICT 表。`projects.location_key` 只用于平台一致的唯一性判断，不通过公开 API 返回。

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

## 9. Validation

### Go workspace

根目录没有 `go.mod`，因此不要把根级 `go test ./...` 当作三 module 全量验证。规范命令是：

```powershell
go test ./apps/core/... ./apps/launcher/... ./agent/...
golangci-lint run ./apps/core/... ./apps/launcher/... ./agent/...
```

检查格式：

```powershell
$phase2GoFiles = Get-ChildItem -Recurse -File -Include *.go -Path '.\apps\core','.\apps\launcher','.\agent'
$phase2Unformatted = $phase2GoFiles | ForEach-Object { gofmt -l $_.FullName }
if ($phase2Unformatted) { throw "Unformatted Go files: $phase2Unformatted" }
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

GitHub Actions 在 Windows、macOS、Linux 原生 runner 上执行包含 Project repository/API 持久化场景的 Core 测试，并构建 Core/Launcher。Phase 2 的工作流定义已经同步，但必须等当前提交实际推送并运行后才能声明三平台通过。交叉编译、打包、签名、公证和发布不属于 Phase 2。

## 10. Launcher 与 Agent 边界

```powershell
go run ./apps/launcher
```

Launcher 仍只显示最小 Fyne 窗口，不启停 Core、不提供托盘功能，也不承载业务 UI。`agent/` 仍是独立 Go module 边界；Phase 2 没有 ADK、模型 Provider、workflow、tool 或 prompt 实现。

## 11. Contribution Rules

- 使用 Trunk Based Development；
- 不强制 Conventional Commits；
- 不创建 Issue Template 或 Pull Request Template；
- 当前不创建 `CONTRIBUTING.md`；
- 每个逻辑变更都创建本地提交；
- 每次变更都必须提供可运行验证证据；
- 不得在没有证据时声称测试、构建或远端 CI 通过。

## 12. Troubleshooting

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

### Launcher 无法构建

确认已安装当前平台的 [Fyne prerequisites](https://docs.fyne.io/started/)，并检查 `go env CGO_ENABLED` 与本机 C 编译器。

### pnpm 安装不一致

不要在子目录单独安装或生成锁文件。回到仓库根目录执行：

```powershell
pnpm install --frozen-lockfile
```

如果 lockfile 与 manifest 不一致，不要删除 lockfile 或关闭严格检查；先核对依赖变更。

## 13. 文档导航

1. `README.md`：当前能力和快速验证。
2. `USAGE.md`：Core、Web、存储、验证和排错。
3. `.agent/phase-2-project-foundation.md`：Phase 2 执行计划、进度和证据。
4. `docs/roadmap.md`：阶段状态和验收门槛。
5. `docs/`：产品、架构、API、数据模型、Agent 和工程设计的信息源。
