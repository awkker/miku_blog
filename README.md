# NanaMiku Blog

亮色主题 + 液态玻璃 (Glassmorphism) 设计风格的个人博客，品牌色 `#39c5bb` (miku) + `#c084fc` (lavender)。

## 技术栈

### 前端 (`frontend/`)
- **Astro 6** -- 路由与静态页面生成 (SSG)
- **Vue 3** -- Islands 按需加载交互组件
- **Nano Stores** -- 轻量跨框架状态管理
- **Tailwind CSS** -- 原子化样式

### 后端 (`backend/`)
- **Go 1.21+** -- 核心语言
- **Hertz** -- 字节跳动开源 HTTP 框架
- **PostgreSQL** -- 数据库（GIN 索引全文检索）
- **Redis** -- 缓存 / Lua 脚本限流
- **sqlc** -- 原生 SQL -> 类型安全 Go 代码

## 仓库结构

```text
nanamiku-blog/
├── frontend/               # Astro + Vue 3 前端
│   ├── src/
│   │   ├── components/     # UI 组件（按域分目录）
│   │   ├── content/copy/   # 前端可 DIY 文案集中入口（站点/后台）
│   │   ├── layouts/        # 页面布局
│   │   ├── lib/            # 共享工具（api.ts 等）
│   │   ├── pages/          # 基于文件的路由
│   │   ├── stores/         # Nano Stores 状态
│   │   └── styles/         # 全局样式
│   └── package.json
├── backend/                # Go + Hertz 后端
│   ├── main.go             # 入口
│   ├── cmd/                # CLI 工具（migrate / seed）
│   ├── biz/
│   │   ├── bootstrap/      # 配置、DB、Redis、路由注册
│   │   ├── handler/        # HTTP 处理器（admin / public）
│   │   ├── service/        # 业务逻辑层
│   │   ├── middleware/      # 中间件（Auth/CORS/RateLimit/Visitor 等）
│   │   ├── jobs/           # 后台任务（友链健康检查）
│   │   ├── dto/            # 响应结构体
│   │   └── errcode/        # 错误码
│   ├── query/              # sqlc 生成代码（勿手动编辑）
│   ├── sql/                # 迁移文件 + SQL 查询
│   ├── docker-compose.yml  # PostgreSQL + Redis
│   └── sqlc.yaml
└── agents.md               # AI 协作规范
```

## 页面总览

| 路由 | 说明 |
|------|------|
| `/` | 开屏页，视差封面轮播 + 液态玻璃 Dock 导航 + 音乐播放器 + 标题特效 |
| `/blog` | 博客首页，动态文章列表 + 作者统计卡片 + 站点趋势图 + 最新说说侧栏 |
| `/blog/post?slug=xxx` | 博客文章详情（Markdown 渲染 + 点赞） |
| `/about` | 关于页，GitHub 概览 + 创作者介绍 + 时间线 + 社交链接 |
| `/moments` | 说说页，Twitter/X 风格动态流（仅展示，管理员通过后台发布） |
| `/guestbook` | 留言板，Reddit 风格嵌套评论（投票、回复、排序） |
| `/friends` | 友情链接页，站点信息卡 + 友链墙 |
| `/login` | 登录页 |
| `/admin` | 后台仪表盘（Traffic Overview、Pages/Sources/Environment/Location、Geo Distribution、Traffic 热力矩阵、近期事项） |
| `/admin/posts` | 文章管理（创建 / 编辑 / 草稿 / 发布 / 定时发布） |
| `/admin/comments` | 评论审核 + 反垃圾指标（批准 / 拒绝 / 删除 / 限流可视化） |
| `/admin/friends` | 友链管理 |
| `/admin/moments` | 说说管理（创建 / 编辑 / 草稿 / 发布 / 定时发布） |
| `/admin/backup` | 数据导出与备份（JSON / SQL 一键下载） |

## 快速开始

### 1. 启动后端依赖

```bash
cd backend
docker-compose up -d          # PostgreSQL + Redis
go run cmd/migrate/main.go    # 执行数据库迁移（可重复执行，已应用版本会自动跳过）
go run cmd/seed/main.go       # 创建管理员账号（默认 admin / admin123）
go run main.go                # 启动 API 服务 :8080
```

可选：如需本地也显示国家分布，可在 `backend/.env` 配置 `GEOIP_DB_PATH=/绝对路径/GeoLite2-City.mmdb`，重启后端后生效。

本地验证是否生效：
1. 用 `X-Forwarded-For: 8.8.8.8` 调一次 `POST /api/v1/analytics/collect`（不要带 `CF-IPCountry` 这类 geo header）。
2. 查询 `analytics_pageviews` 最新记录，确认 `country_code` 从 `ZZ` 变为两位国家码（如 `US`）。
3. 若仍是 `ZZ`，检查 `GEOIP_DB_PATH` 是否为有效绝对路径并重启后端。

### 2. 启动前端

```bash
cd frontend
npm install
npm run dev                   # 启动开发服务器 :4321
```

### 3. 访问

- 前端: `http://localhost:4321`
- 后端 API: `http://localhost:8080/api/v1/health`

## 后端 API 概览

**公开接口** (`/api/v1`)

