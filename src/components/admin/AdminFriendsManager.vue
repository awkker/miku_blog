<template>
  <section class="space-y-5">
    <LiquidGlassCard padding="24px">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 class="text-2xl font-semibold text-slate-900">友链管理</h1>
          <p class="mt-1 text-sm text-slate-600">维护友链展示、申请审核和站点健康检查。</p>
        </div>
        <MikuButton variant="solid" aria-label="添加友链">+ 添加友链</MikuButton>
      </div>
    </LiquidGlassCard>

    <div class="grid gap-4 sm:grid-cols-3">
      <LiquidGlassCard padding="16px">
        <p class="text-xs uppercase tracking-[0.2em] text-slate-600">已通过</p>
        <p class="mt-1 font-mono text-2xl font-semibold text-slate-900">{{ approvedCount }}</p>
      </LiquidGlassCard>
      <LiquidGlassCard padding="16px">
        <p class="text-xs uppercase tracking-[0.2em] text-slate-600">待审核</p>
        <p class="mt-1 font-mono text-2xl font-semibold text-amber-600">{{ pendingCount }}</p>
      </LiquidGlassCard>
      <LiquidGlassCard padding="16px">
        <p class="text-xs uppercase tracking-[0.2em] text-slate-600">异常 / 不可达</p>
        <p class="mt-1 font-mono text-2xl font-semibold text-red-600">{{ downCount }}</p>
      </LiquidGlassCard>
    </div>

    <LiquidGlassCard padding="0px">
      <div class="overflow-x-auto">
        <table class="w-full text-left text-sm">
          <thead>
            <tr class="border-b border-slate-200/60">
              <th class="px-5 py-3.5 font-semibold text-slate-700">站点名称</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700">链接</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700">状态</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700">健康检测</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700">添加时间</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700 text-center">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="link in friends"
              :key="link.id"
              class="border-b border-slate-100/60 transition hover:bg-white/40"
            >
              <td class="px-5 py-3.5 font-medium text-slate-900">{{ link.name }}</td>
              <td class="px-5 py-3.5">
                <a
                  :href="link.url"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="text-miku underline decoration-miku/30 transition hover:decoration-miku"
                >
                  {{ link.url }}
                </a>
              </td>
              <td class="px-5 py-3.5">
                <span
                  class="inline-block rounded-full px-2.5 py-0.5 text-xs font-medium"
                  :class="statusClass(link.status)"
                >
                  {{ statusLabel(link.status) }}
                </span>
              </td>
              <td class="px-5 py-3.5">
                <span
                  class="inline-block rounded-full px-2.5 py-0.5 text-xs font-medium"
                  :class="link.health === 'ok' ? 'bg-emerald-100 text-emerald-700' : 'bg-red-100 text-red-600'"
                >
                  {{ link.health === 'ok' ? '正常' : '不可达' }}
                </span>
              </td>
              <td class="px-5 py-3.5 text-slate-600">{{ link.createdAt }}</td>
              <td class="px-5 py-3.5 text-center">
                <div class="flex items-center justify-center gap-2">
                  <button
                    v-if="link.status === 'pending'"
                    type="button"
                    class="rounded-xl border border-emerald-200/80 bg-white/50 px-2.5 py-1 text-xs text-emerald-600 transition hover:bg-emerald-50"
                    aria-label="通过申请"
                  >
                    通过
                  </button>
                  <button
                    type="button"
                    class="rounded-xl border border-slate-200/80 bg-white/50 px-2.5 py-1 text-xs text-slate-700 transition hover:border-miku/40 hover:text-miku"
                    aria-label="编辑友链"
                  >
                    编辑
                  </button>
                  <button
                    type="button"
                    class="rounded-xl border border-red-200/80 bg-white/50 px-2.5 py-1 text-xs text-red-600 transition hover:bg-red-50"
                    aria-label="删除友链"
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
import { computed, ref } from 'vue'

import LiquidGlassCard from '../ui/LiquidGlassCard.vue'
import MikuButton from '../ui/MikuButton.vue'

interface FriendLink {
  id: string
  name: string
  url: string
  status: 'approved' | 'pending' | 'rejected'
  health: 'ok' | 'down'
  createdAt: string
}

const friends = ref<FriendLink[]>([
  { id: '1', name: 'DIYgod - DIYGOD', url: 'https://diygod.cc', status: 'approved', health: 'ok', createdAt: '2026-01-15' },
  { id: '2', name: 'Sukka\'s Blog', url: 'https://blog.skk.moe', status: 'approved', health: 'ok', createdAt: '2026-01-20' },
  { id: '3', name: '椎咲良田', url: 'https://sanshizhiduo.com', status: 'approved', health: 'ok', createdAt: '2026-02-05' },
  { id: '4', name: 'Innei Space', url: 'https://innei.in', status: 'approved', health: 'ok', createdAt: '2026-02-12' },
  { id: '5', name: '某新站', url: 'https://newsite.example.com', status: 'pending', health: 'ok', createdAt: '2026-03-12' },
  { id: '6', name: '已失效站点', url: 'https://gone.example.com', status: 'approved', health: 'down', createdAt: '2025-11-08' },
])

const approvedCount = computed(() => friends.value.filter((f) => f.status === 'approved').length)
const pendingCount = computed(() => friends.value.filter((f) => f.status === 'pending').length)
const downCount = computed(() => friends.value.filter((f) => f.health === 'down').length)

function statusClass(status: FriendLink['status']) {
  if (status === 'approved') return 'bg-emerald-100 text-emerald-700'
  if (status === 'rejected') return 'bg-red-100 text-red-600'
  return 'bg-amber-100 text-amber-700'
}

function statusLabel(status: FriendLink['status']) {
  if (status === 'approved') return '已通过'
  if (status === 'rejected') return '已拒绝'
  return '待审核'
}
</script>
