# NanaMiku Blog Frontend

基于 Astro + Vue 3 Islands + Nano Stores + Tailwind CSS 的个人博客前端工程。
采用亮色主题 + 液态玻璃 (Glassmorphism) 设计风格，品牌色 `#39c5bb` (miku) + `#c084fc` (lavender)。

## 技术栈
- **Astro 6** — 路由与静态页面生成 (SSG)
- **Vue 3** — Islands 按需加载交互组件
- **Nano Stores** — 轻量跨框架状态管理
- **Tailwind CSS** — 原子化样式

## 页面总览

| 路由 | 说明 |
|------|------|
| `/` | 开屏页，封面背景 + 液态玻璃 Dock 导航 + 音乐播放器 + 标题特效 |
| `/blog` | 博客首页，从后端 API 动态加载文章列表 + 作者卡片 + 最新说说侧栏 |
| `/blog/post?slug=xxx` | 博客文章详情（API 加载 + Markdown 渲染 + 点赞） |
| `/about` | 关于页，GitHub 个人概览 + 创作者介绍 + 时间线 + 项目与写作地图 + 社交链接 |
| `/moments` | 说说页，Twitter/X 风格动态流（仅展示，管理员通过后台发布） |
| `/guestbook` | 留言板，Reddit 风格嵌套评论（投票、回复、排序） |
| `/friends` | 友情链接页，站点信息卡 + 友链墙 |
| `/login` | 登录页 |
| `/admin` | 后台仪表盘（需登录） |
| `/admin/posts` | 文章管理 |
| `/admin/comments` | 评论管理 |
| `/admin/friends` | 友链管理 |
| `/admin/moments` | 说说管理（创建 / 删除） |

## 目录结构

```text
frontend/
├── public/
├── src/
│   ├── assets/
│   ├── components/
│   │   ├── admin/      # 后台组件（侧栏、仪表盘、文章/评论/友链/说说管理）
│   │   ├── auth/       # 登录域组件
│   │   ├── base/       # BaseHead/日期等基础组件
│   │   ├── blog/       # 博客组件（BlogFeed / BlogPostView）
│   │   ├── friends/    # 友链页组件（FriendsGrid / FriendLinkCard）
│   │   ├── guestbook/  # 留言板组件（GuestbookBoard / GuestbookMessageCard）
│   │   ├── about/      # 关于页组件（AboutGithubProfile）
│   │   ├── home/       # 首页专用交互组件（MusicPlayer / HeroTitle / SystemClock / TypewriterSubtitle）
│   │   ├── moments/    # 说说页组件（MomentsBoard / MomentCard / LatestMoments，公开页仅展示）
│   │   ├── ui/         # 通用 UI 组件（按钮、输入、空态、玻璃卡等）
│   │   └── README.md   # 组件职责说明
│   ├── content/
│   ├── layouts/
│   ├── pages/
│   │   ├── admin/
│   │   │   ├── index.astro
│   │   │   ├── posts.astro
│   │   │   ├── comments.astro
│   │   │   ├── friends.astro
│   │   │   └── moments.astro
│   │   ├── blog/
│   │   ├── moments.astro
│   │   ├── login.astro
│   │   ├── guestbook.astro
│   │   └── friends.astro
│   ├── lib/            # 共享工具（api.ts HTTP 客户端）
│   ├── stores/         # auth/ui/loading/guestbook/friends/moments
│   ├── styles/
│   └── utils/
├── astro.config.mjs
├── package.json
└── tsconfig.json
```

## 特色组件

### 博客文章列表 (`BlogFeed.vue`)
从后端 `GET /api/v1/posts` 动态加载已发布文章，展示缩略图、分类、标签、真实阅读/点赞计数，支持分页。置顶文章大卡 + 网格布局。

### 博客文章详情 (`BlogPostView.vue`)
通过 URL 参数 `?slug=xxx` 从后端加载单篇文章，使用 `marked` 库客户端渲染 Markdown，支持点赞交互。

### 音乐播放器 (`MusicPlayer.vue`)
开屏页顶部栏内嵌播放器，支持播放/暂停、上一首/下一首、音量控制。点击歌曲名展开卡片，显示专辑封面、进度条、循环/静音控制、LRC 歌词自动滚动高亮（点击歌词可跳转播放）。音乐文件位于 `public/music/`。

### GitHub 个人概览 (`AboutGithubProfile.vue`)
关于页集成的 GitHub 数据展示组件，包含个人信息卡片、Repos/Stars/Followers 统计、ECharts 活动图表、技术栈分析、最近活跃仓库。客户端直接调用 GitHub API，localStorage 缓存 1 小时。

### 标题特效 (`HeroTitle.vue`)
开屏页标题文字拆分组件，鼠标滑过单个字符时触发 squash-stretch 弹跳动画（压扁 -> 弹起 -> 落回），配合薰衣草色光晕与变色。使用 `animationend` 事件防止动画期间重复触发。

## 关键说明
- `src/lib/api.ts` — 统一 HTTP 客户端（JWT 认证、错误处理、类型安全）
- `src/stores/auth.ts` — 登录态、登录/登出、localStorage 同步
- `src/stores/guestbook.ts` — 留言板状态（嵌套评论、投票、排序）
- `src/stores/moments.ts` — 说说状态（加载、点赞、转发、评论）
- `src/components/ui/LiquidGlassCard.vue` — 项目统一液态玻璃容器
- `src/components/admin/AdminRouteGuard.astro` — 后台页面前置鉴权脚本
- `src/layouts/AdminLayout.astro` — 后台统一壳层

## 本地开发

```bash
npm install
npm run dev
```

## 构建

```bash
npm run build
npm run preview
```

## 登录凭据
- 用户名: `admin`
- 密码: 后端 seed 时设置（默认 `admin123`）
