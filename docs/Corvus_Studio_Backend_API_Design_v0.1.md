# Corvus Studio Backend API Design v0.1

> 文档状态：Draft\
> API 风格：REST + SSE

------------------------------------------------------------------------

# 1. API 原则

API 面向业务，不直接暴露数据库。

错误：

    query_table()

正确：

    get_release_status()

------------------------------------------------------------------------

# 2. API Prefix

    /api/v1

模块：

    projects
    releases
    checklist
    resources
    deliverables
    codex
    templates
    knowledge
    reviews
    reports
    settings

------------------------------------------------------------------------

# 3. Project API

创建：

    POST /api/v1/projects

详情：

    GET /api/v1/projects/{id}

Dashboard：

    GET /api/v1/projects/{id}/dashboard

返回：

-   readiness
-   blockers
-   warnings
-   next actions

------------------------------------------------------------------------

# 4. Release API

创建：

    POST /api/v1/releases

获取：

    GET /api/v1/releases/{id}

状态转换：

    POST /api/v1/releases/{id}/transition

------------------------------------------------------------------------

# 5. Checklist API

列表：

    GET /api/v1/releases/{id}/checklist

详情：

    GET /api/v1/checklist/{id}

创建：

    POST /api/v1/checklist

状态：

    POST /api/v1/checklist/{id}/transition

------------------------------------------------------------------------

# 6. Resource API

创建：

    POST /api/v1/resources

详情：

    GET /api/v1/resources/{id}

搜索：

    GET /api/v1/resources/search

检查：

    POST /api/v1/resources/{id}/check

------------------------------------------------------------------------

# 7. Deliverable API

列表：

    GET /api/v1/releases/{id}/deliverables

详情：

    GET /api/v1/deliverables/{id}

验证：

    POST /api/v1/deliverables/{id}/validate

------------------------------------------------------------------------

# 8. Asset Map API

获取关系图：

    GET /api/v1/projects/{id}/graph

创建关系：

    POST /api/v1/relationships

------------------------------------------------------------------------

# 9. Codex API

获取：

    GET /api/v1/projects/{id}/codex

新增：

    POST /api/v1/codex

更新：

    PUT /api/v1/codex/{id}

------------------------------------------------------------------------

# 10. Review API

创建：

    POST /api/v1/reviews

状态：

    GET /api/v1/reviews/{id}

SSE：

    GET /api/v1/reviews/{id}/stream

------------------------------------------------------------------------

# 11. Report API

获取：

    GET /api/v1/reports/{id}

导出：

    POST /api/v1/reports/{id}/export

------------------------------------------------------------------------

# 12. Settings API

Provider：

    GET /api/v1/settings/providers

测试：

    POST /api/v1/settings/providers/{id}/test

------------------------------------------------------------------------

# 13. API 与 Agent

API 与 Tool 不一一对应。

Agent 使用业务 Tool：

    search_resources()

    get_release_status()

    review_asset()

------------------------------------------------------------------------

# 14. 错误规范

统一：

``` json
{
 "error":{
   "code":"",
   "message":"",
   "recoverable":true
 }
}
```

------------------------------------------------------------------------

# 总结

Backend API 为：

-   React Frontend
-   Agent Tool
-   CLI

提供统一业务接口。
