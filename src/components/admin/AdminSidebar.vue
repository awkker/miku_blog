<template>
  <aside class="hidden w-64 shrink-0 lg:block">
    <LiquidGlassCard padding="16px" class="sticky top-6">
      <p class="px-2 text-xs uppercase tracking-[0.2em] text-slate-700">Nanamiku Admin</p>

      <nav class="mt-4 space-y-1">
        <a
          v-for="item in navItems"
          :key="item.key"
          :href="item.href"
          class="flex items-center gap-3 rounded-2xl px-3 py-2.5 text-sm transition duration-300"
          :class="item.key === activeKey
            ? 'border border-miku/40 bg-miku-soft text-miku'
            : 'text-slate-800 hover:bg-white/45 hover:text-slate-900'"
          :aria-label="item.label"
        >
          <span class="h-2 w-1 rounded-full" :class="item.key === activeKey ? 'bg-miku' : 'bg-transparent'" />
          <span>{{ item.label }}</span>
        </a>
      </nav>
    </LiquidGlassCard>
  </aside>

  <Transition name="drawer-fade">
    <div
      v-if="mobileOpen"
      class="fixed inset-0 z-40 bg-slate-950/45 backdrop-blur-sm lg:hidden"
      @click="setSidebarOpen(false)"
    />
  </Transition>

  <Transition name="drawer-slide">
    <aside
      v-if="mobileOpen"
      class="fixed left-3 top-3 z-50 w-[78vw] max-w-xs lg:hidden"
    >
      <LiquidGlassCard padding="16px">
        <div class="mb-4 flex items-center justify-between">
          <p class="text-sm font-semibold text-slate-900">导航菜单</p>
          <button
            type="button"
            class="rounded-xl border border-slate-300/80 bg-white/50 p-1.5 text-slate-900 transition hover:bg-white/75"
            aria-label="关闭导航菜单"
            @click="setSidebarOpen(false)"
          >
            <svg viewBox="0 0 24 24" class="h-4 w-4 fill-none stroke-current stroke-[1.8]">
              <path d="M6 6l12 12M18 6L6 18" />
            </svg>
          </button>
        </div>

        <nav class="space-y-1">
          <a
            v-for="item in navItems"
            :key="`mobile-${item.key}`"
            :href="item.href"
            class="flex items-center gap-3 rounded-2xl px-3 py-2.5 text-sm transition duration-300"
            :class="item.key === activeKey
              ? 'border border-miku/40 bg-miku-soft text-miku'
              : 'text-slate-800 hover:bg-white/45 hover:text-slate-900'"
            :aria-label="item.label"
            @click="setSidebarOpen(false)"
          >
            <span class="h-2 w-1 rounded-full" :class="item.key === activeKey ? 'bg-miku' : 'bg-transparent'" />
            <span>{{ item.label }}</span>
          </a>
        </nav>
      </LiquidGlassCard>
    </aside>
  </Transition>
</template>

<script setup lang="ts">
import { useStore } from '@nanostores/vue'

import { sidebarOpen, setSidebarOpen } from '../../stores/ui'
import LiquidGlassCard from '../ui/LiquidGlassCard.vue'

interface Props {
  activeKey?: string
}

withDefaults(defineProps<Props>(), {
  activeKey: 'dashboard',
})

const navItems = [
  { key: 'dashboard', label: '仪表盘 Dashboard', href: '/admin' },
  { key: 'posts', label: '文章管理', href: '/admin/posts' },
  { key: 'comments', label: '评论审核', href: '/admin/comments' },
  { key: 'friends', label: '友链管理', href: '/admin/friends' },
]

const mobileOpen = useStore(sidebarOpen)
</script>

<style scoped>
.drawer-fade-enter-active,
.drawer-fade-leave-active {
  transition: opacity 0.3s ease;
}

.drawer-fade-enter-from,
.drawer-fade-leave-to {
  opacity: 0;
}

.drawer-slide-enter-active,
.drawer-slide-leave-active {
  transition: transform 0.35s ease, opacity 0.35s ease;
}

.drawer-slide-enter-from,
.drawer-slide-leave-to {
  opacity: 0;
  transform: translateX(-12px);
}
</style>
