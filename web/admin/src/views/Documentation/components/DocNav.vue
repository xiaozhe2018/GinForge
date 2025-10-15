<template>
  <div class="doc-nav">
    <!-- 搜索框 -->
    <div class="nav-search">
      <el-input
        v-model="searchKeyword"
        placeholder="搜索文档..."
        :prefix-icon="Search"
        clearable
        @input="handleSearch"
      />
    </div>

    <!-- 文档分类 -->
    <el-scrollbar class="nav-scrollbar">
      <div v-if="!searchKeyword" class="nav-categories">
        <div
          v-for="category in categories"
          :key="category.id"
          class="category-group"
        >
          <div class="category-title" @click="toggleCategory(category.id)">
            <el-icon class="category-icon">
              <component :is="category.icon" />
            </el-icon>
            <span>{{ category.name }}</span>
            <el-icon class="expand-icon" :class="{ 'is-expanded': expandedCategories.includes(category.id) }">
              <ArrowRight />
            </el-icon>
          </div>

          <el-collapse-transition>
            <div v-show="expandedCategories.includes(category.id)" class="doc-list">
              <div
                v-for="doc in category.docs"
                :key="doc.id"
                class="doc-item"
                :class="{ 'is-active': currentPath === doc.path }"
                @click="handleDocClick(doc)"
              >
                <span class="doc-title">{{ doc.title }}</span>
              </div>
            </div>
          </el-collapse-transition>
        </div>
      </div>

      <!-- 搜索结果 -->
      <div v-else class="search-results">
        <div v-if="searchResults.length === 0" class="no-results">
          <el-icon><DocumentRemove /></el-icon>
          <p>未找到相关文档</p>
        </div>
        <div v-else class="doc-list">
          <div
            v-for="doc in searchResults"
            :key="doc.id"
            class="doc-item"
            :class="{ 'is-active': currentPath === doc.path }"
            @click="handleDocClick(doc)"
          >
            <span class="doc-title">{{ doc.title }}</span>
            <p v-if="doc.description" class="doc-desc">{{ doc.description }}</p>
          </div>
        </div>
      </div>
    </el-scrollbar>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Search, ArrowRight, DocumentRemove, Rocket, Reading, Star, Document, Upload, TrendCharts, Edit, Promotion } from '@element-plus/icons-vue'
import { docConfig, getAllDocs, type DocItem, type DocCategory } from '@/config/docs'

const props = defineProps<{
  currentPath?: string
}>()

const emit = defineEmits<{
  (e: 'doc-change', doc: DocItem): void
}>()

// 分类数据
const categories = ref<DocCategory[]>(docConfig)

// 展开的分类
const expandedCategories = ref<string[]>([])

// 搜索关键词
const searchKeyword = ref('')

// 搜索结果
const searchResults = ref<DocItem[]>([])

// 切换分类展开/收起
const toggleCategory = (categoryId: string) => {
  const index = expandedCategories.value.indexOf(categoryId)
  if (index > -1) {
    expandedCategories.value.splice(index, 1)
  } else {
    expandedCategories.value.push(categoryId)
  }
}

// 搜索文档
const handleSearch = () => {
  if (!searchKeyword.value.trim()) {
    searchResults.value = []
    return
  }

  const keyword = searchKeyword.value.toLowerCase()
  const allDocs = getAllDocs()
  
  searchResults.value = allDocs.filter(doc => 
    doc.title.toLowerCase().includes(keyword) ||
    doc.description?.toLowerCase().includes(keyword)
  )
}

// 点击文档
const handleDocClick = (doc: DocItem) => {
  emit('doc-change', doc)
}

// 组件挂载时展开第一个分类
onMounted(() => {
  if (categories.value.length > 0) {
    expandedCategories.value.push(categories.value[0].id)
  }
})
</script>

<style scoped>
.doc-nav {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #fff;
  border-right: 1px solid #e4e7ed;
}

.nav-search {
  padding: 16px;
  border-bottom: 1px solid #e4e7ed;
}

.nav-scrollbar {
  flex: 1;
  overflow: auto;
}

.nav-categories {
  padding: 8px 0;
}

.category-group {
  margin-bottom: 4px;
}

.category-title {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  cursor: pointer;
  user-select: none;
  transition: all 0.2s;
}

.category-title:hover {
  background: #f5f7fa;
}

.category-icon {
  font-size: 18px;
  margin-right: 8px;
  color: #667eea;
}

.category-title span {
  flex: 1;
  font-size: 14px;
  font-weight: 600;
  color: #303133;
}

.expand-icon {
  font-size: 14px;
  color: #909399;
  transition: transform 0.2s;
}

.expand-icon.is-expanded {
  transform: rotate(90deg);
}

.doc-list {
  padding-left: 12px;
}

.doc-item {
  padding: 10px 16px;
  cursor: pointer;
  border-left: 2px solid transparent;
  transition: all 0.2s;
}

.doc-item:hover {
  background: #f5f7fa;
  border-left-color: #667eea;
}

.doc-item.is-active {
  background: #ecf5ff;
  border-left-color: #667eea;
}

.doc-item.is-active .doc-title {
  color: #667eea;
  font-weight: 500;
}

.doc-title {
  font-size: 14px;
  color: #606266;
  display: block;
}

.doc-desc {
  font-size: 12px;
  color: #909399;
  margin: 4px 0 0 0;
  line-height: 1.4;
}

.search-results {
  padding: 8px 0;
}

.no-results {
  text-align: center;
  padding: 40px 20px;
  color: #909399;
}

.no-results .el-icon {
  font-size: 48px;
  margin-bottom: 12px;
}

.no-results p {
  font-size: 14px;
  margin: 0;
}
</style>

