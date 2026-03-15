import { atom } from 'nanostores'

// Friends-link domain store:
// encapsulates loading state + list data for the friends page.
export interface FriendLink {
  id: string
  name: string
  description: string
  url: string
  domain: string
  avatar: string
}

const fallbackFriends: FriendLink[] = [
  {
    id: 'friend_001',
    name: 'Aurora Notes',
    description: '专注前端工程化与交互细节的技术笔记站。',
    url: 'https://aurora-notes.example',
    domain: 'aurora-notes.example',
    avatar: 'https://api.dicebear.com/9.x/glass/svg?seed=aurora',
  },
  {
    id: 'friend_002',
    name: 'Nebula Journal',
    description: '记录设计系统、产品体验与独立开发实践。',
    url: 'https://nebula-journal.example',
    domain: 'nebula-journal.example',
    avatar: 'https://api.dicebear.com/9.x/glass/svg?seed=nebula',
  },
  {
    id: 'friend_003',
    name: 'Code Rain',
    description: '后端性能调优、Go 工程经验与数据库优化。',
    url: 'https://code-rain.example',
    domain: 'code-rain.example',
    avatar: 'https://api.dicebear.com/9.x/glass/svg?seed=code-rain',
  },
  {
    id: 'friend_004',
    name: 'Morning Tea Lab',
    description: '以生活化视角分享创作、效率与写作工作流。',
    url: 'https://morning-tea.example',
    domain: 'morning-tea.example',
    avatar: 'https://api.dicebear.com/9.x/glass/svg?seed=morning-tea',
  },
]

export const friendLinks = atom<FriendLink[]>([])
export const friendFetchStatus = atom<'idle' | 'loading' | 'success' | 'error'>('idle')
export const friendError = atom('')

function sleep(duration = 720) {
  return new Promise<void>((resolve) => {
    setTimeout(resolve, duration)
  })
}

export async function loadFriendLinks(options?: { forceError?: boolean }) {
  friendFetchStatus.set('loading')
  friendError.set('')

  await sleep()

  if (options?.forceError) {
    friendFetchStatus.set('error')
    friendError.set('友链加载失败，请稍后再试。')
    return
  }

  friendLinks.set(fallbackFriends)
  friendFetchStatus.set('success')
}
