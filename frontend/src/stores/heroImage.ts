import { atom } from 'nanostores'

import { siteCopy } from '../content/copy'

const DEFAULT_IMAGES = [
  '/picture/fengmian/1.webp',
  '/picture/fengmian/2.webp',
  '/picture/fengmian/3.webp',
]

const configuredImages = (siteCopy.home.heroImages || []).filter((item) => item && item.trim().length > 0)

export const heroImages = configuredImages.length > 0 ? configuredImages : DEFAULT_IMAGES

export const heroIndex = atom(0)

export function shuffleHeroImage() {
  const current = heroIndex.get()
  let next: number
  do {
    next = Math.floor(Math.random() * heroImages.length)
  } while (next === current && heroImages.length > 1)
  heroIndex.set(next)
}
