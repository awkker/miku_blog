import { atom } from 'nanostores'

import { api } from '../lib/api'

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

interface ApiFriendLink {
  id: string
  name: string
  description: string
  url: string
  avatar_url: string
  health_status: string
}

function domainFromUrl(url: string): string {
  try {
    return new URL(url).hostname
  } catch {
    return url
  }
}

function mapFriend(item: ApiFriendLink): FriendLink {
  return {
    id: item.id,
    name: item.name,
    description: item.description,
    url: item.url,
    domain: domainFromUrl(item.url),
    avatar: item.avatar_url || `https://api.dicebear.com/9.x/glass/svg?seed=${encodeURIComponent(item.name)}`,
  }
}

export const friendLinks = atom<FriendLink[]>([])
export const friendFetchStatus = atom<'idle' | 'loading' | 'success' | 'error'>('idle')
export const friendError = atom('')

export async function loadFriendLinks() {
  friendFetchStatus.set('loading')
  friendError.set('')

  try {
    const items = await api.get<ApiFriendLink[]>('/friends')
    friendLinks.set((items || []).map(mapFriend))
    friendFetchStatus.set('success')
  } catch {
    friendFetchStatus.set('error')
    friendError.set('友链加载失败，请稍后再试。')
  }
}
