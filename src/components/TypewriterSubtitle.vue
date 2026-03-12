<template>
  <p class="mt-5 text-center text-2xl text-white/95 drop-shadow-[0_4px_16px_rgba(15,23,42,0.5)]">
    {{ visibleText }}<span class="animate-blink">|</span>
  </p>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'

interface Props {
  text: string
  speed?: number
}

const props = withDefaults(defineProps<Props>(), {
  speed: 110,
})

const visibleText = ref('')
let cursor = 0
let timer: ReturnType<typeof setInterval> | null = null

onMounted(() => {
  timer = setInterval(() => {
    cursor += 1
    visibleText.value = props.text.slice(0, cursor)
    if (cursor >= props.text.length && timer) {
      clearInterval(timer)
      timer = null
    }
  }, props.speed)
})

onBeforeUnmount(() => {
  if (timer) {
    clearInterval(timer)
  }
})
</script>
