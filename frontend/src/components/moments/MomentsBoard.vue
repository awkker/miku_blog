<template>
  <section class="space-y-6">
    <!-- Feed Count -->
    <div v-if="list.length > 0" class="flex items-center justify-between px-1">
      <span class="text-xs text-slate-400">{{ list.length }}{{ copy.countSuffix }}</span>
    </div>

    <!-- Feed -->
    <ErrorState v-if="fetchStatus === 'error'" :description="fetchError || copy.loadErrorFallback" @retry="loadFeed" />

    <div v-else class="space-y-4">
      <div v-if="fetchStatus === 'loading'" class="space-y-4">
        <SkeletonCard v-for="item in 3" :key="item" />
      </div>

      <EmptyState v-else-if="list.length === 0" :title="copy.emptyTitle" :description="copy.emptyDescription" />

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
import { siteCopy } from '../../content/copy'

const list = useStore(moments)
const fetchStatus = useStore(momentsFetchStatus)
const fetchError = useStore(momentsError)
const copy = siteCopy.components.momentsBoard

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
