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
                  >
                    通过
                  </button>
                  <button
                    v-if="comment.status === 'pending'"
                    type="button"
                    class="rounded-xl border border-red-200/80 bg-white/50 px-2.5 py-1 text-xs text-red-600 transition hover:bg-red-50"
                    aria-label="驳回评论"
                  >
                    驳回
                  </button>
                  <button
                    v-if="comment.status !== 'pending'"
                    type="button"
                    class="rounded-xl border border-slate-200/80 bg-white/50 px-2.5 py-1 text-xs text-slate-600 transition hover:bg-slate-50"
                    aria-label="删除评论"
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

interface Comment {
  id: string
  author: string
  email: string
  content: string
  postTitle: string
  status: 'pending' | 'approved' | 'rejected'
  createdAt: string
}

const comments = ref<Comment[]>([
  { id: '1', author: '小明', email: 'ming@example.com', content: '这篇 Astro Islands 的文章写得太棒了，收藏了！', postTitle: 'Astro + Vue Islands 实战指南', status: 'pending', createdAt: '2026-03-13 14:22' },
  { id: '2', author: 'DevFan', email: 'devfan@example.com', content: '请问 Hertz 和 Gin 性能对比数据方便分享一下吗？', postTitle: '用 Go + Hertz 构建高性能博客后端', status: 'pending', createdAt: '2026-03-13 11:05' },
  { id: '3', author: '路过的猫', email: 'cat@example.com', content: '液态玻璃效果确实很惊艳，学到了', postTitle: 'Tailwind CSS 液态玻璃设计实录', status: 'approved', createdAt: '2026-03-12 20:38' },
  { id: '4', author: 'Spammer', email: 'spam@bad.com', content: '免费领取xxx，点击链接...', postTitle: 'Astro + Vue Islands 实战指南', status: 'rejected', createdAt: '2026-03-12 09:14' },
  { id: '5', author: '阿水', email: 'shui@example.com', content: 'Nano Stores 比 Pinia 轻量太多了，适合 Islands 架构', postTitle: 'Nano Stores 跨框架状态共享', status: 'pending', createdAt: '2026-03-11 16:47' },
])

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
