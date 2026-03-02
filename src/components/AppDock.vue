<template>
  <LiquidGlassFrame
    class="mac-dock-glass"
    width="fit-content"
    max-width="none"
    padding="10px 15px"
    :border-radius="24"
    :displacement-strength="1.06"
    :edge-refraction-strength="1.25"
    :interactive="false"
  >
    <div class="mac-dock">
      <button
        v-for="item in dockItems"
        :key="item.label"
        class="dock-item"
        :class="{ active: isActive(item.to) }"
        :title="item.label"
        type="button"
        @click="navigate(item.to)"
      >
        <span class="dock-icon">{{ item.icon }}</span>
        <span class="dock-label">{{ item.label }}</span>
      </button>

      <button
        v-if="showCollapse"
        class="dock-item dock-collapse"
        type="button"
        title="收起导航"
        @click="emit('collapse')"
      >
        <span class="dock-icon">⌄</span>
        <span class="dock-label">收起</span>
      </button>
    </div>
  </LiquidGlassFrame>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'

import LiquidGlassFrame from '@/components/LiquidGlassFrame.vue'

interface DockItem {
  label: string
  icon: string
  to: string
}

interface Props {
  showCollapse?: boolean
}

withDefaults(defineProps<Props>(), {
  showCollapse: false,
})

const emit = defineEmits<{
  collapse: []
}>()

const router = useRouter()
const route = useRoute()

const dockItems: DockItem[] = [
  { label: '主页', icon: '⌂', to: '/' },
  { label: '博客', icon: '✦', to: '/blog' },
]

const navigate = (path: string): void => {
  if (route.path === path) {
    return
  }
  void router.push(path)
}

const isActive = (path: string): boolean => {
  if (path === '/blog') {
    return route.path.startsWith('/blog')
  }
  return route.path === path
}
</script>

<style scoped>
.mac-dock-glass {
  margin-bottom: 20px;
}

.mac-dock {
  display: flex;
  gap: 12px;
  align-items: flex-end;
}

.dock-item {
  border: 0;
  width: 50px;
  height: 50px;
  cursor: pointer;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.8);
  transition:
    transform 320ms cubic-bezier(0.22, 0.61, 0.36, 1),
    box-shadow 320ms ease,
    background-color 320ms ease;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 2px;
  color: rgba(13, 82, 69, 0.9);
  font-family: var(--font-body);
  transform-origin: center bottom;
  will-change: transform;
}

.dock-item:hover {
  transform: translateY(-8px) scale(1.16);
  box-shadow:
    0 12px 22px rgba(8, 38, 31, 0.18),
    inset 0 1px 0 rgba(255, 255, 255, 0.58);
}

.dock-item.active {
  background: rgba(226, 255, 248, 0.95);
}

.dock-collapse {
  background: linear-gradient(160deg, rgba(242, 255, 251, 0.9) 0%, rgba(222, 249, 240, 0.76) 100%);
}

.dock-collapse .dock-icon {
  font-size: 1.08rem;
  transform: translateY(-1px);
}

.dock-collapse .dock-label {
  font-size: 0.52rem;
}

.dock-icon {
  line-height: 1;
  font-size: 1rem;
}

.dock-label {
  font-size: 0.56rem;
  line-height: 1;
  transform: translateY(1px);
}
</style>
