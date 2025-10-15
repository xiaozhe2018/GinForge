<template>
  <div class="documentation-page">
    <!-- 导航栏 -->
    <div class="doc-sidebar">
      <DocNav :current-path="currentDocPath" @doc-change="handleDocChange" />
    </div>

    <!-- 内容区域 -->
    <div class="doc-main">
      <DocContent :content="docContent" :loading="loading" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import DocNav from './components/DocNav.vue'
import DocContent from './components/DocContent.vue'
import type { DocItem } from '@/config/docs'

const route = useRoute()
const router = useRouter()

// 使用 Vite 的 glob import 预加载所有文档
const docModules = import.meta.glob('../../docs/**/*.md', { as: 'raw', eager: true })

// 当前文档路径
const currentDocPath = ref<string>('')

// 文档内容
const docContent = ref<string>('')

// 加载状态
const loading = ref(false)

// 加载文档内容
const loadDocContent = async (path: string) => {
  loading.value = true
  
  try {
    // 构建文件路径
    const filePath = `../../docs/${path}.md`
    
    // 从预加载的模块中获取内容
    const content = docModules[filePath]
    
    if (content) {
      docContent.value = content as string
      console.log('文档加载成功:', path)
    } else {
      console.error('文档不存在:', path, '可用文档:', Object.keys(docModules))
      docContent.value = `# 文档加载失败\n\n无法加载文档：${path}\n\n请检查文件路径是否正确。`
    }
  } catch (error) {
    console.error('加载文档失败:', error, '路径:', path)
    docContent.value = `# 文档加载失败\n\n无法加载文档：${path}\n\n错误信息：${error}`
  } finally {
    loading.value = false
  }
}

// 处理文档切换
const handleDocChange = (doc: DocItem) => {
  currentDocPath.value = doc.path
  loadDocContent(doc.path)
  
  // 更新 URL
  router.push({ 
    name: 'Documentation',
    query: { doc: doc.path }
  })
}

// 组件挂载时加载初始文档
onMounted(() => {
  const docPath = (route.query.doc as string) || 'getting-started/introduction'
  currentDocPath.value = docPath
  loadDocContent(docPath)
  
  // 调试：输出所有可用文档
  console.log('可用文档列表:', Object.keys(docModules))
})
</script>

<style scoped>
.documentation-page {
  display: flex;
  height: 100%;
  background: #fff;
}

.doc-sidebar {
  width: 280px;
  flex-shrink: 0;
  height: 100%;
  overflow: hidden;
}

.doc-main {
  flex: 1;
  height: 100%;
  overflow: hidden;
}

/* 响应式 */
@media (max-width: 768px) {
  .documentation-page {
    flex-direction: column;
  }
  
  .doc-sidebar {
    width: 100%;
    height: auto;
    max-height: 300px;
    border-right: none;
    border-bottom: 1px solid #e4e7ed;
  }
}
</style>

<style>
/* 全局导入 Markdown 样式 */
@import '@/styles/markdown.scss';
</style>

