<template>
  <section class="article-hero" :style="{ backgroundImage: `url(${post.cover})` }">
    <div class="hero-mask"></div>
    <div class="hero-inner">
      <p class="hero-meta">{{ post.date }} · {{ post.category }}</p>
      <h1>{{ post.title }}</h1>
      <p>{{ readingMinutes }} 分钟阅读</p>
    </div>
  </section>
</template>

<script setup lang="ts">
import type { BlogPost } from '@/content/blogPosts'

interface Props {
  post: BlogPost
  readingMinutes: number
}

defineProps<Props>()
</script>

<style scoped>
.article-hero {
  position: relative;
  min-height: clamp(220px, 36vh, 320px);
  background-size: cover;
  background-position: center;
}

.hero-mask {
  position: absolute;
  inset: 0;
  background: linear-gradient(180deg, rgba(9, 14, 25, 0.26) 0%, rgba(9, 14, 25, 0.68) 100%);
}

.hero-inner {
  position: relative;
  z-index: 1;
  width: min(980px, calc(100% - 44px));
  margin: 0 auto;
  padding: 52px 0 34px;
  color: rgba(245, 252, 255, 0.96);
}

.hero-meta {
  margin: 0;
  font-size: 0.88rem;
  letter-spacing: 0.4px;
}

.hero-inner h1 {
  margin: 12px 0 10px;
  font-family: var(--font-title);
  font-size: clamp(1.8rem, 3vw, 2.8rem);
  line-height: 1.35;
}

.hero-inner p {
  margin: 0;
  color: rgba(222, 241, 250, 0.92);
}
</style>

