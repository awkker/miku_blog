<template>
  <section class="space-y-6">
    <LiquidGlassCard class="rounded-3xl p-6">
      <form class="space-y-4" @submit.prevent="handleSubmit">
        <div class="grid gap-4 md:grid-cols-2">
          <MikuInput
            v-model="form.nickname"
            label="昵称"
            placeholder="请输入昵称"
            :error="errors.nickname"
            aria-label="昵称"
            required
          />
          <MikuInput
            v-model="form.website"
            label="网址（可选）"
            placeholder="https://example.com"
            :error="errors.website"
            aria-label="个人网址"
          />
        </div>

        <MikuTextarea
          v-model="form.message"
          label="留言内容"
          placeholder="想说点什么？欢迎留下你的想法。"
          :error="errors.message"
          :rows="5"
          aria-label="留言内容"
          required
        />

        <div class="flex justify-end">
          <MikuButton
            type="submit"
            :loading="submitting"
            :disabled="submitting"
            aria-label="发送留言"
          >
            {{ submitting ? '发送中...' : '发送留言' }}
          </MikuButton>
        </div>
      </form>
    </LiquidGlassCard>

    <ErrorState
      v-if="fetchStatus === 'error'"
      :description="fetchError || '留言读取失败，请稍后再试。'"
      @retry="loadMessages"
    />

    <div v-else class="space-y-4">
      <div v-if="fetchStatus === 'loading'" class="grid gap-4 sm:grid-cols-2">
        <SkeletonCard v-for="item in 4" :key="item" />
      </div>

      <EmptyState
        v-else-if="messages.length === 0"
        title="还没有留言"
        description="成为第一个在这里留下足迹的人。"
      />

      <TransitionGroup
        v-else
        name="message-rise"
        tag="div"
        class="grid gap-4 sm:grid-cols-2"
      >
        <GuestbookMessageCard
          v-for="item in messages"
          :key="item.id"
          :message="item"
        />
      </TransitionGroup>
    </div>
  </section>
</template>

<script setup lang="ts">
import { useStore } from '@nanostores/vue'
import { computed, onMounted, reactive } from 'vue'

import {
  guestbookError,
  guestbookFetchStatus,
  guestbookMessages,
  guestbookSubmitStatus,
  loadGuestbookMessages,
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

const form = reactive({
  nickname: '',
  website: '',
  message: '',
})

const errors = reactive({
  nickname: '',
  website: '',
  message: '',
})

const messages = useStore(guestbookMessages)
const fetchStatus = useStore(guestbookFetchStatus)
const submitStatus = useStore(guestbookSubmitStatus)
const fetchError = useStore(guestbookError)

const submitting = computed(() => submitStatus.value === 'loading')

function isValidUrl(url: string) {
  if (!url.trim()) {
    return true
  }

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
  if (!validate()) {
    return
  }

  try {
    await submitGuestbookMessage({
      nickname: form.nickname,
      website: form.website,
      message: form.message,
    })

    form.message = ''
    showToast('留言已发送', 'success')
  } catch {
    showToast('留言发送失败，请稍后重试', 'error')
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
