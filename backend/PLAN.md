# NanaMiku Blog 后端计划书

## 1. 文档目标

本计划书基于当前 `frontend/` 项目的页面、stores、后台 mock 组件与
`frontend/agents.md` 约束整理，目标不是只列接口，而是完整定义后端的：

- 业务范围
- 架构边界
- 数据模型
- 权限与鉴权
- PostgreSQL / Redis 分工
- 公共 API 与后台 API 设计
- 风控、审核、统计与异步任务
- 实施顺序与阶段目标

这份文档可以直接作为后续 `backend/` 目录的实施基线。

## 2. 已确定的技术硬约束

这些约束已经被前端项目锁定，后端必须遵守：

- 核心语言：Go 1.21+
- Web 框架：Hertz
- 数据库：PostgreSQL
- 数据访问：只用 `sqlc`
- 缓存与中间件：Redis
- 搜索：PostgreSQL 全文检索 + GIN 索引，不引入 Elasticsearch
- 高频接口限流：Redis Lua
- 敏感信息：只能来自环境变量
- 本地开发：兼容 macOS Apple Silicon
- 目录结构以 `backend/biz`、`backend/query`、`backend/sql` 为核心

## 3. 从前端反推出来的后端需求总览

### 3.1 前端已经明确存在的业务模块

- 管理员登录与后台受保护路由
- 博客文章列表与文章详情
- 文章点赞
- 后台文章管理
- 后台评论审核
- 留言板：嵌套回复、投票、排序
- 说说：发布、评论、点赞、转发
- 友链：前台展示、后台管理、健康检查
- 后台仪表盘：统计卡片与近 7 日趋势图
- 关于页 GitHub 数据面板

### 3.2 前端虽然没直接写死，但已经被界面隐含出来的需求

- 匿名访客必须有稳定身份，否则以下功能无法落地：
  - 留言板投票
  - 文章点赞
  - 说说点赞
  - 说说转发
  - 说说评论点赞
- 后台统计图不能每次都扫原始表，必须有聚合层
- 既然留言板和说说允许公开写入，审核与风控就是必需项
- 后台已经出现“文章管理”，但当前博客正文仍来自 Astro Content
  Collections，这意味着内容源必须重新设计
- 后台友链管理已经出现“健康检测”，因此需要定时任务
- 后台评论审核页面已经假定“文章评论”存在，即使前台文章评论区尚未接入

## 4. 本计划采用的工作假设

为便于先把后端边界完整定下来，当前先采用以下假设：

- `moments` 是管理员专属动态流
- 留言板和说说都允许匿名访客参与
- 文章评论后端先建设，前台评论 UI 后续再接
- 文章内容最终会迁移为“后端为内容源，Astro 为渲染层”
- GitHub 数据后续建议由后端代理和缓存，不再完全依赖浏览器直连

如果后续产品方向变化，这份计划的大框架仍然成立，只需要收缩对应公开接口。

## 5. 后端范围定义

### 5.1 P0 必做模块

- 鉴权与管理员会话
- 匿名访客身份体系
- 文章与内容管理
- 文章点赞与阅读统计
- 文章评论与审核
- 留言板与投票
- 说说与互动
- 友链、友链申请、友链健康检查
- 后台仪表盘统计
- Redis 限流、热榜、去重
- 审核基础能力与操作审计

### 5.2 P1 强烈建议模块

- GitHub 代理与缓存
- 图片上传 / 媒体服务
- 敏感词与黑名单管理面板
- 定时发布任务
- 文章版本 diff 与回滚 UI

### 5.3 第一版明确不做

- 公开用户注册系统
- 复杂 RBAC
- WebSocket 实时推送
- Elasticsearch
- ORM 框架

## 6. 总体架构

### 6.1 部署形态

```text
Astro 前端
    |
    | HTTP / JSON
    v
Hertz API
    |
    +-- PostgreSQL  数据真源
    |
    +-- Redis       限流 / 缓存 / 热榜 / 去重
```

