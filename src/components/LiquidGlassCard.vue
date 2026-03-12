<template>
  <div ref="glassContainer" class="liquid-glass-wrapper">
    <slot />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { createLiquidGlass } from '../utils/liquidGlass'
import type { LiquidGlassController, LiquidGlassOptions } from '../utils/liquidGlass'

// 允许外部传入玻璃的配置参数
const props = defineProps<{
  options?: LiquidGlassOptions
}>()

const glassContainer = ref<HTMLElement | null>(null)
let controller: LiquidGlassController | null = null

// onMounted 只会在客户端浏览器环境执行，完美避开 Astro 的服务端渲染报错
onMounted(() => {
  if (glassContainer.value) {
    // 实例化你的液态玻璃魔法！
    controller = createLiquidGlass(glassContainer.value, props.options)
  }
})

// 组件销毁时，记得清理绑定的事件和创建的 canvas/svg，防止内存泄漏
onBeforeUnmount(() => {
  if (controller) {
    controller.destroy()
  }
})
</script>

<style scoped>
/* 给容器加一点基础样式，确保 backdrop-filter 能正常生效 */
.liquid-glass-wrapper {
  position: relative;
  overflow: hidden; /* 防止内容溢出破坏玻璃边缘 */
  /* 可以配合 Tailwind 加一些透明背景，比如 bg-white/10 */
  background: rgba(255, 255, 255, 0.05); 
}
</style>