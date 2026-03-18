<template>
  <div>
    <h3 class="text-sm font-semibold text-slate-900">站点趋势</h3>
    <p class="mt-1 text-[11px] text-slate-500">最近 7 天访问热度</p>

    <div v-if="!loaded" class="mt-3 flex h-[96px] items-center justify-center">
      <span class="text-xs text-slate-300">加载中...</span>
    </div>

    <div v-else-if="hasData" class="mt-3">
      <svg :viewBox="`0 0 ${W} ${H}`" class="w-full" role="img" aria-label="站点趋势图">
        <defs>
          <linearGradient id="siteTrendGrad" x1="0%" x2="100%" y1="0%" y2="0%">
            <stop offset="0%" stop-color="#39c5bb" />
            <stop offset="100%" stop-color="#c084fc" />
          </linearGradient>
          <linearGradient id="siteTrendFill" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stop-color="rgba(57,197,187,0.18)" />
            <stop offset="100%" stop-color="rgba(57,197,187,0.01)" />
          </linearGradient>
        </defs>
        <path :d="areaPath" fill="url(#siteTrendFill)" />
        <path :d="linePath" fill="none" stroke="url(#siteTrendGrad)" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" />
        <circle
          v-for="(pt, i) in points"
          :key="i"
          :cx="pt.x"
          :cy="pt.y"
          r="3"
          fill="white"
          stroke="#39c5bb"
          stroke-width="1.5"
        />
      </svg>
      <div class="mt-1.5 flex justify-between text-[10px] text-slate-400">
        <span v-for="(label, i) in dayLabels" :key="i">{{ label }}</span>
      </div>
    </div>

    <div v-else class="mt-3 flex h-[96px] items-center justify-center">
      <span class="text-xs text-slate-300">暂无数据</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

import { api, type PagedData } from '../../lib/api'

interface PostItem {
  published_at?: string
  created_at: string
  view_count: number
}

const W = 280
const H = 96
const PAD_X = 16
const PAD_Y = 12

const dailyViews = ref<number[]>([])
const dayLabels = ref<string[]>([])
const loaded = ref(false)

const hasData = computed(() => dailyViews.value.some((v) => v > 0))

const points = computed(() => {
  const data = dailyViews.value
  if (data.length === 0) return []
  const max = Math.max(...data, 1)
  const usableW = W - PAD_X * 2
  const usableH = H - PAD_Y * 2
  const step = usableW / Math.max(data.length - 1, 1)
  return data.map((v, i) => ({
    x: PAD_X + i * step,
    y: PAD_Y + usableH - (v / max) * usableH,
  }))
})

const linePath = computed(() => {
  const pts = points.value
  if (pts.length < 2) return ''
  let d = `M${pts[0].x},${pts[0].y}`
  for (let i = 1; i < pts.length; i++) {
    const prev = pts[i - 1]
    const curr = pts[i]
    const cpx1 = prev.x + (curr.x - prev.x) * 0.4
    const cpx2 = curr.x - (curr.x - prev.x) * 0.4
    d += ` C${cpx1},${prev.y} ${cpx2},${curr.y} ${curr.x},${curr.y}`
  }
  return d
})

const areaPath = computed(() => {
  const pts = points.value
  if (pts.length < 2) return ''
  const bottom = H - PAD_Y + 4
  return `${linePath.value} L${pts[pts.length - 1].x},${bottom} L${pts[0].x},${bottom} Z`
})

async function loadTrend() {
  try {
    const data = await api.get<PagedData<PostItem>>('/posts?page=1&size=200')
    const items = data.items || []

    const now = new Date()
    const buckets: Record<string, number> = {}
    const labels: string[] = []

    for (let i = 6; i >= 0; i--) {
      const d = new Date(now)
      d.setDate(d.getDate() - i)
      const key = `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
      buckets[key] = 0
      labels.push(`${d.getMonth() + 1}/${d.getDate()}`)
    }

    for (const item of items) {
      const dateStr = item.published_at || item.created_at
      if (!dateStr) continue
      const key = dateStr.slice(0, 10)
      if (key in buckets) {
        buckets[key] += Number(item.view_count) || 0
      }
    }

    dayLabels.value = labels
    dailyViews.value = Object.values(buckets)
  } catch {
    dailyViews.value = []
    dayLabels.value = []
  } finally {
    loaded.value = true
  }
}

onMounted(() => {
  loadTrend()
})
</script>
