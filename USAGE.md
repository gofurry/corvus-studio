# Corvus Studio Usage & Development Guide

> 文档状态：Draft\
> 项目：Corvus Studio\
> 用途：用户使用、开发环境搭建、贡献者快速入门

------------------------------------------------------------------------

# 1. Introduction

Corvus Studio 是一个本地优先的独立游戏 Steam 发布准备工作台。

支持：

-   Desktop Application
-   Local Web Application
-   Docker Deployment

核心能力：

-   Project 管理
-   Steam Coming Soon 准备
-   Checklist
-   Resource Library
-   Deliverables
-   AI Review

AI 是增强能力，不是运行前提。

------------------------------------------------------------------------

# 2. User Quick Start

## Desktop

启动流程：

    Launch Corvus Studio

    ↓

    Create Project

    ↓

    Select Steam Coming Soon

    ↓

    Start Preparing

------------------------------------------------------------------------

## AI 配置

进入：

    Settings

    ↓

    AI Providers

    ↓

    Add Provider Key

    ↓

    Test Connection

AI 未配置时：

仍可使用：

-   Project
-   Checklist
-   Resource
-   Deliverable

------------------------------------------------------------------------

# 3. Developer Environment

## Backend

需要：

    Go 1.26+

检查：

``` bash
go version
```

------------------------------------------------------------------------

## Frontend

需要：

    Node.js

    pnpm

检查：

``` bash
node -v
pnpm -v
```

------------------------------------------------------------------------

## Optional

Docker：

``` bash
docker --version
```

------------------------------------------------------------------------

# 4. Repository Setup

获取代码：

``` bash
git clone <repository>

cd corvus-studio
```

------------------------------------------------------------------------

## Backend

``` bash
cd apps/core

go mod download
```

------------------------------------------------------------------------

## Frontend

``` bash
pnpm install
```

------------------------------------------------------------------------

# 5. Local Development

开发模式：

## Terminal 1

启动 Core：

``` bash
cd apps/core

go run ./cmd/corvus serve
```

------------------------------------------------------------------------

## Terminal 2

启动 Web：

``` bash
cd apps/web

pnpm dev
```

------------------------------------------------------------------------

## Terminal 3（可选）

启动 Launcher：

``` bash
cd apps/launcher

go run .
```

------------------------------------------------------------------------

# 6. Database Development

技术：

-   SQLite
-   goose
-   sqlc

------------------------------------------------------------------------

Migration：

位置：

    apps/core/migrations

------------------------------------------------------------------------

执行迁移：

``` bash
corvus migrate
```

------------------------------------------------------------------------

创建 Migration：

``` bash
goose create add_feature sql
```

------------------------------------------------------------------------

生成 SQL Code：

``` bash
sqlc generate
```

------------------------------------------------------------------------

# 7. Frontend Development

技术：

-   React 19
-   TypeScript
-   Vite
-   Ant Design

------------------------------------------------------------------------

启动：

``` bash
pnpm dev
```

------------------------------------------------------------------------

构建：

``` bash
pnpm build
```

------------------------------------------------------------------------

API Client：

通过 OpenAPI 生成：

    Backend OpenAPI

    ↓

    Generated Client

    ↓

    Frontend

------------------------------------------------------------------------

# 8. Agent Development

Agent 目录：

    agent/

结构：

    agent/

    ├── runtime
    ├── workflows
    ├── tools
    ├── context
    ├── prompts
    └── providers

------------------------------------------------------------------------

设计原则：

Agent 不直接访问数据库。

流程：

    Agent

    ↓

    Workflow

    ↓

    Tool

    ↓

    Domain Service

------------------------------------------------------------------------

Prompt：

    agent/prompts

------------------------------------------------------------------------

Workflow：

    agent/workflows

------------------------------------------------------------------------

Tool：

    agent/tools

------------------------------------------------------------------------

# 9. Testing

## Go

``` bash
go test ./...
```

------------------------------------------------------------------------

## Frontend

``` bash
pnpm test
```

------------------------------------------------------------------------

## Lint

Go：

``` bash
golangci-lint run
```

Frontend：

``` bash
pnpm lint
```

------------------------------------------------------------------------

## Scenario Test

位置：

    tests/scenario

用于：

-   用户流程测试
-   Agent 回归测试

------------------------------------------------------------------------

# 10. Debugging

## Backend

开启 Debug 日志：

``` yaml
log:
  level: debug
```

------------------------------------------------------------------------

## Agent

Debug 模式记录：

-   Context
-   Tool Call
-   Model Call

------------------------------------------------------------------------

## Frontend

使用浏览器 DevTools。

------------------------------------------------------------------------

# 11. Configuration

Corvus 使用 Viper。

配置优先级：

    CLI

    >

    Environment Variables

    >

    Config File

    >

    Default Values

------------------------------------------------------------------------

配置包含：

-   Server
-   Storage
-   AI Provider
-   Logging

------------------------------------------------------------------------

# 12. Build

## Backend

``` bash
go build ./cmd/corvus
```

------------------------------------------------------------------------

## Frontend

``` bash
pnpm build
```

------------------------------------------------------------------------

## Desktop

构建：

-   Go Core
-   React Assets
-   Fyne Launcher

------------------------------------------------------------------------

# 13. Docker Development

启动：

``` bash
docker compose up
```

------------------------------------------------------------------------

用户数据通过挂载保存。

示例：

    project-data/

------------------------------------------------------------------------

# 14. Contribution

Corvus Studio 不提供强制 Issue/PR Template。

贡献者可以自由组织描述。

建议：

Issue：

包含：

-   背景
-   问题
-   期望

PR：

包含：

-   修改内容
-   测试方式
-   影响范围

------------------------------------------------------------------------

# 15. Troubleshooting

## Core 无法启动

检查：

-   配置文件；
-   SQLite；
-   文件权限。

------------------------------------------------------------------------

## AI 不工作

检查：

-   API Key；
-   Provider 配置；
-   网络连接。

------------------------------------------------------------------------

## Resource Missing

检查：

-   文件路径；
-   外部引用；
-   权限。

------------------------------------------------------------------------

# 16. 文档导航

推荐阅读顺序：

    README.md

    ↓

    USAGE.md

    ↓

    docs/

    ├── Product Design
    ├── PRD
    ├── UX
    ├── System Design
    ├── API
    └── Agent Architecture

------------------------------------------------------------------------

# 总结

USAGE.md 是 Corvus Studio 的用户与开发入口。

它回答：

-   如何运行 Corvus Studio；
-   如何配置环境；
-   如何参与开发；
-   如何调试和测试。

详细架构设计请参考 docs 目录。
