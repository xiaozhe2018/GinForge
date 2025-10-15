package generator

import (
	"bytes"
	"text/template"
)

// renderFrontendFormTemplate 渲染前端表单页面模板
func renderFrontendFormTemplate(data *TemplateData) (string, error) {
	tmpl := `<template>
  <div class="{{toKebabCase .ModelName}}-form-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ isEdit ? '编辑{{.Title}}' : '新建{{.Title}}' }}</span>
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
          <el-select v-model="form.{{toSnakeCase .Name}}" placeholder="请选择{{.Label}}" style="width: 100%">
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
import * as {{.ModelNameCamel}}Api from '@/api/{{toSnakeCase .ModelName}}'

// ========== 路由 ==========

const route = useRoute()
const router = useRouter()

// ========== 数据定义 ==========

const submitLoading = ref(false)
const isEdit = ref(false)
const id = ref<number | null>(null)

const formRef = ref<FormInstance>()
const form = reactive<{{.ModelNameCamel}}Api.{{.ModelName}}CreateParams>({
{{- range .Fields}}
{{- if and .FormVisible (not .AutoIncrement) (not .IsPrimaryKey)}}
  {{toSnakeCase .Name}}: {{getDefaultValue .}},
{{- end}}
{{- end}}
})

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
  if (!id.value) return
  
  try {
    const data = await {{.ModelNameCamel}}Api.get{{.ModelName}}(id.value)
{{- range .Fields}}
{{- if and .FormVisible (not .AutoIncrement)}}
    form.{{toSnakeCase .Name}} = data.{{toSnakeCase .Name}}
{{- end}}
{{- end}}
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
        await {{.ModelNameCamel}}Api.update{{.ModelName}}(id.value, form)
        ElMessage.success('更新成功')
      } else {
        await {{.ModelNameCamel}}Api.create{{.ModelName}}(form)
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
.{{toKebabCase .ModelName}}-form-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
`

	funcMap := template.FuncMap{
		"toSnakeCase":     toSnakeCase,
		"toKebabCase":     toKebabCase,
		"getDefaultValue": getDefaultValue,
		"contains":        contains,
		"hasMinLength":    hasMinLength,
		"hasMaxLength":    hasMaxLength,
		"getMinLength":    getMinLength,
		"getMaxLength":    getMaxLength,
	}

	t, err := template.New("frontend-form").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
