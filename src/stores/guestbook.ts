import { atom } from 'nanostores'

// Guestbook domain store:
// message list fetching, submit status, and optimistic insert behavior.
export interface GuestbookMessage {
  id: string
  nickname: string
  website?: string
  message: string
  createdAt: string
  avatar: string
}

export interface GuestbookDraft {
  nickname: string
  website?: string
  message: string
}

const fallbackMessages: GuestbookMessage[] = [
  {
    id: 'msg_001',
    nickname: 'Airi',
    website: 'https://airi.dev',
    message: '首页的液态玻璃氛围非常舒服，期待之后的文章更新。',
    createdAt: '2026-03-10 20:45',
    avatar: 'https://api.dicebear.com/9.x/shapes/svg?seed=airi',
  },
  {
    id: 'msg_002',
    nickname: 'Konomi',
    message: '留言板的排版很清爽，手机端也很好读。',
    createdAt: '2026-03-11 09:12',
    avatar: 'https://api.dicebear.com/9.x/shapes/svg?seed=konomi',
  },
  {
    id: 'msg_003',
    nickname: 'Yuki',
    website: 'https://yuki-note.example',
    message: '路过打卡，祝项目越来越完善。',
    createdAt: '2026-03-12 08:03',
    avatar: 'https://api.dicebear.com/9.x/shapes/svg?seed=yuki',
  },
]

export const guestbookMessages = atom<GuestbookMessage[]>([])
export const guestbookFetchStatus = atom<'idle' | 'loading' | 'success' | 'error'>('idle')
export const guestbookSubmitStatus = atom<'idle' | 'loading' | 'success' | 'error'>('idle')
export const guestbookError = atom('')

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
    createdAt: new Intl.DateTimeFormat('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      hour12: false,
    })
      .format(new Date())
      .replace('/', '-')
      .replace('/', '-'),
    avatar: `https://api.dicebear.com/9.x/shapes/svg?seed=${encodeURIComponent(draft.nickname.trim())}`,
  }

  guestbookMessages.set([newMessage, ...guestbookMessages.get()])
  guestbookSubmitStatus.set('success')

  return newMessage
}
