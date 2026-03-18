<template>
  <section class="space-y-5">
    <!-- Header -->
    <LiquidGlassCard padding="24px">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 class="text-2xl font-semibold text-slate-900">说说管理</h1>
          <p class="mt-1 text-sm text-slate-600">发布和管理说说动态。</p>
        </div>
        <MikuButton variant="solid" aria-label="发布说说" @click="toggleCreateForm">+ 发布说说</MikuButton>
      </div>
    </LiquidGlassCard>

    <!-- ===== Compose Card (Create) ===== -->
    <LiquidGlassCard v-if="showCreateForm" padding="0px">
      <form @submit.prevent="createMoment">
        <!-- Content zone -->
        <div class="compose-content-zone">
          <textarea
            v-model="newMoment.content"
            placeholder="此刻的想法..."
            class="compose-editor"
          />
        </div>

        <!-- Image preview -->
        <div v-if="createImagePreviews.length > 0" class="flex gap-2 px-7 pb-3">
          <div v-for="(url, idx) in createImagePreviews" :key="idx" class="relative h-20 w-20 shrink-0 overflow-hidden rounded-xl border border-slate-200/50">
            <img :src="url" alt="" class="h-full w-full object-cover" @error="($event.target as HTMLImageElement).style.display='none'" />
          </div>
        </div>

        <!-- Metadata drawer -->
        <div class="border-t border-slate-200/40">
          <button
            type="button"
            class="flex w-full items-center gap-2 px-7 py-2.5 text-left text-xs font-medium uppercase tracking-widest text-slate-400 transition hover:text-slate-600"
            @click="showCreateMeta = !showCreateMeta"
          >
            <svg class="h-3.5 w-3.5 transition-transform" :class="showCreateMeta ? 'rotate-90' : ''" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="m9 18 6-6-6-6" /></svg>
            发布设置
          </button>
          <div v-show="showCreateMeta" class="compose-meta-drawer">
            <div class="compose-meta-row">
              <label class="compose-meta-label">署名</label>
              <input v-model="newMoment.author_name" type="text" placeholder="你的昵称" class="compose-meta-input" />
            </div>
            <div class="compose-meta-row">
              <label class="compose-meta-label">图片</label>
              <input v-model="newMoment.image_urls" type="text" placeholder="贴入图片链接, 用逗号分隔 (最多4张)" class="compose-meta-input" />
            </div>
          </div>
        </div>

        <!-- Bottom bar -->
        <div class="flex items-center justify-between border-t border-slate-200/40 px-7 py-3.5">
          <div class="flex items-center gap-2 text-xs text-slate-400">
            <span class="inline-block h-1.5 w-1.5 rounded-full" :class="newMoment.content.length > 0 ? 'bg-miku' : 'bg-slate-300'" />
            {{ charCount(newMoment.content) }}
          </div>
          <div class="flex items-center gap-3">
            <button type="button" class="rounded-xl px-4 py-2 text-sm text-slate-400 transition hover:bg-slate-100/50 hover:text-slate-600" @click="closeCreateForm">取消</button>
            <MikuButton type="submit" variant="solid" :disabled="creating">{{ creating ? '发布中...' : '发布说说' }}</MikuButton>
          </div>
        </div>
      </form>
    </LiquidGlassCard>

    <!-- ===== Compose Card (Edit) ===== -->
    <LiquidGlassCard v-if="showEditForm" padding="0px">
      <form @submit.prevent="updateMoment">
        <div class="compose-content-zone">
          <p class="mb-1 text-xs tracking-wide text-slate-400">编辑中</p>
          <textarea
            v-model="editMoment.content"
            placeholder="此刻的想法..."
            class="compose-editor"
          />
        </div>

        <div v-if="editImagePreviews.length > 0" class="flex gap-2 px-7 pb-3">
          <div v-for="(url, idx) in editImagePreviews" :key="idx" class="relative h-20 w-20 shrink-0 overflow-hidden rounded-xl border border-slate-200/50">
            <img :src="url" alt="" class="h-full w-full object-cover" @error="($event.target as HTMLImageElement).style.display='none'" />
          </div>
        </div>

        <div class="border-t border-slate-200/40">
          <button
            type="button"
            class="flex w-full items-center gap-2 px-7 py-2.5 text-left text-xs font-medium uppercase tracking-widest text-slate-400 transition hover:text-slate-600"
            @click="showEditMeta = !showEditMeta"
          >
            <svg class="h-3.5 w-3.5 transition-transform" :class="showEditMeta ? 'rotate-90' : ''" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="m9 18 6-6-6-6" /></svg>
            发布设置
          </button>
          <div v-show="showEditMeta" class="compose-meta-drawer">
            <div class="compose-meta-row">
              <label class="compose-meta-label">署名</label>
              <input v-model="editMoment.author_name" type="text" placeholder="你的昵称" class="compose-meta-input" />
            </div>
            <div class="compose-meta-row">
              <label class="compose-meta-label">图片</label>
              <input v-model="editMoment.image_urls" type="text" placeholder="贴入图片链接, 用逗号分隔 (最多4张)" class="compose-meta-input" />
            </div>
          </div>
        </div>

        <div class="flex items-center justify-between border-t border-slate-200/40 px-7 py-3.5">
          <div class="flex items-center gap-2 text-xs text-slate-400">
            <span class="inline-block h-1.5 w-1.5 rounded-full" :class="editMoment.content.length > 0 ? 'bg-miku' : 'bg-slate-300'" />
            {{ charCount(editMoment.content) }}
          </div>
          <div class="flex items-center gap-3">
            <button type="button" class="rounded-xl px-4 py-2 text-sm text-slate-400 transition hover:bg-slate-100/50 hover:text-slate-600" @click="closeEditForm">取消</button>
            <MikuButton type="submit" variant="solid" :disabled="editing">{{ editing ? '保存中...' : '保存修改' }}</MikuButton>
          </div>
        </div>
      </form>
    </LiquidGlassCard>

    <!-- Stats -->
    <div class="grid gap-4 sm:grid-cols-3">
      <LiquidGlassCard padding="16px">
        <p class="text-xs uppercase tracking-[0.2em] text-slate-600">说说总数</p>
        <p class="mt-1 font-mono text-2xl font-semibold text-slate-900">{{ momentsList.length }}</p>
      </LiquidGlassCard>
      <LiquidGlassCard padding="16px">
        <p class="text-xs uppercase tracking-[0.2em] text-slate-600">总点赞</p>
        <p class="mt-1 font-mono text-2xl font-semibold text-miku">{{ totalLikes }}</p>
      </LiquidGlassCard>
      <LiquidGlassCard padding="16px">
        <p class="text-xs uppercase tracking-[0.2em] text-slate-600">总评论</p>
        <p class="mt-1 font-mono text-2xl font-semibold text-[#c084fc]">{{ totalComments }}</p>
      </LiquidGlassCard>
    </div>

    <!-- List -->
    <LiquidGlassCard padding="0px">
      <div v-if="loading" class="flex items-center justify-center py-20">
        <div class="loading-dot" /><div class="loading-dot delay-1" /><div class="loading-dot delay-2" />
      </div>
      <div v-else-if="momentsList.length === 0" class="py-20 text-center">
        <p class="text-base text-slate-400">暂无说说</p>
        <button type="button" class="mt-3 text-sm text-miku/80 transition hover:text-miku" @click="toggleCreateForm">发布第一条 &rarr;</button>
      </div>
      <div v-else class="divide-y divide-slate-100/60">
        <div
          v-for="item in momentsList"
          :key="item.id"
          class="group px-5 py-4 transition-colors hover:bg-white/40"
        >
          <div class="flex items-start justify-between gap-3">
            <div class="min-w-0 flex-1">
              <div class="flex items-center gap-2">
                <span class="text-sm font-medium text-slate-900">{{ item.author }}</span>
                <span class="text-xs text-slate-400">{{ item.createdAt }}</span>
              </div>
              <p class="mt-1.5 whitespace-pre-wrap text-sm leading-relaxed text-slate-700">{{ item.content }}</p>
              <div v-if="item.images.length > 0" class="mt-2.5 flex gap-2">
                <img
                  v-for="(img, idx) in item.images.slice(0, 4)"
                  :key="idx"
                  :src="img"
                  :alt="`图片 ${idx + 1}`"
                  class="h-16 w-16 rounded-xl border border-slate-100/80 object-cover transition-transform hover:scale-105"
                  loading="lazy"
                />
              </div>
              <div class="mt-2 flex items-center gap-4 text-xs text-slate-400">
                <span>{{ item.likes }} ❤</span>
                <span>{{ item.reposts }} ↻</span>
                <span>{{ item.comments }} ✉</span>
              </div>
            </div>
            <div class="shrink-0 opacity-0 transition group-hover:opacity-100">
              <button
                type="button"
                class="rounded-xl border border-slate-200/80 bg-white/50 px-2.5 py-1 text-xs text-slate-700 transition hover:border-miku/40 hover:text-miku"
                aria-label="编辑说说"
                @click="startEditMoment(item)"
              >
                编辑
              </button>
            </div>
          </div>
        </div>
      </div>
    </LiquidGlassCard>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