### 6.2 API 分层

- 公共 API：前台读取数据与匿名互动
- 后台 API：管理员登录后使用
- 内部任务：定时发布、统计汇总、友链健康检查

### 6.3 推荐目录结构

```text
backend/
├── biz/
│   ├── handler/
│   │   ├── public/
│   │   ├── admin/
│   │   └── internal/
│   ├── service/
│   ├── middleware/
│   ├── dto/
│   ├── jobs/
│   └── bootstrap/
├── query/
├── sql/
│   ├── migrations/
│   └── queries/
├── scripts/
├── docker-compose.yml
└── PLAN.md
```

## 7. 身份、鉴权与权限设计

### 7.1 管理员鉴权

#### 必备能力

- 登录
- 登出
- 获取当前管理员信息
- 刷新会话
- 受保护后台接口鉴权
- 关键管理行为审计

#### 推荐实现

- `admin_users` 表保存管理员
- 密码使用 `argon2id` 或 `bcrypt` 哈希
- 短期 access token + 长期 refresh token
- refresh token 持久化到数据库
- 后台接口统一走鉴权中间件

#### 与前端的衔接提醒

当前前端登录态是 `localStorage` mock。真实后端上线后，建议改成更安全的
`HttpOnly Cookie` 或至少“短 access + 持久 refresh”的会话模式，这意味着
前端后台守卫逻辑后面要做一轮配套改造。

### 7.2 匿名访客身份

当前前端的大量交互都是匿名可操作，因此后端必须补一个“访客身份层”。

#### 必备设计

- 首次访问互动接口时签发 `visitor_id`
- 使用签名 Cookie 存储
- 点赞、投票、转发都基于 `visitor_id`
- 同时结合 IP Hash 与 UA Hash 做风控

没有这一层，前端现在的互动状态全部只能停留在 mock。

## 8. 业务域详细规划

### 8.1 文章系统

#### 前台能力

- 获取已发布文章列表
- 按 slug 获取文章详情
- 支持最新排序
- 支持热度排序
- 支持分类与标签筛选
- 支持全文搜索
- 支持相关文章推荐
- 返回阅读元信息
- 返回点赞数、评论数、热度值

#### 后台能力

- 新建文章
- 编辑文章
- 删除文章
- 草稿保存
- 立即发布
- 定时发布
- 下线文章
- 分类与标签管理
- 封面图管理
- 摘要与 SEO 字段管理
- 版本记录
- 回滚历史版本

#### 当前项目最大的架构冲突

现在博客正文来自 `frontend/src/content/blog` 的 Astro 内容集合，而后台又已经有
“文章管理”页面。两者不能长期并存为两个真源。

#### 推荐方案

1. 后端成为内容真源。
2. 正文以 Markdown 存在 PostgreSQL。
3. Astro 在构建期从后端拉取已发布文章并生成静态页。
4. 后台写入只触达后端，不再改前端仓库里的 Markdown 文件。

#### 兼容过渡方案

1. 先保留 Astro Markdown 作为只读内容源。
2. 先做登录、留言板、说说、友链、统计这几类后端。
3. 等确认内容迁移方案后，再做真实文章 CRUD。

由于当前后台 UI 已经有文章管理，建议尽早明确内容迁移方案，而不是继续拖。

### 8.2 文章点赞与阅读统计

#### 必备能力

- 文章点赞 / 取消点赞
- 返回当前访客是否已点赞
- 记录文章浏览
- 对重复浏览做短周期去重
- 聚合日级 PV / UV
- 计算文章热度值

#### Redis 责任

- 浏览去重键：`view:post:{post_id}:{visitor_id}:{day}`
- 热榜 ZSET：`rank:posts:hot`
- 首页热门文章可做短 TTL 缓存

### 8.3 文章评论

