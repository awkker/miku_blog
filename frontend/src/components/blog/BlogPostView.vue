<template>
  <div class="mx-auto w-full max-w-[1260px]">
    <div v-if="loading" class="space-y-6 py-12">
      <div class="animate-pulse space-y-4">
        <div class="h-8 w-3/4 rounded bg-slate-200/60" />
        <div class="h-4 w-1/3 rounded bg-slate-100/80" />
        <div class="h-64 rounded-3xl bg-slate-200/40" />
        <div class="space-y-2">
          <div class="h-4 w-full rounded bg-slate-100/60" />
          <div class="h-4 w-5/6 rounded bg-slate-100/60" />
          <div class="h-4 w-4/6 rounded bg-slate-100/60" />
        </div>
      </div>
    </div>

    <div v-else-if="error" class="py-20 text-center">
      <div class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-red-50">
        <svg viewBox="0 0 24 24" class="h-8 w-8 fill-none stroke-red-400 stroke-[1.5]">
          <circle cx="12" cy="12" r="10" />
          <line x1="15" y1="9" x2="9" y2="15" />
          <line x1="9" y1="9" x2="15" y2="15" />
        </svg>
      </div>
      <p class="text-sm text-red-600">{{ error }}</p>
      <a
        href="/blog"
        class="mt-4 inline-flex items-center gap-1 rounded-xl border border-miku/35 bg-miku-soft px-4 py-2 text-sm text-miku transition hover:border-miku/55"
      >
        返回文章列表
      </a>
    </div>

    <article v-else-if="post" class="space-y-8">
      <a
        href="/blog"
        class="inline-flex items-center gap-1.5 rounded-lg border border-slate-200 bg-white/60 px-3 py-1.5 text-xs text-slate-500 transition hover:border-miku/40 hover:text-miku"
      >
        <svg viewBox="0 0 24 24" class="h-3.5 w-3.5 fill-none stroke-current stroke-[2]">
          <path d="M19 12H5M12 19l-7-7 7-7" />
        </svg>
        返回文章列表
      </a>

      <div
        v-if="post.hero_image_url"
        class="overflow-hidden rounded-[30px] border border-white/70 shadow-[0_20px_40px_rgba(15,23,42,0.14)]"
      >
        <img :src="post.hero_image_url" :alt="post.title" class="h-[260px] w-full object-cover md:h-[460px]" loading="lazy" />
      </div>

      <header class="mx-auto w-full max-w-[1100px] space-y-5">
        <div class="flex flex-wrap items-center gap-2">
          <span
            v-if="post.category"
            class="rounded-full border border-miku/35 bg-miku-soft px-3 py-1 text-xs font-semibold tracking-[0.14em] text-miku"
          >
            {{ post.category }}
          </span>
          <span
            v-for="tag in post.tags || []"
            :key="tag.slug || tag.name"
            class="rounded-full border border-slate-200 bg-white px-2.5 py-1 text-[11px] text-slate-500"
          >
            {{ tag.name }}
          </span>
        </div>

        <h1 class="text-3xl font-black leading-tight text-slate-900 md:text-[3rem]">
          {{ post.title }}
        </h1>

        <p v-if="post.excerpt" class="max-w-4xl text-base leading-relaxed text-slate-500 md:text-lg">
          {{ post.excerpt }}
        </p>

        <div class="flex flex-wrap items-center gap-4 text-xs text-slate-400 md:text-sm">
          <span class="inline-flex items-center gap-1">
            <svg viewBox="0 0 24 24" class="h-3.5 w-3.5 fill-none stroke-current stroke-[1.8]">
              <rect x="3" y="4" width="18" height="18" rx="2" ry="2" />
              <line x1="16" y1="2" x2="16" y2="6" />
              <line x1="8" y1="2" x2="8" y2="6" />
              <line x1="3" y1="10" x2="21" y2="10" />
            </svg>
            {{ formatDate(post.published_at || post.created_at) }}
          </span>
          <span class="inline-flex items-center gap-1">
            <svg viewBox="0 0 24 24" class="h-3.5 w-3.5 fill-none stroke-current stroke-[1.6]">
              <path d="M12 6v6l4 2" />
              <circle cx="12" cy="12" r="9" />
            </svg>
            预计 {{ readingMinutes }} 分钟
          </span>
          <span class="inline-flex items-center gap-1">
            <svg viewBox="0 0 24 24" class="h-3.5 w-3.5 fill-none stroke-current stroke-[1.5]">
              <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
              <circle cx="12" cy="12" r="3" />
            </svg>
            {{ post.view_count.toLocaleString() }} 次阅读
          </span>
          <span class="inline-flex items-center gap-1">
            <svg viewBox="0 0 24 24" class="h-3.5 w-3.5 fill-none stroke-current stroke-[1.5]">
              <path d="M20.84 4.61a5.5 5.5 0 00-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 00-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 000-7.78z" />
            </svg>
            {{ post.like_count }} 点赞
          </span>
        </div>
      </header>

      <section class="mx-auto w-full max-w-[1100px]">
        <PostLikeBar
          :post-id="post.id"
          :like-count="post.like_count"
          :liked="post.liked"
          :comment-count="commentCount"
          @like-updated="handleLikeUpdated"
        />
      </section>

      <section class="mx-auto grid w-full max-w-[1100px] gap-6 xl:grid-cols-[minmax(0,1fr)_292px]">
        <article class="min-w-0">
          <div class="rounded-[26px] border border-slate-200/90 bg-white px-[22px] py-6 shadow-[0_14px_34px_rgba(15,23,42,0.08)]">
            <div id="article-content" class="article-prose prose-pre:m-0" v-html="renderedContent" />
          </div>
          <MarkdownCodeEnhancer container-selector="#article-content" />
        </article>

        <aside class="space-y-4 xl:sticky xl:top-[88px] xl:max-h-[calc(100vh-104px)] xl:overflow-y-auto xl:pr-1">
          <ReadingToc v-if="tocHeadings.length > 0" :headings="tocHeadings" />

          <div class="glass-layer rounded-2xl p-4">
            <h3 class="text-xs font-bold tracking-[0.16em] text-slate-700">阅读信息</h3>
            <ul class="mt-3 space-y-2">
              <li class="flex items-center justify-between rounded-lg border border-white/60 bg-white/58 px-3 py-2 text-xs text-slate-600">
                <span>预计阅读</span>
                <span>{{ readingMinutes }} 分钟</span>
              </li>
              <li class="flex items-center justify-between rounded-lg border border-white/60 bg-white/58 px-3 py-2 text-xs text-slate-600">
                <span>全文字数</span>
                <span>{{ wordCount.toLocaleString() }}</span>
              </li>
              <li class="flex items-center justify-between rounded-lg border border-white/60 bg-white/58 px-3 py-2 text-xs text-slate-600">
                <span>目录节点</span>
                <span>{{ tocHeadings.length }}</span>
              </li>
            </ul>
          </div>

          <div class="glass-layer rounded-2xl p-4">
            <h3 class="text-xs font-bold tracking-[0.16em] text-slate-700">标签</h3>
            <div class="mt-3 flex flex-wrap gap-2">
              <span
                v-for="tag in post.tags || []"
                :key="`meta-${tag.slug || tag.name}`"
                class="rounded-full border border-miku/30 bg-miku-soft px-2 py-0.5 text-[11px] text-miku"
              >
                {{ tag.name }}
              </span>
              <span
                v-if="!post.tags || post.tags.length === 0"
                class="rounded-full border border-slate-200 bg-white px-2 py-0.5 text-[11px] text-slate-500"
              >
                暂无标签
              </span>
            </div>
          </div>
        </aside>
      </section>

      <section class="mx-auto w-full max-w-[1100px] space-y-4">
        <div v-if="relatedPosts.length > 0" class="glass-layer rounded-2xl p-4">
          <h3 class="text-xs font-bold tracking-[0.16em] text-slate-700">相关文章</h3>
          <div class="mt-3 space-y-2.5">
            <a
              v-for="item in relatedPosts"
              :key="item.id"
              :href="toPostUrl(item.slug)"
              class="block rounded-xl border border-white/65 bg-white/58 px-3 py-2.5 transition hover:border-miku/30"
            >
              <p class="line-clamp-1 text-sm font-semibold text-slate-800">{{ item.title }}</p>
              <p class="mt-1 line-clamp-2 text-xs text-slate-500">{{ item.excerpt }}</p>
              <p class="mt-1 text-[11px] text-slate-400">{{ formatDate(item.published_at || item.created_at) }}</p>
            </a>
          </div>
        </div>

        <PostCommentsSection
          id="post-comments"
          :post-id="post.id"
          :initial-total="commentCount"
          @count-updated="handleCommentCountUpdated"
        />
      </section>
    </article>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

