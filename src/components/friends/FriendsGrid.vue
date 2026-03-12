<template>
  <section class="space-y-6">
    <div class="flex flex-wrap items-center justify-between gap-4">
      <p class="max-w-2xl text-sm leading-relaxed text-white/75">
        欢迎交换友链。请确保网站持续更新、可稳定访问，并在你的站点中添加本博客链接。
      </p>
      <MikuButton
        variant="ghost"
        aria-label="申请交换友链"
        @click="showToast('友链申请通道将在下一阶段开放。', 'info')"
      >
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
