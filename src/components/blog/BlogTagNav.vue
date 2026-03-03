<template>
  <LiquidGlassFrame
    class="tag-nav-glass"
    width="100%"
    max-width="none"
    padding="14px 18px"
    :border-radius="18"
    :displacement-strength="1.08"
    :edge-refraction-strength="1.25"
    :interactive="false"
  >
    <div class="tag-scroll-wrap">
      <div class="tag-scroll" role="tablist" aria-label="博客分类筛选">
        <button
          v-for="tag in tags"
          :key="tag"
          type="button"
          class="tag-item"
          :class="{ 'is-active': tag === activeTag }"
          role="tab"
          :aria-selected="tag === activeTag"
          @click="emit('select', tag)"
        >
          {{ tag }}
        </button>
      </div>
    </div>
  </LiquidGlassFrame>
</template>

<script setup lang="ts">
import LiquidGlassFrame from '@/components/LiquidGlassFrame.vue'

interface Props {
  tags: string[]
  activeTag: string
}

defineProps<Props>()

const emit = defineEmits<{
  select: [tag: string]
}>()
</script>

<style scoped>
.tag-scroll-wrap {
  position: relative;
}

.tag-scroll {
  display: flex;
  gap: 10px;
  overflow-x: auto;
  padding: 2px 2px 8px;
  -webkit-overflow-scrolling: touch;
  scrollbar-width: none;
  -ms-overflow-style: none;
  mask-image: linear-gradient(to right, transparent 0, black 18px, black calc(100% - 18px), transparent 100%);
}

.tag-scroll::-webkit-scrollbar {
  display: none;
  width: 0;
  height: 0;
}

.tag-item {
  border: 0;
  border-radius: 999px;
  padding: 8px 16px;
  font-family: var(--font-body);
  font-size: 0.92rem;
  color: rgba(20, 95, 82, 0.94);
  background: rgba(240, 255, 252, 0.72);
  white-space: nowrap;
  cursor: pointer;
  transition:
    transform 180ms ease,
    background-color 180ms ease,
    color 180ms ease,
    box-shadow 180ms ease;
}

.tag-item:hover {
  background: rgba(226, 252, 247, 0.92);
}

.tag-item:focus-visible {
  outline: none;
  box-shadow: 0 0 0 2px rgba(116, 211, 188, 0.36);
}

.tag-item.is-active {
  color: rgba(9, 65, 55, 0.98);
  background: rgba(208, 246, 236, 0.96);
  box-shadow:
    0 2px 9px rgba(16, 90, 74, 0.16),
    0 0 0 1px rgba(255, 255, 255, 0.44) inset;
  transform: translateY(-1px);
}
</style>
