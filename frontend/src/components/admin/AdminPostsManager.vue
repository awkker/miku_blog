<template>
  <section class="space-y-5">
    <LiquidGlassCard padding="24px">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 class="text-2xl font-semibold text-slate-900">文章管理</h1>
          <p class="mt-1 text-sm text-slate-600">管理博客文章的发布、草稿与分类。</p>
        </div>
        <MikuButton variant="solid" aria-label="新建文章" @click="toggleCreateForm">+ 新建文章</MikuButton>
      </div>
    </LiquidGlassCard>

    <!-- Create Form -->
    <LiquidGlassCard v-if="showCreateForm" padding="24px">
      <h2 class="mb-4 text-lg font-semibold text-slate-900">新建文章</h2>
      <form class="space-y-3" @submit.prevent="createPost">
        <input v-model="newPost.title" type="text" placeholder="文章标题 *" class="w-full rounded-xl border border-slate-200/80 bg-white/60 px-3 py-2 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-1 focus:ring-miku/30" />
        <div class="grid gap-3 md:grid-cols-2">
          <input v-model="newPost.category" type="text" placeholder="分类" class="rounded-xl border border-slate-200/80 bg-white/60 px-3 py-2 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-1 focus:ring-miku/30" />
          <input v-model="newPost.tags" type="text" placeholder="标签 (逗号分隔)" class="rounded-xl border border-slate-200/80 bg-white/60 px-3 py-2 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-1 focus:ring-miku/30" />
        </div>
        <input v-model="newPost.excerpt" type="text" placeholder="摘要" class="w-full rounded-xl border border-slate-200/80 bg-white/60 px-3 py-2 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-1 focus:ring-miku/30" />
        <input v-model="newPost.hero_image_url" type="text" placeholder="封面图片 URL" class="w-full rounded-xl border border-slate-200/80 bg-white/60 px-3 py-2 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-1 focus:ring-miku/30" />
        <textarea v-model="newPost.content_markdown" rows="8" placeholder="正文内容 (Markdown)" class="w-full rounded-xl border border-slate-200/80 bg-white/60 px-3 py-2 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-1 focus:ring-miku/30" />
        <div class="flex items-center gap-3">
          <select v-model="newPost.status" class="rounded-xl border border-slate-200/80 bg-white/60 px-3 py-2 text-sm text-slate-900 outline-none">
            <option value="draft">草稿</option>
            <option value="published">直接发布</option>
          </select>
          <MikuButton type="submit" variant="solid" :disabled="creating">{{ creating ? '创建中...' : '创建文章' }}</MikuButton>
          <button type="button" class="text-sm text-slate-500 hover:text-slate-700" @click="closeCreateForm">取消</button>
        </div>
      </form>
    </LiquidGlassCard>

    <LiquidGlassCard v-if="showEditForm" padding="24px">
      <h2 class="mb-4 text-lg font-semibold text-slate-900">编辑文章</h2>
      <form class="space-y-3" @submit.prevent="updatePost">
        <input v-model="editPost.title" type="text" placeholder="文章标题 *" class="w-full rounded-xl border border-slate-200/80 bg-white/60 px-3 py-2 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-1 focus:ring-miku/30" />
        <div class="grid gap-3 md:grid-cols-2">
          <input v-model="editPost.category" type="text" placeholder="分类" class="rounded-xl border border-slate-200/80 bg-white/60 px-3 py-2 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-1 focus:ring-miku/30" />
          <input v-model="editPost.tags" type="text" placeholder="标签 (逗号分隔)" class="rounded-xl border border-slate-200/80 bg-white/60 px-3 py-2 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-1 focus:ring-miku/30" />
        </div>
        <input v-model="editPost.excerpt" type="text" placeholder="摘要" class="w-full rounded-xl border border-slate-200/80 bg-white/60 px-3 py-2 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-1 focus:ring-miku/30" />
        <input v-model="editPost.hero_image_url" type="text" placeholder="封面图片 URL" class="w-full rounded-xl border border-slate-200/80 bg-white/60 px-3 py-2 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-1 focus:ring-miku/30" />
        <textarea v-model="editPost.content_markdown" rows="8" placeholder="正文内容 (Markdown)" class="w-full rounded-xl border border-slate-200/80 bg-white/60 px-3 py-2 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-1 focus:ring-miku/30" />
        <div class="flex items-center gap-3">
          <MikuButton type="submit" variant="solid" :disabled="editing">{{ editing ? '保存中...' : '保存修改' }}</MikuButton>
          <button type="button" class="text-sm text-slate-500 hover:text-slate-700" @click="closeEditForm">取消</button>
        </div>
      </form>
    </LiquidGlassCard>

    <LiquidGlassCard padding="0px">
      <div class="overflow-x-auto">
        <table class="w-full text-left text-sm">
          <thead>
            <tr class="border-b border-slate-200/60">
              <th class="px-5 py-3.5 font-semibold text-slate-700">标题</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700">分类</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700">状态</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700 text-right">浏览量</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700">发布时间</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700 text-center">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="post in posts"
              :key="post.id"
              class="border-b border-slate-100/60 transition hover:bg-white/40"
            >
              <td class="px-5 py-3.5 font-medium text-slate-900">{{ post.title }}</td>
              <td class="px-5 py-3.5 text-slate-600">{{ post.category }}</td>
              <td class="px-5 py-3.5">
                <span
                  class="inline-block rounded-full px-2.5 py-0.5 text-xs font-medium"
                  :class="statusClass(post.status)"
                >
                  {{ statusLabel(post.status) }}
                </span>
              </td>
              <td class="px-5 py-3.5 text-right font-mono text-slate-700">{{ post.views.toLocaleString() }}</td>
              <td class="px-5 py-3.5 text-slate-600">{{ post.publishedAt }}</td>
              <td class="px-5 py-3.5 text-center">
                <div class="flex items-center justify-center gap-2">
                  <button
                    type="button"
                    class="rounded-xl border border-slate-200/80 bg-white/50 px-2.5 py-1 text-xs text-slate-700 transition hover:border-miku/40 hover:text-miku"
                    aria-label="编辑文章"
                    @click="startEditPost(post.id)"
                  >
                    编辑
                  </button>
                  <button
                    v-if="post.status !== 'published'"
                    type="button"
                    class="rounded-xl border border-emerald-200/80 bg-white/50 px-2.5 py-1 text-xs text-emerald-600 transition hover:bg-emerald-50"
                    aria-label="发布文章"
                    @click="publishPost(post.id)"
                  >
                    发布
                  </button>
                  <button
                    type="button"
                    class="rounded-xl border border-red-200/80 bg-white/50 px-2.5 py-1 text-xs text-red-600 transition hover:bg-red-50"
                    aria-label="删除文章"
                    @click="deletePost(post.id)"
                  >
                    删除
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </LiquidGlassCard>
  </section>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'

