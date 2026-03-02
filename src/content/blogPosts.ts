export interface BlogPostMeta {
  slug: string
  title: string
  date: string
  category: string
  summary: string
  cover: string
  views: string
}

export interface BlogPost extends BlogPostMeta {
  markdown: string
}

const postMetas: BlogPostMeta[] = [
  {
    slug: 'morning-mist-composition',
    title: '把晨雾装进镜头：一次山谷徒步的构图笔记',
    date: '2026-02-28',
    category: '旅行风景',
    summary: '从逆光到侧逆光，记录我在清晨山谷里如何处理层次、前景和颜色，让画面保留呼吸感。',
    cover: '/photo/封面.avif',
    views: '2,978',
  },
  {
    slug: 'liquid-glass-performance',
    title: '液态玻璃效果实战：从视觉到性能的平衡',
    date: '2026-02-20',
    category: '技术随笔',
    summary: '拆解液态玻璃卡片在 Vue 中的封装思路，包含位移贴图、降级策略和参数调优方式。',
    cover: '/photo/封面.avif',
    views: '4,932',
  },
  {
    slug: 'tokyo-rain-walk-routes',
    title: '东京雨夜散步地图：5 条安静路线',
    date: '2026-02-15',
    category: '旅行风景',
    summary: '避开热门街区后，雨夜的城市会更柔软。这 5 条路线适合一个人慢慢走。',
    cover: '/photo/封面.avif',
    views: '3,846',
  },
  {
    slug: 'weekly-journal-system',
    title: '写给未来自己的周记系统：轻量但可持续',
    date: '2026-02-10',
    category: '日常片段',
    summary: '如何用固定模版与低门槛习惯，稳定记录生活，而不是把记录变成另一种压力。',
    cover: '/photo/封面.avif',
    views: '2,347',
  },
  {
    slug: 'night-photo-lightroom-preset',
    title: '夜景降噪和色彩修正：我的 Lightroom 预设拆解',
    date: '2026-02-04',
    category: '摄影后期',
    summary: '针对高 ISO 夜景样片，分享一套可复用的降噪、曲线和肤色保护流程。',
    cover: '/photo/封面.avif',
    views: '3,215',
  },
  {
    slug: 'deployment-debug-checklist',
    title: '从一次失败部署说起：我的排障流程清单',
    date: '2026-01-30',
    category: '技术随笔',
    summary: '把线上问题拆成可验证步骤，减少凭感觉排错，提升问题定位速度。',
    cover: '/photo/封面.avif',
    views: '2,116',
  },
]

const markdownModules = import.meta.glob('./posts/*.md', {
  eager: true,
  import: 'default',
  query: '?raw',
}) as Record<string, string>

const posts: BlogPost[] = postMetas.map((meta) => ({
  ...meta,
  markdown: markdownModules[`./posts/${meta.slug}.md`] ?? '',
}))

function sortByDateDesc<T extends { date: string }>(items: T[]): T[] {
  return [...items].sort((a, b) => b.date.localeCompare(a.date))
}

export function getAllBlogPostMetas(): BlogPostMeta[] {
  return sortByDateDesc(postMetas)
}

export function getHotBlogPostMetas(limit = 5): BlogPostMeta[] {
  return [...postMetas]
    .sort((a, b) => Number.parseInt(b.views.replace(/,/g, ''), 10) - Number.parseInt(a.views.replace(/,/g, ''), 10))
    .slice(0, limit)
}

export function getBlogPostBySlug(slug: string): BlogPost | null {
  return posts.find((post) => post.slug === slug) ?? null
}

