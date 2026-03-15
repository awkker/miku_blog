<template>
  <section class="space-y-6">
    <div class="flex flex-wrap items-center justify-between gap-4">
      <div class="flex items-center gap-2.5">
        <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-[#c084fc]/10">
          <svg viewBox="0 0 24 24" class="h-4 w-4 fill-none stroke-[#c084fc] stroke-[2]" aria-hidden="true">
            <path d="M10 13a5 5 0 007.07 0l2.12-2.12a5 5 0 00-7.07-7.07L10 5M14 11a5 5 0 00-7.07 0L4.81 13.12a5 5 0 007.07 7.07L14 19" />
          </svg>
        </div>
        <div>
          <h2 class="text-base font-semibold text-slate-800">友链墙</h2>
          <p class="text-xs text-slate-400">
            {{ links.length > 0 ? `已收录 ${links.length} 个站点` : '正在建设中' }}
          </p>
        </div>
      </div>
      <MikuButton
        variant="ghost"
        aria-label="申请交换友链"
        @click="showToast('友链申请通道将在下一阶段开放。', 'info')"
      >
        <svg viewBox="0 0 24 24" class="h-4 w-4 fill-none stroke-current stroke-[1.8]" aria-hidden="true">
          <path d="M12 5v14M5 12h14" />
        </svg>
        申请交换
      </MikuButton>
    </div>

    <ErrorState
      v-if="status === 'error'"
      :description="error || '友链读取失败，请稍后重试。'"
      @retry="loadLinks"
    />

    <div v-else>
      <div v-if="status === 'loading'" class="grid gap-5 grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
        <SkeletonCard v-for="item in 6" :key="item" />
      </div>

      <EmptyState
        v-else-if="links.length === 0"
        title="暂无友链"
        description="友链墙正在建设中，欢迎稍后再来。"
      />

      <div v-else class="grid gap-5 grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
        <FriendLinkCard
          v-for="item in links"
          :key="item.id"
          :friend="item"
        />
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { useStore } from '@nanostores/vue'
import { onMounted } from 'vue'

import { friendError, friendFetchStatus, friendLinks, loadFriendLinks } from '../../stores/friends'
import { showToast } from '../../stores/ui'
import EmptyState from '../ui/EmptyState.vue'
import ErrorState from '../ui/ErrorState.vue'
import FriendLinkCard from './FriendLinkCard.vue'
import MikuButton from '../ui/MikuButton.vue'
import SkeletonCard from '../ui/SkeletonCard.vue'

const links = useStore(friendLinks)
const status = useStore(friendFetchStatus)
const error = useStore(friendError)

async function loadLinks() {
  await loadFriendLinks()
}

onMounted(async () => {
  await loadLinks()
})
</script>
