<template>
  <section class="space-y-5">
    <!-- Section header -->
    <LiquidGlassCard padding="20px 24px" max-width="100%">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div class="flex items-center gap-3">
          <div class="flex h-9 w-9 items-center justify-center rounded-xl bg-gradient-to-br from-[#39c5bb]/20 to-[#c084fc]/20">
            <svg viewBox="0 0 24 24" class="h-[18px] w-[18px] fill-none stroke-[#39c5bb] stroke-[1.8]">
              <path d="M21 15a2 2 0 01-2 2H7l-4 4V5a2 2 0 012-2h14a2 2 0 012 2z" />
            </svg>
          </div>
          <div>
            <h2 class="text-lg font-bold text-slate-900">评论区</h2>
            <p class="mt-0.5 text-xs text-slate-500">评论提交后需审核通过方可展示</p>
          </div>
        </div>
        <span class="inline-flex items-center gap-1.5 rounded-full border border-[#39c5bb]/25 bg-[#39c5bb]/8 px-3 py-1 text-xs font-semibold text-[#39c5bb]">
          <svg viewBox="0 0 24 24" class="h-3 w-3 fill-none stroke-current stroke-[2]">
            <path d="M21 11.5a8.38 8.38 0 01-.9 3.8 8.5 8.5 0 01-7.6 4.7 8.38 8.38 0 01-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 01-.9-3.8 8.5 8.5 0 014.7-7.6 8.38 8.38 0 013.8-.9h.5a8.48 8.48 0 018 8v.5z" />
          </svg>
          {{ total }} 条
        </span>
      </div>
    </LiquidGlassCard>

    <!-- Comment form -->
    <LiquidGlassCard padding="22px 24px" max-width="100%">
      <div class="mb-4 flex items-center gap-2.5">
        <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-miku/10">
          <svg viewBox="0 0 24 24" class="h-4 w-4 fill-none stroke-miku stroke-[2]" aria-hidden="true">
            <path d="M12 20h9M16.5 3.5a2.121 2.121 0 013 3L7 19l-4 1 1-4 12.5-12.5z" />
          </svg>
        </div>
        <h3 class="text-base font-semibold text-slate-800">发布评论</h3>
      </div>
      <form class="space-y-4" @submit.prevent="submitComment">
        <div
          v-if="replyTarget"
          class="flex items-center justify-between gap-3 rounded-xl border border-[#c084fc]/30 bg-[#c084fc]/10 px-3 py-2"
        >
          <p class="text-xs text-slate-600">
            正在回复 <span class="font-semibold text-[#8b5cf6]">@{{ replyTarget.author }}</span>
          </p>
          <button
            type="button"
            class="rounded-lg border border-white/70 bg-white/75 px-2 py-1 text-xs text-slate-500 transition hover:text-slate-700"
            @click="cancelReply"
          >
            取消回复
          </button>
        </div>

        <div class="grid gap-4 md:grid-cols-2">
          <label class="group block">
            <span class="mb-1.5 block text-xs font-semibold text-slate-600">昵称 <span class="text-[#39c5bb]">*</span></span>
            <div class="relative">
              <svg viewBox="0 0 24 24" class="pointer-events-none absolute left-3 top-1/2 h-3.5 w-3.5 -translate-y-1/2 fill-none stroke-slate-400 stroke-[1.8] transition group-focus-within:stroke-[#39c5bb]">
                <path d="M20 21v-2a4 4 0 00-4-4H8a4 4 0 00-4 4v2" />
                <circle cx="12" cy="7" r="4" />
              </svg>
              <input
                v-model="author"
                type="text"
                maxlength="32"
                required
                class="w-full rounded-xl border border-slate-200/80 bg-white/70 py-2.5 pl-9 pr-3 text-sm text-slate-700 outline-none transition-all duration-300 placeholder:text-slate-400 focus:border-[#39c5bb]/50 focus:bg-white/90 focus:shadow-[0_0_0_3px_rgba(57,197,187,0.1)]"
                placeholder="请输入昵称"
                aria-label="评论昵称"
              />
            </div>
          </label>

          <label class="group block">
            <span class="mb-1.5 block text-xs font-semibold text-slate-600">邮箱（可选）</span>
            <div class="relative">
              <svg viewBox="0 0 24 24" class="pointer-events-none absolute left-3 top-1/2 h-3.5 w-3.5 -translate-y-1/2 fill-none stroke-slate-400 stroke-[1.8] transition group-focus-within:stroke-[#39c5bb]">
                <rect x="2" y="4" width="20" height="16" rx="2" />
                <path d="M22 7l-8.97 5.7a1.94 1.94 0 01-2.06 0L2 7" />
              </svg>
              <input
                v-model="email"
                type="email"
                maxlength="96"
                class="w-full rounded-xl border border-slate-200/80 bg-white/70 py-2.5 pl-9 pr-3 text-sm text-slate-700 outline-none transition-all duration-300 placeholder:text-slate-400 focus:border-[#39c5bb]/50 focus:bg-white/90 focus:shadow-[0_0_0_3px_rgba(57,197,187,0.1)]"
                placeholder="name@example.com"
                aria-label="评论邮箱"
              />
            </div>
          </label>
        </div>

        <label class="group block">
          <span class="mb-1.5 block text-xs font-semibold text-slate-600">评论内容 <span class="text-[#39c5bb]">*</span></span>
          <textarea
            v-model="content"
            ref="contentInputRef"
            rows="4"
            maxlength="800"
            required
            class="w-full resize-none rounded-xl border border-slate-200/80 bg-white/70 px-3 py-2.5 text-sm leading-relaxed text-slate-700 outline-none transition-all duration-300 placeholder:text-slate-400 focus:border-[#39c5bb]/50 focus:bg-white/90 focus:shadow-[0_0_0_3px_rgba(57,197,187,0.1)]"
            placeholder="写下你的看法..."
            aria-label="评论内容"
          />
          <span class="mt-1 block text-right text-[11px] tabular-nums" :class="content.length >= 750 ? 'text-red-400' : 'text-slate-400'">
            {{ content.length }} / 800
          </span>
        </label>

        <div class="flex flex-wrap items-center justify-between gap-3 pt-1">
          <Transition name="fade-slide" mode="out-in">
            <p
              v-if="feedback"
              :key="feedback"
              class="flex items-center gap-1.5 text-xs font-medium"
              :class="feedbackType === 'error' ? 'text-red-500' : feedbackType === 'success' ? 'text-[#39c5bb]' : 'text-slate-500'"
            >
              <svg v-if="feedbackType === 'success'" viewBox="0 0 24 24" class="h-3.5 w-3.5 fill-none stroke-current stroke-[2]">
                <path d="M22 11.08V12a10 10 0 11-5.93-9.14" />
                <path d="M22 4L12 14.01l-3-3" />
              </svg>
              <svg v-else-if="feedbackType === 'error'" viewBox="0 0 24 24" class="h-3.5 w-3.5 fill-none stroke-current stroke-[2]">
                <circle cx="12" cy="12" r="10" />
                <line x1="12" y1="8" x2="12" y2="12" />
                <line x1="12" y1="16" x2="12.01" y2="16" />
              </svg>
              {{ feedback }}
            </p>
          </Transition>
          <button
            type="submit"
            :disabled="submitting || !postId"
            class="ml-auto inline-flex items-center gap-2 rounded-xl border border-[#39c5bb]/40 bg-gradient-to-r from-[#39c5bb]/12 to-[#c084fc]/8 px-5 py-2.5 text-sm font-semibold text-[#39c5bb] shadow-sm transition-all duration-300 hover:border-[#39c5bb]/60 hover:shadow-[0_4px_16px_rgba(57,197,187,0.18)] active:scale-[0.97] disabled:cursor-not-allowed disabled:opacity-60"
          >
            <svg v-if="submitting" class="h-4 w-4 animate-spin" viewBox="0 0 24 24" fill="none">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="3" />
              <path class="opacity-80" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z" />
            </svg>
            <svg v-else viewBox="0 0 24 24" class="h-4 w-4 fill-none stroke-current stroke-[2]">
              <line x1="22" y1="2" x2="11" y2="13" />
              <polygon points="22 2 15 22 11 13 2 9 22 2" />
            </svg>
            {{ submitting ? '发送中...' : '发送评论' }}
          </button>
        </div>
      </form>
    </LiquidGlassCard>

    <!-- Loading skeleton -->
    <div v-if="loading" class="space-y-3">
      <LiquidGlassCard v-for="n in 3" :key="n" padding="16px 20px" max-width="100%" class="animate-pulse">
        <div class="flex items-center gap-3">
          <div class="h-9 w-9 rounded-full bg-slate-200/70" />
          <div class="flex-1 space-y-1.5">
            <div class="h-3 w-24 rounded bg-slate-200/60" />
            <div class="h-2.5 w-16 rounded bg-slate-100/80" />
          </div>
        </div>
        <div class="mt-3 space-y-2 pl-12">
          <div class="h-3 w-full rounded bg-slate-100/60" />
          <div class="h-3 w-4/5 rounded bg-slate-100/60" />
        </div>
      </LiquidGlassCard>
    </div>

    <!-- Empty state -->
    <LiquidGlassCard
      v-else-if="comments.length === 0"
      padding="40px 24px"
      max-width="100%"
      class="flex flex-col items-center text-center"
    >
      <div class="mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-gradient-to-br from-[#39c5bb]/15 to-[#c084fc]/15">
        <svg viewBox="0 0 24 24" class="h-7 w-7 fill-none stroke-[#39c5bb]/60 stroke-[1.5]">
          <path d="M21 15a2 2 0 01-2 2H7l-4 4V5a2 2 0 012-2h14a2 2 0 012 2z" />
          <line x1="9" y1="10" x2="15" y2="10" />
        </svg>
      </div>
      <p class="text-sm font-semibold text-slate-700">还没有评论</p>
      <p class="mt-1 max-w-xs text-xs leading-relaxed text-slate-500">来发表第一条评论吧，分享你的想法与见解。</p>
    </LiquidGlassCard>

    <!-- Comment list -->
    <TransitionGroup v-else name="comment-item" tag="div" class="space-y-4">
      <LiquidGlassCard
        v-for="item in topLevelComments"
        :key="item.id"
        padding="16px 20px"
        max-width="100%"
        class="transition duration-300"
      >
        <article class="flex gap-3">
          <div class="flex shrink-0 flex-col items-center gap-0.5 pt-0.5 text-slate-400">
            <button type="button" class="rounded p-0.5 transition hover:bg-miku/10">
              <svg viewBox="0 0 24 24" class="h-4 w-4 fill-current"><path d="M12 4l-8 8h5v8h6v-8h5z" /></svg>
            </button>
            <span class="min-w-[20px] text-center text-xs font-semibold text-slate-500">0</span>
            <button type="button" class="rounded p-0.5 transition hover:bg-red-50">
              <svg viewBox="0 0 24 24" class="h-4 w-4 fill-current"><path d="M12 20l8-8h-5V4H9v8H4z" /></svg>
            </button>
          </div>

          <div class="min-w-0 flex-1">
            <div class="flex flex-wrap items-center gap-x-2 gap-y-1">
              <img :src="avatarUrl(item.author_name)" :alt="item.author_name || '匿名用户'" class="h-5 w-5 rounded-full object-cover" loading="lazy" />
              <span class="text-sm font-semibold text-slate-700">{{ item.author_name || '匿名用户' }}</span>
              <span class="h-0.5 w-0.5 rounded-full bg-slate-300" aria-hidden="true" />
              <span class="text-[11px] text-slate-400">{{ formatDate(item.created_at) }}</span>
            </div>

            <p class="mt-2 whitespace-pre-wrap break-words text-sm leading-relaxed text-slate-600">{{ item.content }}</p>

            <div class="mt-2.5 flex items-center gap-3">
              <button
                type="button"
                class="inline-flex items-center gap-1 rounded px-1.5 py-1 text-[11px] text-slate-400 transition hover:bg-slate-100 hover:text-slate-600"
                @click="startReply(item)"
              >
                <svg viewBox="0 0 24 24" class="h-3.5 w-3.5 fill-none stroke-current stroke-[1.8]" aria-hidden="true">
                  <path d="M21 15a2 2 0 01-2 2H7l-4 4V5a2 2 0 012-2h14a2 2 0 012 2z" />
                </svg>
                回复
              </button>
              <span v-if="childReplies(item.id).length > 0" class="text-[11px] text-slate-400">
                {{ childReplies(item.id).length }} 条回复
              </span>
            </div>

            <div v-if="childReplies(item.id).length > 0" class="mt-3 space-y-2 border-l-2 border-slate-200 pl-4">
              <div v-for="reply in childReplies(item.id)" :key="reply.id" class="rounded-lg bg-slate-50/60 px-3 py-2.5">
                <div class="flex flex-wrap items-center gap-x-2 gap-y-1">
                  <img :src="avatarUrl(reply.author_name)" :alt="reply.author_name || '匿名用户'" class="h-4 w-4 rounded-full object-cover" loading="lazy" />
                  <span class="text-xs font-semibold text-slate-600">{{ reply.author_name || '匿名用户' }}</span>
                  <span class="h-0.5 w-0.5 rounded-full bg-slate-300" aria-hidden="true" />
                  <span class="text-[11px] text-slate-400">{{ formatDate(reply.created_at) }}</span>
                </div>
                <p class="mt-1 text-sm leading-relaxed text-slate-500">{{ reply.content }}</p>
                <button
                  type="button"
                  class="mt-1.5 inline-flex items-center gap-1 rounded px-1.5 py-1 text-[11px] text-slate-400 transition hover:bg-slate-100 hover:text-slate-600"
                  @click="startReply(reply)"
                >
                  <svg viewBox="0 0 24 24" class="h-3.5 w-3.5 fill-none stroke-current stroke-[1.8]" aria-hidden="true">
                    <path d="M21 15a2 2 0 01-2 2H7l-4 4V5a2 2 0 012-2h14a2 2 0 012 2z" />
                  </svg>
                  回复
                </button>
              </div>
            </div>
          </div>
        </article>
      </LiquidGlassCard>
    </TransitionGroup>
  </section>
