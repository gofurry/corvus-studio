# Corvus Studio Data Model Design v0.1

> 文档状态：Draft\
> 产品：Corvus Studio Launch Workspace

## 1. 设计原则

-   领域对象优先，不从数据库表开始设计。
-   重要对象支持历史追踪。
-   Resource、Deliverable、Checklist 通过关系模型连接。
-   用户数据与系统知识分离。

------------------------------------------------------------------------

# 2. 核心实体

    Project
    ReleaseGoal
    ChecklistItem
    Resource
    Deliverable
    Relationship
    CodexEntry
    Template
    KnowledgeEntry
    ReviewJob
    ReviewReport
    ReviewFinding
    AuditLog
    UserConfig
    ProjectConfig

------------------------------------------------------------------------

# 3. Project

表示一个游戏项目。

字段概念：

-   id
-   name
-   description
-   location
-   steam_app_id
-   language
-   stage
-   status
-   created_at
-   updated_at

状态：

    Active
    Archived

------------------------------------------------------------------------

# 4. ReleaseGoal

表示发布目标。

v0.1：

-   Steam Coming Soon Page

未来：

-   Demo Release
-   Playtest
-   Full Release
-   Update

状态：

    Draft
    Preparing
    NeedsAttention
    ReadyForReview
    Ready
    Submitted

------------------------------------------------------------------------

# 5. ChecklistItem

发布准备任务。

来源：

-   PlatformTemplate
-   CorvusTemplate
-   User
-   Agent

状态：

    NotStarted
    InProgress
    NeedsReview
    Done
    Blocked
    NotApplicable

Checklist 与 Deliverable 使用关联实体支持多对多。

------------------------------------------------------------------------

# 6. Resource

资源对象。

类型：

-   Image
-   Video
-   Audio
-   Document
-   Text
-   URL
-   Folder

资源位置：

    Managed
    Referenced
    Linked

状态：

    Available
    Missing
    Changed
    Unavailable

------------------------------------------------------------------------

# 7. ResourceVersion

支持素材版本：

    Resource

    v1
    v2
    v3

保存：

-   hash
-   location
-   created_at

------------------------------------------------------------------------

# 8. Deliverable

发布交付对象。

类型：

-   File
-   Text
-   ExternalAction
-   Decision

状态：

    Missing
    Draft
    Reviewing
    Approved
    Published

------------------------------------------------------------------------

# 9. DeliverableVersion

记录交付物版本。

例如：

Main Capsule：

-   v1
-   v2
-   v3

------------------------------------------------------------------------

# 10. Relationship

Asset Map 核心。

允许实体间关系。

关系：

-   DerivedFrom
-   UsedBy
-   Supports
-   Blocks
-   Replaces
-   RelatedTo

UI 层限制展示范围，数据层保持扩展性。

------------------------------------------------------------------------

# 11. CodexEntry

Project Memory。

分类：

-   Identity
-   Gameplay
-   Audience
-   SellingPoint
-   Constraint
-   Decision

状态：

    Confirmed
    Temporary
    NeedDecision
    AISuggestion

------------------------------------------------------------------------

# 12. Template

模板系统。

类型：

-   Platform
-   Corvus
-   User

模板生成项目副本，不直接绑定。

------------------------------------------------------------------------

# 13. KnowledgeEntry

知识系统。

范围：

    System
    Project
    User

------------------------------------------------------------------------

# 14. ReviewJob

一次 AI 检查任务。

类型：

-   ReleaseDiagnosis
-   DeliverableReview
-   ScreenshotReview
-   CopyReview

状态：

    Created
    Collecting
    Validating
    Analyzing
    Completed
    Failed

------------------------------------------------------------------------

# 15. ReviewReport

保存结构化结果。

包含：

-   Summary
-   Findings

------------------------------------------------------------------------

# 16. ReviewFinding

每个发现：

包含：

-   Severity
-   Title
-   Finding
-   Evidence
-   Recommendation

------------------------------------------------------------------------

# 17. AuditLog

记录：

-   User
-   Agent
-   System

产生的关键行为。

------------------------------------------------------------------------

# 总结

Corvus Studio 数据模型围绕：

    Project

    ↓

    ReleaseGoal

    ↓

    Checklist

    ↓

    Deliverable

    ↓

    Resource

    ↓

    Review

构建。
