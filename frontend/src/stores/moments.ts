import { atom } from 'nanostores'

// Moments domain store:
// Twitter/X-style posts with images, likes, reposts, and comments.
export interface MomentComment {
  id: string
  nickname: string
  avatar: string
  content: string
  createdAt: string
  likes: number
  liked: boolean
}

export interface Moment {
  id: string
  nickname: string
  avatar: string
  content: string
  images: readonly string[]
  createdAt: string
  likes: number
  liked: boolean
  reposts: number
  reposted: boolean
  comments: readonly MomentComment[]
}

export interface MomentDraft {
  nickname: string
  content: string
  images: string[]
}

export interface CommentDraft {
  momentId: string
  nickname: string
  content: string
}

const fallbackMoments: Moment[] = [
  {
    id: 'mt_001',
    nickname: 'NanaMiku',
    avatar: 'https://api.dicebear.com/9.x/shapes/svg?seed=nanamiku',
    content: '博客的液态玻璃风格终于调到满意的效果了。花了整整一个下午在 backdrop-blur 和 border 透明度之间反复横跳，最后发现秘诀是 saturate 拉高一点点。分享几张截图给大家看看 ~',
    images: [
      'https://images.unsplash.com/photo-1555066931-4365d14bab8c?w=600&h=400&fit=crop',
      'https://images.unsplash.com/photo-1461749280684-dccba630e2f6?w=600&h=400&fit=crop',
    ],
    createdAt: '2026-03-14 10:30',
    likes: 24,
    liked: false,
    reposts: 3,
    reposted: false,
    comments: [
      {
        id: 'mc_001_1',
        nickname: 'Airi',
        avatar: 'https://api.dicebear.com/9.x/shapes/svg?seed=airi',
        content: '效果真的很棒，尤其是那个毛玻璃渐变，太好看了！',
        createdAt: '2026-03-14 11:05',
        likes: 6,
        liked: false,
      },
      {
        id: 'mc_001_2',
        nickname: 'Konomi',
        avatar: 'https://api.dicebear.com/9.x/shapes/svg?seed=konomi',
        content: '请问 saturate 值大概调到多少？我也想试试这个效果。',
        createdAt: '2026-03-14 11:42',
        likes: 2,
        liked: false,
      },
    ],
  },
  {
    id: 'mt_002',
    nickname: 'NanaMiku',
    avatar: 'https://api.dicebear.com/9.x/shapes/svg?seed=nanamiku',
    content: '今天把留言板改成了 Reddit 风格的嵌套评论，支持投票和回复。感觉交互体验好了很多，比之前的平铺卡片有层次感多了。\n\n下一步打算加个说说功能，类似推特那种。大家有什么想法可以留言告诉我！',
    images: [],
    createdAt: '2026-03-13 18:20',
    likes: 15,
    liked: false,
    reposts: 1,
    reposted: false,
    comments: [
      {
        id: 'mc_002_1',
        nickname: 'Yuki',
        avatar: 'https://api.dicebear.com/9.x/shapes/svg?seed=yuki',
        content: '说说功能期待了！建议支持图片上传。',
        createdAt: '2026-03-13 19:00',
        likes: 4,
        liked: false,
      },
    ],
  },
  {
    id: 'mt_003',
    nickname: 'NanaMiku',
    avatar: 'https://api.dicebear.com/9.x/shapes/svg?seed=nanamiku',
    content: '周末读完了《深入理解 TypeScript》，强烈推荐给所有前端开发者。类型体操虽然烧脑，但写出来的代码确实更健壮。',
    images: [
      'https://images.unsplash.com/photo-1532012197267-da84d127e765?w=600&h=400&fit=crop',
    ],
    createdAt: '2026-03-12 14:50',
    likes: 31,
    liked: false,
    reposts: 7,
    reposted: false,
    comments: [],
  },
]

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

function sleep(duration = 750) {
  return new Promise<void>((resolve) => {
    setTimeout(resolve, duration)
  })
}

export const moments = atom<Moment[]>([])
export const momentsFetchStatus = atom<'idle' | 'loading' | 'success' | 'error'>('idle')
export const momentsSubmitStatus = atom<'idle' | 'loading' | 'success' | 'error'>('idle')
export const momentsError = atom('')

export async function loadMoments(options?: { forceError?: boolean }) {
  momentsFetchStatus.set('loading')
  momentsError.set('')

  await sleep(650)

  if (options?.forceError) {
    momentsFetchStatus.set('error')
    momentsError.set('说说加载失败，请稍后重试。')
    return
  }

  moments.set(fallbackMoments)
  momentsFetchStatus.set('success')
}

export async function submitMoment(draft: MomentDraft, options?: { forceError?: boolean }) {
  momentsSubmitStatus.set('loading')
  momentsError.set('')

  await sleep(900)

  if (options?.forceError) {
    momentsSubmitStatus.set('error')
    momentsError.set('发布失败，请检查网络后重试。')
    throw new Error('submit failed')
  }

  const newMoment: Moment = {
    id: `mt_${Date.now()}`,
    nickname: draft.nickname.trim(),
    avatar: `https://api.dicebear.com/9.x/shapes/svg?seed=${encodeURIComponent(draft.nickname.trim())}`,
    content: draft.content.trim(),
    images: draft.images,
    createdAt: formatNow(),
    likes: 0,
    liked: false,
    reposts: 0,
    reposted: false,
    comments: [],
  }

  moments.set([newMoment, ...moments.get()])
  momentsSubmitStatus.set('success')
  return newMoment
}

export function toggleLikeMoment(momentId: string) {
  const list = moments.get().map((m) => {
    if (m.id === momentId) {
      return m.liked
        ? { ...m, likes: m.likes - 1, liked: false }
        : { ...m, likes: m.likes + 1, liked: true }
    }
    return m
  })
  moments.set(list)
}

export function toggleRepostMoment(momentId: string) {
  const list = moments.get().map((m) => {
    if (m.id === momentId) {
      return m.reposted
        ? { ...m, reposts: m.reposts - 1, reposted: false }
        : { ...m, reposts: m.reposts + 1, reposted: true }
    }
    return m
  })
  moments.set(list)
}

export function addComment(draft: CommentDraft) {
  const comment: MomentComment = {
    id: `mc_${Date.now()}`,
    nickname: draft.nickname.trim(),
    avatar: `https://api.dicebear.com/9.x/shapes/svg?seed=${encodeURIComponent(draft.nickname.trim())}`,
    content: draft.content.trim(),
    createdAt: formatNow(),
    likes: 0,
    liked: false,
  }

  const list = moments.get().map((m) => {
    if (m.id === draft.momentId) {
      return { ...m, comments: [...m.comments, comment] }
    }
    return m
  })
  moments.set(list)
}

export function toggleLikeComment(momentId: string, commentId: string) {
  const list = moments.get().map((m) => {
    if (m.id === momentId) {
      return {
        ...m,
        comments: m.comments.map((c) =>
          c.id === commentId
            ? c.liked
              ? { ...c, likes: c.likes - 1, liked: false }
              : { ...c, likes: c.likes + 1, liked: true }
            : c,
        ),
      }
    }
    return m
  })
  moments.set(list)
}
