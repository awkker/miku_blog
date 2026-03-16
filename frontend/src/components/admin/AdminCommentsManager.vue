<template>
  <section class="space-y-5">
    <LiquidGlassCard padding="24px">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 class="text-2xl font-semibold text-slate-900">评论审核</h1>
          <p class="mt-1 text-sm text-slate-600">管理评论队列、举报处理与审核策略。</p>
        </div>
        <div class="flex items-center gap-2">
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
              <th class="px-5 py-3.5 font-semibold text-slate-700">所属文章</th>
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
                  <p class="text-xs text-slate-500">{{ comment.email }}</p>
                </div>
              </td>
              <td class="max-w-xs truncate px-5 py-3.5 text-slate-700">{{ comment.content }}</td>
              <td class="px-5 py-3.5 text-slate-600">{{ comment.postTitle }}</td>
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
                    @click="approveComment(comment.id)"
                  >
                    通过
                  </button>
                  <button
                    v-if="comment.status === 'pending'"
                    type="button"
                    class="rounded-xl border border-red-200/80 bg-white/50 px-2.5 py-1 text-xs text-red-600 transition hover:bg-red-50"
                    aria-label="驳回评论"
                    @click="rejectComment(comment.id)"
                  >
                    驳回
                  </button>
                  <button
                    v-if="comment.status !== 'pending'"
                    type="button"
                    class="rounded-xl border border-slate-200/80 bg-white/50 px-2.5 py-1 text-xs text-slate-600 transition hover:bg-slate-50"
                    aria-label="删除评论"
                    @click="deleteComment(comment.id)"
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
import { computed, onMounted, ref } from 'vue'

import { api, ApiError, type PagedData } from '../../lib/api'
import { showToast } from '../../stores/ui'
import LiquidGlassCard from '../ui/LiquidGlassCard.vue'

interface ApiComment {
  id: string
  post_id: string
  post_title: string
  author_name: string
  author_email: string
  content: string
  status: string
  created_at: string
}

interface Comment {
  id: string
  author: string
  email: string
  content: string
  postTitle: string
  status: 'pending' | 'approved' | 'rejected'
  createdAt: string
}

function mapStatus(s: string): 'pending' | 'approved' | 'rejected' {
  if (s === 'approved') return 'approved'
  if (s === 'rejected') return 'rejected'
  return 'pending'
}

function formatDate(iso: string): string {
  try {
    const d = new Date(iso)
    return d.toLocaleString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', hour12: false })
  } catch {
    return iso
  }
}

function mapComment(item: ApiComment): Comment {
  return {
    id: item.id,
    author: item.author_name,
    email: item.author_email,
    content: item.content,
    postTitle: item.post_title || '--',
    status: mapStatus(item.status),
    createdAt: formatDate(item.created_at),
  }
}

const comments = ref<Comment[]>([])
const loading = ref(false)

async function loadComments() {
  loading.value = true
  try {
    const data = await api.get<PagedData<ApiComment>>('/admin/comments?size=50')
    comments.value = (data.items || []).map(mapComment)
  } catch (err) {
    console.error('[AdminComments] loadComments failed:', err)
    showToast('加载评论列表失败', 'error')
    comments.value = []
  } finally {
    loading.value = false
  }
}

async function approveComment(id: string) {
  try {
    await api.post(`/admin/comments/${id}/approve`)
    comments.value = comments.value.map((c) => c.id === id ? { ...c, status: 'approved' as const } : c)
    showToast('评论已通过', 'success')
  } catch (err) {
    const msg = err instanceof ApiError ? err.message : '操作失败'
    console.error('[AdminComments] approveComment failed:', err)
    showToast(msg, 'error')
  }
}

async function rejectComment(id: string) {
  try {
    await api.post(`/admin/comments/${id}/reject`)
    comments.value = comments.value.map((c) => c.id === id ? { ...c, status: 'rejected' as const } : c)
    showToast('评论已驳回', 'success')
  } catch (err) {
    const msg = err instanceof ApiError ? err.message : '操作失败'
    console.error('[AdminComments] rejectComment failed:', err)
    showToast(msg, 'error')
  }
}

async function deleteComment(id: string) {
  try {
    await api.delete(`/admin/comments/${id}`)
    comments.value = comments.value.filter((c) => c.id !== id)
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

function statusClass(status: Comment['status']) {
  if (status === 'approved') return 'bg-emerald-100 text-emerald-700'
  if (status === 'rejected') return 'bg-red-100 text-red-600'
  return 'bg-amber-100 text-amber-700'
}

function statusLabel(status: Comment['status']) {
  if (status === 'approved') return '已通过'
  if (status === 'rejected') return '已驳回'
  return '待审核'
}
</script>
