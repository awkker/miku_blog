<template>
  <div
    ref="frameRef"
    class="glass-layer relative overflow-hidden rounded-3xl"
  >
    <div class="pointer-events-none absolute inset-0 bg-gradient-to-br from-white/60 via-white/20 to-transparent" />
    <div class="relative z-[1]">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'

import { createLiquidGlass, type LiquidGlassController } from '../../utils/liquidGlass'

interface Props {
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

onMounted(() => {
  if (!frameRef.value) {
    return
  }

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
  (next) => controller?.update(next),
)

onBeforeUnmount(() => {
  controller?.destroy()
  controller = null
})
</script>
