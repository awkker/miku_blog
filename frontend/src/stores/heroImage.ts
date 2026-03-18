import { atom } from 'nanostores'

const IMAGES = [
  '/picture/fengmian/1.jpg',
  '/picture/fengmian/2.jpg',
  '/picture/fengmian/3.jpg',
]

export const heroImages = IMAGES

export const heroIndex = atom(0)

export function shuffleHeroImage() {
  const current = heroIndex.get()
  let next: number
  do {
    next = Math.floor(Math.random() * IMAGES.length)
  } while (next === current && IMAGES.length > 1)
  heroIndex.set(next)
}
