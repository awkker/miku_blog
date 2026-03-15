export interface LiquidGlassOptions {
  borderRadius?: number
  cornerSoftness?: number
  displacementStrength?: number
  edgeRefractionStrength?: number
  blur?: number
  contrast?: number
  brightness?: number
  saturate?: number
  interactive?: boolean
}

export interface LiquidGlassController {
  update: (next?: Partial<LiquidGlassOptions>) => void
  destroy: () => void
}

const SVG_NS = 'http://www.w3.org/2000/svg'
const XLINK_NS = 'http://www.w3.org/1999/xlink'

function smoothStep(a: number, b: number, t: number): number {
  const normalized = Math.max(0, Math.min(1, (t - a) / (b - a)))
  return normalized * normalized * (3 - 2 * normalized)
}

function length2d(x: number, y: number): number {
  return Math.sqrt(x * x + y * y)
}

function roundedRectSdf(x: number, y: number, halfWidth: number, halfHeight: number, radius: number): number {
  const qx = Math.abs(x) - halfWidth + radius
  const qy = Math.abs(y) - halfHeight + radius
  return Math.min(Math.max(qx, qy), 0) + length2d(Math.max(qx, 0), Math.max(qy, 0)) - radius
}

function texture(x: number, y: number): { x: number; y: number } {
  return { x, y }
}

function generateId(): string {
  return `liquid-glass-${Math.random().toString(36).slice(2, 11)}`
}

function toCssFilter(options: Required<LiquidGlassOptions>, filterId: string): string {
  return `url(#${filterId}) blur(${options.blur}px) contrast(${options.contrast}) brightness(${options.brightness}) saturate(${options.saturate})`
}

function toFallbackCssFilter(options: Required<LiquidGlassOptions>): string {
  return `blur(${options.blur}px) contrast(${options.contrast}) brightness(${options.brightness}) saturate(${options.saturate})`
}

const defaultOptions: Required<LiquidGlassOptions> = {
  borderRadius: 24,
  cornerSoftness: 0.12,
  displacementStrength: 1,
  edgeRefractionStrength: 0.75,
  blur: 0.3,
  contrast: 1.14,
  brightness: 1.04,
  saturate: 1.08,
  interactive: true,
}