import { api, ApiError, type PagedData } from '../../lib/api'
import { showToast } from '../../stores/ui'
import LiquidGlassCard from '../ui/LiquidGlassCard.vue'
import MikuButton from '../ui/MikuButton.vue'

interface ApiMoment {
  id: string
  author_name: string
  content: string
  image_urls: string[]
  like_count: number
  repost_count: number
  comment_count: number
  created_at: string
}

interface MomentItem {
  id: string
  author: string
  content: string
  images: string[]
  likes: number
  reposts: number
  comments: number
  createdAt: string
}

interface MomentForm {
  author_name: string
  content: string
  image_urls: string
}

function formatDate(iso: string): string {
  try {
    const d = new Date(iso)
    return d.toLocaleString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', hour12: false })
  } catch {
    return iso
  }
}

function mapMoment(item: ApiMoment): MomentItem {
  return {
    id: item.id,
    author: item.author_name,
    content: item.content,
    images: item.image_urls || [],
    likes: Number(item.like_count) || 0,
    reposts: Number(item.repost_count) || 0,
    comments: Number(item.comment_count) || 0,
    createdAt: formatDate(item.created_at),
  }
}

function createEmptyMomentForm(): MomentForm {
  return {
    author_name: '',
    content: '',
    image_urls: '',
  }
}

