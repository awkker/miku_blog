# Liquid Glass 实现总结

## 来源说明
- 效果来源仓库：<https://github.com/shuding/liquid-glass>
- 原始实现主要为 JavaScript
- 本项目在参考原思路基础上，改造为 TypeScript，并结合 Vue 组件生命周期做了工程化封装

## 目标
在页面中提供可复用的“液态玻璃”效果，并用于三处 UI：
- 顶部导航栏
- 中间文案框
- 底部 Dock

当前效果以“边缘折射”为主，默认关闭鼠标交互（`interactive: false` 时不会出现鼠标触发的额外扰动）。

## 文件结构
- `src/utils/liquidGlass.ts`
  - 底层实现，导出 `createLiquidGlass(target, options)`
  - 返回控制器：`update()` / `destroy()`
- `src/components/LiquidGlassFrame.vue`
  - Vue 组件封装，接收 props，并在生命周期里自动创建/销毁液态玻璃效果
- `src/components/AppTopNav.vue`
  - 顶部导航液态玻璃接入示例
- `src/components/TypewriterQuoteGlass.vue`
  - 中部文案液态玻璃接入示例
- `src/components/AppDock.vue`
  - 底部 Dock 液态玻璃接入示例
- `src/App.vue`
  - 页面组装层（组合导航、文案卡片、Dock）

## 核心原理（简述）
1. 使用 `canvas` 生成位移贴图（RG 通道存储 X/Y 位移）。
2. 通过 `SVG feDisplacementMap` 对 `SourceGraphic` 做折射位移。
3. 用 SDF（圆角矩形距离场）控制边缘区域位移强度，实现“边缘更明显”的折射。
4. 叠加 `backdrop-filter`（blur/contrast/brightness/saturate）增强玻璃质感。

## 接口说明
`createLiquidGlass(target, options)` 的 `options`：
- `borderRadius: number` 圆角半径
- `cornerSoftness: number` 边缘过渡宽度
- `displacementStrength: number` 总体折射强度
- `edgeRefractionStrength: number` 边缘额外增强强度
- `blur: number` 背景模糊
- `contrast: number` 对比度
- `brightness: number` 亮度
- `saturate: number` 饱和度
- `interactive: boolean` 是否启用鼠标交互扰动

## 组件调用示例
```vue
<LiquidGlassFrame
  padding="12px 40px"
  :border-radius="50"
  :displacement-strength="1.08"
  :edge-refraction-strength="1.25"
  :interactive="false"
>
  <nav class="pill-nav">...</nav>
</LiquidGlassFrame>
```

## 调参建议
- 想让“边缘折射更明显”：
  - 提高 `edge-refraction-strength`（例如 `1.2 -> 1.5`）
  - 略提高 `displacement-strength`（例如 `1.05 -> 1.2`）
  - 适当降低 `corner-softness`（边缘会更利落）
- 想保持静态边缘折射：
  - 设 `interactive=false`

## 兼容性与降级
- 在测试环境（如 jsdom）没有 `canvas.getContext` 时，会自动降级为普通磨砂滤镜，不抛异常。
- 在无 `ResizeObserver` 环境下会回退到 `window.resize` 监听。
