<template>
  <LiquidGlassCard class="rounded-3xl p-4 transition duration-300 hover:shadow-[0_10px_28px_rgba(15,23,42,0.3)]">
    <article class="flex items-start gap-3">
      <img
        :src="message.avatar"
        :alt="`${message.nickname} 的头像`"
        class="h-12 w-12 shrink-0 rounded-full border border-white/35 object-cover"
        loading="lazy"
      />

      <div class="min-w-0 flex-1 space-y-2">
        <div class="flex flex-wrap items-center gap-x-3 gap-y-1">
          <p class="text-sm font-semibold text-white">{{ message.nickname }}</p>
          <p class="text-xs text-white/60">{{ message.createdAt }}</p>
        </div>

        <p class="whitespace-pre-wrap break-words text-sm leading-relaxed text-white/85">{{ message.message }}</p>

        <a
          v-if="message.website"
          :href="message.website"
          target="_blank"
          rel="noopener noreferrer"
          class="inline-flex items-center gap-1 text-xs text-miku transition hover:text-miku/80"
        >
          访问主页
          <svg viewBox="0 0 24 24" class="h-3.5 w-3.5 fill-none stroke-current stroke-[1.8]">
            <path d="M7 17L17 7M9 7h8v8" />
          </svg>
        </a>
      </div>
    </article>
  </LiquidGlassCard>
</template>

<script setup lang="ts">
import type { GuestbookMessage } from '../../stores/guestbook'
import LiquidGlassCard from '../ui/LiquidGlassCard.vue'

interface Props {
  message: GuestbookMessage
}

defineProps<Props>()
</script>