function toImageURLs(input: string): string[] {
  return input.split(',').map((u: string) => u.trim()).filter(Boolean)
}

const momentsList = ref<MomentItem[]>([])
const loading = ref(false)
const showCreateForm = ref(false)
const creating = ref(false)
const showEditForm = ref(false)
const editing = ref(false)
const editingMomentID = ref<string | null>(null)

const newMoment = ref<MomentForm>(createEmptyMomentForm())
const editMoment = ref<MomentForm>(createEmptyMomentForm())
const showCreateMeta = ref(false)
const showEditMeta = ref(false)

const createImagePreviews = computed(() => {
  return newMoment.value.image_urls.split(',').map((u: string) => u.trim()).filter(Boolean).slice(0, 4)
})
const editImagePreviews = computed(() => {
  return editMoment.value.image_urls.split(',').map((u: string) => u.trim()).filter(Boolean).slice(0, 4)
})

function charCount(text: string): string {
  const len = text.length
  if (len === 0) return '0 字'
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
  editingMomentID.value = null
  editMoment.value = createEmptyMomentForm()
}

const totalLikes = computed(() => momentsList.value.reduce((sum, m) => sum + m.likes, 0))
const totalComments = computed(() => momentsList.value.reduce((sum, m) => sum + m.comments, 0))

