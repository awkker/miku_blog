<template>
  <div class="home-container">
    <FallingPetals />

    <div class="content-layer">
      <!-- 顶部导航（液态玻璃） -->
      <AppTopNav />

      <main class="center-content">
        <!-- 页面主标题 -->
        <h1 class="art-title">薰逸の猫窝</h1>
        <!-- 中部文案卡片（液态玻璃 + 打字机） -->
        <TypewriterQuoteGlass />
      </main>

      <!-- 底部 Dock（液态玻璃） -->
      <AppDock />
    </div>
  </div>
</template>

<script setup lang="ts">
import AppDock from '@/components/AppDock.vue'
import AppTopNav from '@/components/AppTopNav.vue'
import FallingPetals from '@/components/FallingPetals.vue'
import TypewriterQuoteGlass from '@/components/TypewriterQuoteGlass.vue'
</script>

<style scoped>
.home-container {
  position: relative;
  width: 100%;
  height: 100dvh;
  min-height: 100vh;
  background-image: url('/photo/封面.avif');
  background-size: cover;
  background-position: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: space-between;
  overflow: hidden;
  font-family: var(--font-body);
}

.content-layer {
  position: relative;
  z-index: 1;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: space-between;
}

.center-content {
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  transform: translateY(-50px);
}

.art-title {
  font-family: var(--font-title);
  font-size: 5rem;
  color: rgb(112, 228, 205);
  font-weight: 700;
  letter-spacing: 5px;
  transition: all 0.4s ease;
  cursor: default;
  margin-bottom: 30px;
  text-shadow: 0 6px 22px rgba(0, 0, 0, 0.28);
}

.art-title:hover {
  background: linear-gradient(90deg, #4facfe 0%, #00f2fe 50%, #8a2be2 100%);
  background-size: 200% auto;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  animation: shine 3s linear infinite;
}

@keyframes shine {
  to { background-position: 200% center; }
}
</style>
