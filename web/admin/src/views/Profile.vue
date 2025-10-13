<template>
  <div class="profile-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <h2>个人设置</h2>
      <p>管理个人信息和账户设置</p>
    </div>

    <el-row :gutter="20">
      <!-- 个人信息 -->
      <el-col :span="16">
        <el-card class="profile-card">
          <template #header>
            <span>个人信息</span>
          </template>
          
          <el-form
            ref="profileFormRef"
            :model="profileForm"
            :rules="profileRules"
            label-width="100px"
            style="max-width: 600px;"
          >
            <el-form-item label="头像">
              <div class="avatar-upload">
                <el-avatar :size="80" :src="profileForm.avatar">
                  {{ profileForm.username?.charAt(0) }}
                </el-avatar>
                <el-button type="primary" size="small" @click="handleAvatarUpload" style="margin-left: 10px;">
                  更换头像
                </el-button>
              </div>
            </el-form-item>
            
            <el-form-item label="用户名" prop="username">
              <el-input v-model="profileForm.username" disabled />
            </el-form-item>
            
            <el-form-item label="姓名" prop="name">
              <el-input v-model="profileForm.name" placeholder="请输入姓名" />
            </el-form-item>
            
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="profileForm.email" placeholder="请输入邮箱" />
            </el-form-item>
            
            <el-form-item label="手机号" prop="phone">
              <el-input v-model="profileForm.phone" placeholder="请输入手机号" />
            </el-form-item>
            
            <el-form-item>
              <el-button type="primary" @click="handleUpdateProfile" :loading="profileSubmitting">
                保存信息
              </el-button>
              <el-button @click="handleResetProfile">重置</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>

      <!-- 账户安全 -->
      <el-col :span="8">
        <el-card class="security-card">
          <template #header>
            <span>账户安全</span>
          </template>
          
          <div class="security-item">
            <div class="security-info">
              <div class="security-title">登录密码</div>
              <div class="security-desc">定期更换密码有助于保护账户安全</div>
            </div>
            <el-button type="primary" size="small" @click="handleChangePassword">
              修改密码
            </el-button>
          </div>
          
          <div class="security-item">
            <div class="security-info">
              <div class="security-title">两步验证</div>
              <div class="security-desc">为账户添加额外的安全保护</div>
            </div>
            <el-switch v-model="twoFactorEnabled" @change="handleToggleTwoFactor" />
          </div>
          
          <div class="security-item">
            <div class="security-info">
              <div class="security-title">登录设备</div>
              <div class="security-desc">管理已登录的设备</div>
            </div>
            <el-button type="text" size="small" @click="handleManageDevices">
              管理设备
            </el-button>
          </div>
        </el-card>

        <!-- 最近活动 -->
        <el-card class="activity-card" style="margin-top: 20px;">
          <template #header>
            <span>最近活动</span>
          </template>
          
          <div class="activity-list">
            <div v-for="activity in recentActivities" :key="activity.id" class="activity-item">
              <div class="activity-icon">
                <el-icon :color="getActivityColor(activity.type)">
                  <component :is="getActivityIcon(activity.type)" />
                </el-icon>
              </div>
              <div class="activity-content">
                <div class="activity-title">{{ activity.title }}</div>
                <div class="activity-time">{{ formatDate(activity.time) }}</div>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 修改密码对话框 -->
    <el-dialog
      v-model="passwordDialogVisible"
      title="修改密码"
      width="400px"
      @close="handlePasswordDialogClose"
    >
      <el-form
        ref="passwordFormRef"
        :model="passwordForm"
        :rules="passwordRules"
        label-width="100px"
      >
        <el-form-item label="当前密码" prop="currentPassword">
          <el-input
            v-model="passwordForm.currentPassword"
            type="password"
            placeholder="请输入当前密码"
            show-password
          />
        </el-form-item>
        
        <el-form-item label="新密码" prop="newPassword">
          <el-input
            v-model="passwordForm.newPassword"
            type="password"
            placeholder="请输入新密码"
            show-password
          />
        </el-form-item>
        
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input
            v-model="passwordForm.confirmPassword"
            type="password"
            placeholder="请再次输入新密码"
            show-password
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="passwordDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleUpdatePassword" :loading="passwordSubmitting">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { User, Lock, Monitor, Warning, Check } from '@element-plus/icons-vue'
import * as authApi from '@/api/auth'

// 个人信息表单
const profileFormRef = ref()
const profileForm = reactive({
  username: '',
  name: '',
  email: '',
  phone: '',
  avatar: '',
  status: 1,
  role_ids: [] as number[]
})

const profileSubmitting = ref(false)

// 加载个人信息
const loadProfile = async () => {
  try {
    const response: any = await authApi.getProfile()
    Object.assign(profileForm, {
      username: response.username,
      name: response.name || '',
      email: response.email,
      phone: response.phone || '',
      avatar: response.avatar || '',
      status: response.status,
      role_ids: response.roles ? response.roles.map((r: any) => r.id) : []
    })
  } catch (error) {
    console.error('加载个人信息失败:', error)
    ElMessage.error('加载个人信息失败')
  }
}

// 密码修改表单
const passwordFormRef = ref()
const passwordDialogVisible = ref(false)
const passwordForm = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const passwordSubmitting = ref(false)

// 两步验证
const twoFactorEnabled = ref(false)

// 最近活动
const recentActivities = ref([
  {
    id: 1,
    type: 'login',
    title: '登录系统',
    time: new Date().toISOString()
  },
  {
    id: 2,
    type: 'password',
    title: '修改密码',
    time: new Date(Date.now() - 86400000).toISOString()
  },
  {
    id: 3,
    type: 'profile',
    title: '更新个人信息',
    time: new Date(Date.now() - 172800000).toISOString()
  },
  {
    id: 4,
    type: 'security',
    title: '启用两步验证',
    time: new Date(Date.now() - 259200000).toISOString()
  }
])

