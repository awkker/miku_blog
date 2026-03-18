<template>
  <section class="space-y-5">
    <LiquidGlassCard padding="24px">
      <h1 class="text-2xl font-semibold text-slate-900">欢迎回来，{{ userName }}</h1>
      <p class="mt-2 text-sm text-slate-700">后台已就绪，你可以继续管理文章、评论与友链。</p>
    </LiquidGlassCard>

    <div class="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
      <template v-if="status === 'loading'">
        <SkeletonCard v-for="item in 4" :key="item" />
      </template>

      <template v-else-if="status === 'error'">
        <div class="md:col-span-2 xl:col-span-4">
          <ErrorState description="统计数据读取失败，请重试。" @retry="loadDashboard" />
        </div>
      </template>

      <template v-else>
        <DashboardStatCard
          v-for="item in stats"
          :key="item.label"
          :label="item.label"
          :value="item.value"
          :trend="item.trend"
          :icon="item.icon"
        />
      </template>
    </div>

    <div class="grid gap-4 lg:grid-cols-2">
      <DashboardChart title="近 7 日浏览量 / 评论趋势" :labels="viewLabels" :series="viewSeries" :loading="chartLoading" />
      <DashboardChart title="近 7 日点赞趋势" :labels="likeLabels" :series="likeSeries" :loading="chartLoading" />
    </div>

    <div class="grid gap-4 lg:grid-cols-3">
      <LiquidGlassCard padding="20px" class="lg:col-span-2">
        <h2 class="text-lg font-semibold text-slate-900">最近动态</h2>
        <div v-if="activityLoading && status === 'loading'" class="mt-4 flex items-center justify-center py-6">
          <p class="text-sm text-slate-400">加载中...</p>
        </div>
        <div v-else class="mt-4 space-y-3 text-sm text-slate-800">
          <div v-if="rawStats && rawStats.pending_comments > 0" class="flex items-center gap-2.5 rounded-2xl border border-amber-200/60 bg-amber-50/50 px-3 py-2">
            <span class="flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-amber-400/20 text-[11px] text-amber-600">!</span>
            <span class="text-amber-800">评论审核队列中有 <strong>{{ rawStats.pending_comments }}</strong> 条待处理</span>
          </div>
          <div v-if="rawStats && rawStats.draft_count > 0" class="flex items-center gap-2.5 rounded-2xl border border-sky-200/60 bg-sky-50/50 px-3 py-2">
            <span class="flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-sky-400/20 text-[11px] text-sky-600">&#x270E;</span>
            <span class="text-sky-800">有 <strong>{{ rawStats.draft_count }}</strong> 篇草稿待发布</span>
          </div>
          <template v-if="activityItems.length > 0">
            <div v-for="item in activityItems" :key="item.id" class="rounded-2xl border border-white/20 bg-white/10 px-3 py-2">
              <div class="flex items-center justify-between gap-2">
                <span>{{ formatAction(item) }}</span>
                <span class="shrink-0 text-xs text-slate-400">{{ formatTime(item.created_at) }}</span>
              </div>
              <p v-if="item.admin_username" class="mt-0.5 text-xs text-slate-400">操作人：{{ item.admin_username }}</p>
            </div>
          </template>
          <p v-else-if="!rawStats?.pending_comments && !rawStats?.draft_count" class="py-6 text-center text-sm text-slate-400">暂无操作记录</p>
        </div>
      </LiquidGlassCard>

      <LiquidGlassCard padding="20px">
        <h2 class="text-lg font-semibold text-slate-900">快捷操作</h2>
        <div class="mt-4 grid gap-2">
          <MikuButton variant="ghost" class="!border-slate-300/80 !bg-white/65 !text-slate-900 hover:!bg-white/80" aria-label="新建文章">新建文章</MikuButton>
          <MikuButton variant="ghost" class="!border-slate-300/80 !bg-white/65 !text-slate-900 hover:!bg-white/80" aria-label="查看站点">查看站点</MikuButton>
          <MikuButton variant="ghost" class="!border-slate-300/80 !bg-white/65 !text-slate-900 hover:!bg-white/80" aria-label="管理友链">管理友链</MikuButton>
        </div>
      </LiquidGlassCard>
    </div>
  </section>
</template>

<script setup lang="ts">
import { useStore } from '@nanostores/vue'
import { computed, onMounted, ref } from 'vue'

import { api } from '../../lib/api'
import { authState, hydrateAuth } from '../../stores/auth'
import { setScopeStatus } from '../../stores/loading'
import DashboardChart from './DashboardChart.vue'
import type { ChartSeries } from './DashboardChart.vue'
import DashboardStatCard from './DashboardStatCard.vue'
import ErrorState from '../ui/ErrorState.vue'
import LiquidGlassCard from '../ui/LiquidGlassCard.vue'
import MikuButton from '../ui/MikuButton.vue'
import SkeletonCard from '../ui/SkeletonCard.vue'

