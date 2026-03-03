<template>
  <LiquidGlassFrame
    class="post-card"
    :class="{ 'is-cover-expanded': expanded }"
    @pointermove="handlePointerMove"
    @pointerleave="handlePointerLeave"
    width="100%"
    max-width="none"
    padding="0"
    :border-radius="20"
    :displacement-strength="1.08"
    :edge-refraction-strength="1.3"
    :interactive="false"
    role="link"
    tabindex="0"
    @click="emit('open', post.slug)"
    @keydown.enter.prevent="emit('open', post.slug)"
    @keydown.space.prevent="emit('open', post.slug)"
    @mouseleave="emit('reset', post.slug)"
  >
    <article class="post-article">
      <div class="post-cover" @mouseenter="emit('expand', post.slug)">
        <img :src="post.cover" :alt="post.title" loading="lazy" />
      </div>
      <div class="post-body">
        <p class="post-meta">{{ post.date }} · {{ post.category }}</p>
        <h3 class="post-title">{{ post.title }}</h3>
        <p class="post-summary">{{ post.summary }}</p>
        <div class="post-stats">
          <span>♥ {{ compactCount(post.likes) }}</span>
          <span>↗ {{ compactCount(post.shares) }}</span>
          <span>◉ {{ compactCount(viewCount) }}</span>
        </div>
      </div>
    </article>
  </LiquidGlassFrame>
</template>

<script setup lang="ts">
import { computed } from 'vue'

import LiquidGlassFrame from '@/components/LiquidGlassFrame.vue'
import type { BlogPostMeta } from '@/content/blogPosts'

interface Props {
  post: BlogPostMeta
  expanded: boolean
}

const props = defineProps<Props>()

const emit = defineEmits<{
  open: [slug: string]
  expand: [slug: string]
  reset: [slug: string]
}>()

const viewCount = computed(() => Number.parseInt(props.post.views.replace(/,/g, ''), 10) || 0)

function compactCount(value: number): string {
  const absolute = Math.abs(value)
  if (absolute >= 1_000_000_000) {
    return `${trimTrailingZero((value / 1_000_000_000).toFixed(1))}B`
  }
  if (absolute >= 1_000_000) {
    return `${trimTrailingZero((value / 1_000_000).toFixed(1))}M`
  }
  if (absolute >= 1_000) {
    return `${trimTrailingZero((value / 1_000).toFixed(1))}k`
  }
  return String(value)
}

function trimTrailingZero(value: string): string {
  return value.replace(/\.0$/, '')
}

function prefersReducedMotion(): boolean {
  if (typeof window === 'undefined' || typeof window.matchMedia !== 'function') {
    return false
  }
  return window.matchMedia('(prefers-reduced-motion: reduce)').matches
}

function handlePointerMove(event: PointerEvent): void {
  if (prefersReducedMotion()) {
    return
  }
  const currentTarget = event.currentTarget
  if (!(currentTarget instanceof HTMLElement)) {
    return
  }
  const rect = currentTarget.getBoundingClientRect()
  if (rect.width <= 0 || rect.height <= 0) {
    return
  }
  const x = event.clientX - rect.left
  const y = event.clientY - rect.top
  const px = x / rect.width
  const py = y / rect.height
  const tiltX = (0.5 - py) * 6.8
  const tiltY = (px - 0.5) * 8.2
  const depthX = (px - 0.5) * 9
  const depthY = (py - 0.5) * 8

  currentTarget.style.setProperty('--tilt-x', `${tiltX.toFixed(2)}deg`)
  currentTarget.style.setProperty('--tilt-y', `${tiltY.toFixed(2)}deg`)
  currentTarget.style.setProperty('--depth-x', `${depthX.toFixed(2)}px`)
  currentTarget.style.setProperty('--depth-y', `${depthY.toFixed(2)}px`)
}

function handlePointerLeave(event: PointerEvent): void {
  const currentTarget = event.currentTarget
  if (!(currentTarget instanceof HTMLElement)) {
    return
  }
  currentTarget.style.setProperty('--tilt-x', '0deg')
  currentTarget.style.setProperty('--tilt-y', '0deg')
  currentTarget.style.setProperty('--depth-x', '0px')
  currentTarget.style.setProperty('--depth-y', '0px')
}
</script>

