# Corvus Studio Roadmap

> 状态：Active\
> 最近更新：2026-07-31\
> 当前进度：Phase 0、Phase 1、Phase 2 已完成；下一阶段为 Phase 3\
> 首个产品目标：Corvus Studio Launch `v0.1` Alpha

本文从[开发实施计划 v0.1](development/Corvus_Studio_Development_Implementation_Plan_v0.1.md)提取阶段主线，并补充可维护的完成状态、依赖关系和验收门槛。原始设计文档继续作为范围基线；本文负责反映实际进展，不以目录占位或未验证代码代替完成证据。

## 状态说明

- ✅ **Completed**：范围内任务已完成，并存在可复现的验证证据。
- 🚧 **In progress**：已有受控实施工作，但阶段验收尚未全部通过。
- ⬜ **Not started**：尚未开始实现；已有目录骨架不计为业务进展。
- ⏸️ **Blocked**：存在无法在当前权限或环境内解除的阻塞。

## 当前进度

| 阶段                           | 状态           | 完成度 | 结果或下一步                                                           |
| ------------------------------ | -------------- | -----: | ---------------------------------------------------------------------- |
| Phase 0 — Repository Bootstrap | ✅ Completed   |   100% | Monorepo、三 Go module、pnpm workspace、最小应用骨架和三平台 CI 已验证 |
| Phase 1 — Core Runtime         | ✅ Completed   |   100% | Core runtime、本地验证及 Windows、macOS、Linux CI 均已通过             |
| Phase 2 — Project Foundation   | ✅ Completed   |   100% | Project 创建、选择目录、列表、详情、持久化和三平台基线 CI 已验证       |
| Phase 3 — Release + Checklist  | ⬜ Not started |     0% | 建立 Steam 发布目标与任务闭环                                          |
| Phase 4A — Resource            | ⬜ Not started |     0% | 建立素材与引用管理                                                     |
| Phase 4B — Deliverable         | ⬜ Not started |     0% | 建立交付物及其与任务、资源的关系                                       |
| Phase 5 — Asset Map            | ⬜ Not started |     0% | 可视化并保存关系图                                                     |
| Phase 6 — Agent MVP            | ⬜ Not started |     0% | 在确定性业务基础上接入 ADK 与首批工作流                                |
| Phase 7 — Experience           | ⬜ Not started |     0% | 完善引导、导出、备份和离线体验                                         |
| Phase 8 — Alpha Release        | ⬜ Not started |     0% | 发布给 20–50 名独立游戏开发者验证                                      |

## 依据与冲突处理

实施顺序和阶段边界以本路线图的源文档——[开发实施计划 v0.1](development/Corvus_Studio_Development_Implementation_Plan_v0.1.md)——为最高优先级。具体结构、技术和验收细节依次参考：

1. [仓库结构设计](architecture/Corvus_Studio_Repository_Structure_Design_v0.1.md)；
2. [技术栈决策](architecture/Corvus_Studio_Technology_Stack_Decision_v0.1.md)；
3. [工程与部署指南](deployment/Corvus_Studio_Engineering_and_Deployment_Guide_v0.1.md)；
4. 根目录 `USAGE.md`；
5. 其他产品、架构、API 和 Agent 文档。

[Launch MVP 路线图](product/Corvus_Studio_Launch_MVP_Roadmap_v0.1.md)使用了另一套 Phase 编号，并把 SQLite 写入其 Phase 0。这里不采用该编号：SQLite、Echo、goose、sqlc 等 Core Runtime 能力归入 Phase 1，Phase 0 保持纯 Repository Bootstrap。产品路线图仍作为用户价值和后续版本背景。

## 实施路线

### Phase 0 — Repository Bootstrap

**状态：✅ Completed（2026-07-29）**

**目标：** 建立可构建、可测试、可在 Windows、macOS、Linux CI 中验证的仓库基础，不实现产品业务。

**已完成：**