虽然前台文章页还没挂评论区，但后台评论审核已经说明这个域必须存在。

#### 前台能力

- 提交文章评论
- 读取已通过审核的评论
- 后续可扩展回复

#### 后台能力

- 按状态查看评论
- 审核通过
- 驳回
- 删除
- 按文章筛选
- 按时间筛选
- 按风险等级筛选

#### 审核策略

- 默认 `pending`，或按规则低风险自动通过
- 敏感词检测
- IP / UA / 域名黑名单检测
- 基础反垃圾规则

### 8.4 留言板

#### 现有前端已经明确需要的能力

- 留言列表
- 排序：`hot` / `newest` / `oldest`
- 发表主留言
- 发表一层回复
- 赞同 / 反对
- 切换投票状态
- 可选网站字段
- 返回当前访客的投票状态

#### 审核建议

- 第一版可对低风险留言自动通过
- 但表结构里仍然保留审核状态，后续便于接后台审核

### 8.5 说说系统

#### 当前前端需要的能力

- 说说列表
- 发布说说
- 最多 4 张图片 URL
- 点赞 / 取消点赞
- 转发 / 取消转发
- 发表评论
- 评论点赞 / 取消点赞
- 提供最新 3 条摘要给博客侧栏

#### 后端必须补上的能力

- 说说审核状态
- 隐藏违规说说
- 隐藏违规评论
- 发布与评论的限流

如果说说继续保持“公开写入”，这是全站最需要风控的模块之一。

### 8.6 友链系统

#### 前台能力

- 获取已通过审核的友链列表
- 后续可开放友链申请接口

#### 后台能力

- 查看全部友链
- 手动新增
- 编辑
- 删除
- 审核申请
- 调整展示排序
- 触发站点健康重检

#### 异步任务能力

- 定时探测友链可访问性
- 记录 HTTP 状态码、耗时、最近检查时间
- 将友链标记为 `ok` / `down` / `unknown`

### 8.7 后台仪表盘

#### 当前 UI 已经明确的统计卡片

- 文章总数
- 待审评论数
- 累计点赞数
- 友链数量

#### 当前 UI 已经明确的趋势图

- 近 7 日浏览量趋势
- 近 7 日评论趋势
- 近 7 日点赞趋势

#### 当前 UI 已经暗示的活动流

- 新草稿创建
- 待审核评论
- 友链申请
- 审核动作

后台仪表盘必须读聚合数据，不能直接依赖高频全表扫描。

### 8.8 About 页 GitHub 数据

当前前端是浏览器直接请求 GitHub API 并写入 `localStorage` 缓存。

#### 推荐后端化方案

- 后端代理 GitHub 请求
- 缓存 1 小时左右
- 统一输出前端所需结构
- 降低浏览器侧速率限制和失败波动

这个模块不是第一阶段必须项，但非常值得在二期补齐。

## 9. 数据模型规划

下面是建议的首版数据表集合。

### 9.1 鉴权与访客

- `admin_users`
  - `id`
  - `username`
  - `email`
  - `password_hash`
  - `role`
  - `status`
  - `last_login_at`
  - `created_at`
  - `updated_at`

- `admin_refresh_tokens`
  - `id`
  - `admin_user_id`
  - `token_hash`
  - `expires_at`
  - `revoked_at`
  - `created_at`
  - `last_used_at`

- `visitors`
  - `id`
  - `first_seen_at`
  - `last_seen_at`
  - `ip_hash`
  - `ua_hash`

### 9.2 文章内容

- `posts`
  - `id`
  - `slug`
  - `title`
  - `excerpt`
  - `content_markdown`
  - `hero_image_url`
  - `category`
  - `status`：`draft` / `published` / `scheduled` / `archived`
  - `published_at`
  - `scheduled_at`
  - `created_by`
  - `updated_by`
  - `created_at`
  - `updated_at`
  - `search_vector`

- `tags`
  - `id`
  - `name`
  - `slug`

