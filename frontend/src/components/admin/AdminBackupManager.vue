<template>
  <section class="space-y-5">
    <LiquidGlassCard padding="24px">
      <div class="flex flex-wrap items-start justify-between gap-3">
        <div>
          <h1 class="text-2xl font-semibold text-slate-900">{{ copy.title }}</h1>
          <p class="mt-1 text-sm text-slate-600">{{ copy.subtitle }}</p>
        </div>
      </div>
    </LiquidGlassCard>

    <div class="grid gap-4 md:grid-cols-2">
      <LiquidGlassCard padding="20px">
        <h2 class="text-lg font-semibold text-slate-900">{{ copy.jsonCardTitle }}</h2>
        <p class="mt-1 text-sm text-slate-600">{{ copy.jsonCardDescription }}</p>
        <div class="mt-4">
          <MikuButton variant="solid" :disabled="loadingFormat !== null" @click="download('json')">
            {{ loadingFormat === 'json' ? copy.exporting : copy.jsonButton }}
          </MikuButton>
        </div>
      </LiquidGlassCard>

      <LiquidGlassCard padding="20px">
        <h2 class="text-lg font-semibold text-slate-900">{{ copy.sqlCardTitle }}</h2>
        <p class="mt-1 text-sm text-slate-600">{{ copy.sqlCardDescription }}</p>
        <div class="mt-4">
          <MikuButton variant="solid" :disabled="loadingFormat !== null" @click="download('sql')">
            {{ loadingFormat === 'sql' ? copy.exporting : copy.sqlButton }}
          </MikuButton>
        </div>
      </LiquidGlassCard>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref } from 'vue'

import { ApiError } from '../../lib/api'
import { adminCopy } from '../../content/copy'
import { showToast } from '../../stores/ui'
import LiquidGlassCard from '../ui/LiquidGlassCard.vue'
import MikuButton from '../ui/MikuButton.vue'

type BackupFormat = 'json' | 'sql'

const copy = adminCopy.backup
const loadingFormat = ref<BackupFormat | null>(null)

function getToken(): string | null {
  if (typeof window === 'undefined') return null
  try {
    return window.localStorage.getItem('miku_blog_token')
  } catch {
    return null
  }
}

function resolveDownloadFilename(contentDisposition: string | null, fallback: string): string {
  if (!contentDisposition) return fallback
  const match = contentDisposition.match(/filename="?([^";]+)"?/i)
  if (!match || !match[1]) return fallback
  return match[1]
}

async function download(format: BackupFormat) {
  loadingFormat.value = format
  try {
    const headers: Record<string, string> = {}
    const token = getToken()
    if (token) {
      headers.Authorization = `Bearer ${token}`
    }

    const res = await fetch(`/api/v1/admin/backup/export?format=${format}`, {
      method: 'GET',
      headers,
      credentials: 'include',
    })

    if (!res.ok) {
      let message = copy.failed
      try {
        const body = await (res.json() as Promise<{ message?: string }>)
        message = body.message || message
      } catch {
        // ignore json parse error
      }
      throw new ApiError(message, -1, res.status)
    }

    const blob = await res.blob()
    const fallback = `miku-backup.${format}`
    const filename = resolveDownloadFilename(res.headers.get('Content-Disposition'), fallback)

    const url = window.URL.createObjectURL(blob)
    const anchor = document.createElement('a')
    anchor.href = url
    anchor.download = filename
    document.body.appendChild(anchor)
    anchor.click()
    anchor.remove()
    window.URL.revokeObjectURL(url)

    showToast(copy.success, 'success')
  } catch (err) {
    const message = err instanceof ApiError ? err.message : copy.failed
    showToast(message, 'error')
  } finally {
    loadingFormat.value = null
  }
}
</script>
