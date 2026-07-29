# Corvus Studio Technology Stack Decision v0.1

> 文档状态：Decision Record\
> 产品：Corvus Studio\
> 版本：v0.1\
> 目标：确定 Corvus Studio Launch 的技术选型，为后续 Data
> Model、API、Agent Implementation 提供基础。

------------------------------------------------------------------------

# 1. 技术选型原则

Corvus Studio 是：

-   本地优先桌面应用；
-   Go Core + React Web UI；
-   AI 增强型独立游戏发布工具；
-   支持 Desktop / Headless / Docker。

技术选择优先级：

1.  长期维护性；
2.  跨平台能力；
3.  开源生态；
4.  开发效率；
5.  单用户本地部署体验。

不追求：

-   极限性能；
-   大规模 SaaS 架构；
-   过度复杂基础设施。

------------------------------------------------------------------------

# 2. 总体技术架构

    Fyne Launcher

            |

            |

    Go Core Runtime

            |

            |

    React Web UI


            |

    SQLite Storage


            |

    Agent Runtime

            |

    Model Providers

------------------------------------------------------------------------

# 3. Backend 技术栈

## 3.1 编程语言

选择：

    Go 1.26

原因：

-   单二进制部署；
-   跨平台；
-   文件系统能力强；
-   适合本地工具；
-   适合 Docker/systemd。

------------------------------------------------------------------------

# 3.2 HTTP Framework

选择：

    Echo v5

原因：

-   基于标准 net/http；
-   生态兼容性好；
-   Middleware 丰富；
-   适合长期维护；
-   支持本地 API Server 场景。

放弃 Fiber 的原因：

-   基于 fasthttp；
-   部分 Go 生态兼容成本更高。

------------------------------------------------------------------------

# 3.3 数据库

选择：

    SQLite

定位：

本地单用户数据存储。

保存：

-   项目状态；
-   关系；
-   元数据；
-   历史记录。

不保存：

-   大型图片；
-   视频；
-   音频；
-   PSD 等资产文件。

------------------------------------------------------------------------

# 3.4 SQLite Driver

选择：

    modernc.org/sqlite

原因：

-   纯 Go；
-   无 CGO；
-   跨平台编译简单；
-   适合桌面应用。

------------------------------------------------------------------------

# 3.5 Database Migration

选择：

    goose

用途：

-   Schema Migration；
-   数据版本管理。

------------------------------------------------------------------------

# 3.6 SQL Layer

选择：

    sqlc

原因：

-   类型安全；
-   明确 SQL；
-   避免 ORM 黑盒；
-   适合复杂领域模型。

------------------------------------------------------------------------

# 3.7 Logging

选择：

    zap

    +

    lumberjack

用途：

## zap

结构化日志。

## lumberjack

日志轮转。

支持：

-   文件大小限制；
-   历史日志保留；
-   本地工具运行。

------------------------------------------------------------------------

# 3.8 CLI

选择：

    cobra

用途：

未来支持：

-   serve；
-   backup；
-   export；
-   diagnose；
-   migration。

------------------------------------------------------------------------

# 3.9 Configuration

选择：

    viper

配置来源：

    CLI Arguments

    ↓

    Environment Variables

    ↓

    Config File

    ↓

    Default Values

支持：

-   YAML；
-   JSON；
-   ENV。

------------------------------------------------------------------------

# 3.10 Password Hash

选择：

    gofurry/easyhash

用途：

本地用户密码保护。

要求：

支持：

-   salt；
-   算法版本；
-   未来迁移。

------------------------------------------------------------------------

# 3.11 API Contract

选择：

    OpenAPI First

用途：

统一：

-   Backend API；
-   Frontend Client；
-   Agent Tool Contract。

------------------------------------------------------------------------

# 3.12 ID Strategy

选择：

    UUIDv7

原因：

-   时间有序；
-   适合本地项目；
-   方便未来同步和导出。

------------------------------------------------------------------------

# 3.13 Time Strategy

统一：

    UTC Storage

    ↓

    Local Time Display

用于：

-   发布日期；
-   Review 时间；
-   历史记录。

------------------------------------------------------------------------

# 3.14 File Hash

选择：

    SHA-256

用途：

-   文件变化检测；
-   资源版本；
-   去重。

------------------------------------------------------------------------

# 3.15 Streaming

选择：

    Server-Sent Events (SSE)

用途：

-   Agent 状态；
-   Review 进度；
-   长任务反馈。

暂不使用：

-   WebSocket。

原因：

v0.1 不需要复杂双向通信。

------------------------------------------------------------------------

# 4. Frontend 技术栈

## 4.1 Framework

选择：

    React 19

