<template>
  <LiquidGlassCard padding="20px">
    <h2 class="text-lg font-semibold text-slate-900">{{ title }}</h2>
    <div class="mt-4" style="height: 320px">
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

interface Props {
  title?: string
}

withDefaults(defineProps<Props>(), {
  title: '近 7 日趋势',
})

const MIKU = '#39c5bb'
const LAVENDER = '#c084fc'

const days = ['03-07', '03-08', '03-09', '03-10', '03-11', '03-12', '03-13']

const chartOption = computed(() => ({
  tooltip: {
    trigger: 'axis',
    backgroundColor: 'rgba(255,255,255,0.92)',
    borderColor: 'rgba(0,0,0,0.06)',
    textStyle: { color: '#334155', fontSize: 12 },
  },
  legend: {
    data: ['浏览量', '评论数'],
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
    data: days,
    boundaryGap: false,
    axisLine: { lineStyle: { color: '#e2e8f0' } },
    axisLabel: { color: '#94a3b8', fontSize: 11 },
  },
  yAxis: {
    type: 'value',
    splitLine: { lineStyle: { color: '#f1f5f9' } },
    axisLabel: { color: '#94a3b8', fontSize: 11 },
  },
  series: [
    {
      name: '浏览量',
      type: 'line',
      smooth: true,
      symbol: 'circle',
      symbolSize: 6,
      data: [820, 932, 1105, 864, 1230, 1042, 1380],
      lineStyle: { color: MIKU, width: 2.5 },
      itemStyle: { color: MIKU },
      areaStyle: {
        color: {
          type: 'linear',
          x: 0, y: 0, x2: 0, y2: 1,
          colorStops: [
            { offset: 0, color: 'rgba(57,197,187,0.25)' },
            { offset: 1, color: 'rgba(57,197,187,0.02)' },
          ],
        },
      },
    },
    {
      name: '评论数',
      type: 'line',
      smooth: true,
      symbol: 'circle',
      symbolSize: 6,
      data: [12, 18, 9, 22, 15, 28, 20],
      lineStyle: { color: LAVENDER, width: 2.5 },
      itemStyle: { color: LAVENDER },
      areaStyle: {
        color: {
          type: 'linear',
          x: 0, y: 0, x2: 0, y2: 1,
          colorStops: [
            { offset: 0, color: 'rgba(192,132,252,0.2)' },
            { offset: 1, color: 'rgba(192,132,252,0.02)' },
          ],
        },
      },
    },
  ],
}))
</script>