- [x] 建立 Monorepo 顶层结构和长期工程规则。
- [x] 建立根 `go.work`，聚合 `apps/core`、`apps/launcher`、`agent` 三个 Go module。
- [x] 建立 pnpm workspace、唯一 lockfile、根命令和 `apps/web` 最小 React/Vite 工程。
- [x] 建立可编译的 Core、Web、Fyne Launcher 和空 Agent 边界。
- [x] 建立 Go/frontend 质量检查及 Windows、macOS、Linux 原生构建工作流。
- [x] 同步 `README.md`、`USAGE.md` 和许可证状态。
- [x] 验证未引入 Echo、SQLite、goose、sqlc、ADK 或产品领域实现。

**完成证据：**

- 本地执行计划和证据日志见 [Phase 0 ExecPlan](../.agent/phase-0-repository-bootstrap.md)。
- 规范化验收命令为 `go test ./apps/core/... ./apps/launcher/... ./agent/...`、`pnpm format:check`、`pnpm lint`、`pnpm test`、`pnpm build`。
- GitHub Actions [run 30442490029](https://github.com/gofurry/corvus-studio/actions/runs/30442490029) 的 Go quality、Frontend quality、Ubuntu、macOS、Windows 原生构建均成功。

### Phase 1 — Core Runtime

**状态：✅ Completed（2026-07-31）**

**执行计划：** [Phase 1 Core Runtime ExecPlan](../.agent/phase-1-core-runtime.md)

**目标：** 启动可长期运行、可配置、可观测并具备本地持久化基础的 Corvus Core。

**重点：** Core 进程生命周期、配置、日志、HTTP 服务、SQLite schema 工具链和 CLI；不实现 Project 等产品领域。

**已完成任务：**

- [x] 使用 Cobra 建立 `corvus serve` 命令和清晰的退出码。
- [x] 使用 Viper 建立配置文件、环境变量和命令行覆盖规则。
- [x] 使用 Zap 与 lumberjack 建立结构化日志和本地轮转策略。
- [x] 使用 Echo v5 建立最小 HTTP runtime、错误边界和优雅关闭。
- [x] 使用 `modernc.org/sqlite` 建立无 CGO 的 SQLite 连接与启动检查。
- [x] 建立 goose migration 目录、版本执行和失败恢复规则。
- [x] 建立 sqlc 配置、查询边界和可复现生成命令。
- [x] 为配置、日志、数据库和进程生命周期增加单元/集成测试。

**验收门槛：**

- [x] `corvus serve` 可启动并保持运行；进程 smoke 验证真实 HTTP，生命周期测试验证受控取消和优雅退出。
- [x] 临时数据目录上的首次启动、重复启动、migration 和 migration 前快照均已验证。
- [x] 三个 Go module 的规范测试命令在 Windows 本地通过，Core 不依赖桌面图形环境。
- [x] CI 定义与 `README.md`、`USAGE.md` 使用相同的测试、生成和构建命令。
- [x] 当前 Phase 1 提交在 GitHub Actions 的 Windows、macOS、Linux 原生任务中实际通过。

**当前证据：**

- Phase 1 实现提交为 `2456d78`、`6a00ad1`、`1f49c2f`、`0637980`、边界测试修正 `3abf069` 和 macOS 日志修正 `59a478c`；详细操作与发现见 ExecPlan。
- Windows 本地已通过 Core 默认/嵌入标签测试、三个 Go module 测试、sqlc 生成/vet/漂移、golangci-lint、前端 format/lint/test/build、默认/嵌入 Core build 和真实 `/healthz` 进程 smoke。
- GitHub Actions [run 30611808221](https://github.com/gofurry/corvus-studio/actions/runs/30611808221) 针对提交 `59a478c` 的 Go quality、Frontend quality、Ubuntu、macOS 和 Windows 原生任务全部成功。

### Phase 2 — Project Foundation

**状态：✅ Completed（2026-07-31）**

**执行计划：** [Phase 2 Project Foundation ExecPlan](../.agent/phase-2-project-foundation.md)

**目标：** 完成第一个端到端业务闭环，让用户能够创建并重新打开本地项目。

**重点：** Project Domain、持久化、OpenAPI First 接口和最小 Project 页面。

**已完成任务：**

- [x] 定义 Project 领域模型、UUIDv7 标识、校验规则和错误语义。
- [x] 添加 Project schema/migration、repository 和 service 层。
- [x] 先更新 OpenAPI 合约，再实现 Project API 与生成客户端。
- [x] 实现项目创建、项目列表和打开项目的最小 Web 流程。
- [x] 提供原生目录选择按钮，并在系统 picker 不可用时保留手动路径回退。
- [x] 覆盖领域、存储、API 和前端交互测试。

**验收门槛：**

- [x] Windows 本地实进程验证可创建项目、重启 Core，并从同一 SQLite 重新打开项目。
- [x] OpenAPI、Go models、sqlc 和 TypeScript client 重新生成后不存在漂移。
- [x] 自动化测试验证失败输入和重复目录不会留下半创建项目或修改目标目录。
- [x] Phase 2 基线提交在 GitHub Actions 的 Go/frontend quality 及 Windows、macOS、Linux 原生任务中通过。

**当前证据：**

- 里程碑提交：`88ebc3a`、`2fe1ff6`、`b1a9925`、`670522d`、`6f9ffcd`、`80dca49`；目录选择补充提交为 `695ad88`、`6ad6226`。
- Windows 本地通过 OpenAPI/sqlc/TypeScript client 生成漂移检查、44 个 Go 文件格式检查、三个 module 测试、Go vet、golangci-lint v2.12.2、12 个 Web 交互/API 测试、frontend lint/build、默认及嵌入式 Core build 和 Launcher build。
- 实进程 smoke 创建 UUIDv7 Project，重启同一 Core 后 list/get 仍返回同一记录；嵌入式 `/projects/:id` 刷新返回 `200`，被引用目录保持空目录。
- GitHub Actions [run 30617549600](https://github.com/gofurry/corvus-studio/actions/runs/30617549600) 针对 Phase 2 基线提交 `80dca49` 的全部任务成功；目录选择补充目前只有本地证据，等待下一次推送后的回归运行。

### Phase 3 — Release + Checklist

**状态：⬜ Not started**

**目标：** 建立 Steam 发布准备的基础工作流。

**重点：** Release Goal、Steam Coming Soon 模板和可更新状态的 Checklist。

**计划任务：**

- [ ] 实现 Release Goal 领域、存储、API 和最小页面。
- [ ] 实现 Steam Coming Soon 模板及其版本/来源记录。
- [ ] 从发布目标生成 Checklist，并支持查看、筛选和状态变更。
- [ ] 建立 Release Goal 与 Checklist 的关联和一致性约束。
- [ ] 覆盖模板生成、状态转换、重启持久化和错误恢复测试。

**验收门槛：**

- [ ] 用户可为项目创建发布目标并生成 Checklist。
- [ ] 用户可查看任务、修改状态，重启后结果保持一致。
- [ ] 本阶段不依赖 Agent 即可完成确定性工作流。

### Phase 4A — Resource

**状态：⬜ Not started**

**目标：** 统一管理项目所需的素材、文本和外部引用。

**重点：** 图片、URL、文本、外部引用，以及 managed/referenced/linked 资源语义。

**计划任务：**

- [ ] 实现 Resource 领域、类型校验、版本与元数据模型。
- [ ] 实现 Resource 存储、API、生成客户端和 Resource Library 页面。
- [ ] 明确受管文件、外部文件引用和 URL 链接的生命周期。
- [ ] 覆盖丢失引用、重复导入、无效 URL 和文件边界测试。

**验收门槛：**

- [ ] 用户可新增、查看、修改和删除各类 Resource。
- [ ] 外部引用失效时可观察、可恢复，不静默删除用户数据。

### Phase 4B — Deliverable

**状态：⬜ Not started**

**目标：** 把原始资源转化为可跟踪的 Steam 发布交付物。

**重点：** Main Capsule、Screenshots、Store Copy，以及 `Resource → Deliverable → Checklist` 闭环。

**计划任务：**

- [ ] 实现 Deliverable 领域、版本、状态和验证结果。
- [ ] 实现交付物存储、API、客户端和最小工作区。
- [ ] 建立 Resource、Deliverable、Checklist 的显式关系模型。
- [ ] 为尺寸、格式、缺失来源和关系一致性增加验证与测试。

**验收门槛：**

- [ ] 用户可创建三类首批 Deliverable，并追溯其来源 Resource。
- [ ] Checklist 能反映相关 Deliverable 的存在与状态。
- [ ] 关系更新在重启后保持一致，失败操作可回滚。

### Phase 5 — Asset Map

**状态：⬜ Not started**

**目标：** 可视化项目事实、资源、交付物和任务之间的关系。

**重点：** React Flow、节点展示、关系编辑和布局保存。

**计划任务：**

- [ ] 定义图节点、边和布局持久化模型。
- [ ] 接入 React Flow，展示 Resource、Deliverable、Checklist 等节点。
- [ ] 支持创建、修改和删除允许的关系。
- [ ] 保存并恢复用户布局，不改变底层领域事实。
- [ ] 覆盖循环关系、孤立节点、删除影响和布局恢复测试。

**验收门槛：**

- [ ] 用户可编辑关系并在重启后看到相同图和布局。
- [ ] 图形界面与领域关系保持一致，非法关系被明确拒绝。

### Phase 6 — Agent MVP

**状态：⬜ Not started**

**依赖：** Phase 3 Checklist、Phase 4A Resource、Phase 4B Deliverable 全部完成。

**目标：** 在可靠的确定性产品基础上提供首批可审计的 AI 辅助工作流。

**重点：** Google ADK Go、同进程运行、Root Agent、Tool Registry、Context Provider、DeepSeek 与 Seed provider。

**计划任务：**

- [ ] 在独立 `agent/` module 中建立 ADK runtime 和 Core 内启动边界。
- [ ] 实现 Root Agent、工具注册、上下文装配、超时和取消。
- [ ] 按最小权限实现 Project、Checklist、Resource、Deliverable 工具。
- [ ] 接入 DeepSeek 与 Seed，隔离 provider 配置和凭据。
- [ ] 实现 Release Diagnosis 与 Deliverable Review 工作流。
- [ ] 为建议、证据、工具调用和用户确认建立可追踪记录。
- [ ] 覆盖无模型、超时、限流、工具失败和危险写入确认测试。

**验收门槛：**

- [ ] 两个首批工作流都能基于真实项目上下文给出可追溯结果。
- [ ] Agent 不绕过领域服务直接写数据库，副作用操作需要明确授权。
- [ ] provider 不可用时核心手工工作流仍可使用。

### Phase 7 — Experience

**状态：⬜ Not started**

**目标：** 补齐首次使用、日常效率、数据可携带性和离线反馈。

**重点：** Onboarding、Command Palette、Export、Backup、Offline UI 和 Docker。

**计划任务：**

- [ ] 实现首次启动 Onboarding 和可跳过/重入机制。
- [ ] 实现 Command Palette 与键盘可达的核心动作。
- [ ] 实现项目导出、Markdown/JSON 报告和必要的导入兼容策略。
- [ ] 实现 metadata/full backup、恢复验证和损坏提示。
- [ ] 实现断线、Core 不可用和恢复中的明确 UI 状态。
- [ ] 按工程指南提供最小 Docker 运行方式；不引入 Kubernetes。

**验收门槛：**

- [ ] 新用户无需外部说明即可创建并推进第一个项目。
- [ ] 导出内容可读取，备份可在干净环境恢复并通过一致性检查。
- [ ] 离线或 Core 故障不会造成静默数据丢失。

### Phase 8 — Alpha Release

**状态：⬜ Not started**

**目标：** 向 20–50 名独立游戏开发者交付首个可反馈版本。

**重点：** 完整回归、三平台发布证据、文档、反馈渠道和已知限制。

**计划任务：**

- [ ] 冻结 Alpha 范围并关闭所有阻塞级缺陷。
- [ ] 完成 Windows、macOS、Linux 的可复现构建和安装验证。
- [ ] 完成数据升级、备份恢复、离线和核心业务场景回归。
- [ ] 发布使用说明、已知问题、隐私/凭据说明和反馈入口。
- [ ] 建立 Alpha 反馈分类、复现和修复节奏。

**验收门槛：**

- [ ] 20–50 名目标用户能够获得、启动并完成核心发布准备流程。
- [ ] 发布产物可追溯到 Git tag 和成功 CI，校验值可验证。
- [ ] 首个建议标签为 `v0.1.0-alpha.1`；原计划中的 `v0.1.0-alpha` 视为 Alpha 系列名称，实际打 tag 前确认。

## v0.1 优先级

### P0 — 必须形成闭环

- [x] Project
- [ ] Release Goal
- [ ] Checklist
- [ ] Resource
- [ ] Deliverable
- [ ] Export

### P1 — 重要增强

- [ ] Asset Map
- [ ] AI Review
- [ ] Reports

### P2 — 后续增强

- [ ] Knowledge Pack
- [ ] 更多模板
- [ ] 高级部署

## v0.1 明确不做

- Watch
- Steam 数据分析
- 多人协作
- 自动 Steam 提交
- 自动营销
- 通用 Agent 平台
- 自动更新

## 关键依赖与风险

- **阶段编号冲突：** 产品 MVP 路线图与开发实施计划不一致；实现与完成状态只使用本路线图中的 Phase 0–8 编号。
- **基础设施顺序：** SQLite、migration、sqlc 和 OpenAPI 边界必须先稳定，再进入 Project 及后续领域。
- **跨平台：** Fyne 原生依赖和三平台差异必须由 GitHub Actions 原生 runner 提供证据，不能以单平台或交叉编译替代。
- **本地数据安全：** migration、外部文件引用、导入导出和备份都必须有幂等与恢复测试。
- **Agent 顺序：** Agent 必须晚于 Checklist、Resource、Deliverable，且不得成为手工工作流的可用性前提。
- **凭据安全：** 模型 provider 凭据不得进入仓库、日志、导出或诊断包。
- **发布事项：** macOS 签名/公证、正式安装包和分发策略属于发布阶段决策，不是已完成的 Repository Bootstrap 能力。
- **文档许可：** 当前仓库沿用现有 AGPLv3 `LICENSE`；是否对文档单独采用 CC BY 4.0 仍待维护者确认。

## 发布策略

`v0.1` 按 `alpha → beta → stable` 逐步验证，不以一次性完成全部设想为目标。`v1.0.0` 仅保留给 API、数据兼容、跨平台发布、测试、文档和升级策略达到稳定承诺后的首个正式稳定版本；当前不承诺日期。

## 路线图维护规则

- 每次开始阶段时，将状态改为 **In progress**，并链接对应 ExecPlan 或实现追踪文档。
- 只有在验收命令、可观察行为和 CI/测试证据齐全后才能勾选任务或标记 **Completed**。
- 设计范围变化先更新相应源文档并记录决策，再调整本路线图；不得用路线图静默覆盖高优先级设计决策。
- 发现阶段边界冲突时先停止扩展范围，记录冲突、推荐选择和待确认事项。
- 每次状态变更同时更新“当前进度”、具体阶段和证据链接，避免多处状态不一致。
