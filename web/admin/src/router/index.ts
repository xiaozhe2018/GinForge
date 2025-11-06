import { createRouter, createWebHistory } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { title: '管理员登录', requiresAuth: false }
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('@/layout/index.vue'),
    meta: { title: '仪表盘', requiresAuth: true },
    children: [
      {
        path: '',
        name: 'DashboardHome',
        component: () => import('@/views/Dashboard.vue'),
        meta: { title: '仪表盘', requiresAuth: true }
      },
      // 用户管理
      {
        path: 'users',
        name: 'Users',
        component: () => import('@/views/Users.vue'),
        meta: { title: '用户管理', requiresAuth: true, permission: 'user:read' }
      },
      {
        path: 'users/create',
        name: 'UserCreate',
        component: () => import('@/views/UserForm.vue'),
        meta: { title: '创建用户', requiresAuth: true, permission: 'user:create' }
      },
      {
        path: 'users/:id/edit',
        name: 'UserEdit',
        component: () => import('@/views/UserForm.vue'),
        meta: { title: '编辑用户', requiresAuth: true, permission: 'user:update' }
      },
      // 角色管理
      {
        path: 'roles',
        name: 'Roles',
        component: () => import('@/views/Roles.vue'),
        meta: { title: '角色管理', requiresAuth: true, permission: 'role:read' }
      },
      {
        path: 'roles/create',
        name: 'RoleCreate',
        component: () => import('@/views/RoleForm.vue'),
        meta: { title: '创建角色', requiresAuth: true, permission: 'role:create' }
      },
      {
        path: 'roles/:id/edit',
        name: 'RoleEdit',
        component: () => import('@/views/RoleForm.vue'),
        meta: { title: '编辑角色', requiresAuth: true, permission: 'role:update' }
      },
      // 菜单管理
      {
        path: 'menus',
        name: 'Menus',
        component: () => import('@/views/Menus.vue'),
        meta: { title: '菜单管理', requiresAuth: true, permission: 'menu:read' }
      },
      {
        path: 'menus/create',
        name: 'MenuCreate',
        component: () => import('@/views/MenuForm.vue'),
        meta: { title: '创建菜单', requiresAuth: true, permission: 'menu:create' }
      },
      {
        path: 'menus/:id/edit',
        name: 'MenuEdit',
        component: () => import('@/views/MenuForm.vue'),
        meta: { title: '编辑菜单', requiresAuth: true, permission: 'menu:update' }
      },
      // 权限管理
      {
        path: 'permissions',
        name: 'Permissions',
        component: () => import('@/views/Permissions.vue'),
        meta: { title: '权限管理', requiresAuth: true, permission: 'permission:read' }
      },
      // 系统管理
      {
        path: 'system',
        name: 'System',
        component: () => import('@/views/System.vue'),
        meta: { title: '系统管理', requiresAuth: true, permission: 'system:read' }
      },
      // 个人设置
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/Profile.vue'),
        meta: { title: '个人设置', requiresAuth: true }
      },
      // 文档中心
      {
        path: 'docs',
        name: 'Documentation',
        component: () => import('@/views/Documentation/index.vue'),
        meta: { title: '文档中心', requiresAuth: true }
      },
      // Articles管理
      {
        path: 'articleses',
        name: 'ArticlesList',
        component: () => import('@/views/Articles/index.vue'),
        meta: { title: 'Articles管理', requiresAuth: true }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, _from, next) => {
  // 设置页面标题
  if (to.meta.title) {
    document.title = `${to.meta.title} - GinForge 管理后台`
  }
  
  // 检查是否需要认证
  if (to.meta.requiresAuth !== false) {
    const token = localStorage.getItem('admin_token')
    if (!token) {
      next('/login')
      return
    }
  }
  
  // 检查权限
  if (to.meta.permission) {
    // 获取用户信息，检查是否是超级管理员
    // 使用角色 ID 判断，因为超级管理员角色 ID 通常是 1（固定不变）
    // 使用 code 判断不可靠，因为 code 可能被修改
    const userInfoStr = localStorage.getItem('admin_user_info')
    let isSuperAdmin = false
    if (userInfoStr) {
      try {
        const userInfo = JSON.parse(userInfoStr)
        // 检查用户角色中是否有超级管理员（角色 ID = 1）
        if (userInfo.roles && Array.isArray(userInfo.roles)) {
          isSuperAdmin = userInfo.roles.some((role: any) => role.id === 1)
        }
      } catch (e) {
        console.error('解析用户信息失败:', e)
      }
    }
    
    // 超级管理员跳过权限检查
    if (!isSuperAdmin) {
      const permissions = JSON.parse(localStorage.getItem('admin_permissions') || '[]')
      if (!permissions.includes(to.meta.permission)) {
        ElMessage.error('没有权限访问此页面')
        next('/dashboard')
        return
      }
    }
  }
  
  next()
})

export default router


