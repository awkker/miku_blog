import { atom, computed } from 'nanostores'

import { api, type PagedData } from '../lib/api'
import { siteCopy } from '../content/copy'

// Guestbook domain store:
// Reddit-style threaded messages with votes, replies, and sort.
export interface GuestbookMessage {
  id: string
  nickname: string
  website?: string
  message: string
  createdAt: string
  avatar: string
  votes: number
  myVote: -1 | 0 | 1
  replies: readonly GuestbookMessage[]
  parentId?: string
}

export interface GuestbookDraft {
  nickname: string
  website?: string
  message: string
  parentId?: string
}

export type SortMode = 'newest' | 'oldest' | 'hot'

interface ApiGuestbookMessage {
  id: string
  parent_id?: string
  author_name: string
  author_website?: string
  content: string
  vote_score: number
  created_at: string
  replies?: ApiGuestbookMessage[]
  my_vote?: string | null
}

function formatDate(iso: string): string {
  try {
    const d = new Date(iso)
    return new Intl.DateTimeFormat('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      hour12: false,
    })
      .format(d)
      .replace('/', '-')
      .replace('/', '-')
  } catch {
    return iso
  }
}

function mapVote(v?: string | null): -1 | 0 | 1 {
  if (v === 'up') return 1
  if (v === 'down') return -1
  return 0
}

function mapMessage(item: ApiGuestbookMessage): GuestbookMessage {
  return {
    id: item.id,
    nickname: item.author_name,
    website: item.author_website || undefined,
    message: item.content,
    createdAt: formatDate(item.created_at),
    avatar: `https://api.dicebear.com/9.x/shapes/svg?seed=${encodeURIComponent(item.author_name)}`,
    votes: item.vote_score,
    myVote: mapVote(item.my_vote),
    replies: (item.replies || []).map(mapMessage),
    parentId: item.parent_id || undefined,
  }
}

function sortModeToApi(mode: SortMode): string {
  if (mode === 'hot') return 'hot'
  if (mode === 'oldest') return 'oldest'
  return 'newest'
}

export const guestbookMessages = atom<GuestbookMessage[]>([])
export const guestbookFetchStatus = atom<'idle' | 'loading' | 'success' | 'error'>('idle')
export const guestbookSubmitStatus = atom<'idle' | 'loading' | 'success' | 'error'>('idle')
export const guestbookError = atom('')
export const guestbookSortMode = atom<SortMode>('hot')
const copy = siteCopy.components.guestbookBoard

export const guestbookSorted = computed(guestbookMessages, (msgs) => {
  const mode = guestbookSortMode.get()
  const copy = [...msgs]
  if (mode === 'newest') return copy.sort((a, b) => b.createdAt.localeCompare(a.createdAt))
  if (mode === 'oldest') return copy.sort((a, b) => a.createdAt.localeCompare(b.createdAt))
  return copy.sort((a, b) => b.votes - a.votes)
})

export async function loadGuestbookMessages() {
  guestbookFetchStatus.set('loading')
  guestbookError.set('')

  try {
    const sort = sortModeToApi(guestbookSortMode.get())
    const data = await api.get<PagedData<ApiGuestbookMessage>>(`/guestbook/messages?sort=${sort}&size=50`)
    guestbookMessages.set((data.items || []).map(mapMessage))
    guestbookFetchStatus.set('success')
  } catch {
    guestbookFetchStatus.set('error')
    guestbookError.set(copy.loadErrorFallback)
  }
}

export async function submitGuestbookMessage(draft: GuestbookDraft) {
  guestbookSubmitStatus.set('loading')
  guestbookError.set('')

  try {
    const body: Record<string, unknown> = {
      author_name: draft.nickname.trim(),
      author_website: draft.website?.trim() || '',
      content: draft.message.trim(),
    }
    if (draft.parentId) {
      body.parent_id = draft.parentId
    }

    const created = await api.post<ApiGuestbookMessage>('/guestbook/messages', body)
    const newMessage = mapMessage(created)

    // Messages now require moderation; refresh approved list instead of optimistic insert.
    await loadGuestbookMessages()
    guestbookSubmitStatus.set('success')
    return newMessage
  } catch {
    guestbookSubmitStatus.set('error')
    guestbookError.set(copy.submitError)
    throw new Error('submit failed')
  }
}

export async function voteMessage(messageId: string, direction: 1 | -1, parentId?: string) {
  const vote = direction === 1 ? 'up' : 'down'

  // Optimistic update
  const msgs = guestbookMessages.get().map((m) => {
    if (parentId && m.id === parentId) {
      return {
        ...m,
        replies: m.replies.map((r) =>
          r.id === messageId ? applyVote(r, direction) : r,
        ),
      }
    }
    if (m.id === messageId) return applyVote(m, direction)
    return m
  })
  guestbookMessages.set(msgs)

  try {
    await api.post(`/guestbook/messages/${messageId}/vote`, { vote })
  } catch {
    // Revert on failure - reload from server
    await loadGuestbookMessages()
  }
}

function applyVote(msg: GuestbookMessage, direction: 1 | -1): GuestbookMessage {
  if (msg.myVote === direction) {
    return { ...msg, votes: msg.votes - direction, myVote: 0 }
  }
  const delta = direction - msg.myVote
  return { ...msg, votes: msg.votes + delta, myVote: direction }
}

export function setSortMode(mode: SortMode) {
  guestbookSortMode.set(mode)
  loadGuestbookMessages()
}