import { api, ApiError, type PagedData } from '../../lib/api'
import { showToast } from '../../stores/ui'
import LiquidGlassCard from '../ui/LiquidGlassCard.vue'
import MikuButton from '../ui/MikuButton.vue'

interface ApiPost {
  id: string
  slug: string
  title: string
  category: string
  status: string
  published_at?: string
  view_count: number
  like_count: number
  comment_count: number
  created_at: string
  updated_at: string
}

interface ApiTag {
  name: string
  slug: string
}

interface ApiPostDetail {
  id: string
  slug: string
  title: string
  excerpt: string
  content_markdown: string
  hero_image_url: string
  category: string
  status: string
  tags: ApiTag[]
}

interface Post {
  id: string
  title: string
  category: string
  status: 'published' | 'draft' | 'scheduled'
  views: number
  publishedAt: string
}

interface PostForm {
  title: string
  category: string
  excerpt: string
  content_markdown: string
  hero_image_url: string
  tags: string
}

interface NewPostForm extends PostForm {
  status: 'draft' | 'published'
}

function mapStatus(s: string): 'published' | 'draft' | 'scheduled' {
  if (s === 'published') return 'published'
  if (s === 'scheduled') return 'scheduled'
  return 'draft'
}

function formatDate(iso?: string): string {
  if (!iso) return '--'
  try {
    return iso.slice(0, 10)
  } catch {
    return iso
  }
}

function createEmptyPostForm(): PostForm {
  return {
    title: '',
    category: '',
    excerpt: '',
    content_markdown: '',
    hero_image_url: '',
    tags: '',
  }
}

function createEmptyNewPostForm(): NewPostForm {
  return {
    ...createEmptyPostForm(),
    status: 'draft',
  }
}

function toTagArray(input: string): string[] {
  return input.split(',').map((t) => t.trim()).filter(Boolean)
}

function toTagInput(tags: ApiTag[] | undefined): string {
  if (!tags || tags.length === 0) return ''
  return tags.map((t) => t.name).filter(Boolean).join(', ')
}

function generateSlugFromTitle(title: string): string {
  const normalized = title
    .trim()
    .toLowerCase()
    .normalize('NFKD')
    .replace(/[\u0300-\u036f]/g, '')
    .replace(/[^\p{L}\p{N}\s-]/gu, '')
    .replace(/\s+/g, '-')
    .replace(/-+/g, '-')
    .replace(/^-|-$/g, '')

  if (normalized) return normalized
  return `post-${Date.now()}`
}

function toPostPayload(form: PostForm) {
  const title = form.title.trim()
  return {
    title,
    slug: generateSlugFromTitle(title),
    category: form.category.trim(),
    excerpt: form.excerpt.trim(),
    content_markdown: form.content_markdown.trim(),
    hero_image_url: form.hero_image_url.trim(),
    tags: toTagArray(form.tags),
  }
}

