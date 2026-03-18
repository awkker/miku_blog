<template>
  <section class="space-y-5">
    <LiquidGlassCard max-width="100%" padding="16px">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <button
          type="button"
          class="inline-flex items-center gap-2 rounded-xl border border-slate-300/80 bg-white/70 px-3 py-2 text-sm font-semibold text-slate-800 transition hover:border-miku/45 hover:text-miku"
          @click="showFilters = !showFilters"
        >
          <svg viewBox="0 0 24 24" class="h-4 w-4 fill-none stroke-current stroke-[1.8]">
            <path d="M4 6h16M7 12h10M10 18h4" />
          </svg>
          Filter
        </button>

        <div class="flex items-center gap-2">
          <button
            type="button"
            class="inline-flex h-9 w-9 items-center justify-center rounded-xl border border-slate-300/80 bg-white/70 text-slate-700 transition hover:border-miku/45 hover:text-miku"
            aria-label="上一时间窗口"
            @click="goPrevWindow"
          >
            <svg viewBox="0 0 24 24" class="h-4 w-4 fill-none stroke-current stroke-[2]">
              <path d="M15 6l-6 6 6 6" />
            </svg>
          </button>

          <button
            type="button"
            class="inline-flex h-9 w-9 items-center justify-center rounded-xl border border-slate-300/80 bg-white/70 text-slate-700 transition hover:border-miku/45 hover:text-miku disabled:cursor-not-allowed disabled:opacity-45"
            aria-label="下一时间窗口"
            :disabled="offset === 0"
            @click="goNextWindow"
          >
            <svg viewBox="0 0 24 24" class="h-4 w-4 fill-none stroke-current stroke-[2]">
              <path d="M9 6l6 6-6 6" />
            </svg>
          </button>

          <div class="rounded-xl border border-slate-300/80 bg-white/70 px-2 py-1">
            <select v-model="selectedRange" class="border-0 bg-transparent pr-6 text-sm font-semibold text-slate-800 focus:outline-none">
              <option v-for="item in rangeOptions" :key="item.key" :value="item.key">{{ item.label }}</option>
            </select>
          </div>
        </div>
      </div>

      <div v-if="showFilters" class="mt-3 grid gap-3 rounded-2xl border border-slate-200/80 bg-white/58 p-3 md:grid-cols-2">
        <label class="grid gap-1">
          <span class="text-xs font-semibold uppercase tracking-[0.18em] text-slate-500">Path Filter</span>
          <input
            v-model.trim="pathKeyword"
            type="text"
            class="rounded-xl border border-slate-300/80 bg-white/90 px-3 py-2 text-sm text-slate-800 outline-none transition focus:border-miku/60"
            placeholder="输入路径关键字"
          >
        </label>
        <div class="grid content-center rounded-xl border border-slate-200/80 bg-white/78 px-3 py-2 text-sm text-slate-600">
          <p class="font-semibold text-slate-800">{{ windowLabel }}</p>
          <p>{{ formattedWindowRange }}</p>
        </div>
      </div>
    </LiquidGlassCard>

    <div class="grid gap-4 sm:grid-cols-2 xl:grid-cols-5">
      <template v-if="status === 'loading' && !analytics">
        <SkeletonCard v-for="i in 5" :key="i" />
      </template>
      <template v-else>
        <LiquidGlassCard v-for="item in statCards" :key="item.label" max-width="100%" padding="18px">
          <p class="text-xs font-semibold uppercase tracking-[0.2em] text-slate-500">{{ item.label }}</p>
          <p class="mt-2 font-mono text-4xl font-semibold text-slate-900">{{ item.value }}</p>
          <div class="mt-3 inline-flex items-center gap-1 rounded-xl px-2.5 py-1 text-xs font-semibold" :class="deltaClass(item.change, item.inverse)">
            <span>{{ trendArrow(item.change, item.inverse) }}</span>
            <span>{{ formatChange(item.change) }}</span>
          </div>
        </LiquidGlassCard>
      </template>
    </div>

    <LiquidGlassCard max-width="100%" padding="20px">
      <div class="flex items-center justify-between gap-3">
        <div>
          <h2 class="text-lg font-semibold text-slate-900">Traffic Overview</h2>
          <p class="text-sm text-slate-600">Visitors 和 Views 的时间分布</p>
        </div>
        <p class="text-xs font-semibold uppercase tracking-[0.18em] text-slate-500">{{ windowLabel }}</p>
      </div>

      <div v-if="status === 'loading' && !analytics" class="mt-4 flex h-[360px] items-center justify-center text-sm text-slate-400">
        加载中...
      </div>
      <div v-else class="mt-4 h-[360px]">
        <v-chart :option="trendOption" autoresize />
      </div>
    </LiquidGlassCard>

    <LiquidGlassCard max-width="100%" padding="20px">
      <div class="flex items-center justify-between gap-3">
        <div>
          <h3 class="text-lg font-semibold text-slate-900">近期事项</h3>
          <p class="text-sm text-slate-600">管理员最近操作记录</p>
        </div>
      </div>

      <div v-if="activityLoading" class="mt-4 flex h-[140px] items-center justify-center text-sm text-slate-400">
        加载中...
      </div>
      <div v-else-if="activityItems.length === 0" class="mt-4 rounded-2xl border border-slate-200/80 bg-white/40 px-4 py-8 text-center text-sm text-slate-500">
        暂无操作记录
      </div>
      <div v-else class="mt-4 grid gap-2 md:grid-cols-2">
        <div
          v-for="item in activityItems"
          :key="item.id"
          class="rounded-2xl border border-slate-200/80 bg-white/55 px-3 py-2.5"
        >
          <div class="flex items-start justify-between gap-2">
            <p class="text-sm text-slate-700">{{ formatActivity(item) }}</p>
            <span class="shrink-0 text-xs text-slate-500">{{ formatRelativeTime(item.created_at) }}</span>
          </div>
          <p v-if="item.admin_username" class="mt-1 text-xs text-slate-500">操作人：{{ item.admin_username }}</p>
        </div>
      </div>
    </LiquidGlassCard>

    <LiquidGlassCard v-if="degradedMode" max-width="100%" padding="14px">
      <p class="text-sm text-amber-700">
        Analytics 数据暂不可用，当前显示占位数据。请先在 `backend` 目录执行 `go run cmd/migrate/main.go`。
      </p>
    </LiquidGlassCard>

    <template v-if="analytics">
      <div class="grid gap-4 xl:grid-cols-2">
        <LiquidGlassCard max-width="100%" padding="20px">
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold text-slate-900">Pages</h3>
            <div class="inline-flex gap-1 rounded-xl border border-slate-200/80 bg-white/70 p-1 text-xs">
              <button
                v-for="tab in pageTabs"
                :key="tab.key"
                type="button"
                class="rounded-lg px-2.5 py-1 transition"
                :class="pageTab === tab.key ? 'bg-white text-miku shadow-sm' : 'text-slate-500 hover:text-slate-700'"
                @click="pageTab = tab.key"
              >
                {{ tab.label }}
              </button>
            </div>
          </div>

          <div class="mt-4 text-sm">
            <div class="grid grid-cols-[1fr_auto_auto] border-b border-slate-200/80 pb-2 text-xs font-semibold uppercase tracking-[0.16em] text-slate-500">
              <span>Path</span>
              <span class="text-right">Visitors</span>
              <span class="text-right">%</span>
            </div>
            <div v-for="row in visiblePageRows" :key="row.path" class="grid grid-cols-[1fr_auto_auto] items-center gap-2 border-b border-slate-100/90 py-2.5 text-slate-700">
              <span class="truncate">{{ row.path }}</span>
              <span class="text-right font-semibold text-slate-900">{{ formatInteger(pageMetricValue(row)) }}</span>
              <span class="text-right text-slate-500">{{ toPercent(pageMetricValue(row), pageMetricTotal) }}</span>
            </div>
          </div>

          <button
            v-if="pageRows.length > 8"
            type="button"
            class="mt-3 inline-flex items-center gap-1 text-sm font-semibold text-slate-600 transition hover:text-miku"
            @click="showAllPages = !showAllPages"
          >
            <span>{{ showAllPages ? 'Less' : 'More' }}</span>
          </button>
        </LiquidGlassCard>

        <LiquidGlassCard max-width="100%" padding="20px">
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold text-slate-900">Sources</h3>
            <div class="inline-flex gap-1 rounded-xl border border-slate-200/80 bg-white/70 p-1 text-xs">
              <button
                type="button"
                class="rounded-lg px-2.5 py-1 transition"
                :class="sourceTab === 'referrers' ? 'bg-white text-miku shadow-sm' : 'text-slate-500 hover:text-slate-700'"
                @click="sourceTab = 'referrers'"
              >
                Referrers
              </button>
              <button
                type="button"
                class="rounded-lg px-2.5 py-1 transition"
                :class="sourceTab === 'channels' ? 'bg-white text-miku shadow-sm' : 'text-slate-500 hover:text-slate-700'"
                @click="sourceTab = 'channels'"
              >
                Channels
              </button>
            </div>
          </div>

          <div class="mt-4 text-sm">
            <div class="grid grid-cols-[1fr_auto_auto] border-b border-slate-200/80 pb-2 text-xs font-semibold uppercase tracking-[0.16em] text-slate-500">
              <span>Source</span>
              <span class="text-right">Visitors</span>
              <span class="text-right">%</span>
            </div>
            <div v-for="row in visibleSourceRows" :key="row.name" class="grid grid-cols-[1fr_auto_auto] items-center gap-2 border-b border-slate-100/90 py-2.5 text-slate-700">
              <span class="truncate">{{ row.name }}</span>
              <span class="text-right font-semibold text-slate-900">{{ formatInteger(row.visitors) }}</span>
              <span class="text-right text-slate-500">{{ toPercent(row.visitors, sourceTotal) }}</span>
            </div>
          </div>

          <button
            v-if="sourceRows.length > 7"
            type="button"
            class="mt-3 inline-flex items-center gap-1 text-sm font-semibold text-slate-600 transition hover:text-miku"
            @click="showAllSources = !showAllSources"
          >
            <span>{{ showAllSources ? 'Less' : 'More' }}</span>
          </button>
        </LiquidGlassCard>

        <LiquidGlassCard max-width="100%" padding="20px">
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold text-slate-900">Environment</h3>
            <div class="inline-flex gap-1 rounded-xl border border-slate-200/80 bg-white/70 p-1 text-xs">
              <button
                v-for="tab in environmentTabs"
                :key="tab.key"
                type="button"
                class="rounded-lg px-2.5 py-1 transition"
                :class="environmentTab === tab.key ? 'bg-white text-miku shadow-sm' : 'text-slate-500 hover:text-slate-700'"
                @click="environmentTab = tab.key"
              >
                {{ tab.label }}
              </button>
            </div>
          </div>

          <div class="mt-4 text-sm">
            <div class="grid grid-cols-[1fr_auto_auto] border-b border-slate-200/80 pb-2 text-xs font-semibold uppercase tracking-[0.16em] text-slate-500">
              <span>Name</span>
              <span class="text-right">Visitors</span>
              <span class="text-right">%</span>
            </div>
            <div v-for="row in visibleEnvironmentRows" :key="row.name" class="grid grid-cols-[1fr_auto_auto] items-center gap-2 border-b border-slate-100/90 py-2.5 text-slate-700">
              <span class="truncate">{{ row.name }}</span>
              <span class="text-right font-semibold text-slate-900">{{ formatInteger(row.visitors) }}</span>
              <span class="text-right text-slate-500">{{ toPercent(row.visitors, environmentTotal) }}</span>
            </div>
          </div>

          <button
            v-if="environmentRows.length > 7"
            type="button"
            class="mt-3 inline-flex items-center gap-1 text-sm font-semibold text-slate-600 transition hover:text-miku"
            @click="showAllEnvironment = !showAllEnvironment"
          >
            <span>{{ showAllEnvironment ? 'Less' : 'More' }}</span>
          </button>
        </LiquidGlassCard>

        <LiquidGlassCard max-width="100%" padding="20px">
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold text-slate-900">Location</h3>
            <div class="inline-flex gap-1 rounded-xl border border-slate-200/80 bg-white/70 p-1 text-xs">
              <button
                v-for="tab in locationTabs"
                :key="tab.key"
                type="button"
                class="rounded-lg px-2.5 py-1 transition"
                :class="locationTab === tab.key ? 'bg-white text-miku shadow-sm' : 'text-slate-500 hover:text-slate-700'"
                @click="locationTab = tab.key"
              >
                {{ tab.label }}
              </button>
            </div>
          </div>

          <div class="mt-4 text-sm">
            <div class="grid grid-cols-[1fr_auto_auto] border-b border-slate-200/80 pb-2 text-xs font-semibold uppercase tracking-[0.16em] text-slate-500">
              <span>Region</span>
              <span class="text-right">Visitors</span>
              <span class="text-right">%</span>
            </div>
            <div v-for="row in visibleLocationRows" :key="row.name" class="grid grid-cols-[1fr_auto_auto] items-center gap-2 border-b border-slate-100/90 py-2.5 text-slate-700">
              <span class="truncate">{{ row.name }}</span>
              <span class="text-right font-semibold text-slate-900">{{ formatInteger(row.visitors) }}</span>
              <span class="text-right text-slate-500">{{ toPercent(row.visitors, locationTotal) }}</span>
            </div>
          </div>

          <button
            v-if="locationRows.length > 7"
            type="button"
            class="mt-3 inline-flex items-center gap-1 text-sm font-semibold text-slate-600 transition hover:text-miku"
            @click="showAllLocation = !showAllLocation"
          >
            <span>{{ showAllLocation ? 'Less' : 'More' }}</span>
          </button>
        </LiquidGlassCard>
      </div>

      <div class="grid gap-4 xl:grid-cols-[2fr_1fr]">
        <LiquidGlassCard max-width="100%" padding="20px">
          <div class="flex items-center justify-between gap-3">
            <h3 class="text-lg font-semibold text-slate-900">Geo Distribution</h3>
            <span class="text-xs text-slate-500">By Country</span>
          </div>

          <div class="mt-4 h-[360px]">
            <v-chart v-if="worldMapReady && mapSeries.length > 0" :option="mapOption" autoresize />
            <div v-else class="flex h-full items-center justify-center rounded-2xl border border-slate-200/80 bg-white/45 px-4 text-sm text-slate-500">
              {{ mapFallbackText }}
            </div>
          </div>

          <div class="mt-4 grid gap-2 sm:grid-cols-2">
            <div v-for="item in topCountries" :key="item.code" class="flex items-center justify-between rounded-xl border border-slate-200/80 bg-white/55 px-3 py-2 text-sm">
              <span class="truncate text-slate-700">{{ item.name }}</span>
              <span class="font-semibold text-slate-900">{{ formatInteger(item.visitors) }}</span>
            </div>
          </div>
        </LiquidGlassCard>

        <LiquidGlassCard max-width="100%" padding="20px">
          <h3 class="text-lg font-semibold text-slate-900">Traffic</h3>
          <div class="mt-4 grid grid-cols-[42px_repeat(7,minmax(0,1fr))] items-center gap-x-1.5 gap-y-1 text-xs text-slate-600">
            <div />
            <div v-for="day in weekDays" :key="day" class="text-center font-semibold">{{ day }}</div>

            <template v-for="hour in hours" :key="`hour-${hour}`">
              <div class="pr-1 text-right text-slate-500">{{ hourLabel(hour) }}</div>
              <div v-for="dayIndex in 7" :key="`dot-${hour}-${dayIndex}`" class="flex justify-center">
                <span class="rounded-full border border-slate-200/80" :style="trafficDotStyle(dayIndex - 1, hour)" :title="trafficDotTitle(dayIndex - 1, hour)" />
              </div>
            </template>
          </div>
        </LiquidGlassCard>
      </div>
    </template>
  </section>
