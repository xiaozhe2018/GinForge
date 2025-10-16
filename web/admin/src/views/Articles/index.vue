<template>
  <div class="articles-container">
    <div class="page-header">
      <h2>Articles管理</h2>
    </div>

    <!-- 搜索栏 -->
    <el-card class="search-card" shadow="never">
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="搜索">
          <el-input
            v-model="searchForm.keyword"
            placeholder="请输入关键词"
            clearable
            @keyup.enter="handleSearch"
            style="width: 240px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            <span>搜索</span>
          </el-button>
          <el-button @click="handleReset">
            <el-icon><Refresh /></el-icon>
            <span>重置</span>
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 操作栏 -->
    <el-card class="table-card" shadow="never">
      <template #header>
        <div class="card-header">
          <span>Articles管理列表</span>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            <span>新建Articles管理</span>
          </el-button>
        </div>
      </template>

      <!-- 表格 -->
      <el-table
        v-loading="loading"
        :data="tableData"
        border
        stripe
        style="width: 100%"
      >
        <el-table-column
          prop="id"
          label="文章ID"
        />
        <el-table-column
          prop="title"
          label="标题"
        />
        <el-table-column
          prop="summary"
          label="摘要"
        />
        <el-table-column
          prop="author_id"
          label="作者ID"
        />
        <el-table-column
          prop="author_name"
          label="作者名称"
        />
        <el-table-column
          prop="category_id"
          label="分类ID"
        />
        <el-table-column
          prop="cover_image"
          label="封面图片"
        />
        <el-table-column
          prop="tags"
          label="标签（逗号分隔）"
        />
        <el-table-column
          prop="status"
          label="状态:0草稿,1已发布,2已下线"
          width="100"
        >
          <template #default="{ row }">
            <el-tag :type="row.status ? 'success' : 'info'">
              {{ row.status ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="is_top"
          label="是否置顶"
          width="100"
        >
          <template #default="{ row }">
            <el-tag :type="row.is_top ? 'success' : 'info'">
              {{ row.is_top ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="view_count"
          label="浏览次数"
        />
        <el-table-column
          prop="like_count"
          label="点赞次数"
        />
        <el-table-column
          prop="comment_count"
          label="评论次数"
        />
        <el-table-column
          prop="published_at"
          label="发布时间"
          width="180"
        >
          <template #default="{ row }">
            {{ formatDate(row.published_at) }}
          </template>
        </el-table-column>
        <el-table-column
          prop="created_at"
          label="创建时间"
          width="180"
        >
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column
          prop="updated_at"
          label="更新时间"
          width="180"
        >
          <template #default="{ row }">
            {{ formatDate(row.updated_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" link @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button type="danger" size="small" link @click="handleDelete(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.page_size"
        :page-sizes="[10, 20, 50, 100]"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
        style="margin-top: 20px; justify-content: flex-end"
      />
    </el-card>

    <!-- 表单对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
      @close="handleDialogClose"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="formRules"
        label-width="100px"
      >
        <el-form-item label="标题" prop="title">
          <el-input
            v-model="form.title"
            placeholder="请输入标题"
          />
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <el-input
            v-model="form.content"
            placeholder="请输入内容"
          />
        </el-form-item>
        <el-form-item label="摘要" prop="summary">
          <el-input
            v-model="form.summary"
            placeholder="请输入摘要"
          />
        </el-form-item>
        <el-form-item label="作者ID" prop="author_id">
          <el-input
            v-model="form.author_id"
            placeholder="请输入作者ID"
          />
        </el-form-item>
        <el-form-item label="作者名称" prop="author_name">
          <el-input
            v-model="form.author_name"
            placeholder="请输入作者名称"
          />
        </el-form-item>
        <el-form-item label="分类ID" prop="category_id">
          <el-select v-model="form.category_id" placeholder="请选择分类ID">
            <!-- TODO: 添加选项 -->
            <el-option label="选项1" value="1" />
            <el-option label="选项2" value="2" />
          </el-select>
        </el-form-item>
        <el-form-item label="封面图片" prop="cover_image">
          <el-input
            v-model="form.cover_image"
            placeholder="请输入封面图片"
          />
        </el-form-item>
        <el-form-item label="标签（逗号分隔）" prop="tags">
          <el-input
            v-model="form.tags"
            placeholder="请输入标签（逗号分隔）"
          />
        </el-form-item>
        <el-form-item label="状态:0草稿,1已发布,2已下线" prop="status">
          <el-switch v-model="form.status" />
        </el-form-item>
        <el-form-item label="是否置顶" prop="is_top">
          <el-switch v-model="form.is_top" />
        </el-form-item>
        <el-form-item label="浏览次数" prop="view_count">
          <el-input
            v-model="form.view_count"
            placeholder="请输入浏览次数"
          />
        </el-form-item>
        <el-form-item label="点赞次数" prop="like_count">
          <el-input
            v-model="form.like_count"
            placeholder="请输入点赞次数"
          />
        </el-form-item>
        <el-form-item label="评论次数" prop="comment_count">
          <el-input
            v-model="form.comment_count"
            placeholder="请输入评论次数"
          />
        </el-form-item>
        <el-form-item label="发布时间" prop="published_at">
          <el-input
            v-model="form.published_at"
            placeholder="请输入发布时间"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, FormInstance, FormRules } from 'element-plus'
import { Search, Refresh, Plus } from '@element-plus/icons-vue'
import * as articlesApi from '@/api/articles'

// ========== 数据定义 ==========

const loading = ref(false)
const submitLoading = ref(false)
const tableData = ref<articlesApi.Articles[]>([])

// 搜索表单
const searchForm = reactive({
  keyword: '',
})

// 分页
const pagination = reactive({
  page: 1,
  page_size: 10,
  total: 0
})

// 对话框
const dialogVisible = ref(false)
const dialogTitle = ref('')
const isEdit = ref(false)
const currentId = ref<number | null>(null)

// 表单
const formRef = ref<FormInstance>()
const form = reactive<articlesApi.ArticlesCreateParams>({
  title: '',
  content: '',
  summary: '',
  author_id: 0,
  author_name: '',
  category_id: 0,
  cover_image: '',
  tags: '',
  status: 0,
  is_top: 0,
  view_count: 0,
  like_count: 0,
  comment_count: 0,
  published_at: '',
})

// 表单验证规则
const formRules = reactive<FormRules>({
  title: [
    { required: true, message: '请输入标题', trigger: 'blur' },
    { max: 255, message: '长度不能超过255位', trigger: 'blur' },
  ],
  content: [
    { required: true, message: '请输入内容', trigger: 'blur' },
  ],
  summary: [
    { max: 500, message: '长度不能超过500位', trigger: 'blur' },
  ],
  author_id: [
    { required: true, message: '请输入作者ID', trigger: 'blur' },
  ],
  author_name: [
    { max: 100, message: '长度不能超过100位', trigger: 'blur' },
  ],
  cover_image: [
    { max: 500, message: '长度不能超过500位', trigger: 'blur' },
  ],
  tags: [
    { max: 255, message: '长度不能超过255位', trigger: 'blur' },
  ],
  status: [
    { required: true, message: '请输入状态:0草稿,1已发布,2已下线', trigger: 'blur' },
  ],
  is_top: [
    { required: true, message: '请输入是否置顶', trigger: 'blur' },
  ],
  view_count: [
    { required: true, message: '请输入浏览次数', trigger: 'blur' },
  ],
  like_count: [
    { required: true, message: '请输入点赞次数', trigger: 'blur' },
  ],
  comment_count: [
    { required: true, message: '请输入评论次数', trigger: 'blur' },
  ],
})

// ========== 方法 ==========

// 加载数据
const loadData = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.page_size,
      keyword: searchForm.keyword,
    }
    
    const data = await articlesApi.getArticlesList(params)
    tableData.value = data.list
    pagination.total = data.total
  } catch (error: any) {
    ElMessage.error(error?.message || '加载数据失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  loadData()
}

// 重置
const handleReset = () => {
  searchForm.keyword = ''
  handleSearch()
}

// 分页变化
const handlePageChange = (page: number) => {
  pagination.page = page
  loadData()
}

const handleSizeChange = (size: number) => {
  pagination.page_size = size
  pagination.page = 1
  loadData()
}

// 新建
const handleCreate = () => {
  isEdit.value = false
  dialogTitle.value = '新建Articles管理'
  resetForm()
  dialogVisible.value = true
}

// 编辑
const handleEdit = async (row: articlesApi.Articles) => {
  isEdit.value = true
  currentId.value = row.id
  dialogTitle.value = '编辑Articles管理'
  
  try {
    const data = await articlesApi.getArticles(row.id)
    form.title = data.title
    form.content = data.content
    form.summary = data.summary
    form.author_id = data.author_id
    form.author_name = data.author_name
    form.category_id = data.category_id
    form.cover_image = data.cover_image
    form.tags = data.tags
    form.status = data.status
    form.is_top = data.is_top
    form.view_count = data.view_count
    form.like_count = data.like_count
    form.comment_count = data.comment_count
    form.published_at = data.published_at
    dialogVisible.value = true
  } catch (error: any) {
    ElMessage.error(error?.message || '获取数据失败')
  }
}

// 删除
const handleDelete = async (row: articlesApi.Articles) => {
  try {
    await ElMessageBox.confirm('确定要删除这条记录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await articlesApi.deleteArticles(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error?.message || '删除失败')
    }
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    
    submitLoading.value = true
    try {
      if (isEdit.value && currentId.value) {
        await articlesApi.updateArticles(currentId.value, form)
        ElMessage.success('更新成功')
      } else {
        await articlesApi.createArticles(form)
        ElMessage.success('创建成功')
      }
      
      dialogVisible.value = false
      loadData()
    } catch (error: any) {
      ElMessage.error(error?.message || '操作失败')
    } finally {
      submitLoading.value = false
    }
  })
}

// 关闭对话框
const handleDialogClose = () => {
  resetForm()
}

// 重置表单
const resetForm = () => {
  formRef.value?.resetFields()
  form.title = ''
  form.content = ''
  form.summary = ''
  form.author_id = 0
  form.author_name = ''
  form.category_id = 0
  form.cover_image = ''
  form.tags = ''
  form.status = 0
  form.is_top = 0
  form.view_count = 0
  form.like_count = 0
  form.comment_count = 0
  form.published_at = ''
}

// 格式化日期
const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

// ========== 生命周期 ==========

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.articles-container {
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0;
  font-size: 24px;
  font-weight: 500;
}

.search-card {
  margin-bottom: 20px;
}

.search-form {
  margin-bottom: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