</template>

<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue'

import { api, type PagedData } from '../../lib/api'
import LiquidGlassCard from '../ui/LiquidGlassCard.vue'

interface PostCommentItem {
  id: string
  parent_id?: string | null
  author_name: string
  content: string
  created_at: string
}

interface Props {
  postId: string
  initialTotal?: number
}

const props = withDefaults(defineProps<Props>(), {
  initialTotal: 0,
})

const emit = defineEmits<{
  (e: 'count-updated', value: number): void
}>()

const comments = ref<PostCommentItem[]>([])
const total = ref(0)
const loading = ref(false)
const submitting = ref(false)

const author = ref('')
const email = ref('')
const content = ref('')
const contentInputRef = ref<HTMLTextAreaElement | null>(null)
const replyTarget = ref<{ id: string; author: string } | null>(null)

const feedback = ref('')
const feedbackType = ref<'neutral' | 'success' | 'error'>('neutral')
const repliesByParent = computed(() => {
  const map = new Map<string, PostCommentItem[]>()
  comments.value.forEach((item) => {
    if (!item.parent_id) return
    const list = map.get(item.parent_id) || []
    list.push(item)
    map.set(item.parent_id, list)
  })
  return map
})
const topLevelComments = computed(() => {
  const tops = comments.value.filter((item) => !item.parent_id)
  return tops.length > 0 ? tops : comments.value
})

