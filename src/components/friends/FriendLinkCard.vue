<template>
  <a
    :href="friend.url"
    target="_blank"
    rel="noopener noreferrer"
    class="group block transition duration-300 hover:scale-[1.02]"
    :aria-label="`访问友链 ${friend.name}`"
  >
    <LiquidGlassCard padding="20px" class="h-full transition duration-300 group-hover:shadow-[0_14px_36px_rgba(15,23,42,0.35)]">
      <div class="flex items-start gap-4">
        <img
          :src="friend.avatar"
          :alt="`${friend.name} logo`"
          class="h-14 w-14 shrink-0 rounded-2xl border border-slate-200 object-cover ring-2 ring-transparent transition duration-300 group-hover:ring-miku/30"
          loading="lazy"
        />

        <div class="min-w-0 flex-1">
          <h3 class="truncate text-base font-semibold text-slate-800 transition duration-300 group-hover:text-miku">{{ friend.name }}</h3>
          <p class="mt-1.5 line-clamp-2 text-sm leading-relaxed text-slate-500">{{ friend.description }}</p>
          <p class="mt-2.5 inline-flex items-center gap-1 truncate text-xs text-slate-400 transition duration-300 group-hover:text-slate-500">
            <svg viewBox="0 0 24 24" class="h-3 w-3 shrink-0 fill-none stroke-current stroke-[1.8]" aria-hidden="true">
              <circle cx="12" cy="12" r="10" /><path d="M2 12h20M12 2a15.3 15.3 0 014 10 15.3 15.3 0 01-4 10 15.3 15.3 0 01-4-10 15.3 15.3 0 014-10z" />
            </svg>
            {{ friend.domain }}
          </p>
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
