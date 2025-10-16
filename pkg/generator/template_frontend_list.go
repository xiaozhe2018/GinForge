package generator

import (
	"bytes"
	"strings"
	"text/template"
)

// renderFrontendListTemplate 渲染前端列表页面模板
func renderFrontendListTemplate(data *TemplateData) (string, error) {
	tmpl := `<template>
  <div class="{{toKebabCase .ModelName}}-container">
    <div class="page-header">
      <h2>{{.Title}}</h2>
    </div>

    <!-- 搜索栏 -->
    <el-card class="search-card" shadow="never">
      <el-form :inline="true" :model="searchForm" class="search-form">
{{- if .HasSearch}}
        <el-form-item label="搜索">
          <el-input
            v-model="searchForm.keyword"
            placeholder="请输入关键词"
            clearable
            @keyup.enter="handleSearch"
            style="width: 240px"
          />
        </el-form-item>
{{- end}}
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
          <span>{{.Title}}列表</span>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            <span>新建{{.Title}}</span>
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
{{- range .Fields}}
{{- if .ListVisible}}
        <el-table-column
          prop="{{toSnakeCase .Name}}"
          label="{{.Label}}"
{{- if eq .FormType "switch"}}
          width="100"
        >
          <template #default="{ row }">
            <el-tag :type="row.{{toSnakeCase .Name}} ? 'success' : 'info'">
              {{"{{"}} row.{{toSnakeCase .Name}} ? '是' : '否' {{"}}"}}
            </el-tag>
          </template>
        </el-table-column>
{{- else if or (eq .GoType "time.Time") (eq .GoType "*time.Time")}}
          width="180"
        >
          <template #default="{ row }">
            {{"{{"}} formatDate(row.{{toSnakeCase .Name}}) {{"}}"}}
          </template>
        </el-table-column>
{{- else}}
        />
{{- end}}
{{- end}}
{{- end}}
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
{{- if .HasPagination}}
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
{{- end}}
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
{{- range .Fields}}
{{- if .FormVisible}}
        <el-form-item label="{{.Label}}" prop="{{toSnakeCase .Name}}">
{{- if eq .FormType "textarea"}}
          <el-input
            v-model="form.{{toSnakeCase .Name}}"
            type="textarea"
            :rows="4"
            placeholder="请输入{{.Label}}"
          />
{{- else if eq .FormType "switch"}}
          <el-switch v-model="form.{{toSnakeCase .Name}}" />
{{- else if eq .FormType "select"}}
          <el-select v-model="form.{{toSnakeCase .Name}}" placeholder="请选择{{.Label}}">
            <!-- TODO: 添加选项 -->
            <el-option label="选项1" value="1" />
            <el-option label="选项2" value="2" />
          </el-select>
{{- else if eq .FormType "date"}}
          <el-date-picker
            v-model="form.{{toSnakeCase .Name}}"
            type="date"
            placeholder="请选择{{.Label}}"
            style="width: 100%"
          />
{{- else if eq .FormType "datetime"}}
          <el-date-picker
            v-model="form.{{toSnakeCase .Name}}"
            type="datetime"
            placeholder="请选择{{.Label}}"
            style="width: 100%"
          />
{{- else}}
          <el-input
            v-model="form.{{toSnakeCase .Name}}"
{{- if eq .FormType "password"}}
            type="password"
{{- else if eq .FormType "number"}}
            type="number"
{{- end}}
            placeholder="请输入{{.Label}}"
          />
{{- end}}
        </el-form-item>
{{- end}}
{{- end}}
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
import * as {{.ModelNameCamel}}Api from '@/api/{{toSnakeCase .ModelName}}'

// ========== 数据定义 ==========

const loading = ref(false)
const submitLoading = ref(false)
const tableData = ref<{{.ModelNameCamel}}Api.{{.ModelName}}[]>([])

// 搜索表单
const searchForm = reactive({
{{- if .HasSearch}}
  keyword: '',
{{- end}}
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
const form = reactive<{{.ModelNameCamel}}Api.{{.ModelName}}CreateParams>({
{{- range .Fields}}
{{- if and .FormVisible (not .AutoIncrement) (not .IsPrimaryKey)}}
  {{toSnakeCase .Name}}: {{getDefaultValue .}},
{{- end}}
{{- end}}
})

// 表单验证规则
const formRules = reactive<FormRules>({
{{- range .Fields}}
{{- if and .FormVisible .Validations}}
  {{toSnakeCase .Name}}: [
{{- if contains .Validations "required"}}
    { required: true, message: '请输入{{.Label}}', trigger: 'blur' },
{{- end}}
{{- if contains .Validations "email"}}
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' },
{{- end}}
{{- if hasMinLength .Validations}}
    { min: {{getMinLength .Validations}}, message: '长度不能少于{{getMinLength .Validations}}位', trigger: 'blur' },
{{- end}}
{{- if hasMaxLength .Validations}}
    { max: {{getMaxLength .Validations}}, message: '长度不能超过{{getMaxLength .Validations}}位', trigger: 'blur' },
{{- end}}
  ],
{{- end}}
{{- end}}
})

// ========== 方法 ==========

// 加载数据
const loadData = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.page_size,
{{- if .HasSearch}}
      keyword: searchForm.keyword,
{{- end}}
    }
    
    const data = await {{.ModelNameCamel}}Api.get{{.ModelName}}List(params)
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
{{- if .HasSearch}}
  searchForm.keyword = ''
{{- end}}
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
  dialogTitle.value = '新建{{.Title}}'
  resetForm()
  dialogVisible.value = true
}

// 编辑
const handleEdit = async (row: {{.ModelNameCamel}}Api.{{.ModelName}}) => {
  isEdit.value = true
  currentId.value = row.{{getPrimaryKeySnakeName .Fields}}
  dialogTitle.value = '编辑{{.Title}}'
  
  try {
    const data = await {{.ModelNameCamel}}Api.get{{.ModelName}}(row.{{getPrimaryKeySnakeName .Fields}})
{{- range .Fields}}
{{- if and .FormVisible (not .AutoIncrement)}}
    form.{{toSnakeCase .Name}} = data.{{toSnakeCase .Name}}
{{- end}}
{{- end}}
    dialogVisible.value = true
  } catch (error: any) {
    ElMessage.error(error?.message || '获取数据失败')
  }
}

// 删除
const handleDelete = async (row: {{.ModelNameCamel}}Api.{{.ModelName}}) => {
  try {
    await ElMessageBox.confirm('确定要删除这条记录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await {{.ModelNameCamel}}Api.delete{{.ModelName}}(row.{{getPrimaryKeySnakeName .Fields}})
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
        await {{.ModelNameCamel}}Api.update{{.ModelName}}(currentId.value, form)
        ElMessage.success('更新成功')
      } else {
        await {{.ModelNameCamel}}Api.create{{.ModelName}}(form)
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
{{- range .Fields}}
{{- if and .FormVisible (not .AutoIncrement) (not .IsPrimaryKey)}}
  form.{{toSnakeCase .Name}} = {{getDefaultValue .}}
{{- end}}
{{- end}}
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
.{{toKebabCase .ModelName}}-container {
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
`

	funcMap := template.FuncMap{
		"toSnakeCase":            toSnakeCase,
		"toKebabCase":            toKebabCase,
		"getPrimaryKeySnakeName": getPrimaryKeySnakeName,
		"getDefaultValue":        getDefaultValue,
		"contains":               contains,
		"hasMinLength":           hasMinLength,
		"hasMaxLength":           hasMaxLength,
		"getMinLength":           getMinLength,
		"getMaxLength":           getMaxLength,
	}

	t, err := template.New("frontend-list").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// getPrimaryKeySnakeName 获取主键的 snake_case 名称
func getPrimaryKeySnakeName(fields []FieldConfig) string {
	pk := getPrimaryKeyField(fields)
	if pk != nil {
		return toSnakeCase(pk.Name)
	}
	return "id"
}

// getDefaultValue 获取默认值
func getDefaultValue(field FieldConfig) string {
	switch field.TSType {
	case "number":
		return "0"
	case "boolean":
		return "false"
	case "string":
		return "''"
	default:
		if field.Nullable {
			return "undefined"
		}
		return "''"
	}
}

// hasMinLength 检查是否有最小长度验证
func hasMinLength(validations []string) bool {
	for _, v := range validations {
		if strings.HasPrefix(v, "min:") {
			return true
		}
	}
	return false
}

// hasMaxLength 检查是否有最大长度验证
func hasMaxLength(validations []string) bool {
	for _, v := range validations {
		if strings.HasPrefix(v, "max:") {
			return true
		}
	}
	return false
}

// getMinLength 获取最小长度
func getMinLength(validations []string) string {
	for _, v := range validations {
		if strings.HasPrefix(v, "min:") {
			return strings.TrimPrefix(v, "min:")
		}
	}
	return "0"
}

// getMaxLength 获取最大长度
func getMaxLength(validations []string) string {
	for _, v := range validations {
		if strings.HasPrefix(v, "max:") {
			return strings.TrimPrefix(v, "max:")
		}
	}
	return "255"
}
