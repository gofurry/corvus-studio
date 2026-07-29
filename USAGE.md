# Corvus Studio Usage & Development Guide

> 文档状态：Phase 0 已实现命令与未来产品目标的分界说明\
> 项目：Corvus Studio\
> 用途：当前仓库验证、本地开发入门和未来能力导航

## 1. 当前实现状态

Corvus Studio 当前处于 **Phase 0: Repository Bootstrap**。现在可以验证 monorepo、三个 Go module、React/Vite 前端骨架、最小 Fyne Launcher 和 CI 定义。

当前尚未实现：

- `corvus serve` 或 HTTP API；
- Echo、SQLite、goose、sqlc 和数据库迁移；
- OpenAPI 合约与生成客户端；
- ADK Agent、模型 Provider 或 AI 配置；
- Project、Release Goal、Checklist、Resource、Deliverable、Steam 模板、Asset Map 和 Watch；
- 登录、鉴权、Docker、systemd、安装包、签名、公证和自动更新。

本文档中只有标明为“Phase 0 可用”的命令才能当作已实现功能。

## 2. 产品目标（未来 Phase）

Corvus Studio 的产品目标是成为本地优先的独立游戏 Steam 发布准备工作台。设计文档规划了 Desktop Application、Local Web Application、Project 管理、Steam Coming Soon 准备、Checklist、Resource Library、Deliverables 和 AI Review。

这些是后续阶段的产品目标，不是 Phase 0 已可用能力。AI 仍按设计作为增强能力，不是未来产品的运行前提。

## 3. Developer Environment

Phase 0 固定使用：

- Go 1.26.5；
- Node.js 24.15.x；
- pnpm 10.11.0；
- golangci-lint 2.12.2（执行本地 Go lint 时）；
- 构建 Launcher 所需的本机 C 编译器与 Fyne 图形开发依赖。

在仓库根目录检查：

```powershell
go version
node --version
pnpm --version
golangci-lint --version
```

Windows、macOS 和 Linux 都是正式目标平台。Docker 不是 Phase 0 开发环境前提。

## 4. Repository Setup（Phase 0 可用）

获取代码并进入仓库根目录：

```powershell
git clone <repository>
cd corvus-studio
```

确认 Go workspace 指向根 `go.work` 且只聚合 Core、Launcher 和 Agent：

```powershell
go env GOWORK
go work edit -json
go list -m
```

安装前端 workspace 依赖：

```powershell
pnpm install --frozen-lockfile
```

仓库只使用根目录的 `pnpm-lock.yaml`。不要在 `apps/web` 或 `packages/*` 中创建独立锁文件。

## 5. Local Development（Phase 0 可用）

### Core bootstrap

```powershell
go run ./apps/core/cmd/corvus
```

预期输出：

```text
Corvus Studio core bootstrap
```

该命令随后退出；Phase 0 不启动 HTTP 服务。

### Web development server

```powershell
pnpm dev:web
```

Web 由 Vite Dev Server 独立运行，尚未连接 Core。

### Minimal Launcher

```powershell
go run ./apps/launcher
```

Launcher 只显示一个最小 Fyne 窗口；不启停 Core、不提供托盘功能，也不承载业务 UI。

## 6. Validation（Phase 0 可用）

### Go

`go test ./...` 不能在无根 `go.mod` 的三 module workspace 中表达完整覆盖范围。请始终使用显式命令：

```powershell
go test ./apps/core/... ./apps/launcher/... ./agent/...
golangci-lint run ./apps/core/... ./apps/launcher/... ./agent/...
```

检查 Go 格式：

```powershell
$phase0GoFiles = Get-ChildItem -Recurse -File -Include *.go -Path '.\apps\core','.\apps\launcher','.\agent'
$phase0Unformatted = $phase0GoFiles | ForEach-Object { gofmt -l $_.FullName }
if ($phase0Unformatted) { throw "Unformatted Go files: $phase0Unformatted" }
```

### Frontend

```powershell
pnpm install --frozen-lockfile
pnpm format:check
pnpm lint
pnpm test
pnpm build
```

`pnpm build` 成功后可观察到 `apps/web/dist/index.html`。`dist/` 是本地构建产物，不提交。

### Native build

```powershell
New-Item -ItemType Directory -Force .tmp | Out-Null
go build -o .tmp/corvus ./apps/core/cmd/corvus
go build -o .tmp/corvus-launcher ./apps/launcher
```

GitHub Actions 会在 Windows、macOS 和 Linux 原生 runner 上构建这两个入口。Phase 0 不进行交叉编译、打包、签名、公证或产物发布。

## 7. Frontend Development Boundary

Phase 0 已安装 React 19、TypeScript、Vite 8、SCSS、ESLint、Prettier、Vitest 和 Testing Library。

Ant Design、Zustand、TanStack Query、React Router、React Flow 和 `react-markdown` 会在首次有真实产品需求的后续 Phase 引入。`packages/api-client`、`packages/shared-types` 和 `packages/ui` 当前只有边界说明，不包含虚假导出或生成代码。

开发模式下 Web 与 Core 独立运行。Vite 产物未来由 Go `embed` 提供，但 Phase 0 只验证 Vite 自身能构建。

## 8. Database, API, Agent and Configuration（未来 Phase）

以下设计仍是后续实现输入，Phase 0 不提供对应命令：

- SQLite、goose migration、sqlc 查询生成；
- Echo HTTP runtime、SSE 和 OpenAPI-first API；
- Viper 运行时配置、zap/lumberjack 日志；
- ADK Agent runtime、workflows、tools、context、prompts 和 providers。

因此，不要在 Phase 0 执行 `corvus migrate`、`goose create`、`sqlc generate` 或 Agent 调试命令。`agent/` 目前只是一个独立 Go module 边界。

## 9. Deployment and Packaging（未来 Phase）

Docker、systemd、正式桌面安装包、macOS 签名/公证和自动更新都不属于 Phase 0。`deployments/` 目录仅预留仓库边界，当前没有 `docker compose up` 或正式安装命令。

## 10. Contribution

- 使用 Trunk Based Development；
- 不强制 Conventional Commits；
- 不创建 Issue Template 或 Pull Request Template；
- 当前不创建 `CONTRIBUTING.md`；
- 每次变更都必须提供可运行的验证证据；
- 不得在没有证据时声称测试或构建通过。

## 11. Troubleshooting

### Go workspace 未识别

在仓库根目录执行：

```powershell
go env GOWORK
go work edit -json
```

`GOWORK` 应指向根 `go.work`，`Use` 条目应且只应包含 `agent`、`apps/core` 和 `apps/launcher`。

### Launcher 无法构建

确认已安装当前平台的 [Fyne prerequisites](https://docs.fyne.io/started/)，并检查 `go env CGO_ENABLED` 与本机 C 编译器。

### pnpm 安装不一致

不要在子目录单独安装或生成锁文件。回到仓库根目录执行：

```powershell
pnpm install --frozen-lockfile
```

如果锁文件与 manifest 不一致，不要删除锁文件或关闭严格 peer dependency 检查；先记录错误并核对依赖变更。

## 12. 文档导航

1. `README.md`：Phase 0 现状和快速验证。
2. `USAGE.md`：当前命令、边界和排错。
3. `.agent/phase-0-repository-bootstrap.md`：可执行计划、进度和证据。
4. `docs/`：产品、架构、API、数据模型、Agent 和工程设计的信息源。
