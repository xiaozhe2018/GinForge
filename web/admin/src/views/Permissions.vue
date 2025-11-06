<template>
  <div class="permissions-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <h2>权限管理</h2>
      <p>管理系统权限和资源访问控制</p>
    </div>

    <!-- 搜索和操作栏 -->
    <div class="toolbar">
      <div class="search-form">
        <el-input
          v-model="searchForm.keyword"
          placeholder="搜索权限名称或标识"
          style="width: 300px"
          clearable
          @keyup.enter="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        
        <el-select
          v-model="searchForm.type"
          placeholder="权限类型"
          style="width: 120px; margin-left: 10px"
          clearable
        >
          <el-option label="全部" value="" />
          <el-option label="菜单权限" value="menu" />
          <el-option label="按钮权限" value="button" />
          <el-option label="接口权限" value="api" />
        </el-select>
        
        <el-button type="primary" @click="handleSearch" style="margin-left: 10px">
          <el-icon><Search /></el-icon>
          搜索
        </el-button>
        
        <el-button @click="handleReset">
          <el-icon><Refresh /></el-icon>
          重置
        </el-button>
      </div>
      
      <div class="actions">
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          创建权限
        </el-button>
        <el-button @click="handleExpandAll">
          <el-icon><ArrowDown /></el-icon>
          展开全部
        </el-button>
        <el-button @click="handleCollapseAll">
          <el-icon><ArrowUp /></el-icon>
          收起全部
        </el-button>
      </div>
    </div>

    <!-- 权限树形表格 -->
    <el-card class="table-card">
      <el-table
        :data="permissionList"
        v-loading="loading"
        row-key="id"
        :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
        :default-expand-all="false"
        ref="tableRef"
      >
        <el-table-column prop="name" label="权限名称" width="200">
          <template #default="{ row }">
            <el-icon v-if="row.icon" style="margin-right: 5px">
              <component :is="row.icon" />
            </el-icon>
            {{ row.name }}
          </template>
        </el-table-column>
        
        <el-table-column prop="code" label="权限标识" width="200" />
        
        <el-table-column prop="type" label="权限类型" width="120">
          <template #default="{ row }">
            <el-tag :type="getTypeTagType(row.type)">
              {{ getTypeName(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="resource" label="资源路径" width="200" />
        
        <el-table-column prop="method" label="请求方法" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.method" :type="getMethodTagType(row.method)">
              {{ row.method.toUpperCase() }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="sort" label="排序" width="80" />
        
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="1"
              :inactive-value="0"
              @change="handleStatusChange(row)"
            />
          </template>
        </el-table-column>
        
        <el-table-column prop="createdAt" label="创建时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button
              type="primary"
              size="small"
              @click="handleEdit(row)"
            >
              编辑
            </el-button>
            <el-button
              type="success"
              size="small"
              @click="handleAddChild(row)"
            >
              添加子权限
            </el-button>
            <el-button
              type="danger"
              size="small"
              @click="handleDelete(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 权限表单对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
      @close="handleDialogClose"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="父级权限" prop="parentId">
          <el-tree-select
            v-model="form.parentId"
            :data="permissionTreeOptions"
            :props="treeSelectProps"
            placeholder="请选择父级权限"
            clearable
            check-strictly
          />
        </el-form-item>
        
        <el-form-item label="权限名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入权限名称" />
        </el-form-item>
        
        <el-form-item label="权限标识" prop="code">
          <el-input v-model="form.code" placeholder="请输入权限标识" />
        </el-form-item>
        
        <el-form-item label="权限类型" prop="type">
          <el-select v-model="form.type" placeholder="请选择权限类型" style="width: 100%">
            <el-option label="菜单权限" value="menu" />
            <el-option label="按钮权限" value="button" />
            <el-option label="接口权限" value="api" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="资源路径" prop="resource">
          <el-input v-model="form.resource" placeholder="请输入资源路径" />
        </el-form-item>
        
        <el-form-item label="请求方法" prop="method">
          <el-select v-model="form.method" placeholder="请选择请求方法" style="width: 100%">
            <el-option label="GET" value="get" />
            <el-option label="POST" value="post" />
            <el-option label="PUT" value="put" />
            <el-option label="DELETE" value="delete" />
            <el-option label="PATCH" value="patch" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="form.sort" :min="0" :max="999" />
        </el-form-item>
        
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            placeholder="请输入权限描述"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Plus, ArrowDown, ArrowUp } from '@element-plus/icons-vue'
import * as permissionApi from '@/api/permission'

// 搜索表单
const searchForm = reactive({
  keyword: '',
  type: ''
})

// 权限列表
const permissionList = ref<any[]>([])
const loading = ref(false)
const tableRef = ref()

// 对话框
const dialogVisible = ref(false)
const dialogTitle = ref('')
const submitting = ref(false)

// 表单
const formRef = ref()
const form = reactive({
  id: null,
  parentId: null,
  name: '',
  code: '',
  type: '',
  resource: '',
  method: '',
  sort: 0,
  status: 1,  // 默认启用
  description: ''
})

// 树形选择器配置
const treeSelectProps = {
  value: 'id',
  label: 'name',
  children: 'children'
}

// 权限树选项（用于父级权限选择）
const permissionTreeOptions = computed(() => {
  const buildTree = (permissions: any[], parentId: any = null): any[] => {
    return permissions
      .filter(permission => permission.parentId === parentId)
      .map(permission => ({
        ...permission,
        children: buildTree(permissions, permission.id)
      }))
  }
  return buildTree(permissionList.value)
})

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入权限名称', trigger: 'blur' },
    { min: 2, max: 20, message: '权限名称长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入权限标识', trigger: 'blur' },
    { pattern: /^[a-zA-Z_][a-zA-Z0-9_:]*$/, message: '权限标识只能包含字母、数字、下划线和冒号，且不能以数字开头', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择权限类型', trigger: 'change' }
  ],
  sort: [
    { required: true, message: '请输入排序值', trigger: 'blur' }
  ]
}

// 获取权限列表
const fetchPermissionList = async () => {
  loading.value = true
  try {
    // 构建请求参数，不传递的参数会使用后端默认值
    const params: any = {}
    if (searchForm.keyword) params.keyword = searchForm.keyword
    if (searchForm.type) params.type = searchForm.type
    
    const response: any = await permissionApi.getPermissionList(params)
    
    // 转换后端数据格式
    permissionList.value = (response.list || []).map((perm: any) => ({
      id: perm.id,
      parentId: null,  // 权限列表是扁平的，没有父子关系
      name: perm.name,
      code: perm.code,
      type: perm.type,
      resource: '',  // 后端没有resource字段
      method: '',    // 后端没有method字段
      sort: 0,
      status: perm.status !== undefined ? Number(perm.status) : 1,  // 从后端读取status，默认为1（启用）
      description: perm.description || '',
      createdAt: perm.created_at
    }))
  } catch (error) {
    console.error('获取权限列表失败:', error)
    ElMessage.error('获取权限列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  fetchPermissionList()
}

// 重置搜索
const handleReset = () => {
  Object.assign(searchForm, {
    keyword: '',
    type: ''
  })
  fetchPermissionList()
}

// 创建权限
const handleCreate = () => {
  dialogTitle.value = '创建权限'
  Object.assign(form, {
    id: null,
    parentId: null,
    name: '',
    code: '',
    type: '',
    resource: '',
    method: '',
    sort: 0,
    status: 1,  // 默认启用（数字类型）
    description: ''
  })
  dialogVisible.value = true
}

// 添加子权限
const handleAddChild = (row: any) => {
  dialogTitle.value = '添加子权限'
  Object.assign(form, {
    id: null,
    parentId: row.id,
    name: '',
    code: '',
    type: '',
    resource: '',
    method: '',
    sort: 0,
    status: 1,  // 默认启用（数字类型）
    description: ''
  })
  dialogVisible.value = true
}

// 编辑权限
const handleEdit = (row: any) => {
  dialogTitle.value = '编辑权限'
  Object.assign(form, {
    id: row.id,
    parentId: row.parentId,
    name: row.name,
    code: row.code,
    type: row.type,
    resource: row.resource,
    method: row.method,
    sort: row.sort,
    status: row.status,
    description: row.description
  })
  dialogVisible.value = true
}

// 删除权限
const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定要删除该权限吗？删除后子权限也会被删除！', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await permissionApi.deletePermission(row.id)
    ElMessage.success('删除成功')
    fetchPermissionList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 状态变更
const handleStatusChange = async (row: any) => {
  const oldStatus = row.status === 1 ? 0 : 1 // 记录原状态（因为 v-model 已经改变了）
  
  try {
    await permissionApi.updatePermissionStatus(row.id, row.status)
    ElMessage.success(`权限已${row.status === 1 ? '启用' : '禁用'}`)
  } catch (error: any) {
    console.error('更新权限状态失败:', error)
    ElMessage.error(error.message || '更新状态失败')
    // 恢复原状态
    row.status = oldStatus
  }
}

// 展开全部
const handleExpandAll = () => {
  const allKeys = getAllKeys(permissionList.value)
  allKeys.forEach(key => {
    tableRef.value.toggleRowExpansion(permissionList.value.find(item => item.id === key), true)
  })
}

// 收起全部
const handleCollapseAll = () => {
  const allKeys = getAllKeys(permissionList.value)
  allKeys.forEach(key => {
    tableRef.value.toggleRowExpansion(permissionList.value.find(item => item.id === key), false)
  })
}

// 获取所有节点ID
const getAllKeys = (permissions: any[]): any[] => {
  let keys: any[] = []
  permissions.forEach(permission => {
    keys.push(permission.id)
    if (permission.children && permission.children.length > 0) {
      keys = keys.concat(getAllKeys(permission.children))
    }
  })
  return keys
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    submitting.value = true
    
    // 准备后端需要的数据格式
    const permData = {
      name: form.name,
      code: form.code,
      type: form.type,
      description: form.description
    }
    
    if (form.id) {
      await permissionApi.updatePermission(form.id, permData)
      ElMessage.success('更新成功')
    } else {
      await permissionApi.createPermission(permData)
      ElMessage.success('创建成功')
    }
    
    dialogVisible.value = false
    fetchPermissionList()
  } catch (error) {
    console.error('保存失败:', error)
    ElMessage.error(form.id ? '更新失败' : '创建失败')
  } finally {
    submitting.value = false
  }
}

// 关闭对话框
const handleDialogClose = () => {
  formRef.value?.resetFields()
}

// 获取类型标签类型
const getTypeTagType = (type: string) => {
  const typeMap: Record<string, string> = {
    menu: 'primary',
    button: 'success',
    api: 'warning'
  }
  return typeMap[type] || 'info'
}

// 获取类型名称
const getTypeName = (type: string) => {
  const nameMap: Record<string, string> = {
    menu: '菜单权限',
    button: '按钮权限',
    api: '接口权限'
  }
  return nameMap[type] || type
}

// 获取方法标签类型
const getMethodTagType = (method: string) => {
  const typeMap: Record<string, string> = {
    get: 'success',
    post: 'primary',
    put: 'warning',
    delete: 'danger',
    patch: 'info'
  }
  return typeMap[method.toLowerCase()] || 'info'
}

// 格式化日期
const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

// 组件挂载时获取数据
onMounted(() => {
  fetchPermissionList()
})
</script>

<style scoped>
.permissions-page {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
  padding-bottom: 20px;
  border-bottom: 1px solid #e6e6e6;
}

.page-header h2 {
  margin: 0 0 8px 0;
  font-size: 20px;
  color: #333;
}

.page-header p {
  margin: 0;
  color: #666;
  font-size: 14px;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.search-form {
  display: flex;
  align-items: center;
}

.actions {
  display: flex;
  gap: 10px;
}

.table-card {
  margin-bottom: 20px;
}
</style>
