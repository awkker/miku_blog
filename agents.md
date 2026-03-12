# Agent Guide (miku-blog-ai)

生成或修改代码前，必须严格遵守以下架构规范与约束。

## 🛠 技术栈硬约束

### 前端生态 (Frontend)
- **包管理器**：只用 `npm` 
  - **允许**：`npm install` / `npm dev` / `npm build`
  - **禁止**：使用 `pnpm` 或 `yarn`，坚决避免 lockfile 冲突。
- **核心框架**：Astro 负责路由与静态页面生成 (SSG)，Vue.js 3 仅作为互动岛屿 (Islands) 按需加载。
- **状态管理**：只用 Nano Stores。
  - **禁止**：引入 Pinia 或 Vuex。所有跨 Astro/Vue 的状态（如壁纸设定、暗黑模式）必须通过 Nano Stores 共享。
- **样式方案**：Tailwind CSS。
  - **约束**：所有样式必须原子化，避免编写独立 `.css` 文件（全局基础变量除外）。

### 后端生态 (Backend)
- **核心语言**：Go (1.21+)
- **Web 框架**：Hertz (字节跳动开源)。
  - **禁止**：引入 Gin、Fiber 等其他 HTTP 框架。
- **数据库交互**：只用 `sqlc`。
  - **约束**：通过编写原生 SQL 生成类型安全的 Go 代码。
  - **禁止**：使用 GORM 或 ent 等重型 ORM（为了追求极致性能与零反射损耗）。
- **数据库**：PostgreSQL (利用原生 GIN 索引进行全文检索，**禁止**引入 Elasticsearch)。
- **缓存/中间件**：Redis (用于 Lua 限流、ZSET 热榜、Bitmap 点赞状态)。

## 🎨 视觉与 UI 约束 (核心特色)

### 主题色与字体
- **全局变量**：必须在 CSS 根节点定义并使用 `--miku-color: rgb(102, 205, 170)` 作为核心主题色（如按钮、高亮、链接等）。
- **字体**：排版优先使用「霞鹜文楷」(LXGW WenKai Screen)，做好字体的 Fallback 设置。

### 液态玻璃风格 (Glassmorphism)
- **设计规范**：UI 组件（如导航栏、评论区卡片、动态壁纸控制面板）必须采用“液态玻璃”拟物风格。
- **Tailwind 实现要求**：
  - 使用：frontend/src/components/LiquidGlassCard.vue和frontend/src/utils/liquidGlass.ts
  - **禁止**：使用扁平化 (Flat Design) 的纯色不透明背景卡片。

## ⚙️ 业务与目录结构约束

### 目录约定
```text
miku-blog/
├── frontend/           # Astro + Vue.js 3
│   ├── src/
│   │   ├── content/    # Content Collections (Markdown/MDX)
│   │   ├── components/ # 纯 Astro 组件 & Vue 互动岛屿组件
│   │   ├── store/      # Nano Stores 状态
│   │   └── pages/      # 基于文件的路由
├── backend/            # Go + Hertz
│   ├── biz/            # 业务逻辑与 Hertz 路由
│   ├── query/          # sqlc 生成的代码存放处
│   └── sql/            # 原生 SQL 语句 (用于 sqlc 生成)
└── docker-compose.yml  # 本地 DB/Redis 一键拉起环境

## 开发与集成流
- **前端交互**：评论系统、动态壁纸切换面板必须封装为 Vue 3 组件，并在 Astro 中使用 client:load 或 client:visible 指令按需激活。
- **环境适配**：后端本地开发需兼容 macOS (Apple Silicon 架构) 的工具链 (如 Homebrew 安装的依赖)。
- **API 通信**：前端请求后端高频 API (如点赞、评论) 必须处理防抖，后端必须有 Redis Lua 脚本限流。

## 行为准则 (Agent 必读)
- **先思考后写代码**：遇到复杂交互或全栈打通逻辑，先输出实现思路和步骤，等待用户确认后再生成代码。
- **只做必要修改**：不要随意重构与当前需求无关的历史代码。

## 安全与文档
- **安全第一**：代码中严禁硬编码数据库密码、JWT Secret 等敏感信息，必须通过环境变量 (.env) 读取。
- **遇事不决查文档**：如果遇到 Hertz 或 Astro 的最新 API 变动，请优先使用 Web Search 检索最新官方文档，不要凭历史语料猜测。