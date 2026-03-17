const API_BASE = '/api/v1'

export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data?: T
}

export interface PagedData<T = unknown> {
  items: T[]
  total: number
  page: number
  size: number
}

function getToken(): string | null {
  if (typeof window === 'undefined') return null
  try {
    return window.localStorage.getItem('miku_blog_token')
  } catch {
    return null
  }
}

async function request<T>(
  path: string,
  options: RequestInit = {},
): Promise<T> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string>),
  }

  const token = getToken()
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  const res = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers,
    credentials: 'include',
  })

  const contentType = res.headers.get('content-type') || ''
  let body: ApiResponse<T> | null = null
  let rawText = ''

  if (contentType.includes('application/json')) {
    try {
      body = await res.json()
    } catch {
      body = null
    }
  } else {
    try {
      rawText = (await res.text()).trim()
    } catch {
      rawText = ''
    }
  }

  if (!res.ok) {
    const message = body?.message || rawText || `Request failed (${res.status})`
    throw new ApiError(message, body?.code ?? -1, res.status)
  }

  if (!body || body.code !== 0) {
    throw new ApiError(body?.message || rawText || `Request failed (${res.status})`, body?.code ?? -1, res.status)
  }

  return body.data as T
}

export class ApiError extends Error {
  constructor(
    message: string,
    public code: number,
    public status: number,
  ) {
    super(message)
    this.name = 'ApiError'
  }
}

export const api = {
  get<T>(path: string): Promise<T> {
    return request<T>(path, { method: 'GET' })
  },

  post<T>(path: string, body?: unknown): Promise<T> {
    return request<T>(path, {
      method: 'POST',
      body: body ? JSON.stringify(body) : undefined,
    })
  },

  put<T>(path: string, body?: unknown): Promise<T> {
    return request<T>(path, {
      method: 'PUT',
      body: body ? JSON.stringify(body) : undefined,
    })
  },

  delete<T>(path: string): Promise<T> {
    return request<T>(path, { method: 'DELETE' })
  },
}
