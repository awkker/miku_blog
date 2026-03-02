<template>
  <LiquidGlassFrame
    class="pill-nav-glass"
    width="fit-content"
    max-width="none"
    padding="12px 40px"
    :border-radius="50"
    :displacement-strength="1.08"
    :edge-refraction-strength="1.25"
    :interactive="false"
  >
    <nav class="pill-nav">
      <a v-for="item in items" :key="item.label" :href="item.href">{{ item.label }}</a>
    </nav>
  </LiquidGlassFrame>
</template>

<script setup lang="ts">
import LiquidGlassFrame from '@/components/LiquidGlassFrame.vue'

interface NavItem {
  label: string
  href: string
}

const items: NavItem[] = [
  // 后续如果接路由，可把 href 改为对应 path / router-link
  { label: '首页', href: '#' },
  { label: '关于', href: '#' },
  { label: '留言板', href: '#' },
  { label: '友情链接', href: '#' },
]
</script>

<style scoped>
.pill-nav-glass {
  margin-top: 30px;
}

.pill-nav {
  display: flex;
  gap: 30px;
}

.pill-nav a {
  text-decoration: none;
  color: var(--miku-color);
  font-weight: 600;
  font-size: 1.1rem;
  transition: opacity 0.3s;
}

.pill-nav a:hover {
  opacity: 0.7;
}
</style>
