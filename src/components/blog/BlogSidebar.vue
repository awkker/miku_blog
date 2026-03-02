<template>
  <div class="sidebar-wrap">
    <LiquidGlassFrame
      class="profile-card"
      width="100%"
      max-width="none"
      padding="18px 16px"
      :border-radius="20"
      :displacement-strength="1.08"
      :edge-refraction-strength="1.3"
      :interactive="false"
    >
      <div class="profile-head">
        <div class="avatar">薰</div>
        <div>
          <h3>薰逸</h3>
          <p>旅途记录者 / 前端开发者</p>
        </div>
      </div>
      <p class="profile-intro">喜欢把漫长的旅途切成小片段，存进文字和照片里。偶尔也写点代码，让页面像风景一样呼吸。</p>
    </LiquidGlassFrame>

    <LiquidGlassFrame
      class="hot-card"
      width="100%"
      max-width="none"
      padding="20px"
      :border-radius="22"
      :displacement-strength="1.08"
      :edge-refraction-strength="1.3"
      :interactive="false"
    >
      <h3 class="hot-title">热门文章</h3>
      <ul class="hot-list">
        <li v-for="item in hotPosts" :key="item.slug">
          <button type="button" class="hot-link" @click="emit('openPost', item.slug)">
            <p class="hot-name">{{ item.title }}</p>
            <p class="hot-meta">{{ item.views }} 次阅读</p>
          </button>
        </li>
      </ul>
    </LiquidGlassFrame>
  </div>
</template>

<script setup lang="ts">
import LiquidGlassFrame from '@/components/LiquidGlassFrame.vue'
import type { BlogPostMeta } from '@/content/blogPosts'

interface Props {
  hotPosts: BlogPostMeta[]
}

defineProps<Props>()

const emit = defineEmits<{
  openPost: [slug: string]
}>()
</script>

<style scoped>
.sidebar-wrap {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.profile-head {
  display: flex;
  align-items: center;
  gap: 12px;
}

.avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--font-title);
  font-size: 1.3rem;
  color: rgba(15, 72, 60, 0.95);
  background: rgba(234, 255, 248, 0.75);
}

.profile-head h3 {
  margin: 0;
  font-size: 1.16rem;
  color: rgba(8, 52, 43, 0.92);
}

.profile-head p {
  margin: 4px 0 0;
  font-size: 0.88rem;
  color: rgba(22, 84, 71, 0.82);
}

.profile-intro {
  margin: 14px 0 0;
  line-height: 1.7;
  color: rgba(20, 78, 66, 0.88);
  font-size: 0.92rem;
}

.hot-title {
  margin: 0 0 12px;
  font-size: 1.08rem;
  color: rgba(9, 54, 45, 0.92);
}

.hot-list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 11px;
}

.hot-list li {
  border-radius: 12px;
  padding: 0;
  background: rgba(244, 255, 252, 0.58);
  overflow: hidden;
}

.hot-link {
  width: 100%;
  border: 0;
  background: transparent;
  text-align: left;
  padding: 10px 12px;
  cursor: pointer;
  transition: background-color 220ms ease;
}

.hot-link:hover {
  background: rgba(225, 249, 241, 0.56);
}

.hot-name {
  margin: 0;
  font-size: 0.92rem;
  line-height: 1.5;
  color: rgba(8, 56, 46, 0.92);
}

.hot-meta {
  margin: 4px 0 0;
  font-size: 0.8rem;
  color: rgba(24, 92, 78, 0.72);
}
</style>