- `post_tags`
  - `post_id`
  - `tag_id`

- `post_revisions`
  - `id`
  - `post_id`
  - `title`
  - `excerpt`
  - `content_markdown`
  - `hero_image_url`
  - `category`
  - `editor_id`
  - `created_at`

### 9.3 文章互动

- `post_likes`
  - `post_id`
  - `visitor_id`
  - `created_at`

- `post_view_daily`
  - `post_id`
  - `day`
  - `pv`
  - `uv`

- `post_comments`
  - `id`
  - `post_id`
  - `parent_id`
  - `author_name`
  - `author_email`
  - `author_website`
  - `content`
  - `status`
  - `ip_hash`
  - `ua_hash`
  - `created_at`
  - `approved_at`
  - `reviewed_by`

### 9.4 留言板

- `guestbook_messages`
  - `id`
  - `parent_id`
  - `author_name`
  - `author_website`
  - `content`
  - `status`
  - `vote_score`
  - `ip_hash`
  - `ua_hash`
  - `created_at`
  - `reviewed_by`

- `guestbook_votes`
  - `message_id`
  - `visitor_id`
  - `vote`
  - `created_at`
  - `updated_at`

### 9.5 说说

- `moments`
  - `id`
  - `author_name`
  - `content`
  - `image_urls`（`jsonb`）
  - `status`
  - `like_count`
  - `repost_count`
  - `comment_count`
  - `ip_hash`
  - `ua_hash`
  - `created_at`
  - `reviewed_by`

- `moment_likes`
  - `moment_id`
  - `visitor_id`
  - `created_at`

- `moment_reposts`
  - `moment_id`
  - `visitor_id`
  - `created_at`

- `moment_comments`
  - `id`
  - `moment_id`
  - `author_name`
  - `content`
  - `status`
  - `like_count`
  - `ip_hash`
  - `ua_hash`
  - `created_at`
  - `reviewed_by`

- `moment_comment_likes`
  - `comment_id`
  - `visitor_id`
  - `created_at`

### 9.6 友链

- `friend_links`
  - `id`
  - `name`
  - `description`
  - `url`
  - `domain`
  - `avatar_url`
  - `status`：`approved` / `pending` / `rejected`
  - `health_status`：`ok` / `down` / `unknown`
  - `sort_order`
  - `created_at`
  - `approved_at`
  - `reviewed_by`
  - `last_checked_at`

- `friend_link_health_logs`
  - `id`
  - `friend_link_id`
  - `http_status`
  - `latency_ms`
  - `result`
  - `checked_at`

#### 如果后续开放公开申请，再补

- `friend_link_applications`
  - `id`
  - `site_name`
  - `site_url`
  - `avatar_url`
  - `description`
  - `contact_email`
  - `contact_note`
  - `status`
  - `created_at`
  - `reviewed_at`
  - `review_note`

### 9.7 审核与审计

- `sensitive_words`
- `blocked_ips`
- `audit_logs`

第一版可以先做轻量结构，但表建议尽早建立。

## 10. API 规划

默认统一前缀：`/api/v1`

### 10.1 鉴权

- `POST /auth/login`
- `POST /auth/logout`
- `POST /auth/refresh`
- `GET /auth/me`

### 10.2 前台文章

- `GET /posts`
- `GET /posts/:slug`
- `GET /posts/:slug/related`
- `POST /posts/:slug/view`
- `POST /posts/:slug/like`
- `GET /search/posts`

### 10.3 后台文章管理

- `GET /admin/posts`
- `POST /admin/posts`
- `GET /admin/posts/:id`
- `PATCH /admin/posts/:id`
- `DELETE /admin/posts/:id`
- `POST /admin/posts/:id/publish`
- `POST /admin/posts/:id/schedule`
- `POST /admin/posts/:id/unpublish`
- `GET /admin/posts/:id/revisions`
- `POST /admin/posts/:id/rollback`

