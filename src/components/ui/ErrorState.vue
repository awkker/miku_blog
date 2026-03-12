<template>
  <LiquidGlassCard class="rounded-3xl border-red-200/50 p-6 text-center">
    <div class="mx-auto mb-3 flex h-12 w-12 items-center justify-center rounded-full border border-red-200/50 bg-red-300/15 text-red-100">
      <svg viewBox="0 0 24 24" class="h-6 w-6 fill-none stroke-current stroke-[1.8]">
        <path d="M12 8v5" />
        <path d="M12 16h.01" />
        <path d="M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3l-8.47-14.14a2 2 0 00-3.42 0z" />
      </svg>
    </div>
    <h3 class="text-lg font-semibold text-white">{{ title }}</h3>
    <p class="mt-2 text-sm text-white/75">{{ description }}</p>
    <MikuButton
      class="mt-4"
      variant="ghost"
      :aria-label="retryLabel"
      @click="$emit('retry')"
    >
      {{ retryLabel }}
    </MikuButton>
  </LiquidGlassCard>
</template>

<script setup lang="ts">
import LiquidGlassCard from './LiquidGlassCard.vue'
import MikuButton from './MikuButton.vue'

interface Props {
  title?: string
  description: string
  retryLabel?: string
}

withDefaults(defineProps<Props>(), {
  title: '加载失败',
  retryLabel: '重试',
})

defineEmits<{
  (e: 'retry'): void
}>()
</script>
