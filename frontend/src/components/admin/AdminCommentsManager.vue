<template>
  <section class="space-y-5">
    <LiquidGlassCard padding="24px">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 class="text-2xl font-semibold text-slate-900">评论审核</h1>
          <p class="mt-1 text-sm text-slate-600">管理文章评论与留言板留言的审核队列。</p>
        </div>
        <div class="flex flex-wrap items-center gap-2">
          <div class="flex items-center gap-1 rounded-xl border border-slate-200/80 bg-white/60 p-1">
            <button
              v-for="tab in sourceTabs"
              :key="tab.value"
              type="button"
              class="rounded-lg px-3 py-1.5 text-xs font-medium transition"
              :class="activeSource === tab.value ? 'bg-miku text-white shadow-sm' : 'text-slate-500 hover:text-slate-700'"
              @click="changeSource(tab.value)"
            >
              {{ tab.label }}
            </button>
          </div>
          <span class="rounded-full bg-amber-100 px-3 py-1 text-xs font-medium text-amber-700">{{ pendingCount }} 条待审核</span>
        </div>
      </div>
    </LiquidGlassCard>

    <LiquidGlassCard padding="0px">
      <div class="overflow-x-auto">
        <table class="w-full text-left text-sm">
          <thead>
            <tr class="border-b border-slate-200/60">
              <th class="px-5 py-3.5 font-semibold text-slate-700">评论者</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700">内容</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700">来源</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700">状态</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700">提交时间</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700 text-center">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="comment in comments"
              :key="comment.id"
              class="border-b border-slate-100/60 transition hover:bg-white/40"
            >
              <td class="px-5 py-3.5">
                <div>
                  <p class="font-medium text-slate-900">{{ comment.author }}</p>
                  <p class="text-xs text-slate-500">{{ comment.secondary }}</p>
                </div>
              </td>
              <td class="max-w-xs truncate px-5 py-3.5 text-slate-700">{{ comment.content }}</td>
              <td class="px-5 py-3.5 text-slate-600">
                <p>{{ comment.context }}</p>
                <p v-if="comment.contextHint" class="text-xs text-slate-400">{{ comment.contextHint }}</p>
              </td>
              <td class="px-5 py-3.5">
                <span
                  class="inline-block rounded-full px-2.5 py-0.5 text-xs font-medium"
                  :class="statusClass(comment.status)"
                >
                  {{ statusLabel(comment.status) }}
                </span>
              </td>
              <td class="px-5 py-3.5 text-slate-600">{{ comment.createdAt }}</td>
              <td class="px-5 py-3.5 text-center">
                <div class="flex items-center justify-center gap-2">
                  <button
                    v-if="comment.status === 'pending'"
                    type="button"
                    class="rounded-xl border border-emerald-200/80 bg-white/50 px-2.5 py-1 text-xs text-emerald-600 transition hover:bg-emerald-50"
                    aria-label="通过评论"
                    @click="approveComment(comment)"
                  >
                    通过
                  </button>
                  <button
                    v-if="comment.status === 'pending'"
                    type="button"
                    class="rounded-xl border border-red-200/80 bg-white/50 px-2.5 py-1 text-xs text-red-600 transition hover:bg-red-50"
                    aria-label="驳回评论"
                    @click="rejectComment(comment)"
                  >
                    驳回
                  </button>
                  <button
                    v-if="comment.status !== 'pending'"
                    type="button"
                    class="rounded-xl border border-slate-200/80 bg-white/50 px-2.5 py-1 text-xs text-slate-600 transition hover:bg-slate-50"
                    aria-label="删除评论"
                    @click="deleteComment(comment)"
                  >
                    删除
                  </button>
                </div>
              </td>
            </tr>
            <tr v-if="comments.length === 0 && !loading">
              <td colspan="6" class="px-5 py-8 text-center text-sm text-slate-500">暂无数据</td>
            </tr>
          </tbody>
        </table>
      </div>
      <div class="flex flex-wrap items-center justify-between gap-2 border-t border-slate-200/60 px-5 py-3 text-xs text-slate-500">
        <span>第 {{ page }} / {{ totalPages }} 页，共 {{ total }} 条</span>
        <div class="flex items-center gap-2">
          <button
            type="button"
            class="rounded-lg border border-slate-200/80 bg-white/60 px-2.5 py-1 transition hover:bg-white disabled:cursor-not-allowed disabled:opacity-50"
            :disabled="loading || page <= 1"
            @click="changePage(page - 1)"
          >
            上一页
          </button>
          <button
            type="button"
            class="rounded-lg border border-slate-200/80 bg-white/60 px-2.5 py-1 transition hover:bg-white disabled:cursor-not-allowed disabled:opacity-50"
            :disabled="loading || page >= totalPages"
            @click="changePage(page + 1)"
          >
            下一页
          </button>
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

type SourceType = 'post' | 'guestbook'
type CommentStatus = 'pending' | 'approved' | 'rejected'

interface ApiPostComment {
  id: string
  post_id: string
  post_title: string
  author_name: string
  author_email: string
  content: string
  status: string
  created_at: string
}

interface ApiGuestbookMessage {
  id: string
  parent_id?: string
  parent_author_name?: string
  author_name: string
  author_website: string
  content: string
  status: string
  vote_score: number
  created_at: string
}