</template>

<script setup lang="ts">
import { useStore } from '@nanostores/vue'
import { computed, onMounted, ref, watch } from 'vue'
import * as echarts from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { BarChart, MapChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, LegendComponent, VisualMapComponent, GeoComponent } from 'echarts/components'
import VChart from 'vue-echarts'

import { api } from '../../lib/api'
import { authState, hydrateAuth } from '../../stores/auth'
import { setScopeStatus } from '../../stores/loading'
import LiquidGlassCard from '../ui/LiquidGlassCard.vue'
import SkeletonCard from '../ui/SkeletonCard.vue'

echarts.use([CanvasRenderer, BarChart, MapChart, GridComponent, TooltipComponent, LegendComponent, VisualMapComponent, GeoComponent])

type RangeKey = '24h' | '7d' | '30d'
type PageTabKey = 'path' | 'entry' | 'exit'
type SourceTabKey = 'referrers' | 'channels'
type EnvironmentTabKey = 'browsers' | 'os' | 'devices'
type LocationTabKey = 'countries' | 'regions' | 'cities'

interface MetricValue {
  value: number
  change: number
}

interface AnalyticsNamedCount {
  name: string
  visitors: number
}

interface AnalyticsCountryCount {
  code: string
  visitors: number
}

interface AnalyticsPageItem {
  path: string
  visitors: number
  views: number
  entries: number
  exits: number
}

