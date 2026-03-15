<template>
  <span class="jump-letter__container">
    <span
      v-for="(char, i) in chars"
      :key="i"
      class="jump-letter__alphabet"
      @mouseenter="onEnter"
      @animationend="onAnimEnd"
    >{{ char }}</span>
  </span>
</template>

<script setup lang="ts">
interface Props {
  text: string
}

const props = defineProps<Props>()

const chars = [...props.text]

function onEnter(e: Event) {
  const el = e.target as HTMLElement
  if (el.classList.contains('is-active')) return
  el.classList.add('is-active')
}

function onAnimEnd(e: AnimationEvent) {
  ;(e.target as HTMLElement).classList.remove('is-active')
}
</script>

<style scoped>
.jump-letter__container {
  display: inline-flex;
}

.jump-letter__alphabet {
  display: inline-block;
  transform-origin: center bottom;
  cursor: default;
}

.jump-letter__alphabet.is-active {
  animation: jump-letter 1.5s 0s cubic-bezier(0.165, 0.85, 0.45, 1);
}

@keyframes jump-letter {
  0%, 100% {
    transform: translateY(0) scale(1, 1);
    text-shadow: none;
  }
  25%, 75% {
    transform: translateY(0) scale(1.2, 0.8);
    text-shadow: 0 -3px 2px rgba(192, 132, 252, 0.5);
  }
  50% {
    transform: translateY(-24px) scale(0.8, 1.2);
    text-shadow: 0 -8px 4px rgba(192, 132, 252, 0.4);
    color: #c084fc;
  }
}
</style>
