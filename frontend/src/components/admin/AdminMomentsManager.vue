<template>
  <section class="space-y-5">
    <LiquidGlassCard padding="24px">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 class="text-2xl font-semibold text-slate-900">说说管理</h1>
          <p class="mt-1 text-sm text-slate-600">发布和管理说说动态。</p>
        </div>
        <MikuButton variant="solid" aria-label="发布说说" @click="toggleCreateForm">+ 发布说说</MikuButton>
      </div>
    </LiquidGlassCard>

    <!-- Create Form -->
    <LiquidGlassCard v-if="showCreateForm" padding="24px">
      <h2 class="mb-4 text-lg font-semibold text-slate-900">发布新说说</h2>
      <form class="space-y-3" @submit.prevent="createMoment">
        <input v-model="newMoment.author_name" type="text" placeholder="作者昵称 *" class="w-full rounded-xl border border-slate-200/80 bg-white/60 px-3 py-2 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-1 focus:ring-miku/30" />
        <textarea v-model="newMoment.content" rows="4" placeholder="说说内容 *" class="w-full rounded-xl border border-slate-200/80 bg-white/60 px-3 py-2 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-1 focus:ring-miku/30" />
        <input v-model="newMoment.image_urls" type="text" placeholder="图片 URL (逗号分隔, 最多 4 张)" class="w-full rounded-xl border border-slate-200/80 bg-white/60 px-3 py-2 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-1 focus:ring-miku/30" />
        <div class="flex items-center gap-3">
          <MikuButton type="submit" variant="solid" :disabled="creating">{{ creating ? '发布中...' : '发布说说' }}</MikuButton>
          <button type="button" class="text-sm text-slate-500 hover:text-slate-700" @click="closeCreateForm">取消</button>
        </div>
      </form>
    </LiquidGlassCard>

    <LiquidGlassCard v-if="showEditForm" padding="24px">
      <h2 class="mb-4 text-lg font-semibold text-slate-900">编辑说说</h2>
      <form class="space-y-3" @submit.prevent="updateMoment">
        <input v-model="editMoment.author_name" type="text" placeholder="作者昵称 *" class="w-full rounded-xl border border-slate-200/80 bg-white/60 px-3 py-2 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-1 focus:ring-miku/30" />
        <textarea v-model="editMoment.content" rows="4" placeholder="说说内容 *" class="w-full rounded-xl border border-slate-200/80 bg-white/60 px-3 py-2 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-1 focus:ring-miku/30" />
        <input v-model="editMoment.image_urls" type="text" placeholder="图片 URL (逗号分隔, 最多 4 张)" class="w-full rounded-xl border border-slate-200/80 bg-white/60 px-3 py-2 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-1 focus:ring-miku/30" />
        <div class="flex items-center gap-3">
          <MikuButton type="submit" variant="solid" :disabled="editing">{{ editing ? '保存中...' : '保存修改' }}</MikuButton>
          <button type="button" class="text-sm text-slate-500 hover:text-slate-700" @click="closeEditForm">取消</button>
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
      <div v-if="loading" class="px-5 py-8 text-center text-sm text-slate-400">加载中...</div>
      <div v-else-if="momentsList.length === 0" class="px-5 py-8 text-center text-sm text-slate-400">暂无说说</div>
      <div v-else class="divide-y divide-slate-100/60">
        <div
          v-for="item in momentsList"
          :key="item.id"
          class="px-5 py-4 transition hover:bg-white/40"
        >
          <div class="flex items-start justify-between gap-3">
            <div class="min-w-0 flex-1">
              <div class="flex items-center gap-2">
                <span class="text-sm font-medium text-slate-900">{{ item.author }}</span>
                <span class="text-xs text-slate-400">{{ item.createdAt }}</span>
              </div>
              <p class="mt-1.5 whitespace-pre-wrap text-sm leading-relaxed text-slate-700">{{ item.content }}</p>
              <div v-if="item.images.length > 0" class="mt-2 flex gap-2">
                <img
                  v-for="(img, idx) in item.images.slice(0, 4)"
                  :key="idx"
                  :src="img"
                  :alt="`\u56fe\u7247 ${idx + 1}`"
                  class="h-16 w-16 rounded-lg border border-slate-100 object-cover"
                  loading="lazy"
                />
              </div>
              <div class="mt-2 flex items-center gap-4 text-xs text-slate-400">
                <span>{{ item.likes }} \u2764</span>
                <span>{{ item.reposts }} \u21BB</span>
                <span>{{ item.comments }} \u2709</span>
              </div>
            </div>
            <div class="shrink-0">
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
