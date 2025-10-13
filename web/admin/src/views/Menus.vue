<template>
  <div class="menus-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <h2>菜单管理</h2>
      <p>管理系统菜单结构和权限配置</p>
    </div>

    <!-- 搜索和操作栏 -->
    <div class="toolbar">
      <div class="search-form">
        <el-input
          v-model="searchForm.keyword"
          placeholder="搜索菜单名称"
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
          placeholder="菜单状态"
          style="width: 120px; margin-left: 10px"
          clearable
        >
          <el-option label="全部" value="" />
          <el-option label="显示" value="show" />
          <el-option label="隐藏" value="hide" />
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
          创建菜单
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

    <!-- 菜单树形表格 -->
    <el-card class="table-card">
      <el-table
        :data="menuList"
        v-loading="loading"
        row-key="id"
        :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
        :default-expand-all="false"
        ref="tableRef"
      >
        <el-table-column prop="name" label="菜单名称" width="200">
          <template #default="{ row }">
            <el-icon v-if="row.icon" style="margin-right: 5px">
              <component :is="row.icon" />
            </el-icon>
            {{ row.name }}
          </template>
        </el-table-column>
        
        <el-table-column prop="path" label="路由路径" width="200" />
        
        <el-table-column prop="component" label="组件路径" width="200" />
        
        <el-table-column prop="permission" label="权限标识" width="150" />
        
        <el-table-column prop="sort" label="排序" width="80" />
        
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="'show'"
              :inactive-value="'hide'"
              @change="handleStatusChange(row)"
            />
          </template>
        </el-table-column>
        
        <el-table-column prop="createdAt" label="创建时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="200" fixed="right">
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
              添加子菜单
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

    <!-- 菜单表单对话框 -->
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
        <el-form-item label="父级菜单" prop="parentId">
          <el-tree-select
            v-model="form.parentId"
            :data="menuTreeOptions"
            :props="treeSelectProps"
            placeholder="请选择父级菜单"
            clearable
            check-strictly
          />
        </el-form-item>
        
        <el-form-item label="菜单名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入菜单名称" />
        </el-form-item>
        
        <el-form-item label="菜单图标" prop="icon">
          <el-input v-model="form.icon" placeholder="请输入图标名称" />
        </el-form-item>
        
        <el-form-item label="路由路径" prop="path">
          <el-input v-model="form.path" placeholder="请输入路由路径" />
        </el-form-item>
        
        <el-form-item label="组件路径" prop="component">
          <el-input v-model="form.component" placeholder="请输入组件路径" />
        </el-form-item>
        
        <el-form-item label="权限标识" prop="permission">
          <el-input v-model="form.permission" placeholder="请输入权限标识" />
        </el-form-item>
        
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="form.sort" :min="0" :max="999" />
        </el-form-item>
        
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio label="show">显示</el-radio>
            <el-radio label="hide">隐藏</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            placeholder="请输入菜单描述"
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
import * as menuApi from '@/api/menu'

// 搜索表单
const searchForm = reactive({
  keyword: '',
  status: ''
})

// 菜单列表
const menuList = ref<any[]>([])
const loading = ref(false)
const tableRef = ref()

// 对话框
const dialogVisible = ref(false)
const dialogTitle = ref('')
const submitting = ref(false)

// 表单
const formRef = ref()
const form = reactive({
  id: null as number | null,
  parentId: null as number | null,
  name: '',
  icon: '',
  path: '',
  component: '',
  permission: '',
  sort: 0,
  status: 'show' as 'show' | 'hide',
  description: ''
})

// 树形选择器配置
const treeSelectProps = {
  value: 'id',
  label: 'name',
  children: 'children'
}

// 菜单树选项（用于父级菜单选择）
const menuTreeOptions = computed((): any[] => {
  const buildTree = (menus: any[], parentId: any = null): any[] => {
    return menus
      .filter(menu => menu.parentId === parentId)
      .map(menu => ({
        ...menu,
        children: buildTree(menus, menu.id)
      }))
  }
  return buildTree(menuList.value)
})

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入菜单名称', trigger: 'blur' },
    { min: 2, max: 20, message: '菜单名称长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  path: [
    { required: true, message: '请输入路由路径', trigger: 'blur' },
    { pattern: /^\/[a-zA-Z0-9\/\-_]*$/, message: '路由路径必须以/开头，只能包含字母、数字、/、-、_', trigger: 'blur' }
  ],
  component: [
    { required: true, message: '请输入组件路径', trigger: 'blur' }
  ],
  sort: [
    { required: true, message: '请输入排序值', trigger: 'blur' }
  ]
}