function avatarUrl(name: string): string {
  return `https://api.dicebear.com/9.x/shapes/svg?seed=${encodeURIComponent(name || '匿名用户')}`
}

function childReplies(parentId: string): PostCommentItem[] {
  return repliesByParent.value.get(parentId) || []
}

async function startReply(item: PostCommentItem) {
  replyTarget.value = {
    id: item.id,
    author: item.author_name || '匿名用户',
  }
  await nextTick()
  contentInputRef.value?.focus()
}

function cancelReply() {
  replyTarget.value = null
}

watch(
  () => props.initialTotal,
  (next) => {
    total.value = Number(next || 0)
  },
  { immediate: true },
)

watch(
  () => props.postId,
  async () => {
    replyTarget.value = null
    await loadComments()
  },
  { immediate: true },
)

function setFeedback(message: string, type: 'neutral' | 'success' | 'error' = 'neutral') {
  feedback.value = message
  feedbackType.value = type
}

function formatDate(iso: string): string {
  try {
    const d = new Date(iso)
    if (Number.isNaN(d.getTime())) return '--'
    const yyyy = d.getFullYear()
    const mm = String(d.getMonth() + 1).padStart(2, '0')
    const dd = String(d.getDate()).padStart(2, '0')
    const hh = String(d.getHours()).padStart(2, '0')
    const mi = String(d.getMinutes()).padStart(2, '0')
    return `${yyyy}-${mm}-${dd} ${hh}:${mi}`
  } catch {
    return '--'
  }
}