------------------------------------------------------------------------

# 4.2 Language

选择：

    TypeScript

原因：

-   类型安全；
-   复杂领域模型友好；
-   降低维护成本。

------------------------------------------------------------------------

# 4.3 Build Tool

选择：

    Vite 8

------------------------------------------------------------------------

# 4.4 Package Manager

选择：

    pnpm

原因：

-   React/Vite 生态支持好；
-   磁盘占用低；
-   Monorepo 友好。

------------------------------------------------------------------------

# 4.5 UI Framework

选择：

    Ant Design

原因：

适合：

-   表格；
-   表单；
-   Drawer；
-   Modal；
-   Tree；
-   Upload；
-   Timeline。

------------------------------------------------------------------------

# 4.6 CSS

选择：

    SCSS

原因：

-   主题管理；
-   全局变量；
-   复杂布局。

------------------------------------------------------------------------

# 4.7 State Management

选择：

    Zustand

用途：

-   UI 状态；
-   用户状态；
-   项目状态。

------------------------------------------------------------------------

# 4.8 Server State

选择：

    TanStack Query

用途：

-   API Cache；
-   Loading；
-   Retry；
-   数据刷新。

------------------------------------------------------------------------

# 4.9 Routing

选择：

    React Router

------------------------------------------------------------------------

# 4.10 Relationship Graph

选择：

    React Flow

用途：

Asset Map。

------------------------------------------------------------------------

# 4.11 Markdown

选择：

    react-markdown

    +

    remark

    +

    rehype

用途：

-   Report；
-   Codex；
-   文档展示。

------------------------------------------------------------------------

# 4.12 Frontend Testing

选择：

    Vitest

    +

    React Testing Library

------------------------------------------------------------------------

# 5. Desktop 技术栈

选择：

    Fyne

定位：

仅作为 Launcher。

负责：

-   启动 Go Core；
-   停止 Core；
-   系统托盘；
-   打开 Web UI。

不负责：

-   业务页面；
-   Asset Map；
-   Checklist；
-   Agent UI。

------------------------------------------------------------------------

# 6. Agent 技术栈

## Runtime

选择：

    Google ADK Go

------------------------------------------------------------------------

## 架构

    Agent Runtime

    ↓

    Workflow

    ↓

    Tool Layer

    ↓

    Domain Services

------------------------------------------------------------------------

## 模型能力

### Text Reasoning

默认：

DeepSeek

------------------------------------------------------------------------

### Vision / Audio / Video

默认：

Seed

------------------------------------------------------------------------

### Long Context / Video

支持：

Kimi K3

------------------------------------------------------------------------

### General Multimodal

支持：

Gemini

------------------------------------------------------------------------

# 7. Agent Model Abstraction

不做通用模型网关。

采用：

    Capability Layer

    ↓

    Provider Adapter

    ↓

    Model

能力：

-   Text Reasoning；
-   Image Review；
-   Video Understanding；
-   Audio Understanding；
-   Embedding。

------------------------------------------------------------------------

# 8. Deployment

支持：

## Desktop

    Fyne

    +

    Go Core

    +

    React

------------------------------------------------------------------------

## Binary

    Go Core

    +

    Browser

------------------------------------------------------------------------

## Docker

    Container

    +

    Mounted Project Data

------------------------------------------------------------------------

# 9. Quality Tools

## Go

选择：

    golangci-lint

------------------------------------------------------------------------

## Frontend

选择：

    ESLint

    +

    Prettier

------------------------------------------------------------------------

# 10. Final Technology Stack

## Backend

    Go 1.26

    Echo v5

    SQLite

    modernc.org/sqlite

    goose

    sqlc

    zap

    lumberjack

    cobra

    viper

    gofurry/easyhash

    OpenAPI

    UUIDv7

    SSE

------------------------------------------------------------------------

## Frontend

    React 19

    TypeScript

    Vite 8

    pnpm

    Ant Design

    SCSS

    Zustand

    TanStack Query

    React Router

    React Flow

    react-markdown

------------------------------------------------------------------------

## Desktop

    Fyne Launcher

------------------------------------------------------------------------

## Agent

    Google ADK Go

    DeepSeek

    Seed

    Kimi

    Gemini

------------------------------------------------------------------------

# 11. 后续设计依赖

基于本技术栈继续：

1.  Corvus Studio Data Model Design v0.1
2.  Corvus Studio Backend API Design v0.1
3.  Corvus Studio ADK Implementation Plan v0.1
4.  Corvus Studio Launch Frontend Specification v0.1
5.  Corvus Studio Engineering & Deployment Guide v0.1
6.  Corvus Studio Launch MVP Roadmap v0.1
