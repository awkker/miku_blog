<template>
  <LiquidGlassFrame
    class="toc-card"
    width="100%"
    max-width="none"
    padding="18px 16px"
    :border-radius="20"
    :displacement-strength="1.08"
    :edge-refraction-strength="1.3"
    :interactive="false"
  >
    <h3>目录导航</h3>
    <ul class="toc-list">
      <li v-for="heading in headings" :key="heading.id">
        <button
          type="button"
          class="toc-link"
          :class="`level-${heading.level}`"
          @click="emit('scrollToHeading', heading.id)"
        >
          {{ heading.text }}
        </button>
      </li>
    </ul>
    <button type="button" class="reader-btn top-btn" @click="emit('scrollToTop')">回到顶部</button>
  </LiquidGlassFrame>
</template>

<script setup lang="ts">
import LiquidGlassFrame from '@/components/LiquidGlassFrame.vue'
import type { MarkdownHeading } from '@/utils/markdown'

interface Props {
  headings: MarkdownHeading[]
}

defineProps<Props>()

const emit = defineEmits<{
  scrollToHeading: [id: string]
  scrollToTop: []
}>()
</script>

<style scoped>
.toc-card h3 {
  margin: 0;
  font-size: 1rem;
  color: rgba(9, 54, 45, 0.94);
}

.toc-list {
  margin: 12px 0 14px;
  padding: 0;
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.toc-link {
  width: 100%;
  border: 0;
  border-radius: 10px;
  padding: 6px 8px;
  text-align: left;
  font-family: var(--font-body);
  background: rgba(242, 255, 251, 0.72);
  color: rgba(14, 78, 65, 0.9);
  cursor: pointer;
  transition: background-color 180ms ease;
}

.toc-link:hover {
  background: rgba(224, 249, 241, 0.86);
}

.toc-link.level-3 {
  margin-left: 12px;
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

.top-btn {
  width: 100%;
}
</style>

