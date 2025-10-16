<template>
  <div class="articles-form-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ isEdit ? '编辑Articles管理' : '新建Articles管理' }}</span>
          <el-button @click="handleBack">返回</el-button>
        </div>
      </template>

      <el-form
        ref="formRef"
        :model="form"
        :rules="formRules"
        label-width="120px"
        style="max-width: 600px"
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
          <el-select v-model="form.category_id" placeholder="请选择分类ID" style="width: 100%">
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

        <el-form-item>
          <el-button type="primary" :loading="submitLoading" @click="handleSubmit">
            保存
          </el-button>
          <el-button @click="handleBack">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, FormInstance, FormRules } from 'element-plus'
import * as articlesApi from '@/api/articles'

// ========== 路由 ==========

const route = useRoute()
const router = useRouter()

// ========== 数据定义 ==========

const submitLoading = ref(false)
const isEdit = ref(false)
const id = ref<number | null>(null)

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
  if (!id.value) return
  
  try {
    const data = await articlesApi.getArticles(id.value)
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
  } catch (error: any) {
    ElMessage.error(error?.message || '加载数据失败')
    handleBack()
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    
    submitLoading.value = true
    try {
      if (isEdit.value && id.value) {
        await articlesApi.updateArticles(id.value, form)
        ElMessage.success('更新成功')
      } else {
        await articlesApi.createArticles(form)
        ElMessage.success('创建成功')
      }
      
      handleBack()
    } catch (error: any) {
      ElMessage.error(error?.message || '操作失败')
    } finally {
      submitLoading.value = false
    }
  })
}

// 返回列表
const handleBack = () => {
  router.back()
}

// ========== 生命周期 ==========

onMounted(() => {
  const routeId = route.params.id
  if (routeId && routeId !== 'create') {
    isEdit.value = true
    id.value = Number(routeId)
    loadData()
  }
})
</script>

<style scoped>
.articles-form-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
