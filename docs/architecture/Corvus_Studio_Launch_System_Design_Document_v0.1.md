# Corvus Studio Launch System Design Document v0.1

> 文档状态：Draft\
> 产品：Corvus Studio Launch Workspace

## 1. 系统定位

Corvus Studio Launch 是一个本地优先的独立游戏 Steam 发布准备工作台。

核心目标：

> 帮助独立游戏开发者管理发布准备过程，并通过 AI 增强发现风险。

------------------------------------------------------------------------

## 2. 核心原则

### Local-first

核心数据保存在用户本机。

AI 是增强能力，不是运行前提。

离线可使用：

-   项目管理
-   Checklist
-   Resource
-   Deliverable
-   Asset Map
-   Export

------------------------------------------------------------------------

### 业务驱动架构

系统按业务领域划分：

    Project
    Release
    Checklist
    Resource
    Deliverable
    Template
    Knowledge
    Review
    Export
    Scheduler

------------------------------------------------------------------------

## 3. 总体架构

    User

    ↓

    React Web UI

    ↓

    Go Core Runtime

    ↓

    Domain Services

    ↓

    Storage / Review / Template / Event

    ↓

    Agent Runtime

    ↓

    Capability Layer

    ↓

    Model Providers

------------------------------------------------------------------------

## 4. 运行模式

### Desktop

    Fyne Launcher

    ↓

    Go Core

    ↓

    React UI

Fyne 只负责：

-   启动
-   停止
-   托盘
-   打开界面

------------------------------------------------------------------------

### Headless

    Go Core

    ↓

    Browser

------------------------------------------------------------------------

### Docker

    Container

    ├── Go Core
    ├── React Static
    └── Project Data

------------------------------------------------------------------------

## 5. Domain Service

### Project Domain

管理：

-   游戏项目
-   项目生命周期

### Release Domain

管理：

-   Release Goal
-   发布状态

### Checklist Domain

管理：

-   任务
-   状态
-   依赖

### Resource Domain

管理：

-   文件
-   URL
-   外部引用

### Deliverable Domain

管理：

-   Steam 交付对象
-   验证状态

### Template Domain

管理：

-   模板
-   模板版本

### Knowledge Domain

管理：

-   Corvus Knowledge
-   Project Codex
-   User Knowledge

### Review Domain

管理：

-   Review 流程
-   报告生成

------------------------------------------------------------------------

## 6. Storage 原则

SQLite 保存：

-   状态
-   关系
-   元数据
-   历史

文件系统保存：

-   图片
-   视频
-   音频
-   PSD
-   Deliverables

不将大型资产存入数据库。

------------------------------------------------------------------------

## 7. Project Folder

项目目录是一级公民。

支持：

### Corvus Project Folder

适合新项目。

### External Workspace

适合已有游戏项目。

Corvus 不强制迁移原始项目。

------------------------------------------------------------------------

## 8. Resource 系统

Resource 类型：

-   Managed Resource
-   Referenced Resource
-   Linked Resource

状态：

    Available
    Missing
    Changed
    Unavailable

------------------------------------------------------------------------

## 9. Deliverable 系统

Deliverable：

> 发布过程中需要完成并确认的交付对象。

类型：

-   文件
-   文本
-   外部操作
-   决策

状态：

    Missing
    Draft
    Reviewing
    Approved
    Published

------------------------------------------------------------------------

## 10. AI Integration Boundary

Agent 不直接访问数据库。

流程：

    Agent

    ↓

    Context Provider

    ↓

    Domain Service

    ↓

    Storage

------------------------------------------------------------------------

## 11. Review Pipeline

统一流程：

    Review Request

    ↓

    Collect Context

    ↓

    Deterministic Validation

    ↓

    AI Analysis

    ↓

    Normalize Result

    ↓

    Generate Report

------------------------------------------------------------------------

## 12. Validation Service

负责确定性检查：

-   文件存在
-   格式
-   尺寸
-   大小
-   文本限制

不交给模型。

------------------------------------------------------------------------

## 13. Capability Layer

不是通用模型网关。

能力：

-   Text Reasoning
-   Image Review
-   Video Understanding
-   Audio Understanding
-   Embedding

结构：

    Capability

    ↓

    Provider Adapter

    ↓

    Model

支持：

-   DeepSeek
-   Seed
-   Kimi
-   Gemini
-   Local Models

------------------------------------------------------------------------

## 14. Event System

预留内部事件：

-   ResourceChanged
-   DeliverableUpdated
-   ReviewCompleted
-   TemplateUpdated

v0.1 不实现复杂消息系统。

------------------------------------------------------------------------

## 15. 数据一致性

采用：

> 本地单写，多读。

所有修改经过 Go Core。

React 不直接修改数据库。

------------------------------------------------------------------------

## 16. 配置系统

分为：

### User Config

-   模型
-   API Key
-   UI 设置

### Project Config

-   Release Goal
-   Template Version

### Runtime Config

-   Port
-   Log

------------------------------------------------------------------------

## 17. 安全

API Key：

-   不进入项目
-   不进入导出
-   使用系统安全存储

AI 调用：

-   用户主动确认
-   明确发送内容

------------------------------------------------------------------------

## 18. 扩展点

预留：

-   Corvus Watch
-   Team Collaboration
-   Plugin System
-   Provider Extension

------------------------------------------------------------------------

## 19. Non-goals

v0.1 不做：

-   SaaS
-   多租户
-   自动 Steam 提交
-   通用 AI 平台
-   在线协作
-   Watch 数据分析

------------------------------------------------------------------------

# 总结

Corvus Studio Launch 的系统核心：

> 一个本地优先、领域驱动、AI 增强的 Steam 独立游戏发布准备系统。
