<template>
  <div ref="pageRef" class="post-page" :class="{ 'dock-open': isDockOpen }">
    <div class="page-bg"></div>
    <div class="page-mask"></div>

    <div class="reading-progress" :style="{ width: `${readingProgress}%` }"></div>

    <BlogPostHero v-if="post" :post="post" :reading-minutes="readingMinutes" />

    <section v-if="post" class="reader-shell">
      <aside class="reader-left">
        <BlogPostMetaCard
          :date="post.date"
          :category="post.category"
          :reading-minutes="readingMinutes"
          :views="post.views"
          :is-liked="isLiked"
          :like-count="likeCount"
          :share-count="shareCount"
          :share-hint="shareHint"
          @back="backToList"
          @toggle-like="toggleLike"
          @share="sharePost"
        />
      </aside>

      <main class="reader-main">
        <BlogArticleContent
          :html="rendered.html"
          :font-scale="fontScale"
          @back="backToList"
          @increase-font-size="increaseFontSize"
          @decrease-font-size="decreaseFontSize"
          @reset-font-size="resetFontSize"
        />

      </main>

      <aside class="reader-side">
        <BlogPostToc
          :headings="rendered.headings"
          :active-heading-id="activeHeadingId"
          :reading-progress="readingProgress"
          @scroll-to-heading="scrollToHeading"
          @scroll-to-top="scrollToTop"
        />
      </aside>
    </section>

    <section v-else class="reader-shell is-empty">
      <main class="reader-main">
        <BlogPostNotFound @back="backToList" />
      </main>
    </section>

    <BlogDockToggle :is-open="isDockOpen" @show="showDock" @hide="hideDock" />
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import BlogArticleContent from '@/components/blog/BlogArticleContent.vue'
import BlogDockToggle from '@/components/blog/BlogDockToggle.vue'
import BlogPostHero from '@/components/blog/BlogPostHero.vue'
import BlogPostMetaCard from '@/components/blog/BlogPostMetaCard.vue'
import BlogPostNotFound from '@/components/blog/BlogPostNotFound.vue'
import BlogPostToc from '@/components/blog/BlogPostToc.vue'
import { getBlogPostBySlug, incrementBlogPostViews } from '@/content/blogPosts'
import { estimateReadingMinutes, renderMarkdown } from '@/utils/markdown'

const LIKE_STATE_STORAGE_KEY = 'miku_blog_post_likes'
const LIKE_COUNT_STORAGE_KEY = 'miku_blog_post_like_count'
const SHARE_COUNT_STORAGE_KEY = 'miku_blog_post_share_count'
const route = useRoute()
const router = useRouter()

const pageRef = ref<HTMLElement | null>(null)
const readingProgress = ref(0)
const fontScale = ref(1)
const isDockOpen = ref(false)
const isLiked = ref(false)
const likeCount = ref(0)
const shareCount = ref(0)
const activeHeadingId = ref('')
const shareHint = ref('')
let shareHintTimer: number | null = null
let headingObserver: IntersectionObserver | null = null
let headingVisibility = new Map<string, number>()

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
  const container = pageRef.value
  if (!container) {
    return
  }
  container.scrollTo({ top: 0, behavior: 'smooth' })
  window.setTimeout(() => {
    container.scrollTo({ top: 0, behavior: 'auto' })
  }, 320)
  window.setTimeout(() => {
    container.scrollTop = 0
  }, 520)
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
  activeHeadingId.value = id
  const containerRect = container.getBoundingClientRect()
  const targetRect = target.getBoundingClientRect()
  const currentTop = container.scrollTop
  const safeTopOffset = 104
  const destination = Math.max(0, currentTop + targetRect.top - containerRect.top - safeTopOffset)

  container.scrollTo({
    top: destination,
    behavior: 'smooth',
  })
  window.setTimeout(() => {
    container.scrollTo({
      top: destination,
      behavior: 'auto',
    })
  }, 420)
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

function showDock(): void {
  isDockOpen.value = true
}

function hideDock(): void {
  isDockOpen.value = false
}

function recordPostView(): void {
  if (!post.value) {
    return
  }
  incrementBlogPostViews(post.value.slug)
}

function toggleLike(): void {
  const currentSlug = post.value?.slug
  if (!currentSlug) {
    return
  }
  const stateMap = readLikeMap()
  const countMap = readLikeCountMap()
  const alreadyLiked = Boolean(stateMap[currentSlug])
  const currentCount = getStoredCount(countMap[currentSlug])

  if (alreadyLiked) {
    stateMap[currentSlug] = false
    countMap[currentSlug] = Math.max(0, currentCount - 1)
    isLiked.value = false
    likeCount.value = countMap[currentSlug]
  } else {
    stateMap[currentSlug] = true
    countMap[currentSlug] = currentCount + 1
    isLiked.value = true
    likeCount.value = countMap[currentSlug]
  }

  writeLikeMap(stateMap)
  writeLikeCountMap(countMap)
}

