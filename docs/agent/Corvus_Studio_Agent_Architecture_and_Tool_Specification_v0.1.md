# Corvus Studio Agent Architecture & Tool Specification v0.1

> 文档状态：Draft\
> 产品：Corvus Studio Launch Workspace\
> 目标：定义 Agent Runtime、Tool、Workflow、Memory、Policy 和治理模型

------------------------------------------------------------------------

# 1. Agent 定位

Corvus Agent 不是：

-   Chatbot；
-   通用 AI 助手；
-   自动执行机器人。

定位：

> 独立游戏发行顾问 + 严谨审查员。

核心职责：

-   理解发布状态；
-   发现风险；
-   分析资产；
-   提供建议；
-   生成结构化报告。

不负责：

-   自动修改项目；
-   自动上传 Steam；
-   自动替用户决策。

------------------------------------------------------------------------

# 2. Agent 总体架构

    User

    ↓

    Corvus Agent

    ↓

    Workflow Selection

    ↓

    Context Provider

    ↓

    Tool Layer

    ↓

    Domain Services

    ↓

    Review Pipeline

    ↓

    Capability Layer

    ↓

    Models

------------------------------------------------------------------------

# 3. Agent Runtime

## 设计原则

Corvus 使用：

> Agent + Workflow

而不是自由规划型 Agent。

------------------------------------------------------------------------

## Agent 负责

-   理解用户意图；
-   选择 Workflow；
-   调用工具；
-   解释结果。

------------------------------------------------------------------------

## Workflow 负责

-   固定检查流程；
-   数据收集；
-   状态验证；
-   报告生成。

------------------------------------------------------------------------

# 4. Agent 类型

v0.1：

单主 Agent。

    Corvus Advisor

    ↓

    Tools

    ↓

    Specialized Models

不实现 Multi-Agent。

------------------------------------------------------------------------

# 5. Session 与 Memory

## Session Memory

保存：

-   当前对话；
-   当前 Workflow；
-   工具调用过程。

生命周期：

短期。

------------------------------------------------------------------------

## Project Memory

对应：

Project Codex。

保存：

-   项目事实；
-   项目约束；
-   用户确认决策。

------------------------------------------------------------------------

## Review Memory

保存：

-   历史检查；
-   已解决问题；
-   重复风险。

------------------------------------------------------------------------

## User Knowledge

保存：

-   用户自己的经验；
-   工作流程偏好。

------------------------------------------------------------------------

# 6. Memory 优先级

    User Confirmed Knowledge

    >

    Platform Requirement

    >

    Project Codex

    >

    Review History

    >

    User Knowledge

    >

    Corvus Recommendation

    >

    General Model Knowledge

------------------------------------------------------------------------

# 7. Prompt Architecture

Prompt 分层：

    Identity Prompt

    +

    Policy Prompt

    +

    Workflow Prompt

    +

    Tool Prompt

    +

    Domain Prompt

    +

    User Context

------------------------------------------------------------------------

## Identity Prompt

定义：

Agent 身份和目标。

------------------------------------------------------------------------

## Policy Prompt

定义：

-   不编造规则；
-   不未经确认修改；
-   必须提供证据。

------------------------------------------------------------------------

## Workflow Prompt

定义：

具体检查目标。

例如：

Release Diagnosis。

------------------------------------------------------------------------

# 8. Tool 设计原则

Tool 是业务能力，不是 API 包装。

错误：

    query_database()
    read_file()

正确：

    get_release_status()

    review_asset()

    suggest_task()

------------------------------------------------------------------------

# 9. Tool 权限分类

## READ

无副作用。

例如：

-   get_project_info
-   get_release_status

------------------------------------------------------------------------

## ANALYZE

执行分析。

例如：

-   validate_deliverable
-   request_visual_review

------------------------------------------------------------------------

## SUGGEST

产生建议。

例如：

-   suggest_task
-   suggest_resource_relation

------------------------------------------------------------------------

## MUTATE

修改数据。

必须用户确认。

例如：

-   create_task
-   update_codex

------------------------------------------------------------------------

# 10. Tool Schema 标准

每个 Tool 包含：

    Identity

    Description

    Input Schema

    Output Schema

    Permission

    Side Effect

    Confirmation Policy

    Error Handling

------------------------------------------------------------------------

# 11. 核心 Tool

## Project Tools

包括：

-   get_project_info
-   get_codex_entries
-   search_project_context

------------------------------------------------------------------------

## Release Tools

包括：

-   get_release_goal
-   get_release_readiness

------------------------------------------------------------------------

## Checklist Tools

包括：

-   list_tasks
-   get_task_detail
-   suggest_task

