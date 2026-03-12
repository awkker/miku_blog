# NanaMiku Blog Frontend

基于 Astro + Vue 3 Islands + Nano Stores + Tailwind CSS 的前端工程。

## 技术栈
- Astro 6
- Vue 3 (Islands)
- Nano Stores
- Tailwind CSS

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
│   │   ├── friends/    # 友链页组件
│   │   ├── guestbook/  # 留言板组件
│   │   ├── home/       # 首页专用交互组件
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
│   │   ├── login.astro
│   │   ├── guestbook.astro
│   │   └── friends.astro
│   ├── stores/         # auth/ui/loading/guestbook/friends
│   ├── styles/
│   └── utils/
├── astro.config.mjs
├── package.json
└── tsconfig.json
```

## 关键说明
- `src/stores/auth.ts`: 登录态、登录/登出、localStorage 同步。
- `src/components/ui/LiquidGlassCard.vue`: 项目统一玻璃风容器。
- `src/components/admin/AdminRouteGuard.astro`: 后台页面前置鉴权脚本。
- `src/layouts/AdminLayout.astro`: 后台统一壳层。

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