interface AnalyticsTrendItem {
  bucket: string
  visitors: number
  views: number
}

interface AnalyticsTrafficPoint {
  dow: number
  hour: number
  value: number
}

interface AnalyticsOverview {
  window: {
    range: string
    label: string
    start: string
    end: string
    previous_from: string
    previous_to: string
    granularity: string
  }
  summary: {
    visitors: MetricValue
    visits: MetricValue
    views: MetricValue
    bounce_rate: MetricValue
    visit_duration: MetricValue
  }
  trend: AnalyticsTrendItem[]
  pages: AnalyticsPageItem[]
  sources: {
    referrers: AnalyticsNamedCount[]
    channels: AnalyticsNamedCount[]
  }
  environment: {
    browsers: AnalyticsNamedCount[]
    os: AnalyticsNamedCount[]
    devices: AnalyticsNamedCount[]
  }
  location: {
    countries: AnalyticsCountryCount[]
    regions: AnalyticsNamedCount[]
    cities: AnalyticsNamedCount[]
  }
  traffic: AnalyticsTrafficPoint[]
}

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

const auth = useStore(authState)
const status = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const analytics = ref<AnalyticsOverview | null>(null)
const degradedMode = ref(false)
const selectedRange = ref<RangeKey>('24h')
const offset = ref(0)
const showFilters = ref(false)
const pathKeyword = ref('')

