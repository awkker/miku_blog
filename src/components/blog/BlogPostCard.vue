<template>
  <LiquidGlassFrame
    class="post-card"
    :class="{ 'is-cover-expanded': expanded }"
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
      </div>
    </article>
  </LiquidGlassFrame>
</template>

<script setup lang="ts">
import LiquidGlassFrame from '@/components/LiquidGlassFrame.vue'
import type { BlogPostMeta } from '@/content/blogPosts'

interface Props {
  post: BlogPostMeta
  expanded: boolean
}

defineProps<Props>()

const emit = defineEmits<{
  open: [slug: string]
  expand: [slug: string]
  reset: [slug: string]
}>()
</script>

<style scoped>
.post-card {
  position: relative;
  height: auto;
  min-height: 320px;
  align-self: start;
  cursor: pointer;
  transition:
    transform 560ms cubic-bezier(0.2, 0.7, 0.2, 1),
    box-shadow 560ms ease;
  outline: none;
}

.post-card :deep(.liquid-glass-content) {
  height: 100%;
}

.post-card:focus-visible {
  transform: translateY(-6px);
  box-shadow:
    0 0 0 2px rgba(var(--miku-color-rgb), 0.46),
    0 12px 24px rgba(8, 35, 28, 0.18);
}

.post-card:hover {
  transform: translateY(-6px);
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
  height: 56%;
  overflow: hidden;
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
  padding: 14px 16px 16px;
  background: linear-gradient(170deg, rgba(249, 255, 253, 0.78) 0%, rgba(237, 251, 247, 0.62) 100%);
  border-top: 1px solid rgba(255, 255, 255, 0.35);
  backdrop-filter: blur(3px);
  -webkit-backdrop-filter: blur(3px);
  transition: background 720ms ease;
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
</style>

