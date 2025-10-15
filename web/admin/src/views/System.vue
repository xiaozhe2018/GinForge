<template>
  <div class="system-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <h2>系统管理</h2>
      <p>系统配置、监控和日志管理</p>
    </div>

    <!-- 系统概览 -->
    <el-row :gutter="20" class="overview-cards">
      <el-col :span="6">
        <el-card class="overview-card">
          <div class="card-content">
            <div class="card-icon">
              <el-icon size="24" color="#409EFF"><User /></el-icon>
            </div>
            <div class="card-info">
              <div class="card-title">在线用户</div>
              <div class="card-value">{{ systemInfo.onlineUsers }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="overview-card">
          <div class="card-content">
            <div class="card-icon">
              <el-icon size="24" color="#67C23A"><Monitor /></el-icon>
            </div>
            <div class="card-info">
              <div class="card-title">系统负载</div>
              <div class="card-value">{{ systemInfo.cpuUsage }}%</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="overview-card">
          <div class="card-content">
            <div class="card-icon">
              <el-icon size="24" color="#E6A23C"><Coin /></el-icon>
            </div>
            <div class="card-info">
              <div class="card-title">内存使用</div>
              <div class="card-value">{{ systemInfo.memoryUsage }}%</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="overview-card">
          <div class="card-content">
            <div class="card-icon">
              <el-icon size="24" color="#F56C6C"><Warning /></el-icon>
            </div>
            <div class="card-info">
              <div class="card-title">错误日志</div>
              <div class="card-value">{{ systemInfo.errorCount }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 系统配置 -->
    <el-card class="config-card" style="margin-top: 20px;">
      <template #header>
        <div class="card-header">
          <span>系统配置</span>
          <el-button type="primary" size="small" @click="handleSaveConfig" :loading="saveConfigLoading">
            保存配置
          </el-button>
        </div>
      </template>
      
      <el-tabs v-model="activeTab" type="border-card" @tab-click="handleTabClick">
        <!-- 基本配置 -->
        <el-tab-pane label="基本配置" name="basic">
          <el-form :model="basicConfig" label-width="120px" style="max-width: 600px;" v-loading="configLoading">
            <el-form-item label="系统名称">
              <el-input v-model="basicConfig.systemName" />
            </el-form-item>
            <el-form-item label="系统版本">
              <el-input v-model="basicConfig.systemVersion" disabled />
            </el-form-item>
            <el-form-item label="系统描述">
              <el-input
                v-model="basicConfig.systemDescription"
                type="textarea"
                :rows="3"
              />
            </el-form-item>
            <el-form-item label="系统Logo">
              <el-input v-model="basicConfig.systemLogo" />
            </el-form-item>
            <el-form-item label="默认语言">
              <el-select v-model="basicConfig.defaultLanguage" style="width: 200px;">
                <el-option label="中文" value="zh-CN" />
                <el-option label="English" value="en-US" />
              </el-select>
            </el-form-item>
          </el-form>
        </el-tab-pane>
        
        <!-- 安全配置 -->
        <el-tab-pane label="安全配置" name="security">
          <el-form :model="securityConfig" label-width="120px" style="max-width: 600px;" v-loading="configLoading">
            <el-form-item label="密码最小长度">
              <el-input-number v-model="securityConfig.minPasswordLength" :min="6" :max="20" />
            </el-form-item>
            <el-form-item label="密码复杂度">
              <el-checkbox-group v-model="securityConfig.passwordComplexity">
                <el-checkbox label="uppercase">包含大写字母</el-checkbox>
                <el-checkbox label="lowercase">包含小写字母</el-checkbox>
                <el-checkbox label="numbers">包含数字</el-checkbox>
                <el-checkbox label="symbols">包含特殊字符</el-checkbox>
              </el-checkbox-group>
            </el-form-item>
            <el-form-item label="登录失败次数">
              <el-input-number v-model="securityConfig.maxLoginAttempts" :min="3" :max="10" />
            </el-form-item>
            <el-form-item label="账户锁定时间(分钟)">
              <el-input-number v-model="securityConfig.lockoutDuration" :min="5" :max="60" />
            </el-form-item>
            <el-form-item label="会话超时时间(分钟)">
              <el-input-number v-model="securityConfig.sessionTimeout" :min="30" :max="480" />
            </el-form-item>
          </el-form>
        </el-tab-pane>
        
        <!-- 邮件配置 -->
        <el-tab-pane label="邮件配置" name="email">
          <el-form :model="emailConfig" label-width="120px" style="max-width: 600px;" v-loading="configLoading">
            <el-form-item label="SMTP服务器">
              <el-input v-model="emailConfig.smtpHost" />
            </el-form-item>
            <el-form-item label="SMTP端口">
              <el-input-number v-model="emailConfig.smtpPort" :min="1" :max="65535" />
            </el-form-item>
            <el-form-item label="发送邮箱">
              <el-input v-model="emailConfig.fromEmail" />
            </el-form-item>
            <el-form-item label="邮箱密码">
              <el-input v-model="emailConfig.emailPassword" type="password" show-password />
            </el-form-item>
            <el-form-item label="启用SSL">
              <el-switch v-model="emailConfig.enableSSL" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="testEmail" :loading="testEmailLoading">测试邮件发送</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
        
        <!-- 缓存配置 -->
        <el-tab-pane label="缓存配置" name="cache">
          <el-form :model="cacheConfig" label-width="120px" style="max-width: 600px;" v-loading="configLoading">
            <el-form-item label="缓存类型">
              <el-select v-model="cacheConfig.cacheType" style="width: 200px;">
                <el-option label="Redis" value="redis" />
                <el-option label="内存" value="memory" />
              </el-select>
            </el-form-item>
            <el-form-item label="Redis地址">
              <el-input v-model="cacheConfig.redisHost" :disabled="cacheConfig.cacheType !== 'redis'" />
            </el-form-item>
            <el-form-item label="Redis端口">
              <el-input-number v-model="cacheConfig.redisPort" :disabled="cacheConfig.cacheType !== 'redis'" />
            </el-form-item>
            <el-form-item label="Redis密码">
              <el-input v-model="cacheConfig.redisPassword" type="password" :disabled="cacheConfig.cacheType !== 'redis'" />
            </el-form-item>
            <el-form-item label="默认过期时间(秒)">
              <el-input-number v-model="cacheConfig.defaultExpiration" :min="60" :max="86400" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="testCache" :loading="testCacheLoading">测试缓存连接</el-button>
              <el-button type="danger" @click="clearCache" :loading="clearCacheLoading">清空缓存</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- 系统日志 -->
    <el-card class="log-card" style="margin-top: 20px;">
      <template #header>
        <div class="card-header">
          <span>系统日志</span>
          <div>
            <el-select v-model="logLevel" placeholder="日志级别" style="width: 120px; margin-right: 10px;">
              <el-option label="全部" value="" />
              <el-option label="错误" value="error" />
              <el-option label="警告" value="warn" />
              <el-option label="信息" value="info" />
              <el-option label="调试" value="debug" />
            </el-select>
            <el-button type="primary" size="small" @click="fetchLogs" :loading="logLoading">
              刷新
            </el-button>
            <el-button type="danger" size="small" @click="clearLogs" :loading="clearLogsLoading">
              清空日志
            </el-button>
          </div>
        </div>
      </template>
      
      <el-table :data="logList" v-loading="logLoading" max-height="400">
        <el-table-column prop="timestamp" label="时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.created_at || row.timestamp) }}
          </template>
        </el-table-column>
        <el-table-column prop="level" label="级别" width="80">
          <template #default="{ row }">
            <el-tag :type="getLogLevelType(row.level || getLogLevelFromStatusCode(row.status_code))">
              {{ (row.level || getLogLevelFromStatusCode(row.status_code)).toUpperCase() }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="message" label="消息" min-width="300">
          <template #default="{ row }">
            {{ row.message || `${row.method} ${row.path}` }}
          </template>
        </el-table-column>
        <el-table-column prop="source" label="来源" width="120">
          <template #default="{ row }">
            {{ row.source || row.username || '系统' }}
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { User, Monitor, Coin, Warning } from '@element-plus/icons-vue'
import * as systemApi from '@/api/system'

// 当前激活的标签页
const activeTab = ref('basic')

// 加载状态
const configLoading = ref(false)
const saveConfigLoading = ref(false)
const testEmailLoading = ref(false)
const testCacheLoading = ref(false)
const clearCacheLoading = ref(false)
const logLoading = ref(false)
const clearLogsLoading = ref(false)

// 系统信息
const systemInfo = reactive({
  onlineUsers: 0,
  cpuUsage: 0,
  memoryUsage: 0,
  errorCount: 0
})

// 基本配置
const basicConfig = reactive({
  systemName: 'GinForge 管理后台',
  systemVersion: '1.0.0',
  systemDescription: '基于 Go + Gin 的企业级微服务开发框架',
  systemLogo: '/logo.svg',
  defaultLanguage: 'zh-CN'
})

// 安全配置
const securityConfig = reactive({
  minPasswordLength: 8,
  passwordComplexity: ['lowercase', 'numbers'],
  maxLoginAttempts: 5,
  lockoutDuration: 15,
  sessionTimeout: 120
})

// 邮件配置
const emailConfig = reactive({
  smtpHost: 'smtp.example.com',
  smtpPort: 587,
  fromEmail: 'noreply@example.com',
  emailPassword: '',
  enableSSL: true
})

// 缓存配置
const cacheConfig = reactive({
  cacheType: 'redis',
  redisHost: 'localhost',
  redisPort: 6379,
  redisPassword: '',
  defaultExpiration: 3600
})

// 日志相关
const logLevel = ref('')
const logList = ref<any[]>([])

// 配置键映射
const configKeyMap = {
  basic: {
    systemName: 'system.name',
    systemVersion: 'system.version',
    systemDescription: 'system.description',
    systemLogo: 'system.logo',
    defaultLanguage: 'system.default_language'
  },
  security: {
    minPasswordLength: 'security.min_password_length',
    passwordComplexity: 'security.password_complexity',
    maxLoginAttempts: 'security.max_login_attempts',
    lockoutDuration: 'security.lockout_duration',
    sessionTimeout: 'security.session_timeout'
  },
  email: {
    smtpHost: 'email.smtp_host',
    smtpPort: 'email.smtp_port',
    fromEmail: 'email.from_email',
    emailPassword: 'email.email_password',
    enableSSL: 'email.enable_ssl'
  },
  cache: {
    cacheType: 'cache.type',
    redisHost: 'cache.redis_host',
    redisPort: 'cache.redis_port',
    redisPassword: 'cache.redis_password',
    defaultExpiration: 'cache.default_expiration'
  }
}

// 获取系统信息
const fetchSystemInfo = async () => {
  try {
    const data = await systemApi.getSystemInfo()
    
    systemInfo.onlineUsers = data.online_users || 0
    systemInfo.cpuUsage = data.cpu_usage || 0
    systemInfo.memoryUsage = data.memory_usage || 0
    systemInfo.errorCount = 0 // 错误日志数量需要另外获取
  } catch (error) {
    console.error('获取系统信息失败:', error)
    ElMessage.error('获取系统信息失败')
  }
}

// 加载配置
const loadConfigs = async (group: string) => {
  configLoading.value = true
  try {
    const response = await systemApi.getSystemConfigList({ group })
    const configs = response.list
    
    // 根据当前标签页更新配置
    switch (group) {
      case 'basic':
        updateBasicConfig(configs)
        break
      case 'security':
        updateSecurityConfig(configs)
        break
      case 'email':
        updateEmailConfig(configs)
        break
      case 'cache':
        updateCacheConfig(configs)
        break
    }
  } catch (error) {
    console.error('获取配置失败:', error)
    ElMessage.error('获取配置失败')
  } finally {
    configLoading.value = false
  }
}

// 更新基本配置
const updateBasicConfig = (configs: systemApi.SystemConfig[]) => {
  configs.forEach(config => {
    switch (config.key) {
      case 'system.name':
        basicConfig.systemName = config.value
        break
      case 'system.version':
        basicConfig.systemVersion = config.value
        break
      case 'system.description':
        basicConfig.systemDescription = config.value
        break
      case 'system.logo':
        basicConfig.systemLogo = config.value
        break
      case 'system.default_language':
        basicConfig.defaultLanguage = config.value
        break
    }
  })
}

// 更新安全配置
const updateSecurityConfig = (configs: systemApi.SystemConfig[]) => {
  configs.forEach(config => {
    switch (config.key) {
      case 'security.min_password_length':
        securityConfig.minPasswordLength = parseInt(config.value)
        break
      case 'security.password_complexity':
        try {
          securityConfig.passwordComplexity = JSON.parse(config.value)
        } catch (e) {
          securityConfig.passwordComplexity = ['lowercase', 'numbers']
        }
        break
      case 'security.max_login_attempts':
        securityConfig.maxLoginAttempts = parseInt(config.value)
        break
      case 'security.lockout_duration':
        securityConfig.lockoutDuration = parseInt(config.value)
        break
      case 'security.session_timeout':
        securityConfig.sessionTimeout = parseInt(config.value)
        break
    }
  })
}

// 更新邮件配置
const updateEmailConfig = (configs: systemApi.SystemConfig[]) => {
  configs.forEach(config => {
    switch (config.key) {
      case 'email.smtp_host':
        emailConfig.smtpHost = config.value
        break
      case 'email.smtp_port':
        emailConfig.smtpPort = parseInt(config.value)
        break
      case 'email.from_email':
        emailConfig.fromEmail = config.value
        break
      case 'email.email_password':
        emailConfig.emailPassword = config.value
        break
      case 'email.enable_ssl':
        emailConfig.enableSSL = config.value === 'true'
        break
    }
  })
}

// 更新缓存配置
const updateCacheConfig = (configs: systemApi.SystemConfig[]) => {
  configs.forEach(config => {
    switch (config.key) {
      case 'cache.type':
        cacheConfig.cacheType = config.value
        break
      case 'cache.redis_host':
        cacheConfig.redisHost = config.value
        break
      case 'cache.redis_port':
        cacheConfig.redisPort = parseInt(config.value)
        break
      case 'cache.redis_password':
        cacheConfig.redisPassword = config.value
        break
      case 'cache.default_expiration':
        cacheConfig.defaultExpiration = parseInt(config.value)
        break
    }
  })
}

// 标签页切换
const handleTabClick = () => {
  loadConfigs(activeTab.value)
}

// 获取日志列表
const fetchLogs = async () => {
  logLoading.value = true
  try {
    const params: any = {
      page: 1,
      page_size: 50
    }
    
    if (logLevel.value) {
      // 这里需要根据后端API调整，目前后端没有level字段
      // params.level = logLevel.value
    }
    
    const response = await systemApi.getSystemLogList(params)
    logList.value = response.list
  } catch (error) {
    console.error('获取日志失败:', error)
    ElMessage.error('获取日志失败')
  } finally {
    logLoading.value = false
  }
}

// 保存配置
const handleSaveConfig = async () => {
  saveConfigLoading.value = true
  try {
    const configMap = configKeyMap[activeTab.value as keyof typeof configKeyMap]
    const config = activeTab.value === 'basic' ? basicConfig :
                  activeTab.value === 'security' ? securityConfig :
                  activeTab.value === 'email' ? emailConfig :
                  cacheConfig
    
    // 保存当前标签页的所有配置
    const promises = Object.entries(configMap).map(([key, configKey]) => {
      let value = (config as any)[key]
      
      // 特殊处理某些类型
      if (key === 'passwordComplexity') {
        value = JSON.stringify(value)
      } else if (key === 'enableSSL') {
        value = value.toString()
      }
      
      return systemApi.updateSystemConfig(configKey as string, { value: value.toString() })
    })
    
    await Promise.all(promises)
    ElMessage.success('配置保存成功')
  } catch (error) {
    console.error('配置保存失败:', error)
    ElMessage.error('配置保存失败')
  } finally {
    saveConfigLoading.value = false
  }
}

// 测试邮件发送
const testEmail = async () => {
  // 先保存邮件配置
  await handleSaveConfig()
  
  testEmailLoading.value = true
  try {
    const email = await ElMessageBox.prompt('请输入测试邮箱地址', '测试邮件发送', {
      confirmButtonText: '发送',
      cancelButtonText: '取消',
      inputPattern: /[\w!#$%&'*+/=?^_`{|}~-]+(?:\.[\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\w](?:[\w-]*[\w])?\.)+[\w](?:[\w-]*[\w])?/,
      inputErrorMessage: '邮箱格式不正确'
    })
    
    await systemApi.testEmailConfig(email.value)
    ElMessage.success('测试邮件发送成功')
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('测试邮件发送失败:', error)
      ElMessage.error(error.message || '测试邮件发送失败')
    }
  } finally {
    testEmailLoading.value = false
  }
}

// 测试缓存连接
const testCache = async () => {
  // 先保存缓存配置
  await handleSaveConfig()
  
  testCacheLoading.value = true
  try {
    await systemApi.testCacheConnection()
    ElMessage.success('缓存连接测试成功')
  } catch (error) {
    console.error('缓存连接测试失败:', error)
    ElMessage.error('缓存连接测试失败')
  } finally {
    testCacheLoading.value = false
  }
}

// 清空缓存
const clearCache = async () => {
  try {
    await ElMessageBox.confirm('确定要清空所有缓存吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    clearCacheLoading.value = true
    await systemApi.clearCache()
    ElMessage.success('缓存清空成功')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('缓存清空失败:', error)
      ElMessage.error('缓存清空失败')
    }
  } finally {
    clearCacheLoading.value = false
  }
}

// 清空日志
const clearLogs = async () => {
  try {
    await ElMessageBox.confirm('确定要清空所有日志吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    clearLogsLoading.value = true
    await systemApi.clearSystemLogs()
    logList.value = []
    ElMessage.success('日志清空成功')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('日志清空失败:', error)
      ElMessage.error('日志清空失败')
    }
  } finally {
    clearLogsLoading.value = false
  }
}

// 获取日志级别标签类型
const getLogLevelType = (level: string) => {
  const typeMap: Record<string, string> = {
    error: 'danger',
    warn: 'warning',
    info: 'success',
    debug: 'info'
  }
  return typeMap[level] || 'info'
}

// 根据状态码获取日志级别
const getLogLevelFromStatusCode = (statusCode: number) => {
  if (statusCode >= 500) return 'error'
  if (statusCode >= 400) return 'warn'
  return 'info'
}

// 格式化日期
const formatDate = (date: string) => {
  return new Date(date).toLocaleString('zh-CN')
}

// 组件挂载时获取数据
onMounted(() => {
  fetchSystemInfo()
  loadConfigs(activeTab.value)
  fetchLogs()
  
  // 定时更新系统信息
  setInterval(() => {
    fetchSystemInfo()
  }, 30000)
})
</script>

<style scoped>
.system-page {
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

.overview-cards {
  margin-bottom: 20px;
}

.overview-card {
  border: none;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.card-content {
  display: flex;
  align-items: center;
}

.card-icon {
  margin-right: 15px;
  padding: 10px;
  background-color: #f5f7fa;
  border-radius: 8px;
}

.card-info {
  flex: 1;
}

.card-title {
  font-size: 14px;
  color: #666;
  margin-bottom: 5px;
}

.card-value {
  font-size: 24px;
  font-weight: bold;
  color: #333;
}

.config-card,
.log-card {
  border: none;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: bold;
  color: #333;
}
</style>