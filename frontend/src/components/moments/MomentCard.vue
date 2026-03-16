<template>
  <LiquidGlassCard padding="20px 24px" class="transition duration-300">
    <article class="space-y-3">
      <!-- Author Header -->
      <div class="flex items-center gap-3">
        <img
          :src="moment.avatar"
          :alt="moment.nickname"
          class="h-10 w-10 shrink-0 rounded-full object-cover"
          loading="lazy"
        />
        <div class="min-w-0 flex-1">
          <span class="text-sm font-semibold text-slate-700">{{ moment.nickname }}</span>
          <p class="text-[11px] text-slate-400">{{ moment.createdAt }}</p>
        </div>
      </div>

      <!-- Content -->
      <p class="whitespace-pre-wrap break-words text-[15px] leading-relaxed text-slate-600">{{ moment.content }}</p>

      <!-- Image Grid -->
      <div v-if="moment.images.length > 0" class="grid gap-2" :class="imageGridClass">
        <div
          v-for="(img, idx) in moment.images"
          :key="idx"
          class="group/img relative cursor-pointer overflow-hidden rounded-xl border border-slate-200"
          @click="openPreview(idx)"
        >
          <img
            :src="img"
            :alt="`图片 ${idx + 1}`"
            class="h-full w-full object-cover transition duration-300 group-hover/img:scale-105"
            :class="moment.images.length === 1 ? 'max-h-[400px]' : 'aspect-square'"
            loading="lazy"
          />
        </div>
      </div>

      <!-- Action Bar -->
      <div class="flex items-center gap-1 border-t border-slate-100 pt-3">
        <!-- Comment -->
        <button
          type="button"
          class="inline-flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-xs transition"
          :class="showComments ? 'bg-miku/10 text-miku' : 'text-slate-400 hover:bg-slate-100 hover:text-slate-600'"
          @click="showComments = !showComments"
        >
          <svg viewBox="0 0 24 24" class="h-4 w-4 fill-none stroke-current stroke-[1.8]" aria-hidden="true">
            <path d="M21 15a2 2 0 01-2 2H7l-4 4V5a2 2 0 012-2h14a2 2 0 012 2z" />
          </svg>
          {{ moment.comments.length || '' }}
        </button>

        <!-- Repost -->
        <button
          type="button"
          class="inline-flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-xs transition"
          :class="moment.reposted ? 'bg-emerald-50 text-emerald-500' : 'text-slate-400 hover:bg-slate-100 hover:text-emerald-500'"
          @click="handleRepost"
        >
          <svg viewBox="0 0 24 24" class="h-4 w-4 fill-none stroke-current stroke-[1.8]" aria-hidden="true">
            <path d="M17 1l4 4-4 4" /><path d="M3 11V9a4 4 0 014-4h14" /><path d="M7 23l-4-4 4-4" /><path d="M21 13v2a4 4 0 01-4 4H3" />
          </svg>
          {{ moment.reposts || '' }}
        </button>

        <!-- Like -->
        <button
          type="button"
          class="inline-flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-xs transition"
          :class="moment.liked ? 'bg-red-50 text-red-500' : 'text-slate-400 hover:bg-slate-100 hover:text-red-500'"
          @click="handleLike"
        >
          <svg viewBox="0 0 24 24" class="h-4 w-4 stroke-current stroke-[1.8]" :class="moment.liked ? 'fill-current' : 'fill-none'" aria-hidden="true">
            <path d="M20.84 4.61a5.5 5.5 0 00-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 00-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 000-7.78z" />
          </svg>
          {{ moment.likes || '' }}
        </button>

        <!-- Share -->
        <button
          type="button"
          class="ml-auto inline-flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-xs text-slate-400 transition hover:bg-slate-100 hover:text-slate-600"
          @click="handleShare"
        >
          <svg viewBox="0 0 24 24" class="h-4 w-4 fill-none stroke-current stroke-[1.8]" aria-hidden="true">
            <path d="M4 12v8a2 2 0 002 2h12a2 2 0 002-2v-8" /><polyline points="16,6 12,2 8,6" /><line x1="12" y1="2" x2="12" y2="15" />
          </svg>
        </button>
      </div>

      <!-- Comments Section -->
      <div v-if="showComments" class="space-y-3 border-t border-slate-100 pt-3">
        <!-- Comment List -->
        <div v-if="moment.comments.length > 0" class="space-y-2.5">
          <div v-for="comment in moment.comments" :key="comment.id" class="flex gap-2.5">
            <img
              :src="comment.avatar"
              :alt="comment.nickname"
              class="mt-0.5 h-7 w-7 shrink-0 rounded-full object-cover"
              loading="lazy"
            />
            <div class="min-w-0 flex-1">
              <div class="rounded-lg bg-slate-50 px-3 py-2">
                <div class="flex items-center gap-2">
                  <span class="text-xs font-semibold text-slate-600">{{ comment.nickname }}</span>
                  <span class="text-[10px] text-slate-400">{{ comment.createdAt }}</span>
                </div>
                <p class="mt-0.5 text-sm leading-relaxed text-slate-500">{{ comment.content }}</p>
              </div>
              <button
                type="button"
                class="mt-1 inline-flex items-center gap-1 px-1 text-[11px] transition"
                :class="comment.liked ? 'text-red-400' : 'text-slate-400 hover:text-red-400'"
                @click="handleLikeComment(comment.id)"
              >
                <svg viewBox="0 0 24 24" class="h-3 w-3 stroke-current stroke-[2]" :class="comment.liked ? 'fill-current' : 'fill-none'" aria-hidden="true">
                  <path d="M20.84 4.61a5.5 5.5 0 00-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 00-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 000-7.78z" />
                </svg>
                {{ comment.likes || '' }}
              </button>
            </div>
          </div>
        </div>

        <!-- Comment Input -->
        <div class="flex gap-2.5">
          <div class="mt-0.5 h-7 w-7 shrink-0 rounded-full bg-slate-200" />
          <div class="min-w-0 flex-1 space-y-2">
            <input
              v-model="commentNickname"
              type="text"
              placeholder="你的昵称"
              class="w-full rounded-lg border border-slate-200 bg-white px-3 py-1.5 text-sm text-slate-700 outline-none transition placeholder:text-slate-400 focus:border-miku/50 focus:ring-1 focus:ring-miku/30"
            />
            <div class="flex gap-2">
              <input
                v-model="commentContent"
                type="text"
                placeholder="写下评论..."
                class="min-w-0 flex-1 rounded-lg border border-slate-200 bg-white px-3 py-1.5 text-sm text-slate-700 outline-none transition placeholder:text-slate-400 focus:border-miku/50 focus:ring-1 focus:ring-miku/30"
                @keydown.enter="submitComment"
              />
              <button
                type="button"
                class="shrink-0 rounded-lg bg-miku px-3 py-1.5 text-xs font-medium text-white transition hover:bg-miku/85 disabled:opacity-50"
                :disabled="!commentNickname.trim() || !commentContent.trim()"
                @click="submitComment"
              >
                发送
              </button>
            </div>
          </div>
        </div>
      </div>
    </article>

    <!-- Image Preview Overlay -->
    <Teleport to="body">
      <div
        v-if="previewIndex !== null"
        class="fixed inset-0 z-[100] flex items-center justify-center bg-black/70 backdrop-blur-sm"
        @click.self="previewIndex = null"
      >
        <button
          type="button"
          class="absolute right-4 top-4 rounded-full bg-white/20 p-2 text-white transition hover:bg-white/30"
          @click="previewIndex = null"
        >
          <svg viewBox="0 0 24 24" class="h-5 w-5 fill-none stroke-current stroke-[2]" aria-hidden="true">
            <line x1="18" y1="6" x2="6" y2="18" /><line x1="6" y1="6" x2="18" y2="18" />
          </svg>
        </button>
        <img
          :src="moment.images[previewIndex]"
          alt="预览大图"
          class="max-h-[85vh] max-w-[90vw] rounded-lg object-contain"
        />
        <!-- Prev/Next -->
        <button
          v-if="moment.images.length > 1 && previewIndex > 0"
          type="button"
          class="absolute left-4 rounded-full bg-white/20 p-2 text-white transition hover:bg-white/30"
          @click.stop="previewIndex--"
        >
          <svg viewBox="0 0 24 24" class="h-5 w-5 fill-none stroke-current stroke-[2]" aria-hidden="true"><polyline points="15,18 9,12 15,6" /></svg>
        </button>
        <button
          v-if="moment.images.length > 1 && previewIndex < moment.images.length - 1"
          type="button"
          class="absolute right-4 rounded-full bg-white/20 p-2 text-white transition hover:bg-white/30"
          @click.stop="previewIndex++"
        >
          <svg viewBox="0 0 24 24" class="h-5 w-5 fill-none stroke-current stroke-[2]" aria-hidden="true"><polyline points="9,18 15,12 9,6" /></svg>
        </button>
      </div>
    </Teleport>
  </LiquidGlassCard>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import type { Moment } from '../../stores/moments'
