<template>
  <div class="space-y-7">
    <!-- Loading -->
    <div v-if="loading" class="space-y-6">
      <div v-for="i in 3" :key="i" class="animate-pulse rounded-3xl border border-white/55 bg-white/72 p-6">
        <div class="h-52 rounded-2xl bg-slate-200/60" />
        <div class="mt-4 h-5 w-2/3 rounded bg-slate-200/60" />
        <div class="mt-2 h-4 w-1/2 rounded bg-slate-100/80" />
      </div>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="rounded-2xl border border-red-200/60 bg-red-50/60 p-8 text-center">
      <p class="text-sm text-red-600">{{ error }}</p>
      <button type="button" class="mt-3 rounded-xl border border-red-200 bg-white px-4 py-2 text-xs text-red-600 transition hover:bg-red-50" @click="load">
        {{ copy.retry }}
      </button>
    </div>

    <!-- Empty -->
    <div v-else-if="posts.length === 0" class="rounded-2xl border border-slate-200/60 bg-white/60 p-12 text-center">
      <p class="text-sm text-slate-500">{{ copy.empty }}</p>
    </div>

    <!-- Post list -->
    <template v-else>
      <!-- Featured post -->
      <article v-if="featured" class="blog-card group overflow-hidden rounded-3xl border border-white/65 bg-white/78 shadow-[0_20px_46px_rgba(15,23,42,0.14)] backdrop-blur">
        <a :href="`/blog/post?slug=${featured.slug}`" class="block lg:grid lg:grid-cols-[1.2fr_1fr]">
          <div class="relative min-h-[260px] overflow-hidden lg:min-h-full">
            <img
              v-if="featured.hero_image_url"
              :src="featured.hero_image_url"
              :alt="featured.title"
              class="h-full w-full object-cover transition duration-700 group-hover:scale-105"
              loading="lazy"
            />
            <div v-else class="flex h-full min-h-[260px] items-center justify-center bg-gradient-to-br from-miku/10 to-[#c084fc]/10">
              <svg viewBox="0 0 24 24" class="h-16 w-16 fill-none stroke-miku/30 stroke-[1]"><path d="M4 19.5A2.5 2.5 0 016.5 17H20" /><path d="M6.5 2H20v20H6.5A2.5 2.5 0 014 19.5v-15A2.5 2.5 0 016.5 2z" /></svg>
            </div>
            <div class="absolute inset-0 bg-gradient-to-t from-slate-900/35 via-transparent to-transparent" />
            <div class="absolute left-4 top-4 rounded-full border border-white/45 bg-white/22 px-3 py-1 text-xs font-semibold text-white backdrop-blur">
              {{ copy.featuredBadge }}
            </div>
          </div>

          <div class="space-y-4 p-5 md:p-6">
            <div class="flex flex-wrap items-center gap-2">
              <span v-if="featured.category" class="inline-flex items-center gap-1 rounded-full border border-miku/35 bg-miku-soft px-3 py-1 text-xs font-semibold text-miku">
                {{ featured.category }}
              </span>
              <span v-if="featured.tags && featured.tags.length > 0" v-for="tag in featured.tags.slice(0, 2)" :key="tag.name" class="rounded-full border border-slate-200 bg-white px-2.5 py-1 text-[11px] text-slate-500">
                {{ tag.name }}
              </span>
            </div>

            <h2 class="text-2xl font-black leading-tight text-slate-900 transition group-hover:text-miku">
              {{ featured.title }}
            </h2>
            <p class="text-sm leading-relaxed text-slate-600 md:text-base">
              {{ featured.excerpt }}
            </p>

            <div class="grid gap-2 text-xs text-slate-500 sm:grid-cols-3">
              <span class="rounded-lg border border-slate-200 bg-white/75 px-2.5 py-1.5">
                {{ featured.view_count.toLocaleString() }}{{ copy.readSuffix }}
              </span>
              <span class="rounded-lg border border-slate-200 bg-white/75 px-2.5 py-1.5">
                {{ featured.like_count }}{{ copy.likeSuffix }}
              </span>
              <span class="rounded-lg border border-slate-200 bg-white/75 px-2.5 py-1.5">
                {{ copy.publishedPrefix }}{{ formatDate(featured.published_at || featured.created_at) }}
              </span>
            </div>
          </div>
        </a>
      </article>

      <!-- Grid -->
      <div class="grid grid-cols-1 gap-7 md:grid-cols-2">
        <article
          v-for="(post, index) in restPosts"
          :key="post.id"
          :class="[
            'blog-card group relative overflow-hidden rounded-3xl border border-white/55 bg-white/72 shadow-lg backdrop-blur transition duration-500 hover:border-miku/25 hover:shadow-[0_16px_40px_rgba(57,197,187,0.12)]',
            index > 0 && index % 4 === 0 && 'md:col-span-2',
          ]"
        >
          <a
            :href="`/blog/post?slug=${post.slug}`"
            :class="[
              'block h-full',
              index > 0 && index % 4 === 0 && 'md:grid md:grid-cols-[320px_1fr] md:items-stretch',
            ]"
          >
            <div :class="['relative overflow-hidden', index > 0 && index % 4 === 0 ? 'h-56 md:h-full' : 'h-52']">
              <img
                v-if="post.hero_image_url"
                :src="post.hero_image_url"
                :alt="post.title"
                class="h-full w-full object-cover transition duration-700 group-hover:scale-105"
                loading="lazy"
              />
              <div v-else class="flex h-full items-center justify-center bg-gradient-to-br from-miku/8 to-[#c084fc]/8">
                <svg viewBox="0 0 24 24" class="h-12 w-12 fill-none stroke-miku/25 stroke-[1]"><path d="M4 19.5A2.5 2.5 0 016.5 17H20" /><path d="M6.5 2H20v20H6.5A2.5 2.5 0 014 19.5v-15A2.5 2.5 0 016.5 2z" /></svg>
              </div>
              <div class="absolute inset-0 bg-gradient-to-t from-black/20 via-transparent to-transparent opacity-0 transition duration-500 group-hover:opacity-100" />
            </div>

            <div class="space-y-3 px-5 py-5">
              <div class="flex flex-wrap items-center gap-2">
                <span v-if="post.category" class="inline-flex items-center gap-1 rounded-full border border-miku/35 bg-miku-soft px-3 py-1 text-xs font-semibold text-miku">
                  {{ post.category }}
                </span>
                <span v-for="tag in (post.tags || []).slice(0, 2)" :key="tag.name" class="inline-flex rounded-full border border-slate-200 bg-white px-2.5 py-1 text-[11px] text-slate-500">
                  {{ tag.name }}
                </span>
              </div>
              <h2 class="text-xl font-bold text-slate-900 transition duration-300 group-hover:text-miku">{{ post.title }}</h2>
              <p class="overflow-hidden text-sm leading-relaxed text-slate-500 [display:-webkit-box] [-webkit-box-orient:vertical] [-webkit-line-clamp:2]">
                {{ post.excerpt }}
              </p>
              <div class="flex flex-wrap items-center gap-3 text-xs text-slate-400">
                <span class="inline-flex items-center gap-1">
                  <svg viewBox="0 0 24 24" class="h-3.5 w-3.5 fill-none stroke-current stroke-[1.8]" aria-hidden="true">
                    <path d="M20 21V7a2 2 0 00-2-2h-3V3H9v2H6a2 2 0 00-2 2v14M16 11h-8M16 15h-8" />
                  </svg>
                  {{ formatDate(post.published_at || post.created_at) }}
                </span>
                <span class="h-0.5 w-0.5 rounded-full bg-slate-300" />
                <span class="inline-flex items-center gap-1 text-slate-500">
                  <svg viewBox="0 0 24 24" class="h-3.5 w-3.5 fill-none stroke-current stroke-[1.5]" aria-hidden="true"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" /><circle cx="12" cy="12" r="3" /></svg>
                  {{ post.view_count.toLocaleString() }}
                </span>
                <span class="h-0.5 w-0.5 rounded-full bg-slate-300" />
                <span>{{ post.like_count }}{{ copy.shortLikeSuffix }}</span>
              </div>
            </div>
          </a>
        </article>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="flex items-center justify-center gap-2 pt-4">
        <button
          v-for="p in totalPages"
          :key="p"
          type="button"
          :class="[
            'rounded-xl border px-3.5 py-2 text-sm transition',
            p === page
              ? 'border-miku/45 bg-miku-soft text-miku'
              : 'border-slate-200 bg-white/60 text-slate-500 hover:border-miku/30 hover:text-miku',
          ]"
          @click="goToPage(p)"
        >
          {{ p }}
        </button>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

