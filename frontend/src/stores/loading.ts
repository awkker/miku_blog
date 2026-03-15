import { atom, computed, map } from 'nanostores'

// Generic async status registry used by pages/components to render
// idle/loading/success/error states in a consistent way.
export type AsyncStatus = 'idle' | 'loading' | 'success' | 'error'

export const globalLoadingStatus = atom<AsyncStatus>('idle')
export const scopedLoadingStatus = map<Record<string, AsyncStatus>>({})

export function getScopeStatus(scope: string) {
  return computed(scopedLoadingStatus, (all) => all[scope] ?? 'idle')
}

export function getScopeLoading(scope: string) {
  return computed(scopedLoadingStatus, (all) => (all[scope] ?? 'idle') === 'loading')
}

export function setScopeStatus(scope: string, status: AsyncStatus) {
  scopedLoadingStatus.setKey(scope, status)

  if (status === 'loading') {
    globalLoadingStatus.set('loading')
    return
  }

  const hasLoading = Object.values(scopedLoadingStatus.get()).some((value) => value === 'loading')
  globalLoadingStatus.set(hasLoading ? 'loading' : status)
}