interface ModerationItem {
  id: string
  source: SourceType
  author: string
  secondary: string
  content: string
  context: string
  contextHint: string
  status: CommentStatus
  createdAt: string
}

const sourceTabs: Array<{ label: string; value: SourceType }> = [
  { label: '文章评论', value: 'post' },
  { label: '留言板留言', value: 'guestbook' },
]

const activeSource = ref<SourceType>('post')
const comments = ref<ModerationItem[]>([])
const loading = ref(false)
const page = ref(1)
const pageSize = 20
const total = ref(0)

function mapStatus(s: string): CommentStatus {
  if (s === 'approved') return 'approved'
  if (s === 'rejected') return 'rejected'
  return 'pending'
}

function formatDate(iso: string): string {
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

function mapPostComment(item: ApiPostComment): ModerationItem {
  return {
    id: item.id,
    source: 'post',
    author: item.author_name || '匿名用户',
    secondary: item.author_email || '--',
    content: item.content,
    context: item.post_title || '--',
    contextHint: '',
    status: mapStatus(item.status),
    createdAt: formatDate(item.created_at),
  }
}

function mapGuestbook(item: ApiGuestbookMessage): ModerationItem {
  const context = item.parent_id ? '留言板回复' : '留言板主留言'
  const contextHint = item.parent_id
    ? `回复 ${item.parent_author_name ? `@${item.parent_author_name}` : '上一条留言'}`
    : `热度分 ${Number(item.vote_score || 0)}`

  return {
    id: item.id,
    source: 'guestbook',
    author: item.author_name || '匿名用户',
    secondary: item.author_website || '--',
    content: item.content,
    context,
    contextHint,
    status: mapStatus(item.status),
    createdAt: formatDate(item.created_at),
  }
}

async function loadComments() {
  loading.value = true
  try {
    if (activeSource.value === 'guestbook') {
      const data = await api.get<PagedData<ApiGuestbookMessage>>(`/admin/guestbook/messages?page=${page.value}&size=${pageSize}`)
      comments.value = (data.items || []).map(mapGuestbook)
      total.value = Number(data.total || 0)
      return
    }

    const data = await api.get<PagedData<ApiPostComment>>(`/admin/comments?page=${page.value}&size=${pageSize}`)
    comments.value = (data.items || []).map(mapPostComment)
    total.value = Number(data.total || 0)
  } catch (err) {
    console.error('[AdminComments] loadComments failed:', err)
    if (err instanceof ApiError && err.status === 404 && activeSource.value === 'guestbook') {
      showToast('后端缺少留言审核接口，请重启并更新 backend 服务', 'error')
    } else {
      showToast('加载评论列表失败', 'error')
    }
    comments.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

function changeSource(source: SourceType) {
  if (activeSource.value === source) return
  activeSource.value = source
  page.value = 1
  loadComments()
}

function changePage(nextPage: number) {
  if (nextPage < 1 || nextPage > totalPages.value || nextPage === page.value) return
  page.value = nextPage
  loadComments()
}

async function approveComment(item: ModerationItem) {
  const endpoint = item.source === 'guestbook'
    ? `/admin/guestbook/messages/${item.id}/approve`
    : `/admin/comments/${item.id}/approve`

  try {
    await api.post(endpoint)
    comments.value = comments.value.map((c) => c.id === item.id ? { ...c, status: 'approved' as const } : c)
    await loadComments()
    showToast('评论已通过', 'success')
  } catch (err) {
    const msg = err instanceof ApiError ? err.message : '操作失败'
    console.error('[AdminComments] approveComment failed:', err)
    showToast(msg, 'error')
  }
}

async function rejectComment(item: ModerationItem) {
  const endpoint = item.source === 'guestbook'
    ? `/admin/guestbook/messages/${item.id}/reject`
    : `/admin/comments/${item.id}/reject`

  try {
    await api.post(endpoint)
    comments.value = comments.value.map((c) => c.id === item.id ? { ...c, status: 'rejected' as const } : c)
    await loadComments()
    showToast('评论已驳回', 'success')
  } catch (err) {
    const msg = err instanceof ApiError ? err.message : '操作失败'
    console.error('[AdminComments] rejectComment failed:', err)
    showToast(msg, 'error')
  }
}

async function deleteComment(item: ModerationItem) {
  const endpoint = item.source === 'guestbook'
    ? `/admin/guestbook/messages/${item.id}`
    : `/admin/comments/${item.id}`

  try {
    await api.delete(endpoint)
    comments.value = comments.value.filter((c) => c.id !== item.id)
    if (comments.value.length === 0 && page.value > 1) {
      page.value -= 1
    }
    await loadComments()
    showToast('评论已删除', 'success')
  } catch (err) {
    const msg = err instanceof ApiError ? err.message : '删除失败'
    console.error('[AdminComments] deleteComment failed:', err)
    showToast(msg, 'error')
  }
}

onMounted(() => {
  loadComments()
})

const pendingCount = computed(() => comments.value.filter((c) => c.status === 'pending').length)
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize)))

function statusClass(status: CommentStatus) {
  if (status === 'approved') return 'bg-emerald-100 text-emerald-700'
  if (status === 'rejected') return 'bg-red-100 text-red-600'
  return 'bg-amber-100 text-amber-700'
}

function statusLabel(status: CommentStatus) {
  if (status === 'approved') return '已通过'
  if (status === 'rejected') return '已驳回'
  return '待审核'
}
</script>
