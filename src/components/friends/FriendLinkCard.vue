<template>
  <a
    :href="friend.url"
    target="_blank"
    rel="noopener noreferrer"
    class="block transition duration-300 hover:scale-105"
    :aria-label="`访问友链 ${friend.name}`"
  >
    <LiquidGlassCard class="rounded-3xl p-4 transition duration-300 hover:border-miku/45 hover:shadow-xl">
      <div class="flex items-center gap-3">
        <img
          :src="friend.avatar"
          :alt="`${friend.name} logo`"
          class="h-14 w-14 rounded-2xl border border-white/35 object-cover"
          loading="lazy"
        />

        <div class="min-w-0 flex-1">
          <h3 class="truncate text-base font-semibold text-white">{{ friend.name }}</h3>
          <p class="mt-1 max-h-10 overflow-hidden text-sm text-white/70">{{ friend.description }}</p>
          <p class="mt-2 truncate font-mono text-xs text-white/55">{{ friend.domain }}</p>
        </div>
      </div>
    </LiquidGlassCard>
  </a>
</template>

<script setup lang="ts">
import type { FriendLink } from '../../stores/friends'
import LiquidGlassCard from '../ui/LiquidGlassCard.vue'

interface Props {
  friend: FriendLink
}

defineProps<Props>()
</script>
