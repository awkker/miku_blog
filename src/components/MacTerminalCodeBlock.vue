<template>
  <section class="mb-8 overflow-hidden rounded-2xl border border-slate-700/70 bg-[#1e1e1e] shadow-2xl">
    <header class="flex h-8 items-center justify-between bg-[#2d2d2d] px-3">
      <div class="flex items-center gap-2">
        <span class="h-3 w-3 rounded-full bg-[#ff5f56]" />
        <span class="h-3 w-3 rounded-full bg-[#ffbd2e]" />
        <span class="h-3 w-3 rounded-full bg-[#27c93f]" />
      </div>
      <p class="text-[11px] font-medium uppercase tracking-[0.18em] text-slate-400">{{ languageLabel }}</p>
      <button
        type="button"
        class="rounded-md px-2 py-1 text-[11px] text-slate-400 transition hover:bg-slate-600/50 hover:text-slate-100"
        @click="copyCode"
      >
        {{ copied ? '已复制' : '复制代码' }}
      </button>
    </header>

    <pre class="m-0 overflow-x-auto bg-[#1e1e1e] px-4 py-5 text-[0.9rem] leading-7 text-slate-100 [scrollbar-width:none] [&::-webkit-scrollbar]:hidden"><code class="font-mono" v-html="codeHtml" /></pre>
  </section>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

interface Props {
  language?: string
  codeHtml: string
  codeRaw: string
}

const props = withDefaults(defineProps<Props>(), {
  language: 'text',
})

const copied = ref(false)

const languageLabel = computed(() => {
  const aliasMap: Record<string, string> = {
    ts: 'TypeScript',
    js: 'JavaScript',
    vue: 'Vue',
    sh: 'Shell',
    bash: 'Bash',
    md: 'Markdown',
  }

  const normalized = props.language.toLowerCase()
  return aliasMap[normalized] ?? `${normalized.charAt(0).toUpperCase()}${normalized.slice(1)}`
})

const copyCode = async (): Promise<void> => {
  try {
    await navigator.clipboard.writeText(props.codeRaw)
    copied.value = true
    setTimeout(() => {
      copied.value = false
    }, 1000)
  } catch {
    copied.value = false
  }
}
</script>
