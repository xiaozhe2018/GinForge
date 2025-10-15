<template>
  <div class="login-container">
    <!-- 动态背景 -->
    <div class="background-animation">
      <div class="gradient-bg"></div>
      <div class="shape shape1"></div>
      <div class="shape shape2"></div>
      <div class="shape shape3"></div>
      <div class="shape shape4"></div>
      <div class="shape shape5"></div>
      <div class="particles">
        <div v-for="i in 50" :key="i" class="particle" :style="getParticleStyle(i)"></div>
      </div>
    </div>
    
    <!-- 登录框 -->
    <div class="login-box" :class="{ 'login-box-active': isActive }">
      <div class="login-header">
        <!-- Logo 暂时隐藏 -->
        <!-- <div class="logo-wrapper">
          <img v-if="systemStore.systemLogo" :src="systemStore.systemLogo" alt="Logo" class="logo-image" />
          <div v-else class="logo-placeholder">
            <svg width="70" height="70" viewBox="0 0 60 60" fill="none">
              <circle cx="30" cy="30" r="28" stroke="url(#gradient)" stroke-width="3"/>
              <path d="M20 30L27 37L40 24" stroke="url(#gradient)" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"/>
              <defs>
                <linearGradient id="gradient" x1="0%" y1="0%" x2="100%" y2="100%">
                  <stop offset="0%" style="stop-color:#667eea;stop-opacity:1" />
                  <stop offset="100%" style="stop-color:#764ba2;stop-opacity:1" />
                </linearGradient>
              </defs>
            </svg>
          </div>
        </div> -->
        <h2 class="title">{{ systemStore.systemName }}</h2>
        <p class="subtitle">{{ systemStore.systemDescription }}</p>
        <div class="title-underline"></div>
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
            class="input-field"
            @focus="handleInputFocus"
            @blur="handleInputBlur"
          >
            <template #prefix>
              <el-icon class="input-icon"><User /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        
        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="请输入密码"
            size="large"
            class="input-field"
            show-password
            @focus="handleInputFocus"
            @blur="handleInputBlur"
            @keyup.enter="handleLogin"
          >
            <template #prefix>
              <el-icon class="input-icon"><Lock /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        
        <el-form-item>
          <el-button
            type="primary"
            size="large"
            class="login-button"
            :loading="loading"
            @click="handleLogin"
          >
            <span v-if="!loading" class="button-content">
              <span>登录</span>
              <el-icon class="button-icon"><Right /></el-icon>
            </span>
            <span v-else class="button-content">
              <span>登录中</span>
            </span>
          </el-button>
        </el-form-item>
      </el-form>
      
      <div class="login-footer">
        <p>{{ systemStore.systemVersion }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock, Right } from '@element-plus/icons-vue'
import * as authApi from '@/api/auth'
import { resetRedirectFlag } from '@/api/index'
import { useWebSocketStore } from '@/stores/websocket'
import { useSystemStore } from '@/stores/system'

const router = useRouter()
const wsStore = useWebSocketStore()
const systemStore = useSystemStore()

const isActive = ref(false)

