<template>
  <section class="hero-banner">
    <div class="hero-overlay"></div>
    <div class="hero-slogan">
      <h1>风过山川，字里有光</h1>
      <p>记录风景、代码与日常灵感</p>
    </div>
  </section>
</template>

<style scoped>
.hero-banner {
  position: relative;
  flex: 0 0 clamp(240px, 38vh, 360px);
  width: 100%;
  background-image: url('/photo/bloglist.jpg');
  background-size: cover;
  background-position: center;
  overflow: hidden;
}

.hero-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(180deg, rgba(8, 12, 19, 0.2) 0%, rgba(8, 12, 19, 0.62) 100%);
}

.hero-slogan {
  position: absolute;
  z-index: 1;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: rgba(240, 255, 250, 0.96);
  text-align: center;
  padding: 0 20px;
}

.hero-slogan h1 {
  margin: 0;
  font-family: var(--font-title);
  font-size: clamp(2.1rem, 4vw, 3.4rem);
  letter-spacing: 4px;
  text-shadow: 0 8px 20px rgba(0, 0, 0, 0.45);
}

.hero-slogan p {
  margin: 14px 0 0;
  font-size: clamp(1rem, 1.5vw, 1.2rem);
  letter-spacing: 1px;
}
</style>

