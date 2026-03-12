<template>
  <aside class="glass-layer rounded-3xl p-5">
    <h3 class="mb-4 text-sm font-bold tracking-[0.2em] text-slate-700">CONTENTS</h3>
    <ul class="space-y-2">
      <li v-for="heading in filteredHeadings" :key="heading.slug" class="relative">
        <a
          :href="`#${heading.slug}`"
          class="block rounded-lg border-l-2 py-1 pl-3 text-sm transition"
          :class="
            activeId === heading.slug
              ? 'border-miku text-miku'
              : 'border-transparent text-slate-500 hover:text-slate-700'
          "
        >
          {{ heading.text }}
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

const filteredHeadings = computed(() => props.headings.filter((item) => item.depth <= 3))

onMounted(() => {
  activeId.value = filteredHeadings.value[0]?.slug ?? ''

  observer = new IntersectionObserver(
    (entries) => {
      const visibleEntries = entries
        .filter((entry) => entry.isIntersecting)
        .sort((a, b) => b.intersectionRatio - a.intersectionRatio)

      if (visibleEntries[0]?.target.id) {
        activeId.value = visibleEntries[0].target.id
      }
    },
    {
      rootMargin: '-25% 0px -60% 0px',
      threshold: [0.1, 0.3, 0.6],
    },
  )

  filteredHeadings.value.forEach((heading) => {
    const element = document.getElementById(heading.slug)
    if (element) {
      observer?.observe(element)
    }
  })
})

onBeforeUnmount(() => {
  observer?.disconnect()
})
</script>