const worldMapReady = ref(false)
const worldMapLoading = ref(false)

const pageTab = ref<PageTabKey>('path')
const sourceTab = ref<SourceTabKey>('referrers')
const environmentTab = ref<EnvironmentTabKey>('browsers')
const locationTab = ref<LocationTabKey>('countries')

const showAllPages = ref(false)
const showAllSources = ref(false)
const showAllEnvironment = ref(false)
const showAllLocation = ref(false)

const activityItems = ref<AuditLogItem[]>([])
const activityLoading = ref(true)

const rangeOptions: Array<{ key: RangeKey; label: string }> = [
  { key: '24h', label: 'Last 24 hours' },
  { key: '7d', label: 'Last 7 days' },
  { key: '30d', label: 'Last 30 days' },
]

const pageTabs: Array<{ key: PageTabKey; label: string }> = [
  { key: 'path', label: 'Path' },
  { key: 'entry', label: 'Entry' },
  { key: 'exit', label: 'Exit' },
]

const environmentTabs: Array<{ key: EnvironmentTabKey; label: string }> = [
  { key: 'browsers', label: 'Browsers' },
  { key: 'os', label: 'OS' },
  { key: 'devices', label: 'Devices' },
]

const locationTabs: Array<{ key: LocationTabKey; label: string }> = [
  { key: 'countries', label: 'Countries' },
  { key: 'regions', label: 'Regions' },
  { key: 'cities', label: 'Cities' },
]

