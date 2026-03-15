<template>
  <section class="space-y-5">
    <LiquidGlassCard padding="24px">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 class="text-2xl font-semibold text-slate-900">文章管理</h1>
          <p class="mt-1 text-sm text-slate-600">管理博客文章的发布、草稿与分类。</p>
        </div>
        <MikuButton variant="solid" aria-label="新建文章">+ 新建文章</MikuButton>
      </div>
    </LiquidGlassCard>

    <LiquidGlassCard padding="0px">
      <div class="overflow-x-auto">
        <table class="w-full text-left text-sm">
          <thead>
            <tr class="border-b border-slate-200/60">
              <th class="px-5 py-3.5 font-semibold text-slate-700">标题</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700">分类</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700">状态</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700 text-right">浏览量</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700">发布时间</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700 text-center">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="post in posts"
              :key="post.id"
              class="border-b border-slate-100/60 transition hover:bg-white/40"
            >
              <td class="px-5 py-3.5 font-medium text-slate-900">{{ post.title }}</td>
              <td class="px-5 py-3.5 text-slate-600">{{ post.category }}</td>
              <td class="px-5 py-3.5">
                <span
                  class="inline-block rounded-full px-2.5 py-0.5 text-xs font-medium"
                  :class="statusClass(post.status)"
                >
                  {{ statusLabel(post.status) }}
                </span>
              </td>
              <td class="px-5 py-3.5 text-right font-mono text-slate-700">{{ post.views.toLocaleString() }}</td>
              <td class="px-5 py-3.5 text-slate-600">{{ post.publishedAt }}</td>
              <td class="px-5 py-3.5 text-center">
                <div class="flex items-center justify-center gap-2">
                  <button
                    type="button"
                    class="rounded-xl border border-slate-200/80 bg-white/50 px-2.5 py-1 text-xs text-slate-700 transition hover:border-miku/40 hover:text-miku"
                    aria-label="编辑文章"
                  >
                    编辑
                  </button>
                  <button
                    type="button"
                    class="rounded-xl border border-red-200/80 bg-white/50 px-2.5 py-1 text-xs text-red-600 transition hover:bg-red-50"
                    aria-label="删除文章"
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
import { ref } from 'vue'

import LiquidGlassCard from '../ui/LiquidGlassCard.vue'
import MikuButton from '../ui/MikuButton.vue'

interface Post {
  id: string
  title: string
  category: string
  status: 'published' | 'draft' | 'scheduled'
  views: number
  publishedAt: string
}

const posts = ref<Post[]>([
  { id: '1', title: 'Astro + Vue Islands 实战指南', category: '前端技术', status: 'published', views: 2340, publishedAt: '2026-03-10' },
  { id: '2', title: '用 Go + Hertz 构建高性能博客后端', category: '后端架构', status: 'published', views: 1876, publishedAt: '2026-03-08' },
  { id: '3', title: 'Tailwind CSS 液态玻璃设计实录', category: 'UI/UX', status: 'published', views: 1523, publishedAt: '2026-03-05' },
  { id: '4', title: 'PostgreSQL 全文检索优化笔记', category: '数据库', status: 'draft', views: 0, publishedAt: '--' },
  { id: '5', title: 'Redis Lua 限流方案详解', category: '后端架构', status: 'scheduled', views: 0, publishedAt: '2026-03-15' },
  { id: '6', title: 'Nano Stores 跨框架状态共享', category: '前端技术', status: 'published', views: 968, publishedAt: '2026-02-28' },
])

function statusClass(status: Post['status']) {
  if (status === 'published') return 'bg-emerald-100 text-emerald-700'
  if (status === 'draft') return 'bg-slate-100 text-slate-600'
  return 'bg-purple-100 text-purple-700'
}

function statusLabel(status: Post['status']) {
  if (status === 'published') return '已发布'
  if (status === 'draft') return '草稿'
  return '定时发布'
}
</script>
