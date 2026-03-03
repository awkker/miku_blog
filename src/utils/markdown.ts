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
const jsLikeLanguages = new Set(['js', 'jsx', 'ts', 'tsx', 'javascript', 'typescript'])
const jsonLikeLanguages = new Set(['json', 'json5'])
const shellLikeLanguages = new Set(['sh', 'bash', 'zsh', 'shell'])
const htmlLikeLanguages = new Set(['html', 'xml', 'vue'])
const cssLikeLanguages = new Set(['css', 'scss', 'sass', 'less'])

const jsKeywords = new Set([
  'const',
  'let',
  'var',
  'function',
  'return',
  'if',
  'else',
  'switch',
  'case',
  'break',
  'continue',
  'for',
  'while',
  'do',
  'new',
  'class',
  'extends',
  'implements',
  'interface',
  'type',
  'enum',
  'import',
  'from',
  'export',
  'default',
  'async',
  'await',
  'try',
  'catch',
  'finally',
  'throw',
  'typeof',
  'instanceof',
  'in',
  'of',
  'as',
])

const jsLiterals = new Set(['true', 'false', 'null', 'undefined'])
const shellKeywords = new Set([
  'if',
  'then',
  'fi',
  'for',
  'in',
  'do',
  'done',
  'case',
  'esac',
  'function',
  'echo',
  'export',
  'cd',
  'ls',
  'cat',
  'grep',
  'sed',
  'awk',
])

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
    const block = renderCodeBlock(trimmedCode, language || 'text')
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

function renderCodeBlock(rawCode: string, rawLanguage: string): string {
  const language = normalizeLanguage(rawLanguage)
  const safeLanguage = escapeHtml(language)
  const languageLabel = escapeHtml(getLanguageLabel(language))
  const encodedCode = encodeURIComponent(rawCode)
  const linesHtml = renderCodeLines(rawCode, language)

  return `<div class="code-block language-${safeLanguage}"><div class="code-toolbar"><div class="window-dots" aria-hidden="true"><span class="dot red"></span><span class="dot yellow"></span><span class="dot green"></span></div><span class="code-lang">${languageLabel}</span><button type="button" class="code-copy-btn" data-code="${encodedCode}" aria-label="复制代码">复制</button></div><pre><code class="language-${safeLanguage}">${linesHtml}</code></pre></div>`
}

function renderCodeLines(rawCode: string, language: string): string {
  const lines = rawCode.split('\n')
  const normalizedLines = lines.length > 0 ? lines : ['']
  return normalizedLines
    .map((line, index) => {
      const highlighted = highlightCodeLine(line, language)
      const lineContent = highlighted || '&nbsp;'
      return `<span class="code-line"><span class="line-no">${index + 1}</span><span class="line-content">${lineContent}</span></span>`
    })
    .join('')
}

function normalizeLanguage(rawLanguage: string): string {
  const language = rawLanguage.trim().toLowerCase()
  if (!language) {
    return 'text'
  }
  if (language === 'javascript') {
    return 'js'
  }
  if (language === 'typescript') {
    return 'ts'
  }
  if (language === 'shell') {
    return 'bash'
  }
  if (language === 'yml') {
    return 'yaml'
  }
  return language
}

function getLanguageLabel(language: string): string {
  const labelMap: Record<string, string> = {
    js: 'JavaScript',
    jsx: 'JSX',
    ts: 'TypeScript',
    tsx: 'TSX',
    json: 'JSON',
    json5: 'JSON5',
    bash: 'Bash',
    sh: 'Shell',
    zsh: 'Zsh',
    shell: 'Shell',
    html: 'HTML',
    vue: 'Vue',
    xml: 'XML',
    css: 'CSS',
    scss: 'SCSS',
    sass: 'Sass',
    less: 'Less',
    yaml: 'YAML',
    yml: 'YAML',
    md: 'Markdown',
    text: 'Text',
    txt: 'Text',
  }
  return labelMap[language] ?? language.toUpperCase()
}

