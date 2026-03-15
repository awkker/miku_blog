<template>
  <div class="flex items-center gap-2">
    <button @click="expanded = !expanded" class="max-w-[140px] truncate rounded-full bg-white/15 px-2 py-0.5 transition hover:bg-white/25">
      {{ currentTrack.title }}
    </button>
    <button @click="prev" class="rounded p-1 transition hover:bg-white/15" aria-label="上一首">
      <svg viewBox="0 0 24 24" class="h-3.5 w-3.5 fill-current"><path d="M6 6h2v12H6zm3 6l9-6v12z" /></svg>
    </button>
    <button @click="togglePlay" class="rounded p-1 transition hover:bg-white/15" :aria-label="isPlaying ? '暂停' : '播放'">
      <svg v-if="isPlaying" viewBox="0 0 24 24" class="h-3.5 w-3.5 fill-current"><path d="M7 5h4v14H7zm6 0h4v14h-4z" /></svg>
      <svg v-else viewBox="0 0 24 24" class="h-3.5 w-3.5 fill-current"><path d="M8 5v14l11-7z" /></svg>
    </button>
    <button @click="next" class="rounded p-1 transition hover:bg-white/15" aria-label="下一首">
      <svg viewBox="0 0 24 24" class="h-3.5 w-3.5 fill-current"><path d="M16 6h2v12h-2zM7 6l9 6-9 6z" /></svg>
    </button>
    <input type="range" min="0" max="100" :value="Math.round(volume * 100)" @input="onVolumeInput" class="player-range h-[3px] w-14 cursor-pointer" aria-label="音量" />

    <Transition name="fade">
      <div v-if="expanded" class="fixed inset-0 z-40 bg-black/20" @click="expanded = false" />
    </Transition>

    <Transition name="card-slide">
      <div v-if="expanded" class="fixed right-4 top-10 z-50 w-80 rounded-2xl border border-white/15 bg-slate-900/88 p-5 shadow-[0_20px_60px_rgba(0,0,0,0.5)] backdrop-blur-2xl">
        <button @click="expanded = false" class="absolute right-3 top-3 rounded-full p-1 text-white/40 transition hover:bg-white/10 hover:text-white/70" aria-label="关闭">
          <svg viewBox="0 0 24 24" class="h-4 w-4 fill-none stroke-current stroke-[2]"><path d="M18 6L6 18M6 6l12 12" /></svg>
        </button>

        <div class="flex items-center gap-4">
          <img :src="currentTrack.cover" :alt="currentTrack.title" class="h-16 w-16 rounded-xl border border-white/10 object-cover shadow-lg" />
          <div class="min-w-0 flex-1">
            <h3 class="truncate text-sm font-semibold text-white">{{ currentTrack.title }}</h3>
            <p class="mt-0.5 truncate text-xs text-white/50">{{ currentTrack.artist }}</p>
          </div>
        </div>

        <div class="mt-4">
          <input type="range" min="0" :max="duration || 0" :value="currentTime" @input="onSeek" step="0.1" class="player-range-lg w-full cursor-pointer" />
          <div class="mt-1 flex justify-between text-[10px] text-white/35">
            <span>{{ formatTime(currentTime) }}</span>
            <span>{{ formatTime(duration) }}</span>
          </div>
        </div>

        <div class="mt-1 flex items-center justify-center gap-4">
          <button @click="loopMode = !loopMode" class="rounded-full p-1.5 transition" :class="loopMode ? 'text-[#39c5bb]' : 'text-white/35 hover:text-white/60'" aria-label="循环">
            <svg viewBox="0 0 24 24" class="h-4 w-4 fill-current"><path d="M7 7h10v3l4-4-4-4v3H5v6h2V7zm10 10H7v-3l-4 4 4 4v-3h12v-6h-2v4z" /></svg>
          </button>
          <button @click="prev" class="rounded-full p-2 text-white/55 transition hover:text-white" aria-label="上一首">
            <svg viewBox="0 0 24 24" class="h-5 w-5 fill-current"><path d="M6 6h2v12H6zm3 6l9-6v12z" /></svg>
          </button>
          <button @click="togglePlay" class="rounded-full bg-[#39c5bb] p-3 text-slate-900 shadow-lg shadow-[#39c5bb]/25 transition hover:bg-[#4dd4c8]" :aria-label="isPlaying ? '暂停' : '播放'">
            <svg v-if="isPlaying" viewBox="0 0 24 24" class="h-5 w-5 fill-current"><path d="M7 5h4v14H7zm6 0h4v14h-4z" /></svg>
            <svg v-else viewBox="0 0 24 24" class="h-5 w-5 fill-current"><path d="M8 5v14l11-7z" /></svg>
          </button>
          <button @click="next" class="rounded-full p-2 text-white/55 transition hover:text-white" aria-label="下一首">
            <svg viewBox="0 0 24 24" class="h-5 w-5 fill-current"><path d="M16 6h2v12h-2zM7 6l9 6-9 6z" /></svg>
          </button>
          <button @click="muted = !muted" class="rounded-full p-1.5 transition" :class="muted ? 'text-[#39c5bb]' : 'text-white/35 hover:text-white/60'" :aria-label="muted ? '取消静音' : '静音'">
            <svg v-if="muted" viewBox="0 0 24 24" class="h-4 w-4 fill-none stroke-current stroke-[2]"><path d="M11 5L6 9H2v6h4l5 4V5z" /><line x1="23" y1="9" x2="17" y2="15" /><line x1="17" y1="9" x2="23" y2="15" /></svg>
            <svg v-else viewBox="0 0 24 24" class="h-4 w-4 fill-none stroke-current stroke-[2]"><path d="M11 5L6 9H2v6h4l5 4V5z" /><path d="M15.54 8.46a5 5 0 010 7.07" /></svg>
          </button>
        </div>

        <div ref="lyricsRef" class="scrollbar-hide mt-3 h-36 overflow-y-auto rounded-xl bg-white/5 px-3 py-2">
          <div v-if="lyrics.length === 0" class="flex h-full items-center justify-center text-xs text-white/25">暂无歌词</div>
          <div v-else class="space-y-1.5 py-14 text-center">
            <p
              v-for="(line, i) in lyrics" :key="i" :data-idx="i"
              class="cursor-pointer px-1 py-0.5 text-xs leading-relaxed transition-all duration-300"
              :class="i === currentLyricIndex ? 'scale-[1.05] font-medium text-[#39c5bb]' : Math.abs(i - currentLyricIndex) <= 1 ? 'text-white/35' : 'text-white/18'"
              @click="seekToLyric(line.time)"
            >{{ line.text }}</p>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'