| 模块 | 端点 |
|------|------|
| 认证 | `POST /auth/login` `POST /auth/refresh` `POST /auth/logout` |
| 文章 | `GET /posts` `GET /posts/hot` `GET /posts/search?q=` `GET /posts/:slug` `POST /posts/:id/like` |
| 评论 | `GET /posts/:id/comments` `POST /posts/:id/comments` |
| 留言板 | `GET /guestbook/messages` `POST /guestbook/messages` `POST /guestbook/messages/:id/vote` |
| 说说 | `GET /moments` `GET /moments/latest` `POST /moments` `POST /moments/:id/like` `POST /moments/:id/repost` `POST /moments/:id/comments` |
| 友链 | `GET /friends` |
| Analytics | `POST /analytics/collect` |

**管理接口** (`/api/v1/admin`，需 JWT)

| 模块 | 端点 |
|------|------|
| 仪表盘 | `GET /dashboard/stats` `GET /dashboard/trend/views\|comments\|likes` `GET /dashboard/analytics?range=24h\|7d\|30d&offset=0` |
| 文章管理 | `GET\|POST /posts` `PUT\|DELETE /posts/:id` `POST /posts/:id/publish\|unpublish\|schedule` |
| 说说管理 | `GET\|POST /moments` `PUT /moments/:id` `POST /moments/:id/publish\|unpublish\|schedule` |
| 评论审核 | `GET /comments` `POST /comments/:id/approve\|reject` `DELETE /comments/:id` |
| 反垃圾看板 | `GET /moderation/rate-limit-metrics?minutes=60` |
| 友链管理 | `GET\|POST /friends` `PUT\|DELETE /friends/:id` |
| 备份导出 | `GET /backup/export?format=json\|sql` |
| 审计日志 | `GET /audit-logs` |

## 最近更新

- **草稿 / 定时发布**：文章与说说均支持 `draft` / `published` / `scheduled` 生命周期，并由后端定时任务自动发布到点内容。
- **反垃圾增强**：留言与评论提交新增敏感词拦截；Redis 限流命中数据可在后台评论页可视化查看（总量 + 趋势 + 规则拆分）。
- **数据导出与备份**：新增后台一键导出 `JSON` / `SQL`，用于迁移和灾备。

## Analytics 仪表盘

- 顶部 KPI 卡片：Visitors / Visits / Views / Bounce rate / Visit duration（含环比变化）
- 时间窗口：`Last 24 hours`、`Last 7 days`、`Last 30 days`，支持 `offset` 查看历史窗口
- 趋势图：Visitors 与 Views 柱状对比（ECharts）
- 细分看板：Pages、Sources（Referrers/Channels）、Environment（Browsers/OS/Devices）、Location（Countries/Regions/Cities）
- 空间分布：世界地图 Geo Distribution + 小时/星期 Traffic 热力矩阵
- 运营辅助：保留「近期事项」管理员操作记录

## 文案 DIY

- 站点级文案入口：`frontend/src/content/copy/site.ts`
- 后台文案入口：`frontend/src/content/copy/admin.ts`
- 统一导出：`frontend/src/content/copy/index.ts`
- 建议：如果要改页面可见文案，优先改 copy 文件，再由页面/组件读取；不要在页面模板里直接写死字符串。

### `site.ts` 分区速查

- `home`：开屏页（`/`）标题、副标题、Dock 文案
- `blogTopNav`：博客导航（博客首页/说说/留言板/友情链接/关于/后台/登录）
- `blogIndex`：博客首页（`/blog`）Hero 区、右侧作者卡、近况/歌单、按钮文案
- `aboutPage`：关于我（`/about`）整页文案（简介、时间线、项目、写作地图、联系区等）
- `momentsPage`：说说页（`/moments`）标题、副标题、页面内快捷导航文案
- `guestbookPage`：留言板（`/guestbook`）标题、副标题、页面内快捷导航文案
- `components.blogFeed`：博客首页文章流组件文案（空态/错误/标签等）
- `components.aboutGithubProfile`：关于页 GitHub 模块文案
- `components.latestMoments`：博客侧栏 Latest Moments 文案
- `components.momentsBoard` / `components.momentCard`：说说流及卡片文案
- `components.guestbookBoard` / `components.guestbookMessageCard`：留言板表单/列表/回复文案

## 特色组件

- **LiquidGlassCard** -- 项目统一液态玻璃容器
- **HeroParallax** -- 开屏视差封面（鼠标追踪 3D 旋转 + 交叉淡入切换）
- **HeroShuffleBtn** -- Dock 栏随机换图按钮
- **MusicPlayer** -- 开屏页内嵌播放器（LRC 歌词滚动）
- **HeroTitle** -- 鼠标悬停 squash-stretch 弹跳动画
- **AuthorStats** -- 博客侧栏动态统计（文章数 / 分类 / 总浏览）
- **SiteTrend** -- 博客侧栏 SVG 趋势图（近 7 天访问热度）
- **AboutGithubProfile** -- GitHub 数据可视化（ECharts）
- **BlogFeed** -- 博客文章列表（从后端 API 动态加载，分页）
- **BlogPostView** -- 博客文章详情（API 加载 + marked 渲染 Markdown）

## 构建

```bash
# 前端
cd frontend && npm run build

# 后端
cd backend && go build -o miku-blog .
```
