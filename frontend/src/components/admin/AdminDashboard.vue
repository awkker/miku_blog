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
        <ul class="mt-4 space-y-3 text-sm text-slate-800">
          <li class="rounded-2xl border border-white/20 bg-white/10 px-3 py-2">新文章草稿《Astro + Vue Islands 实战》待发布</li>
          <li class="rounded-2xl border border-white/20 bg-white/10 px-3 py-2">评论审核队列中有 3 条待处理内容</li>
          <li class="rounded-2xl border border-white/20 bg-white/10 px-3 py-2">友链申请箱新增 1 条请求</li>
        </ul>
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
const status = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const mounted = ref(false)

const viewLabels = ref<string[]>([])
const viewSeries = ref<ChartSeries[]>([])
const likeLabels = ref<string[]>([])
const likeSeries = ref<ChartSeries[]>([])
const chartLoading = ref(true)

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

  await Promise.all([loadDashboard(), loadTrends()])
})
</script>