interface Track {
  title: string
  artist: string
  src: string
  cover: string
  lrcUrl: string
}

interface LyricLine {
  time: number
  text: string
}

const playlist: Track[] = [
  {
    title: 'からくりピエロ',
    artist: '40mP, 初音ミク',
    src: '/music/karakuri-pierrot.mp3',
    cover: '/music/musicimage/40mp.jpg',
    lrcUrl: '/music/lrc/karakuri-pierrot.lrc',
  },
  {
    title: 'ODDS&ENDS',
    artist: 'ryo (supercell), 初音ミク',
    src: '/music/odds-and-ends.mp3',
    cover: '/music/musicimage/ryo.jpg',
    lrcUrl: '/music/lrc/odds-and-ends.lrc',
  },
]

const expanded = ref(false)
const isPlaying = ref(false)
const currentIndex = ref(0)
const currentTime = ref(0)
const duration = ref(0)
const volume = ref(0.45)
const loopMode = ref(false)
const muted = ref(false)
const lyrics = ref<LyricLine[]>([])
const currentLyricIndex = ref(-1)
const lyricsRef = ref<HTMLElement | null>(null)

let audio: HTMLAudioElement | null = null

const currentTrack = computed(() => playlist[currentIndex.value])

function createAudio() {
  audio = new Audio()
  audio.volume = volume.value
  audio.src = currentTrack.value.src
  audio.addEventListener('loadedmetadata', () => { duration.value = audio!.duration })
  audio.addEventListener('timeupdate', () => {
    currentTime.value = audio!.currentTime
    updateCurrentLyric()
  })
  audio.addEventListener('ended', () => {
    if (loopMode.value) { audio!.currentTime = 0; audio!.play() }
    else { next() }
  })
  audio.addEventListener('play', () => { isPlaying.value = true })
  audio.addEventListener('pause', () => { isPlaying.value = false })
}

