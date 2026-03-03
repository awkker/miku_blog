<template>
  <div class="blog-page" :class="{ 'dock-open': isDockOpen }">
    <FallingPetals />

    <BlogListHero />

    <section class="content-shell">
      <main class="main-pane">
        <BlogTagNav :tags="tags" :active-tag="selectedTag" @select="selectTag" />

        <div class="post-grid">
          <BlogPostCard
            v-for="post in filteredPosts"
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
        <BlogSidebar :hot-posts="hotPosts" :author-stats="authorStats" />
      </aside>
    </section>

    <BlogDockToggle :is-open="isDockOpen" @show="showDock" @hide="hideDock" />
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

import FallingPetals from '@/components/FallingPetals.vue'
import BlogDockToggle from '@/components/blog/BlogDockToggle.vue'
import BlogListHero from '@/components/blog/BlogListHero.vue'
import BlogPostCard from '@/components/blog/BlogPostCard.vue'
import BlogSidebar from '@/components/blog/BlogSidebar.vue'
import BlogTagNav from '@/components/blog/BlogTagNav.vue'
import {
  getAllBlogPostMetas,
  getBlogAuthorStats,
  getBlogStatsUpdatedEventName,
  getHotBlogPostMetas,
  incrementBlogListPageVisits,
} from '@/content/blogPosts'

const router = useRouter()
const allPosts = ref(getAllBlogPostMetas())
const hotPosts = ref(getHotBlogPostMetas())
const authorStats = ref(getBlogAuthorStats())
const blogStatsUpdatedEventName = getBlogStatsUpdatedEventName()

const isDockOpen = ref(false)
const expandedCoverSlug = ref<string | null>(null)
const selectedTag = ref('全部')

const tags = computed(() => {
  const uniqueCategories = new Set(allPosts.value.map((post) => post.category))
  return ['全部', ...Array.from(uniqueCategories)]
})

const filteredPosts = computed(() => {
  if (selectedTag.value === '全部') {
    return allPosts.value
  }
  return allPosts.value.filter((post) => post.category === selectedTag.value)
})

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

function selectTag(tag: string): void {
  selectedTag.value = tag
  expandedCoverSlug.value = null
}

function refreshBlogData(): void {
  allPosts.value = getAllBlogPostMetas()
  hotPosts.value = getHotBlogPostMetas()
  authorStats.value = getBlogAuthorStats()
}

function handleStatsUpdated(): void {
  refreshBlogData()
}

function handleStorageChange(event: StorageEvent): void {
  if (!event.key || event.key.startsWith('miku_blog_')) {
    refreshBlogData()
  }
}

onMounted(() => {
  incrementBlogListPageVisits()
  refreshBlogData()
  window.addEventListener(blogStatsUpdatedEventName, handleStatsUpdated)
  window.addEventListener('storage', handleStorageChange)
})

onBeforeUnmount(() => {
  window.removeEventListener(blogStatsUpdatedEventName, handleStatsUpdated)
  window.removeEventListener('storage', handleStorageChange)
})
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
