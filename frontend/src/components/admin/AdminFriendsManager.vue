<template>
  <section class="space-y-5">
    <LiquidGlassCard padding="24px">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 class="text-2xl font-semibold text-slate-900">友链管理</h1>
          <p class="mt-1 text-sm text-slate-600">维护友链展示、申请审核和站点健康检查。</p>
        </div>
        <MikuButton variant="solid" aria-label="添加友链" @click="toggleCreateForm">+ 添加友链</MikuButton>
      </div>
    </LiquidGlassCard>

    <!-- Create Form -->
    <LiquidGlassCard v-if="showCreateForm" padding="24px">
      <h2 class="mb-4 text-lg font-semibold text-slate-900">添加友链</h2>
      <form class="space-y-3" @submit.prevent="createFriend">
        <div class="grid gap-3 md:grid-cols-2">
          <input v-model="newFriend.name" type="text" placeholder="站点名称 *" class="rounded-xl border border-slate-200/80 bg-white/60 px-3 py-2 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-1 focus:ring-miku/30" />
          <input v-model="newFriend.url" type="text" placeholder="站点 URL *" class="rounded-xl border border-slate-200/80 bg-white/60 px-3 py-2 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-1 focus:ring-miku/30" />
        </div>
        <input v-model="newFriend.avatar_url" type="text" placeholder="头像 URL" class="w-full rounded-xl border border-slate-200/80 bg-white/60 px-3 py-2 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-1 focus:ring-miku/30" />
        <input v-model="newFriend.description" type="text" placeholder="站点描述" class="w-full rounded-xl border border-slate-200/80 bg-white/60 px-3 py-2 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-1 focus:ring-miku/30" />
        <div class="flex items-center gap-3">
          <MikuButton type="submit" variant="solid" :disabled="creatingFriend">{{ creatingFriend ? '添加中...' : '添加友链' }}</MikuButton>
          <button type="button" class="text-sm text-slate-500 hover:text-slate-700" @click="closeCreateForm">取消</button>
        </div>
      </form>
    </LiquidGlassCard>

    <LiquidGlassCard v-if="showEditForm" padding="24px">
      <h2 class="mb-4 text-lg font-semibold text-slate-900">编辑友链</h2>
      <form class="space-y-3" @submit.prevent="updateFriend">
        <div class="grid gap-3 md:grid-cols-2">
          <input v-model="editFriend.name" type="text" placeholder="站点名称 *" class="rounded-xl border border-slate-200/80 bg-white/60 px-3 py-2 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-1 focus:ring-miku/30" />
          <input v-model="editFriend.url" type="text" placeholder="站点 URL *" class="rounded-xl border border-slate-200/80 bg-white/60 px-3 py-2 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-1 focus:ring-miku/30" />
        </div>
        <input v-model="editFriend.avatar_url" type="text" placeholder="头像 URL" class="w-full rounded-xl border border-slate-200/80 bg-white/60 px-3 py-2 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-1 focus:ring-miku/30" />
        <input v-model="editFriend.description" type="text" placeholder="站点描述" class="w-full rounded-xl border border-slate-200/80 bg-white/60 px-3 py-2 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-1 focus:ring-miku/30" />
        <input v-model.number="editFriend.sort_order" type="number" min="0" placeholder="排序值 (越小越靠前)" class="w-full rounded-xl border border-slate-200/80 bg-white/60 px-3 py-2 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-1 focus:ring-miku/30" />
        <div class="flex items-center gap-3">
          <MikuButton type="submit" variant="solid" :disabled="updatingFriend">{{ updatingFriend ? '保存中...' : '保存修改' }}</MikuButton>
          <button type="button" class="text-sm text-slate-500 hover:text-slate-700" @click="closeEditForm">取消</button>
        </div>
      </form>
    </LiquidGlassCard>

    <div class="grid gap-4 sm:grid-cols-3">
      <LiquidGlassCard padding="16px">
        <p class="text-xs uppercase tracking-[0.2em] text-slate-600">已通过</p>
        <p class="mt-1 font-mono text-2xl font-semibold text-slate-900">{{ approvedCount }}</p>
      </LiquidGlassCard>
      <LiquidGlassCard padding="16px">
        <p class="text-xs uppercase tracking-[0.2em] text-slate-600">待审核</p>
        <p class="mt-1 font-mono text-2xl font-semibold text-amber-600">{{ pendingCount }}</p>
      </LiquidGlassCard>
      <LiquidGlassCard padding="16px">
        <p class="text-xs uppercase tracking-[0.2em] text-slate-600">异常 / 不可达</p>
        <p class="mt-1 font-mono text-2xl font-semibold text-red-600">{{ downCount }}</p>
      </LiquidGlassCard>
    </div>

    <LiquidGlassCard padding="0px">
      <div class="overflow-x-auto">
        <table class="w-full text-left text-sm">
          <thead>
            <tr class="border-b border-slate-200/60">
              <th class="px-5 py-3.5 font-semibold text-slate-700">站点名称</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700">链接</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700">状态</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700">健康检测</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700">添加时间</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700 text-center">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="link in friends"
              :key="link.id"
              class="border-b border-slate-100/60 transition hover:bg-white/40"
            >
              <td class="px-5 py-3.5 font-medium text-slate-900">{{ link.name }}</td>
              <td class="px-5 py-3.5">
                <a
                  :href="link.url"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="text-miku underline decoration-miku/30 transition hover:decoration-miku"
                >
                  {{ link.url }}
                </a>
              </td>
              <td class="px-5 py-3.5">
                <span
                  class="inline-block rounded-full px-2.5 py-0.5 text-xs font-medium"
                  :class="statusClass(link.status)"
                >
                  {{ statusLabel(link.status) }}
                </span>
              </td>
              <td class="px-5 py-3.5">
                <span
                  class="inline-block rounded-full px-2.5 py-0.5 text-xs font-medium"
                  :class="link.health === 'ok' ? 'bg-emerald-100 text-emerald-700' : 'bg-red-100 text-red-600'"
                >
                  {{ link.health === 'ok' ? '正常' : '不可达' }}
                </span>
              </td>
              <td class="px-5 py-3.5 text-slate-600">{{ link.createdAt }}</td>
              <td class="px-5 py-3.5 text-center">
                <div class="flex items-center justify-center gap-2">
                  <button
                    v-if="link.status === 'pending'"
                    type="button"
                    class="rounded-xl border border-emerald-200/80 bg-white/50 px-2.5 py-1 text-xs text-emerald-600 transition hover:bg-emerald-50"
                    aria-label="通过申请"
                  >
                    通过
                  </button>
                  <button
                    type="button"
                    class="rounded-xl border border-slate-200/80 bg-white/50 px-2.5 py-1 text-xs text-slate-700 transition hover:border-miku/40 hover:text-miku"
                    aria-label="编辑友链"
                    @click="startEditFriend(link)"
                  >
                    编辑
                  </button>
                  <button
                    type="button"
                    class="rounded-xl border border-red-200/80 bg-white/50 px-2.5 py-1 text-xs text-red-600 transition hover:bg-red-50"
                    aria-label="删除友链"
                    @click="deleteFriend(link.id)"
                  >
                    删除
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </LiquidGlassCard>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

