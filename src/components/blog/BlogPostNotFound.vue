<template>
  <LiquidGlassFrame
    class="not-found-card"
    width="100%"
    max-width="none"
    padding="28px 22px"
    :border-radius="22"
    :displacement-strength="1.08"
    :edge-refraction-strength="1.3"
    :interactive="false"
  >
    <h2>文章不存在</h2>
    <p>可能已被移动或 slug 不正确。</p>
    <button type="button" class="not-found-btn" @click="emit('back')">返回博客列表</button>
  </LiquidGlassFrame>
</template>

<script setup lang="ts">
import LiquidGlassFrame from '@/components/LiquidGlassFrame.vue'

const emit = defineEmits<{
  back: []
}>()
</script>

<style scoped>
.not-found-card h2 {
  margin: 0 0 8px;
  color: rgba(10, 67, 55, 0.95);
}

.not-found-card p {
  margin: 0 0 14px;
  color: rgba(20, 87, 72, 0.88);
}

.not-found-btn {
  border: 0;
  border-radius: 10px;
  padding: 6px 10px;
  background: rgba(241, 255, 251, 0.78);
  color: rgba(10, 73, 60, 0.92);
  font-family: var(--font-body);
  cursor: pointer;
  transition: transform 180ms ease, background-color 180ms ease;
}

.not-found-btn:hover {
  transform: translateY(-1px);
  background: rgba(222, 249, 241, 0.92);
}
</style>