const weekDays = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat']
const hours = Array.from({ length: 24 }, (_, idx) => idx)

const windowLabel = computed(() => {
  const base = analytics.value?.window.label || rangeOptions.find((x) => x.key === selectedRange.value)?.label || 'Window'
  if (offset.value === 0) return base
  return `${base} · -${offset.value}`
})

const formattedWindowRange = computed(() => {
  const start = analytics.value?.window.start
  const end = analytics.value?.window.end
  if (!start || !end) return '--'
  return `${formatDateTime(start)} ~ ${formatDateTime(end)}`
})

const statCards = computed(() => {
  const s = analytics.value?.summary
  if (!s) return []
  return [
    { label: 'Visitors', value: formatInteger(s.visitors.value), change: s.visitors.change, inverse: false },
    { label: 'Visits', value: formatInteger(s.visits.value), change: s.visits.change, inverse: false },
    { label: 'Views', value: formatInteger(s.views.value), change: s.views.change, inverse: false },
    { label: 'Bounce rate', value: `${s.bounce_rate.value.toFixed(1)}%`, change: s.bounce_rate.change, inverse: true },
    { label: 'Visit duration', value: formatDuration(s.visit_duration.value), change: s.visit_duration.change, inverse: false },
  ]
})

const trendOption = computed(() => {
  const trend = analytics.value?.trend || []
  const labels = trend.map((item) => formatBucketLabel(item.bucket))
  return {
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'shadow' },
      backgroundColor: 'rgba(255,255,255,0.96)',
      borderColor: 'rgba(148,163,184,0.25)',
      textStyle: { color: '#334155', fontSize: 12 },
    },
    legend: {
      right: 0,
      top: 0,
      textStyle: { color: '#64748b', fontSize: 12 },
      data: ['Visitors', 'Views'],
    },
    grid: {
      left: 12,
      right: 16,
      top: 42,
      bottom: 8,
      containLabel: true,
    },
    xAxis: {
      type: 'category',
      data: labels,
      axisLine: { lineStyle: { color: '#dbe4ef' } },
      axisLabel: { color: '#94a3b8', fontSize: 11 },
      splitLine: { show: false },
    },
    yAxis: {
      type: 'value',
      axisLabel: { color: '#94a3b8', fontSize: 11 },
      splitLine: { lineStyle: { color: '#edf2f9' } },
    },
    series: [
      {
        name: 'Visitors',
        type: 'bar',
        data: trend.map((item) => item.visitors),
        barMaxWidth: 28,
        itemStyle: {
          borderRadius: [8, 8, 0, 0],
          color: '#39c5bb',
        },
      },
      {
        name: 'Views',
        type: 'bar',
        data: trend.map((item) => item.views),
        barMaxWidth: 28,
        itemStyle: {
          borderRadius: [8, 8, 0, 0],
          color: '#c084fc',
        },
      },
    ],
  }
})