### 10.4 文章评论

- `GET /posts/:slug/comments`
- `POST /posts/:slug/comments`
- `GET /admin/comments`
- `POST /admin/comments/:id/approve`
- `POST /admin/comments/:id/reject`
- `DELETE /admin/comments/:id`

### 10.5 留言板

- `GET /guestbook/messages?sort=hot|newest|oldest`
- `POST /guestbook/messages`
- `POST /guestbook/messages/:id/vote`

### 10.6 说说

- `GET /moments`
- `GET /moments/latest`
- `POST /moments`
- `POST /moments/:id/like`
- `POST /moments/:id/repost`
- `POST /moments/:id/comments`
- `POST /moments/comments/:id/like`

#### 建议补上的后台接口

- `GET /admin/moments`
- `POST /admin/moments/:id/hide`
- `POST /admin/moment-comments/:id/hide`

### 10.7 友链

- `GET /friends`
- `POST /friends/applications`
- `GET /admin/friends`
- `POST /admin/friends`
- `PATCH /admin/friends/:id`
- `DELETE /admin/friends/:id`
- `POST /admin/friends/:id/approve`
- `POST /admin/friends/:id/reject`
- `POST /admin/friends/:id/recheck`

### 10.8 后台仪表盘

- `GET /admin/dashboard/overview`
- `GET /admin/dashboard/trends?days=7`
- `GET /admin/dashboard/activities`

### 10.9 About / 集成

- `GET /integrations/github/:username`

## 11. SQLC 组织建议

建议按业务域拆 SQL 文件，而不是做一个大杂烩：

```text
sql/queries/
├── admin_auth.sql
├── visitors.sql
├── posts.sql
├── post_comments.sql
├── guestbook.sql
├── moments.sql
├── friends.sql
├── dashboard.sql
└── moderation.sql
```

原则是：

- SQL 尽量显式、可读
- 业务逻辑放到 service 层
- 查询文件按域分开，便于 sqlc 生成和维护

## 12. Redis 规划

### 12.1 必用场景

- 登录、留言、评论、说说发布、点赞、转发、友链申请的 Lua 限流
- 文章和说说热榜 ZSET
- 仪表盘短 TTL 缓存
- 浏览去重
- 互动状态加速查询

### 12.2 建议 Key 设计

- `rl:login:{ip}`
- `rl:guestbook:create:{ip}`
- `rl:moment:create:{ip}`
- `rl:moment:comment:{visitor_id}`
- `view:post:{post_id}:{visitor_id}:{day}`
- `rank:posts:hot`
- `rank:moments:hot`
- `cache:dashboard:overview`
- `cache:github:{username}`

## 13. 风控与审核

### 13.1 第一版必须具备

- 登录限流
- 公开写接口限流
- 短周期重复内容拦截
- URL 格式校验
- 图片数量限制
- 基础敏感词匹配
- 管理端审核动作审计日志

### 13.2 第二步建议补齐

- 风险分数
- IP / 域名黑名单
- 可疑昵称与高频行为检测
- 触发式验证码

## 14. 统计设计

### 14.1 后端必须支撑的指标

- 文章总数
- 待审核评论数
- 累计点赞数
- 友链数量
- 每日浏览量
- 每日评论量
- 每日点赞量
- 最近管理动作

### 14.2 聚合策略

建议采用：

1. 互动先写主业务表。
2. 能安全增量更新的计数器同步更新。
3. 日级统计由定时任务汇总。
4. 仪表盘优先读聚合结果。

这样既能保证写路径简单，又能保证后台查询快。

## 15. 异步任务规划

后端从第一版起就需要具备基础任务能力。

### 必做任务

- 定时发布文章
- 友链健康检查
- 日级统计汇总
- 清理过期 refresh token

### 建议任务

- GitHub 缓存刷新
- 热榜重算
- 长时间未处理审核项提醒

