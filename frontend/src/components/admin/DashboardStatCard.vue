<template>
  <LiquidGlassCard padding="20px">
    <div class="flex items-start justify-between gap-3">
      <div>
        <p class="text-xs uppercase tracking-[0.2em] text-slate-600">{{ label }}</p>
        <p class="mt-2 font-mono text-3xl font-semibold text-slate-900">{{ value }}</p>
      </div>
      <span class="rounded-xl border border-white/30 bg-white/20 p-2 text-miku" aria-hidden="true">
        <svg viewBox="0 0 24 24" class="h-4 w-4 fill-none stroke-current stroke-[1.8]">
          <path :d="iconPath" />
        </svg>
      </span>
    </div>

    <p
      v-if="trend"
      class="mt-3 text-xs"
      :class="trend.startsWith('-') ? 'text-red-700' : 'text-emerald-700'"
    >
      {{ trend }}
    </p>
  </LiquidGlassCard>
</template>

<script setup lang="ts">
import { computed } from 'vue'

import LiquidGlassCard from '../ui/LiquidGlassCard.vue'

interface Props {
  label: string
  value: string
  trend?: string
  icon?: 'article' | 'comment' | 'like' | 'link'
}

const props = withDefaults(defineProps<Props>(), {
  trend: '',
  icon: 'article',
})

const iconPath = computed(() => {
  if (props.icon === 'comment') {
    return 'M7 10h10M7 14h6M5 4h14a2 2 0 012 2v9a2 2 0 01-2 2H9l-4 3v-3H5a2 2 0 01-2-2V6a2 2 0 012-2z'
  }

  if (props.icon === 'like') {
    return 'M12 21s-7-4.35-7-10a4 4 0 017-2.65A4 4 0 0119 11c0 5.65-7 10-7 10z'
  }

  if (props.icon === 'link') {
    return 'M10 13a5 5 0 007.07 0l2.12-2.12a5 5 0 10-7.07-7.07L10 5M14 11a5 5 0 00-7.07 0L4.81 13.12a5 5 0 007.07 7.07L14 19'
  }

  return 'M7 4h10a2 2 0 012 2v12a2 2 0 01-2 2H7a2 2 0 01-2-2V6a2 2 0 012-2zm2 4h6M9 12h6M9 16h4'
})
</script>