const pageRows = computed(() => {
  const keyword = pathKeyword.value.toLowerCase()
  let rows = analytics.value?.pages || []
  if (keyword) {
    rows = rows.filter((item) => item.path.toLowerCase().includes(keyword))
  }
  return [...rows].sort((a, b) => pageMetricValue(b) - pageMetricValue(a))
})

const pageMetricTotal = computed(() => pageRows.value.reduce((sum, item) => sum + pageMetricValue(item), 0))
const visiblePageRows = computed(() => (showAllPages.value ? pageRows.value : pageRows.value.slice(0, 8)))

const sourceRows = computed(() => {
  const rows = sourceTab.value === 'referrers'
    ? analytics.value?.sources.referrers || []
    : analytics.value?.sources.channels || []
  return [...rows].sort((a, b) => b.visitors - a.visitors)
})

const sourceTotal = computed(() => sourceRows.value.reduce((sum, item) => sum + item.visitors, 0))
const visibleSourceRows = computed(() => (showAllSources.value ? sourceRows.value : sourceRows.value.slice(0, 7)))

const environmentRows = computed(() => {
  if (!analytics.value) return []
  if (environmentTab.value === 'os') return [...analytics.value.environment.os].sort((a, b) => b.visitors - a.visitors)
  if (environmentTab.value === 'devices') return [...analytics.value.environment.devices].sort((a, b) => b.visitors - a.visitors)
  return [...analytics.value.environment.browsers].sort((a, b) => b.visitors - a.visitors)
})

const environmentTotal = computed(() => environmentRows.value.reduce((sum, item) => sum + item.visitors, 0))
const visibleEnvironmentRows = computed(() => (showAllEnvironment.value ? environmentRows.value : environmentRows.value.slice(0, 7)))

const locationRows = computed(() => {
  if (!analytics.value) return []
  if (locationTab.value === 'regions') {
    return [...analytics.value.location.regions].sort((a, b) => b.visitors - a.visitors).map((item) => ({ name: item.name, visitors: item.visitors }))
  }
  if (locationTab.value === 'cities') {
    return [...analytics.value.location.cities].sort((a, b) => b.visitors - a.visitors).map((item) => ({ name: item.name, visitors: item.visitors }))
  }
  return [...analytics.value.location.countries]
    .sort((a, b) => b.visitors - a.visitors)
    .map((item) => ({ name: countryName(item.code), visitors: item.visitors }))
})

const locationTotal = computed(() => locationRows.value.reduce((sum, item) => sum + item.visitors, 0))
const visibleLocationRows = computed(() => (showAllLocation.value ? locationRows.value : locationRows.value.slice(0, 7)))

const mapSeries = computed(() => {
  if (!analytics.value) return []
  return analytics.value.location.countries
    .filter((item) => item.code !== 'ZZ' && item.visitors > 0)
    .map((item) => ({ name: countryName(item.code), value: item.visitors }))
})

const mapOption = computed(() => {
  const maxValue = mapSeries.value.reduce((m, item) => Math.max(m, Number(item.value) || 0), 0)
  return {
    tooltip: {
      trigger: 'item',
      formatter: (params: { name?: string; value?: number }) => `${params.name || 'Unknown'}: ${formatInteger(params.value || 0)}`,
      backgroundColor: 'rgba(255,255,255,0.96)',
      borderColor: 'rgba(148,163,184,0.25)',
      textStyle: { color: '#334155' },
    },
    visualMap: {
      min: 0,
      max: maxValue || 1,
      show: false,
      inRange: {
        color: ['#d9f5f2', '#39c5bb'],
      },
    },
    series: [
      {
        type: 'map',
        map: 'world',
        roam: true,
        zoom: 1.08,
        itemStyle: {
          areaColor: '#edf7ff',
          borderColor: '#7ad6ce',
          borderWidth: 0.8,
        },
        emphasis: {
          itemStyle: { areaColor: '#c9f1ec' },
          label: { show: false },
        },
        data: mapSeries.value,
      },
    ],
  }
})

