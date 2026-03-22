<template>
  <section class="space-y-5">
    <!-- Header -->
    <LiquidGlassCard padding="24px">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 class="text-2xl font-semibold text-slate-900">文章管理</h1>
          <p class="mt-1 text-sm text-slate-600">管理博客文章的发布、草稿与分类。</p>
        </div>
        <MikuButton variant="solid" aria-label="新建文章" @click="toggleCreateForm">+ 新建文章</MikuButton>
      </div>
    </LiquidGlassCard>

    <!-- ===== Writing Studio (Create) ===== -->
    <LiquidGlassCard v-if="showCreateForm" padding="0px">
      <form @submit.prevent="createPost">
        <!-- Title zone -->
        <div class="writing-title-zone">
          <input
            v-model="newPost.title"
            type="text"
            placeholder="无题..."
            class="writing-title-input"
            autocomplete="off"
          />
          <p class="mt-2 text-xs tracking-wide text-slate-400">新文章</p>
        </div>

        <!-- Metadata drawer -->
        <div class="border-t border-slate-200/40">
          <button
            type="button"
            class="flex w-full items-center gap-2 px-8 py-3 text-left text-xs font-medium uppercase tracking-widest text-slate-400 transition hover:text-slate-600"
            @click="showCreateMeta = !showCreateMeta"
          >
            <svg class="h-3.5 w-3.5 transition-transform" :class="showCreateMeta ? 'rotate-90' : ''" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="m9 18 6-6-6-6" /></svg>
            文章属性
          </button>
          <div v-show="showCreateMeta" class="meta-drawer">
            <div class="meta-grid">
              <div class="meta-field">
                <label class="meta-label">分类</label>
                <input v-model="newPost.category" type="text" placeholder="技术 / 随笔 / 教程" class="meta-input" />
              </div>
              <div class="meta-field">
                <label class="meta-label">标签</label>
                <input v-model="newPost.tags" type="text" placeholder="用逗号分隔" class="meta-input" />
              </div>
              <div class="meta-field">
                <label class="meta-label">摘要</label>
                <input v-model="newPost.excerpt" type="text" placeholder="简短描述文章内容" class="meta-input" />
              </div>
              <div class="meta-field">
                <label class="meta-label">封面图片</label>
                <input v-model="newPost.hero_image_url" type="text" placeholder="https://..." class="meta-input" />
              </div>
            </div>
            <div v-if="newPost.hero_image_url" class="mx-8 mb-4 overflow-hidden rounded-xl border border-slate-200/50">
              <img :src="newPost.hero_image_url" alt="" class="h-28 w-full object-cover" @error="($event.target as HTMLImageElement).style.display='none'" />
            </div>
          </div>
        </div>

        <!-- Editor zone -->
        <div class="writing-editor-zone">
          <textarea
            v-model="newPost.content_markdown"
            placeholder="在此撰写你的想法..."
            class="writing-editor"
          />
        </div>

        <!-- Bottom bar -->
        <div class="flex items-center justify-between border-t border-slate-200/40 px-8 py-4">
          <div class="flex items-center gap-2 text-xs text-slate-400">
            <span class="inline-block h-1.5 w-1.5 rounded-full" :class="newPost.content_markdown.length > 0 ? 'bg-miku' : 'bg-slate-300'" />
            {{ charCount(newPost.content_markdown) }}
          </div>
          <div class="flex items-center gap-3">
            <select v-model="newPost.status" class="meta-input cursor-pointer py-1.5 pr-7 text-xs">
              <option value="draft">草稿</option>
              <option value="published">直接发布</option>
              <option value="scheduled">定时发布</option>
            </select>
            <input
              v-if="newPost.status === 'scheduled'"
              v-model="newPost.scheduled_at"
              type="datetime-local"
              class="meta-input py-1.5 text-xs"
            />
            <button type="button" class="rounded-xl px-4 py-2 text-sm text-slate-400 transition hover:bg-slate-100/50 hover:text-slate-600" @click="closeCreateForm">取消</button>
            <MikuButton type="submit" variant="solid" :disabled="creating">{{ creating ? '创建中...' : '创建文章' }}</MikuButton>
          </div>
        </div>
      </form>
    </LiquidGlassCard>

    <!-- ===== Writing Studio (Edit) ===== -->
    <LiquidGlassCard v-if="showEditForm" padding="0px">
      <form @submit.prevent="updatePost">
        <div class="writing-title-zone">
          <input
            v-model="editPost.title"
            type="text"
            placeholder="无题..."
            class="writing-title-input"
            autocomplete="off"
          />
          <p class="mt-2 text-xs tracking-wide text-slate-400">编辑中</p>
        </div>

        <div class="border-t border-slate-200/40">
          <button
            type="button"
            class="flex w-full items-center gap-2 px-8 py-3 text-left text-xs font-medium uppercase tracking-widest text-slate-400 transition hover:text-slate-600"
            @click="showEditMeta = !showEditMeta"
          >
            <svg class="h-3.5 w-3.5 transition-transform" :class="showEditMeta ? 'rotate-90' : ''" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="m9 18 6-6-6-6" /></svg>
            文章属性
          </button>
          <div v-show="showEditMeta" class="meta-drawer">
            <div class="meta-grid">
              <div class="meta-field">
                <label class="meta-label">分类</label>
                <input v-model="editPost.category" type="text" placeholder="技术 / 随笔 / 教程" class="meta-input" />
              </div>
              <div class="meta-field">
                <label class="meta-label">标签</label>
                <input v-model="editPost.tags" type="text" placeholder="用逗号分隔" class="meta-input" />
              </div>
              <div class="meta-field">
                <label class="meta-label">摘要</label>
                <input v-model="editPost.excerpt" type="text" placeholder="简短描述文章内容" class="meta-input" />
              </div>
              <div class="meta-field">
                <label class="meta-label">封面图片</label>
                <input v-model="editPost.hero_image_url" type="text" placeholder="https://..." class="meta-input" />
              </div>
            </div>
            <div v-if="editPost.hero_image_url" class="mx-8 mb-4 overflow-hidden rounded-xl border border-slate-200/50">
              <img :src="editPost.hero_image_url" alt="" class="h-28 w-full object-cover" @error="($event.target as HTMLImageElement).style.display='none'" />
            </div>
          </div>
        </div>

        <div class="writing-editor-zone">
          <textarea
            v-model="editPost.content_markdown"
            placeholder="在此撰写你的想法..."
            class="writing-editor"
          />
        </div>

        <div class="flex items-center justify-between border-t border-slate-200/40 px-8 py-4">
          <div class="flex items-center gap-2 text-xs text-slate-400">
            <span class="inline-block h-1.5 w-1.5 rounded-full" :class="editPost.content_markdown.length > 0 ? 'bg-miku' : 'bg-slate-300'" />
            {{ charCount(editPost.content_markdown) }}
          </div>
          <div class="flex items-center gap-3">
            <button type="button" class="rounded-xl px-4 py-2 text-sm text-slate-400 transition hover:bg-slate-100/50 hover:text-slate-600" @click="closeEditForm">取消</button>
            <MikuButton type="submit" variant="solid" :disabled="editing">{{ editing ? '保存中...' : '保存修改' }}</MikuButton>
          </div>
        </div>
      </form>
    </LiquidGlassCard>

    <!-- ===== Posts Table ===== -->
    <LiquidGlassCard padding="0px">
      <div v-if="loading" class="flex items-center justify-center py-20">
        <div class="loading-dot" /><div class="loading-dot delay-1" /><div class="loading-dot delay-2" />
      </div>
      <div v-else-if="posts.length === 0" class="py-20 text-center">
        <p class="text-base text-slate-400">尚无文章</p>
        <button type="button" class="mt-3 text-sm text-miku/80 transition hover:text-miku" @click="toggleCreateForm">撰写第一篇 &rarr;</button>
      </div>
      <div v-else class="overflow-x-auto">
        <table class="w-full text-left text-sm">
          <thead>
            <tr class="border-b border-slate-200/60">
              <th class="px-5 py-3.5 font-semibold text-slate-700">标题</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700">分类</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700">状态</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700 text-right">浏览量</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700">发布时间/计划时间</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700 text-center">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="post in posts"
              :key="post.id"
              class="group border-b border-slate-100/60 transition-colors hover:bg-white/40"
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
              <td class="px-5 py-3.5 text-slate-600">{{ post.displayTime }}</td>
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
                    class="rounded-xl border border-[#e9d5ff]/80 bg-white/50 px-2.5 py-1 text-xs text-[#9333ea] transition hover:bg-[#faf5ff]"
                    aria-label="定时发布文章"
                    @click="schedulePost(post.id)"
                  >
                    定时
                  </button>
                  <button
                    v-if="post.status !== 'draft'"
                    type="button"
                    class="rounded-xl border border-slate-200/80 bg-white/50 px-2.5 py-1 text-xs text-slate-600 transition hover:bg-slate-50"
                    aria-label="转为草稿"
                    @click="unpublishPost(post.id)"
                  >
                    转草稿
                  </button>
                  <button
                    type="button"
                    class="rounded-xl border border-red-200/80 bg-white/50 px-2.5 py-1 text-xs text-red-600 transition hover:bg-red-50 opacity-0 group-hover:opacity-100"
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
  scheduled_at?: string
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
  scheduledAt: string
  displayTime: string
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
  status: 'draft' | 'published' | 'scheduled'
  scheduled_at: string
}