import { api, ApiError, type PagedData } from '../../lib/api'
import { showToast } from '../../stores/ui'
import LiquidGlassCard from '../ui/LiquidGlassCard.vue'
import MikuButton from '../ui/MikuButton.vue'

interface ApiFriendLink {
  id: string
  name: string
  description: string
  url: string
  domain: string
  avatar_url: string
  status: string
  health_status: string
  sort_order: number
  created_at: string
}

interface FriendLink {
  id: string
  name: string
  description: string
  url: string
  domain: string
  avatar_url: string
  sortOrder: number
  status: 'approved' | 'pending' | 'rejected'
  health: 'ok' | 'down'
  createdAt: string
}

interface FriendForm {
  name: string
  url: string
  avatar_url: string
  description: string
  sort_order: number
}

function mapStatus(s: string): 'approved' | 'pending' | 'rejected' {
  if (s === 'approved') return 'approved'
  if (s === 'rejected') return 'rejected'
  return 'pending'
}

function formatDate(iso: string): string {
  try {
    return iso.slice(0, 10)
  } catch {
    return iso
  }
}

function mapFriend(item: ApiFriendLink): FriendLink {
  return {
    id: item.id,
    name: item.name,
    description: item.description || '',
    url: item.url,
    domain: item.domain || '',
    avatar_url: item.avatar_url || '',
    sortOrder: Number(item.sort_order) || 0,
    status: mapStatus(item.status),
    health: item.health_status === 'ok' ? 'ok' : 'down',
    createdAt: formatDate(item.created_at),
  }
}

function createEmptyFriendForm(): FriendForm {
  return {
    name: '',
    url: '',
    avatar_url: '',
    description: '',
    sort_order: 0,
  }
}

