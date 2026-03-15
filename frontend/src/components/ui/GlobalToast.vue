<template>
  <Transition name="toast-up">
    <div
      v-if="toast.open"
      class="fixed bottom-5 right-5 z-50 rounded-2xl border px-4 py-3 text-sm text-white shadow-lg backdrop-blur-md"
      :class="toast.type === 'error'
        ? 'border-red-300/50 bg-red-400/25'
        : toast.type === 'success'
          ? 'border-emerald-300/50 bg-emerald-400/25'
          : 'border-white/30 bg-white/20'"
      role="status"
      aria-live="polite"
    >
      {{ toast.message }}
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { useStore } from '@nanostores/vue'

import { toastState } from '../../stores/ui'

const toast = useStore(toastState)
</script>

<style scoped>
.toast-up-enter-active,
.toast-up-leave-active {
  transition: transform 0.3s ease, opacity 0.3s ease;
}

.toast-up-enter-from,
.toast-up-leave-to {
  opacity: 0;
  transform: translateY(8px);
}
</style>
