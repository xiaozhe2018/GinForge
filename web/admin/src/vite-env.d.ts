/// <reference types="vite/client" />

// Markdown 文件导入类型
declare module '*.md' {
  const content: string
  export default content
}

declare module '*.md?raw' {
  const content: string
  export default content
}
