export interface BlogPostMeta {
  slug: string
  title: string
  date: string
  category: string
  summary: string
  cover: string
  views: string
  likes: number
  shares: number
}

export interface BlogPost extends BlogPostMeta {
  markdown: string
}

export interface BlogAuthorStats {
  postCount: number
  totalViews: number
  totalLikes: number
  totalBlogVisits: number
}

interface FrontmatterMeta {
  title?: string
  date?: string
  category?: string
  summary?: string
  cover?: string
  views?: string
}

interface BlogPostRecord {
  slug: string
  title: string
  date: string
  category: string
  summary: string
  cover: string
  baseViews: number
  markdown: string
}

const VIEW_STORAGE_KEY = 'miku_blog_post_views'
const LIKE_COUNT_STORAGE_KEY = 'miku_blog_post_like_count'
const SHARE_COUNT_STORAGE_KEY = 'miku_blog_post_share_count'
const BLOG_LIST_VISIT_STORAGE_KEY = 'miku_blog_list_page_visits'
const BLOG_STATS_UPDATED_EVENT = 'miku_blog_stats_updated'
const NUMBER_FORMATTER = new Intl.NumberFormat('zh-CN')
const SUMMARY_EXCERPT_LENGTH = 72

const markdownModules = import.meta.glob('./posts/*.md', {
  eager: true,
  import: 'default',
  query: '?raw',
}) as Record<string, string>

const postRecords: BlogPostRecord[] = Object.entries(markdownModules)
  .map(([path, raw]) => createPostRecord(path, raw))
  .filter((record): record is BlogPostRecord => record !== null)

const postRecordMap = new Map(postRecords.map((record) => [record.slug, record]))

function sortByDateDesc<T extends { date: string }>(items: T[]): T[] {
  return [...items].sort((a, b) => b.date.localeCompare(a.date))
}

export function getAllBlogPostMetas(): BlogPostMeta[] {
  const viewMap = readViewMap()
  const likeMap = readStoredCountMap(LIKE_COUNT_STORAGE_KEY)
  const shareMap = readStoredCountMap(SHARE_COUNT_STORAGE_KEY)
  return sortByDateDesc(postRecords.map((record) => toBlogPostMeta(record, viewMap, likeMap, shareMap)))
}

export function getHotBlogPostMetas(limit = 5): BlogPostMeta[] {
  const viewMap = readViewMap()
  const likeMap = readStoredCountMap(LIKE_COUNT_STORAGE_KEY)
  const shareMap = readStoredCountMap(SHARE_COUNT_STORAGE_KEY)
  return [...postRecords]
    .sort((a, b) => getViewCount(b.slug, b.baseViews, viewMap) - getViewCount(a.slug, a.baseViews, viewMap))
    .slice(0, limit)
    .map((record) => toBlogPostMeta(record, viewMap, likeMap, shareMap))
}

export function getBlogPostBySlug(slug: string): BlogPost | null {
  const record = postRecordMap.get(slug)
  if (!record) {
    return null
  }
  const viewMap = readViewMap()
  const likeMap = readStoredCountMap(LIKE_COUNT_STORAGE_KEY)
  const shareMap = readStoredCountMap(SHARE_COUNT_STORAGE_KEY)
  return {
    ...toBlogPostMeta(record, viewMap, likeMap, shareMap),
    markdown: record.markdown,
  }
}

export function incrementBlogPostViews(slug: string): void {
  const record = postRecordMap.get(slug)
  if (!record) {
    return
  }
  const map = readViewMap()
  const current = getStoredViewCount(map[slug], record.baseViews)
  map[slug] = current + 1
  writeViewMap(map)
  notifyStatsUpdated()
}

export function incrementBlogListPageVisits(): void {
  const storage = getLocalStorage()
  if (!storage) {
    return
  }
  const current = Number.parseInt(storage.getItem(BLOG_LIST_VISIT_STORAGE_KEY) ?? '0', 10)
  const next = Number.isFinite(current) && current >= 0 ? current + 1 : 1
  storage.setItem(BLOG_LIST_VISIT_STORAGE_KEY, String(next))
  notifyStatsUpdated()
}

export function getBlogAuthorStats(): BlogAuthorStats {
  const viewMap = readViewMap()
  const likeMap = readStoredCountMap(LIKE_COUNT_STORAGE_KEY)
  const storage = getLocalStorage()
  const rawBlogVisits = storage?.getItem(BLOG_LIST_VISIT_STORAGE_KEY) ?? '0'
  const totalBlogVisits = Number.parseInt(rawBlogVisits, 10)

  const totalViews = postRecords.reduce((sum, record) => sum + getViewCount(record.slug, record.baseViews, viewMap), 0)
  const totalLikes = postRecords.reduce((sum, record) => sum + getStoredViewCount(likeMap[record.slug], 0), 0)

  return {
    postCount: postRecords.length,
    totalViews,
    totalLikes,
    totalBlogVisits: Number.isFinite(totalBlogVisits) && totalBlogVisits > 0 ? totalBlogVisits : 0,
  }
}

function createPostRecord(path: string, rawMarkdown: string): BlogPostRecord | null {
  const segments = path.split('/')
  const filename = segments[segments.length - 1]
  const slug = filename?.replace(/\.md$/, '')
  if (!slug) {
    return null
  }

  const { frontmatter, content } = parseMarkdownFile(rawMarkdown)
  const markdown = content.trim()

  return {
    slug,
    title: frontmatter.title ?? extractTitle(markdown) ?? slug,
    date: frontmatter.date ?? '1970-01-01',
    category: frontmatter.category ?? '未分类',
    summary: frontmatter.summary ?? extractSummary(markdown),
    cover: frontmatter.cover ?? '/photo/封面.avif',
    baseViews: parseViewCount(frontmatter.views),
    markdown,
  }
}

