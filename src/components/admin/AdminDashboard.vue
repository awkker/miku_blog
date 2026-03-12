<template>
  <section class="space-y-5">
    <LiquidGlassCard padding="24px">
      <h1 class="text-2xl font-semibold text-white">欢迎回来，{{ userName }}</h1>
      <p class="mt-2 text-sm text-white/75">后台已就绪，你可以继续管理文章、评论与友链。</p>
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

    <div class="grid gap-4 lg:grid-cols-3">
      <LiquidGlassCard padding="20px" class="lg:col-span-2">
        <h2 class="text-lg font-semibold text-white">最近动态</h2>
        <ul class="mt-4 space-y-3 text-sm text-white/80">
          <li class="rounded-2xl border border-white/20 bg-white/10 px-3 py-2">新文章草稿《Astro + Vue Islands 实战》待发布</li>
          <li class="rounded-2xl border border-white/20 bg-white/10 px-3 py-2">评论审核队列中有 3 条待处理内容</li>
          <li class="rounded-2xl border border-white/20 bg-white/10 px-3 py-2">友链申请箱新增 1 条请求</li>
        </ul>
      </LiquidGlassCard>

      <LiquidGlassCard padding="20px">
        <h2 class="text-lg font-semibold text-white">快捷操作</h2>
        <div class="mt-4 grid gap-2">
          <MikuButton variant="ghost" aria-label="新建文章">新建文章</MikuButton>
          <MikuButton variant="ghost" aria-label="查看站点">查看站点</MikuButton>
          <MikuButton variant="ghost" aria-label="管理友链">管理友链</MikuButton>
        </div>
      </LiquidGlassCard>
    </div>
  </section>
</template>

<script setup lang="ts">
import { useStore } from '@nanostores/vue'
import { computed, onMounted, ref } from 'vue'

import { authState, hydrateAuth } from '../../stores/auth'
import { setScopeStatus } from '../../stores/loading'
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

const auth = useStore(authState)

const stats = ref<StatItem[]>([])
const status = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const mounted = ref(false)

const userName = computed(() => {
  if (!mounted.value) {
    return '管理员'
  }

  return auth.value.user?.name ?? '管理员'
})

function sleep(duration = 680) {
  return new Promise<void>((resolve) => {
    setTimeout(resolve, duration)
  })
}

async function loadDashboard() {
  status.value = 'loading'
  setScopeStatus('adminDashboard', 'loading')

  try {
    await sleep()
    stats.value = [
      { label: '文章总数', value: '128', trend: '+8.2% 本周', icon: 'article' },
      { label: '待审评论', value: '3', trend: '-2 本日', icon: 'comment' },
      { label: '累计点赞', value: '4,298', trend: '+6.7% 本周', icon: 'like' },
      { label: '友链数量', value: '24', trend: '+1 本月', icon: 'link' },
    ]
    status.value = 'success'
    setScopeStatus('adminDashboard', 'success')
  } catch {
    status.value = 'error'
    setScopeStatus('adminDashboard', 'error')
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

  await loadDashboard()
})
</script>
