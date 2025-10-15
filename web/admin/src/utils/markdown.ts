import MarkdownIt from 'markdown-it'
import anchor from 'markdown-it-anchor'
import hljs from 'highlight.js'

// 创建 Markdown 渲染器
const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
  highlight: function (str: string, lang: string) {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return `<pre class="hljs"><code class="language-${lang}">${
          hljs.highlight(str, { language: lang, ignoreIllegals: true }).value
        }</code></pre>`
      } catch (__) {}
    }
    return `<pre class="hljs"><code>${md.utils.escapeHtml(str)}</code></pre>`
  }
})

// 添加锚点插件
md.use(anchor, {
  permalink: anchor.permalink.headerLink({
    safariReaderFix: true,
  })
})

/**
 * 渲染 Markdown 为 HTML
 */
export function renderMarkdown(markdown: string): string {
  return md.render(markdown)
}

/**
 * 提取 Markdown 标题作为目录
 */
export interface TocItem {
  level: number
  title: string
  id: string
  children?: TocItem[]
}

export function extractToc(markdown: string): TocItem[] {
  const toc: TocItem[] = []
  const lines = markdown.split('\n')
  
  for (const line of lines) {
    const match = line.match(/^(#{1,6})\s+(.+)$/)
    if (match) {
      const level = match[1].length
      const title = match[2].trim()
      const id = title
        .toLowerCase()
        .replace(/[^\w\u4e00-\u9fa5]+/g, '-')
        .replace(/^-+|-+$/g, '')
      
      toc.push({
        level,
        title,
        id
      })
    }
  }
  
  return toc
}

/**
 * 构建层级目录树
 */
export function buildTocTree(toc: TocItem[]): TocItem[] {
  const tree: TocItem[] = []
  const stack: TocItem[] = []
  
  for (const item of toc) {
    const newItem = { ...item, children: [] }
    
    // 找到合适的父节点
    while (stack.length > 0 && stack[stack.length - 1].level >= item.level) {
      stack.pop()
    }
    
    if (stack.length === 0) {
      tree.push(newItem)
    } else {
      const parent = stack[stack.length - 1]
      if (!parent.children) {
        parent.children = []
      }
      parent.children.push(newItem)
    }
    
    stack.push(newItem)
  }
  
  return tree
}

export default md

