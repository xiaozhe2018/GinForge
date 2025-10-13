<template>
  <div class="dashboard">
    <!-- 欢迎信息 -->
    <div class="welcome-section">
      <div class="welcome-content">
        <h1>欢迎回来，{{ userInfo.name }}！</h1>
        <p>今天是 {{ currentDate }}，祝您工作愉快！</p>
      </div>
      <div class="weather-info">
        <el-icon size="24" color="#409EFF"><Sunny /></el-icon>
        <span>晴 22°C</span>
      </div>
    </div>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-cards">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon users">
              <el-icon size="32"><User /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.totalUsers }}</div>
              <div class="stat-label">总用户数</div>
              <div class="stat-change positive">
                <el-icon><ArrowUp /></el-icon>
                +12.5%
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon orders">
              <el-icon size="32"><ShoppingCart /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.totalOrders }}</div>
              <div class="stat-label">总订单数</div>
              <div class="stat-change positive">
                <el-icon><ArrowUp /></el-icon>
                +8.2%
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon revenue">
              <el-icon size="32"><Money /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">¥{{ stats.totalRevenue }}</div>
              <div class="stat-label">总收入</div>
              <div class="stat-change positive">
                <el-icon><ArrowUp /></el-icon>
                +15.3%
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon online">
              <el-icon size="32"><Monitor /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.onlineUsers }}</div>
              <div class="stat-label">在线用户</div>
              <div class="stat-change negative">
                <el-icon><ArrowDown /></el-icon>
                -2.1%
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 图表和表格区域 -->
    <el-row :gutter="20" class="charts-section">
      <!-- 访问趋势图 -->
      <el-col :span="16">
        <el-card class="chart-card">
          <template #header>
            <div class="card-header">
              <span>访问趋势</span>
              <el-radio-group v-model="chartPeriod" size="small">
                <el-radio-button label="7d">7天</el-radio-button>
                <el-radio-button label="30d">30天</el-radio-button>
                <el-radio-button label="90d">90天</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div class="chart-container">
            <div class="chart-placeholder">
              <el-icon size="48" color="#ddd"><TrendCharts /></el-icon>
              <p>访问趋势图表</p>
              <small>这里将显示访问量趋势图</small>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 系统状态 -->
      <el-col :span="8">
        <el-card class="status-card">
          <template #header>
            <span>系统状态</span>
          </template>
          <div class="status-list">
            <div class="status-item">
              <div class="status-label">CPU使用率</div>
              <div class="status-value">
                <el-progress :percentage="systemStatus.cpu" :color="getProgressColor(systemStatus.cpu)" />
                <span>{{ systemStatus.cpu }}%</span>
              </div>
            </div>
            <div class="status-item">
              <div class="status-label">内存使用率</div>
              <div class="status-value">
                <el-progress :percentage="systemStatus.memory" :color="getProgressColor(systemStatus.memory)" />
                <span>{{ systemStatus.memory }}%</span>
              </div>
            </div>
            <div class="status-item">
              <div class="status-label">磁盘使用率</div>
              <div class="status-value">
                <el-progress :percentage="systemStatus.disk" :color="getProgressColor(systemStatus.disk)" />
                <span>{{ systemStatus.disk }}%</span>
              </div>
            </div>
            <div class="status-item">
              <div class="status-label">网络状态</div>
              <div class="status-value">
                <el-tag :type="systemStatus.network ? 'success' : 'danger'">
                  {{ systemStatus.network ? '正常' : '异常' }}
                </el-tag>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 数据表格区域 -->
    <el-row :gutter="20" class="tables-section">
      <!-- 最近用户 -->
      <el-col :span="12">
        <el-card class="table-card">
          <template #header>
            <div class="card-header">
              <span>最近注册用户</span>
              <el-button type="text" @click="$router.push('/dashboard/users')">查看全部</el-button>
            </div>
          </template>
          <el-table :data="recentUsers" style="width: 100%">
            <el-table-column prop="username" label="用户名" width="120" />
            <el-table-column prop="email" label="邮箱" min-width="150" />
            <el-table-column prop="status" label="状态" width="80">
              <template #default="{ row }">
                <el-tag :type="row.status === 'active' ? 'success' : 'danger'">
                  {{ row.status === 'active' ? '正常' : '禁用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="createdAt" label="注册时间" width="120">
              <template #default="{ row }">
                {{ formatDate(row.createdAt) }}
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>

      <!-- 系统日志 -->
      <el-col :span="12">
        <el-card class="table-card">
          <template #header>
            <div class="card-header">
              <span>系统日志</span>
              <el-button type="text" @click="$router.push('/dashboard/system')">查看全部</el-button>
            </div>
          </template>
          <div class="log-list">
            <div v-for="log in systemLogs" :key="log.id" class="log-item">
              <div class="log-icon">
                <el-icon :color="getLogColor(log.level)">
                  <component :is="getLogIcon(log.level)" />
                </el-icon>
              </div>
              <div class="log-content">
                <div class="log-message">{{ log.message }}</div>
                <div class="log-time">{{ formatTime(log.time) }}</div>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 快捷操作 -->
    <el-card class="quick-actions">
      <template #header>
        <span>快捷操作</span>
      </template>
      <div class="actions-grid">
        <el-button type="primary" @click="$router.push('/dashboard/users/create')">
          <el-icon><UserFilled /></el-icon>
          创建用户
        </el-button>
        <el-button type="success" @click="$router.push('/dashboard/roles/create')">
          <el-icon><UserFilled /></el-icon>
          创建角色
        </el-button>
        <el-button type="warning" @click="$router.push('/dashboard/menus/create')">
          <el-icon><Menu /></el-icon>
          创建菜单
        </el-button>
        <el-button type="info" @click="$router.push('/dashboard/system')">
          <el-icon><Setting /></el-icon>
          系统设置
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { 
  User, 
  ShoppingCart, 
  Money, 
  Monitor, 
  Sunny, 
  ArrowUp, 
  ArrowDown, 
  TrendCharts,
  UserFilled,
  Menu,
  Setting,
  Warning,
  Check,
  InfoFilled
} from '@element-plus/icons-vue'

// 用户信息
const userInfo = reactive({
  name: '管理员',
  role: '超级管理员'
})

// 当前日期
const currentDate = computed(() => {
  return new Date().toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    weekday: 'long'
  })
})

