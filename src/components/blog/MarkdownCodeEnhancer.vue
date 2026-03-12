<template>
  <div class="hidden" aria-hidden="true" />
</template>

<script setup lang="ts">
import { createApp, onBeforeUnmount, onMounted } from 'vue'

import MacTerminalCodeBlock from './MacTerminalCodeBlock.vue'

interface Props {
  containerSelector?: string
}

const props = withDefaults(defineProps<Props>(), {
  containerSelector: '#article-content',
})

const mountedApps: Array<{ unmount: () => void; element: HTMLElement }> = []

const inferLanguage = (className: string): string => {
  const match = className.match(/language-([\w-]+)/)
  return match?.[1] ?? 'text'
}

onMounted(() => {
  const container = document.querySelector(props.containerSelector)
  if (!container) {
    return
  }

  const preElements = Array.from(container.querySelectorAll('pre'))

  preElements.forEach((pre) => {
    if (pre.dataset.terminalized === 'true') {
      return
    }

    const codeEl = pre.querySelector('code')
    if (!codeEl) {
      return
    }

    pre.dataset.terminalized = 'true'
    const mountEl = document.createElement('div')
    pre.replaceWith(mountEl)

    const codeHtml = codeEl.innerHTML
    const codeRaw = codeEl.textContent ?? ''
    const language = inferLanguage(`${pre.className} ${codeEl.className}`)

    const app = createApp(MacTerminalCodeBlock, {
      language,
      codeHtml,
      codeRaw,
    })

    app.mount(mountEl)
    mountedApps.push({
      unmount: () => app.unmount(),
      element: mountEl,
    })
  })
})

onBeforeUnmount(() => {
  mountedApps.forEach(({ unmount, element }) => {
    unmount()
    element.remove()
  })
  mountedApps.length = 0
})
</script>
