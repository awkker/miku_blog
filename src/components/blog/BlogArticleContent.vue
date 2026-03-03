<template>
  <LiquidGlassFrame
    class="article-card"
    width="100%"
    max-width="none"
    padding="18px 22px 24px"
    :border-radius="24"
    :displacement-strength="1.08"
    :edge-refraction-strength="1.3"
    :interactive="false"
  >
    <div class="article-actions">
      <button type="button" class="reader-btn" @click="emit('back')">← 返回列表</button>
      <div class="reader-size">
        <button type="button" class="reader-btn" @click="emit('decreaseFontSize')">A-</button>
        <button type="button" class="reader-btn" @click="emit('resetFontSize')">A</button>
        <button type="button" class="reader-btn" @click="emit('increaseFontSize')">A+</button>
      </div>
    </div>

    <article ref="articleRef" class="markdown-body" :style="markdownStyle" v-html="html"></article>
  </LiquidGlassFrame>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'

import LiquidGlassFrame from '@/components/LiquidGlassFrame.vue'

interface Props {
  html: string
  fontScale: number
}

const props = defineProps<Props>()

const emit = defineEmits<{
  back: []
  increaseFontSize: []
  decreaseFontSize: []
  resetFontSize: []
}>()

const articleRef = ref<HTMLElement | null>(null)
const resetLabelTimers = new Map<HTMLButtonElement, number>()

const markdownStyle = computed(() => ({
  '--reader-font-size': `${props.fontScale.toFixed(2)}rem`,
}))

function decodeCodePayload(payload: string): string {
  try {
    return decodeURIComponent(payload)
  } catch {
    return payload
  }
}

async function copyToClipboard(content: string): Promise<boolean> {
  if (typeof navigator !== 'undefined' && navigator.clipboard?.writeText) {
    try {
      await navigator.clipboard.writeText(content)
      return true
    } catch {
      return false
    }
  }
  if (typeof document === 'undefined') {
    return false
  }
  const textarea = document.createElement('textarea')
  textarea.value = content
  textarea.setAttribute('readonly', '')
  textarea.style.position = 'fixed'
  textarea.style.opacity = '0'
  document.body.appendChild(textarea)
  textarea.select()
  try {
    return document.execCommand('copy')
  } catch {
    return false
  } finally {
    document.body.removeChild(textarea)
  }
}

function updateCopyButtonLabel(button: HTMLButtonElement, label: string): void {
  button.textContent = label
  const previousTimer = resetLabelTimers.get(button)
  if (previousTimer !== undefined) {
    window.clearTimeout(previousTimer)
  }
  const timer = window.setTimeout(() => {
    button.textContent = '复制'
    resetLabelTimers.delete(button)
  }, 1400)
  resetLabelTimers.set(button, timer)
}

async function handleArticleClick(event: Event): Promise<void> {
  const target = event.target
  if (!(target instanceof HTMLElement)) {
    return
  }
  const copyButton = target.closest('.code-copy-btn')
  if (!(copyButton instanceof HTMLButtonElement)) {
    return
  }
  const payload = copyButton.dataset.code ?? ''
  const rawCode = decodeCodePayload(payload)
  const copied = await copyToClipboard(rawCode)
  updateCopyButtonLabel(copyButton, copied ? '已复制' : '复制失败')
}

function bindArticleEvents(): void {
  articleRef.value?.addEventListener('click', handleArticleClick)
}

function unbindArticleEvents(): void {
  articleRef.value?.removeEventListener('click', handleArticleClick)
}

onMounted(() => {
  bindArticleEvents()
})

onBeforeUnmount(() => {
  unbindArticleEvents()
  for (const timer of resetLabelTimers.values()) {
    window.clearTimeout(timer)
  }
  resetLabelTimers.clear()
})

watch(
  () => props.html,
  () => {
    unbindArticleEvents()
    nextTick(() => bindArticleEvents())
  },
)
</script>

<style scoped>
.article-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid rgba(173, 219, 205, 0.35);
}

.reader-size {
  display: flex;
  gap: 8px;
}

.reader-btn {
  border: 0;
  border-radius: 10px;
  padding: 5px 9px;
  background: rgba(241, 255, 251, 0.58);
  color: rgba(10, 73, 60, 0.78);
  font-size: 0.82rem;
  font-family: var(--font-body);
  cursor: pointer;
  transition: background-color 180ms ease;
}

.reader-btn:hover {
  background: rgba(226, 250, 243, 0.84);
}

.markdown-body {
  --reader-font-size: 1rem;
  font-size: var(--reader-font-size);
  line-height: 1.74;
  color: rgba(15, 67, 56, 0.92);
}