export function createLiquidGlass(target: HTMLElement, options: LiquidGlassOptions = {}): LiquidGlassController {
  let mergedOptions: Required<LiquidGlassOptions> = { ...defaultOptions, ...options }
  const id = generateId()

  const svg = document.createElementNS(SVG_NS, 'svg')
  svg.setAttribute('xmlns', SVG_NS)
  svg.setAttribute('width', '0')
  svg.setAttribute('height', '0')
  svg.style.position = 'fixed'
  svg.style.top = '0'
  svg.style.left = '0'
  svg.style.pointerEvents = 'none'
  svg.style.zIndex = '-1'

  const defs = document.createElementNS(SVG_NS, 'defs')
  const filter = document.createElementNS(SVG_NS, 'filter')
  filter.setAttribute('id', id)
  filter.setAttribute('filterUnits', 'userSpaceOnUse')
  filter.setAttribute('color-interpolation-filters', 'sRGB')

  const feImage = document.createElementNS(SVG_NS, 'feImage')
  feImage.setAttribute('id', `${id}-map`)

  const feDisplacementMap = document.createElementNS(SVG_NS, 'feDisplacementMap')
  feDisplacementMap.setAttribute('in', 'SourceGraphic')
  feDisplacementMap.setAttribute('in2', `${id}-map`)
  feDisplacementMap.setAttribute('xChannelSelector', 'R')
  feDisplacementMap.setAttribute('yChannelSelector', 'G')

  filter.appendChild(feImage)
  filter.appendChild(feDisplacementMap)
  defs.appendChild(filter)
  svg.appendChild(defs)
  document.body.appendChild(svg)

  const canvas = document.createElement('canvas')
  const context = canvas.getContext('2d')
  const previousInlineStyle = {
    backdropFilter: target.style.backdropFilter,
    webkitBackdropFilter: target.style.getPropertyValue('-webkit-backdrop-filter'),
  }
  const setBackdropFilter = (value: string): void => {
    target.style.backdropFilter = value
    target.style.setProperty('-webkit-backdrop-filter', value)
  }
  const restoreBackdropFilter = (): void => {
    target.style.backdropFilter = previousInlineStyle.backdropFilter
    if (previousInlineStyle.webkitBackdropFilter) {
      target.style.setProperty('-webkit-backdrop-filter', previousInlineStyle.webkitBackdropFilter)
    } else {
      target.style.removeProperty('-webkit-backdrop-filter')
    }
  }

  if (!context) {
    const applyFallbackStyle = (): void => {
      const fallbackFilter = toFallbackCssFilter(mergedOptions)
      setBackdropFilter(fallbackFilter)
    }

    svg.remove()
    applyFallbackStyle()

    return {
      update(next: Partial<LiquidGlassOptions> = {}): void {
        mergedOptions = { ...mergedOptions, ...next }
        applyFallbackStyle()
      },
      destroy(): void {
        restoreBackdropFilter()
        canvas.remove()
      },
    }
  }

  setBackdropFilter(toCssFilter(mergedOptions, id))

  let mouseX = 0.5
  let mouseY = 0.5
  const mouseMoveHandler = (event: MouseEvent): void => {
    const rect = target.getBoundingClientRect()
    if (!rect.width || !rect.height) {
      return
    }
    mouseX = Math.max(0, Math.min(1, (event.clientX - rect.left) / rect.width))
    mouseY = Math.max(0, Math.min(1, (event.clientY - rect.top) / rect.height))
    render()
  }

  if (mergedOptions.interactive) {
    target.addEventListener('mousemove', mouseMoveHandler)
  }

  let disposed = false

  const render = (): void => {
    if (disposed) {
      return
    }
    const rect = target.getBoundingClientRect()
    const dpr = Math.min(window.devicePixelRatio || 1, 2)
    const width = Math.max(1, Math.floor(rect.width * dpr))
    const height = Math.max(1, Math.floor(rect.height * dpr))

    canvas.width = width
    canvas.height = height

    filter.setAttribute('x', '0')
    filter.setAttribute('y', '0')
    filter.setAttribute('width', `${rect.width}`)
    filter.setAttribute('height', `${rect.height}`)
    feImage.setAttribute('width', `${rect.width}`)
    feImage.setAttribute('height', `${rect.height}`)

    const data = new Uint8ClampedArray(width * height * 4)
    const displacements: number[] = []
    let maxScale = 0

    for (let i = 0; i < data.length; i += 4) {
      const x = (i / 4) % width
      const y = Math.floor(i / 4 / width)
      const uvx = x / width
      const uvy = y / height

      const ix = uvx - 0.5
      const iy = uvy - 0.5
      const halfWidth = 0.5
      const halfHeight = 0.5
      const maxRadius = Math.min(halfWidth, halfHeight)
      const normalizedRadius = Math.min(mergedOptions.borderRadius / Math.max(rect.width, 1), maxRadius)
      const edgeDistance = roundedRectSdf(ix, iy, halfWidth, halfHeight, normalizedRadius)
      const edgeDisplacement = smoothStep(1, 0, edgeDistance - mergedOptions.cornerSoftness)
      const scaled = smoothStep(0, 1, edgeDisplacement)
      const edgeRing = 1 - smoothStep(-0.16, 0.24, edgeDistance)
      const edgeBoost = 1 + edgeRing * mergedOptions.edgeRefractionStrength

      let mapped = texture(ix * scaled + 0.5, iy * scaled + 0.5)

      if (mergedOptions.interactive) {
        const dxMouse = uvx - mouseX
        const dyMouse = uvy - mouseY
        const mouseDistance = Math.sqrt(dxMouse * dxMouse + dyMouse * dyMouse)
        const mouseInfluence = Math.exp(-mouseDistance * 16) * 0.06
        mapped = texture(mapped.x + dxMouse * mouseInfluence, mapped.y + dyMouse * mouseInfluence)
      }

      const dx = (mapped.x * width - x) * edgeBoost
      const dy = (mapped.y * height - y) * edgeBoost
      maxScale = Math.max(maxScale, Math.abs(dx), Math.abs(dy))
      displacements.push(dx, dy)
    }

    maxScale = Math.max(maxScale * 0.5 * mergedOptions.displacementStrength, 0.0001)

    let displacementIndex = 0
    for (let i = 0; i < data.length; i += 4) {
      const r = (displacements[displacementIndex++] ?? 0) / maxScale + 0.5
      const g = (displacements[displacementIndex++] ?? 0) / maxScale + 0.5
      data[i] = Math.round(r * 255)
      data[i + 1] = Math.round(g * 255)
      data[i + 2] = 0
      data[i + 3] = 255
    }

    context.putImageData(new ImageData(data, width, height), 0, 0)
    feImage.setAttributeNS(XLINK_NS, 'href', canvas.toDataURL())
    feDisplacementMap.setAttribute('scale', `${maxScale / dpr}`)
  }

  const resizeObserver = typeof ResizeObserver !== 'undefined' ? new ResizeObserver(() => render()) : null
  if (resizeObserver) {
    resizeObserver.observe(target)
  } else {
    window.addEventListener('resize', render)
  }
  render()

  return {
    update(next: Partial<LiquidGlassOptions> = {}): void {
      if (disposed) {
        return
      }

      const interactiveBefore = mergedOptions.interactive
      mergedOptions = { ...mergedOptions, ...next }
      setBackdropFilter(toCssFilter(mergedOptions, id))

      if (!interactiveBefore && mergedOptions.interactive) {
        target.addEventListener('mousemove', mouseMoveHandler)
      }
      if (interactiveBefore && !mergedOptions.interactive) {
        target.removeEventListener('mousemove', mouseMoveHandler)
      }

      render()
    },
    destroy(): void {
      if (disposed) {
        return
      }
      disposed = true
      resizeObserver?.disconnect()
      if (!resizeObserver) {
        window.removeEventListener('resize', render)
      }
      target.removeEventListener('mousemove', mouseMoveHandler)
      restoreBackdropFilter()
      canvas.remove()
      svg.remove()
    },
  }
}
