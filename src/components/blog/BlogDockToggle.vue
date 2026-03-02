<template>
  <div class="dock-layer">
    <Transition name="dock-switch" mode="out-in">
      <LiquidGlassFrame
        v-if="!isOpen"
        key="peek"
        class="dock-peek"
        width="68px"
        max-width="none"
        padding="7px 0 11px"
        :border-radius="17"
        :displacement-strength="1.04"
        :edge-refraction-strength="1.2"
        :interactive="false"
      >
        <button class="dock-handle-btn" type="button" aria-label="展开底部导航" @click="emit('show')">
          <span class="dock-handle-grip"></span>
          <span class="dock-handle-arrow">⌃</span>
        </button>
      </LiquidGlassFrame>

      <div v-else key="open" class="dock-open-state">
        <div class="dock-wrap">
          <AppDock show-collapse @collapse="emit('hide')" />
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import AppDock from '@/components/AppDock.vue'
import LiquidGlassFrame from '@/components/LiquidGlassFrame.vue'

interface Props {
  isOpen: boolean
}

defineProps<Props>()

const emit = defineEmits<{
  show: []
  hide: []
}>()
</script>

<style scoped>
.dock-layer {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 6;
  display: flex;
  justify-content: center;
  pointer-events: none;
}

.dock-open-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  pointer-events: none;
}

.dock-peek {
  pointer-events: auto;
  transform: translateY(11px);
  transition:
    transform 360ms cubic-bezier(0.22, 0.61, 0.36, 1),
    box-shadow 360ms cubic-bezier(0.22, 0.61, 0.36, 1);
  box-shadow:
    0 8px 24px rgba(6, 33, 28, 0.2),
    inset 0 1px 0 rgba(255, 255, 255, 0.48);
  animation: dock-handle-breathe 2.8s ease-in-out infinite;
}

.dock-peek:hover {
  transform: translateY(7px);
  box-shadow:
    0 10px 26px rgba(6, 33, 28, 0.22),
    inset 0 1px 0 rgba(255, 255, 255, 0.56);
  animation-play-state: paused;
}

.dock-handle-btn {
  width: 48px;
  height: 28px;
  border: 0;
  background: transparent;
  color: rgba(14, 83, 70, 0.96);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  gap: 1px;
  padding: 0;
}

.dock-handle-grip {
  width: 18px;
  height: 3px;
  border-radius: 999px;
  background: linear-gradient(90deg, rgba(147, 220, 203, 0.68) 0%, rgba(108, 198, 177, 0.9) 100%);
  box-shadow: 0 0 6px rgba(102, 205, 170, 0.26);
}

.dock-handle-arrow {
  font-size: 0.98rem;
  line-height: 1;
  transform: translateY(-1px);
  color: rgba(14, 83, 70, 0.9);
}

.dock-wrap {
  position: static;
  z-index: 5;
  display: flex;
  justify-content: center;
  pointer-events: none;
  animation: dock-panel-rise 420ms cubic-bezier(0.22, 0.61, 0.36, 1) both;
}

.dock-wrap :deep(.mac-dock-glass) {
  pointer-events: auto;
  margin-bottom: 8px;
}

@keyframes dock-handle-breathe {
  0%,
  100% {
    transform: translateY(11px);
  }
  50% {
    transform: translateY(8px);
  }
}

@keyframes dock-panel-rise {
  from {
    transform: translateY(16px) scale(0.96);
    opacity: 0;
  }
  to {
    transform: translateY(0) scale(1);
    opacity: 1;
  }
}

.dock-switch-enter-active,
.dock-switch-leave-active {
  transition:
    transform 320ms cubic-bezier(0.22, 0.61, 0.36, 1),
    opacity 320ms ease;
}

.dock-switch-enter-from,
.dock-switch-leave-to {
  opacity: 0;
  transform: translateY(12px) scale(0.96);
}
</style>

