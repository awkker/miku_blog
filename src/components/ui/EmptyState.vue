<template>
  <LiquidGlassCard class="rounded-3xl p-6 text-center">
    <div class="mx-auto mb-3 flex h-12 w-12 items-center justify-center rounded-full border border-white/25 bg-white/15 text-miku">
      <svg viewBox="0 0 24 24" class="h-6 w-6 fill-none stroke-current stroke-[1.8]">
        <path d="M4 12h16M12 4v16" />
      </svg>
    </div>
    <h3 class="text-lg font-semibold text-white">{{ title }}</h3>
    <p class="mt-2 text-sm text-white/70">{{ description }}</p>
    <MikuButton
      v-if="actionLabel"
      class="mt-4"
      variant="ghost"
      :aria-label="actionLabel"
      @click="$emit('action')"
    >
      {{ actionLabel }}
    </MikuButton>
  </LiquidGlassCard>
</template>

<script setup lang="ts">
import LiquidGlassCard from './LiquidGlassCard.vue'
import MikuButton from './MikuButton.vue'

interface Props {
  title: string
  description: string
  actionLabel?: string
}

withDefaults(defineProps<Props>(), {
  actionLabel: undefined,
})

defineEmits<{
  (e: 'action'): void
}>()
</script>
