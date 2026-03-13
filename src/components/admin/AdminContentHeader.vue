<template>
  <div class="flex items-center justify-between gap-4">
    <div class="flex items-center gap-3">
      <button
        type="button"
        class="inline-flex h-9 w-9 items-center justify-center rounded-lg text-slate-500 transition hover:bg-slate-100 hover:text-slate-900 lg:hidden"
        aria-label="打开导航菜单"
        @click="toggleSidebar"
      >
        <svg viewBox="0 0 24 24" class="h-5 w-5 fill-none stroke-current stroke-[1.8]">
          <path d="M4 7h16M4 12h16M4 17h16" />
        </svg>
      </button>

      <h1 class="text-xl font-semibold text-slate-900">{{ pageTitle }}</h1>
    </div>

    <div class="flex items-center gap-3">
      <div class="hidden text-right sm:block">
        <p class="text-sm font-medium text-slate-900">{{ userName }}</p>
        <p class="text-xs text-slate-500">管理员</p>
      </div>
      <div
        class="flex h-9 w-9 items-center justify-center rounded-full bg-miku/15 text-sm font-semibold text-miku"
        aria-hidden="true"
      >
        {{ avatarLetter }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useStore } from '@nanostores/vue'

import { authState, hydrateAuth } from '../../stores/auth'
import { toggleSidebar } from '../../stores/ui'

interface Props {
  pageTitle?: string
}

withDefaults(defineProps<Props>(), {
  pageTitle: '仪表盘',
})

const auth = useStore(authState)
const mounted = ref(false)

const userName = computed(() => {
  if (!mounted.value) return 'Admin'
  return auth.value.user?.name ?? 'Admin'
})

const avatarLetter = computed(() => {
  const name = userName.value
  return name.charAt(0).toUpperCase()
})

onMounted(() => {
  mounted.value = true
  hydrateAuth()
})
</script>
