<template>
  <span class="tabular-nums text-[11px] tracking-wide text-[#39c5bb]/90">
    {{ displayText }}
  </span>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { api } from '../../lib/api'

interface WeatherData {
  temp: string
  feels_like: string
  humidity: string
  desc: string
  icon: string
  wind_speed: string
  location: string
}

const iconMap: Record<string, string> = {
  sunny: '\u2600',
  partly_cloudy: '\u26C5',
  cloudy: '\u2601',
  light_rain: '\u2602',
  rain: '\u2602',
  snow: '\u2744',
  thunderstorm: '\u26A1',
}

const displayText = ref('--\u00B0C')

onMounted(async () => {
  try {
    const data = await api.get<WeatherData>('/weather')
    const icon = iconMap[data.icon] || '\u2601'
    displayText.value = `${icon} ${data.temp}\u00B0C`
  } catch {
    displayText.value = '--\u00B0C'
  }
})
</script>