async function sharePost(): Promise<void> {
  if (!post.value || typeof window === 'undefined') {
    return
  }
  const shareUrl = window.location.href
  if (typeof navigator !== 'undefined' && typeof navigator.share === 'function') {
    try {
      await navigator.share({
        title: post.value.title,
        url: shareUrl,
      })
      incrementShareCount(post.value.slug)
      setShareHint('已打开系统分享面板')
      return
    } catch {
      // User canceled share panel; fallback to copy link.
    }
  }

  if (typeof navigator !== 'undefined' && navigator.clipboard?.writeText) {
    try {
      await navigator.clipboard.writeText(shareUrl)
      incrementShareCount(post.value.slug)
      setShareHint('链接已复制到剪贴板')
      return
    } catch {
      setShareHint('复制失败，请手动复制地址栏链接')
      return
    }
  }

  setShareHint('当前环境不支持自动分享，请手动复制链接')
}

function setShareHint(message: string): void {
  shareHint.value = message
  if (shareHintTimer !== null) {
    window.clearTimeout(shareHintTimer)
  }
  shareHintTimer = window.setTimeout(() => {
    shareHint.value = ''
    shareHintTimer = null
  }, 2000)
}

function syncInteractionState(): void {
  const currentSlug = post.value?.slug
  if (!currentSlug) {
    isLiked.value = false
    likeCount.value = 0
    shareCount.value = 0
    return
  }
  const likeStateMap = readLikeMap()
  const likeCountMap = readLikeCountMap()
  const shareCountMap = readShareCountMap()

  isLiked.value = Boolean(likeStateMap[currentSlug])
  likeCount.value = getStoredCount(likeCountMap[currentSlug])
  shareCount.value = getStoredCount(shareCountMap[currentSlug])
}

function readLikeMap(): Record<string, boolean> {
  const storage = getLocalStorage()
  if (!storage) {
    return {}
  }
  const raw = storage.getItem(LIKE_STATE_STORAGE_KEY)
  if (!raw) {
    return {}
  }
  try {
    const parsed = JSON.parse(raw)
    if (!parsed || typeof parsed !== 'object') {
      return {}
    }
    return parsed as Record<string, boolean>
  } catch {
    return {}
  }
}

function writeLikeMap(map: Record<string, boolean>): void {
  const storage = getLocalStorage()
  if (!storage) {
    return
  }
  storage.setItem(LIKE_STATE_STORAGE_KEY, JSON.stringify(map))
}

function incrementShareCount(currentSlug: string): void {
  const map = readShareCountMap()
  const current = getStoredCount(map[currentSlug])
  map[currentSlug] = current + 1
  shareCount.value = map[currentSlug]
  writeShareCountMap(map)
}

function readLikeCountMap(): Record<string, number> {
  const storage = getLocalStorage()
  if (!storage) {
    return {}
  }
  const raw = storage.getItem(LIKE_COUNT_STORAGE_KEY)
  if (!raw) {
    return {}
  }
  try {
    const parsed = JSON.parse(raw)
    if (!parsed || typeof parsed !== 'object') {
      return {}
    }
    return parsed as Record<string, number>
  } catch {
    return {}
  }
}

function writeLikeCountMap(map: Record<string, number>): void {
  const storage = getLocalStorage()
  if (!storage) {
    return
  }
  storage.setItem(LIKE_COUNT_STORAGE_KEY, JSON.stringify(map))
}

function readShareCountMap(): Record<string, number> {
  const storage = getLocalStorage()
  if (!storage) {
    return {}
  }
  const raw = storage.getItem(SHARE_COUNT_STORAGE_KEY)
  if (!raw) {
    return {}
  }
  try {
    const parsed = JSON.parse(raw)
    if (!parsed || typeof parsed !== 'object') {
      return {}
    }
    return parsed as Record<string, number>
  } catch {
    return {}
  }
}

function writeShareCountMap(map: Record<string, number>): void {
  const storage = getLocalStorage()
  if (!storage) {
    return
  }
  storage.setItem(SHARE_COUNT_STORAGE_KEY, JSON.stringify(map))
}

function getStoredCount(raw: unknown): number {
  if (typeof raw !== 'number' || !Number.isFinite(raw) || raw < 0) {
    return 0
  }
  return Math.floor(raw)
}

function disconnectHeadingObserver(): void {
  headingObserver?.disconnect()
  headingObserver = null
  headingVisibility = new Map<string, number>()
}

function updateActiveHeadingByPosition(headingElements: HTMLElement[], container: HTMLElement): void {
  if (!headingElements.length) {
    activeHeadingId.value = ''
    return
  }
  const top = container.getBoundingClientRect().top
  const anchorOffset = 136
  let currentId = headingElements[0]?.id ?? ''
  for (const heading of headingElements) {
    if (heading.getBoundingClientRect().top - top <= anchorOffset) {
      currentId = heading.id
    } else {
      break
    }
  }
  activeHeadingId.value = currentId
}