// 组件挂载时加载系统信息和动画
onMounted(() => {
  resetRedirectFlag()
  // 加载系统基本信息
  if (!systemStore.loaded) {
    systemStore.loadSystemInfo()
  }
  // 登录框进入动画
  setTimeout(() => {
    isActive.value = true
  }, 100)
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

// 输入框聚焦处理
const handleInputFocus = (e: FocusEvent) => {
  const input = (e.target as HTMLElement).closest('.input-field')
  if (input) {
    input.classList.add('input-focused')
  }
}

const handleInputBlur = (e: FocusEvent) => {
  const input = (e.target as HTMLElement).closest('.input-field')
  if (input) {
    input.classList.remove('input-focused')
  }
}

// 粒子样式生成
const getParticleStyle = (index: number) => {
  const size = Math.random() * 4 + 2
  const x = Math.random() * 100
  const y = Math.random() * 100
  const duration = Math.random() * 20 + 10
  const delay = Math.random() * 5
  
  return {
    width: `${size}px`,
    height: `${size}px`,
    left: `${x}%`,
    top: `${y}%`,
    animationDuration: `${duration}s`,
    animationDelay: `${delay}s`
  }
}

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
    
    // 连接WebSocket
    wsStore.connect(result.token)
    
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
/* ========== 登录容器 ========== */
.login-container {
  position: relative;
  height: 100vh;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* ========== 动态背景 ========== */
.background-animation {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 0;
}

/* 渐变背景动画 */
.gradient-bg {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 50%, #f093fb 100%);
  background-size: 400% 400%;
  animation: gradientShift 15s ease infinite;
}

@keyframes gradientShift {
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
}

/* 浮动几何图形 */
.shape {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  animation: float 20s ease-in-out infinite;
}

.shape1 {
  width: 300px;
  height: 300px;
  top: -150px;
  left: -100px;
  animation-delay: 0s;
  animation-duration: 25s;
}

.shape2 {
  width: 200px;
  height: 200px;
  top: 60%;
  right: -100px;
  animation-delay: 2s;
  animation-duration: 20s;
}

.shape3 {
  width: 150px;
  height: 150px;
  bottom: -50px;
  left: 30%;
  animation-delay: 4s;
  animation-duration: 22s;
  border-radius: 30% 70% 70% 30% / 30% 30% 70% 70%;
}

.shape4 {
  width: 250px;
  height: 250px;
  top: 20%;
  right: 20%;
  animation-delay: 1s;
  animation-duration: 18s;
}

.shape5 {
  width: 180px;
  height: 180px;
  bottom: 20%;
  left: 10%;
  animation-delay: 3s;
  animation-duration: 24s;
  border-radius: 63% 37% 54% 46% / 55% 48% 52% 45%;
}

@keyframes float {
  0%, 100% {
    transform: translate(0, 0) rotate(0deg);
  }
  33% {
    transform: translate(30px, -50px) rotate(120deg);
  }
  66% {
    transform: translate(-20px, 20px) rotate(240deg);
  }
}

/* 粒子效果 */
.particles {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
  pointer-events: none;
}

.particle {
  position: absolute;
  background: rgba(255, 255, 255, 0.5);
  border-radius: 50%;
  animation: particleFloat 15s linear infinite;
}

@keyframes particleFloat {
  0% {
    transform: translateY(0) translateX(0);
    opacity: 0;
  }
  10% {
    opacity: 1;
  }
  90% {
    opacity: 1;
  }
  100% {
    transform: translateY(-100vh) translateX(100px);
    opacity: 0;
  }
}

/* ========== 登录框 - 商务风格 ========== */
.login-box {
  position: relative;
  z-index: 10;
  width: 450px;
  padding: 50px;
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  box-shadow: 
    0 10px 40px rgba(0, 0, 0, 0.12),
    0 0 0 1px rgba(255, 255, 255, 0.8) inset;
  opacity: 0;
  transform: translateY(30px) scale(0.95);
  transition: all 0.6s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.login-box-active {
  opacity: 1;
  transform: translateY(0) scale(1);
}

/* ========== Logo和标题 - 商务风格 ========== */
.login-header {
  text-align: center;
  margin-bottom: 40px;
}

.logo-wrapper {
  display: flex;
  justify-content: center;
  margin-bottom: 20px;
}

.logo-image,
.logo-placeholder {
  width: 70px;
  height: 70px;
  object-fit: contain;
}

.title {
  color: #333;
  margin: 10px 0 10px 0;
  font-size: 28px;
  font-weight: 600;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.subtitle {
  color: #666;
  margin: 0;
  font-size: 14px;
  opacity: 0.85;
}

.title-underline {
  width: 50px;
  height: 3px;
  background: linear-gradient(90deg, #667eea, #764ba2);
  margin: 15px auto 0;
  border-radius: 2px;
}

/* ========== 表单 ========== */
.login-form {
  margin-top: 30px;
}

/* 表单项对齐 */
:deep(.el-form-item) {
  margin-bottom: 25px;
}

:deep(.el-form-item__content) {
  display: flex;
  flex-direction: column;
}

/* 输入框样式 - 商务风格 */
.input-field {
  position: relative;
}

.input-icon {
  font-size: 18px;
  color: #999;
  transition: color 0.2s ease;
}

.input-field.input-focused .input-icon {
  color: #667eea;
}

:deep(.el-input__wrapper) {
  height: 50px;
  border-radius: 8px;
  transition: border-color 0.2s ease;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  padding: 0 15px;
  background: #fff;
}

:deep(.el-input__wrapper:hover) {
  border-color: #c0c4cc;
}

:deep(.el-input__wrapper.is-focus) {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.1);
}

:deep(.el-input__prefix) {
  display: flex;
  align-items: center;
  margin-right: 8px;
}

:deep(.el-input__inner) {
  font-size: 15px;
  color: #333;
}

:deep(.el-input__inner::placeholder) {
  color: #c0c4cc;
}

/* ========== 登录按钮 - 商务风格 ========== */
.login-button {
  width: 100%;
  height: 50px;
  font-size: 16px;
  font-weight: 600;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  border-radius: 8px;
  transition: opacity 0.2s ease;
  margin-top: 5px;
  cursor: pointer;
}

.login-button:hover {
  opacity: 0.9;
}

.login-button:active {
  opacity: 0.8;
}

.button-content {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.button-icon {
  font-size: 16px;
}

/* ========== 页脚 ========== */
.login-footer {
  text-align: center;
  margin-top: 30px;
  padding-top: 20px;
  border-top: 1px solid rgba(0, 0, 0, 0.05);
}

.login-footer p {
  color: #999;
  font-size: 12px;
  margin: 0;
}

/* ========== 响应式设计 ========== */
@media (max-width: 768px) {
  .login-box {
    width: 90%;
    padding: 30px;
  }
  
  .title {
    font-size: 24px;
  }
  
  .shape {
    display: none;
  }
}

/* ========== 加载动画优化 ========== */
.login-button.is-loading {
  pointer-events: none;
}

:deep(.el-loading-spinner) {
  margin-top: -10px;
}

:deep(.el-loading-spinner .path) {
  stroke: #fff;
}
</style>
