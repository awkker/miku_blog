import { atom } from 'nanostores'

import { api, type PagedData } from '../lib/api'
import { siteCopy } from '../content/copy'

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

interface ApiMoment {
  id: string
  author_name: string
  content: string
  image_urls: string[]
  like_count: number
  repost_count: number
  comment_count: number
  created_at: string
  liked: boolean
  reposted: boolean
}

interface ApiMomentComment {
  id: string
  author_name: string
  content: string
  like_count: number
  created_at: string
  liked: boolean
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

function mapComment(c: ApiMomentComment): MomentComment {
  return {
    id: c.id,
    nickname: c.author_name,
    avatar: `https://api.dicebear.com/9.x/shapes/svg?seed=${encodeURIComponent(c.author_name)}`,
    content: c.content,
    createdAt: formatDate(c.created_at),
    likes: Number(c.like_count) || 0,
    liked: c.liked,
  }
}

function mapMoment(item: ApiMoment, comments: MomentComment[] = []): Moment {
  return {
    id: item.id,
    nickname: item.author_name,
    avatar: `https://api.dicebear.com/9.x/shapes/svg?seed=${encodeURIComponent(item.author_name)}`,
    content: item.content,
    images: item.image_urls || [],
    createdAt: formatDate(item.created_at),
    likes: Number(item.like_count) || 0,
    liked: item.liked,
    reposts: Number(item.repost_count) || 0,
    reposted: item.reposted,
    comments,
  }
}

async function fetchComments(momentId: string): Promise<MomentComment[]> {
  try {
    const data = await api.get<PagedData<ApiMomentComment>>(`/moments/${momentId}/comments?size=50`)
    return (data.items || []).map(mapComment)
  } catch {
    return []
  }
}

export const moments = atom<Moment[]>([])
export const momentsFetchStatus = atom<'idle' | 'loading' | 'success' | 'error'>('idle')
export const momentsSubmitStatus = atom<'idle' | 'loading' | 'success' | 'error'>('idle')
export const momentsError = atom('')
const copy = siteCopy.components.momentsBoard

export async function loadMoments() {
  momentsFetchStatus.set('loading')
  momentsError.set('')

  try {
    const data = await api.get<PagedData<ApiMoment>>('/moments?size=50')
    const apiItems = data.items || []

    // Fetch comments for each moment in parallel
    const commentLists = await Promise.all(
      apiItems.map((item) =>
        item.comment_count > 0 ? fetchComments(item.id) : Promise.resolve([]),
      ),
    )

    const mapped = apiItems.map((item, i) => mapMoment(item, commentLists[i]))
    moments.set(mapped)
    momentsFetchStatus.set('success')
  } catch {
    momentsFetchStatus.set('error')
    momentsError.set(copy.loadErrorFallback)
  }
}

export async function submitMoment(draft: MomentDraft) {
  momentsSubmitStatus.set('loading')
  momentsError.set('')

  try {
    const created = await api.post<ApiMoment>('/moments', {
      author_name: draft.nickname.trim(),
      content: draft.content.trim(),
      image_urls: draft.images.filter((u) => u.trim()),
    })

    const newMoment = mapMoment(created)
    moments.set([newMoment, ...moments.get()])
    momentsSubmitStatus.set('success')
    return newMoment
  } catch {
    momentsSubmitStatus.set('error')
    momentsError.set(copy.submitError)
    throw new Error('submit failed')
  }
}

export async function toggleLikeMoment(momentId: string) {
  // Optimistic update
  const list = moments.get().map((m) => {
    if (m.id === momentId) {
      return m.liked
        ? { ...m, likes: m.likes - 1, liked: false }
        : { ...m, likes: m.likes + 1, liked: true }
    }
    return m
  })
  moments.set(list)

  try {
    await api.post(`/moments/${momentId}/like`)
  } catch {
    // Revert on failure
    await loadMoments()
  }
}

export async function toggleRepostMoment(momentId: string) {
  // Optimistic update
  const list = moments.get().map((m) => {
    if (m.id === momentId) {
      return m.reposted
        ? { ...m, reposts: m.reposts - 1, reposted: false }
        : { ...m, reposts: m.reposts + 1, reposted: true }
    }
    return m
  })
  moments.set(list)

  try {
    await api.post(`/moments/${momentId}/repost`)
  } catch {
    await loadMoments()
  }
}

export async function addComment(draft: CommentDraft) {
  try {
    const created = await api.post<ApiMomentComment>(`/moments/${draft.momentId}/comments`, {
      author_name: draft.nickname.trim(),
      content: draft.content.trim(),
    })

    const comment = mapComment(created)
    const list = moments.get().map((m) => {
      if (m.id === draft.momentId) {
        return { ...m, comments: [...m.comments, comment] }
      }
      return m
    })
    moments.set(list)
  } catch {
    throw new Error('comment failed')
  }
}

export async function toggleLikeComment(momentId: string, commentId: string) {
  // Optimistic update
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

  try {
    await api.post(`/moments/comments/${commentId}/like`)
  } catch {
    await loadMoments()
  }
}
