import { atom } from 'nanostores'

// Shared UI states for cross-page UX elements such as:
// - admin mobile sidebar drawer
// - global toast notifications
export interface ToastState {
  open: boolean
  message: string
  type: 'success' | 'error' | 'info'
}

export const sidebarOpen = atom(false)
export const toastState = atom<ToastState>({
  open: false,
  message: '',
  type: 'info',
})

let toastTimer: ReturnType<typeof setTimeout> | null = null

export function setSidebarOpen(open: boolean) {
  sidebarOpen.set(open)
}

export function toggleSidebar() {
  sidebarOpen.set(!sidebarOpen.get())
}

export function showToast(message: string, type: ToastState['type'] = 'info', duration = 2200) {
  toastState.set({ open: true, message, type })

  if (toastTimer) {
    clearTimeout(toastTimer)
  }

  toastTimer = setTimeout(() => {
    toastState.set({ open: false, message: '', type: 'info' })
  }, duration)
}

export function closeToast() {
  toastState.set({ open: false, message: '', type: 'info' })

  if (toastTimer) {
    clearTimeout(toastTimer)
    toastTimer = null
  }
}
