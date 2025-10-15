import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as systemApi from '@/api/system'

export const useSystemStore = defineStore('system', () => {
  // 系统基本信息
  const systemName = ref('GinForge 管理后台')
  const systemVersion = ref('1.0.0')
  const systemDescription = ref('基于 Go + Gin 的企业级微服务开发框架')
  const systemLogo = ref('/logo.svg')
  const systemLanguage = ref('zh-CN')
  
  // 是否已加载
  const loaded = ref(false)

  // 加载系统基本信息
  const loadSystemInfo = async () => {
    try {
      const data = await systemApi.getSystemBasicInfo()
      
      systemName.value = data['system.name'] || systemName.value
      systemVersion.value = data['system.version'] || systemVersion.value
      systemDescription.value = data['system.description'] || systemDescription.value
      systemLogo.value = data['system.logo'] || systemLogo.value
      systemLanguage.value = data['system.default_language'] || systemLanguage.value
      
      loaded.value = true
      
      // 更新浏览器标签页标题
      document.title = systemName.value
    } catch (error) {
      console.error('加载系统信息失败:', error)
      // 使用默认值
      loaded.value = true
    }
  }

  return {
    systemName,
    systemVersion,
    systemDescription,
    systemLogo,
    systemLanguage,
    loaded,
    loadSystemInfo
  }
})

