# Corvus Studio ADK Implementation Plan v0.1

> 文档状态：Draft\
> Agent Runtime：Google ADK Go

------------------------------------------------------------------------

# 1. ADK 定位

ADK 是 Corvus 的 Agent 工作流执行层。

不是整个 Backend。

架构：

    React

    ↓

    Go Core

    ↓

    Agent Service

    ↓

    ADK Runtime

    ↓

    Tools / Workflow

------------------------------------------------------------------------

# 2. Agent Service

逻辑独立，运行在 Go Core 内。

负责：

-   Agent 生命周期
-   Session
-   Workflow
-   Tool Registry
-   Context

------------------------------------------------------------------------

# 3. Root Agent

v0.1：

单 Root Agent。

名称：

Corvus Advisor

职责：

-   理解用户意图
-   选择 Workflow
-   解释结果

------------------------------------------------------------------------

# 4. Workflow

采用：

代码 Workflow + Agent 选择。

不使用自由规划。

------------------------------------------------------------------------

Workflow：

-   Release Diagnosis
-   Deliverable Review
-   Screenshot Review
-   Copy Review

------------------------------------------------------------------------

# 5. Session

分为：

## Chat Session

短期：

-   对话
-   工具调用

## Review Job

长期：

-   输入
-   模型
-   结果
-   报告

------------------------------------------------------------------------

# 6. Tool Registry

静态注册：

    Project Tools
    Release Tools
    Checklist Tools
    Resource Tools
    Deliverable Tools
    Review Tools

------------------------------------------------------------------------

# 7. Context Provider

按需构建。

包含：

-   Project Context
-   Release Context
-   Resource Context
-   Deliverable Context
-   Knowledge Context
-   Review History

------------------------------------------------------------------------

# 8. Review Job

使用：

SQLite + Go Worker。

流程：

    Create Job

    ↓

    Worker

    ↓

    Workflow

    ↓

    Report

------------------------------------------------------------------------

# 9. Capability Layer

不直接调用模型。

流程：

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
-   Local Model

------------------------------------------------------------------------

# 10. Prompt Registry

Prompt 独立管理：

    Identity
    Policy
    Workflow
    Tool Instructions
    Domain

支持版本化。

------------------------------------------------------------------------

# 11. Trace

开发模式保存：

-   Prompt
-   Context
-   Tool Calls
-   Model Calls

生产保存摘要。

------------------------------------------------------------------------

# 12. Error Recovery

支持：

-   Tool Retry
-   Provider Failure
-   Context Error

禁止编造结果。

------------------------------------------------------------------------

# 13. 测试

包括：

-   Unit Test
-   Scenario Test
-   Regression Test

------------------------------------------------------------------------

# 14. v0.1 范围

实现：

-   Root Agent
-   Tool Registry
-   Context Provider
-   Workflow
-   DeepSeek
-   Seed
-   Report Generation

不实现：

-   Multi Agent
-   Autonomous Loop
-   Agent Marketplace

------------------------------------------------------------------------

# 总结

Corvus ADK 实现：

    Agent

    ↓

    Workflow

    ↓

    Tools

    ↓

    Domain Services

    ↓

    Models

形成可控、可解释的领域 Agent。