const topCountries = computed(() => {
  if (!analytics.value) return []
  return analytics.value.location.countries
    .slice()
    .sort((a, b) => b.visitors - a.visitors)
    .slice(0, 8)
    .map((item) => ({ code: item.code, name: countryName(item.code), visitors: item.visitors }))
})

const mapFallbackText = computed(() => {
  if (worldMapLoading.value) return '世界地图加载中...'
  if (mapSeries.value.length === 0) return '暂无国家分布数据'
  return '无法加载地图底图，已在下方显示国家排行'
})

const trafficLookup = computed(() => {
  const map = new Map<string, number>()
  ;(analytics.value?.traffic || []).forEach((item) => {
    map.set(`${item.dow}-${item.hour}`, item.value)
  })
  return map
})

const trafficMax = computed(() => {
  let max = 0
  trafficLookup.value.forEach((value) => {
    if (value > max) max = value
  })
  return max
})

function buildWindow(range: RangeKey, pageOffset: number) {
  const safeOffset = Math.max(0, pageOffset || 0)
  const now = new Date()

  if (range === '24h') {
    const end = new Date(now)
    end.setMinutes(0, 0, 0)
    end.setHours(end.getHours() + 1 - safeOffset * 24)
    const start = new Date(end)
    start.setHours(start.getHours() - 24)

    return {
      label: 'Last 24 hours',
      granularity: 'hour',
      start,
      end,
      previousFrom: new Date(start.getTime() - 24 * 60 * 60 * 1000),
      previousTo: start,
    }
  }

  const days = range === '7d' ? 7 : 30
  const end = new Date(now)
  end.setHours(0, 0, 0, 0)
  end.setDate(end.getDate() + 1 - safeOffset * days)
  const start = new Date(end)
  start.setDate(start.getDate() - days)

  return {
    label: range === '7d' ? 'Last 7 days' : 'Last 30 days',
    granularity: 'day',
    start,
    end,
    previousFrom: new Date(start.getTime() - days * 24 * 60 * 60 * 1000),
    previousTo: start,
  }
}

function buildEmptyOverview(range: RangeKey, pageOffset: number): AnalyticsOverview {
  const w = buildWindow(range, pageOffset)
  return {
    window: {
      range,
      label: w.label,
      start: w.start.toISOString(),
      end: w.end.toISOString(),
      previous_from: w.previousFrom.toISOString(),
      previous_to: w.previousTo.toISOString(),
      granularity: w.granularity,
    },
    summary: {
      visitors: { value: 0, change: 0 },
      visits: { value: 0, change: 0 },
      views: { value: 0, change: 0 },
      bounce_rate: { value: 0, change: 0 },
      visit_duration: { value: 0, change: 0 },
    },
    trend: [],
    pages: [],
    sources: { referrers: [], channels: [] },
    environment: { browsers: [], os: [], devices: [] },
    location: { countries: [], regions: [], cities: [] },
    traffic: [],
  }
}

watch(selectedRange, async () => {
  offset.value = 0
  showAllPages.value = false
  showAllSources.value = false
  showAllEnvironment.value = false
  showAllLocation.value = false
  await loadAnalytics()
})