function mapStatus(s: string): 'published' | 'draft' | 'scheduled' {
  if (s === 'published') return 'published'
  if (s === 'scheduled') return 'scheduled'
  return 'draft'
}

function formatDate(iso?: string): string {
  if (!iso) return '--'
  try {
    const d = new Date(iso)
    return d.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      hour12: false,
    })
  } catch {
    return iso
  }
}

function formatDateInputLocal(date: Date): string {
  const local = new Date(date.getTime() - date.getTimezoneOffset() * 60 * 1000)
  return local.toISOString().slice(0, 16)
}

function localInputToRFC3339(value: string): string {
  const date = new Date(value)
  return date.toISOString()
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
    scheduled_at: '',
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
  const publishedAt = formatDate(item.published_at)
  const scheduledAt = formatDate(item.scheduled_at)
  return {
    id: item.id,
    title: item.title,
    category: item.category || '--',
    status: mapStatus(item.status),
    views: Number(item.view_count) || 0,
    publishedAt,
    scheduledAt,
    displayTime: item.status === 'scheduled' ? scheduledAt : publishedAt,
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
const showCreateMeta = ref(false)
const showEditMeta = ref(false)

function charCount(text: string): string {
  const len = text.length
  if (len === 0) return '0 字'
  if (len >= 10000) return `${(len / 10000).toFixed(1)} 万字`
  return `${len} 字`
}

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
  if (newPost.value.status === 'scheduled' && !newPost.value.scheduled_at) {
    showToast('请选择定时发布时间', 'error')
    return
  }
  creating.value = true
  try {
    await api.post('/admin/posts', {
      ...toPostPayload(newPost.value),
      status: newPost.value.status,
      scheduled_at: newPost.value.status === 'scheduled'
        ? localInputToRFC3339(newPost.value.scheduled_at)
        : undefined,
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
    await loadPosts()
    showToast('文章发布成功', 'success')
  } catch (err) {
    const msg = err instanceof ApiError ? err.message : '发布失败，请稍后重试'
    console.error('[AdminPosts] publishPost failed:', err)
    showToast(msg, 'error')
  }
}

async function schedulePost(id: string) {
  const defaultTime = formatDateInputLocal(new Date(Date.now() + 30 * 60 * 1000))
  const next = window.prompt('请输入计划发布时间（格式：YYYY-MM-DDTHH:mm）', defaultTime)
  if (!next) return

  try {
    await api.post(`/admin/posts/${id}/schedule`, {
      scheduled_at: localInputToRFC3339(next),
    })
    await loadPosts()
    showToast('文章已设为定时发布', 'success')
  } catch (err) {
    const msg = err instanceof ApiError ? err.message : '定时发布设置失败，请稍后重试'
    console.error('[AdminPosts] schedulePost failed:', err)
    showToast(msg, 'error')
  }
}

async function unpublishPost(id: string) {
  try {
    await api.post(`/admin/posts/${id}/unpublish`)
    await loadPosts()
    showToast('文章已转为草稿', 'success')
  } catch (err) {
    const msg = err instanceof ApiError ? err.message : '转草稿失败，请稍后重试'
    console.error('[AdminPosts] unpublishPost failed:', err)
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
  return 'bg-[#f3e8ff] text-[#9333ea]'
}

function statusLabel(status: Post['status']) {
  if (status === 'published') return '已发布'
  if (status === 'draft') return '草稿'
  return '定时发布'
}
</script>

<style scoped>
/* ---- Writing Studio: Title Zone ---- */
.writing-title-zone {
  padding: 40px 32px 24px;
}

.writing-title-input {
  width: 100%;
  border: none;
  outline: none;
  background: transparent;
  font-family: 'Noto Serif SC', 'Noto Serif JP', 'Songti SC', serif;
  font-size: 1.75rem;
  font-weight: 700;
  line-height: 1.3;
  color: #0f172a;
  caret-color: rgb(57, 197, 187);
}

.writing-title-input::placeholder {
  color: #cbd5e1;
  font-weight: 600;
}

/* ---- Metadata Drawer ---- */
.meta-drawer {
  padding-bottom: 4px;
}

.meta-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px 20px;
  padding: 0 32px 16px;
}

.meta-field {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.meta-label {
  font-size: 0.6875rem;
  font-weight: 600;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  color: #94a3b8;
}

.meta-input {
  border-radius: 10px;
  border: 1px solid rgba(203, 213, 225, 0.5);
  background: rgba(255, 255, 255, 0.4);
  padding: 7px 12px;
  font-size: 0.8125rem;
  color: #1e293b;
  outline: none;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.meta-input:focus {
  border-color: rgba(57, 197, 187, 0.4);
  box-shadow: 0 0 0 2px rgba(57, 197, 187, 0.08);
}

.meta-input::placeholder {
  color: #b0bec5;
}

/* ---- Editor Zone ---- */
.writing-editor-zone {
  border-top: 1px solid rgba(226, 232, 240, 0.4);
  padding: 0;
}

.writing-editor {
  width: 100%;
  min-height: 420px;
  resize: vertical;
  border: none;
  outline: none;
  background: transparent;
  padding: 28px 32px;
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
  font-size: 0.8125rem;
  line-height: 1.9;
  color: #334155;
  caret-color: rgb(57, 197, 187);
  tab-size: 2;
}

.writing-editor::placeholder {
  color: #cbd5e1;
}

.writing-editor::selection {
  background: rgba(57, 197, 187, 0.18);
}

/* ---- Loading Dots ---- */
.loading-dot {
  width: 6px;
  height: 6px;
  margin: 0 4px;
  border-radius: 50%;
  background: rgb(57, 197, 187);
  opacity: 0.35;
  animation: dot-pulse 1.2s ease-in-out infinite;
}

.loading-dot.delay-1 { animation-delay: 0.2s; }
.loading-dot.delay-2 { animation-delay: 0.4s; }

@keyframes dot-pulse {
  0%, 80%, 100% { opacity: 0.2; transform: scale(0.8); }
  40% { opacity: 1; transform: scale(1.1); }
}

/* Small screen: stack meta grid */
@media (max-width: 640px) {
  .meta-grid {
    grid-template-columns: 1fr;
  }
  .writing-title-zone {
    padding: 28px 20px 16px;
  }
  .writing-editor {
    padding: 20px;
  }
}
</style>