function highlightCodeLine(line: string, language: string): string {
  if (!line) {
    return ''
  }

  if (jsLikeLanguages.has(language)) {
    return highlightWithPattern(
      line,
      /\/\/.*|\/\*.*?\*\/|"(?:\\.|[^"\\])*"|'(?:\\.|[^'\\])*'|`(?:\\.|[^`\\])*`|\b(?:const|let|var|function|return|if|else|switch|case|break|continue|for|while|do|new|class|extends|implements|interface|type|enum|import|from|export|default|async|await|try|catch|finally|throw|typeof|instanceof|in|of|as)\b|\b(?:true|false|null|undefined)\b|\b\d+(?:\.\d+)?\b|\b[A-Za-z_$][\w$]*(?=\s*\()|[=+\-*/%<>!&|^~?:]+/g,
      (token) => {
        if (token.startsWith('//') || token.startsWith('/*')) {
          return 'token-comment'
        }
        if (token.startsWith('"') || token.startsWith("'") || token.startsWith('`')) {
          return 'token-string'
        }
        if (jsKeywords.has(token)) {
          return 'token-keyword'
        }
        if (jsLiterals.has(token)) {
          return 'token-literal'
        }
        if (/^\d/.test(token)) {
          return 'token-number'
        }
        if (/^[A-Za-z_$]/.test(token)) {
          return 'token-function'
        }
        return 'token-operator'
      },
    )
  }

  if (jsonLikeLanguages.has(language)) {
    return highlightWithPattern(
      line,
      /"(?:\\.|[^"\\])*"(?=\s*:)|"(?:\\.|[^"\\])*"|\b-?\d+(?:\.\d+)?(?:e[+-]?\d+)?\b|\b(?:true|false|null)\b|[{}\[\],:]/gi,
      (token, index, source) => {
        if (/^"/.test(token) && source.slice(index + token.length).trimStart().startsWith(':')) {
          return 'token-key'
        }
        if (/^"/.test(token)) {
          return 'token-string'
        }
        if (/^(true|false|null)$/i.test(token)) {
          return 'token-literal'
        }
        if (/^-?\d/.test(token)) {
          return 'token-number'
        }
        return 'token-operator'
      },
    )
  }

  if (shellLikeLanguages.has(language)) {
    return highlightWithPattern(
      line,
      /#.*$|"(?:\\.|[^"\\])*"|'(?:\\.|[^'\\])*'|\$\{?[A-Za-z_][\w]*\}?|\b(?:if|then|fi|for|in|do|done|case|esac|function|echo|export|cd|ls|cat|grep|sed|awk)\b|--?[a-zA-Z-]+|\b\d+\b/g,
      (token) => {
        if (token.startsWith('#')) {
          return 'token-comment'
        }
        if (token.startsWith('"') || token.startsWith("'")) {
          return 'token-string'
        }
        if (token.startsWith('$')) {
          return 'token-literal'
        }
        if (shellKeywords.has(token)) {
          return 'token-keyword'
        }
        if (token.startsWith('-')) {
          return 'token-operator'
        }
        if (/^\d/.test(token)) {
          return 'token-number'
        }
        return 'token-function'
      },
    )
  }

  if (htmlLikeLanguages.has(language)) {
    return highlightWithPattern(
      line,
      /<!--.*?-->|<\/?[A-Za-z][\w:-]*|\/?>|"(?:\\.|[^"\\])*"|'(?:\\.|[^'\\])*'|[=]/g,
      (token) => {
        if (token.startsWith('<!--')) {
          return 'token-comment'
        }
        if (token.startsWith('"') || token.startsWith("'")) {
          return 'token-string'
        }
        if (token === '=' || token === '/>' || token === '>') {
          return 'token-operator'
        }
        return 'token-keyword'
      },
    )
  }

  if (cssLikeLanguages.has(language)) {
    return highlightWithPattern(
      line,
      /\/\*.*?\*\/|"(?:\\.|[^"\\])*"|'(?:\\.|[^'\\])*'|#[A-Za-z0-9_-]+|\.[A-Za-z0-9_-]+|@[A-Za-z-]+|[A-Za-z-]+(?=\s*:)|\b\d+(?:\.\d+)?(?:px|em|rem|%)?\b|[{}:;(),]/g,
      (token) => {
        if (token.startsWith('/*')) {
          return 'token-comment'
        }
        if (token.startsWith('"') || token.startsWith("'")) {
          return 'token-string'
        }
        if (token.startsWith('#') || token.startsWith('.')) {
          return 'token-function'
        }
        if (token.startsWith('@') || /[{}:;(),]/.test(token)) {
          return 'token-operator'
        }
        if (/^\d/.test(token)) {
          return 'token-number'
        }
        return 'token-keyword'
      },
    )
  }

  return escapeHtml(line)
}

function highlightWithPattern(line: string, pattern: RegExp, getClassName: (token: string, index: number, source: string) => string): string {
  let output = ''
  let lastIndex = 0
  for (const match of line.matchAll(pattern)) {
    const token = match[0] ?? ''
    const index = match.index ?? 0
    if (index > lastIndex) {
      output += escapeHtml(line.slice(lastIndex, index))
    }
    const className = getClassName(token, index, line)
    output += `<span class="${className}">${escapeHtml(token)}</span>`
    lastIndex = index + token.length
  }
  if (lastIndex < line.length) {
    output += escapeHtml(line.slice(lastIndex))
  }
  return output
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
