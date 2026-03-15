<template>
  <section class="space-y-6">
    <!-- Post Form -->
    <LiquidGlassCard padding="20px 24px">
      <div class="mb-4 flex items-center gap-2.5">
        <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-miku/10">
          <svg viewBox="0 0 24 24" class="h-4 w-4 fill-none stroke-miku stroke-[2]" aria-hidden="true">
            <path d="M12 20h9M16.5 3.5a2.121 2.121 0 013 3L7 19l-4 1 1-4 12.5-12.5z" />
          </svg>
        </div>
        <h2 class="text-lg font-semibold text-slate-800">发布说说</h2>
      </div>
      <form class="space-y-3" @submit.prevent="handleSubmit">
        <MikuInput v-model="form.nickname" label="昵称" placeholder="请输入昵称" :error="errors.nickname" aria-label="昵称" required />
        <MikuTextarea v-model="form.content" label="内容" placeholder="此刻的想法..." :error="errors.content" :rows="3" aria-label="说说内容" required />

        <!-- Image URLs -->
        <div class="space-y-2">
          <label class="text-sm font-medium text-slate-600">图片（可选，输入图片链接）</label>
          <div v-for="(_, idx) in form.images" :key="idx" class="flex items-center gap-2">
            <input
              v-model="form.images[idx]"
              type="text"
              placeholder="https://example.com/photo.jpg"
              class="min-w-0 flex-1 rounded-lg border border-slate-200 bg-white/80 px-3 py-1.5 text-sm text-slate-700 outline-none transition placeholder:text-slate-400 focus:border-miku/50 focus:ring-1 focus:ring-miku/30"
            />
            <button
              type="button"
              class="shrink-0 rounded-lg p-1.5 text-slate-400 transition hover:bg-red-50 hover:text-red-500"
              @click="removeImage(idx)"
            >
              <svg viewBox="0 0 24 24" class="h-4 w-4 fill-none stroke-current stroke-[2]" aria-hidden="true">
                <line x1="18" y1="6" x2="6" y2="18" /><line x1="6" y1="6" x2="18" y2="18" />
              </svg>
            </button>
          </div>
          <button
            v-if="form.images.length < 4"
            type="button"
            class="inline-flex items-center gap-1.5 rounded-lg border border-dashed border-slate-300 px-3 py-1.5 text-xs text-slate-400 transition hover:border-miku/40 hover:text-miku"
            @click="addImageSlot"
          >
            <svg viewBox="0 0 24 24" class="h-3.5 w-3.5 fill-none stroke-current stroke-[1.8]" aria-hidden="true">
              <rect x="3" y="3" width="18" height="18" rx="2" ry="2" /><circle cx="8.5" cy="8.5" r="1.5" /><polyline points="21,15 16,10 5,21" />
            </svg>
            添加图片链接
          </button>
        </div>

        <div class="flex justify-end">
          <MikuButton type="submit" :loading="submitting" :disabled="submitting" aria-label="发布说说">
            {{ submitting ? '发布中...' : '发布' }}
          </MikuButton>
        </div>
      </form>
    </LiquidGlassCard>

    <!-- Feed Count -->
    <div v-if="list.length > 0" class="flex items-center justify-between px-1">
      <span class="text-xs text-slate-400">{{ list.length }} 条说说</span>
    </div>

    <!-- Feed -->
    <ErrorState v-if="fetchStatus === 'error'" :description="fetchError || '说说加载失败，请稍后再试。'" @retry="loadFeed" />

    <div v-else class="space-y-4">
      <div v-if="fetchStatus === 'loading'" class="space-y-4">
        <SkeletonCard v-for="item in 3" :key="item" />
      </div>

      <EmptyState v-else-if="list.length === 0" title="还没有说说" description="发布第一条说说，记录此刻的想法。" />

      <TransitionGroup v-else name="feed-rise" tag="div" class="space-y-4">
        <MomentCard v-for="item in list" :key="item.id" :moment="item" />
      </TransitionGroup>
    </div>
  </section>
</template>

<script setup lang="ts">
import { useStore } from '@nanostores/vue'
import { computed, onMounted, reactive } from 'vue'

import {
  loadMoments,
  moments,
  momentsError,
  momentsFetchStatus,
  momentsSubmitStatus,
  submitMoment,
} from '../../stores/moments'
import { showToast } from '../../stores/ui'
import EmptyState from '../ui/EmptyState.vue'
import ErrorState from '../ui/ErrorState.vue'
import LiquidGlassCard from '../ui/LiquidGlassCard.vue'
import MikuButton from '../ui/MikuButton.vue'
import MikuInput from '../ui/MikuInput.vue'
import MikuTextarea from '../ui/MikuTextarea.vue'
import MomentCard from './MomentCard.vue'
import SkeletonCard from '../ui/SkeletonCard.vue'

const form = reactive({ nickname: '', content: '', images: [] as string[] })
const errors = reactive({ nickname: '', content: '' })

const list = useStore(moments)
const fetchStatus = useStore(momentsFetchStatus)
const submitStatus = useStore(momentsSubmitStatus)
const fetchError = useStore(momentsError)

const submitting = computed(() => submitStatus.value === 'loading')

function addImageSlot() {
  if (form.images.length < 4) form.images.push('')
}

function removeImage(idx: number) {
  form.images.splice(idx, 1)
}

function validate() {
  errors.nickname = form.nickname.trim() ? '' : '昵称不能为空'
  errors.content = form.content.trim() ? '' : '内容不能为空'
  return !errors.nickname && !errors.content
}

async function loadFeed() {
  await loadMoments()
}

async function handleSubmit() {
  if (!validate()) return
  try {
    const validImages = form.images.filter((u) => u.trim())
    await submitMoment({ nickname: form.nickname, content: form.content, images: validImages })
    form.content = ''
    form.images = []
    showToast('说说已发布', 'success')
  } catch {
    showToast('发布失败，请稍后重试', 'error')
  }
}

onMounted(async () => {
  await loadFeed()
})
</script>

<style scoped>
.feed-rise-enter-active,
.feed-rise-leave-active {
  transition: transform 0.35s ease, opacity 0.35s ease;
}
.feed-rise-enter-from,
.feed-rise-leave-to {
  opacity: 0;
  transform: translateY(10px);
}
.feed-rise-move {
  transition: transform 0.35s ease;
}
</style>
