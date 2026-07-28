# Corvus Studio Launch Functional Specification v0.1

> 文档状态：Draft\
> 产品：Corvus Studio Launch Workspace\
> 版本：v0.1

## 1. 功能原则

Corvus Studio Launch 围绕 Steam 发布准备状态设计。

核心链路：

    Release Goal
        ↓
    Checklist
        ↓
    Resource
        ↓
    Deliverable
        ↓
    Review
        ↓
    Diagnosis

原则：

-   用户拥有最终控制权；
-   AI 提供建议，不默认修改；
-   所有状态必须可解释。

------------------------------------------------------------------------

# 2. 核心对象

## Game Project

管理独立游戏项目。

包含：

-   项目信息
-   Release Goal
-   Project Codex
-   Resources
-   Checklist
-   Reports

------------------------------------------------------------------------

## Release Goal

表示当前发布目标。

v0.1：

-   Steam Coming Soon Page

未来：

-   Demo
-   Playtest
-   Full Release
-   Update

------------------------------------------------------------------------

## Checklist Item

管理发布准备任务。

来源：

-   Platform Template
-   Corvus Template
-   User Task
-   Agent Suggestion

状态：

    Not Started
    In Progress
    Needs Review
    Done
    Not Applicable
    Blocked

Steam 硬性任务允许删除，但保留风险记录。

------------------------------------------------------------------------

## Resource

发布相关资源。

类型：

-   Image
-   Video
-   Audio
-   Document
-   Text
-   URL
-   Folder Reference

Resource 是原始素材和参考资料。

------------------------------------------------------------------------

## Deliverable

发布交付对象。

类型：

-   File Deliverable
-   Text Deliverable
-   External Action Deliverable
-   Decision Deliverable

Resource 与 Deliverable 分离。

支持多对多关系。

------------------------------------------------------------------------

# 3. Asset Map

Asset Map 是领域关系图，不是通用白板。

节点：

-   Resource
-   Deliverable
-   Checklist Item
-   Codex Fact
-   Release Goal
-   Store Page Section

关系：

-   来源于
-   用于
-   支持
-   阻塞
-   替代
-   版本关系

Agent 可以建议关系，但需要用户确认。

------------------------------------------------------------------------

# 4. Project Codex

Project Codex 是 Agent 理解项目的事实层。

内容：

-   游戏定位
-   核心玩法
-   目标玩家
-   卖点
-   美术风格
-   参考作品
-   限制条件
-   决策记录

状态：

    Confirmed
    Temporary
    Need Decision
    AI Suggestion

------------------------------------------------------------------------

# 5. Agent Review

Agent 定位：

> 独游发行顾问 + 严谨审查员

不是聊天机器人。

Review 类型：

## Release Readiness Review

检查整个发布目标。

## Deliverable Review

检查单个交付物。

## Checklist Review

检查任务完成情况。

## Codex Review

检查项目定位一致性。

------------------------------------------------------------------------

Agent 输出：

    Blockers

    Warnings

    Suggestions

    Next Actions

每条结论包含：

-   Finding
-   Evidence
-   Reason
-   Suggestion

------------------------------------------------------------------------

# 6. 模型能力

## DeepSeek

负责：

-   主 Agent
-   推理
-   工具调用
-   报告生成

## Seed

负责：

-   图片
-   视频
-   音频分析

## Kimi K3

负责：

-   长视频
-   长上下文分析

## Gemini

负责：

-   通用多模态
-   备用分析

------------------------------------------------------------------------

# 7. Template System

Template 是发布经验结构化载体。

类型：

## Platform Template

Steam 等平台规则。

## Corvus Knowledge Template

独游发行经验。

## User Template

用户个人经验。

流程：

    Release Goal
     ↓
    Select Templates
     ↓
    Merge
     ↓
    Generate Checklist

模板生成项目副本，不直接绑定。

支持：

-   导入
-   导出

不做在线市场。

------------------------------------------------------------------------

# 8. Knowledge System

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

# 9. Local-first

支持：

## Managed Resource

复制到项目管理。

## Referenced Resource

保存外部路径。

## Linked Resource

保存 URL。

资源状态：

    Available
    Missing
    Changed
    Unavailable

------------------------------------------------------------------------

# 10. Import / Export / Backup

支持：

-   Project Export
-   Markdown Report Export
-   JSON Export
-   Template Import/Export

备份：

-   Metadata Backup
-   Full Backup
-   Archive Export

不保存：

-   API Key
-   临时 AI 数据

------------------------------------------------------------------------

# 11. Settings

包含：

-   General
-   Projects
-   AI Models
-   Privacy
-   Templates
-   Backup

AI 设置：

-   Provider
-   API Key
-   Capability Test

------------------------------------------------------------------------

# 12. AI 隐私

默认：

-   不自动上传资源；
-   不后台调用模型；
-   不自动扫描项目。

调用前显示：

-   模型；
-   发送内容；
-   数据类型。

用户确认后执行。

------------------------------------------------------------------------

# 13. v0.1 范围

必须：

-   Game Project
-   Steam Coming Soon Goal
-   Checklist
-   Resource Library
-   Steam Deliverables
-   Asset Map
-   Project Codex
-   Agent Review
-   Template
-   Export

暂不：

-   多人协作
-   SaaS
-   自动 Steam 提交
-   通用 Agent 平台
-   Watch 数据分析

------------------------------------------------------------------------

# 14. 验收标准

一个没有 Steam 发布经验的独立开发者，可以：

1.  创建项目；
2.  选择 Coming Soon；
3.  获得准备清单；
4.  管理资源；
5.  生成交付物；
6.  使用 AI Review；
7.  获得可追溯诊断；
8.  导出发布档案。

------------------------------------------------------------------------

# 总结

Corvus Studio Launch 的核心：

> 将 Steam
> 独立游戏发布过程转化为可管理的目标、清单、资源、交付物、知识和 AI
> 辅助诊断系统。