// 获取菜单列表
const fetchMenuList = async () => {
  loading.value = true
  try {
    // 使用菜单树接口
    const response: any = await menuApi.getMenuTree()
    
    // 转换后端数据格式为前端需要的格式
    const convertMenu = (menu: any): any => ({
      id: menu.id,
      parentId: menu.parent_id,
      name: menu.name,
      icon: menu.icon,
      path: menu.path,
      component: menu.component,
      permission: menu.code,  // 使用code作为permission标识
      sort: menu.sort,
      status: menu.status === 1 ? 'show' : 'hide',
      visible: menu.visible,
      description: menu.description,
      createdAt: menu.created_at,
      children: menu.children ? menu.children.map(convertMenu) : []
    })
    
    menuList.value = (response.list || []).map(convertMenu)
  } catch (error) {
    console.error('获取菜单列表失败:', error)
    ElMessage.error('获取菜单列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  fetchMenuList()
}

// 重置搜索
const handleReset = () => {
  Object.assign(searchForm, {
    keyword: '',
    status: ''
  })
  fetchMenuList()
}

// 创建菜单
const handleCreate = () => {
  dialogTitle.value = '创建菜单'
  Object.assign(form, {
    id: null,
    parentId: null,
    name: '',
    icon: '',
    path: '',
    component: '',
    permission: '',
    sort: 0,
    status: 'show',
    description: ''
  })
  dialogVisible.value = true
}

// 添加子菜单
const handleAddChild = (row: any) => {
  dialogTitle.value = '添加子菜单'
  Object.assign(form, {
    id: null,
    parentId: row.id,
    name: '',
    icon: '',
    path: '',
    component: '',
    permission: '',
    sort: 0,
    status: 'show',
    description: ''
  })
  dialogVisible.value = true
}

// 编辑菜单
const handleEdit = (row: any) => {
  dialogTitle.value = '编辑菜单'
  Object.assign(form, {
    id: row.id,
    parentId: row.parentId,
    name: row.name,
    icon: row.icon,
    path: row.path,
    component: row.component,
    permission: row.permission,
    sort: row.sort,
    status: row.status,
    description: row.description
  })
  dialogVisible.value = true
}

// 删除菜单
const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定要删除该菜单吗？删除后子菜单也会被删除！', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await menuApi.deleteMenu(row.id)
    ElMessage.success('删除成功')
    fetchMenuList()
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
    // 准备完整的菜单数据
    const menuData = {
      parent_id: row.parentId || 0,
      name: row.name,
      code: row.permission || row.name.toLowerCase(),
      type: 'menu',
      path: row.path,
      component: row.component,
      icon: row.icon,
      sort: row.sort,
      visible: row.status === 'show' ? 1 : 0,
      status: row.status === 'show' ? 1 : 0,
      description: row.description || ''
    }
    
    await menuApi.updateMenu(row.id, menuData)
    ElMessage.success('状态更新成功')
  } catch (error) {
    console.error('状态更新失败:', error)
    ElMessage.error('状态更新失败')
    // 恢复原状态
    row.status = row.status === 'show' ? 'hide' : 'show'
  }
}

// 展开全部
const handleExpandAll = () => {
  const allKeys = getAllKeys(menuList.value)
  allKeys.forEach(key => {
    tableRef.value.toggleRowExpansion(menuList.value.find(item => item.id === key), true)
  })
}

// 收起全部
const handleCollapseAll = () => {
  const allKeys = getAllKeys(menuList.value)
  allKeys.forEach(key => {
    tableRef.value.toggleRowExpansion(menuList.value.find(item => item.id === key), false)
  })
}

// 获取所有节点ID
const getAllKeys = (menus: any[]): any[] => {
  let keys: any[] = []
  menus.forEach((menu: any) => {
    keys.push(menu.id)
    if (menu.children && menu.children.length > 0) {
      keys = keys.concat(getAllKeys(menu.children))
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
    const menuData = {
      parent_id: form.parentId || 0,
      name: form.name,
      code: form.permission || form.name.toLowerCase(),  // 使用permission作为code
      type: 'menu',  // 默认为menu类型
      path: form.path,
      component: form.component,
      icon: form.icon,
      sort: form.sort,
      visible: form.status === 'show' ? 1 : 0,
      status: form.status === 'show' ? 1 : 0,
      description: form.description
    }
    
    if (form.id) {
      await menuApi.updateMenu(form.id, menuData)
      ElMessage.success('更新成功')
    } else {
      await menuApi.createMenu(menuData)
      ElMessage.success('创建成功')
    }
    
    dialogVisible.value = false
    fetchMenuList()
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

// 格式化日期
const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

// 组件挂载时获取数据
onMounted(() => {
  fetchMenuList()
})
</script>

<style scoped>
.menus-page {
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