async function loadActivities() {
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

async function loadAnalytics() {
  status.value = 'loading'
  degradedMode.value = false
  setScopeStatus('adminDashboard', 'loading')
  try {
    const data = await api.get<AnalyticsOverview>(`/admin/dashboard/analytics?range=${selectedRange.value}&offset=${offset.value}`)
    analytics.value = data
    degradedMode.value = false
    status.value = 'success'
    setScopeStatus('adminDashboard', 'success')
    if (data.location.countries.length > 0) {
      void ensureWorldMap()
    }
  } catch {
    analytics.value = buildEmptyOverview(selectedRange.value, offset.value)
    degradedMode.value = true
    status.value = 'success'
    setScopeStatus('adminDashboard', 'success')
  }
}

async function ensureWorldMap() {
  if (worldMapReady.value || worldMapLoading.value) return
  try {
    worldMapLoading.value = true
    const hasMap = typeof (echarts as unknown as { getMap?: (name: string) => unknown }).getMap === 'function'
      ? (echarts as unknown as { getMap: (name: string) => unknown }).getMap('world')
      : null
    if (!hasMap) {
      const response = await fetch('https://cdn.jsdelivr.net/npm/echarts@5/map/json/world.json')
      if (!response.ok) return
      const geoJson = await response.json()
      echarts.registerMap('world', geoJson)
    }
    worldMapReady.value = true
  } catch {
    worldMapReady.value = false
  } finally {
    worldMapLoading.value = false
  }
}

function goPrevWindow() {
  offset.value += 1
  void loadAnalytics()
}

function goNextWindow() {
  if (offset.value === 0) return
  offset.value -= 1
  void loadAnalytics()
}

function pageMetricValue(item: AnalyticsPageItem): number {
  if (pageTab.value === 'entry') return item.entries
  if (pageTab.value === 'exit') return item.exits
  return item.visitors
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

function formatActivity(item: AuditLogItem): string {
  const action = ACTION_MAP[item.action] || item.action
  const target = TARGET_MAP[item.target_type] || item.target_type
  return `${action}一条${target}`
}

function formatRelativeTime(iso: string): string {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return iso

  const now = Date.now()
  const diffMs = now - d.getTime()
  const mins = Math.floor(diffMs / 60000)
  if (mins < 1) return '刚刚'
  if (mins < 60) return `${mins} 分钟前`

  const hoursDiff = Math.floor(mins / 60)
  if (hoursDiff < 24) return `${hoursDiff} 小时前`

  const days = Math.floor(hoursDiff / 24)
  if (days < 7) return `${days} 天前`

  return d.toLocaleString('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  })
}

function deltaClass(change: number, inverse: boolean): string {
  if (change === 0) return 'bg-slate-100 text-slate-600'
  const positive = change > 0
  const good = inverse ? !positive : positive
  return good ? 'bg-emerald-100/80 text-emerald-700' : 'bg-rose-100/80 text-rose-700'
}

function trendArrow(change: number, inverse: boolean): string {
  if (change === 0) return '·'
  const positive = change > 0
  const good = inverse ? !positive : positive
  return good ? '↑' : '↓'
}

function formatChange(change: number): string {
  return `${Math.abs(change).toFixed(1)}%`
}

function formatInteger(value: number): string {
  return Number(value || 0).toLocaleString('en-US')
}

function formatDuration(seconds: number): string {
  const total = Math.max(0, Math.floor(seconds || 0))
  const mins = Math.floor(total / 60)
  const secs = total % 60
  if (mins <= 0) return `${secs}s`
  return `${mins}m ${secs}s`
}

function toPercent(value: number, total: number): string {
  if (total <= 0) return '0%'
  return `${((value / total) * 100).toFixed(0)}%`
}

function formatBucketLabel(bucket: string): string {
  if (!bucket) return ''
  if (analytics.value?.window.granularity === 'hour') return bucket.slice(11)
  return bucket.slice(5)
}

function formatDateTime(iso: string): string {
  const date = new Date(iso)
  if (Number.isNaN(date.getTime())) return iso
  return date.toLocaleString('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  })
}

function countryName(code: string): string {
  if (!code || code === 'ZZ') return 'Unknown'
  try {
    if (typeof Intl !== 'undefined' && 'DisplayNames' in Intl) {
      const displayNames = new Intl.DisplayNames(['en'], { type: 'region' })
      return displayNames.of(code) || code
    }
  } catch {
    // ignore
  }
  return code
}

function hourLabel(hour: number): string {
  if (hour === 0) return '12am'
  if (hour < 12) return `${hour}am`
  if (hour === 12) return '12pm'
  return `${hour - 12}pm`
}

function trafficDotStyle(dow: number, hour: number): Record<string, string> {
  const value = trafficLookup.value.get(`${dow}-${hour}`) || 0
  const max = trafficMax.value
  if (value <= 0 || max <= 0) {
    return {
      width: '8px',
      height: '8px',
      opacity: '0.25',
      backgroundColor: '#cbd5e1',
      display: 'block',
    }
  }
  const ratio = value / max
  const size = Math.round(8 + ratio * 12)
  const opacity = (0.38 + ratio * 0.62).toFixed(2)
  return {
    width: `${size}px`,
    height: `${size}px`,
    opacity,
    backgroundColor: ratio > 0.55 ? '#39c5bb' : '#7ccaf6',
    display: 'block',
  }
}

function trafficDotTitle(dow: number, hour: number): string {
  const count = trafficLookup.value.get(`${dow}-${hour}`) || 0
  return `${weekDays[dow]} ${hourLabel(hour)}: ${formatInteger(count)}`
}

onMounted(async () => {
  hydrateAuth()
  const latestAuth = authState.get()
  if (latestAuth.status !== 'authenticated') {
    window.location.replace('/login')
    return
  }
  await Promise.all([loadAnalytics(), loadActivities()])
})
</script>