import { addComment, toggleLikeComment, toggleLikeMoment, toggleRepostMoment } from '../../stores/moments'
import { showToast } from '../../stores/ui'
import LiquidGlassCard from '../ui/LiquidGlassCard.vue'

interface Props {
  moment: Moment
}

const props = defineProps<Props>()

const showComments = ref(false)
const commentNickname = ref('')
const commentContent = ref('')
const previewIndex = ref<number | null>(null)

const imageGridClass = computed(() => {
  const count = props.moment.images.length
  if (count === 1) return 'grid-cols-1'
  if (count === 2) return 'grid-cols-2'
  return 'grid-cols-2 sm:grid-cols-3'
})

function handleLike() {
  toggleLikeMoment(props.moment.id)
}

function handleRepost() {
  toggleRepostMoment(props.moment.id)
  if (!props.moment.reposted) {
    navigator.clipboard?.writeText(`${window.location.origin}/moments#${props.moment.id}`)
    showToast('链接已复制到剪贴板', 'success')
  }
}

function handleShare() {
  const url = `${window.location.origin}/moments#${props.moment.id}`
  if (navigator.share) {
    navigator.share({ title: `${props.moment.nickname} 的说说`, url })
  } else {
    navigator.clipboard?.writeText(url)
    showToast('链接已复制到剪贴板', 'success')
  }
}

function handleLikeComment(commentId: string) {
  toggleLikeComment(props.moment.id, commentId)
}

function openPreview(idx: number) {
  previewIndex.value = idx
}

async function submitComment() {
  if (!commentNickname.value.trim() || !commentContent.value.trim()) return
  try {
    await addComment({ momentId: props.moment.id, nickname: commentNickname.value, content: commentContent.value })
    commentContent.value = ''
    showToast('评论已发送', 'success')
  } catch {
    showToast('评论发送失败，请稍后重试', 'error')
  }
}
</script>
