# Corvus Studio Repository Structure Design v0.1

> 文档状态：Draft\
> 项目：Corvus Studio

------------------------------------------------------------------------

# 1. Repository 原则

Corvus Studio 采用 Monorepo。

目标：

-   统一管理 Go Core；
-   React Frontend；
-   Fyne Launcher；
-   Agent；
-   文档；
-   部署配置。

原则：

-   业务入口与共享能力分离；
-   文档是一等公民；
-   支持长期开源维护。

------------------------------------------------------------------------

# 2. Repository 总结构

    corvus-studio/

    ├── apps/
    │   ├── core/
    │   ├── web/
    │   └── launcher/
    │
    ├── agent/
    │
    ├── packages/
    │   ├── api-client/
    │   ├── shared-types/
    │   └── ui/
    │
    ├── docs/
    │
    ├── deployments/
    │
    ├── configs/
    │
    ├── scripts/
    │
    ├── tests/
    │   ├── integration/
    │   ├── scenario/
    │   └── fixtures/
    │
    ├── tools/
    │
    ├── .github/
    │
    ├── go.work
    ├── pnpm-workspace.yaml
    ├── package.json
    ├── README.md
    └── LICENSE

------------------------------------------------------------------------

# 3. apps/core

Go Core 主程序。

负责：

-   HTTP Server；
-   Domain Services；
-   Storage；
-   Agent Service；
-   Scheduler。

结构：

    apps/core/

    ├── cmd/
    ├── internal/
    ├── migrations/
    ├── configs/
    └── go.mod

------------------------------------------------------------------------

## cmd

程序入口。

支持：

    corvus serve

    corvus migrate

    corvus doctor

    corvus backup

------------------------------------------------------------------------

## internal

业务实现。

包含：

-   domain；
-   application；
-   infrastructure；
-   interfaces。

------------------------------------------------------------------------

# 4. Domain 组织

    domain/

    ├── project
    ├── release
    ├── checklist
    ├── resource
    ├── deliverable
    ├── codex
    ├── template
    ├── knowledge
    └── review

------------------------------------------------------------------------

# 5. Agent 目录

独立顶层：

    agent/

    ├── runtime
    ├── workflows
    ├── tools
    ├── context
    ├── prompts
    └── providers

原因：

Agent 是核心能力，不只是内部实现。

------------------------------------------------------------------------

# 6. apps/web

React 前端。

    apps/web/

    ├── src/
    ├── pages/
    ├── features/
    ├── components/
    ├── api/
    ├── hooks/
    ├── stores/
    └── styles/

------------------------------------------------------------------------

Feature 按业务划分：

    features/

    ├── project
    ├── release
    ├── checklist
    ├── resource
    ├── asset-map
    ├── codex
    ├── review
    └── settings

------------------------------------------------------------------------

# 7. apps/launcher

Fyne Launcher。

只负责：

-   启动 Core；
-   停止 Core；
-   系统托盘；
-   打开 Web UI。

不包含业务逻辑。

------------------------------------------------------------------------

# 8. packages

共享能力。

## api-client

OpenAPI 生成客户端。

## shared-types

共享类型。

## ui

未来共享组件。

------------------------------------------------------------------------

# 9. docs

设计文档：

    docs/

    ├── product
    ├── architecture
    ├── api
    ├── agent
    ├── development
    └── deployment

------------------------------------------------------------------------

# 10. deployments

部署：

    deployments/

    ├── docker
    ├── systemd
    └── packaging

------------------------------------------------------------------------

# 11. tests

测试分层：

## 模块测试

靠近代码。

## 集成测试

    tests/integration

## 场景测试

    tests/scenario

用于：

-   Agent 回归；
-   用户流程测试。

------------------------------------------------------------------------

# 12. GitHub 工程

    .github/

    ├── workflows
    └── SECURITY.md

Issue/PR 不使用强制模板。

贡献者自行组织描述。

------------------------------------------------------------------------

# 13. 开发模式

采用：

Trunk Based Development。

主分支：

    main

短期：

    feature/*

不采用复杂 Git Flow。
