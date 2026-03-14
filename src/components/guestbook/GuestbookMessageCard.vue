<template>
  <LiquidGlassCard padding="16px 20px" class="transition duration-300">
    <article class="flex gap-3">
      <!-- Vote Column -->
      <div class="flex shrink-0 flex-col items-center gap-0.5 pt-0.5">
        <button
          type="button"
          class="rounded p-0.5 transition hover:bg-miku/10"
          :class="message.myVote === 1 ? 'text-miku' : 'text-slate-400'"
          aria-label="赞同"
          @click="vote(1)"
        >
          <svg viewBox="0 0 24 24" class="h-4 w-4 fill-current"><path d="M12 4l-8 8h5v8h6v-8h5z" /></svg>
        </button>
        <span class="min-w-[20px] text-center text-xs font-semibold" :class="voteColor">{{ message.votes }}</span>
        <button
          type="button"
          class="rounded p-0.5 transition hover:bg-red-50"
          :class="message.myVote === -1 ? 'text-red-500' : 'text-slate-400'"
          aria-label="反对"
          @click="vote(-1)"
        >
          <svg viewBox="0 0 24 24" class="h-4 w-4 fill-current"><path d="M12 20l8-8h-5V4H9v8H4z" /></svg>
        </button>
      </div>

      <!-- Content -->
      <div class="min-w-0 flex-1">
        <div class="flex flex-wrap items-center gap-x-2 gap-y-1">
          <img :src="message.avatar" :alt="message.nickname" class="h-5 w-5 rounded-full object-cover" loading="lazy" />
          <span class="text-sm font-semibold text-slate-700">{{ message.nickname }}</span>
          <a
            v-if="message.website"
            :href="message.website"
            target="_blank"
            rel="noopener noreferrer"
            class="inline-flex items-center gap-0.5 text-[11px] text-miku/70 transition hover:text-miku"
          >
            <svg viewBox="0 0 24 24" class="h-2.5 w-2.5 fill-none stroke-current stroke-[2]" aria-hidden="true">
              <circle cx="12" cy="12" r="10" /><path d="M2 12h20M12 2a15.3 15.3 0 014 10 15.3 15.3 0 01-4 10 15.3 15.3 0 01-4-10 15.3 15.3 0 014-10z" />
            </svg>
            {{ message.website.replace(/^https?:\/\//, '').replace(/\/$/, '') }}
          </a>
          <span class="h-0.5 w-0.5 rounded-full bg-slate-300" aria-hidden="true" />
          <span class="text-[11px] text-slate-400">{{ message.createdAt }}</span>
        </div>

        <p class="mt-2 whitespace-pre-wrap break-words text-sm leading-relaxed text-slate-600">{{ message.message }}</p>

        <!-- Action Bar -->
        <div class="mt-2.5 flex items-center gap-3">
          <button
            type="button"
            class="inline-flex items-center gap-1 rounded px-1.5 py-1 text-[11px] text-slate-400 transition hover:bg-slate-100 hover:text-slate-600"
            @click="showReplyForm = !showReplyForm"
          >
            <svg viewBox="0 0 24 24" class="h-3.5 w-3.5 fill-none stroke-current stroke-[1.8]" aria-hidden="true">
              <path d="M21 15a2 2 0 01-2 2H7l-4 4V5a2 2 0 012-2h14a2 2 0 012 2z" />
            </svg>
            回复
          </button>
          <span v-if="message.replies.length > 0" class="text-[11px] text-slate-400">
            {{ message.replies.length }} 条回复
          </span>
        </div>

        <!-- Inline Reply Form -->
        <div v-if="showReplyForm" class="mt-3 rounded-lg border border-slate-200 bg-slate-50/80 p-3">
          <div class="mb-2">
            <input
              v-model="replyNickname"
              type="text"
              placeholder="你的昵称"
              class="w-full rounded-md border border-slate-200 bg-white px-3 py-1.5 text-sm text-slate-700 outline-none transition placeholder:text-slate-400 focus:border-miku/50 focus:ring-1 focus:ring-miku/30"
            />
          </div>
          <textarea
            v-model="replyContent"
            rows="2"
            placeholder="写下你的回复..."
            class="w-full resize-none rounded-md border border-slate-200 bg-white px-3 py-1.5 text-sm text-slate-700 outline-none transition placeholder:text-slate-400 focus:border-miku/50 focus:ring-1 focus:ring-miku/30"
          />
          <div class="mt-2 flex items-center justify-end gap-2">
            <button
              type="button"
              class="rounded-lg px-3 py-1.5 text-xs text-slate-500 transition hover:bg-slate-200"
              @click="showReplyForm = false"
            >
              取消
            </button>
            <button
              type="button"
              class="rounded-lg bg-miku px-3 py-1.5 text-xs font-medium text-white transition hover:bg-miku/85 disabled:opacity-50"
              :disabled="!replyNickname.trim() || !replyContent.trim()"
              @click="submitReply"
            >
              回复
            </button>
          </div>
        </div>

        <!-- Nested Replies -->
        <div v-if="message.replies.length > 0" class="mt-3 space-y-2 border-l-2 border-slate-200 pl-4">
          <div v-for="reply in message.replies" :key="reply.id" class="rounded-lg bg-slate-50/60 px-3 py-2.5">
            <div class="flex flex-wrap items-center gap-x-2 gap-y-1">
              <img :src="reply.avatar" :alt="reply.nickname" class="h-4 w-4 rounded-full object-cover" loading="lazy" />
              <span class="text-xs font-semibold text-slate-600">{{ reply.nickname }}</span>
              <span class="h-0.5 w-0.5 rounded-full bg-slate-300" aria-hidden="true" />
              <span class="text-[11px] text-slate-400">{{ reply.createdAt }}</span>
            </div>
            <p class="mt-1 text-sm leading-relaxed text-slate-500">{{ reply.message }}</p>
            <div class="mt-1.5 flex items-center gap-2">
              <button
                type="button"
                class="inline-flex items-center gap-0.5 rounded p-0.5 text-[11px] transition"
                :class="reply.myVote === 1 ? 'text-miku' : 'text-slate-400 hover:text-miku'"
                @click="voteReply(reply.id, 1)"
              >
                <svg viewBox="0 0 24 24" class="h-3 w-3 fill-current"><path d="M12 4l-8 8h5v8h6v-8h5z" /></svg>
                {{ reply.votes }}
              </button>
              <button
                type="button"
                class="rounded p-0.5 text-[11px] transition"
                :class="reply.myVote === -1 ? 'text-red-500' : 'text-slate-400 hover:text-red-400'"
                @click="voteReply(reply.id, -1)"
              >
                <svg viewBox="0 0 24 24" class="h-3 w-3 fill-current"><path d="M12 20l8-8h-5V4H9v8H4z" /></svg>
              </button>
            </div>
          </div>
        </div>
      </div>
    </article>
  </LiquidGlassCard>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import type { GuestbookMessage } from '../../stores/guestbook'
import { voteMessage } from '../../stores/guestbook'
import LiquidGlassCard from '../ui/LiquidGlassCard.vue'

interface Props {
  message: GuestbookMessage
}

const props = defineProps<Props>()
const emit = defineEmits<{ reply: [payload: { parentId: string; nickname: string; message: string }] }>()

const showReplyForm = ref(false)
const replyNickname = ref('')
const replyContent = ref('')

const voteColor = computed(() => {
  if (props.message.myVote === 1) return 'text-miku'
  if (props.message.myVote === -1) return 'text-red-500'
  return 'text-slate-500'
})

function vote(direction: 1 | -1) {
  voteMessage(props.message.id, direction)
}

function voteReply(replyId: string, direction: 1 | -1) {
  voteMessage(replyId, direction, props.message.id)
}

function submitReply() {
  if (!replyNickname.value.trim() || !replyContent.value.trim()) return
  emit('reply', { parentId: props.message.id, nickname: replyNickname.value, message: replyContent.value })
  replyContent.value = ''
  showReplyForm.value = false
}
</script>
