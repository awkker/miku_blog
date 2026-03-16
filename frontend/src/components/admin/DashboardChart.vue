<template>
  <LiquidGlassCard padding="20px">
    <h2 class="text-lg font-semibold text-slate-900">{{ title }}</h2>
    <div v-if="loading" class="mt-4 flex items-center justify-center" style="height: 320px">
      <p class="text-sm text-slate-400">加载中...</p>
    </div>
    <div v-else class="mt-4" style="height: 320px">
      <v-chart :option="chartOption" autoresize />
    </div>
  </LiquidGlassCard>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart } from 'echarts/charts'
import {
  GridComponent,
  TooltipComponent,
  LegendComponent,
} from 'echarts/components'
import VChart from 'vue-echarts'

import LiquidGlassCard from '../ui/LiquidGlassCard.vue'

use([CanvasRenderer, LineChart, GridComponent, TooltipComponent, LegendComponent])

export interface ChartSeries {
  name: string
  data: number[]
  color: string
}

interface Props {
  title?: string
  labels?: string[]
  series?: ChartSeries[]
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  title: '近 7 日趋势',
  labels: () => [],
  series: () => [],
  loading: false,
})

const MIKU = '#39c5bb'
const LAVENDER = '#c084fc'

const defaultColors = [MIKU, LAVENDER, '#f59e0b', '#ef4444']

const chartOption = computed(() => {
  const seriesData = props.series.map((s, i) => ({
    name: s.name,
    type: 'line' as const,
    smooth: true,
    symbol: 'circle' as const,
    symbolSize: 6,
    data: s.data,
    lineStyle: { color: s.color || defaultColors[i % defaultColors.length], width: 2.5 },
    itemStyle: { color: s.color || defaultColors[i % defaultColors.length] },
    areaStyle: {
      color: {
        type: 'linear' as const,
        x: 0, y: 0, x2: 0, y2: 1,
        colorStops: [
          { offset: 0, color: (s.color || defaultColors[i % defaultColors.length]).replace(')', ',0.25)').replace('rgb', 'rgba') },
          { offset: 1, color: (s.color || defaultColors[i % defaultColors.length]).replace(')', ',0.02)').replace('rgb', 'rgba') },
        ],
      },
    },
  }))

  return {
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'rgba(255,255,255,0.92)',
      borderColor: 'rgba(0,0,0,0.06)',
      textStyle: { color: '#334155', fontSize: 12 },
    },
    legend: {
      data: props.series.map((s) => s.name),
      top: 0,
      right: 0,
      textStyle: { color: '#64748b', fontSize: 12 },
    },
    grid: {
      top: 36,
      left: 12,
      right: 12,
      bottom: 0,
      containLabel: true,
    },
    xAxis: {
      type: 'category',
      data: props.labels,
      boundaryGap: false,
      axisLine: { lineStyle: { color: '#e2e8f0' } },
      axisLabel: { color: '#94a3b8', fontSize: 11 },
    },
    yAxis: {
      type: 'value',
      splitLine: { lineStyle: { color: '#f1f5f9' } },
      axisLabel: { color: '#94a3b8', fontSize: 11 },
    },
    series: seriesData,
  }
})
</script>
