<template>
  <aside class="rounded-3xl border border-slate-200/90 bg-white p-4 shadow-[0_12px_30px_rgba(15,23,42,0.08)]">
    <h3 class="mb-3 text-xs font-bold tracking-[0.18em] text-slate-700">阅读目录</h3>
    <ul class="space-y-1.5">
      <li v-for="heading in filteredHeadings" :key="heading.slug" class="relative">
        <a
          :href="`#${heading.slug}`"
          class="group block rounded-lg border-l-2 py-1.5 pr-2 text-[13px] leading-5 transition"
          :class="
            activeId === heading.slug
              ? 'border-miku bg-miku-soft text-miku'
              : 'border-transparent text-slate-500 hover:border-miku/40 hover:bg-white/45 hover:text-slate-700'
          "
          :style="{ paddingLeft: `${heading.depth === 3 ? 22 : 12}px` }"
          @click.prevent="scrollToHeading(heading.slug)"
        >
          <span class="line-clamp-2">{{ heading.text }}</span>
        </a>
      </li>
    </ul>
  </aside>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'

interface HeadingItem {
  depth: number
  slug: string
  text: string
}

interface Props {
  headings?: HeadingItem[]
}

const props = withDefaults(defineProps<Props>(), {
  headings: () => [],
})
const activeId = ref('')
let observer: IntersectionObserver | null = null
let rafId = 0

const filteredHeadings = computed(() => props.headings.filter((item) => item.depth <= 3))

function updateActiveFromScroll() {
  const headingElements = filteredHeadings.value
    .map((heading) => document.getElementById(heading.slug))
    .filter((element): element is HTMLElement => Boolean(element))

  if (headingElements.length === 0) {
    return
  }

  const triggerLine = 160
  let candidate = headingElements[0]

  for (const element of headingElements) {
    if (element.getBoundingClientRect().top - triggerLine <= 0) {
      candidate = element
    } else {
      break
    }
  }

  activeId.value = candidate.id
}

function onScroll() {
  cancelAnimationFrame(rafId)
  rafId = requestAnimationFrame(updateActiveFromScroll)
}

function scrollToHeading(slug: string) {
  const target = document.getElementById(slug)
  if (!target) {
    return
  }

  target.scrollIntoView({ behavior: 'smooth', block: 'start' })
  history.replaceState(null, '', `#${slug}`)
  activeId.value = slug
}

onMounted(() => {
  activeId.value = filteredHeadings.value[0]?.slug ?? ''

  observer = new IntersectionObserver(
    (entries) => {
      const visibleEntries = entries.filter((entry) => entry.isIntersecting)
      if (visibleEntries.length > 0) {
        const topEntry = visibleEntries.reduce((closest, entry) =>
          entry.boundingClientRect.top < closest.boundingClientRect.top ? entry : closest,
        )
        if (topEntry.target.id) {
          activeId.value = topEntry.target.id
        }
      }
    },
    {
      rootMargin: '-20% 0px -62% 0px',
      threshold: [0, 0.2, 0.5, 1],
    },
  )

  filteredHeadings.value.forEach((heading) => {
    const element = document.getElementById(heading.slug)
    if (element) {
      observer?.observe(element)
    }
  })

  window.addEventListener('scroll', onScroll, { passive: true })
  updateActiveFromScroll()
})

onBeforeUnmount(() => {
  cancelAnimationFrame(rafId)
  window.removeEventListener('scroll', onScroll)
  observer?.disconnect()
})
</script>