<style scoped>
.post-card {
  --tilt-x: 0deg;
  --tilt-y: 0deg;
  --depth-x: 0px;
  --depth-y: 0px;
  position: relative;
  height: auto;
  min-height: 320px;
  align-self: start;
  cursor: pointer;
  transform-style: preserve-3d;
  transition:
    transform 560ms cubic-bezier(0.2, 0.7, 0.2, 1),
    box-shadow 560ms ease;
  transform: perspective(980px) translateY(0) rotateX(var(--tilt-x)) rotateY(var(--tilt-y));
  outline: none;
}

.post-card :deep(.liquid-glass-content) {
  height: 100%;
  transform-style: preserve-3d;
}

.post-card:focus-visible {
  transform: perspective(980px) translateY(-6px) rotateX(var(--tilt-x)) rotateY(var(--tilt-y));
  box-shadow:
    0 0 0 2px rgba(var(--miku-color-rgb), 0.46),
    0 12px 24px rgba(8, 35, 28, 0.18);
}

.post-card:hover {
  transform: perspective(980px) translateY(-6px) rotateX(var(--tilt-x)) rotateY(var(--tilt-y));
  box-shadow:
    0 12px 24px rgba(8, 35, 28, 0.18),
    0 3px 10px rgba(255, 255, 255, 0.2) inset;
}

.post-article {
  position: relative;
  min-height: 320px;
  height: 100%;
  overflow: hidden;
}

.post-cover {
  position: absolute;
  inset: 0 0 auto 0;
  height: 51%;
  overflow: hidden;
  transform: translate3d(calc(var(--depth-x) * 0.18), calc(var(--depth-y) * 0.2), 12px);
  transition: height 880ms cubic-bezier(0.22, 0.62, 0.36, 1);
  z-index: 0;
}

.post-cover img {
  width: 100%;
  height: 100%;
  display: block;
  object-fit: cover;
  transition: transform 880ms cubic-bezier(0.22, 0.62, 0.36, 1);
}

.post-body {
  position: absolute;
  inset: auto 0 0 0;
  z-index: 1;
  transform: translate3d(calc(var(--depth-x) * 0.34), calc(var(--depth-y) * 0.3), 18px);
  padding: 14px 16px 16px;
  background: linear-gradient(170deg, rgba(249, 255, 253, 0.78) 0%, rgba(237, 251, 247, 0.62) 100%);
  border-top: 1px solid rgba(255, 255, 255, 0.35);
  backdrop-filter: blur(3px);
  -webkit-backdrop-filter: blur(3px);
  transition:
    background 720ms ease,
    transform 220ms ease;
}

.post-card.is-cover-expanded .post-cover {
  height: 100%;
}

.post-card.is-cover-expanded .post-cover img {
  transform: scale(1.03);
}

.post-card.is-cover-expanded .post-body {
  background: linear-gradient(180deg, rgba(247, 255, 252, 0.12) 0%, rgba(247, 255, 252, 0.72) 56%, rgba(234, 250, 245, 0.82) 100%);
}

.post-meta {
  margin: 0 0 8px;
  font-size: 0.78rem;
  color: rgba(29, 105, 89, 0.8);
}

.post-title {
  margin: 0 0 8px;
  font-size: 1.02rem;
  color: rgba(8, 48, 39, 0.94);
  line-height: 1.42;
}

.post-summary {
  margin: 0;
  color: rgba(18, 72, 60, 0.86);
  line-height: 1.55;
  font-size: 0.86rem;
  display: -webkit-box;
  overflow: hidden;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
}

.post-stats {
  margin-top: 10px;
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 0.76rem;
  color: rgba(16, 80, 67, 0.78);
}

.post-stats span {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

@media (prefers-reduced-motion: reduce) {
  .post-card,
  .post-cover,
  .post-body {
    transition: none;
    transform: none;
  }
}
</style>
