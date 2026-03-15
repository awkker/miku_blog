<template>
  <div class="glass-layer w-16 rounded-2xl p-3">
    <button
      type="button"
      class="flex w-full flex-col items-center gap-2 rounded-xl py-2 text-slate-600 transition hover:bg-white/25"
      @click="toggleLike"
    >
      <svg
        viewBox="0 0 24 24"
        class="h-6 w-6 transition"
        :class="[
          pulse ? 'animate-heartbeat' : '',
          liked ? 'fill-miku text-miku' : 'fill-transparent text-slate-500',
        ]"
      >
        <path
          d="M12.1 21.35l-1.1-1.01C5.14 14.99 2 12.05 2 8.5 2 5.57 4.42 3.2 7.4 3.2c1.74 0 3.4.82 4.5 2.12 1.1-1.3 2.76-2.12 4.5-2.12C19.58 3.2 22 5.57 22 8.5c0 3.55-3.14 6.49-8.9 11.85l-1 .99z"
          stroke="currentColor"
          stroke-width="1.5"
        />
      </svg>
      <span class="text-xs font-semibold" :class="liked ? 'text-miku' : 'text-slate-600'">{{ count }}</span>
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

interface Props {
  initialCount?: number
}

const props = withDefaults(defineProps<Props>(), {
  initialCount: 88,
})

const count = ref(props.initialCount)
const liked = ref(false)
const pulse = ref(false)

const toggleLike = (): void => {
  liked.value = !liked.value
  count.value += liked.value ? 1 : -1
  pulse.value = true
  setTimeout(() => {
    pulse.value = false
  }, 320)
}
</script>
