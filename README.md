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
| `/` | 开屏页，封面背景 + 液态玻璃 Dock 导航 |
| `/blog` | 博客首页，文章列表 + 作者卡片 + 最新说说侧栏 |
| `/blog/:id` | 博客文章详情 |
| `/moments` | 说说页，Twitter/X 风格动态流（发图、点赞、转发、评论） |
| `/guestbook` | 留言板，Reddit 风格嵌套评论（投票、回复、排序） |
| `/friends` | 友情链接页，站点信息卡 + 友链墙 |
| `/login` | 登录页 |
| `/admin` | 后台仪表盘（需登录） |
| `/admin/posts` | 文章管理 |
| `/admin/comments` | 评论管理 |
| `/admin/friends` | 友链管理 |

## 目录结构

```text
frontend/
├── public/
├── src/
│   ├── assets/
│   ├── components/
│   │   ├── admin/      # 后台相关组件（侧栏、顶栏、仪表盘等）
│   │   ├── auth/       # 登录域组件
│   │   ├── base/       # BaseHead/日期等基础组件
│   │   ├── blog/       # 博客阅读与总览组件
│   │   ├── friends/    # 友链页组件（FriendsGrid / FriendLinkCard）
│   │   ├── guestbook/  # 留言板组件（GuestbookBoard / GuestbookMessageCard）
│   │   ├── home/       # 首页专用交互组件
│   │   ├── moments/    # 说说页组件（MomentsBoard / MomentCard / LatestMoments）
│   │   ├── ui/         # 通用 UI 组件（按钮、输入、空态、玻璃卡等）
│   │   └── README.md   # 组件职责说明
│   ├── content/
│   ├── layouts/
│   ├── pages/
│   │   ├── admin/
│   │   │   ├── index.astro
│   │   │   ├── posts.astro
│   │   │   ├── comments.astro
│   │   │   └── friends.astro
│   │   ├── blog/
│   │   ├── moments.astro
│   │   ├── login.astro
│   │   ├── guestbook.astro
│   │   └── friends.astro
│   ├── stores/         # auth/ui/loading/guestbook/friends/moments
│   ├── styles/
│   └── utils/
├── astro.config.mjs
├── package.json
└── tsconfig.json
```

## 关键说明
- `src/stores/auth.ts` — 登录态、登录/登出、localStorage 同步
- `src/stores/guestbook.ts` — 留言板状态（嵌套评论、投票、排序）
- `src/stores/moments.ts` — 说说状态（发布、点赞、转发、评论）
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

## 当前登录凭据（前端 mock）
- 用户名: `admin`
- 密码: `miku1234`