// 统计数据
const stats = reactive({
  totalUsers: 12580,
  totalOrders: 8965,
  totalRevenue: 1256800,
  onlineUsers: 156
})

// 图表周期
const chartPeriod = ref('7d')

// 系统状态
const systemStatus = reactive({
  cpu: 45,
  memory: 68,
  disk: 32,
  network: true
})

// 最近用户
const recentUsers = ref([
  {
    id: 1,
    username: 'user001',
    email: 'user001@example.com',
    status: 'active',
    createdAt: new Date().toISOString()
  },
  {
    id: 2,
    username: 'user002',
    email: 'user002@example.com',
    status: 'active',
    createdAt: new Date(Date.now() - 3600000).toISOString()
  },
  {
    id: 3,
    username: 'user003',
    email: 'user003@example.com',
    status: 'disabled',
    createdAt: new Date(Date.now() - 7200000).toISOString()
  }
])

// 系统日志
const systemLogs = ref([
  {
    id: 1,
    level: 'info',
    message: '用户 admin 登录成功',
    time: new Date().toISOString()
  },
  {
    id: 2,
    level: 'warning',
    message: '数据库连接池使用率较高',
    time: new Date(Date.now() - 300000).toISOString()
  },
  {
    id: 3,
    level: 'error',
    message: 'Redis 连接失败，正在重试',
    time: new Date(Date.now() - 600000).toISOString()
  },
  {
    id: 4,
    level: 'info',
    message: '系统启动完成',
    time: new Date(Date.now() - 900000).toISOString()
  }
])

