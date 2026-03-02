<template>
  <LiquidGlassFrame
    class="article-card"
    width="100%"
    max-width="none"
    padding="18px 22px 24px"
    :border-radius="24"
    :displacement-strength="1.08"
    :edge-refraction-strength="1.3"
    :interactive="false"
  >
    <div class="article-actions">
      <button type="button" class="reader-btn" @click="emit('back')">← 返回列表</button>
      <div class="reader-size">
        <button type="button" class="reader-btn" @click="emit('decreaseFontSize')">A-</button>
        <button type="button" class="reader-btn" @click="emit('resetFontSize')">A</button>
        <button type="button" class="reader-btn" @click="emit('increaseFontSize')">A+</button>
      </div>
    </div>

    <article class="markdown-body" :style="markdownStyle" v-html="html"></article>
  </LiquidGlassFrame>
</template>

<script setup lang="ts">
import { computed } from 'vue'

import LiquidGlassFrame from '@/components/LiquidGlassFrame.vue'

interface Props {
  html: string
  fontScale: number
}

const props = defineProps<Props>()

const emit = defineEmits<{
  back: []
  increaseFontSize: []
  decreaseFontSize: []
  resetFontSize: []
}>()

const markdownStyle = computed(() => ({
  '--reader-font-size': `${props.fontScale.toFixed(2)}rem`,
}))
</script>

<style scoped>
.article-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  margin-bottom: 14px;
  padding-bottom: 12px;
  border-bottom: 1px solid rgba(173, 219, 205, 0.5);
}

.reader-size {
  display: flex;
  gap: 8px;
}

.reader-btn {
  border: 0;
  border-radius: 10px;
  padding: 6px 10px;
  background: rgba(241, 255, 251, 0.78);
  color: rgba(10, 73, 60, 0.92);
  font-family: var(--font-body);
  cursor: pointer;
  transition: transform 180ms ease, background-color 180ms ease;
}

.reader-btn:hover {
  transform: translateY(-1px);
  background: rgba(222, 249, 241, 0.92);
}

.markdown-body {
  --reader-font-size: 1rem;
  font-size: var(--reader-font-size);
  line-height: 1.9;
  color: rgba(15, 67, 56, 0.92);
}

.markdown-body :deep(h1),
.markdown-body :deep(h2),
.markdown-body :deep(h3) {
  font-family: var(--font-title);
  color: rgba(8, 53, 44, 0.96);
  line-height: 1.45;
}

.markdown-body :deep(h1) {
  font-size: 2rem;
  margin: 24px 0 12px;
}

.markdown-body :deep(h2) {
  font-size: 1.45rem;
  margin: 22px 0 10px;
}

.markdown-body :deep(h3) {
  font-size: 1.16rem;
  margin: 18px 0 8px;
}

.markdown-body :deep(p) {
  margin: 0 0 14px;
}

.markdown-body :deep(ul),
.markdown-body :deep(ol) {
  margin: 0 0 14px;
  padding-left: 20px;
}

.markdown-body :deep(li) {
  margin: 6px 0;
}

.markdown-body :deep(code) {
  font-size: 0.92em;
  padding: 2px 6px;
  border-radius: 6px;
  background: rgba(226, 246, 239, 0.85);
}

.markdown-body :deep(pre) {
  margin: 0 0 16px;
  overflow: auto;
  padding: 12px 14px;
  border-radius: 12px;
  background: rgba(7, 31, 36, 0.9);
}

.markdown-body :deep(pre code) {
  padding: 0;
  background: transparent;
  color: rgba(236, 251, 247, 0.95);
}

.markdown-body :deep(a) {
  color: rgba(24, 122, 102, 0.94);
}

.markdown-body :deep(hr) {
  margin: 20px 0;
  border: 0;
  border-top: 1px solid rgba(161, 216, 200, 0.48);
}
</style>

