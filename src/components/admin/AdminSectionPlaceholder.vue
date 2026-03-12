<template>
  <section class="space-y-5">
    <LiquidGlassCard class="rounded-3xl p-6">
      <h1 class="text-2xl font-semibold text-white">{{ title }}</h1>
      <p class="mt-2 text-sm text-white/75">{{ description }}</p>
    </LiquidGlassCard>

    <div class="grid gap-4 lg:grid-cols-3">
      <LiquidGlassCard
        v-for="item in cards"
        :key="item.title"
        class="rounded-3xl p-5"
      >
        <h2 class="text-base font-semibold text-white">{{ item.title }}</h2>
        <p class="mt-2 text-sm text-white/70">{{ item.desc }}</p>
      </LiquidGlassCard>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'

import LiquidGlassCard from '../ui/LiquidGlassCard.vue'

interface Props {
  title: string
  description: string
  module: 'posts' | 'comments' | 'friends'
}

const props = defineProps<Props>()

const cardMap = {
  posts: [
    { title: '文章列表', desc: '后续接入发布状态、分类筛选和全文检索。' },
    { title: '草稿箱', desc: '支持草稿保存、自动恢复与版本记录。' },
    { title: '发布计划', desc: '预留定时发布与回滚入口。' },
  ],
  comments: [
    { title: '待审核评论', desc: '支持批量通过、驳回与敏感词标记。' },
    { title: '举报处理', desc: '统一管理用户举报与风险内容。' },
    { title: '黑名单策略', desc: '预留 IP/关键词规则配置面板。' },
  ],
  friends: [
    { title: '友链列表', desc: '维护已通过友链与展示排序。' },
    { title: '申请队列', desc: '集中处理交换申请与状态通知。' },
    { title: '校验任务', desc: '定期检测友链可用性与内容质量。' },
  ],
}

const cards = computed(() => cardMap[props.module])
</script>