## 16. 前端联动点

### 16.1 后续会接后端的前端状态域

- `src/stores/auth.ts`
- `src/stores/guestbook.ts`
- `src/stores/moments.ts`
- `src/stores/friends.ts`
- admin 仪表盘与管理组件

### 16.2 后端落地后，前端需要配合改造的地方

- 登录改为真实鉴权流程
- 后台守卫不能再只依赖 `localStorage`
- mock 数组改为真实 API
- 文章评论区前端补接
- 确认说说是否继续保持“游客可发布”

## 17. 实施阶段建议

### Phase 0：基础脚手架

- 初始化 Go module
- 初始化 Hertz
- 加入配置加载
- 接 PostgreSQL 与 Redis
- 配好 migration 与 sqlc
- 写本地 `docker-compose`

### Phase 1：平台基础能力

- 管理员表
- 登录 / 登出 / me / refresh
- 鉴权中间件
- 访客身份 Cookie
- 请求 ID、日志、recover 中间件
- Redis Lua 限流

### Phase 2：先替换前台互动 mock

- 留言板
- 说说
- 友链前台接口
- 友链申请接口
- 基础审核状态字段

这一阶段完成后，前台主要互动页就能摆脱 mock store。

### Phase 3：后台运营能力

- 仪表盘总览
- 评论审核
- 友链管理
- 说说审核
- 审计日志

### Phase 4：文章系统迁移

- 文章表与版本表
- 后台文章 CRUD
- 发布与定时发布
- 搜索
- 文章点赞与统计
- Astro 从后端拉内容构建

这是整个项目最关键的一期，因为它会真正解决“静态 Markdown 与后台文章管理”的冲突。

### Phase 5：稳定性与增强

- 健康检查任务完善
- 统计聚合优化
- GitHub 代理缓存
- 黑名单与敏感词配置
- 性能测试与压测

## 18. 验收标准

满足以下条件后，可以认为后端首版具备可用性：

- 管理员可以安全登录后台
- 留言板可真实持久化并支持投票
- 说说可真实持久化并支持互动
- 友链列表与后台审核可用
- 仪表盘展示真实聚合数据
- 后台评论审核可用
- 高频接口已做限流
- PostgreSQL schema 与 sqlc 查询结构稳定
- Redis 已实际承担限流、热榜与去重职责
- 文章内容迁移路径已经明确并排期

## 19. 需要尽早确认的关键决策

### 19.1 内容真源到底是谁

当前博客正文是 Astro Markdown，后台又要做文章管理，这个问题必须先定。否则文章系统会一直处在“前台一套真源，后台一套假管理”的状态。

### 19.2 说说是否继续开放游客发布

如果继续开放，风控与审核就不能省；如果改成管理员发布，接口和后台需求都要相应收缩。

### 19.3 登录态最终采用什么安全模型

当前前端是 `localStorage` mock，真实后端建议改成 Cookie 或 refresh token 体系，这需要前端配合。

### 19.4 文章评论前台接入时机

后台评论审核已经要求这个域存在，但前台文章页还没评论 UI，需要明确接入优先级。

### 19.5 媒体策略是否要一步到位

现在很多地方还是图片 URL 或静态资源。如果后台文章管理要做得完整，最终需要上传方案，或者至少先定好外链媒体策略。

## 20. 推荐的实际开工顺序

如果现在马上开始做，最务实的顺序是：

1. Hertz 脚手架 + 配置 + PostgreSQL + Redis + sqlc
2. 管理员鉴权 + 访客身份 + 限流中间件
3. 留言板与说说接口
4. 友链接口 + 健康检查任务
5. 仪表盘 + 审核接口
6. 文章系统迁移与后台文章 CRUD
7. 统计增强与 GitHub 代理

这个顺序最贴合当前前端现状，也能把“文章内容迁移”这个最大架构决策单独留在一阶段处理，避免前期全线耦合。
