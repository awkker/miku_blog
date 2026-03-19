<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-2">
        <div class="flex h-6 w-6 items-center justify-center rounded-md bg-[#c084fc]/10">
          <svg viewBox="0 0 24 24" class="h-3.5 w-3.5 fill-none stroke-[#c084fc] stroke-[2]" aria-hidden="true">
            <path d="M23 3a10.9 10.9 0 01-3.14 1.53 4.48 4.48 0 00-7.86 3v1A10.66 10.66 0 013 4s-4 9 5 13a11.64 11.64 0 01-7 2c9 5 20 0 20-11.5a4.5 4.5 0 00-.08-.83A7.72 7.72 0 0023 3z" />
          </svg>
        </div>
        <h3 class="text-sm font-bold tracking-[0.12em] text-slate-700">{{ copy.title }}</h3>
      </div>
      <a href="/moments" class="text-[11px] text-miku transition hover:underline">{{ copy.viewAll }}</a>
    </div>

    <!-- Loading -->
    <div v-if="status === 'loading'" class="space-y-3">
      <div v-for="i in 3" :key="i" class="animate-pulse space-y-2 rounded-xl bg-white/40 p-3">
        <div class="h-3 w-20 rounded bg-slate-200" />
        <div class="h-3 w-full rounded bg-slate-100" />
        <div class="h-3 w-3/4 rounded bg-slate-100" />
      </div>
    </div>

    <!-- Error -->
    <div v-else-if="status === 'error'" class="rounded-xl bg-red-50/60 px-3 py-4 text-center text-xs text-red-400">
      {{ copy.loadFailed }}
      <button type="button" class="ml-1 text-miku underline" @click="load">{{ copy.retry }}</button>
    </div>

    <!-- Feed -->
    <div v-else class="space-y-2.5">
      <a
        v-for="item in latestList"
        :key="item.id"
        :href="`/moments#${item.id}`"
        class="group block rounded-xl border border-transparent bg-white/40 px-3.5 py-3 transition duration-300 hover:border-miku/20 hover:bg-white/70 hover:shadow-sm"
      >
        <div class="flex items-center gap-2">
          <img :src="item.avatar" :alt="item.nickname" class="h-5 w-5 rounded-full object-cover" loading="lazy" />
          <span class="text-xs font-semibold text-slate-600">{{ item.nickname }}</span>
          <span class="ml-auto text-[10px] text-slate-400">{{ item.createdAt }}</span>
        </div>
        <p class="mt-1.5 line-clamp-2 text-[13px] leading-relaxed text-slate-500 transition group-hover:text-slate-700">{{ item.content }}</p>
        <div v-if="item.images.length > 0" class="mt-2 flex gap-1.5">
          <img
            v-for="(img, idx) in item.images.slice(0, 3)"
            :key="idx"
            :src="img"
            :alt="`${copy.imageAltPrefix}${idx + 1}`"
            class="h-12 w-12 rounded-lg border border-slate-100 object-cover"
            loading="lazy"
          />
          <span v-if="item.images.length > 3" class="flex h-12 w-12 items-center justify-center rounded-lg border border-slate-100 bg-slate-50 text-[11px] text-slate-400">
            +{{ item.images.length - 3 }}
          </span>
        </div>
        <div class="mt-2 flex items-center gap-3 text-[11px] text-slate-400">
          <span class="inline-flex items-center gap-0.5">
            <svg viewBox="0 0 24 24" class="h-3 w-3 fill-none stroke-current stroke-[2]" aria-hidden="true"><path d="M20.84 4.61a5.5 5.5 0 00-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 00-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 000-7.78z" /></svg>
            {{ item.likes }}
          </span>
          <span class="inline-flex items-center gap-0.5">
            <svg viewBox="0 0 24 24" class="h-3 w-3 fill-none stroke-current stroke-[2]" aria-hidden="true"><path d="M21 15a2 2 0 01-2 2H7l-4 4V5a2 2 0 012-2h14a2 2 0 012 2z" /></svg>
            {{ item.comments.length }}
          </span>
        </div>
      </a>

      <div v-if="latestList.length === 0" class="py-4 text-center text-xs text-slate-400">
        {{ copy.empty }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useStore } from '@nanostores/vue'
import { computed, onMounted } from 'vue'
import { loadMoments, moments, momentsFetchStatus } from '../../stores/moments'
import { siteCopy } from '../../content/copy'

const list = useStore(moments)
const status = useStore(momentsFetchStatus)
const copy = siteCopy.components.latestMoments

const latestList = computed(() => [...list.value].slice(0, 3))

async function load() {
  await loadMoments()
}

onMounted(async () => {
  if (list.value.length === 0) {
    await load()
  }
})
</script>