async function loadMoments() {
  loading.value = true
  try {
    const data = await api.get<PagedData<ApiMoment>>('/moments?size=50')
    momentsList.value = (data.items || []).map(mapMoment)
  } catch (err) {
    console.error('[AdminMoments] loadMoments failed:', err)
    showToast('加载说说列表失败', 'error')
    momentsList.value = []
  } finally {
    loading.value = false
  }
}

async function createMoment() {
  if (!newMoment.value.author_name.trim() || !newMoment.value.content.trim()) return
  creating.value = true
  try {
    await api.post('/moments', {
      author_name: newMoment.value.author_name.trim(),
      content: newMoment.value.content.trim(),
      image_urls: toImageURLs(newMoment.value.image_urls),
    })
    closeCreateForm()
    newMoment.value = createEmptyMomentForm()
    showToast('说说发布成功', 'success')
    await loadMoments()
  } catch (err) {
    const msg = err instanceof ApiError ? err.message : '发布说说失败，请稍后重试'
    console.error('[AdminMoments] createMoment failed:', err)
    showToast(msg, 'error')
  } finally {
    creating.value = false
  }
}

function startEditMoment(item: MomentItem) {
  editingMomentID.value = item.id
  editMoment.value = {
    author_name: item.author,
    content: item.content,
    image_urls: (item.images || []).join(', '),
  }
  showEditForm.value = true
  showCreateForm.value = false
}

async function updateMoment() {
  if (!editingMomentID.value) return
  if (!editMoment.value.author_name.trim() || !editMoment.value.content.trim()) return
  editing.value = true
  try {
    await api.put(`/admin/moments/${editingMomentID.value}`, {
      author_name: editMoment.value.author_name.trim(),
      content: editMoment.value.content.trim(),
      image_urls: toImageURLs(editMoment.value.image_urls),
    })
    showToast('说说更新成功', 'success')
    closeEditForm()
    await loadMoments()
  } catch (err) {
    const msg = err instanceof ApiError ? err.message : '更新说说失败，请稍后重试'
    console.error('[AdminMoments] updateMoment failed:', err)
    showToast(msg, 'error')
  } finally {
    editing.value = false
  }
}

onMounted(() => {
  loadMoments()
})
</script>

<style scoped>
/* ---- Compose Content Zone ---- */
.compose-content-zone {
  padding: 28px 28px 12px;
}

.compose-editor {
  width: 100%;
  min-height: 140px;
  resize: vertical;
  border: none;
  outline: none;
  background: transparent;
  font-size: 0.9375rem;
  line-height: 1.75;
  color: #1e293b;
  caret-color: rgb(57, 197, 187);
}

.compose-editor::placeholder {
  color: #cbd5e1;
}

.compose-editor::selection {
  background: rgba(57, 197, 187, 0.18);
}

/* ---- Compose Metadata ---- */
.compose-meta-drawer {
  padding: 0 28px 12px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.compose-meta-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.compose-meta-label {
  flex-shrink: 0;
  width: 40px;
  font-size: 0.6875rem;
  font-weight: 600;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  color: #94a3b8;
  text-align: right;
}

.compose-meta-input {
  flex: 1;
  border-radius: 10px;
  border: 1px solid rgba(203, 213, 225, 0.5);
  background: rgba(255, 255, 255, 0.4);
  padding: 7px 12px;
  font-size: 0.8125rem;
  color: #1e293b;
  outline: none;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.compose-meta-input:focus {
  border-color: rgba(57, 197, 187, 0.4);
  box-shadow: 0 0 0 2px rgba(57, 197, 187, 0.08);
}

.compose-meta-input::placeholder {
  color: #b0bec5;
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
</style>
