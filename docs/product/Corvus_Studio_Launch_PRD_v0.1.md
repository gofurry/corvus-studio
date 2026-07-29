# Corvus Studio Launch PRD v0.1

> 文档状态：Draft\
> 产品：Corvus Studio Launch Workspace\
> 首发目标：Steam 商店页与 Coming Soon 准备闭环

## 1. 产品定位

Corvus Studio Launch 是面向独立游戏开发者的本地优先发售运营工作台。

核心目标：

> 帮助开发者准备 Steam 商店页，发现发布前缺失项，并通过 AI
> 辅助完成发售准备。

不替代： - 游戏引擎 - 代码管理 - 通用项目管理工具

------------------------------------------------------------------------

# 2. v0.1 范围

## 实现

-   Steam 商店页准备
-   Coming Soon 准备
-   Release Goal
-   Checklist
-   Project Codex
-   Creative Library
-   Steam Deliverables
-   Asset Map
-   AI 发布准备诊断

## 暂不实现

-   Demo 发布
-   Playtest
-   Steam Next Fest
-   抢先体验
-   正式发售流程
-   Watch 数据分析

后续能力基于相同 Release Goal 机制扩展。

------------------------------------------------------------------------

# 3. 核心工作流

    Release Goal
        ↓
    Checklist
        ↓
    Resources
        ↓
    Steam Deliverables
        ↓
    Asset Map
        ↓
    AI Diagnosis

首页不是聊天窗口，而是发布目标状态页。

------------------------------------------------------------------------

# 4. 核心模型

## Game Project

游戏项目。

包含： - 游戏信息 - 发布目标 - 项目资源 - Project Codex

## Project Codex

项目知识档案。

包含： - 核心玩法 - 卖点 - 目标玩家 - 美术风格 - 对标作品 - 品牌语气 -
已确认事实 - 待决定事项

用于 Agent 理解项目。

## Release Goal

发布目标。

v0.1： - Steam Coming Soon

未来： - Demo - Playtest - Early Access - Full Release - Update

------------------------------------------------------------------------

# 5. Checklist

清单由模板生成，但属于用户项目数据。

来源：

-   Steam Requirement
-   Corvus Template
-   User Task
-   Agent Suggestion

状态：

-   Not Started
-   In Progress
-   Needs Review
-   Done
-   Not Applicable
-   Blocked

Steam 硬性任务允许删除，但：

-   保留风险提示
-   记录来源
-   可查看 Ignored Requirements

------------------------------------------------------------------------

# 6. Resource 系统

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

自定义素材库。

保存： - 原始图片 - PSD - Logo - Trailer 素材 - 宣传资源

## Steam Deliverables

Steam 最终交付物。

保存： - Capsule - Screenshots - Trailer - Store Copy - Localization

交付物必须： - 可检查 - 可上传 - 可版本管理 - 可关联源素材

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

------------------------------------------------------------------------

# 9. Agent 设计

Agent 定位：

> 发布准备审查员

不是聊天机器人。

## DeepSeek

负责： - 主 Agent - 工具调用 - 文案生成 - 项目分析 - 最终诊断

## Seed 多模态

负责： - 图片检查 - 视频检查 - 音频检查 - 视觉风险分析

## Kimi K3

负责： - 长视频理解 - 长上下文分析

## Gemini

负责： - 通用多模态 - 备用分析

------------------------------------------------------------------------

# 10. 模型策略

不建设通用模型聚合平台。

采用：

能力层 + 少量官方适配模型。

能力包括：

-   Text Reasoning
-   Image Review
-   Video Review
-   Audio Review
-   Long Context

所有云模型调用：

-   用户主动触发
-   显示发送内容
-   用户确认后执行

------------------------------------------------------------------------

# 11. AI 检查报告

报告结构：

## Blockers

阻止发布的问题。

## Warnings

潜在风险。

## Suggestions

优化建议。

## Next Actions

下一步行动。

结果区分：

-   程序确定性检查
-   多模态观察
-   Agent 推理建议

------------------------------------------------------------------------

# 12. 首页设计

首页展示：

-   当前 Release Goal
-   当前状态
-   Blockers
-   Next Actions
-   Checklist 状态
-   Asset 状态
-   最近诊断

主要操作：

> Continue Preparation

------------------------------------------------------------------------

# 13. 本地优先与隐私

默认：

-   数据保存在本机
-   不自动上传素材
-   不后台调用模型
-   不自动扫描项目

模型调用前显示：

-   使用模型
-   发送资源
-   数据类型

------------------------------------------------------------------------

# 14. Steam 规则模板

规则作为数据管理。

包含：

-   模板版本
-   来源
-   更新时间
-   官方参考
-   变更记录

升级：

-   不覆盖用户修改
-   提供差异
-   用户选择合并

------------------------------------------------------------------------

# 15. v0.1 验收标准

用户能够：

1.  创建游戏项目
2.  选择 Steam Coming Soon 目标
3.  生成准备清单
4.  管理 Project Codex
5.  管理素材资源
6.  准备 Steam Deliverables
7.  查看 Asset Map
8.  执行 AI 诊断
9.  获得可追溯问题和建议

------------------------------------------------------------------------

# 16. 后续路线

Launch：

-   Demo
-   Playtest
-   宣传活动
-   正式发售
-   更新运营

Watch：

-   Steam 公开数据
-   竞品监控
-   市场分析

------------------------------------------------------------------------

# 总结

Corvus Studio Launch 的核心价值：

> 将独立游戏 Steam 发布过程转化为可管理的目标、清单、资源、交付物和 AI
> 辅助诊断系统。
