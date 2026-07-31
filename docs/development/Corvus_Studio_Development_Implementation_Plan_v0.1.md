# Corvus Studio Development Implementation Plan v0.1

> 文档状态：Draft\
> 目标：定义 Corvus Studio v0.1 开发实施路线。

------------------------------------------------------------------------

# 1. 开发原则

每个阶段必须产出：

-   可运行成果；
-   可验证功能；
-   用户价值。

避免只完成代码而没有闭环。

------------------------------------------------------------------------

# 2. Phase 0：Repository Bootstrap

目标：

建立开发基础。

完成：

-   Monorepo；
-   Go Module；
-   pnpm Workspace；
-   GitHub Actions；
-   README；
-   License。

验收：

    go test ./...

    pnpm build

------------------------------------------------------------------------

# 3. Phase 1：Core Runtime

目标：

启动 Corvus Core。

完成：

-   Echo；
-   Viper；
-   Zap；
-   SQLite；
-   goose；
-   sqlc；
-   Cobra。

验收：

    corvus serve

可以运行。

------------------------------------------------------------------------

# 4. Phase 2：Project Foundation

目标：

第一个业务闭环。

完成：

-   Project Domain；
-   Project API；
-   Project 页面；
-   本地存储。

用户可以：

-   创建项目；
-   打开项目。

------------------------------------------------------------------------

# 5. Phase 3：Release + Checklist

目标：

形成 Steam 发布准备基础。

完成：

-   Release Goal；
-   Steam Coming Soon Template；
-   Checklist。

用户可以：

-   创建发布目标；
-   查看任务；
-   修改状态。

------------------------------------------------------------------------

# 6. Phase 4：Resource System

拆分：

## Phase 4A Resource

完成：

-   图片；
-   URL；
-   文本；
-   外部引用。

------------------------------------------------------------------------

## Phase 4B Deliverable

完成：

-   Main Capsule；
-   Screenshots；
-   Store Copy。

建立：

Resource → Deliverable → Checklist

闭环。

------------------------------------------------------------------------

# 7. Phase 5：Asset Map

完成：

-   React Flow；
-   节点展示；
-   关系编辑；
-   布局保存。

------------------------------------------------------------------------

# 8. Phase 6：Agent MVP

依赖：

-   Checklist；
-   Resource；
-   Deliverable。

完成：

-   ADK Runtime；
-   Root Agent；
-   Tool Registry；
-   Context Provider。

首批 Workflow：

## Release Diagnosis

## Deliverable Review

模型：

-   DeepSeek；
-   Seed。

------------------------------------------------------------------------

# 9. Phase 7：体验完善

完成：

-   Onboarding；
-   Command Palette；
-   Export；
-   Backup；
-   Offline UI；
-   Docker。

------------------------------------------------------------------------

# 10. Phase 8：Alpha Release

版本：

    v0.1.0-alpha

目标：

20-50 名独立游戏开发者测试。

------------------------------------------------------------------------

# 11. 优先级

## P0

必须：

-   Project；
-   Release Goal；
-   Checklist；
-   Resource；
-   Deliverable；
-   Export。

------------------------------------------------------------------------

## P1

重要：

-   Asset Map；
-   AI Review；
-   Reports。

------------------------------------------------------------------------

## P2

增强：

-   Knowledge Pack；
-   更多模板；
-   高级部署。

------------------------------------------------------------------------

# 12. 明确不做

v0.1：

不做：

-   Watch；
-   Steam 数据分析；
-   多人协作；
-   自动 Steam 提交；
-   自动营销；
-   通用 Agent 平台。

------------------------------------------------------------------------

# 13. Git Workflow

采用：

Trunk Based Development。

分支：

    main

    feature/*

------------------------------------------------------------------------

# 14. Commit

不强制 Conventional Commits。

推荐：

清晰描述：

例如：

    Add project dashboard

    Fix resource validation

------------------------------------------------------------------------

# 15. 发布策略

不追求一次完成。

通过：

alpha → beta → stable

逐步验证。

------------------------------------------------------------------------

# 总结

开发顺序：

    Core

    ↓

    Project

    ↓

    Release

    ↓

    Checklist

    ↓

    Resource

    ↓

    Deliverable

    ↓

    Asset Map

    ↓

    Agent

    ↓

    Alpha Release

保证先建立可靠产品基础，再增加 AI 能力。