import { api, type PagedData } from '../../lib/api'
import { siteCopy } from '../../content/copy'

interface TagItem {
  name: string
  slug: string
}

interface PostItem {
  id: string
  slug: string
  title: string
  excerpt: string
  hero_image_url: string
  category: string
  published_at?: string
  view_count: number
  like_count: number
  comment_count: number
  created_at: string
  tags?: TagItem[]
}

const posts = ref<PostItem[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = 10
const loading = ref(false)
const error = ref('')
const copy = siteCopy.components.blogFeed

const featured = computed(() => posts.value[0] || null)
const restPosts = computed(() => posts.value.slice(1))
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize)))

function formatDate(iso?: string): string {
  if (!iso) return '--'
  try {
    const d = new Date(iso)
    return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
  } catch {
    return iso.slice(0, 10)
  }
}

async function load() {
  loading.value = true
  error.value = ''
  try {
    const data = await api.get<PagedData<PostItem>>(`/posts?page=${page.value}&size=${pageSize}`)
    posts.value = data.items || []
    total.value = data.total || 0
  } catch {
    error.value = copy.loadError
    posts.value = []
  } finally {
    loading.value = false
  }
}

function goToPage(p: number) {
  page.value = p
  load()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

defineExpose({ postCount: total })

onMounted(() => {
  load()
})
</script>

<style scoped>
.blog-card {
  transform: translateY(0);
  transition: transform 0.5s cubic-bezier(0.22, 1, 0.36, 1), box-shadow 0.5s ease, border-color 0.5s ease;
}
.blog-card:hover {
  transform: translateY(-4px);
}
</style>
