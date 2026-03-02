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
    <div class="tag-scroll">
      <button v-for="tag in tags" :key="tag" class="tag-item">{{ tag }}</button>
    </div>
  </LiquidGlassFrame>
</template>

<script setup lang="ts">
import LiquidGlassFrame from '@/components/LiquidGlassFrame.vue'

interface Props {
  tags: string[]
}

defineProps<Props>()
</script>

<style scoped>
.tag-scroll {
  display: flex;
  gap: 10px;
  overflow-x: auto;
  padding-bottom: 4px;
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
}

.tag-item:hover {
  background: rgba(226, 252, 247, 0.92);
}
</style>

