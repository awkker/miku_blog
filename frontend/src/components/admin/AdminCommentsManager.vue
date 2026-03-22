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

    <LiquidGlassCard padding="20px">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h2 class="text-lg font-semibold text-slate-900">{{ rateCopy.title }}</h2>
          <p class="mt-1 text-sm text-slate-600">{{ rateCopy.subtitle }}</p>
        </div>
        <button
          type="button"
          class="rounded-xl border border-slate-200/80 bg-white/60 px-3 py-1.5 text-xs font-medium text-slate-600 transition hover:bg-white"
          :disabled="rateLoading"
          @click="loadRateLimitMetrics"
        >
          {{ rateCopy.refresh }}
        </button>
      </div>

      <div class="mt-4 grid gap-3 sm:grid-cols-2">
        <div class="rounded-xl border border-slate-200/70 bg-white/50 px-4 py-3">
          <p class="text-xs uppercase tracking-[0.18em] text-slate-500">{{ rateCopy.totalAllowed }}</p>
          <p class="mt-1 text-2xl font-semibold text-emerald-600">{{ totalAllowed }}</p>
        </div>
        <div class="rounded-xl border border-slate-200/70 bg-white/50 px-4 py-3">
          <p class="text-xs uppercase tracking-[0.18em] text-slate-500">{{ rateCopy.totalBlocked }}</p>
          <p class="mt-1 text-2xl font-semibold text-red-500">{{ totalBlocked }}</p>
        </div>
      </div>

      <div v-if="rateLoading" class="mt-4 rounded-xl border border-slate-200/70 bg-white/40 px-4 py-8 text-center text-sm text-slate-500">
        {{ rateCopy.loading }}
      </div>
      <div v-else-if="!rateMetrics || rateMetrics.trend.length === 0" class="mt-4 rounded-xl border border-slate-200/70 bg-white/40 px-4 py-8 text-center text-sm text-slate-500">
        {{ rateCopy.noData }}
      </div>
      <div v-else class="mt-4 grid gap-4 xl:grid-cols-[1.4fr_1fr]">
        <div class="rounded-xl border border-slate-200/70 bg-white/50 p-3">
          <p class="mb-2 text-xs font-semibold uppercase tracking-[0.16em] text-slate-500">{{ rateCopy.trendTitle }}</p>
          <div class="h-64">
            <v-chart :option="rateTrendOption" autoresize />
          </div>
        </div>
        <div class="rounded-xl border border-slate-200/70 bg-white/50 p-3">
          <p class="mb-2 text-xs font-semibold uppercase tracking-[0.16em] text-slate-500">{{ rateCopy.ruleTitle }}</p>
          <div class="text-xs text-slate-600">
            <div class="grid grid-cols-[1fr_auto_auto_auto] gap-2 border-b border-slate-200/70 pb-2 font-semibold text-slate-500">
              <span>{{ rateCopy.ruleColumn }}</span>
              <span class="text-right">{{ rateCopy.allowedColumn }}</span>
              <span class="text-right">{{ rateCopy.blockedColumn }}</span>
              <span class="text-right">{{ rateCopy.totalColumn }}</span>
            </div>
            <div
              v-for="item in rateRules"
              :key="item.rule"
              class="grid grid-cols-[1fr_auto_auto_auto] gap-2 border-b border-slate-100/80 py-2"
            >
              <span class="truncate">{{ item.rule }}</span>
              <span class="text-right text-emerald-600">{{ item.allowed }}</span>
              <span class="text-right text-red-500">{{ item.blocked }}</span>
              <span class="text-right">{{ item.total }}</span>
            </div>
          </div>
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
import { use } from 'echarts/core'
import { LineChart } from 'echarts/charts'
import { GridComponent, LegendComponent, TooltipComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import { computed, onMounted, ref } from 'vue'
import VChart from 'vue-echarts'

import { api, ApiError, type PagedData } from '../../lib/api'
import { adminCopy } from '../../content/copy'
import { showToast } from '../../stores/ui'
import LiquidGlassCard from '../ui/LiquidGlassCard.vue'

use([CanvasRenderer, LineChart, GridComponent, TooltipComponent, LegendComponent])

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

interface RateLimitRuleMetric {
  rule: string
  allowed: number
  blocked: number
  total: number
}

interface RateLimitTrendPoint {
  bucket: string
  allowed: number
  blocked: number
}

interface RateLimitMetrics {
  window_minutes: number
  total_allowed: number
  total_blocked: number
  rules: RateLimitRuleMetric[]
  trend: RateLimitTrendPoint[]
}

const sourceTabs: Array<{ label: string; value: SourceType }> = [
  { label: '文章评论', value: 'post' },
  { label: '留言板留言', value: 'guestbook' },
]
const rateCopy = adminCopy.moderation.rateLimit

const activeSource = ref<SourceType>('post')
const comments = ref<ModerationItem[]>([])
const loading = ref(false)
const page = ref(1)
const pageSize = 20
const total = ref(0)
const rateMetrics = ref<RateLimitMetrics | null>(null)
const rateLoading = ref(false)

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

async function loadRateLimitMetrics() {
  rateLoading.value = true
  try {
    rateMetrics.value = await api.get<RateLimitMetrics>('/admin/moderation/rate-limit-metrics?minutes=60')
  } catch (err) {
    const msg = err instanceof ApiError ? err.message : '加载限流指标失败'
    console.error('[AdminComments] loadRateLimitMetrics failed:', err)
    showToast(msg, 'error')
    rateMetrics.value = null
  } finally {
    rateLoading.value = false
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
  loadRateLimitMetrics()
})

const pendingCount = computed(() => comments.value.filter((c) => c.status === 'pending').length)
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize)))
const totalAllowed = computed(() => Number(rateMetrics.value?.total_allowed || 0))
const totalBlocked = computed(() => Number(rateMetrics.value?.total_blocked || 0))
const rateRules = computed(() => rateMetrics.value?.rules || [])
const rateTrendOption = computed(() => {
  const trend = rateMetrics.value?.trend || []
  return {
    backgroundColor: 'transparent',
    tooltip: { trigger: 'axis' },
    legend: {
      top: 0,
      textStyle: { color: '#64748b', fontSize: 11 },
      data: [rateCopy.allowedLegend, rateCopy.blockedLegend],
    },
    grid: { top: 36, left: 40, right: 16, bottom: 30 },
    xAxis: {
      type: 'category',
      data: trend.map((item) => item.bucket.slice(-5)),
      axisLine: { lineStyle: { color: '#cbd5e1' } },
      axisLabel: { color: '#64748b', fontSize: 10 },
    },
    yAxis: {
      type: 'value',
      axisLine: { show: false },
      splitLine: { lineStyle: { color: 'rgba(148,163,184,0.2)' } },
      axisLabel: { color: '#64748b', fontSize: 10 },
    },
    series: [
      {
        name: rateCopy.allowedLegend,
        type: 'line',
        smooth: true,
        data: trend.map((item) => item.allowed),
        lineStyle: { color: '#10b981', width: 2 },
        itemStyle: { color: '#10b981' },
      },
      {
        name: rateCopy.blockedLegend,
        type: 'line',
        smooth: true,
        data: trend.map((item) => item.blocked),
        lineStyle: { color: '#ef4444', width: 2 },
        itemStyle: { color: '#ef4444' },
      },
    ],
  }
})

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
