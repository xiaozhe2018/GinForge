<template>
  <div class="admin-layout">
    <!-- 侧边栏 -->
    <el-aside :width="isCollapse ? '64px' : '200px'" class="sidebar">
      <div class="logo" @click="router.push('/dashboard')">
        <span v-if="!isCollapse" class="logo-text">{{ systemName }}</span>
        <span v-else class="logo-text-mini">{{ systemNameShort }}</span>
      </div>
      
      <el-menu
        :default-active="activeMenu"
        :collapse="isCollapse"
        :unique-opened="true"
        router
        class="sidebar-menu"
      >
        <template v-for="menu in menuTree" :key="menu.id">
          <!-- 有子菜单的菜单项（directory 类型或 children 不为空） -->
          <el-sub-menu v-if="(menu.children && menu.children.length > 0) || menu.type === 'directory'" :index="menu.path || String(menu.id)">
            <template #title>
              <el-icon v-if="menu.icon"><component :is="getIcon(menu.icon)" /></el-icon>
              <span>{{ menu.name }}</span>
            </template>
            <template v-if="menu.children && menu.children.length > 0">
              <template v-for="child in menu.children" :key="child.id">
                <el-menu-item v-if="(!child.children || child.children.length === 0) && child.type !== 'directory'" :index="child.path || String(child.id)">
                  <el-icon v-if="child.icon"><component :is="getIcon(child.icon)" /></el-icon>
                  <span>{{ child.name }}</span>
        </el-menu-item>
                <!-- 支持三级菜单 -->
                <el-sub-menu v-else-if="child.children && child.children.length > 0" :index="child.path || String(child.id)">
                  <template #title>
                    <el-icon v-if="child.icon"><component :is="getIcon(child.icon)" /></el-icon>
                    <span>{{ child.name }}</span>
                  </template>
                  <el-menu-item v-for="grandchild in child.children" :key="grandchild.id" :index="grandchild.path || String(grandchild.id)">
                    <el-icon v-if="grandchild.icon"><component :is="getIcon(grandchild.icon)" /></el-icon>
                    <span>{{ grandchild.name }}</span>
        </el-menu-item>
                </el-sub-menu>
              </template>
            </template>
          </el-sub-menu>
          <!-- 没有子菜单的菜单项 -->
          <el-menu-item v-else :index="menu.path || String(menu.id)">
            <el-icon v-if="menu.icon"><component :is="getIcon(menu.icon)" /></el-icon>
            <span>{{ menu.name }}</span>
        </el-menu-item>
        </template>
      </el-menu>
    </el-aside>

    <!-- 主内容区 -->
    <el-container class="main-container">
      <!-- 顶部导航 -->
      <el-header class="header">
        <div class="header-left">
          <el-button
            type="text"
            @click="toggleCollapse"
            class="collapse-btn"
          >
            <el-icon><Fold v-if="!isCollapse" /><Expand v-else /></el-icon>
          </el-button>
          
          <el-breadcrumb separator="/" class="breadcrumb">
            <el-breadcrumb-item
              v-for="item in breadcrumbList"
              :key="item.path"
              :to="item.path"
            >
              {{ item.title }}
            </el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        
        <div class="header-right">
          <!-- 实时通知中心 -->
          <NotificationCenter />
          
          <!-- 用户信息 -->
          <el-dropdown @command="handleUserCommand">
            <div class="user-info">
              <el-avatar :size="32" :src="userInfo.avatar">
                {{ userInfo.name.charAt(0) }}
              </el-avatar>
              <span class="user-name">{{ userInfo.name }}</span>
              <el-icon><ArrowDown /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">
                  <el-icon><User /></el-icon>
                  个人设置
                </el-dropdown-item>
                <el-dropdown-item command="logout" divided>
                  <el-icon><SwitchButton /></el-icon>
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <!-- 主内容 -->
      <el-main class="main-content">
        <router-view />
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { logout } from '@/api/auth'
import NotificationCenter from '@/components/NotificationCenter.vue'
import { useSystemStore } from '@/stores/system'
import { getMenuTree } from '@/api/menu'
import {
  House,
  User,
  UserFilled,
  Menu,
  Lock,
  Setting,
  Reading,
  Fold,
  Expand,
  ArrowDown,
  SwitchButton,
  Document
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const systemStore = useSystemStore()

// 侧边栏折叠状态
const isCollapse = ref(false)

// 系统名称
const systemName = computed(() => {
  return systemStore.systemName || 'GinForge 管理后台'
})

// 系统名称简写（取首字母或前两个字符）
const systemNameShort = computed(() => {
  const name = systemName.value.trim()
  if (name.length <= 2) {
    return name
  }
  
  // 优先提取中文字符的第一个字
  const chineseMatch = name.match(/[\u4e00-\u9fa5]/)
  if (chineseMatch) {
    return chineseMatch[0]
  }
  
  // 提取英文字母，如果有多个单词，取每个单词的首字母
  const words = name.match(/[A-Za-z]+/g)
  if (words && words.length > 0) {
    if (words.length === 1) {
      // 单个单词，取前2个字母
      return words[0].substring(0, 2).toUpperCase()
    } else {
      // 多个单词，取每个单词的首字母
      return words.map(w => w[0].toUpperCase()).join('').substring(0, 2)
    }
  }
  
  // 默认取前2个字符
  return name.substring(0, 2).toUpperCase()
})

// 用户信息
const userInfo = ref({
  name: '管理员',
  avatar: '',
  role: '超级管理员'
})

// 菜单树
const menuTree = ref<any[]>([])

// 当前激活的菜单
const activeMenu = computed(() => {
  return route.path
})

// 获取图标组件
const getIcon = (iconName?: string) => {
  if (!iconName) return House
  const iconMap: Record<string, any> = {
    House,
    User,
    UserFilled,
    Menu,
    Lock,
    Setting,
    Reading,
    Document,
    Key: Lock
  }
  return iconMap[iconName] || House
}

// 加载菜单树
const loadMenuTree = async () => {
  try {
    // 先尝试从 localStorage 读取（登录时保存的）
    const savedMenus = localStorage.getItem('admin_menus')
    if (savedMenus) {
      try {
        const parsed = JSON.parse(savedMenus)
        if (Array.isArray(parsed) && parsed.length > 0) {
          // 过滤并处理菜单数据
          const filterMenu = (menus: any[]): any[] => {
            return menus
              .filter(menu => menu.visible === 1 && menu.status === 1)
              .map(menu => ({
                ...menu,
                path: menu.path || (menu.type === 'directory' ? '' : undefined),
                children: menu.children ? filterMenu(menu.children) : []
              }))
          }
          menuTree.value = filterMenu(parsed)
          // 如果菜单数据有效，直接返回
          if (menuTree.value.length > 0) {
            return
          }
        }
      } catch (e) {
        console.error('解析保存的菜单失败:', e)
      }
    }
    
    // 从后端获取菜单树
    const response = await getMenuTree()
    if (response && response.list) {
      // 过滤掉隐藏的菜单和禁用的菜单
      const filterMenu = (menus: any[]): any[] => {
        return menus
          .filter(menu => menu.visible === 1 && menu.status === 1)
          .map(menu => ({
            ...menu,
            path: menu.path || (menu.type === 'directory' ? '' : undefined),
            children: menu.children ? filterMenu(menu.children) : []
          }))
      }
      menuTree.value = filterMenu(response.list)
      // 保存到 localStorage
      localStorage.setItem('admin_menus', JSON.stringify(menuTree.value))
    }
  } catch (error) {
    console.error('加载菜单树失败:', error)
    // 如果加载失败，使用默认菜单
    menuTree.value = [
      {
        id: 1,
        name: '仪表盘',
        path: '/dashboard',
        icon: 'House',
        children: []
      }
    ]
  }
}

// 面包屑导航
const breadcrumbList = computed(() => {
  const matched = route.matched.filter(item => item.meta && item.meta.title)
  return matched.map(item => ({
    title: item.meta?.title,
    path: item.path
  }))
})

// 切换侧边栏折叠状态
const toggleCollapse = () => {
  isCollapse.value = !isCollapse.value
}

// 加载用户信息
onMounted(() => {
  // 加载系统信息
  if (!systemStore.loaded) {
    systemStore.loadSystemInfo()
  }
  
  // 加载菜单树
  loadMenuTree()
  
  const savedUserInfo = localStorage.getItem('admin_user_info')
  if (savedUserInfo) {
    try {
      const parsed = JSON.parse(savedUserInfo)
      userInfo.value = {
        name: parsed.name || parsed.username || '管理员',
        avatar: parsed.avatar || '',
        role: parsed.roles?.[0]?.name || '管理员'
      }
    } catch (e) {
      console.error('加载用户信息失败:', e)
    }
  }
})

// 处理用户下拉菜单命令
const handleUserCommand = (command: string) => {
  switch (command) {
    case 'profile':
      router.push('/dashboard/profile')
      break
    case 'logout':
      handleLogout()
      break
  }
}

// 退出登录
const handleLogout = async () => {
  try {
    await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    // 调用后端登出API
    try {
      await logout()
    } catch (error) {
      console.error('登出API调用失败:', error)
      // 即使API调用失败，也继续清除本地数据
    }
    
    // 清除本地存储
    localStorage.removeItem('admin_token')
    localStorage.removeItem('admin_permissions')
    localStorage.removeItem('admin_user_info')
    localStorage.removeItem('admin_menus')
    
    // 跳转到登录页
    router.push('/login')
    ElMessage.success('已退出登录')
  } catch {
    // 用户取消
  }
}

// 监听路由变化，更新面包屑
watch(route, () => {
  // 路由变化时的处理逻辑
}, { immediate: true })
</script>

<style scoped>
.admin-layout {
  display: flex;
  height: 100vh;
}

.sidebar {
  background-color: #304156;
  transition: width 0.3s;
  overflow: hidden;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #2b3a4b;
  border-bottom: 1px solid #434a50;
  cursor: pointer;
  transition: all 0.3s;
}

.logo:hover {
  background-color: #263445;
}

.logo-text {
  color: #fff;
  font-size: 18px;
  font-weight: bold;
  letter-spacing: 1px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  padding: 0 10px;
}

.logo-text-mini {
  color: #fff;
  font-size: 20px;
  font-weight: bold;
  letter-spacing: 2px;
}

.sidebar-menu {
  border: none;
  background-color: #304156;
}

/* 统一菜单样式 - 使用 Element Plus 标准深色主题 */
:deep(.sidebar-menu .el-menu-item),
:deep(.sidebar-menu .el-sub-menu__title) {
  color: #bfcbd9;
  background-color: transparent;
}

/* 菜单项悬停效果 */
:deep(.sidebar-menu .el-menu-item:hover),
:deep(.sidebar-menu .el-sub-menu__title:hover) {
  background-color: #263445;
  color: #fff;
}

/* 激活的菜单项 */
:deep(.sidebar-menu .el-menu-item.is-active) {
  background-color: #409eff;
  color: #fff;
}

/* 子菜单背景 */
:deep(.sidebar-menu .el-sub-menu .el-menu) {
  background-color: #1f2d3d;
}

/* 子菜单项样式 */
:deep(.sidebar-menu .el-sub-menu .el-menu-item) {
  background-color: transparent;
  color: #bfcbd9;
}

:deep(.sidebar-menu .el-sub-menu .el-menu-item:hover) {
  background-color: #263445;
  color: #fff;
}

:deep(.sidebar-menu .el-sub-menu .el-menu-item.is-active) {
  background-color: #409eff;
  color: #fff;
}

.main-container {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.header {
  background-color: #fff;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
}

.header-left {
  display: flex;
  align-items: center;
}

.collapse-btn {
  margin-right: 20px;
  font-size: 18px;
}

.breadcrumb {
  font-size: 14px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 20px;
}

.notification-badge {
  cursor: pointer;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 8px 12px;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.user-info:hover {
  background-color: #f5f5f5;
}

.user-name {
  font-size: 14px;
  color: #333;
}

.main-content {
  background-color: #f5f5f5;
  padding: 20px;
  overflow-y: auto;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .sidebar {
    position: fixed;
    left: 0;
    top: 0;
    height: 100vh;
    z-index: 1000;
    transform: translateX(-100%);
    transition: transform 0.3s;
  }
  
  .sidebar.show {
    transform: translateX(0);
  }
  
  .main-container {
    margin-left: 0;
  }
}
</style>


