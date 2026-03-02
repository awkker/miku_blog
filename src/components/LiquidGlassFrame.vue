<template>
  <div ref="frameRef" class="liquid-glass-frame" :style="frameStyle">
    <div class="liquid-glass-content">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'

import { createLiquidGlass, type LiquidGlassController } from '@/utils/liquidGlass'

interface Props {
  width?: string
  maxWidth?: string
  padding?: string
  borderRadius?: number
  cornerSoftness?: number
  displacementStrength?: number
  edgeRefractionStrength?: number
  blur?: number
  contrast?: number
  brightness?: number
  saturate?: number
  interactive?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  width: '100%',
  maxWidth: '800px',
  padding: '30px 50px',
  borderRadius: 24,
  cornerSoftness: 0.12,
  displacementStrength: 1,
  edgeRefractionStrength: 0.75,
  blur: 0.3,
  contrast: 1.14,
  brightness: 1.04,
  saturate: 1.08,
  interactive: true,
})

const frameRef = ref<HTMLElement | null>(null)
let controller: LiquidGlassController | null = null

const frameStyle = computed(() => ({
  width: props.width,
  maxWidth: props.maxWidth,
  padding: props.padding,
  borderRadius: `${props.borderRadius}px`,
}))

onMounted(() => {
  if (!frameRef.value) {
    return
  }

  // 组件挂载时初始化液态玻璃效果
  controller = createLiquidGlass(frameRef.value, {
    borderRadius: props.borderRadius,
    cornerSoftness: props.cornerSoftness,
    displacementStrength: props.displacementStrength,
    edgeRefractionStrength: props.edgeRefractionStrength,
    blur: props.blur,
    contrast: props.contrast,
    brightness: props.brightness,
    saturate: props.saturate,
    interactive: props.interactive,
  })
})

watch(
  () => ({
    borderRadius: props.borderRadius,
    cornerSoftness: props.cornerSoftness,
    displacementStrength: props.displacementStrength,
    edgeRefractionStrength: props.edgeRefractionStrength,
    blur: props.blur,
    contrast: props.contrast,
    brightness: props.brightness,
    saturate: props.saturate,
    interactive: props.interactive,
  }),
  // 外部 props 变化时同步更新液态玻璃参数
  (next) => controller?.update(next),
)

onBeforeUnmount(() => {
  // 组件卸载时释放监听与滤镜资源
  controller?.destroy()
  controller = null
})
</script>

<style scoped>
.liquid-glass-frame {
  position: relative;
  box-sizing: border-box;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.62);
  background: linear-gradient(140deg, rgba(255, 255, 255, 0.45) 0%, rgba(255, 255, 255, 0.12) 100%);
  box-shadow:
    0 10px 36px rgba(10, 18, 34, 0.12),
    inset 0 1px 1px rgba(255, 255, 255, 0.75),
    inset 0 -10px 26px rgba(255, 255, 255, 0.1);
}

.liquid-glass-frame::before {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(120deg, rgba(255, 255, 255, 0.56) 5%, rgba(255, 255, 255, 0) 40%);
  pointer-events: none;
}

.liquid-glass-content {
  position: relative;
  z-index: 1;
}
</style>