function deriveDomainFromURL(rawURL: string): string {
  const input = rawURL.trim()
  if (!input) return ''
  try {
    const withScheme = /^[a-zA-Z][a-zA-Z\d+\-.]*:\/\//.test(input) ? input : `https://${input}`
    return new URL(withScheme).hostname.replace(/^www\./, '')
  } catch {
    return ''
  }
}

function toFriendPayload(form: FriendForm) {
  const url = form.url.trim()
  return {
    name: form.name.trim(),
    url,
    domain: deriveDomainFromURL(url),
    avatar_url: form.avatar_url.trim(),
    description: form.description.trim(),
    sort_order: Number(form.sort_order) || 0,
  }
}

const friends = ref<FriendLink[]>([])
const loading = ref(false)
const showCreateForm = ref(false)
const creatingFriend = ref(false)
const showEditForm = ref(false)
const updatingFriend = ref(false)
const editingFriendID = ref<string | null>(null)

const newFriend = ref<FriendForm>(createEmptyFriendForm())
const editFriend = ref<FriendForm>(createEmptyFriendForm())

function toggleCreateForm() {
  showCreateForm.value = !showCreateForm.value
  if (showCreateForm.value) {
    closeEditForm()
  }
}

function closeCreateForm() {
  showCreateForm.value = false
}

function closeEditForm() {
  showEditForm.value = false
  editingFriendID.value = null
  editFriend.value = createEmptyFriendForm()
}

async function loadFriends() {
  loading.value = true
  try {
    const data = await api.get<PagedData<ApiFriendLink>>('/admin/friends?size=100')
    friends.value = (data.items || []).map(mapFriend)
  } catch (err) {
    console.error('[AdminFriends] loadFriends failed:', err)
    showToast('加载友链列表失败', 'error')
    friends.value = []
  } finally {
    loading.value = false
  }
}

async function createFriend() {
  if (!newFriend.value.name.trim() || !newFriend.value.url.trim()) return
  creatingFriend.value = true
  try {
    await api.post('/admin/friends', toFriendPayload(newFriend.value))
    closeCreateForm()
    newFriend.value = createEmptyFriendForm()
    showToast('友链添加成功', 'success')
    await loadFriends()
  } catch (err) {
    const msg = err instanceof ApiError ? err.message : '添加友链失败'
    console.error('[AdminFriends] createFriend failed:', err)
    showToast(msg, 'error')
  } finally {
    creatingFriend.value = false
  }
}

function startEditFriend(link: FriendLink) {
  editingFriendID.value = link.id
  editFriend.value = {
    name: link.name,
    url: link.url,
    avatar_url: link.avatar_url,
    description: link.description,
    sort_order: link.sortOrder,
  }
  showEditForm.value = true
  showCreateForm.value = false
}

async function updateFriend() {
  if (!editingFriendID.value) return
  if (!editFriend.value.name.trim() || !editFriend.value.url.trim()) return
  updatingFriend.value = true
  try {
    await api.put(`/admin/friends/${editingFriendID.value}`, toFriendPayload(editFriend.value))
    showToast('友链更新成功', 'success')
    closeEditForm()
    await loadFriends()
  } catch (err) {
    const msg = err instanceof ApiError ? err.message : '更新友链失败'
    console.error('[AdminFriends] updateFriend failed:', err)
    showToast(msg, 'error')
  } finally {
    updatingFriend.value = false
  }
}

async function deleteFriend(id: string) {
  try {
    await api.delete(`/admin/friends/${id}`)
    friends.value = friends.value.filter((f) => f.id !== id)
    showToast('友链已删除', 'success')
  } catch (err) {
    const msg = err instanceof ApiError ? err.message : '删除失败'
    console.error('[AdminFriends] deleteFriend failed:', err)
    showToast(msg, 'error')
  }
}

onMounted(() => {
  loadFriends()
})

const approvedCount = computed(() => friends.value.filter((f) => f.status === 'approved').length)
const pendingCount = computed(() => friends.value.filter((f) => f.status === 'pending').length)
const downCount = computed(() => friends.value.filter((f) => f.health === 'down').length)

function statusClass(status: FriendLink['status']) {
  if (status === 'approved') return 'bg-emerald-100 text-emerald-700'
  if (status === 'rejected') return 'bg-red-100 text-red-600'
  return 'bg-amber-100 text-amber-700'
}

function statusLabel(status: FriendLink['status']) {
  if (status === 'approved') return '已通过'
  if (status === 'rejected') return '已拒绝'
  return '待审核'
}
</script>
