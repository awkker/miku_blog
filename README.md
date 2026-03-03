# miku_blog

一个基于 Vue 3 + TypeScript + Vite 的个人站点项目，当前包含：
- 首页：液态玻璃（Liquid Glass）组件化效果（导航、文案框、Dock）
- 博客列表页：卡片化列表、侧栏信息、Dock 显隐交互
- 博客阅读页：Markdown 渲染、目录导航、阅读进度与字号调节
- 背景动效：左上向右下飘落的五彩花瓣动画

## 致谢
- 液态玻璃效果思路与实现参考自：<https://github.com/shuding/liquid-glass>
- 原仓库实现以 JavaScript 为主，本项目已按当前工程规范改造为 TypeScript 版本并做了组件化封装

## 技术栈
- Vue 3 (`<script setup>`)
- TypeScript（`vue-tsc`）
- Vite
- Vitest + Vue Test Utils
- Playwright

## 目录结构（核心）
```text
src/
  components/
    AppTopNav.vue
    TypewriterQuoteGlass.vue
    AppDock.vue
    FallingPetals.vue
    LiquidGlassFrame.vue
    blog/
      BlogListHero.vue
      BlogTagNav.vue
      BlogPostCard.vue
      BlogSidebar.vue
      BlogDockToggle.vue
      BlogPostHero.vue
      BlogArticleContent.vue
      BlogPostToc.vue
  composables/
    useTypewriter.ts
  content/
    blogPosts.ts
    posts/
      *.md
  utils/
    liquidGlass.ts
    markdown.ts
  styles/
    base.css
  views/
    HomeView.vue
    BlogListView.vue
    BlogPostView.vue
  App.vue
  main.ts
```

## 快速开始
```sh
npm install
npm run dev
```

## 常用命令
```sh
npm run dev         # 本地开发
npm run type-check  # TS 类型检查
npm run test:unit   # 单元测试
npm run test:e2e    # E2E 测试
npm run build       # 构建
npm run lint        # 代码检查并修复
```

## 参数调整入口
- 花瓣数量/速度/颜色：`src/components/FallingPetals.vue`
- 液态玻璃折射参数：`src/components/LiquidGlassFrame.vue` + `src/utils/liquidGlass.ts`
- 打字机速度：`src/composables/useTypewriter.ts`
- 博客文章元数据（自动读取 Frontmatter）：`src/content/posts/*.md`
- Markdown 渲染逻辑：`src/utils/markdown.ts`

## 博客数据说明
- 文章列表会自动扫描 `src/content/posts/*.md`，无需在 `blogPosts.ts` 手动维护数组。
- 每篇文章请在 Markdown 顶部添加 frontmatter：

```md
---
title: 文章标题
date: 2026-03-03
category: 技术随笔
summary: 列表摘要
cover: /photo/封面.avif
views: 1,000
---
```

- 阅读数会在进入文章详情页时自动累加，统计结果保存到浏览器 `localStorage`（key：`miku_blog_post_views`）。

## 文档
- 液态玻璃实现说明：`docs/liquid-glass/README.md`
- 代码书写规范：`docs/CODING_STYLE.md`
- 博客系统重构总结：`docs/BLOG_REFACTOR_SUMMARY.md`
