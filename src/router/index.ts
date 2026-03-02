import { createRouter, createWebHistory } from 'vue-router'
import BlogListView from '@/views/BlogListView.vue'
import BlogPostView from '@/views/BlogPostView.vue'
import HomeView from '@/views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/blog',
      name: 'blog',
      component: BlogListView,
    },
    {
      path: '/blog/:slug',
      name: 'blog-post',
      component: BlogPostView,
    },
  ],
})

export default router
