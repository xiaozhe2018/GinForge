<template>
  <div class="doc-content">
    <el-scrollbar ref="scrollbarRef" class="content-scrollbar">
      <div class="content-container">
        <!-- 文档内容 -->
        <article 
          v-if="content"
          class="markdown-body"
          v-html="renderedContent"
        ></article>

        <!-- 加载中 -->
        <div v-else-if="loading" class="loading-container">
          <el-icon class="is-loading"><Loading /></el-icon>
          <p>加载中...</p>
        </div>

        <!-- 空状态 -->
        <div v-else class="empty-container">
          <el-icon><Document /></el-icon>
          <p>选择左侧文档开始阅读</p>
        </div>
      </div>
    </el-scrollbar>

    <!-- 浮动目录 -->
    <div v-if="toc.length > 0" class="doc-toc">
      <div class="toc-title">目录</div>
      <div class="toc-list">
        <a
          v-for="item in toc"
          :key="item.id"
          :href="`#${item.id}`"
          :class="['toc-item', `toc-level-${item.level}`, { 'is-active': activeId === item.id }]"
          @click.prevent="scrollToAnchor(item.id)"
        >
          {{ item.title }}
        </a>
      </div>
    </div>

    <!-- 返回顶部 -->
    <el-backtop :right="40" :bottom="100" target=".content-scrollbar .el-scrollbar__wrap" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick, onMounted, onUnmounted } from 'vue'
import { Loading, Document } from '@element-plus/icons-vue'
import { renderMarkdown, extractToc, type TocItem } from '@/utils/markdown'

const props = defineProps<{
  content?: string
  loading?: boolean
}>()

// 滚动容器
const scrollbarRef = ref()

// 渲染的 HTML
const renderedContent = computed(() => {
  if (!props.content) return ''
  return renderMarkdown(props.content)
})

// 目录
const toc = ref<TocItem[]>([])

// 当前激活的标题
const activeId = ref('')

// 提取目录
watch(() => props.content, (newContent) => {
  if (newContent) {
    toc.value = extractToc(newContent)
  } else {
    toc.value = []
  }
}, { immediate: true })

// 滚动到锚点
const scrollToAnchor = (id: string) => {
  const element = document.getElementById(id)
  if (element && scrollbarRef.value) {
    const scrollWrap = scrollbarRef.value.$el.querySelector('.el-scrollbar__wrap')
    if (scrollWrap) {
      const offsetTop = element.offsetTop - 80
      scrollWrap.scrollTo({
        top: offsetTop,
        behavior: 'smooth'
      })
      activeId.value = id
    }
  }
}

// 监听滚动事件，更新当前激活的标题
let scrollTimer: any = null
const handleScroll = () => {
  if (scrollTimer) {
    clearTimeout(scrollTimer)
  }
  
  scrollTimer = setTimeout(() => {
    const scrollWrap = scrollbarRef.value?.$el.querySelector('.el-scrollbar__wrap')
    if (!scrollWrap || toc.value.length === 0) return
    
    const scrollTop = scrollWrap.scrollTop
    const headings = toc.value.map(item => ({
      id: item.id,
      offset: document.getElementById(item.id)?.offsetTop || 0
    }))
    
    // 找到当前滚动位置对应的标题
    for (let i = headings.length - 1; i >= 0; i--) {
      if (scrollTop + 100 >= headings[i].offset) {
        activeId.value = headings[i].id
        break
      }
    }
  }, 100)
}

// 添加滚动监听
onMounted(() => {
  nextTick(() => {
    const scrollWrap = scrollbarRef.value?.$el.querySelector('.el-scrollbar__wrap')
    if (scrollWrap) {
      scrollWrap.addEventListener('scroll', handleScroll)
    }
    
    // 添加代码复制按钮
    addCopyButtons()
  })
})

// 移除滚动监听
onUnmounted(() => {
  const scrollWrap = scrollbarRef.value?.$el.querySelector('.el-scrollbar__wrap')
  if (scrollWrap) {
    scrollWrap.removeEventListener('scroll', handleScroll)
  }
  if (scrollTimer) {
    clearTimeout(scrollTimer)
  }
})

// 为代码块添加复制按钮
const addCopyButtons = () => {
  nextTick(() => {
    const codeBlocks = document.querySelectorAll('.markdown-body pre code')
    codeBlocks.forEach((block) => {
      const pre = block.parentElement
      if (pre && !pre.querySelector('.copy-btn')) {
        const button = document.createElement('button')
        button.className = 'copy-btn'
        button.textContent = '复制'
        button.onclick = () => copyCode(block.textContent || '', button)
        pre.appendChild(button)
      }
    })
  })
}

// 复制代码
const copyCode = async (text: string, button: HTMLButtonElement) => {
  try {
    await navigator.clipboard.writeText(text)
    button.textContent = '已复制!'
    button.classList.add('copied')
    setTimeout(() => {
      button.textContent = '复制'
      button.classList.remove('copied')
    }, 2000)
  } catch (err) {
    console.error('复制失败:', err)
  }
}

// 监听内容变化，重新添加复制按钮
watch(renderedContent, () => {
  nextTick(() => {
    addCopyButtons()
    // 滚动到顶部
    const scrollWrap = scrollbarRef.value?.$el.querySelector('.el-scrollbar__wrap')
    if (scrollWrap) {
      scrollWrap.scrollTo({ top: 0, behavior: 'smooth' })
    }
  })
})
</script>

<style scoped>
.doc-content {
  position: relative;
  height: 100%;
  display: flex;
}

.content-scrollbar {
  flex: 1;
  overflow: auto;
}

.content-container {
  max-width: 900px;
  margin: 0 auto;
  padding: 40px 40px 100px;
  min-height: 100%;
}

.loading-container,
.empty-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  color: #909399;
}

.loading-container .el-icon,
.empty-container .el-icon {
  font-size: 64px;
  margin-bottom: 16px;
}

.loading-container p,
.empty-container p {
  font-size: 16px;
  margin: 0;
}

/* 浮动目录 */
.doc-toc {
  width: 200px;
  padding: 40px 20px;
  border-left: 1px solid #e4e7ed;
  overflow-y: auto;
}

.toc-title {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 12px;
}

.toc-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.toc-item {
  font-size: 13px;
  color: #606266;
  text-decoration: none;
  padding: 4px 8px;
  border-left: 2px solid transparent;
  transition: all 0.2s;
  line-height: 1.5;
}

.toc-item:hover {
  color: #667eea;
  border-left-color: #667eea;
}

.toc-item.is-active {
  color: #667eea;
  border-left-color: #667eea;
  background: #f5f7fa;
}

.toc-level-1 {
  padding-left: 8px;
}

.toc-level-2 {
  padding-left: 20px;
}

.toc-level-3 {
  padding-left: 32px;
}

/* 响应式 */
@media (max-width: 1200px) {
  .doc-toc {
    display: none;
  }
}
</style>

