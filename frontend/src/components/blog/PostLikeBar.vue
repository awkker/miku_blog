<template>
  <section class="glass-layer rounded-[22px] p-4">
    <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <div class="space-y-1">
        <p class="text-xs font-semibold tracking-[0.14em] text-slate-500">互动</p>
        <p class="text-sm text-slate-600">喜欢这篇文章的话，先点个赞，再到评论区聊聊你的看法。</p>
      </div>

      <div class="flex flex-wrap items-center gap-2">
        <button
          type="button"
          :disabled="!postId || liking"
          :class="[
            'inline-flex items-center gap-2 rounded-xl border px-4 py-2 text-sm font-semibold transition disabled:cursor-not-allowed disabled:opacity-70',
            likedState
              ? 'border-red-200 bg-red-50 text-red-500'
              : 'border-slate-200 bg-white/80 text-slate-600 hover:border-red-200 hover:text-red-500',
          ]"
          @click="toggleLike"
        >
          <svg
            viewBox="0 0 24 24"
            :class="[
              'h-4 w-4 stroke-[1.6] transition',
              likedState ? 'fill-red-500 stroke-red-500' : 'fill-none stroke-current',
            ]"
            aria-hidden="true"
          >
            <path d="M20.84 4.61a5.5 5.5 0 00-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 00-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 000-7.78z" />
          </svg>
          {{ likedState ? '已点赞' : '点赞' }} · {{ likeCountState }}
        </button>

        <a
          href="#post-comments"
          class="inline-flex items-center gap-1.5 rounded-xl border border-miku/35 bg-miku-soft px-4 py-2 text-sm font-semibold text-miku transition hover:border-miku/55"
          aria-label="跳转到评论区"
        >
          去评论区
          <span class="rounded-full border border-miku/35 bg-white/70 px-1.5 py-0.5 text-[11px] text-miku">
            {{ commentCount }}
          </span>
        </a>
      </div>
    </div>

    <p v-if="message" class="mt-2 text-xs text-red-500">{{ message }}</p>
  </section>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

import { api } from '../../lib/api'

interface Props {
  postId: string
  likeCount: number
  liked: boolean
  commentCount: number
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'like-updated', payload: { liked: boolean; likeCount: number }): void
}>()

const likedState = ref(false)
const likeCountState = ref(0)
const liking = ref(false)
const message = ref('')
let lastActionAt = 0

watch(
  () => [props.postId, props.likeCount, props.liked],
  () => {
    likedState.value = Boolean(props.liked)
    likeCountState.value = Number(props.likeCount || 0)
    message.value = ''
  },
  { immediate: true },
)

async function toggleLike() {
  if (!props.postId || liking.value) return

  const now = Date.now()
  if (now - lastActionAt < 500) return
  lastActionAt = now

  liking.value = true
  message.value = ''

  try {
    const result = await api.post<{ liked: boolean }>(`/posts/${encodeURIComponent(props.postId)}/like`)
    const nextLiked = Boolean(result?.liked)
    if (nextLiked !== likedState.value) {
      likeCountState.value = Math.max(0, likeCountState.value + (nextLiked ? 1 : -1))
    }
    likedState.value = nextLiked
    emit('like-updated', {
      liked: likedState.value,
      likeCount: likeCountState.value,
    })
  } catch {
    message.value = '点赞失败，请稍后重试。'
  } finally {
    liking.value = false
  }
}
</script>
