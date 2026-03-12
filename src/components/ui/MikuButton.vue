<template>
  <button
    :type="type"
    :disabled="isDisabled"
    :aria-label="ariaLabel"
    class="inline-flex items-center justify-center gap-2 rounded-2xl border px-4 py-2.5 text-sm font-semibold transition duration-300 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-miku/70 focus-visible:ring-offset-2 focus-visible:ring-offset-transparent active:scale-[0.98]"
    :class="[
      fullWidth ? 'w-full' : 'w-auto',
      variant === 'solid'
        ? 'border-transparent bg-miku text-slate-900 shadow-lg shadow-miku/30 hover:bg-miku/90 disabled:bg-miku/45 disabled:text-slate-700'
        : 'border-white/35 bg-white/15 text-slate-100 hover:bg-white/25 disabled:text-slate-400',
      isDisabled ? 'cursor-not-allowed' : '',
    ]"
  >
    <LoadingSpinner v-if="loading" size="sm" />
    <slot />
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue'

import LoadingSpinner from './LoadingSpinner.vue'

interface Props {
  type?: 'button' | 'submit' | 'reset'
  variant?: 'solid' | 'ghost'
  loading?: boolean
  disabled?: boolean
  fullWidth?: boolean
  ariaLabel?: string
}

const props = withDefaults(defineProps<Props>(), {
  type: 'button',
  variant: 'solid',
  loading: false,
  disabled: false,
  fullWidth: false,
  ariaLabel: undefined,
})

const isDisabled = computed(() => props.loading || props.disabled)
</script>