import { api, type PagedData } from '../../lib/api'
import MarkdownCodeEnhancer from './MarkdownCodeEnhancer.vue'
import PostCommentsSection from './PostCommentsSection.vue'
import PostLikeBar from './PostLikeBar.vue'
import ReadingToc from './ReadingToc.vue'

interface TagItem {
  name: string
  slug: string
}

interface PostDetail {
  id: string
  slug: string
  title: string
  excerpt: string
  content_markdown: string
  hero_image_url: string
  category: string
  status: string
  published_at?: string
  view_count: number
  like_count: number
  comment_count: number
  created_at: string
  updated_at: string
  tags: TagItem[]
  liked: boolean
}

interface PostItem {
  id: string
  slug: string
  title: string
  excerpt: string
  published_at?: string
  created_at: string
}

interface HeadingItem {
  depth: number
  slug: string
  text: string
}

const post = ref<PostDetail | null>(null)
const loading = ref(true)
const error = ref('')
const markdownHtml = ref('')
const tocHeadings = ref<HeadingItem[]>([])
const relatedPosts = ref<PostItem[]>([])
const wordCount = ref(0)
const readingMinutes = ref(1)
const commentCount = ref(0)

const renderedContent = computed(() => markdownHtml.value)

let markedModule: Awaited<ReturnType<typeof import('marked')>> | null = null