// 表单验证规则
const profileRules = {
  name: [
    { required: true, message: '请输入姓名', trigger: 'blur' },
    { min: 2, max: 20, message: '姓名长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ]
}

const passwordRules = {
  currentPassword: [
    { required: true, message: '请输入当前密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (_rule: any, value: any, callback: any) => {
        if (value !== passwordForm.newPassword) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

// 更新个人信息
const handleUpdateProfile = async () => {
  if (!profileFormRef.value) return
  
  try {
    await profileFormRef.value.validate()
    profileSubmitting.value = true
    
    // 调用真实API更新个人信息
    const updateData = {
      email: profileForm.email,
      phone: profileForm.phone,
      name: profileForm.name,
      status: profileForm.status,
      role_ids: profileForm.role_ids
    }
    
    await authApi.updateProfile(updateData)
    
    // 更新localStorage中的用户信息
    const updatedUser = {
      ...JSON.parse(localStorage.getItem('admin_user_info') || '{}'),
      email: profileForm.email,
      phone: profileForm.phone,
      name: profileForm.name
    }
    localStorage.setItem('admin_user_info', JSON.stringify(updatedUser))
    
    ElMessage.success('个人信息更新成功')
  } catch (error) {
    console.error('个人信息更新失败:', error)
    ElMessage.error('个人信息更新失败')
  } finally {
    profileSubmitting.value = false
  }
}

// 重置个人信息
const handleResetProfile = () => {
  loadProfile()
  profileFormRef.value?.resetFields()
}

// 头像上传
const handleAvatarUpload = () => {
  ElMessage.info('头像上传功能待实现')
}

// 修改密码
const handleChangePassword = () => {
  passwordDialogVisible.value = true
}

// 更新密码
const handleUpdatePassword = async () => {
  if (!passwordFormRef.value) return
  
  try {
    await passwordFormRef.value.validate()
    passwordSubmitting.value = true
    
    // 调用真实API修改密码
    await authApi.changePassword({
      old_password: passwordForm.currentPassword,
      new_password: passwordForm.newPassword
    })
    
    ElMessage.success('密码修改成功，请重新登录')
    passwordDialogVisible.value = false
    handleResetPassword()
    
    // 清除token，跳转到登录页
    setTimeout(() => {
      localStorage.removeItem('admin_token')
      localStorage.removeItem('admin_user_info')
      localStorage.removeItem('admin_permissions')
      localStorage.removeItem('admin_menus')
      window.location.href = '/login'
    }, 1500)
  } catch (error) {
    console.error('密码修改失败:', error)
    ElMessage.error('密码修改失败，请检查当前密码是否正确')
  } finally {
    passwordSubmitting.value = false
  }
}

// 关闭密码对话框
const handlePasswordDialogClose = () => {
  Object.assign(passwordForm, {
    currentPassword: '',
    newPassword: '',
    confirmPassword: ''
  })
  passwordFormRef.value?.resetFields()
}

// 切换两步验证
const handleToggleTwoFactor = async (enabled: boolean) => {
  try {
    // 这里应该调用API切换两步验证
    await new Promise(resolve => setTimeout(resolve, 1000)) // 模拟API调用
    
    ElMessage.success(enabled ? '两步验证已启用' : '两步验证已禁用')
  } catch (error) {
    ElMessage.error('操作失败')
    twoFactorEnabled.value = !enabled // 恢复原状态
  }
}

// 管理设备
const handleManageDevices = () => {
  ElMessage.info('设备管理功能待实现')
}

// 获取活动图标
const getActivityIcon = (type: string) => {
  const iconMap: Record<string, any> = {
    login: User,
    password: Lock,
    profile: User,
    security: Lock,
    warning: Warning,
    success: Check
  }
  return iconMap[type] || Monitor
}

// 获取活动颜色
const getActivityColor = (type: string) => {
  const colorMap: Record<string, string> = {
    login: '#409EFF',
    password: '#67C23A',
    profile: '#E6A23C',
    security: '#F56C6C',
    warning: '#E6A23C',
    success: '#67C23A'
  }
  return colorMap[type] || '#909399'
}

// 格式化日期
const formatDate = (date: string) => {
  return new Date(date).toLocaleString('zh-CN')
}

// 组件挂载时获取数据
// 组件挂载时加载数据
onMounted(() => {
  loadProfile()
})
</script>

<style scoped>
.profile-page {
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

.profile-card,
.security-card,
.activity-card {
  border: none;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.avatar-upload {
  display: flex;
  align-items: center;
}

.security-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 0;
  border-bottom: 1px solid #f0f0f0;
}

.security-item:last-child {
  border-bottom: none;
}

.security-info {
  flex: 1;
}

.security-title {
  font-weight: bold;
  color: #333;
  margin-bottom: 5px;
}

.security-desc {
  font-size: 12px;
  color: #666;
}

.activity-list {
  max-height: 300px;
  overflow-y: auto;
}

.activity-item {
  display: flex;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid #f0f0f0;
}

.activity-item:last-child {
  border-bottom: none;
}

.activity-icon {
  margin-right: 15px;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f7fa;
  border-radius: 50%;
}

.activity-content {
  flex: 1;
}

.activity-title {
  font-size: 14px;
  color: #333;
  margin-bottom: 5px;
}

.activity-time {
  font-size: 12px;
  color: #666;
}
</style>



