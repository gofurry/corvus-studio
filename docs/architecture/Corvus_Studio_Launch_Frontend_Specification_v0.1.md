# Corvus Studio Launch Frontend Specification v0.1

## 1. 定位

React 前端负责用户交互、状态展示、资源浏览和 Agent 结果展示。

不负责业务判断、AI 调用和数据持久化。

技术栈：

-   React 19
-   TypeScript
-   Vite 8
-   pnpm
-   Ant Design
-   SCSS
-   Zustand
-   TanStack Query
-   React Router
-   React Flow
-   react-markdown

------------------------------------------------------------------------

## 2. 页面结构

-   Home
-   Release
-   Checklist
-   Resources
-   Asset Map
-   Codex
-   Reports
-   Settings
-   Advisor

------------------------------------------------------------------------

## 3. 数据层

API Layer： 负责后端通信。

Query Layer： 使用 TanStack Query 管理服务器状态。

UI State： 使用 Zustand 管理页面和用户状态。

------------------------------------------------------------------------

## 4. 核心页面

### Home

发布准备驾驶舱：

-   Release Goal
-   Readiness
-   Blockers
-   Next Actions
-   Activity

### Checklist

三栏：

分类 / 任务列表 / 任务详情

### Resources

管理：

-   图片
-   视频
-   音频
-   文档
-   URL

### Asset Map

使用 React Flow 展示 Resource、Deliverable、Checklist、Codex 关系。

------------------------------------------------------------------------

## 5. AI 交互

Agent 作为 Advisor 存在。

不是首页聊天框。

支持：

-   Review
-   Explain
-   Suggest

------------------------------------------------------------------------

## 6. 体验能力

支持：

-   Command Palette
-   Onboarding
-   Offline UI
-   Global Search
-   Empty State
-   Error State
-   国际化

------------------------------------------------------------------------

## 7. 测试

-   Vitest
-   React Testing Library
-   Playwright（未来）