function togglePlay() {
  if (!audio) return
  if (isPlaying.value) audio.pause()
  else audio.play().catch(() => {})
}

function prev() { loadTrack((currentIndex.value - 1 + playlist.length) % playlist.length) }
function next() { loadTrack((currentIndex.value + 1) % playlist.length) }

function loadTrack(index: number) {
  if (!audio) return
  const wasPlaying = isPlaying.value
  audio.pause()
  currentIndex.value = index
  currentTime.value = 0
  duration.value = 0
  currentLyricIndex.value = -1
  audio.src = playlist[index].src
  audio.load()
  fetchLyrics(playlist[index].lrcUrl)
  if (wasPlaying) audio.play().catch(() => {})
}

function onSeek(e: Event) {
  if (!audio) return
  audio.currentTime = parseFloat((e.target as HTMLInputElement).value)
}

function seekToLyric(time: number) {
  if (!audio) return
  audio.currentTime = time
  if (!isPlaying.value) audio.play().catch(() => {})
}

function onVolumeInput(e: Event) {
  const v = parseInt((e.target as HTMLInputElement).value) / 100
  volume.value = v
  if (audio) { audio.volume = v; audio.muted = false; muted.value = false }
}

watch(muted, (m) => { if (audio) audio.muted = m })

function parseLRC(text: string): LyricLine[] {
  const result: LyricLine[] = []
  for (const raw of text.split('\n')) {
    const m = raw.match(/\[(\d{2}):(\d{2})\.(\d{2,3})\](.*)/)
    if (!m) continue
    const time = parseInt(m[1]) * 60 + parseInt(m[2]) + parseInt(m[3].padEnd(3, '0')) / 1000
    const t = m[4].trim()
    if (t) result.push({ time, text: t })
  }
  return result.sort((a, b) => a.time - b.time)
}

async function fetchLyrics(url: string) {
  try {
    const res = await fetch(url)
    lyrics.value = res.ok ? parseLRC(await res.text()) : []
  } catch { lyrics.value = [] }
}

function updateCurrentLyric() {
  const t = currentTime.value
  const ls = lyrics.value
  if (!ls.length) return
  let idx = -1
  for (let i = ls.length - 1; i >= 0; i--) {
    if (t >= ls[i].time) { idx = i; break }
  }
  if (idx !== currentLyricIndex.value) {
    currentLyricIndex.value = idx
    nextTick(() => {
      if (!lyricsRef.value || idx < 0) return
      const el = lyricsRef.value.querySelector(`[data-idx="${idx}"]`)
      el?.scrollIntoView({ behavior: 'smooth', block: 'center' })
    })
  }
}

function formatTime(s: number): string {
  if (!s || isNaN(s)) return '0:00'
  return `${Math.floor(s / 60)}:${Math.floor(s % 60).toString().padStart(2, '0')}`
}

onMounted(() => { createAudio(); fetchLyrics(currentTrack.value.lrcUrl) })
onBeforeUnmount(() => { if (audio) { audio.pause(); audio.src = '' } })
</script>

<style scoped>
.player-range {
  -webkit-appearance: none;
  appearance: none;
  background: rgba(255, 255, 255, 0.15);
  border-radius: 9999px;
}
.player-range::-webkit-slider-thumb {
  -webkit-appearance: none;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #39c5bb;
  cursor: pointer;
}
.player-range-lg {
  -webkit-appearance: none;
  appearance: none;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 9999px;
  height: 4px;
}
.player-range-lg::-webkit-slider-thumb {
  -webkit-appearance: none;
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: #39c5bb;
  cursor: pointer;
  box-shadow: 0 0 8px rgba(57, 197, 187, 0.4);
}
.card-slide-enter-active,
.card-slide-leave-active {
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}
.card-slide-enter-from,
.card-slide-leave-to {
  opacity: 0;
  transform: translateY(-8px) scale(0.96);
}
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.25s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
.scrollbar-hide::-webkit-scrollbar {
  display: none;
}
.scrollbar-hide {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
