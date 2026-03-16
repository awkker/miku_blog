<template>
  <div class="mx-auto max-w-3xl">
    <!-- Loading -->
    <div v-if="loading" class="space-y-6 py-12">
      <div class="animate-pulse space-y-4">
        <div class="h-8 w-3/4 rounded bg-slate-200/60" />
        <div class="h-4 w-1/3 rounded bg-slate-100/80" />
        <div class="h-64 rounded-2xl bg-slate-200/40" />
        <div class="space-y-2">
          <div class="h-4 w-full rounded bg-slate-100/60" />
          <div class="h-4 w-5/6 rounded bg-slate-100/60" />
          <div class="h-4 w-4/6 rounded bg-slate-100/60" />
        </div>
      </div>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="py-20 text-center">
      <div class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-red-50">
        <svg viewBox="0 0 24 24" class="h-8 w-8 fill-none stroke-red-400 stroke-[1.5]">
          <circle cx="12" cy="12" r="10" /><line x1="15" y1="9" x2="9" y2="15" /><line x1="9" y1="9" x2="15" y2="15" />
        </svg>
      </div>
      <p class="text-sm text-red-600">{{ error }}</p>
      <a href="/blog" class="mt-4 inline-flex items-center gap-1 rounded-xl border border-miku/35 bg-miku-soft px-4 py-2 text-sm text-miku transition hover:border-miku/55">
        返回文章列表
      </a>
    </div>

    <!-- Post content -->
    <article v-else-if="post" class="space-y-8">
      <!-- Back link -->
      <a href="/blog" class="inline-flex items-center gap-1.5 rounded-lg border border-slate-200 bg-white/60 px-3 py-1.5 text-xs text-slate-500 transition hover:border-miku/40 hover:text-miku">
        <svg viewBox="0 0 24 24" class="h-3.5 w-3.5 fill-none stroke-current stroke-[2]"><path d="M19 12H5M12 19l-7-7 7-7" /></svg>
        返回文章列表
      </a>

      <!-- Hero image -->
      <div v-if="post.hero_image_url" class="overflow-hidden rounded-2xl">
        <img :src="post.hero_image_url" :alt="post.title" class="w-full object-cover" loading="lazy" />
      </div>

      <!-- Meta -->
      <header class="space-y-4">
        <div class="flex flex-wrap items-center gap-2">
          <span v-if="post.category" class="rounded-full border border-miku/35 bg-miku-soft px-3 py-1 text-xs font-semibold text-miku">
            {{ post.category }}
          </span>
          <span v-for="tag in (post.tags || [])" :key="tag.name" class="rounded-full border border-slate-200 bg-white px-2.5 py-1 text-[11px] text-slate-500">
            {{ tag.name }}
          </span>
        </div>

        <h1 class="text-3xl font-black leading-tight text-slate-900 md:text-4xl">
          {{ post.title }}
        </h1>

        <p v-if="post.excerpt" class="text-base leading-relaxed text-slate-500">
          {{ post.excerpt }}
        </p>

        <div class="flex flex-wrap items-center gap-4 text-xs text-slate-400">
          <span class="inline-flex items-center gap-1">
            <svg viewBox="0 0 24 24" class="h-3.5 w-3.5 fill-none stroke-current stroke-[1.8]"><rect x="3" y="4" width="18" height="18" rx="2" ry="2" /><line x1="16" y1="2" x2="16" y2="6" /><line x1="8" y1="2" x2="8" y2="6" /><line x1="3" y1="10" x2="21" y2="10" /></svg>
            {{ formatDate(post.published_at || post.created_at) }}
          </span>
          <span class="inline-flex items-center gap-1">
            <svg viewBox="0 0 24 24" class="h-3.5 w-3.5 fill-none stroke-current stroke-[1.5]"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" /><circle cx="12" cy="12" r="3" /></svg>
            {{ post.view_count.toLocaleString() }} 次阅读
          </span>
          <span class="inline-flex items-center gap-1">
            <svg viewBox="0 0 24 24" class="h-3.5 w-3.5 fill-none stroke-current stroke-[1.5]"><path d="M20.84 4.61a5.5 5.5 0 00-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 00-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 000-7.78z" /></svg>
            {{ post.like_count }} 点赞
          </span>
        </div>
      </header>

      <!-- Content -->
      <div class="prose prose-slate max-w-none prose-headings:font-bold prose-a:text-miku prose-a:no-underline hover:prose-a:underline prose-img:rounded-xl" v-html="renderedContent" />

      <!-- Like button -->
      <div class="flex items-center justify-center border-t border-slate-200/60 pt-8">
        <button
          type="button"
          :class="[
            'inline-flex items-center gap-2 rounded-2xl border px-6 py-3 text-sm font-semibold transition',
            post.liked
              ? 'border-red-200 bg-red-50 text-red-500'
              : 'border-slate-200 bg-white/60 text-slate-500 hover:border-red-200 hover:text-red-500',
          ]"
          @click="toggleLike"
        >
          <svg viewBox="0 0 24 24" :class="['h-5 w-5 stroke-[1.5] transition', post.liked ? 'fill-red-500 stroke-red-500' : 'fill-none stroke-current']">
            <path d="M20.84 4.61a5.5 5.5 0 00-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 00-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 000-7.78z" />
          </svg>
          {{ post.liked ? '已点赞' : '点赞' }} {{ post.like_count }}
        </button>
      </div>
    </article>
  </div>
</template>

<script setup lang="ts">
import { marked } from 'marked'
import { computed, onMounted, ref } from 'vue'

import { api } from '../../lib/api'

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

const post = ref<PostDetail | null>(null)
const loading = ref(true)
const error = ref('')

const renderedContent = computed(() => {
  if (!post.value?.content_markdown) return ''
  return marked.parse(post.value.content_markdown) as string
})

function formatDate(iso?: string): string {
  if (!iso) return '--'
  try {
    const d = new Date(iso)
    return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
  } catch {
    return iso.slice(0, 10)
  }
}

function getSlugFromUrl(): string {
  if (typeof window === 'undefined') return ''
  const params = new URLSearchParams(window.location.search)
  return params.get('slug') || ''
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
  try {
    post.value = await api.get<PostDetail>(`/posts/${slug}`)
  } catch {
    error.value = '文章加载失败，请检查链接是否正确'
  } finally {
    loading.value = false
  }
}

async function toggleLike() {
  if (!post.value) return
  try {
    const res = await api.post<{ liked: boolean }>(`/posts/${post.value.id}/like`)
    post.value.liked = res.liked
    post.value.like_count += res.liked ? 1 : -1
  } catch {
    // silently fail for like toggle
  }
}

onMounted(() => {
  loadPost()
})
</script>
