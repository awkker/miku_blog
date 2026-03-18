<template>
  <div class="grid grid-cols-3 gap-2.5 text-center">
    <div class="rounded-xl border border-slate-100 bg-white/50 px-2 py-2.5 transition duration-300 hover:border-miku/20 hover:shadow-sm">
      <p class="text-lg font-bold text-miku">{{ loaded ? totalPosts : '--' }}</p>
      <p class="text-[11px] text-slate-500">文章</p>
    </div>
    <div class="rounded-xl border border-slate-100 bg-white/50 px-2 py-2.5 transition duration-300 hover:border-miku/20 hover:shadow-sm">
      <p class="text-lg font-bold text-miku">{{ loaded ? totalCategories : '--' }}</p>
      <p class="text-[11px] text-slate-500">系列</p>
    </div>
    <div class="rounded-xl border border-slate-100 bg-white/50 px-2 py-2.5 transition duration-300 hover:border-miku/20 hover:shadow-sm">
      <p class="text-lg font-bold text-miku">{{ loaded ? totalViews : '--' }}</p>
      <p class="text-[11px] text-slate-500">访问</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'

import { api, type PagedData } from '../../lib/api'

interface PostStat {
  category: string
  view_count: number
}

const totalPosts = ref(0)
const totalCategories = ref(0)
const totalViews = ref(0)
const loaded = ref(false)

async function loadStats() {
  try {
    const data = await api.get<PagedData<PostStat>>('/posts?page=1&size=200')
    totalPosts.value = data.total || 0
    const categories = new Set<string>()
    let views = 0
    for (const item of data.items || []) {
      if (item.category) categories.add(item.category)
      views += Number(item.view_count) || 0
    }
    totalCategories.value = categories.size
    totalViews.value = views
    loaded.value = true
  } catch {
    loaded.value = true
  }
}

onMounted(() => {
  loadStats()
})
</script>