function mapPost(item: ApiPost): Post {
  return {
    id: item.id,
    title: item.title,
    category: item.category || '--',
    status: mapStatus(item.status),
    views: Number(item.view_count) || 0,
    publishedAt: item.status === 'published' ? formatDate(item.published_at) : '--',
  }
}

const posts = ref<Post[]>([])
const loading = ref(false)
const showCreateForm = ref(false)
const creating = ref(false)
const showEditForm = ref(false)
const editing = ref(false)
const editingPostID = ref<string | null>(null)

const newPost = ref<NewPostForm>(createEmptyNewPostForm())
const editPost = ref<PostForm>(createEmptyPostForm())

function toggleCreateForm() {
  showCreateForm.value = !showCreateForm.value
  if (showCreateForm.value) {
    closeEditForm()
  }
}

function closeCreateForm() {
  showCreateForm.value = false
}

function closeEditForm() {
  showEditForm.value = false
  editingPostID.value = null
  editPost.value = createEmptyPostForm()
}

async function loadPosts() {
  loading.value = true
  try {
    const data = await api.get<PagedData<ApiPost>>('/admin/posts?size=50')
    posts.value = (data.items || []).map(mapPost)
  } catch (err) {
    console.error('[AdminPosts] loadPosts failed:', err)
    showToast('加载文章列表失败', 'error')
    posts.value = []
  } finally {
    loading.value = false
  }
}

async function createPost() {
  if (!newPost.value.title.trim()) return
  creating.value = true
  try {
    await api.post('/admin/posts', {
      ...toPostPayload(newPost.value),
      status: newPost.value.status,
    })
    closeCreateForm()
    newPost.value = createEmptyNewPostForm()
    showToast('文章创建成功', 'success')
    await loadPosts()
  } catch (err) {
    const msg = err instanceof ApiError ? err.message : '创建文章失败，请稍后重试'
    console.error('[AdminPosts] createPost failed:', err)
    showToast(msg, 'error')
  } finally {
    creating.value = false
  }
}

async function startEditPost(id: string) {
  try {
    const detail = await api.get<ApiPostDetail>(`/admin/posts/${id}`)
    editingPostID.value = detail.id
    editPost.value = {
      title: detail.title || '',
      category: detail.category || '',
      excerpt: detail.excerpt || '',
      content_markdown: detail.content_markdown || '',
      hero_image_url: detail.hero_image_url || '',
      tags: toTagInput(detail.tags),
    }
    showEditForm.value = true
    showCreateForm.value = false
  } catch (err) {
    const msg = err instanceof ApiError ? err.message : '加载文章详情失败'
    console.error('[AdminPosts] startEditPost failed:', err)
    showToast(msg, 'error')
  }
}

async function updatePost() {
  if (!editingPostID.value) return
  if (!editPost.value.title.trim()) return
  editing.value = true
  try {
    await api.put(`/admin/posts/${editingPostID.value}`, toPostPayload(editPost.value))
    showToast('文章更新成功', 'success')
    closeEditForm()
    await loadPosts()
  } catch (err) {
    const msg = err instanceof ApiError ? err.message : '更新文章失败，请稍后重试'
    console.error('[AdminPosts] updatePost failed:', err)
    showToast(msg, 'error')
  } finally {
    editing.value = false
  }
}

async function publishPost(id: string) {
  try {
    await api.post(`/admin/posts/${id}/publish`)
    posts.value = posts.value.map((p) => p.id === id ? { ...p, status: 'published' as const } : p)
    showToast('文章发布成功', 'success')
  } catch (err) {
    const msg = err instanceof ApiError ? err.message : '发布失败，请稍后重试'
    console.error('[AdminPosts] publishPost failed:', err)
    showToast(msg, 'error')
  }
}

async function deletePost(id: string) {
  try {
    await api.delete(`/admin/posts/${id}`)
    posts.value = posts.value.filter((p) => p.id !== id)
    showToast('文章已删除', 'success')
  } catch (err) {
    const msg = err instanceof ApiError ? err.message : '删除失败，请稍后重试'
    console.error('[AdminPosts] deletePost failed:', err)
    showToast(msg, 'error')
  }
}

onMounted(() => {
  loadPosts()
})

function statusClass(status: Post['status']) {
  if (status === 'published') return 'bg-emerald-100 text-emerald-700'
  if (status === 'draft') return 'bg-slate-100 text-slate-600'
  return 'bg-purple-100 text-purple-700'
}

function statusLabel(status: Post['status']) {
  if (status === 'published') return '已发布'
  if (status === 'draft') return '草稿'
  return '定时发布'
}
</script>
