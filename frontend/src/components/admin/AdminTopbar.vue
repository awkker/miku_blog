<template>
  <div class="space-y-3">
    <LiquidGlassCard padding="12px">
      <div class="flex items-center gap-3">
        <button
          type="button"
          class="inline-flex h-10 w-10 items-center justify-center rounded-2xl border border-slate-300/80 bg-white/60 text-slate-900 transition duration-300 hover:bg-white/80 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-miku/70 lg:hidden"
          aria-label="打开导航菜单"
          @click="toggleSidebar"
        >
          <svg viewBox="0 0 24 24" class="h-5 w-5 fill-none stroke-current stroke-[1.8]">
            <path d="M4 7h16M4 12h16M4 17h16" />
          </svg>
        </button>

        <div class="min-w-0 flex-1">
          <nav class="flex items-center gap-1.5 text-sm" aria-label="面包屑导航">
            <span class="text-slate-500">Nanamiku Admin</span>
            <span class="text-slate-400">/</span>
            <span class="font-semibold text-slate-900">{{ pageTitle }}</span>
          </nav>
        </div>

        <div class="hidden items-center gap-2 md:flex">
          <a
            href="/"
            class="inline-flex items-center justify-center rounded-2xl border border-slate-300/80 bg-white/60 px-3 py-2 text-sm text-slate-900 transition duration-300 hover:border-miku/40 hover:text-miku"
            aria-label="前往前台首页"
          >
            前台首页
          </a>
          <div class="rounded-2xl border border-slate-300/80 bg-white/60 px-3 py-2 text-right">
            <p class="text-sm font-semibold text-slate-900">{{ userName }}</p>
            <p class="text-xs text-slate-600">管理员</p>
          </div>

          <MikuButton
            variant="ghost"
            class="!border-slate-300/80 !bg-white/65 !text-slate-900 hover:!bg-white/80"
            aria-label="退出登录"
            @click="handleLogout"
          >
            退出
          </MikuButton>
        </div>
      </div>
    </LiquidGlassCard>

  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

import { authState, hydrateAuth, logout } from '../../stores/auth'
import { useStore } from '@nanostores/vue'
import { toggleSidebar } from '../../stores/ui'
import LiquidGlassCard from '../ui/LiquidGlassCard.vue'
import MikuButton from '../ui/MikuButton.vue'

interface Props {
  pageTitle?: string
}

withDefaults(defineProps<Props>(), {
  pageTitle: '仪表盘',
})

const auth = useStore(authState)
const mounted = ref(false)

const userName = computed(() => {
  if (!mounted.value) {
    return 'Admin'
  }

  return auth.value.user?.name ?? 'Admin'
})

onMounted(() => {
  mounted.value = true
  hydrateAuth()
})

function handleLogout() {
  logout()
  window.location.replace('/login')
}
</script>