function parseMarkdownFile(rawMarkdown: string): { frontmatter: FrontmatterMeta; content: string } {
  const normalized = rawMarkdown.replace(/\r\n?/g, '\n')
  if (!normalized.startsWith('---\n')) {
    return {
      frontmatter: {},
      content: normalized,
    }
  }

  const frontmatterEnd = normalized.indexOf('\n---\n', 4)
  if (frontmatterEnd === -1) {
    return {
      frontmatter: {},
      content: normalized,
    }
  }

  const frontmatterRaw = normalized.slice(4, frontmatterEnd)
  const content = normalized.slice(frontmatterEnd + '\n---\n'.length)
  return {
    frontmatter: parseFrontmatter(frontmatterRaw),
    content,
  }
}

function parseFrontmatter(rawFrontmatter: string): FrontmatterMeta {
  const frontmatter: FrontmatterMeta = {}
  const lines = rawFrontmatter.split('\n')
  for (const rawLine of lines) {
    const line = rawLine.trim()
    if (!line || line.startsWith('#')) {
      continue
    }
    const separatorIndex = line.indexOf(':')
    if (separatorIndex === -1) {
      continue
    }
    const rawKey = line.slice(0, separatorIndex).trim()
    const key = rawKey.toLowerCase()
    const value = stripWrappedQuotes(line.slice(separatorIndex + 1).trim())
    if (!value) {
      continue
    }
    if (key === 'title' || key === 'date' || key === 'category' || key === 'summary' || key === 'cover' || key === 'views') {
      frontmatter[key] = value
    }
  }
  return frontmatter
}

function stripWrappedQuotes(input: string): string {
  if (input.length < 2) {
    return input
  }
  const first = input[0]
  const last = input[input.length - 1]
  if ((first === '"' && last === '"') || (first === "'" && last === "'")) {
    return input.slice(1, -1)
  }
  return input
}

function extractTitle(markdown: string): string | null {
  const headingMatch = markdown.match(/^#\s+(.+)$/m)
  return headingMatch?.[1]?.trim() ?? null
}

function extractSummary(markdown: string): string {
  const plainText = markdown
    .replace(/```[\s\S]*?```/g, ' ')
    .replace(/^#{1,6}\s+/gm, '')
    .replace(/^\s*[-*+]\s+/gm, '')
    .replace(/^\s*\d+\.\s+/gm, '')
    .replace(/^>\s?/gm, '')
    .replace(/!\[[^\]]*]\([^)]+\)/g, ' ')
    .replace(/\[([^\]]+)\]\([^)]+\)/g, '$1')
    .replace(/[*_~`]/g, '')
    .replace(/\s+/g, ' ')
    .trim()

  if (!plainText) {
    return '暂无摘要'
  }
  if (plainText.length <= SUMMARY_EXCERPT_LENGTH) {
    return plainText
  }
  return `${plainText.slice(0, SUMMARY_EXCERPT_LENGTH).trimEnd()}...`
}

function parseViewCount(rawViews: string | undefined): number {
  if (!rawViews) {
    return 0
  }
  const parsed = Number.parseInt(rawViews.replace(/,/g, ''), 10)
  return Number.isFinite(parsed) && parsed > 0 ? parsed : 0
}

function toBlogPostMeta(
  record: BlogPostRecord,
  viewMap: Record<string, number>,
  likeMap: Record<string, number>,
  shareMap: Record<string, number>,
): BlogPostMeta {
  return {
    slug: record.slug,
    title: record.title,
    date: record.date,
    category: record.category,
    summary: record.summary,
    cover: record.cover,
    views: formatViewCount(getViewCount(record.slug, record.baseViews, viewMap)),
    likes: getStoredViewCount(likeMap[record.slug], 0),
    shares: getStoredViewCount(shareMap[record.slug], 0),
  }
}

function formatViewCount(views: number): string {
  return NUMBER_FORMATTER.format(views)
}

function getViewCount(slug: string, baseViews: number, map: Record<string, number>): number {
  return getStoredViewCount(map[slug], baseViews)
}

function getStoredViewCount(rawStoredViews: unknown, fallback: number): number {
  if (typeof rawStoredViews === 'number' && Number.isFinite(rawStoredViews) && rawStoredViews >= 0) {
    return Math.max(Math.floor(rawStoredViews), fallback)
  }
  return fallback
}

function readViewMap(): Record<string, number> {
  return readStoredCountMap(VIEW_STORAGE_KEY)
}

function readStoredCountMap(storageKey: string): Record<string, number> {
  const storage = getLocalStorage()
  if (!storage) {
    return {}
  }
  const rawMap = storage.getItem(storageKey)
  if (!rawMap) {
    return {}
  }
  try {
    const parsed = JSON.parse(rawMap)
    if (!parsed || typeof parsed !== 'object') {
      return {}
    }
    return parsed as Record<string, number>
  } catch {
    return {}
  }
}

function writeViewMap(map: Record<string, number>): void {
  const storage = getLocalStorage()
  if (!storage) {
    return
  }
  storage.setItem(VIEW_STORAGE_KEY, JSON.stringify(map))
}

function getLocalStorage(): Storage | null {
  if (typeof window === 'undefined') {
    return null
  }
  try {
    return window.localStorage
  } catch {
    return null
  }
}

function notifyStatsUpdated(): void {
  if (typeof window === 'undefined') {
    return
  }
  window.dispatchEvent(new CustomEvent(BLOG_STATS_UPDATED_EVENT))
}

export function getBlogStatsUpdatedEventName(): string {
  return BLOG_STATS_UPDATED_EVENT
}
