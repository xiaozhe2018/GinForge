<template>
  <div class="role-form">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ isEdit ? '编辑角色' : '创建角色' }}</span>
          <el-button @click="goBack">返回</el-button>
        </div>
      </template>
      
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
        @submit.prevent="handleSubmit"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="角色名称" prop="name">
              <el-input v-model="form.name" placeholder="请输入角色名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="角色编码" prop="code">
              <el-input v-model="form.code" placeholder="请输入角色编码" />
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="排序" prop="sort">
              <el-input-number v-model="form.sort" :min="0" :max="999" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态" prop="status">
              <el-radio-group v-model="form.status">
                <el-radio :label="1">启用</el-radio>
                <el-radio :label="0">禁用</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-form-item label="角色描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            placeholder="请输入角色描述"
          />
        </el-form-item>
        
        <el-form-item label="权限配置" prop="permissions">
          <el-tree
            ref="treeRef"
            :data="permissionTree"
            :props="treeProps"
            show-checkbox
            node-key="id"
            :default-checked-keys="form.permissions"
            @check="handlePermissionCheck"
          />
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="handleSubmit" :loading="loading">
            {{ isEdit ? '更新' : '创建' }}
          </el-button>
          <el-button @click="goBack">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'

const router = useRouter()
const route = useRoute()

const formRef = ref()
const treeRef = ref()
const loading = ref(false)

const isEdit = computed(() => route.name === 'RoleEdit')

const form = reactive({
  name: '',
  code: '',
  sort: 0,
  status: 1,
  description: '',
  permissions: []
})

const permissionTree = ref([
  {
    id: 1,
    label: '用户管理',
    children: [
      { id: 11, label: '用户查看' },
      { id: 12, label: '用户创建' },
      { id: 13, label: '用户编辑' },
      { id: 14, label: '用户删除' }
    ]
  },
  {
    id: 2,
    label: '角色管理',
    children: [
      { id: 21, label: '角色查看' },
      { id: 22, label: '角色创建' },
      { id: 23, label: '角色编辑' },
      { id: 24, label: '角色删除' }
    ]
  },
  {
    id: 3,
    label: '菜单管理',
    children: [
      { id: 31, label: '菜单查看' },
      { id: 32, label: '菜单创建' },
      { id: 33, label: '菜单编辑' },
      { id: 34, label: '菜单删除' }
    ]
  }
])

const treeProps = {
  children: 'children',
  label: 'label'
}

const rules = {
  name: [
    { required: true, message: '请输入角色名称', trigger: 'blur' },
    { min: 2, max: 20, message: '角色名称长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入角色编码', trigger: 'blur' },
    { min: 2, max: 20, message: '角色编码长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  sort: [
    { required: true, message: '请输入排序', trigger: 'blur' }
  ]
}

const handlePermissionCheck = (_data: any, checked: any) => {
  form.permissions = checked.checkedKeys
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    loading.value = true
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    ElMessage.success(isEdit.value ? '角色更新成功' : '角色创建成功')
    goBack()
  } catch (error) {
    ElMessage.error('操作失败')
  } finally {
    loading.value = false
  }
}

const goBack = () => {
  router.push('/dashboard/roles')
}

onMounted(() => {
  if (isEdit.value) {
    // 模拟加载角色数据
    const roleId = route.params.id
    // 这里应该根据ID加载角色数据
    console.log('编辑角色ID:', roleId)
  }
})
</script>

<style scoped>
.role-form {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
