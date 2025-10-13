<template>
  <div class="roles-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <h2>角色管理</h2>
      <p>管理系统角色和权限分配</p>
    </div>

    <!-- 搜索和操作栏 -->
    <div class="toolbar">
      <div class="search-form">
        <el-input
          v-model="searchForm.keyword"
          placeholder="搜索角色名称"
          style="width: 300px"
          clearable
          @keyup.enter="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        
        <el-select
          v-model="searchForm.status"
          placeholder="角色状态"
          style="width: 120px; margin-left: 10px"
          clearable
        >
          <el-option label="全部" value="" />
          <el-option label="启用" value="active" />
          <el-option label="禁用" value="disabled" />
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
          创建角色
        </el-button>
      </div>
    </div>

    <!-- 角色列表 -->
    <el-card class="table-card">
      <el-table
        :data="roleList"
        v-loading="loading"
        stripe
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="角色名称" width="150" />
        <el-table-column prop="code" label="角色代码" width="150" />
        <el-table-column prop="description" label="描述" min-width="200" />
        <el-table-column prop="userCount" label="用户数量" width="100">
          <template #default="{ row }">
            <el-tag type="info">{{ row.userCount || 0 }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="'active'"
              :inactive-value="'disabled'"
              @change="handleStatusChange(row)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="250" fixed="right">
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
              @click="handlePermissions(row)"
            >
              权限
            </el-button>
            <el-button
              type="danger"
              size="small"
              @click="handleDelete(row)"
              :disabled="row.code === 'super_admin'"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination">
        <el-pagination
          :current-page="pagination.current"
          :page-size="pagination.size"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 角色表单对话框 -->
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
        <el-form-item label="角色名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入角色名称" />
        </el-form-item>
        
        <el-form-item label="角色代码" prop="code">
          <el-input v-model="form.code" placeholder="请输入角色代码" />
        </el-form-item>
        
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            placeholder="请输入角色描述"
          />
        </el-form-item>
        
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio label="active">启用</el-radio>
            <el-radio label="disabled">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">
          确定
        </el-button>
      </template>
    </el-dialog>

    <!-- 权限分配对话框 -->
    <el-dialog
      v-model="permissionDialogVisible"
      title="权限分配"
      width="800px"
      @close="handlePermissionDialogClose"
    >
      <div class="permission-tree">
        <el-tree
          ref="permissionTreeRef"
          :data="permissionTree"
          :props="treeProps"
          show-checkbox
          node-key="id"
          :default-checked-keys="checkedPermissions"
          check-strictly
        />
      </div>
      
      <template #footer>
        <el-button @click="permissionDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handlePermissionSubmit" :loading="permissionSubmitting">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Plus } from '@element-plus/icons-vue'
import * as roleApi from '@/api/role'
import * as menuApi from '@/api/menu'
import * as permissionApi from '@/api/permission'

// 搜索表单
const searchForm = reactive({
  keyword: '',
  status: ''
})

// 分页信息
const pagination = reactive({
  current: 1,
  size: 20,
  total: 0
})

// 角色列表
const roleList = ref<any[]>([])
const loading = ref(false)

// 对话框
const dialogVisible = ref(false)
const dialogTitle = ref('')
const submitting = ref(false)

// 权限对话框
const permissionDialogVisible = ref(false)
const permissionSubmitting = ref(false)
const currentRole = ref(null)

// 表单
const formRef = ref()
const form = reactive({
  id: null,
  name: '',
  code: '',
  description: '',
  status: 'active'
})

// 权限树
const permissionTreeRef = ref()
const permissionTree = ref<any[]>([])
const checkedPermissions = ref<any[]>([])

// 树形组件配置
const treeProps = {
  children: 'children',
  label: 'name'
}

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入角色名称', trigger: 'blur' },
    { min: 2, max: 20, message: '角色名称长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入角色代码', trigger: 'blur' },
    { pattern: /^[a-zA-Z_][a-zA-Z0-9_]*$/, message: '角色代码只能包含字母、数字和下划线，且不能以数字开头', trigger: 'blur' }
  ],
  description: [
    { max: 200, message: '描述不能超过 200 个字符', trigger: 'blur' }
  ]
}

