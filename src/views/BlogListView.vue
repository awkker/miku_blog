<template>
  <div class="blog-page" :class="{ 'dock-open': isDockOpen }">
    <FallingPetals />

    <BlogListHero />

    <section class="content-shell">
      <main class="main-pane">
        <BlogTagNav :tags="tags" />

        <div class="post-grid">
          <BlogPostCard
            v-for="post in posts"
            :key="post.slug"
            :post="post"
            :expanded="expandedCoverSlug === post.slug"
            @open="openPost"
            @expand="expandCover"
            @reset="resetCover"
          />
        </div>
      </main>

      <aside class="side-pane">
        <BlogSidebar :hot-posts="hotPosts" @open-post="openPost" />
      </aside>
    </section>

    <BlogDockToggle :is-open="isDockOpen" @show="showDock" @hide="hideDock" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

import FallingPetals from '@/components/FallingPetals.vue'
import BlogDockToggle from '@/components/blog/BlogDockToggle.vue'
import BlogListHero from '@/components/blog/BlogListHero.vue'
import BlogPostCard from '@/components/blog/BlogPostCard.vue'
import BlogSidebar from '@/components/blog/BlogSidebar.vue'
import BlogTagNav from '@/components/blog/BlogTagNav.vue'
import { getAllBlogPostMetas, getHotBlogPostMetas } from '@/content/blogPosts'

const tags = ['全部', '旅行风景', '技术随笔', '日常片段', '摄影后期', '阅读札记', '音乐记录', '生活感悟']

const router = useRouter()
const posts = getAllBlogPostMetas()
const hotPosts = getHotBlogPostMetas()

const isDockOpen = ref(false)
const expandedCoverSlug = ref<string | null>(null)

function showDock(): void {
  isDockOpen.value = true
}

function hideDock(): void {
  isDockOpen.value = false
}

function openPost(slug: string): void {
  void router.push({ name: 'blog-post', params: { slug } })
}

function expandCover(slug: string): void {
  expandedCoverSlug.value = slug
}

function resetCover(slug: string): void {
  if (expandedCoverSlug.value === slug) {
    expandedCoverSlug.value = null
  }
}
</script>

<style scoped>
.blog-page {
  position: relative;
  width: 100%;
  height: 100%;
  max-height: 100dvh;
  display: flex;
  flex-direction: column;
  overflow-x: hidden;
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;
  padding-bottom: 68px;
  background-image: url('/photo/ui.png');
  background-size: cover;
  background-position: center;
}

.blog-page.dock-open {
  padding-bottom: 120px;
}

.content-shell {
  position: relative;
  z-index: 1;
  flex: 0 0 auto;
  width: min(1120px, calc(100% - 64px));
  margin: 0 auto;
  display: grid;
  grid-template-columns: minmax(0, 1fr) 300px;
  gap: 24px;
  padding: 18px 20px 12px;
}

.main-pane,
.side-pane {
  min-width: 0;
  scrollbar-width: thin;
}

.main-pane {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.side-pane {
  position: sticky;
  top: 18px;
  align-self: start;
}

.post-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 14px;
  align-items: start;
  grid-auto-rows: max-content;
}

@media (max-width: 1080px) {
  .content-shell {
    width: calc(100% - 20px);
    grid-template-columns: 1fr;
    padding: 16px 14px 10px;
  }

  .post-grid {
    grid-template-columns: 1fr;
  }

  .side-pane {
    position: static;
    top: auto;
  }
}
</style>
