<template>
  <span class="tabular-nums text-[11px] tracking-wide text-[#39c5bb]/90">{{ clockText }}</span>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'

const clockText = ref('')
let timer: ReturnType<typeof setInterval> | null = null

const formatter = new Intl.DateTimeFormat('zh-CN', {
  month: '2-digit',
  day: '2-digit',
  hour: '2-digit',
  minute: '2-digit',
  weekday: 'short',
  hour12: false,
})

const updateClock = (): void => {
  clockText.value = formatter.format(new Date())
}

onMounted(() => {
  updateClock()
  timer = setInterval(updateClock, 1000)
})

onBeforeUnmount(() => {
  if (timer) {
    clearInterval(timer)
  }
})
</script>