function setupHeadingObserver(): void {
  disconnectHeadingObserver()

  const container = pageRef.value
  if (!container) {
    return
  }

  const headingElements = rendered.value.headings
    .map((heading) => document.getElementById(heading.id))
    .filter((item): item is HTMLElement => Boolean(item))

  if (!headingElements.length) {
    activeHeadingId.value = ''
    return
  }

  headingObserver = new IntersectionObserver(
    (entries) => {
      for (const entry of entries) {
        const id = (entry.target as HTMLElement).id
        if (entry.isIntersecting) {
          headingVisibility.set(id, entry.intersectionRatio)
        } else {
          headingVisibility.delete(id)
        }
      }

      if (headingVisibility.size > 0) {
        let bestId = ''
        let bestScore = -1
        for (let index = 0; index < headingElements.length; index += 1) {
          const heading = headingElements[index]
          if (!heading) {
            continue
          }
          const ratio = headingVisibility.get(heading.id)
          if (ratio === undefined) {
            continue
          }
          const score = ratio * 10 + (1 - index / 1000)
          if (score > bestScore) {
            bestScore = score
            bestId = heading.id
          }
        }
        if (bestId) {
          activeHeadingId.value = bestId
          return
        }
      }

      updateActiveHeadingByPosition(headingElements, container)
    },
    {
      root: container,
      rootMargin: '-110px 0px -58% 0px',
      threshold: [0, 0.01, 0.1, 0.3, 0.6, 1],
    },
  )

  for (const heading of headingElements) {
    headingObserver.observe(heading)
  }

  updateActiveHeadingByPosition(headingElements, container)
}

function getLocalStorage(): Storage | null {
  if (typeof window === 'undefined') {
    return null
  }
  try {
    return window.localStorage
  } catch {
    return null
  }
}

onMounted(() => {
  recordPostView()
  syncInteractionState()
  const container = pageRef.value
  container?.addEventListener('scroll', updateReadingProgress, { passive: true })
  window.addEventListener('resize', updateReadingProgress)
  nextTick(() => {
    updateReadingProgress()
    setupHeadingObserver()
  })
})

onBeforeUnmount(() => {
  const container = pageRef.value
  container?.removeEventListener('scroll', updateReadingProgress)
  window.removeEventListener('resize', updateReadingProgress)
  disconnectHeadingObserver()
  if (shareHintTimer !== null) {
    window.clearTimeout(shareHintTimer)
  }
})

watch(
  () => route.params.slug,
  () => {
    fontScale.value = 1
    recordPostView()
    syncInteractionState()
    shareHint.value = ''
    pageRef.value?.scrollTo({ top: 0, behavior: 'auto' })
    nextTick(() => {
      updateReadingProgress()
      setupHeadingObserver()
    })
  },
)

watch(
  () => rendered.value.html,
  () => {
    nextTick(() => setupHeadingObserver())
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
  padding-bottom: 68px;
}

.post-page.dock-open {
  padding-bottom: 120px;
}

.page-bg {
  position: fixed;
  inset: -28px;
  z-index: 0;
  pointer-events: none;
  background-image: url('/photo/ui.png');
  background-size: cover;
  background-position: center;
  filter: blur(18px) saturate(1.08);
  transform: scale(1.08);
}

.page-mask {
  position: fixed;
  inset: 0;
  z-index: 0;
  pointer-events: none;
  background:
    radial-gradient(circle at 18% 16%, rgba(236, 255, 250, 0.42) 0%, rgba(236, 255, 250, 0) 54%),
    linear-gradient(180deg, rgba(244, 251, 248, 0.58) 0%, rgba(247, 252, 250, 0.76) 100%);
}

.post-page > :not(.page-bg):not(.page-mask) {
  position: relative;
  z-index: 1;
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
  width: min(1320px, calc(100% - 42px));
  margin: 20px auto 0;
  display: grid;
  justify-content: center;
  grid-template-columns: 220px minmax(0, 760px) 248px;
  gap: 20px;
}

.reader-left,
.reader-main,
.reader-side {
  min-width: 0;
}

.reader-left,
.reader-side {
  position: sticky;
  top: 18px;
  align-self: start;
}

.reader-main {
  width: min(100%, 760px);
  justify-self: center;
}

.reader-shell.is-empty {
  width: min(860px, calc(100% - 28px));
  grid-template-columns: 1fr;
}

@media (max-width: 1320px) {
  .reader-shell {
    grid-template-columns: 200px minmax(0, 740px) 228px;
  }
}

@media (max-width: 1180px) {
  .reader-shell {
    width: min(840px, calc(100% - 24px));
    grid-template-columns: 1fr;
    gap: 14px;
  }

  .reader-left,
  .reader-side {
    position: static;
    top: auto;
  }

  .reader-main {
    width: 100%;
  }

}

@media (max-width: 640px) {
  .reader-shell {
    width: calc(100% - 14px);
  }
}
</style>