// 获取进度条颜色
const getProgressColor = (percentage: any) => {
  if (percentage < 50) return '#67C23A'
  if (percentage < 80) return '#E6A23C'
  return '#F56C6C'
}

// 获取日志图标
const getLogIcon = (level: any) => {
  const iconMap: any = {
    info: InfoFilled,
    warning: Warning,
    error: Warning,
    success: Check
  }
  return iconMap[level] || InfoFilled
}

// 获取日志颜色
const getLogColor = (level: any) => {
  const colorMap: any = {
    info: '#409EFF',
    warning: '#E6A23C',
    error: '#F56C6C',
    success: '#67C23A'
  }
  return colorMap[level] || '#909399'
}

// 格式化日期
const formatDate = (date: any) => {
  return new Date(date).toLocaleDateString('zh-CN')
}

// 格式化时间
const formatTime = (date: any) => {
  return new Date(date).toLocaleTimeString('zh-CN')
}

// 组件挂载时获取数据
onMounted(() => {
  // 这里应该从API获取仪表盘数据
  console.log('Dashboard mounted')
})
</script>

<style scoped>
.dashboard {
  padding: 0;
}

.welcome-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding: 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  color: white;
}

.welcome-content h1 {
  margin: 0 0 8px 0;
  font-size: 24px;
  font-weight: bold;
}

.welcome-content p {
  margin: 0;
  opacity: 0.9;
}

.weather-info {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
}

.stats-cards {
  margin-bottom: 24px;
}

.stat-card {
  border: none;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.stat-content {
  display: flex;
  align-items: center;
}

.stat-icon {
  width: 64px;
  height: 64px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16px;
}

.stat-icon.users {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.stat-icon.orders {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  color: white;
}

.stat-icon.revenue {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  color: white;
}

.stat-icon.online {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
  color: white;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  color: #333;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: #666;
  margin-bottom: 8px;
}

.stat-change {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  font-weight: bold;
}

.stat-change.positive {
  color: #67C23A;
}

.stat-change.negative {
  color: #F56C6C;
}

.charts-section {
  margin-bottom: 24px;
}

.chart-card,
.status-card,
.table-card {
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

.chart-container {
  height: 300px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.chart-placeholder {
  text-align: center;
  color: #999;
}

.chart-placeholder p {
  margin: 8px 0 4px 0;
  font-size: 16px;
}

.chart-placeholder small {
  font-size: 12px;
}

.status-list {
  /* 状态列表样式 */
  margin: 0;
  padding: 0;
}

.status-item {
  margin-bottom: 16px;
}

.status-label {
  font-size: 14px;
  color: #666;
  margin-bottom: 8px;
}

.status-value {
  display: flex;
  align-items: center;
  gap: 12px;
}

.status-value span {
  font-size: 14px;
  font-weight: bold;
  color: #333;
  min-width: 40px;
}

.tables-section {
  margin-bottom: 24px;
}

.log-list {
  max-height: 300px;
  overflow-y: auto;
}

.log-item {
  display: flex;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
}

.log-item:last-child {
  border-bottom: none;
}

.log-icon {
  margin-right: 12px;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f7fa;
  border-radius: 50%;
}

.log-content {
  flex: 1;
}

.log-message {
  font-size: 14px;
  color: #333;
  margin-bottom: 4px;
}

.log-time {
  font-size: 12px;
  color: #666;
}

.quick-actions {
  border: none;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.actions-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.actions-grid .el-button {
  height: 48px;
  font-size: 14px;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .charts-section .el-col:first-child {
    margin-bottom: 20px;
  }
}

@media (max-width: 768px) {
  .welcome-section {
    flex-direction: column;
    text-align: center;
    gap: 16px;
  }
  
  .stats-cards .el-col {
    margin-bottom: 16px;
  }
  
  .tables-section .el-col {
    margin-bottom: 20px;
  }
}
</style>
