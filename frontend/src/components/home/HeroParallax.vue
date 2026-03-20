<template>
  <div
    ref="containerRef"
    class="absolute inset-0 overflow-hidden"
    style="perspective: 1200px"
  >
    <div
      ref="layerRef"
      class="parallax-layer absolute"
      :style="{ transform: layerTransform }"
    >
      <img
        v-for="(src, i) in heroImages"
        :key="src"
        :src="src"
        :alt="`封面 ${i + 1}`"
        class="parallax-img absolute inset-0 h-full w-full object-cover"
        :class="mounted && i === currentIndex ? 'is-active' : ''"
        :loading="i === currentIndex ? 'eager' : 'lazy'"
        :fetchpriority="i === currentIndex ? 'high' : 'low'"
        decoding="async"
        draggable="false"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useStore } from '@nanostores/vue'
import { onMounted, onUnmounted, ref, watch } from 'vue'

import { heroImages, heroIndex, shuffleHeroImage } from '../../stores/heroImage'

const $heroIndex = useStore(heroIndex)
const currentIndex = ref($heroIndex.value)
const mounted = ref(false)

watch($heroIndex, (v) => {
  currentIndex.value = v
})

const containerRef = ref<HTMLElement | null>(null)
const layerRef = ref<HTMLElement | null>(null)

const RANGE = 10
const SHIFT = 18

let rx = 0
let ry = 0
let tx = 0
let ty = 0
let targetRx = 0
let targetRy = 0
let targetTx = 0
let targetTy = 0
let rafId = 0

const layerTransform = ref('rotateX(0deg) rotateY(0deg) translateX(0px) translateY(0px)')

function onMouseMove(e: MouseEvent) {
  const w = window.innerWidth
  const h = window.innerHeight
  const xRatio = (e.clientX / w - 0.5) * 2
  const yRatio = (e.clientY / h - 0.5) * 2

  targetRy = xRatio * RANGE
  targetRx = -yRatio * RANGE
  targetTx = -xRatio * SHIFT
  targetTy = -yRatio * SHIFT
}

function onMouseLeave() {
  targetRx = 0
  targetRy = 0
  targetTx = 0
  targetTy = 0
}

function lerp(a: number, b: number, t: number) {
  return a + (b - a) * t
}

function loop() {
  rx = lerp(rx, targetRx, 0.06)
  ry = lerp(ry, targetRy, 0.06)
  tx = lerp(tx, targetTx, 0.06)
  ty = lerp(ty, targetTy, 0.06)

  layerTransform.value =
    `rotateX(${rx.toFixed(2)}deg) rotateY(${ry.toFixed(2)}deg) translateX(${tx.toFixed(1)}px) translateY(${ty.toFixed(1)}px)`

  rafId = requestAnimationFrame(loop)
}

onMounted(() => {
  mounted.value = true
  shuffleHeroImage()
  currentIndex.value = heroIndex.get()

  const el = containerRef.value
  if (el) {
    el.addEventListener('mousemove', onMouseMove, { passive: true })
    el.addEventListener('mouseleave', onMouseLeave)
  }
  rafId = requestAnimationFrame(loop)
})

onUnmounted(() => {
  const el = containerRef.value
  if (el) {
    el.removeEventListener('mousemove', onMouseMove)
    el.removeEventListener('mouseleave', onMouseLeave)
  }
  if (rafId) cancelAnimationFrame(rafId)
})
</script>

<style scoped>
.parallax-layer {
  inset: -30px;
  transform-origin: 50% 50%;
  transform-style: preserve-3d;
  will-change: transform;
}

.parallax-img {
  opacity: 0;
  transition: opacity 1.2s cubic-bezier(0.4, 0, 0.2, 1);
  will-change: opacity;
}

.parallax-img.is-active {
  opacity: 1;
}
</style>
