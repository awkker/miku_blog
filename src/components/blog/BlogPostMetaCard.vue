<template>
  <LiquidGlassFrame
    class="meta-card"
    width="100%"
    max-width="none"
    padding="16px"
    :border-radius="18"
    :displacement-strength="1.08"
    :edge-refraction-strength="1.24"
    :interactive="false"
  >
    <h3>阅读工具</h3>
    <button type="button" class="meta-btn back-btn" @click="emit('back')">返回博客列表</button>

    <div class="meta-list">
      <p><span>发布日期</span><strong>{{ date }}</strong></p>
      <p><span>分类标签</span><strong>{{ category }}</strong></p>
      <p><span>预计阅读</span><strong>{{ readingMinutes }} 分钟</strong></p>
      <p><span>阅读次数</span><strong>{{ views }}</strong></p>
      <p><span>点赞总数</span><strong>{{ formatCount(likeCount) }}</strong></p>
      <p><span>分享次数</span><strong>{{ formatCount(shareCount) }}</strong></p>
    </div>

    <div class="meta-actions">
      <button type="button" class="meta-btn" :class="{ 'is-active': isLiked }" @click="emit('toggle-like')">
        {{ isLiked ? '已点赞' : '点赞' }} · {{ formatCount(likeCount) }}
      </button>
      <button type="button" class="meta-btn" @click="emit('share')">分享 · {{ formatCount(shareCount) }}</button>
    </div>

    <p v-if="shareHint" class="share-hint">{{ shareHint }}</p>
  </LiquidGlassFrame>
</template>

<script setup lang="ts">
import LiquidGlassFrame from '@/components/LiquidGlassFrame.vue'

interface Props {
  date: string
  category: string
  readingMinutes: number
  views: string
  isLiked: boolean
  likeCount: number
  shareCount: number
  shareHint: string
}

defineProps<Props>()

const emit = defineEmits<{
  back: []
  'toggle-like': []
  share: []
}>()

const NUMBER_FORMATTER = new Intl.NumberFormat('zh-CN')

function formatCount(value: number): string {
  return NUMBER_FORMATTER.format(Math.max(0, Math.floor(value)))
}
</script>

<style scoped>
.meta-card h3 {
  margin: 0 0 10px;
  font-size: 0.95rem;
  color: rgba(10, 64, 54, 0.86);
}

.meta-list {
  margin: 10px 0 12px;
  display: grid;
  gap: 8px;
}

.meta-list p {
  margin: 0;
  display: flex;
  justify-content: space-between;
  gap: 10px;
  font-size: 0.84rem;
  color: rgba(19, 78, 65, 0.74);
}

.meta-list p span {
  flex: 0 0 auto;
}

.meta-list p strong {
  margin: 0;
  font-weight: 500;
  color: rgba(8, 60, 50, 0.9);
  text-align: right;
}

.meta-actions {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
}

.meta-btn {
  border: 0;
  border-radius: 10px;
  padding: 6px 10px;
  background: rgba(236, 251, 246, 0.68);
  color: rgba(10, 73, 60, 0.86);
  font-size: 0.82rem;
  font-family: var(--font-body);
  cursor: pointer;
  transition: background-color 180ms ease;
}

.meta-btn:hover {
  background: rgba(223, 248, 239, 0.92);
}

.meta-btn.is-active {
  color: rgba(7, 63, 52, 0.96);
  background: rgba(205, 243, 232, 0.96);
}

.back-btn {
  width: 100%;
}

.share-hint {
  margin: 8px 2px 0;
  min-height: 18px;
  color: rgba(22, 86, 72, 0.74);
  font-size: 0.76rem;
}

@media (max-width: 1180px) {
  .meta-card {
    order: -1;
  }

  .meta-actions {
    grid-template-columns: repeat(2, minmax(0, max-content));
    justify-content: start;
  }
}
</style>
