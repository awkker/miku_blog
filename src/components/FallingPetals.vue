<template>
  <canvas ref="canvasRef" class="falling-petals-canvas"></canvas>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'

// 单片花瓣的运动与绘制参数
interface Petal {
  x: number
  y: number
  size: number
  color: string
  speedX: number
  speedY: number
  angle: number
  angleSpeed: number
  swing: number
  swingSpeed: number
}

const canvasRef = ref<HTMLCanvasElement | null>(null)
let animationId = 0

// 当前屏幕内所有花瓣实例
const petals: Petal[] = []
const config = {
  count: 30,
  colors: ['#ffb3ba', '#ffdfba', '#ffffba', '#baffc9', '#bae1ff', '#e8baff'],
}

// 生成新花瓣：初始时散布全屏，后续重生时从左上区域进入
const createPetal = (width: number, height: number, isInitial = false): Petal => {
  const startX = isInitial ? Math.random() * width : Math.random() * width - width * 0.3
  const startY = isInitial ? Math.random() * height : -30
  const colorIndex = Math.floor(Math.random() * config.colors.length)

  return {
    x: startX,
    y: startY,
    size: Math.random() * 6 + 8,
    color: config.colors[colorIndex] ?? config.colors[0] ?? '#ffb3ba',
    speedX: Math.random() * 1.2 + 0.8,
    speedY: Math.random() * 1.5 + 1.2,
    angle: Math.random() * 360,
    angleSpeed: (Math.random() - 0.5) * 3,
    swing: Math.random() * Math.PI * 2,
    swingSpeed: Math.random() * 0.03 + 0.01,
  }
}

let cleanupResize: (() => void) | null = null

onMounted(() => {
  const canvas = canvasRef.value
  if (!canvas) {
    return
  }

  const ctx = canvas.getContext('2d')
  if (!ctx) {
    return
  }

  const updateCanvasSize = () => {
    // 按设备像素比放大画布，避免高分屏发虚
    const dpr = Math.min(window.devicePixelRatio || 1, 2)
    const width = window.innerWidth
    const height = window.innerHeight
    canvas.width = Math.floor(width * dpr)
    canvas.height = Math.floor(height * dpr)
    ctx.setTransform(dpr, 0, 0, dpr, 0, 0)
  }

  window.addEventListener('resize', updateCanvasSize)
  cleanupResize = () => window.removeEventListener('resize', updateCanvasSize)
  updateCanvasSize()

  petals.length = 0
  // 首屏先铺满一批花瓣
  for (let i = 0; i < config.count; i += 1) {
    petals.push(createPetal(window.innerWidth, window.innerHeight, true))
  }

  const render = () => {
    // 每帧先清屏，再重绘全部花瓣
    ctx.clearRect(0, 0, window.innerWidth, window.innerHeight)

    petals.forEach((petal, index) => {
      // 位移：整体向右下飘，同时附加小幅摆动
      petal.swing += petal.swingSpeed
      petal.x += petal.speedX + Math.sin(petal.swing) * 0.8
      petal.y += petal.speedY
      petal.angle += petal.angleSpeed

      // 以花瓣中心点旋转后绘制两段贝塞尔曲线
      ctx.save()
      ctx.translate(petal.x, petal.y)
      ctx.rotate((petal.angle * Math.PI) / 180)

      ctx.fillStyle = petal.color
      ctx.beginPath()
      ctx.moveTo(0, 0)
      ctx.bezierCurveTo(-petal.size, -petal.size * 0.5, -petal.size, petal.size * 1.2, 0, petal.size)
      ctx.bezierCurveTo(petal.size, petal.size * 1.2, petal.size, -petal.size * 0.5, 0, 0)
      ctx.fill()
      ctx.restore()

      // 飞出画面后立即重生，维持固定数量
      if (petal.x > window.innerWidth + 30 || petal.y > window.innerHeight + 30) {
        petals[index] = createPetal(window.innerWidth, window.innerHeight, false)
      }
    })

    animationId = window.requestAnimationFrame(render)
  }

  render()
})

onBeforeUnmount(() => {
  // 组件销毁时释放监听与动画循环
  cleanupResize?.()
  cleanupResize = null
  window.cancelAnimationFrame(animationId)
})
</script>

<style scoped>
.falling-petals-canvas {
  /* 固定铺满全屏，且不拦截鼠标事件 */
  position: fixed;
  inset: 0;
  width: 100vw;
  height: 100vh;
  pointer-events: none;
  z-index: 0;
}
</style>
