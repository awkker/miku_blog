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
    <div class="toc-progress-track" aria-hidden="true">
      <div class="toc-progress-bar" :style="{ width: `${normalizedProgress}%` }"></div>
    </div>
    <h3>目录导航</h3>
    <div ref="tocScrollRef" class="toc-scroll-wrap">
      <ul class="toc-list">
        <li v-for="heading in headings" :key="heading.id">
          <button
            type="button"
            class="toc-link"
            :class="[`level-${heading.level}`, { 'is-active': heading.id === activeHeadingId }]"
            :data-heading-id="heading.id"
            @click="emit('scrollToHeading', heading.id)"
          >
            {{ heading.text }}
          </button>
        </li>
      </ul>
    </div>
    <button type="button" class="reader-btn top-btn" @click="emit('scrollToTop')">回到顶部</button>
  </LiquidGlassFrame>
</template>

<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue'

import LiquidGlassFrame from '@/components/LiquidGlassFrame.vue'
import type { MarkdownHeading } from '@/utils/markdown'

interface Props {
  headings: MarkdownHeading[]
  activeHeadingId: string
  readingProgress: number
}

const props = defineProps<Props>()

const emit = defineEmits<{
  scrollToHeading: [id: string]
  scrollToTop: []
}>()

const tocScrollRef = ref<HTMLElement | null>(null)

const normalizedProgress = computed(() => Math.min(100, Math.max(0, props.readingProgress)))

function ensureActiveItemVisible(): void {
  const container = tocScrollRef.value
  if (!container || !props.activeHeadingId) {
    return
  }
  if (container.scrollHeight <= container.clientHeight + 1) {
    return
  }
  const activeButton = container.querySelector<HTMLButtonElement>('.toc-link.is-active')
  if (!activeButton) {
    return
  }
  const containerRect = container.getBoundingClientRect()
  const activeRect = activeButton.getBoundingClientRect()
  const safePadding = 12

  if (activeRect.top < containerRect.top + safePadding) {
    const delta = containerRect.top + safePadding - activeRect.top
    container.scrollTo({
      top: Math.max(0, container.scrollTop - delta),
      behavior: 'smooth',
    })
    return
  }

  if (activeRect.bottom > containerRect.bottom - safePadding) {
    const delta = activeRect.bottom - (containerRect.bottom - safePadding)
    container.scrollTo({
      top: container.scrollTop + delta,
      behavior: 'smooth',
    })
  }
}

watch(
  () => props.activeHeadingId,
  () => {
    nextTick(() => ensureActiveItemVisible())
  },
)
</script>

<style scoped>
.toc-progress-track {
  width: 100%;
  height: 3px;
  border-radius: 999px;
  margin: -6px 0 12px;
  overflow: hidden;
  background: rgba(126, 180, 163, 0.24);
}

.toc-progress-bar {
  width: 0;
  height: 100%;
  border-radius: inherit;
  background: linear-gradient(90deg, rgba(124, 226, 199, 0.95) 0%, rgba(93, 197, 173, 0.95) 100%);
  box-shadow: 0 0 8px rgba(94, 201, 175, 0.38);
  transition: width 180ms linear;
}

.toc-card h3 {
  margin: 0;
  font-size: 1rem;
  color: rgba(9, 54, 45, 0.94);
}

.toc-scroll-wrap {
  margin: 12px 0 14px;
  max-height: min(50vh, 420px);
  overflow-y: auto;
  scrollbar-width: thin;
  scrollbar-color: rgba(116, 184, 165, 0.48) rgba(202, 237, 227, 0.28);
}

.toc-scroll-wrap::-webkit-scrollbar {
  width: 7px;
}

.toc-scroll-wrap::-webkit-scrollbar-track {
  background: rgba(202, 237, 227, 0.2);
  border-radius: 999px;
}

.toc-scroll-wrap::-webkit-scrollbar-thumb {
  background: rgba(116, 184, 165, 0.48);
  border-radius: 999px;
}

.toc-list {
  margin: 0;
  padding: 0;
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: 7px;
}

.toc-link {
  position: relative;
  width: 100%;
  border: 0;
  border-radius: 10px;
  padding: 7px 8px 7px 12px;
  text-align: left;
  font-family: var(--font-body);
  background: rgba(242, 255, 251, 0.72);
  color: rgba(14, 78, 65, 0.9);
  cursor: pointer;
  transition:
    background-color 180ms ease,
    color 180ms ease,
    font-weight 180ms ease;
}

.toc-link::before {
  content: '';
  position: absolute;
  left: 0;
  top: 8px;
  bottom: 8px;
  width: 3px;
  border-radius: 999px;
  background: transparent;
}

.toc-link:hover {
  background: rgba(224, 249, 241, 0.86);
}

.toc-link.is-active {
  background: rgba(198, 242, 230, 0.92);
  color: rgba(8, 66, 55, 0.98);
  font-weight: 600;
}

.toc-link.is-active::before {
  background: rgba(var(--miku-color-rgb), 0.96);
  box-shadow: 0 0 6px rgba(var(--miku-color-rgb), 0.46);
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
