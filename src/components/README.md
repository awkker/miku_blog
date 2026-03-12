# Components Guide

组件已按业务域和复用层拆分，便于定位与维护。

## `base/`
- 页面基础骨架与 SEO 元信息组件。
- `BaseHead.astro`: 全站 head 元信息与全局样式注入。
- `FormattedDate.astro`: 日期格式化展示。

## `ui/`
- 与业务无关的可复用基础 UI。
- 包含玻璃容器、按钮、输入框、空态/错态、Toast、Skeleton 等。

## `home/`
- 首页专用交互组件。
- `SystemClock.vue`: 顶栏实时钟表。
- `TypewriterSubtitle.vue`: 打字机副标题。

## `blog/`
- 博客阅读与总览页相关组件。
- 顶部导航、目录、点赞面板、代码增强等。

## `auth/`
- 登录域组件。
- `LoginForm.vue`: 登录表单、校验、反馈与跳转逻辑。

## `admin/`
- 后台壳层与管理模块组件。
- 侧栏、顶栏、仪表盘、路由守卫、模块占位组件。

## `guestbook/`
- 留言板输入与卡片流组件。

## `friends/`
- 友链网格页与友链卡片组件。
