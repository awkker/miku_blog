<template>
  <LiquidGlassFrame
    class="quote-glass"
    max-width="800px"
    padding="30px 50px"
    :border-radius="24"
    :displacement-strength="1.1"
    :edge-refraction-strength="1.35"
    :interactive="false"
  >
    <p
      v-for="(line, index) in typedLines"
      :key="index"
      class="miku-text"
      :class="{ typing: activeLineIndex === index }"
    >
      {{ line }}
    </p>
  </LiquidGlassFrame>
</template>

<script setup lang="ts">
import LiquidGlassFrame from '@/components/LiquidGlassFrame.vue'
import { useTypewriter } from '@/composables/useTypewriter'

// 可替换为多行文案，打字机会按顺序逐行输出
const quoteLines = ['「月が綺麗ですね, 風も優しいですね」。']

// 复用打字机逻辑：返回当前已打出的文本和光标所在行
const { typedLines, activeLineIndex } = useTypewriter(quoteLines, {
  startDelay: 380,
  charDelay: 68,
  lineDelay: 320,
})
</script>

<style scoped>
.quote-glass {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.miku-text {
  font-family: var(--font-body);
  margin: 0;
  min-height: 1.5em;
  font-size: 1.2rem;
  font-weight: 500;
  letter-spacing: 1px;
  color: var(--miku-color);
  text-shadow: 0 1px 2px rgba(255, 255, 255, 0.5);
}

.miku-text.typing::after {
  content: '|';
  display: inline-block;
  margin-left: 4px;
  color: rgba(255, 255, 255, 0.92);
  animation: caret-blink 0.9s steps(1) infinite;
}

@keyframes caret-blink {
  50% { opacity: 0; }
}
</style>
