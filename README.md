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
│   │   ├── layouts/        # 页面布局
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
| `/` | 开屏页，封面背景 + 液态玻璃 Dock 导航 + 音乐播放器 + 标题特效 |
| `/blog` | 博客首页，文章列表 + 作者卡片 + 最新说说侧栏 |
| `/blog/:slug` | 博客文章详情（Markdown 渲染 + 评论） |
| `/about` | 关于页，GitHub 概览 + 创作者介绍 + 时间线 + 社交链接 |
| `/moments` | 说说页，Twitter/X 风格动态流（发图、点赞、转发、评论） |
| `/guestbook` | 留言板，Reddit 风格嵌套评论（投票、回复、排序） |
| `/friends` | 友情链接页，站点信息卡 + 友链墙 |
| `/login` | 登录页 |
| `/admin` | 后台仪表盘（需登录） |
| `/admin/posts` | 文章管理（创建 / 编辑 / 发布 / 定时发布） |
| `/admin/comments` | 评论审核（批准 / 拒绝 / 删除） |
| `/admin/friends` | 友链管理 |

## 快速开始

### 1. 启动后端依赖

```bash
cd backend
docker-compose up -d          # PostgreSQL + Redis
go run cmd/migrate/main.go    # 执行数据库迁移
go run cmd/seed/main.go       # 创建管理员账号（默认 admin / admin123）
go run main.go                # 启动 API 服务 :8080
```

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

**管理接口** (`/api/v1/admin`，需 JWT)

| 模块 | 端点 |
|------|------|
| 仪表盘 | `GET /dashboard/stats` `GET /dashboard/trend/views\|comments\|likes` |
| 文章管理 | `GET\|POST /posts` `PUT\|DELETE /posts/:id` `POST /posts/:id/publish\|unpublish\|schedule` |
| 评论审核 | `GET /comments` `POST /comments/:id/approve\|reject` `DELETE /comments/:id` |
| 友链管理 | `GET\|POST /friends` `PUT\|DELETE /friends/:id` |
| 审计日志 | `GET /audit-logs` |

## 特色组件

- **LiquidGlassCard** -- 项目统一液态玻璃容器
- **MusicPlayer** -- 开屏页内嵌播放器（LRC 歌词滚动）
- **HeroTitle** -- 鼠标悬停 squash-stretch 弹跳动画
- **AboutGithubProfile** -- GitHub 数据可视化（ECharts）

## 构建

```bash
# 前端
cd frontend && npm run build

# 后端
cd backend && go build -o miku-blog .
```