------------------------------------------------------------------------

## Resource Tools

包括：

-   search_resources
-   get_resource_metadata
-   check_resource_availability
-   suggest_resource_relation

------------------------------------------------------------------------

## Deliverable Tools

包括：

-   list_deliverables
-   get_deliverable_detail
-   validate_deliverable
-   request_visual_review

------------------------------------------------------------------------

## Review Tools

包括：

-   create_review_request
-   get_review_result
-   generate_report

------------------------------------------------------------------------

# 12. Human-in-the-loop

所有修改：

必须：

    Agent Suggestion

    ↓

    User Confirm

    ↓

    Domain Service

    ↓

    State Change

------------------------------------------------------------------------

# 13. Review Pipeline

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

    ↓

    Persist

------------------------------------------------------------------------

# 14. Review Context

Agent 不直接读取数据库。

通过 Context Provider 获取：

    Project Context

    Release Context

    Checklist Context

    Resource Context

    Deliverable Context

    Knowledge Context

    Review History

------------------------------------------------------------------------

# 15. Deterministic Validation

程序负责：

-   文件存在；
-   格式；
-   尺寸；
-   大小；
-   文本限制。

模型不负责机械检查。

------------------------------------------------------------------------

# 16. Capability Layer

不是模型聚合网关。

抽象能力：

    Text Reasoning

    Image Review

    Video Understanding

    Audio Understanding

    Embedding

结构：

    Capability

    ↓

    Provider Adapter

    ↓

    Model

------------------------------------------------------------------------

支持：

-   DeepSeek；
-   Seed；
-   Kimi；
-   Gemini；
-   Local Models。

------------------------------------------------------------------------

# 17. 模型分工

## DeepSeek

负责：

-   主 Agent；
-   推理；
-   工具调用；
-   报告生成。

------------------------------------------------------------------------

## Seed

负责：

-   图片；
-   视频；
-   音频分析。

------------------------------------------------------------------------

## Kimi K3

负责：

-   长视频；
-   长上下文分析。

------------------------------------------------------------------------

## Gemini

负责：

-   通用多模态；
-   备用分析。

------------------------------------------------------------------------

# 18. Evidence Architecture

每条发现必须包含：

    Finding

    ↓

    Evidence

    ↓

    Reason

    ↓

    Recommendation

------------------------------------------------------------------------

禁止：

无依据结论。

------------------------------------------------------------------------

# 19. Agent Security

原则：

-   最小数据发送；
-   用户控制；
-   可追踪。

------------------------------------------------------------------------

数据等级：

## Level 0

公开信息。

## Level 1

项目描述。

## Level 2

未发布运营资料。

## Level 3

商业敏感资料。

------------------------------------------------------------------------

# 20. AI Consent Flow

用户触发：

    Review Request

    ↓

    显示发送内容

    ↓

    用户确认

    ↓

    模型调用

------------------------------------------------------------------------

# 21. API Key 安全

禁止：

-   保存到项目；
-   保存到导出文件。

使用系统安全存储。

------------------------------------------------------------------------

# 22. Tool Security

分级：

    Read

    Analyze

    Suggest

    Mutate

v0.1 禁止：

-   删除项目；
-   自动上传 Steam；
-   执行系统命令。

------------------------------------------------------------------------

# 23. Prompt Injection 防护

外部内容永远是：

Data

不是：

Instruction。

优先级：

    System Policy

    >

    Workflow Rules

    >

    User Facts

    >

    Imported Content

    >

    Model Suggestions

------------------------------------------------------------------------

# 24. Audit Log

记录：

-   用户操作；
-   AI 调用；
-   权限确认；
-   Review 历史。

用途：

-   调试；
-   复盘；
-   建立信任。

------------------------------------------------------------------------

# 25. Agent Quality Evaluation

测试：

-   缺失资产检测；
-   素材风险发现；
-   Codex 约束遵守；
-   模型失败处理；
-   证据完整性。

------------------------------------------------------------------------

# 26. v0.1 Agent 范围

支持：

-   Release Diagnosis；
-   Deliverable Review；
-   Screenshot Review；
-   Store Copy Review；
-   Task Suggestion。

不支持：

-   自主执行任务；
-   自动营销；
-   自动联系媒体；
-   自动 Steam 操作。

------------------------------------------------------------------------

# 总结

Corvus Studio Agent 的核心：

> 一个受领域模型、工作流、工具和证据约束的独立游戏发行顾问系统。

它不是聊天机器人，而是连接：

-   Project Codex；
-   Checklist；
-   Resource；
-   Deliverable；
-   Review Pipeline；

的智能决策层。