function isValidEmail(input: string): boolean {
  if (!input.trim()) return true
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(input.trim())
}

async function loadComments() {
  if (!props.postId) {
    comments.value = []
    total.value = 0
    emit('count-updated', 0)
    return
  }

  loading.value = true
  try {
    const data = await api.get<PagedData<PostCommentItem>>(`/posts/${encodeURIComponent(props.postId)}/comments?page=1&size=20`)
    comments.value = data.items || []
    total.value = Number(data.total || 0)
    emit('count-updated', total.value)
  } catch {
    comments.value = []
    setFeedback('评论加载失败，请稍后重试。', 'error')
  } finally {
    loading.value = false
  }
}

async function submitComment() {
  if (!props.postId || submitting.value) return

  const finalAuthor = author.value.trim()
  const finalEmail = email.value.trim()
  const finalContent = content.value.trim()

  if (!finalAuthor || !finalContent) {
    setFeedback('昵称和评论内容不能为空。', 'error')
    return
  }
  if (!isValidEmail(finalEmail)) {
    setFeedback('邮箱格式不正确，请检查后重试。', 'error')
    return
  }

  submitting.value = true
  setFeedback('')
  try {
    const payload: Record<string, string> = {
      author_name: finalAuthor,
      author_email: finalEmail,
      content: finalContent,
    }
    if (replyTarget.value?.id) {
      payload.parent_id = replyTarget.value.id
    }

    await api.post(`/posts/${encodeURIComponent(props.postId)}/comments`, payload)
    content.value = ''
    replyTarget.value = null
    setFeedback('评论已提交，审核通过后会显示在列表中。', 'success')
    await loadComments()
  } catch {
    setFeedback('评论提交失败，请稍后重试。', 'error')
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.25s ease;
}
.fade-slide-enter-from {
  opacity: 0;
  transform: translateY(-4px);
}
.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(4px);
}

.comment-item-enter-active {
  transition: all 0.4s cubic-bezier(0.16, 1, 0.3, 1);
}
.comment-item-leave-active {
  transition: all 0.25s ease;
}
.comment-item-enter-from {
  opacity: 0;
  transform: translateY(12px);
}
.comment-item-leave-to {
  opacity: 0;
  transform: translateX(-12px);
}
.comment-item-move {
  transition: transform 0.35s ease;
}
</style>
