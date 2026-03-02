import { onBeforeUnmount, onMounted, ref } from 'vue'

export interface TypewriterOptions {
  startDelay?: number
  charDelay?: number
  lineDelay?: number
}

const defaultOptions: Required<TypewriterOptions> = {
  startDelay: 380,
  charDelay: 68,
  lineDelay: 320,
}

export function useTypewriter(lines: string[], options: TypewriterOptions = {}) {
  const mergedOptions = { ...defaultOptions, ...options }
  // 每一行“当前已经打出来”的文本
  const typedLines = ref(lines.map(() => ''))
  // 当前光标所在的行（用于显示闪烁光标）
  const activeLineIndex = ref(0)
  let timer: number | null = null

  const clearTimer = (): void => {
    if (timer !== null) {
      window.clearTimeout(timer)
      timer = null
    }
  }

  const start = (): void => {
    clearTimer()
    typedLines.value = lines.map(() => '')
    activeLineIndex.value = 0

    let lineIndex = 0
    let charIndex = 0

    const tick = (): void => {
      if (lineIndex >= lines.length) {
        activeLineIndex.value = Math.max(0, lines.length - 1)
        return
      }

      const currentLine = lines[lineIndex]
      if (!currentLine) {
        activeLineIndex.value = Math.max(0, lines.length - 1)
        return
      }

      charIndex += 1
      typedLines.value[lineIndex] = currentLine.slice(0, charIndex)

      if (charIndex < currentLine.length) {
        timer = window.setTimeout(tick, mergedOptions.charDelay)
        return
      }

      if (lineIndex === lines.length - 1) {
        activeLineIndex.value = lineIndex
        return
      }

      lineIndex += 1
      charIndex = 0
      activeLineIndex.value = lineIndex
      // 当前行打完后，停顿一小段再进入下一行
      timer = window.setTimeout(tick, mergedOptions.lineDelay)
    }

    // 初始延时，避免页面刚加载时过于突兀
    timer = window.setTimeout(tick, mergedOptions.startDelay)
  }

  onMounted(() => {
    start()
  })

  onBeforeUnmount(() => {
    clearTimer()
  })

  return {
    typedLines,
    activeLineIndex,
    restart: start,
  }
}
