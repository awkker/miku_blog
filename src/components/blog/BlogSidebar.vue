<template>
  <div class="sidebar-wrap">
    <div class="interactive-shell" @pointermove="handleSurfacePointerMove" @pointerleave="handleSurfacePointerLeave">
      <LiquidGlassFrame
        class="profile-card tilt-card"
        width="100%"
        max-width="none"
        padding="18px 16px"
        :border-radius="20"
        :displacement-strength="1.08"
        :edge-refraction-strength="1.3"
        :interactive="false"
      >
        <div class="profile-top depth-2">
          <div class="avatar-orbit depth-3" aria-hidden="true"></div>
          <img class="avatar depth-3" src="/photo/author.jpg" alt="薰逸头像" />
          <h3>薰逸</h3>
          <p>旅途记录者 / 前端开发者</p>
        </div>
        <p class="profile-intro depth-1">喜欢把漫长的旅途切成小片段，存进文字和照片里。偶尔也写点代码，让页面像风景一样呼吸。</p>
        <div class="profile-stats depth-1">
          <div class="stat-item">
            <span>博客数</span>
            <strong>{{ formatMetric(authorStats.postCount) }}</strong>
          </div>
          <div class="stat-item">
            <span>总阅读</span>
            <strong>{{ formatMetric(authorStats.totalViews) }}</strong>
          </div>
          <div class="stat-item">
            <span>获赞数</span>
            <strong>{{ formatMetric(authorStats.totalLikes) }}</strong>
          </div>
          <div class="stat-item">
            <span>博客访问</span>
            <strong>{{ formatMetric(authorStats.totalBlogVisits) }}</strong>
          </div>
        </div>
      </LiquidGlassFrame>
    </div>

    <div class="interactive-shell hot-shell">
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
        <h3 class="hot-title depth-2">热门文章</h3>
        <TransitionGroup name="hot-stagger" tag="ul" class="hot-list">
          <li v-for="(item, index) in hotPosts" :key="item.slug" :style="getStaggerStyle(index)">
            <RouterLink class="hot-link" :to="{ name: 'blog-post', params: { slug: item.slug } }">
              <p class="hot-name">{{ item.title }}</p>
              <p class="hot-meta">{{ item.views }} 次阅读</p>
            </RouterLink>
          </li>
        </TransitionGroup>
      </LiquidGlassFrame>
    </div>
  </div>
</template>

<script setup lang="ts">
import LiquidGlassFrame from '@/components/LiquidGlassFrame.vue'
import type { BlogAuthorStats, BlogPostMeta } from '@/content/blogPosts'

interface Props {
  hotPosts: BlogPostMeta[]
  authorStats: BlogAuthorStats
}

defineProps<Props>()

function getStaggerStyle(index: number): Record<string, string> {
  return {
    '--stagger-delay': `${index * 72}ms`,
  }
}

function formatMetric(value: number): string {
  const absolute = Math.abs(value)
  if (absolute >= 1_000_000_000) {
    return `${trimZero((value / 1_000_000_000).toFixed(1))}B`
  }
  if (absolute >= 1_000_000) {
    return `${trimZero((value / 1_000_000).toFixed(1))}M`
  }
  if (absolute >= 1_000) {
    return `${trimZero((value / 1_000).toFixed(1))}k`
  }
  return String(value)
}

function trimZero(value: string): string {
  return value.replace(/\.0$/, '')
}

function prefersReducedMotion(): boolean {
  if (typeof window === 'undefined' || typeof window.matchMedia !== 'function') {
    return false
  }
  return window.matchMedia('(prefers-reduced-motion: reduce)').matches
}

function handleSurfacePointerMove(event: PointerEvent): void {
  if (prefersReducedMotion()) {
    return
  }
  const target = event.currentTarget
  if (!(target instanceof HTMLElement)) {
    return
  }
  const rect = target.getBoundingClientRect()
  if (rect.width <= 0 || rect.height <= 0) {
    return
  }
  const x = event.clientX - rect.left
  const y = event.clientY - rect.top
  const px = x / rect.width
  const py = y / rect.height
  const rx = (0.5 - py) * 7.5
  const ry = (px - 0.5) * 9
  const offsetX = (px - 0.5) * 9
  const offsetY = (py - 0.5) * 8

  target.style.setProperty('--rx', `${rx.toFixed(2)}deg`)
  target.style.setProperty('--ry', `${ry.toFixed(2)}deg`)
  target.style.setProperty('--px', `${offsetX.toFixed(2)}px`)
  target.style.setProperty('--py', `${offsetY.toFixed(2)}px`)
}

