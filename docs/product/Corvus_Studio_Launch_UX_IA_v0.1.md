# Corvus Studio Launch UX Flow & Information Architecture v0.1

> 文档状态：Draft\
> 产品：Corvus Studio Launch Workspace\
> 版本：v0.1

## 1. UX 核心理念

Corvus Studio Launch 是面向独立游戏开发者的 Steam 发布准备工作台。

它不是： - 通用项目管理工具 - 文件管理器 - AI 聊天应用

核心问题：

> 我的游戏距离公开发布还缺少什么？

核心模型：

    Release Goal
        ↓
    Status
        ↓
    Action
        ↓
    Evidence
        ↓
    Diagnosis

------------------------------------------------------------------------

# 2. 首次启动流程

采用强流程向导。

流程：

    Welcome
     ↓
    Create Game Project
     ↓
    Select Release Goal
     ↓
    Generate Workspace
     ↓
    Complete Preparation

首次创建最少信息：

-   Game Name
-   Project Location
-   Steam AppID（可选）
-   Primary Language
-   Current Stage

Release Goal：

v0.1: - Steam Coming Soon Page

未来： - Demo Release - Playtest - Full Release - Update

------------------------------------------------------------------------

# 3. 信息架构

    Corvus Studio

    ├── Home
    ├── Release
    ├── Checklist
    ├── Resources
    ├── Asset Map
    ├── Codex
    ├── Reports
    └── Watch (Future)

------------------------------------------------------------------------

# 4. Home 发布准备驾驶舱

首页不是聊天入口。

目标：

回答：

1.  当前发布目标是什么？
2.  当前状态如何？
3.  最大风险是什么？
4.  下一步做什么？

状态模型：

    Draft
     ↓
    Preparing
     ↓
    Needs Attention
     ↓
    Ready for Review
     ↓
    Ready
     ↓
    Submitted
     ↓
    Published

Readiness 不是简单完成率，而是综合：

-   Blockers
-   Required 项完成情况
-   Review 结果
-   资产完整性
-   时间风险

------------------------------------------------------------------------

# 5. Checklist 工作区

Checklist 是第一核心工作区。

任务状态：

    Not Started
    In Progress
    Needs Review
    Done
    Not Applicable
    Blocked

任务详情：

-   What
-   Why
-   Requirement
-   Evidence
-   Related Resources
-   Agent Review
-   Actions

Steam 硬性任务允许删除：

-   保留风险提示
-   记录来源
-   进入 Ignored Requirements

------------------------------------------------------------------------

# 6. Resource Library

Resource 不等于文件。

支持：

-   本地文件
-   本地目录
-   URL
-   云资源链接
-   图片
-   视频
-   音频
-   文档
-   文本

------------------------------------------------------------------------

# 7. 资产体系

## Creative Library

保存：

-   原始图片
-   PSD
-   Logo
-   Trailer 素材
-   宣传资源

## Steam Deliverables

保存：

-   Capsule
-   Screenshots
-   Trailer
-   Store Copy
-   Localization

交付物：

-   可检查
-   可上传
-   可版本管理
-   关联源素材

------------------------------------------------------------------------

# 8. Asset Map

Asset Map 是领域关系图，不是通用白板。

节点：

-   Resource
-   Deliverable
-   Checklist Item
-   Project Codex Fact
-   Release Goal
-   Store Page Section

关系：

-   来源于
-   用于
-   支持
-   阻塞
-   替代
-   版本关系

Agent 可以建议关系，但必须用户确认。

------------------------------------------------------------------------

# 9. Project Codex

Project Codex 是 AI 的项目知识层。

内容：

-   Game Identity
-   Core Gameplay
-   Audience
-   Selling Points
-   Visual Style
-   References
-   Restrictions
-   Decisions

状态：

    Confirmed
    Temporary
    Need Decision
    AI Suggestion

------------------------------------------------------------------------

# 10. Template System

模板是发布经验的结构化载体。

三层：

## Platform Template

来源： - Steam 官方规则

## Corvus Knowledge Template

来源： - 独游发行经验 - 社区经验

## User Template

用户自己的经验。

应用：

    Release Goal
     ↓
    Select Template Pack
     ↓
    Generate Checklist
     ↓
    User Owns Checklist

模板支持：

-   导入
-   导出

不做在线模板市场。

------------------------------------------------------------------------

# 11. Knowledge System

三层：

    Corvus Knowledge
    +
    Project Codex
    +
    User Knowledge

优先级：

    User Confirmed Fact
    >
    Platform Requirement
    >
    Project Codex
    >
    User Knowledge
    >
    Corvus Recommendation
    >
    General Model Knowledge

------------------------------------------------------------------------

# 12. Agent UX

Agent 定位：

> 独游发行顾问 + 严谨审查员

不是：

-   Chatbot
-   自动执行机器人

输入：

-   Project Codex
-   Release Goal
-   Checklist State
-   Resource Graph
-   Knowledge Base
-   Model Tools

输出：

    Blockers
    Warnings
    Suggestions
    Next Actions

------------------------------------------------------------------------

# 13. AI Review 流程

    Run Review

    ↓

    选择范围

    ↓

    显示发送内容

    ↓

    用户确认

    ↓

    模型分析

    ↓

    生成报告

所有云模型调用默认用户主动触发。

------------------------------------------------------------------------

# 14. Reports

保存：

-   Release Diagnosis
-   Visual Review
-   Text Review
-   历史检查结果

每条结果包含：

-   Reason
-   Evidence
-   Suggested Action
-   Related Item

------------------------------------------------------------------------

# 15. 最终 UX 模型

    Release Goal

        ↓

    Checklist

        ↓

    Resources + Codex

        ↓

    Steam Deliverables

        ↓

    Asset Map

        ↓

    Agent Review

        ↓

    Reports

------------------------------------------------------------------------

# 总结

Corvus Studio Launch 的核心不是 AI 聊天。

而是：

> 将独立游戏 Steam 发布过程转化为可管理的目标、清单、资源、关系、知识和
> AI 辅助诊断系统。