function escapeHtml(input: string): string {
  return input
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#39;')
}

function normalizeMarkdown(source: string): string {
  return source
    .replace(/```[\s\S]*?```/g, ' ')
    .replace(/`[^`]*`/g, ' ')
    .replace(/!\[[^\]]*]\([^)]*\)/g, ' ')
    .replace(/\[([^\]]+)]\([^)]*\)/g, '$1')
    .replace(/<\/?[^>]+(>|$)/g, ' ')
    .replace(/[#>*_~\-]/g, ' ')
}

function getReadingStats(source: string) {
  const normalized = normalizeMarkdown(source)
  const latinWords = normalized.match(/[A-Za-z0-9]+/g) ?? []
  const cjkChars = normalized.match(/[\u3400-\u9fff]/g) ?? []
  const words = latinWords.length + cjkChars.length
  return {
    wordCount: words,
    readingMinutes: Math.max(1, Math.ceil(words / 260)),
  }
}

function formatDate(iso?: string): string {
  if (!iso) return '--'
  try {
    const d = new Date(iso)
    return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
  } catch {
    return iso.slice(0, 10)
  }
}

function slugifyHeading(text: string): string {
  const normalized = text
    .trim()
    .toLowerCase()
    .normalize('NFKD')
    .replace(/[\u0300-\u036f]/g, '')
    .replace(/[^\p{Letter}\p{Number}\u4e00-\u9fff\s-]/gu, '')
    .replace(/\s+/g, '-')
    .replace(/-+/g, '-')
    .replace(/^-|-$/g, '')

  return normalized || 'section'
}

function collectHeadings(markdown: string, markedLib: Awaited<ReturnType<typeof import('marked')>>): HeadingItem[] {
  const counters = new Map<string, number>()
  const collected: HeadingItem[] = []
  const tokens = markedLib.marked.lexer(markdown)

  const visit = (list: unknown[]) => {
    list.forEach((token) => {
      if (!token || typeof token !== 'object') {
        return
      }

      const t = token as {
        type?: string
        depth?: number
        text?: string
        tokens?: unknown[]
        items?: Array<{ tokens?: unknown[] }>
      }

      if (t.type === 'heading') {
        const rawText = t.text?.trim() || `章节 ${collected.length + 1}`
        const base = slugifyHeading(rawText)
        const nextIndex = (counters.get(base) ?? 0) + 1
        counters.set(base, nextIndex)
        const slug = nextIndex === 1 ? base : `${base}-${nextIndex}`
        collected.push({
          depth: Number(t.depth) || 2,
          slug,
          text: rawText,
        })
        return
      }

      if (t.type === 'blockquote' && Array.isArray(t.tokens)) {
        visit(t.tokens)
      }

      if (t.type === 'list' && Array.isArray(t.items)) {
        t.items.forEach((item) => {
          if (Array.isArray(item.tokens)) {
            visit(item.tokens)
          }
        })
      }
    })
  }

  visit(tokens)
  return collected
}

async function renderMarkdown(markdown: string): Promise<{ html: string; headings: HeadingItem[] }> {
  if (!markdown) {
    return {
      html: '',
      headings: [],
    }
  }

  try {
    if (!markedModule) {
      markedModule = await import('marked')
    }

    const headings = collectHeadings(markdown, markedModule)
    let headingCursor = 0

    const renderer = new markedModule.marked.Renderer()
    renderer.heading = function ({ tokens, depth }: { tokens: unknown[]; depth: number }) {
      const content = this.parser.parseInline(tokens)
      const fallbackSlug = `section-${headingCursor + 1}`
      const slug = headings[headingCursor]?.slug || fallbackSlug
      headingCursor += 1
      return `<h${depth} id="${slug}" data-heading-id="${slug}">${content}</h${depth}>`
    }

    const html = await markedModule.marked.parse(markdown, {
      gfm: true,
      breaks: false,
      renderer,
    })

    return {
      html: String(html),
      headings: headings.filter((item) => item.depth <= 3),
    }
  } catch {
    const readable = escapeHtml(markdown).replace(/\n/g, '<br>')
    return {
      html: `<p>${readable}</p>`,
      headings: [],
    }
  }
}

function getSlugFromUrl(): string {
  if (typeof window === 'undefined') return ''
  const params = new URLSearchParams(window.location.search)
  return params.get('slug') || ''
}

function toPostUrl(slug: string): string {
  return `/blog/post?slug=${encodeURIComponent(slug)}`
}

async function loadRelatedPosts(currentSlug: string) {
  try {
    const data = await api.get<PagedData<PostItem>>('/posts?page=1&size=8')
    relatedPosts.value = (data.items || []).filter((item) => item.slug !== currentSlug).slice(0, 3)
  } catch {
    relatedPosts.value = []
  }
}

async function loadPost() {
  const slug = getSlugFromUrl()
  if (!slug) {
    error.value = '缺少文章标识'
    loading.value = false
    return
  }

  loading.value = true
  error.value = ''
  markdownHtml.value = ''
  tocHeadings.value = []
  relatedPosts.value = []

  try {
    post.value = await api.get<PostDetail>(`/posts/${encodeURIComponent(slug)}`)
    commentCount.value = Number(post.value.comment_count || 0)
    const markdownSource = post.value.content_markdown || ''
    const rendered = await renderMarkdown(markdownSource)
    markdownHtml.value = rendered.html
    tocHeadings.value = rendered.headings

    const stats = getReadingStats(markdownSource || post.value.excerpt || '')
    wordCount.value = stats.wordCount
    readingMinutes.value = stats.readingMinutes

    await loadRelatedPosts(post.value.slug || slug)
  } catch {
    error.value = '文章加载失败，请检查链接是否正确'
  } finally {
    loading.value = false
  }
}

function handleLikeUpdated(payload: { liked: boolean; likeCount: number }) {
  if (!post.value) return
  post.value.liked = payload.liked
  post.value.like_count = payload.likeCount
}

function handleCommentCountUpdated(nextCount: number) {
  commentCount.value = Number(nextCount || 0)
}

onMounted(() => {
  loadPost()
})
</script>
