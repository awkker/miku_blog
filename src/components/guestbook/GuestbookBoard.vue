<template>
  <section class="space-y-6">
    <!-- Post Form -->
    <LiquidGlassCard padding="24px 28px">
      <div class="mb-5 flex items-center gap-2.5">
        <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-miku/10">
          <svg viewBox="0 0 24 24" class="h-4 w-4 fill-none stroke-miku stroke-[2]" aria-hidden="true">
            <path d="M12 20h9M16.5 3.5a2.121 2.121 0 013 3L7 19l-4 1 1-4 12.5-12.5z" />
          </svg>
        </div>
        <h2 class="text-lg font-semibold text-slate-800">发布留言</h2>
      </div>
      <form class="space-y-4" @submit.prevent="handleSubmit">
        <div class="grid gap-4 md:grid-cols-2">
          <MikuInput v-model="form.nickname" label="昵称" placeholder="请输入昵称" :error="errors.nickname" aria-label="昵称" required />
          <MikuInput v-model="form.website" label="网址（可选）" placeholder="https://example.com" :error="errors.website" aria-label="个人网址" />
        </div>
        <MikuTextarea v-model="form.message" label="留言内容" placeholder="想说点什么？欢迎留下你的想法。" :error="errors.message" :rows="4" aria-label="留言内容" required />
        <div class="flex justify-end">
          <MikuButton type="submit" :loading="submitting" :disabled="submitting" aria-label="发送留言">
            {{ submitting ? '发送中...' : '发布' }}
          </MikuButton>
        </div>
      </form>
    </LiquidGlassCard>

    <!-- Sort Tabs + Count -->
    <div class="flex flex-wrap items-center justify-between gap-3">
      <div class="flex items-center gap-1.5 rounded-xl border border-slate-200 bg-white/60 p-1">
        <button
          v-for="tab in sortTabs"
          :key="tab.value"
          type="button"
          class="rounded-lg px-3 py-1.5 text-xs font-medium transition"
          :class="currentSort === tab.value ? 'bg-miku text-white shadow-sm' : 'text-slate-500 hover:text-slate-700'"
          @click="changeSort(tab.value)"
        >
          {{ tab.label }}
        </button>
      </div>
      <span v-if="messages.length > 0" class="text-xs text-slate-400">
        {{ messages.length }} 条留言
      </span>
    </div>

    <!-- Messages -->
    <ErrorState v-if="fetchStatus === 'error'" :description="fetchError || '留言读取失败，请稍后再试。'" @retry="loadMessages" />

    <div v-else class="space-y-4">
      <div v-if="fetchStatus === 'loading'" class="space-y-4">
        <SkeletonCard v-for="item in 3" :key="item" />
      </div>

      <EmptyState v-else-if="messages.length === 0" title="还没有留言" description="成为第一个在这里留下足迹的人。" />

      <TransitionGroup v-else name="message-rise" tag="div" class="space-y-4">
        <GuestbookMessageCard
          v-for="item in sortedMessages"
          :key="item.id"
          :message="item"
          @reply="handleReply"
        />
      </TransitionGroup>
    </div>
  </section>
</template>

<script setup lang="ts">
import { useStore } from '@nanostores/vue'
import { computed, onMounted, reactive } from 'vue'

import {
  type SortMode,
  guestbookError,
  guestbookFetchStatus,
  guestbookMessages,
  guestbookSortMode,
  guestbookSorted,
  guestbookSubmitStatus,
  loadGuestbookMessages,
  setSortMode,
  submitGuestbookMessage,
} from '../../stores/guestbook'
import { showToast } from '../../stores/ui'
import EmptyState from '../ui/EmptyState.vue'
import ErrorState from '../ui/ErrorState.vue'
import GuestbookMessageCard from './GuestbookMessageCard.vue'
import LiquidGlassCard from '../ui/LiquidGlassCard.vue'
import MikuButton from '../ui/MikuButton.vue'
import MikuInput from '../ui/MikuInput.vue'
import MikuTextarea from '../ui/MikuTextarea.vue'
import SkeletonCard from '../ui/SkeletonCard.vue'

const sortTabs = [
  { label: '最热', value: 'hot' as SortMode },
  { label: '最新', value: 'newest' as SortMode },
  { label: '最早', value: 'oldest' as SortMode },
]

const form = reactive({ nickname: '', website: '', message: '' })
const errors = reactive({ nickname: '', website: '', message: '' })

const messages = useStore(guestbookMessages)
const sortedMessages = useStore(guestbookSorted)
const fetchStatus = useStore(guestbookFetchStatus)
const submitStatus = useStore(guestbookSubmitStatus)
const fetchError = useStore(guestbookError)
const currentSort = useStore(guestbookSortMode)

const submitting = computed(() => submitStatus.value === 'loading')

function changeSort(mode: SortMode) {
  setSortMode(mode)
}

function isValidUrl(url: string) {
  if (!url.trim()) return true
  try {
    const parsed = new URL(url)
    return parsed.protocol === 'http:' || parsed.protocol === 'https:'
  } catch {
    return false
  }
}

function validate() {
  errors.nickname = form.nickname.trim() ? '' : '昵称不能为空'
  errors.message = form.message.trim() ? '' : '留言内容不能为空'
  errors.website = isValidUrl(form.website) ? '' : '网址格式不合法，请输入 http(s) 链接'
  return !errors.nickname && !errors.message && !errors.website
}

async function loadMessages() {
  await loadGuestbookMessages()
}

async function handleSubmit() {
  if (!validate()) return
  try {
    await submitGuestbookMessage({ nickname: form.nickname, website: form.website, message: form.message })
    form.message = ''
    showToast('留言已发送', 'success')
  } catch {
    showToast('留言发送失败，请稍后重试', 'error')
  }
}

async function handleReply(payload: { parentId: string; nickname: string; message: string }) {
  try {
    await submitGuestbookMessage({ nickname: payload.nickname, message: payload.message, parentId: payload.parentId })
    showToast('回复已发送', 'success')
  } catch {
    showToast('回复发送失败，请稍后重试', 'error')
  }
}

onMounted(async () => {
  await loadMessages()
})
</script>

<style scoped>
.message-rise-enter-active,
.message-rise-leave-active {
  transition: transform 0.35s ease, opacity 0.35s ease;
}
.message-rise-enter-from,
.message-rise-leave-to {
  opacity: 0;
  transform: translateY(10px);
}
.message-rise-move {
  transition: transform 0.35s ease;
}
</style>
