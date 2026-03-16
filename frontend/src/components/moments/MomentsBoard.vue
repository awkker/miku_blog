<template>
  <section class="space-y-6">
    <!-- Feed Count -->
    <div v-if="list.length > 0" class="flex items-center justify-between px-1">
      <span class="text-xs text-slate-400">{{ list.length }} 条说说</span>
    </div>

    <!-- Feed -->
    <ErrorState v-if="fetchStatus === 'error'" :description="fetchError || '说说加载失败，请稍后再试。'" @retry="loadFeed" />

    <div v-else class="space-y-4">
      <div v-if="fetchStatus === 'loading'" class="space-y-4">
        <SkeletonCard v-for="item in 3" :key="item" />
      </div>

      <EmptyState v-else-if="list.length === 0" title="还没有说说" description="发布第一条说说，记录此刻的想法。" />

      <TransitionGroup v-else name="feed-rise" tag="div" class="space-y-4">
        <MomentCard v-for="item in list" :key="item.id" :moment="item" />
      </TransitionGroup>
    </div>
  </section>
</template>

<script setup lang="ts">
import { useStore } from '@nanostores/vue'
import { onMounted } from 'vue'

import {
  loadMoments,
  moments,
  momentsError,
  momentsFetchStatus,
} from '../../stores/moments'
import EmptyState from '../ui/EmptyState.vue'
import ErrorState from '../ui/ErrorState.vue'
import MomentCard from './MomentCard.vue'
import SkeletonCard from '../ui/SkeletonCard.vue'

const list = useStore(moments)
const fetchStatus = useStore(momentsFetchStatus)
const fetchError = useStore(momentsError)

async function loadFeed() {
  await loadMoments()
}

onMounted(async () => {
  await loadFeed()
})
</script>

<style scoped>
.feed-rise-enter-active,
.feed-rise-leave-active {
  transition: transform 0.35s ease, opacity 0.35s ease;
}
.feed-rise-enter-from,
.feed-rise-leave-to {
  opacity: 0;
  transform: translateY(10px);
}
.feed-rise-move {
  transition: transform 0.35s ease;
}
</style>