interface StatItem {
  label: string
  value: string
  trend: string
  icon: 'article' | 'comment' | 'like' | 'link'
}

interface ApiStats {
  total_posts: number
  total_likes: number
  pending_comments: number
  friend_count: number
  draft_count: number
}

interface ViewTrendPoint {
  day: string
  pv: number
  uv: number
}

interface TrendPoint {
  day: string
  value: number
}

const auth = useStore(authState)

const stats = ref<StatItem[]>([])
const rawStats = ref<ApiStats | null>(null)
const status = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const mounted = ref(false)

const viewLabels = ref<string[]>([])
const viewSeries = ref<ChartSeries[]>([])
const likeLabels = ref<string[]>([])
const likeSeries = ref<ChartSeries[]>([])
const chartLoading = ref(true)

interface AuditLogItem {
  id: string
  action: string
  target_type: string
  target_id: string
  detail: unknown
  ip: string
  admin_username?: string
  created_at: string
}

const activityItems = ref<AuditLogItem[]>([])
const activityLoading = ref(true)

const userName = computed(() => {
  if (!mounted.value) {
    return '管理员'
  }

  return auth.value.user?.name ?? '管理员'
})

function fmtDay(iso: string) {
  return iso.slice(5)
}

async function loadDashboard() {
  status.value = 'loading'
  setScopeStatus('adminDashboard', 'loading')

  try {
    const data = await api.get<ApiStats>('/admin/dashboard/stats')
    rawStats.value = data
    stats.value = [
      { label: '文章总数', value: String(data.total_posts), trend: '', icon: 'article' },
      { label: '待审评论', value: String(data.pending_comments), trend: '', icon: 'comment' },
      { label: '累计点赞', value: data.total_likes.toLocaleString(), trend: '', icon: 'like' },
      { label: '友链数量', value: String(data.friend_count), trend: '', icon: 'link' },
    ]
    status.value = 'success'
    setScopeStatus('adminDashboard', 'success')
  } catch {
    status.value = 'error'
    setScopeStatus('adminDashboard', 'error')
  }
}

async function loadTrends() {
  chartLoading.value = true
  try {
    const [viewData, commentData, likeData] = await Promise.all([
      api.get<ViewTrendPoint[]>('/admin/dashboard/trend/views?days=7'),
      api.get<TrendPoint[]>('/admin/dashboard/trend/comments?days=7'),
      api.get<TrendPoint[]>('/admin/dashboard/trend/likes?days=7'),
    ])

    viewLabels.value = (viewData || []).map((p) => fmtDay(p.day))
    viewSeries.value = [
      { name: '浏览量 (PV)', data: (viewData || []).map((p) => p.pv), color: '#39c5bb' },
      { name: '评论数', data: (commentData || []).map((p) => p.value), color: '#c084fc' },
    ]

    likeLabels.value = (likeData || []).map((p) => fmtDay(p.day))
    likeSeries.value = [
      { name: '点赞数', data: (likeData || []).map((p) => p.value), color: '#39c5bb' },
    ]
  } catch {
    // silent fail for charts
  } finally {
    chartLoading.value = false
  }
}

const ACTION_MAP: Record<string, string> = {
  approve: '通过了',
  reject: '拒绝了',
  delete: '删除了',
  create: '创建了',
  update: '更新了',
  publish: '发布了',
  unpublish: '下架了',
  schedule: '定时发布了',
}

const TARGET_MAP: Record<string, string> = {
  comment: '评论',
  post: '文章',
  friend_link: '友链',
  guestbook: '留言',
  moment: '说说',
}

function formatAction(item: AuditLogItem): string {
  const action = ACTION_MAP[item.action] || item.action
  const target = TARGET_MAP[item.target_type] || item.target_type
  return `${action}一条${target}`
}

function formatTime(iso: string): string {
  try {
    const d = new Date(iso)
    const now = new Date()
    const diff = now.getTime() - d.getTime()
    const mins = Math.floor(diff / 60000)
    if (mins < 1) return '刚刚'
    if (mins < 60) return `${mins} 分钟前`
    const hours = Math.floor(mins / 60)
    if (hours < 24) return `${hours} 小时前`
    const days = Math.floor(hours / 24)
    if (days < 7) return `${days} 天前`
    return d.toLocaleString('zh-CN', { month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', hour12: false })
  } catch {
    return iso
  }
}

async function loadActivity() {
  activityLoading.value = true
  try {
    const items = await api.get<AuditLogItem[]>('/admin/audit-logs?page=1&size=8')
    activityItems.value = items || []
  } catch {
    activityItems.value = []
  } finally {
    activityLoading.value = false
  }
}

onMounted(async () => {
  mounted.value = true
  hydrateAuth()

  // Read the latest auth snapshot directly to avoid stale ref value
  // during the first mount tick.
  const latestAuth = authState.get()

  if (latestAuth.status !== 'authenticated') {
    window.location.replace('/login')
    return
  }

  await Promise.all([loadDashboard(), loadTrends(), loadActivity()])
})
</script>
