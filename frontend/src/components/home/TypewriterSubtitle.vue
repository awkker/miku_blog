<template>
  <div class="subtitle-shell">
    <p class="subtitle-text">
      {{ visibleText }}<span class="subtitle-cursor animate-blink">|</span>
    </p>
  </div>
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

<style scoped>
.subtitle-shell {
  margin-top: 1.35rem;
  width: min(88vw, 700px);
  border: 1px solid rgba(240, 255, 252, 0.62);
  border-radius: 1.5rem;
  background: rgba(248, 254, 252, 0.34);
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.58),
    0 10px 28px rgba(8, 25, 32, 0.16);
}

.subtitle-text {
  margin: 0;
  padding: clamp(0.72rem, 1.25vw, 0.92rem) clamp(1rem, 2.2vw, 1.6rem);
  text-align: center;
  color: rgba(57, 197, 187, 0.9);
  font-family: 'Noto Serif SC', 'Noto Serif JP', 'Songti SC', 'STSong', serif;
  font-size: clamp(1rem, 1.35vw, 1.22rem);
  font-weight: 700;
  line-height: 1.4;
  letter-spacing: 0.02em;
}

.subtitle-cursor {
  opacity: 0.65;
}
</style>
