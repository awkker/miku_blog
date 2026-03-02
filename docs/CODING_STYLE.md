# 代码书写规范（Vue + TypeScript）

适用范围：本仓库 `src/` 下的组件、工具函数、样式与测试代码。

## 1. 基本原则
- 可读性优先于炫技写法。
- 单一职责：一个文件尽量只做一件事。
- 优先复用：重复逻辑抽到 `composables/` 或 `utils/`。
- 注释只解释“为什么”，不解释显而易见的“做了什么”。

## 2. 文件与目录
- 页面组装写在 `App.vue` 或页面级组件。
- 业务组件放在 `src/components/`。
- 可复用逻辑放在 `src/composables/`（命名以 `use` 开头）。
- 纯函数或底层能力放在 `src/utils/`。
- 全局样式放在 `src/styles/`，组件样式放在各自 `.vue` 文件内。

## 3. Vue 组件规范
- 使用 `<script setup lang="ts">`。
- 模板中复杂表达式提取到 `computed` 或函数。
- `props` 必须声明类型，必要时用 `withDefaults` 指定默认值。
- 组件必须在卸载时清理副作用（事件监听、定时器、`requestAnimationFrame`）。
- 组件职责清晰：页面组装组件不要塞底层算法。

## 4. TypeScript 规范
- 避免 `any`；优先显式接口与联合类型。
- 对可空值做显式判断（`null` / `undefined`）。
- 数组索引访问要考虑 `noUncheckedIndexedAccess`。
- 导出的 API（函数、接口）命名清晰，参数语义明确。

## 5. 样式规范
- 优先使用 CSS 变量统一主题色（如 `--miku-color`）。
- 组件样式默认 `scoped`；全局 reset 放 `src/styles/base.css`。
- 动画必须可控，避免过高频率导致性能问题。
- 图层管理统一使用 `z-index` 语义分层，不随意写极大数值。

## 6. 动效与 Canvas 规范
- 所有动画循环必须可停止（保存并 `cancelAnimationFrame`）。
- 监听 `resize` 的逻辑必须在卸载时移除。
- Canvas 需考虑 DPR（设备像素比）防止模糊。
- 动效参数（数量、速度、颜色）集中在配置区，便于调参。

## 7. 测试与验证
- 修改组件行为后，至少运行：
  - `npm run type-check`
  - `npm run test:unit -- --run`
- 测试断言应关注可见行为，不绑定内部实现细节。

## 8. 命名建议
- 组件名：`PascalCase`，如 `TypewriterQuoteGlass.vue`。
- composable：`useXxx`，如 `useTypewriter.ts`。
- 工具函数：动词开头，如 `createLiquidGlass`。
- 变量命名表达语义，避免 `data1`、`tmp` 等弱语义命名。

## 9. 提交前检查清单
- 是否清理了定时器、监听器、动画帧？
- 是否有重复代码可抽离？
- 是否增加了必要注释（尤其复杂动效/算法）？
- 是否通过 `type-check` 和单测？
