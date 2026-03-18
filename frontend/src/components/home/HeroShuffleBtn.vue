<template>
  <button
    class="group flex w-12 flex-col items-center gap-1 rounded-2xl py-2 text-[10px] text-[#39c5bb]/90 transition duration-300 hover:scale-110 hover:bg-white/20 sm:w-16 sm:text-xs sm:hover:scale-125"
    :title="btnTitle"
    @click="onShuffle"
  >
    <span
      class="flex h-8 w-8 items-center justify-center rounded-xl border border-white/30 bg-white/15 transition-transform duration-500 sm:h-9 sm:w-9"
      :class="{ 'rotate-180': spinning }"
    >
      <svg viewBox="0 0 24 24" class="h-4 w-4 fill-none stroke-current stroke-[1.5] sm:h-5 sm:w-5">
        <path d="M16 3h5v5M4 20L21 3M21 16v5h-5M4 4l17 17" />
      </svg>
    </span>
    <span>换图</span>
    <span class="h-1.5 w-1.5 rounded-full bg-white opacity-0" />
  </button>
</template>

<script setup lang="ts">
import { useStore } from '@nanostores/vue'
import { computed, onMounted, ref, watch } from 'vue'

import { heroImages, heroIndex, shuffleHeroImage } from '../../stores/heroImage'

const $heroIndex = useStore(heroIndex)
const total = heroImages.length
const currentIdx = ref(0)
const mounted = ref(false)
const spinning = ref(false)

const btnTitle = computed(() =>
  mounted.value ? `切换封面 (${currentIdx.value + 1}/${total})` : '切换封面'
)

watch($heroIndex, (v) => {
  currentIdx.value = v
})

onMounted(() => {
  mounted.value = true
  currentIdx.value = heroIndex.get()
})

function onShuffle() {
  spinning.value = true
  shuffleHeroImage()
  setTimeout(() => {
    spinning.value = false
  }, 500)
}
</script>