.markdown-body :deep(h1),
.markdown-body :deep(h2),
.markdown-body :deep(h3) {
  font-family: var(--font-title);
  color: rgba(8, 53, 44, 0.96);
  line-height: 1.45;
  scroll-margin-top: 106px;
}

.markdown-body :deep(h1) {
  font-size: 2rem;
  margin: 24px 0 12px;
}

.markdown-body :deep(h2) {
  font-size: 1.45rem;
  margin: 22px 0 10px;
}

.markdown-body :deep(h3) {
  font-size: 1.16rem;
  margin: 18px 0 8px;
}

.markdown-body :deep(p) {
  margin: 0 0 18px;
}

.markdown-body :deep(ul),
.markdown-body :deep(ol) {
  margin: 0 0 18px;
  padding-left: 20px;
}

.markdown-body :deep(li) {
  margin: 6px 0;
}

.markdown-body :deep(code) {
  font-size: 0.92em;
  padding: 2px 6px;
  border-radius: 6px;
  background: rgba(226, 246, 239, 0.85);
  font-family: var(--font-code);
}

.markdown-body :deep(.code-block) {
  margin: 0 0 22px;
  border-radius: 14px;
  overflow: hidden;
  border: 1px solid rgba(106, 149, 162, 0.28);
  background: linear-gradient(180deg, rgba(15, 24, 39, 0.98) 0%, rgba(12, 19, 34, 0.99) 100%);
  box-shadow:
    0 16px 30px rgba(9, 20, 35, 0.26),
    0 2px 0 rgba(255, 255, 255, 0.05) inset;
}

.markdown-body :deep(.code-toolbar) {
  height: 38px;
  display: grid;
  grid-template-columns: auto 1fr auto;
  align-items: center;
  gap: 10px;
  padding: 0 10px;
  border-bottom: 1px solid rgba(138, 174, 186, 0.2);
  background: linear-gradient(180deg, rgba(31, 43, 63, 0.98) 0%, rgba(22, 34, 51, 0.96) 100%);
}

.markdown-body :deep(.window-dots) {
  display: flex;
  align-items: center;
  gap: 6px;
}

.markdown-body :deep(.dot) {
  width: 10px;
  height: 10px;
  border-radius: 999px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.35) inset;
}

.markdown-body :deep(.dot.red) {
  background: #ff5f57;
}

.markdown-body :deep(.dot.yellow) {
  background: #febc2e;
}

.markdown-body :deep(.dot.green) {
  background: #28c840;
}

.markdown-body :deep(.code-lang) {
  justify-self: center;
  font-family: var(--font-code);
  font-size: 0.75rem;
  letter-spacing: 0.03em;
  color: rgba(198, 219, 255, 0.72);
}

.markdown-body :deep(.code-copy-btn) {
  border: 0;
  border-radius: 8px;
  padding: 5px 10px;
  background: rgba(82, 169, 228, 0.16);
  color: rgba(198, 229, 255, 0.92);
  font-size: 0.75rem;
  font-family: var(--font-body);
  cursor: pointer;
  transition: background-color 160ms ease;
}

.markdown-body :deep(.code-copy-btn:hover) {
  background: rgba(97, 188, 250, 0.24);
}

.markdown-body :deep(.code-block pre) {
  margin: 0;
  padding: 10px 0 12px;
  overflow: auto;
}

.markdown-body :deep(.code-block code) {
  display: block;
  padding: 0;
  background: transparent;
  color: rgba(231, 240, 255, 0.92);
  font-family: var(--font-code);
  font-size: 0.88rem;
  line-height: 1.68;
}

.markdown-body :deep(.code-line) {
  display: grid;
  grid-template-columns: 44px 1fr;
  align-items: baseline;
}

.markdown-body :deep(.line-no) {
  padding: 0 10px 0 0;
  text-align: right;
  user-select: none;
  color: rgba(131, 153, 184, 0.62);
}

.markdown-body :deep(.line-content) {
  padding: 0 14px 0 0;
  white-space: pre;
}

.markdown-body :deep(.token-keyword) {
  color: #c792ea;
}

.markdown-body :deep(.token-string) {
  color: #ecc48d;
}

.markdown-body :deep(.token-number) {
  color: #f78c6c;
}

.markdown-body :deep(.token-comment) {
  color: #7f8ca5;
}

.markdown-body :deep(.token-function) {
  color: #82aaff;
}

.markdown-body :deep(.token-key) {
  color: #7dd3fc;
}

.markdown-body :deep(.token-literal) {
  color: #89ddff;
}

.markdown-body :deep(.token-operator) {
  color: #b8c3db;
}

.markdown-body :deep(a) {
  color: rgba(24, 122, 102, 0.94);
}

.markdown-body :deep(hr) {
  margin: 20px 0;
  border: 0;
  border-top: 1px solid rgba(161, 216, 200, 0.48);
}
</style>
