<template>
  <div class="flex items-center justify-between gap-4 rounded-2xl border border-white/30 bg-white/20 px-4 py-2.5 backdrop-blur-md">
    <div class="flex items-center gap-2" aria-hidden="true">
      <span class="h-3 w-3 rounded-full bg-[#ff5f57]" />
      <span class="h-3 w-3 rounded-full bg-[#febc2e]" />
      <span class="h-3 w-3 rounded-full bg-[#28c840]" />
    </div>
    <div class="min-w-0 text-center">
      <p class="truncate text-sm font-semibold text-white/90">{{ title }}</p>
      <p v-if="subtitle" class="truncate text-xs text-white/60">{{ subtitle }}</p>
    </div>
    <div class="h-3 w-14" aria-hidden="true" />
  </div>
</template>

<script setup lang="ts">
interface Props {
  title: string
  subtitle?: string
}

withDefaults(defineProps<Props>(), {
  subtitle: '',
})
</script>
