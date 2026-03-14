import { atom, computed } from 'nanostores'

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

function formatNow(): string {
  return new Intl.DateTimeFormat('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  })
    .format(new Date())
    .replace('/', '-')
    .replace('/', '-')
}

const fallbackMessages: GuestbookMessage[] = [
  {
    id: 'msg_001',
    nickname: 'Airi',
    website: 'https://airi.dev',
    message: '首页的液态玻璃氛围非常舒服，期待之后的文章更新。',
    createdAt: '2026-03-10 20:45',
    avatar: 'https://api.dicebear.com/9.x/shapes/svg?seed=airi',
    votes: 12,
    myVote: 0,
    replies: [
      {
        id: 'msg_001_r1',
        nickname: 'NanaMiku',
        message: '谢谢喜欢！后续会继续优化视觉体验的。',
        createdAt: '2026-03-10 21:30',
        avatar: 'https://api.dicebear.com/9.x/shapes/svg?seed=nanamiku',
        votes: 5,
        myVote: 0,
        replies: [],
        parentId: 'msg_001',
      },
    ],
  },
  {
    id: 'msg_002',
    nickname: 'Konomi',
    message: '留言板的排版很清爽，手机端也很好读。建议加一个暗色模式切换按钮，这样晚上看起来更舒服。',
    createdAt: '2026-03-11 09:12',
    avatar: 'https://api.dicebear.com/9.x/shapes/svg?seed=konomi',
    votes: 8,
    myVote: 0,
    replies: [],
  },
  {
    id: 'msg_003',
    nickname: 'Yuki',
    website: 'https://yuki-note.example',
    message: '路过打卡，祝项目越来越完善。',
    createdAt: '2026-03-12 08:03',
    avatar: 'https://api.dicebear.com/9.x/shapes/svg?seed=yuki',
    votes: 3,
    myVote: 0,
    replies: [
      {
        id: 'msg_003_r1',
        nickname: 'Airi',
        message: '同路过 +1',
        createdAt: '2026-03-12 10:15',
        avatar: 'https://api.dicebear.com/9.x/shapes/svg?seed=airi',
        votes: 2,
        myVote: 0,
        replies: [],
        parentId: 'msg_003',
      },
    ],
  },
]

export const guestbookMessages = atom<GuestbookMessage[]>([])
export const guestbookFetchStatus = atom<'idle' | 'loading' | 'success' | 'error'>('idle')
export const guestbookSubmitStatus = atom<'idle' | 'loading' | 'success' | 'error'>('idle')
export const guestbookError = atom('')
export const guestbookSortMode = atom<SortMode>('hot')

export const guestbookSorted = computed(guestbookMessages, (msgs) => {
  const mode = guestbookSortMode.get()
  const copy = [...msgs]
  if (mode === 'newest') return copy.sort((a, b) => b.createdAt.localeCompare(a.createdAt))
  if (mode === 'oldest') return copy.sort((a, b) => a.createdAt.localeCompare(b.createdAt))
  return copy.sort((a, b) => b.votes - a.votes)
})

function sleep(duration = 750) {
  return new Promise<void>((resolve) => {
    setTimeout(resolve, duration)
  })
}

export async function loadGuestbookMessages(options?: { forceError?: boolean }) {
  guestbookFetchStatus.set('loading')
  guestbookError.set('')

  await sleep(650)

  if (options?.forceError) {
    guestbookFetchStatus.set('error')
    guestbookError.set('留言加载失败，请稍后重试。')
    return
  }

  guestbookMessages.set(fallbackMessages)
  guestbookFetchStatus.set('success')
}

export async function submitGuestbookMessage(draft: GuestbookDraft, options?: { forceError?: boolean }) {
  guestbookSubmitStatus.set('loading')
  guestbookError.set('')

  await sleep(900)

  if (options?.forceError) {
    guestbookSubmitStatus.set('error')
    guestbookError.set('留言提交失败，请检查网络后重试。')
    throw new Error('submit failed')
  }

  const newMessage: GuestbookMessage = {
    id: `msg_${Date.now()}`,
    nickname: draft.nickname.trim(),
    website: draft.website?.trim() || undefined,
    message: draft.message.trim(),
    createdAt: formatNow(),
    avatar: `https://api.dicebear.com/9.x/shapes/svg?seed=${encodeURIComponent(draft.nickname.trim())}`,
    votes: 1,
    myVote: 1,
    replies: [],
    parentId: draft.parentId,
  }

  if (draft.parentId) {
    const msgs = guestbookMessages.get().map((m) => {
      if (m.id === draft.parentId) {
        return { ...m, replies: [...m.replies, newMessage] }
      }
      return m
    })
    guestbookMessages.set(msgs)
  } else {
    guestbookMessages.set([newMessage, ...guestbookMessages.get()])
  }

  guestbookSubmitStatus.set('success')
  return newMessage
}

export function voteMessage(messageId: string, direction: 1 | -1, parentId?: string) {
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
  // Re-trigger computed by touching the atom
  guestbookMessages.set([...guestbookMessages.get()])
}
