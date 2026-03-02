export interface MarkdownHeading {
  id: string
  level: number
  text: string
}

export interface MarkdownRenderResult {
  html: string
  headings: MarkdownHeading[]
}

const codeBlockTokenPrefix = '__CODE_BLOCK_'

function replacePlain(input: string, search: string, replacement: string): string {
  return input.split(search).join(replacement)
}

function escapeHtml(input: string): string {
  return input
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;')
}

function sanitizeUrl(url: string): string {
  if (/^https?:\/\//i.test(url)) {
    return url
  }
  return '#'
}

function stripInlineMarkdown(input: string): string {
  return input
    .replace(/`([^`]+)`/g, '$1')
    .replace(/\*\*([^*]+)\*\*/g, '$1')
    .replace(/\*([^*]+)\*/g, '$1')
    .replace(/\[([^\]]+)\]\(([^)]+)\)/g, '$1')
    .trim()
}

function toHeadingId(text: string, idCountMap: Map<string, number>): string {
  const normalized = text
    .toLowerCase()
    .replace(/[^\p{L}\p{N}\s-]/gu, '')
    .trim()
    .replace(/\s+/g, '-')
  const base = normalized || 'section'
  const count = idCountMap.get(base) ?? 0
  idCountMap.set(base, count + 1)
  return count === 0 ? base : `${base}-${count + 1}`
}

function parseInline(input: string): string {
  let output = escapeHtml(input)
  output = output.replace(/`([^`]+)`/g, '<code>$1</code>')
  output = output.replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>')
  output = output.replace(/\*([^*]+)\*/g, '<em>$1</em>')
  output = output.replace(/\[([^\]]+)\]\(([^)]+)\)/g, (_match: string, label: string, rawUrl: string) => {
    const safeUrl = sanitizeUrl(rawUrl.trim())
    return `<a href="${safeUrl}" target="_blank" rel="noopener noreferrer">${label}</a>`
  })
  return output
}

function isBlockStarter(line: string): boolean {
  return (
    /^#{1,6}\s+/.test(line) ||
    /^>\s?/.test(line) ||
    /^[-*]\s+/.test(line) ||
    /^\d+\.\s+/.test(line) ||
    /^-{3,}\s*$/.test(line) ||
    line.startsWith(codeBlockTokenPrefix)
  )
}

function collectParagraph(lines: string[], start: number): { html: string; next: number } {
  const paragraphLines: string[] = []
  let index = start
  while (index < lines.length) {
    const current = lines[index] ?? ''
    const trimmed = current.trim()
    if (!trimmed) {
      break
    }
    if (index !== start && isBlockStarter(trimmed)) {
      break
    }
    paragraphLines.push(trimmed)
    index += 1
  }
  return {
    html: `<p>${parseInline(paragraphLines.join(' '))}</p>`,
    next: index,
  }
}

function restoreCodeTokens(html: string, codeBlocks: string[]): string {
  let restored = html
  for (let index = 0; index < codeBlocks.length; index += 1) {
    restored = restored.replace(`${codeBlockTokenPrefix}${index}__`, codeBlocks[index] ?? '')
  }
  return restored
}

export function renderMarkdown(markdown: string): MarkdownRenderResult {
  const normalized = markdown.replace(/\r\n?/g, '\n')
  const codeBlocks: string[] = []

  const content = normalized.replace(/```([\w-]*)\n([\s\S]*?)```/g, (_match: string, language: string, code: string) => {
    const trimmedCode = code.replace(/\n$/, '')
    const escapedCode = escapeHtml(trimmedCode)
    const safeLanguage = escapeHtml(language || 'text')
    const block = `<pre><code class="language-${safeLanguage}">${escapedCode}</code></pre>`
    const token = `${codeBlockTokenPrefix}${codeBlocks.length}__`
    codeBlocks.push(block)
    return token
  })

  const lines = content.split('\n')
  const htmlParts: string[] = []
  const headings: MarkdownHeading[] = []
  const idCountMap = new Map<string, number>()

  let index = 0
  while (index < lines.length) {
    const rawLine = lines[index] ?? ''
    const line = rawLine.trim()

    if (!line) {
      index += 1
      continue
    }

    if (line.startsWith(codeBlockTokenPrefix)) {
      htmlParts.push(line)
      index += 1
      continue
    }

    const headingMatch = line.match(/^(#{1,6})\s+(.+)$/)
    if (headingMatch) {
      const level = headingMatch[1]?.length ?? 1
      const headingText = stripInlineMarkdown(headingMatch[2] ?? '')
      const headingId = toHeadingId(headingText, idCountMap)
      htmlParts.push(`<h${level} id="${headingId}">${parseInline(headingMatch[2] ?? '')}</h${level}>`)
      if (level <= 3) {
        headings.push({ id: headingId, level, text: headingText })
      }
      index += 1
      continue
    }

    if (/^-{3,}\s*$/.test(line)) {
      htmlParts.push('<hr />')
      index += 1
      continue
    }

    if (/^>\s?/.test(line)) {
      const quoteLines: string[] = []
      while (index < lines.length) {
        const candidate = (lines[index] ?? '').trim()
        if (!/^>\s?/.test(candidate)) {
          break
        }
        quoteLines.push(candidate.replace(/^>\s?/, ''))
        index += 1
      }
      htmlParts.push(`<blockquote>${parseInline(quoteLines.join('<br />'))}</blockquote>`)
      continue
    }

    if (/^[-*]\s+/.test(line)) {
      const items: string[] = []
      while (index < lines.length) {
        const candidate = (lines[index] ?? '').trim()
        if (!/^[-*]\s+/.test(candidate)) {
          break
        }
        items.push(`<li>${parseInline(candidate.replace(/^[-*]\s+/, ''))}</li>`)
        index += 1
      }
      htmlParts.push(`<ul>${items.join('')}</ul>`)
      continue
    }

    if (/^\d+\.\s+/.test(line)) {
      const items: string[] = []
      while (index < lines.length) {
        const candidate = (lines[index] ?? '').trim()
        if (!/^\d+\.\s+/.test(candidate)) {
          break
        }
        items.push(`<li>${parseInline(candidate.replace(/^\d+\.\s+/, ''))}</li>`)
        index += 1
      }
      htmlParts.push(`<ol>${items.join('')}</ol>`)
      continue
    }

    const paragraph = collectParagraph(lines, index)
    htmlParts.push(paragraph.html)
    index = paragraph.next
  }

  return {
    html: restoreCodeTokens(htmlParts.join('\n'), codeBlocks),
    headings,
  }
}

export function estimateReadingMinutes(markdown: string): number {
  const plainText = markdown
    .replace(/```[\s\S]*?```/g, ' ')
    .replace(/`([^`]+)`/g, '$1')
    .replace(/\[([^\]]+)\]\(([^)]+)\)/g, '$1')
    .replace(/[>#*_~`[\]()!-]/g, ' ')
  const chineseCount = (plainText.match(/[\u4e00-\u9fff]/g) ?? []).length
  const wordCount = (replacePlain(plainText, '\u3000', ' ').replace(/[\u4e00-\u9fff]/g, ' ').match(/\b[\w-]+\b/g) ?? []).length
  return Math.max(1, Math.ceil((chineseCount + wordCount) / 320))
}
