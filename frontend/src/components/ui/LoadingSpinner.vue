<template>
  <svg
    class="animate-spin text-current"
    :class="sizeClass"
    viewBox="0 0 24 24"
    fill="none"
    aria-hidden="true"
  >
    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
    <path
      class="opacity-90"
      fill="currentColor"
      d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z"
    />
  </svg>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  size?: 'sm' | 'md' | 'lg'
}

const props = withDefaults(defineProps<Props>(), {
  size: 'md',
})

const sizeClass = computed(() => {
  if (props.size === 'sm') {
    return 'h-4 w-4'
  }

  if (props.size === 'lg') {
    return 'h-7 w-7'
  }

  return 'h-5 w-5'
})
</script>