function handleSurfacePointerLeave(event: PointerEvent): void {
  const target = event.currentTarget
  if (!(target instanceof HTMLElement)) {
    return
  }
  target.style.setProperty('--rx', '0deg')
  target.style.setProperty('--ry', '0deg')
  target.style.setProperty('--px', '0px')
  target.style.setProperty('--py', '0px')
}
</script>

<style scoped>
.sidebar-wrap {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.interactive-shell {
  --rx: 0deg;
  --ry: 0deg;
  --px: 0px;
  --py: 0px;
  position: relative;
  border-radius: 22px;
  transform-style: preserve-3d;
  perspective: 1100px;
}

.hot-shell {
  transform-style: flat;
  perspective: none;
}

.tilt-card {
  position: relative;
  z-index: 1;
  transform-style: preserve-3d;
  transform: rotateX(var(--rx)) rotateY(var(--ry));
  transition: transform 220ms ease;
}

.tilt-card :deep(.liquid-glass-content) {
  transform-style: preserve-3d;
}

.depth-1 {
  transform: translate3d(calc(var(--px) * 0.2), calc(var(--py) * 0.18), 10px);
  transition: transform 220ms ease;
}

.depth-2 {
  transform: translate3d(calc(var(--px) * 0.35), calc(var(--py) * 0.3), 16px);
  transition: transform 220ms ease;
}

.depth-3 {
  transform: translate3d(calc(var(--px) * 0.5), calc(var(--py) * 0.45), 24px);
  transition: transform 220ms ease;
}

.profile-top {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
}

.avatar-orbit {
  position: absolute;
  top: -8px;
  width: 84px;
  height: 84px;
  border-radius: 50%;
  background: conic-gradient(
    from 0deg,
    rgba(255, 94, 180, 0.72) 0deg,
    rgba(138, 126, 255, 0.7) 95deg,
    rgba(90, 218, 236, 0.68) 190deg,
    rgba(121, 237, 169, 0.7) 290deg,
    rgba(255, 94, 180, 0.72) 360deg
  );
  filter: blur(9px);
  opacity: 0.82;
  animation: orbit-spin 7s linear infinite;
}

.avatar {
  position: relative;
  width: 62px;
  height: 62px;
  border-radius: 50%;
  display: block;
  object-fit: cover;
  object-position: center;
  background: rgba(234, 255, 248, 0.86);
  box-shadow:
    0 8px 16px rgba(25, 105, 87, 0.16),
    inset 0 1px 0 rgba(255, 255, 255, 0.58);
}

.profile-top h3 {
  margin: 0;
  font-size: 1.24rem;
  color: rgba(8, 52, 43, 0.92);
  margin-top: 10px;
}

.profile-top p {
  margin: 6px 0 0;
  font-size: 0.85rem;
  color: rgba(22, 84, 71, 0.82);
}

.profile-intro {
  margin: 14px 0 0;
  line-height: 1.62;
  color: rgba(20, 78, 66, 0.88);
  font-size: 0.88rem;
  text-align: center;
}

.profile-stats {
  margin-top: 14px;
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 8px;
}

.stat-item {
  border-radius: 11px;
  padding: 8px 8px 9px;
  background: rgba(237, 253, 247, 0.62);
  border: 1px solid rgba(255, 255, 255, 0.46);
  text-align: center;
}

.stat-item span {
  display: block;
  font-size: 0.74rem;
  color: rgba(23, 84, 71, 0.72);
  letter-spacing: 0.02em;
}

.stat-item strong {
  display: block;
  margin-top: 4px;
  font-size: 0.96rem;
  color: rgba(10, 68, 56, 0.95);
  font-weight: 700;
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
  transform-origin: center bottom;
}

.hot-link {
  position: relative;
  z-index: 2;
  display: block;
  width: 100%;
  text-align: left;
  padding: 10px 12px;
  cursor: pointer;
  text-decoration: none;
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

.hot-stagger-enter-active {
  transition:
    opacity 420ms ease,
    transform 540ms cubic-bezier(0.2, 0.7, 0.2, 1);
  transition-delay: var(--stagger-delay, 0ms);
}

.hot-stagger-enter-from {
  opacity: 0;
  transform: translateY(16px) scale(0.97);
}

.hot-stagger-move {
  transition: transform 380ms ease;
}

@keyframes orbit-spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

@media (prefers-reduced-motion: reduce) {
  .tilt-card,
  .depth-1,
  .depth-2,
  .depth-3,
  .hot-stagger-enter-active,
  .hot-stagger-move {
    transition: none;
  }

  .avatar-orbit {
    animation: none;
  }
}
</style>
