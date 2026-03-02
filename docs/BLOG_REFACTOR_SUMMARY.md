# 博客系统重构总结（2026-03-02）

## 目标
- 增加博客阅读页，支持 Markdown 渲染。
- 在博客列表页点击卡片可跳转阅读页。
- 保持液态玻璃视觉风格一致。
- 将超长页面文件拆分为可维护的组件结构。

## 结果概览
- 新增路由：`/blog/:slug`
- 新增内容层：`src/content/blogPosts.ts` + `src/content/posts/*.md`
- 新增 Markdown 渲染工具：`src/utils/markdown.ts`
- 博客列表页与阅读页完成组件化拆分

## 页面与路由
- 列表页：`src/views/BlogListView.vue`
- 阅读页：`src/views/BlogPostView.vue`
- 路由定义：`src/router/index.ts`

## 组件拆分结果

### 列表页组件
- `src/components/blog/BlogListHero.vue`
- `src/components/blog/BlogTagNav.vue`
- `src/components/blog/BlogPostCard.vue`
- `src/components/blog/BlogSidebar.vue`
- `src/components/blog/BlogDockToggle.vue`

### 阅读页组件
- `src/components/blog/BlogPostHero.vue`
- `src/components/blog/BlogArticleContent.vue`
- `src/components/blog/BlogPostToc.vue`

## Markdown 渲染说明
- 当前使用项目内置轻量解析器（`src/utils/markdown.ts`），支持：
  - 标题（h1-h6）
  - 段落
  - 有序/无序列表
  - 引用
  - 代码块与行内代码
  - 链接
  - 分割线
- 同时提供：
  - 标题目录提取（TOC）
  - 预计阅读时长估算

## 阅读增强功能
- 阅读进度条
- 字号调节（A- / A / A+）
- 目录跳转
- 回到顶部
- 文章不存在时的兜底提示与返回入口

## 维护收益
- `BlogListView.vue`：`644 -> 152` 行
- `BlogPostView.vue`：`456 -> 239` 行
- 页面层只负责“数据与状态编排”，样式与结构下沉到组件

## 验证
- `npm run type-check` 通过
- `npm run test:unit -- --run` 通过

