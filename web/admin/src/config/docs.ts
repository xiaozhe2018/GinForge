// 文档配置
export interface DocItem {
  id: string
  title: string
  path: string
  description?: string
}

export interface DocCategory {
  id: string
  name: string
  icon: string
  docs: DocItem[]
}

export const docConfig: DocCategory[] = [
  {
    id: 'getting-started',
    name: '快速开始',
    icon: 'Rocket',
    docs: [
      {
        id: 'introduction',
        title: '框架介绍',
        path: 'getting-started/introduction',
        description: '了解 GinForge 框架的核心特性和优势'
      },
      {
        id: 'installation',
        title: '安装指南',
        path: 'getting-started/installation',
        description: '快速安装和配置 GinForge 框架'
      },
      {
        id: 'quick-start',
        title: '快速开始',
        path: 'getting-started/quick-start',
        description: '5分钟创建你的第一个应用'
      },
      {
        id: 'project-structure',
        title: '项目结构',
        path: 'getting-started/project-structure',
        description: '了解框架的目录结构和文件组织'
      }
    ]
  },
  {
    id: 'tutorials',
    name: '实战教程',
    icon: 'Edit',
    docs: [
      {
        id: 'create-module',
        title: '创建完整业务模块',
        path: 'tutorials/create-module',
        description: '从零到一创建文章管理模块'
      },
      {
        id: 'faq',
        title: '常见问题（FAQ）',
        path: 'tutorials/faq',
        description: '开发中的常见问题和解决方案'
      }
    ]
  },
  {
    id: 'generator',
    name: 'CRUD 代码生成器',
    icon: 'MagicStick',
    docs: [
      {
        id: 'introduction',
        title: '代码生成器介绍',
        path: 'generator/introduction',
        description: '了解强大的脚手架工具，效率提升10倍'
      },
      {
        id: 'usage',
        title: '详细使用指南',
        path: 'generator/usage',
        description: '掌握所有命令和选项'
      },
      {
        id: 'auto-register',
        title: '自动注册功能详解',
        path: 'generator/auto-register',
        description: '一键生成开箱即用，效率再提升6倍'
      },
      {
        id: 'configuration',
        title: '配置文件详解',
        path: 'generator/configuration',
        description: '自定义生成规则和字段配置'
      },
      {
        id: 'examples',
        title: '实战示例',
        path: 'generator/examples',
        description: '完整的文章管理模块开发流程'
      }
    ]
  },
  {
    id: 'core-concepts',
    name: '核心概念',
    icon: 'Reading',
    docs: [
      {
        id: 'configuration',
        title: '配置系统',
        path: 'core-concepts/configuration',
        description: '学习如何配置和管理应用'
      },
      {
        id: 'routing',
        title: '路由管理',
        path: 'core-concepts/routing',
        description: '定义和管理 API 路由'
      },
      {
        id: 'middleware',
        title: '中间件',
        path: 'core-concepts/middleware',
        description: '使用和编写自定义中间件'
      },
      {
        id: 'database',
        title: '数据库操作',
        path: 'core-concepts/database',
        description: 'GORM 数据库操作指南'
      }
    ]
  },
  {
    id: 'features',
    name: '功能特性',
    icon: 'Star',
    docs: [
      {
        id: 'authentication',
        title: '认证授权',
        path: 'features/authentication',
        description: 'JWT 认证和 RBAC 权限控制'
      },
      {
        id: 'file-upload',
        title: '文件上传',
        path: 'features/file-upload',
        description: '本地和云存储文件上传'
      },
      {
        id: 'cache',
        title: '缓存系统',
        path: 'features/cache',
        description: 'Redis 缓存使用指南'
      },
      {
        id: 'websocket',
        title: 'WebSocket',
        path: 'features/websocket',
        description: '实时通信和消息推送'
      }
    ]
  },
  {
    id: 'api-reference',
    name: 'API 参考',
    icon: 'Document',
    docs: [
      {
        id: 'base-classes',
        title: '基础类',
        path: 'api-reference/base-classes',
        description: 'Controller、Service、Repository 基类'
      },
      {
        id: 'utilities',
        title: '工具函数',
        path: 'api-reference/utilities',
        description: '常用工具函数参考'
      },
      {
        id: 'config-options',
        title: '配置选项',
        path: 'api-reference/config-options',
        description: '完整的配置项说明'
      }
    ]
  },
  {
    id: 'deployment',
    name: '部署指南',
    icon: 'Upload',
    docs: [
      {
        id: 'development',
        title: '开发环境',
        path: 'deployment/development',
        description: '搭建本地开发环境'
      },
      {
        id: 'production',
        title: '生产部署',
        path: 'deployment/production',
        description: '生产环境部署最佳实践'
      },
      {
        id: 'docker',
        title: 'Docker 部署',
        path: 'deployment/docker',
        description: '使用 Docker 容器化部署'
      }
    ]
  },
  {
    id: 'advanced',
    name: '高级功能',
    icon: 'Promotion',
    docs: [
      {
        id: 'message-queue',
        title: '消息队列详解',
        path: 'advanced/message-queue',
        description: '深入掌握异步任务处理'
      },
      {
        id: 'distributed-lock',
        title: '分布式锁',
        path: 'advanced/distributed-lock',
        description: '使用 Redis 解决并发问题'
      },
      {
        id: 'performance',
        title: '性能优化',
        path: 'advanced/performance',
        description: '系统性能优化完整指南'
      }
    ]
  },
  {
    id: 'best-practices',
    name: '最佳实践',
    icon: 'TrendCharts',
    docs: [
      {
        id: 'code-style',
        title: '代码规范',
        path: 'best-practices/code-style',
        description: 'Go 代码风格和命名规范'
      },
      {
        id: 'error-handling',
        title: '错误处理',
        path: 'best-practices/error-handling',
        description: '优雅的错误处理方案'
      },
      {
        id: 'security',
        title: '安全建议',
        path: 'best-practices/security',
        description: '应用安全最佳实践'
      }
    ]
  }
]

// 根据路径查找文档
export function findDocByPath(path: string): DocItem | undefined {
  for (const category of docConfig) {
    const doc = category.docs.find(d => d.path === path)
    if (doc) {
      return doc
    }
  }
  return undefined
}

// 获取所有文档列表（用于搜索）
export function getAllDocs(): DocItem[] {
  return docConfig.flatMap(category => category.docs)
}

