<template>
  <div class="login-container">
    <div class="login-box">
      <div class="login-header">
        <h2>GinForge 管理后台</h2>
        <p>让开发更加简单</p>
      </div>
      
      <el-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginRules"
        class="login-form"
        @submit.prevent="handleLogin"
      >
        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="请输入用户名"
            size="large"
            prefix-icon="User"
          />
        </el-form-item>
        
        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="请输入密码"
            size="large"
            prefix-icon="Lock"
            show-password
            @keyup.enter="handleLogin"
          />
        </el-form-item>
        
        <el-form-item>
          <el-button
            type="primary"
            size="large"
            class="login-button"
            :loading="loading"
            @click="handleLogin"
          >
            {{ loading ? '登录中...' : '登录' }}
          </el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import * as authApi from '@/api/auth'
import { resetRedirectFlag } from '@/api/index'

const router = useRouter()

// 组件挂载时重置重定向标志
onMounted(() => {
  resetRedirectFlag()
})

// 表单数据
const loginForm = reactive({
  username: 'admin',
  password: 'admin123'
})

// 表单验证规则
const loginRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ]
}

const loginFormRef = ref()
const loading = ref(false)

// 处理登录
const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  try {
    await loginFormRef.value.validate()
    loading.value = true
    
    // 调用真实的后端登录API
    const result: any = await authApi.login({
      username: loginForm.username,
      password: loginForm.password
    })
    
    // 保存token和用户信息
    localStorage.setItem('admin_token', result.token)
    localStorage.setItem('admin_user_info', JSON.stringify(result.user))
    localStorage.setItem('admin_permissions', JSON.stringify(result.permissions || []))
    localStorage.setItem('admin_menus', JSON.stringify(result.menus || []))
    
    ElMessage.success('登录成功')
    router.push('/dashboard')
  } catch (error: any) {
    console.error('登录失败:', error)
    ElMessage.error(error?.message || '登录失败，请检查用户名和密码')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}

.login-box {
  width: 400px;
  padding: 40px;
  background: white;
  border-radius: 10px;
  box-shadow: 0 15px 35px rgba(0, 0, 0, 0.1);
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.login-header h2 {
  color: #333;
  margin: 0 0 10px 0;
  font-size: 28px;
}

.login-header p {
  color: #666;
  margin: 0;
  font-size: 14px;
}

.login-form {
  margin-top: 20px;
}

.login-button {
  width: 100%;
  height: 45px;
  font-size: 16px;
}

:deep(.el-input__wrapper) {
  height: 45px;
}

:deep(.el-form-item) {
  margin-bottom: 20px;
}
</style>
