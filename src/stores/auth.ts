import { atom, computed } from 'nanostores'

// Centralized authentication state:
// - hydrate from localStorage
// - mock login/logout flows
// - shared auth status for route guards and admin islands
const TOKEN_KEY = 'miku_blog_token'
const USER_KEY = 'miku_blog_user'

export interface AuthUser {
  id: string
  name: string
  email: string
  role: 'admin'
}

export interface AuthState {
  status: 'checking' | 'guest' | 'authenticated'
  token: string | null
  user: AuthUser | null
}

export const authState = atom<AuthState>({
  status: 'checking',
  token: null,
  user: null,
})

export const isAuthenticated = computed(authState, (state) => state.status === 'authenticated')

function isBrowser() {
  return typeof window !== 'undefined'
}

function persistAuth(token: string, user: AuthUser) {
  if (!isBrowser()) {
    return
  }

  try {
    window.localStorage.setItem(TOKEN_KEY, token)
    window.localStorage.setItem(USER_KEY, JSON.stringify(user))
  } catch {
    throw new Error('浏览器阻止了本地存储，无法保持登录状态。请关闭无痕模式后重试。')
  }
}

function clearPersistedAuth() {
  if (!isBrowser()) {
    return
  }

  window.localStorage.removeItem(TOKEN_KEY)
  window.localStorage.removeItem(USER_KEY)
}

export function hydrateAuth() {
  if (!isBrowser()) {
    return
  }

  let token: string | null = null
  let rawUser: string | null = null

  try {
    token = window.localStorage.getItem(TOKEN_KEY)
    rawUser = window.localStorage.getItem(USER_KEY)
  } catch {
    authState.set({ status: 'guest', token: null, user: null })
    return
  }

  if (!token || !rawUser) {
    authState.set({ status: 'guest', token: null, user: null })
    return
  }

  try {
    const user = JSON.parse(rawUser) as AuthUser
    authState.set({ status: 'authenticated', token, user })
  } catch {
    clearPersistedAuth()
    authState.set({ status: 'guest', token: null, user: null })
  }
}

function sleep(duration = 900) {
  return new Promise<void>((resolve) => {
    setTimeout(resolve, duration)
  })
}

export async function loginWithPassword(identifier: string, password: string) {
  await sleep()

  const normalized = identifier.trim().toLowerCase()
  const normalizedPassword = password.trim()
  const validIdentifier = normalized === 'admin'
  const validPassword = normalizedPassword === 'miku1234'

  if (!validIdentifier || !validPassword) {
    throw new Error('账号或密码错误，请检查后重试。')
  }

  const token = `miku_token_${Date.now()}`
  const user: AuthUser = {
    id: 'admin-001',
    name: 'Nanamiku Admin',
    email: 'admin@miku.blog',
    role: 'admin',
  }

  persistAuth(token, user)
  authState.set({ status: 'authenticated', token, user })

  return user
}

export function logout() {
  clearPersistedAuth()
  authState.set({ status: 'guest', token: null, user: null })
}

export function getStoredToken() {
  if (!isBrowser()) {
    return null
  }

  return window.localStorage.getItem(TOKEN_KEY)
}