// 获取角色列表
const fetchRoleList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.current,
      page_size: pagination.size,
      keyword: searchForm.keyword || undefined,
      status: searchForm.status === 'active' ? 1 : searchForm.status === 'disabled' ? 0 : undefined
    }
    const response: any = await roleApi.getRoleList(params)
    roleList.value = (response.list || []).map((role: any) => ({
      ...role,
      status: role.status === 1 ? 'active' : 'disabled',
      userCount: 0,  // TODO: 后端需要返回用户数量
      createdAt: role.created_at
    }))
    pagination.total = response.total || 0
  } catch (error) {
    console.error('获取角色列表失败:', error)
    ElMessage.error('获取角色列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.current = 1
  fetchRoleList()
}

// 重置搜索
const handleReset = () => {
  Object.assign(searchForm, {
    keyword: '',
    status: ''
  })
  pagination.current = 1
  fetchRoleList()
}

// 创建角色
const handleCreate = () => {
  dialogTitle.value = '创建角色'
  Object.assign(form, {
    id: null,
    name: '',
    code: '',
    description: '',
    status: 'active'
  })
  dialogVisible.value = true
}

// 编辑角色
const handleEdit = (row: any) => {
  dialogTitle.value = '编辑角色'
  Object.assign(form, {
    id: row.id,
    name: row.name,
    code: row.code,
    description: row.description,
    status: row.status
  })
  dialogVisible.value = true
}

// 删除角色
const handleDelete = async (row: any) => {
  if (row.code === 'super_admin') {
    ElMessage.warning('超级管理员角色不能删除')
    return
  }
  
  try {
    await ElMessageBox.confirm('确定要删除该角色吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await roleApi.deleteRole(row.id)
    ElMessage.success('删除成功')
    fetchRoleList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 状态变更
const handleStatusChange = async (row: any) => {
  try {
    const statusNum = row.status === 'active' ? 1 : 0
    const updateData = {
      name: row.name,
      code: row.code,
      description: row.description || '',
      sort: row.sort || 0,
      status: statusNum,
      permission_ids: [],
      menu_ids: []
    }
    await roleApi.updateRole(row.id, updateData)
    ElMessage.success('状态更新成功')
  } catch (error) {
    console.error('状态更新失败:', error)
    ElMessage.error('状态更新失败')
    // 恢复原状态
    row.status = row.status === 'active' ? 'disabled' : 'active'
  }
}

// 权限管理
const handlePermissions = async (row: any) => {
  currentRole.value = row
  permissionDialogVisible.value = true
  
  try {
    // 获取菜单树作为权限树
    const menuTreeResponse: any = await menuApi.getMenuTree()
    permissionTree.value = menuTreeResponse.list || []
    
    // 获取角色详情（包含已分配的权限和菜单）
    const roleDetail: any = await roleApi.getRoleDetail(row.id)
    
    // 使用菜单ID作为已选中的权限
    if (roleDetail.menus && Array.isArray(roleDetail.menus)) {
      checkedPermissions.value = roleDetail.menus.map((m: any) => m.id)
    } else {
      checkedPermissions.value = []
    }
  } catch (error) {
    console.error('获取权限信息失败:', error)
    ElMessage.error('获取权限信息失败')
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    submitting.value = true
    
    const roleData = {
      name: form.name,
      code: form.code,
      description: form.description || '',
      sort: 0,
      status: form.status === 'active' ? 1 : 0,
      permission_ids: [],
      menu_ids: []
    }
    
    if (form.id) {
      await roleApi.updateRole(form.id, roleData)
      ElMessage.success('更新成功')
    } else {
      await roleApi.createRole(roleData)
      ElMessage.success('创建成功')
    }
    
    dialogVisible.value = false
    fetchRoleList()
  } catch (error) {
    console.error('保存失败:', error)
    ElMessage.error(form.id ? '更新失败' : '创建失败')
  } finally {
    submitting.value = false
  }
}

// 提交权限
const handlePermissionSubmit = async () => {
  const role = currentRole.value
  if (!role || !role.id) return
  
  try {
    permissionSubmitting.value = true
    const checkedKeys = permissionTreeRef.value.getCheckedKeys()
    const halfCheckedKeys = permissionTreeRef.value.getHalfCheckedKeys()
    const allMenuIds = [...checkedKeys, ...halfCheckedKeys]
    
    // 获取所有权限ID
    const permissionsResponse: any = await permissionApi.getPermissionList({})
    const allPermissionIds = permissionsResponse.list ? permissionsResponse.list.map((p: any) => p.id) : []
    
    // 更新角色权限和菜单
    const updateData = {
      name: role.name,
      code: role.code,
      description: role.description || '',
      sort: role.sort || 0,
      status: role.status === 'active' ? 1 : 0,
      permission_ids: allPermissionIds,
      menu_ids: allMenuIds
    }
    
    await roleApi.updateRole(role.id, updateData)
    ElMessage.success('权限分配成功')
    permissionDialogVisible.value = false
    fetchRoleList()
  } catch (error) {
    console.error('权限分配失败:', error)
    ElMessage.error('权限分配失败')
  } finally {
    permissionSubmitting.value = false
  }
}

// 关闭对话框
const handleDialogClose = () => {
  formRef.value?.resetFields()
}

// 关闭权限对话框
const handlePermissionDialogClose = () => {
  currentRole.value = null
  checkedPermissions.value = []
}

// 分页大小变更
const handleSizeChange = (size: any) => {
  pagination.size = size
  pagination.current = 1
  fetchRoleList()
}

// 当前页变更
const handleCurrentChange = (page: any) => {
  pagination.current = page
  fetchRoleList()
}

// 格式化日期
const formatDate = (date: any) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}


// 组件挂载时获取数据
onMounted(() => {
  fetchRoleList()
})
</script>

<style scoped>
.roles-page {
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

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.permission-tree {
  max-height: 400px;
  overflow-y: auto;
  border: 1px solid #e6e6e6;
  border-radius: 4px;
  padding: 10px;
}
</style>
