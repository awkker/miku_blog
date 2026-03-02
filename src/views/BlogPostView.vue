<template>
  <div ref="pageRef" class="post-page">
    <FallingPetals />

    <div class="reading-progress" :style="{ width: `${readingProgress}%` }"></div>

    <BlogPostHero v-if="post" :post="post" :reading-minutes="readingMinutes" />

    <section class="reader-shell">
      <main class="reader-main">
        <BlogArticleContent
          v-if="post"
          :html="rendered.html"
          :font-scale="fontScale"
          @back="backToList"
          @increase-font-size="increaseFontSize"
          @decrease-font-size="decreaseFontSize"
          @reset-font-size="resetFontSize"
        />

        <LiquidGlassFrame
          v-else
          class="not-found-card"
          width="100%"
          max-width="none"
          padding="28px 22px"
          :border-radius="22"
          :displacement-strength="1.08"
          :edge-refraction-strength="1.3"
          :interactive="false"
        >
          <h2>文章不存在</h2>
          <p>可能已被移动或 slug 不正确。</p>
          <button type="button" class="not-found-btn" @click="backToList">返回博客列表</button>
        </LiquidGlassFrame>
      </main>

      <aside v-if="post" class="reader-side">
        <BlogPostToc :headings="rendered.headings" @scroll-to-heading="scrollToHeading" @scroll-to-top="scrollToTop" />
      </aside>
    </section>

    <div class="post-dock-wrap">
      <AppDock />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import AppDock from '@/components/AppDock.vue'
import FallingPetals from '@/components/FallingPetals.vue'
import LiquidGlassFrame from '@/components/LiquidGlassFrame.vue'
import BlogArticleContent from '@/components/blog/BlogArticleContent.vue'
import BlogPostHero from '@/components/blog/BlogPostHero.vue'
import BlogPostToc from '@/components/blog/BlogPostToc.vue'
import { getBlogPostBySlug } from '@/content/blogPosts'
import { estimateReadingMinutes, renderMarkdown } from '@/utils/markdown'

const route = useRoute()
const router = useRouter()

const pageRef = ref<HTMLElement | null>(null)
const readingProgress = ref(0)
const fontScale = ref(1)

const slug = computed(() => String(route.params.slug ?? ''))
const post = computed(() => getBlogPostBySlug(slug.value))
const rendered = computed(() => renderMarkdown(post.value?.markdown ?? ''))
const readingMinutes = computed(() => estimateReadingMinutes(post.value?.markdown ?? ''))

function updateReadingProgress(): void {
  const container = pageRef.value
  if (!container) {
    return
  }
  const maxScrollTop = container.scrollHeight - container.clientHeight
  if (maxScrollTop <= 0) {
    readingProgress.value = 0
    return
  }
  readingProgress.value = Math.min(100, Math.max(0, (container.scrollTop / maxScrollTop) * 100))
}

function scrollToTop(): void {
  pageRef.value?.scrollTo({ top: 0, behavior: 'smooth' })
}

function backToList(): void {
  void router.push({ name: 'blog' })
}

function scrollToHeading(id: string): void {
  const container = pageRef.value
  const target = document.getElementById(id)
  if (!container || !target) {
    return
  }
  const containerTop = container.getBoundingClientRect().top
  const targetTop = target.getBoundingClientRect().top
  const offsetTop = container.scrollTop + targetTop - containerTop - 96
  container.scrollTo({ top: Math.max(0, offsetTop), behavior: 'smooth' })
}

function increaseFontSize(): void {
  fontScale.value = Math.min(1.3, Number((fontScale.value + 0.06).toFixed(2)))
}

function decreaseFontSize(): void {
  fontScale.value = Math.max(0.9, Number((fontScale.value - 0.06).toFixed(2)))
}

function resetFontSize(): void {
  fontScale.value = 1
}

onMounted(() => {
  const container = pageRef.value
  container?.addEventListener('scroll', updateReadingProgress, { passive: true })
  window.addEventListener('resize', updateReadingProgress)
  nextTick(() => updateReadingProgress())
})

onBeforeUnmount(() => {
  const container = pageRef.value
  container?.removeEventListener('scroll', updateReadingProgress)
  window.removeEventListener('resize', updateReadingProgress)
})

watch(
  () => route.params.slug,
  () => {
    fontScale.value = 1
    pageRef.value?.scrollTo({ top: 0, behavior: 'auto' })
    nextTick(() => updateReadingProgress())
  },
)
</script>

<style scoped>
.post-page {
  position: relative;
  width: 100%;
  height: 100%;
  max-height: 100dvh;
  overflow-x: hidden;
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;
  background-image: url('/photo/ui.png');
  background-size: cover;
  background-position: center;
  padding-bottom: 120px;
}

.reading-progress {
  position: fixed;
  top: 0;
  left: 0;
  height: 3px;
  width: 0;
  z-index: 8;
  background: linear-gradient(90deg, rgba(137, 231, 207, 0.95) 0%, rgba(91, 191, 169, 0.95) 100%);
  box-shadow: 0 0 12px rgba(102, 205, 170, 0.42);
  transition: width 120ms linear;
}

.reader-shell {
  width: min(1080px, calc(100% - 36px));
  margin: 18px auto 0;
  display: grid;
  grid-template-columns: minmax(0, 1fr) 256px;
  gap: 20px;
}

.reader-main,
.reader-side {
  min-width: 0;
}

.reader-side {
  position: sticky;
  top: 18px;
  align-self: start;
}

.not-found-card h2 {
  margin: 0 0 8px;
  color: rgba(10, 67, 55, 0.95);
}

.not-found-card p {
  margin: 0 0 14px;
  color: rgba(20, 87, 72, 0.88);
}

.not-found-btn {
  border: 0;
  border-radius: 10px;
  padding: 6px 10px;
  background: rgba(241, 255, 251, 0.78);
  color: rgba(10, 73, 60, 0.92);
  font-family: var(--font-body);
  cursor: pointer;
  transition: transform 180ms ease, background-color 180ms ease;
}

.not-found-btn:hover {
  transform: translateY(-1px);
  background: rgba(222, 249, 241, 0.92);
}

.post-dock-wrap {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 10px;
  z-index: 6;
  display: flex;
  justify-content: center;
  pointer-events: none;
}

.post-dock-wrap :deep(.mac-dock-glass) {
  pointer-events: auto;
}

@media (max-width: 1040px) {
  .reader-shell {
    grid-template-columns: 1fr;
  }

  .reader-side {
    position: static;
    top: auto;
  }
}
</style>
