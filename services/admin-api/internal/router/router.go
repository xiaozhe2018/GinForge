package router

import (
	"goweb/pkg/config"
	"goweb/pkg/logger"
	"goweb/pkg/middleware"
	"goweb/pkg/notification"
	"goweb/pkg/redis"
	"goweb/pkg/response"
	"goweb/services/admin-api/internal/handler"
	"goweb/services/admin-api/internal/repository"
	"goweb/services/admin-api/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// NewRouter 创建路由
func NewRouter(db *gorm.DB, redisClient *redis.Client, notifyService *notification.Service, log logger.Logger, cfg *config.Config) *gin.Engine {
	r := gin.New()

	// 中间件
	r.Use(gin.Recovery())
	r.Use(middleware.CORS()) // 添加 CORS 中间件，必须放在最前面
	r.Use(middleware.SecurityHeaders(nil)) // 添加安全响应头中间件
	r.Use(middleware.RequestID())
	// 开发环境使用 gin.Logger()，生产环境使用 middleware.AccessLogger(log)
	if cfg.IsProduction() {
		r.Use(middleware.AccessLogger(log))
	} else {
		r.Use(gin.Logger()) // 开发环境：简洁的请求日志
	}
	r.Use(middleware.OperationLogger(db)) // 添加操作日志中间件

	// 健康检查
	r.GET("/healthz", func(c *gin.Context) {
		response.Success(c, "OK")
	})

	// 初始化服务
	// 先初始化系统服务，其他服务需要依赖它
	adminSystemService := service.NewAdminSystemService(db, redisClient, notifyService, log)

	userService := service.NewUserService(db, cfg, redisClient)
	userService.SetLogger(log)
	userService.SetSystemService(adminSystemService) // 注入系统服务用于安全配置

	roleService := service.NewRoleService(db, cfg, redisClient)
	roleService.SetLogger(log)
	permissionService := service.NewPermissionService(db, cfg, redisClient)
	permissionService.SetLogger(log)
	menuService := service.NewMenuService(db, cfg, redisClient)
	menuService.SetLogger(log)
	notificationService := service.NewNotificationService(db, redisClient, log)

	// 初始化处理器
	adminAuthHandler := handler.NewAdminAuthHandler(userService)
	adminAuthHandler.SetLogger(log)
	adminUserHandler := handler.NewAdminUserHandler(userService)
	adminUserHandler.SetLogger(log)
	adminRoleHandler := handler.NewAdminRoleHandler(roleService)
	adminRoleHandler.SetLogger(log)
	adminPermissionHandler := handler.NewAdminPermissionHandler(permissionService)
	adminPermissionHandler.SetLogger(log)
	adminMenuHandler := handler.NewAdminMenuHandler(menuService)
	adminMenuHandler.SetLogger(log)
	adminSystemHandler := handler.NewAdminSystemHandler(adminSystemService, notifyService, log)
	notificationHandler := handler.NewNotificationHandler(notificationService)

	// 初始化 Articles
	articlesRepo := repository.NewArticlesRepository(db)
	articlesService := service.NewArticlesService(articlesRepo, log)
	articlesHandler := handler.NewArticlesHandler(articlesService, log)
	notificationHandler.SetLogger(log)

	// API路由组
	api := r.Group("/api/v1/admin")

	// 无需认证的路由
	api.POST("/login", adminAuthHandler.Login)
	api.GET("/system/basic-info", adminSystemHandler.GetSystemBasicInfo) // 获取系统基本信息（公开接口）
	// 获取 CSRF Token（公开接口，用于登录前获取Token）
	// 注意：CSRF Token 会在首次访问时自动生成并设置到 Cookie
	// 此接口主要用于前端获取 Token 值
	api.GET("/csrf-token", func(c *gin.Context) {
		// 生成临时CSRF配置来获取Token
		csrfConfig := middleware.DefaultCSRFConfig()
		csrfConfig.CookieSecure = cfg.IsProduction() // 生产环境使用 HTTPS Cookie
		token, err := csrfConfig.GenerateToken()
		if err != nil {
			response.Error(c, 500, "生成CSRF令牌失败")
			return
		}
		// 设置Cookie
		c.SetCookie(
			csrfConfig.CookieName,
			token,
			3600, // 1小时过期
			csrfConfig.CookiePath,
			csrfConfig.CookieDomain,
			csrfConfig.CookieSecure,
			true, // HttpOnly
		)
		response.Success(c, gin.H{"csrf_token": token})
	})

	// 需要认证的路由
	auth := api.Group("")
	// 使用 JWT 中间件（从配置读取 secret）
	jwtSecret := cfg.GetString("jwt.secret")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-in-production" // 默认值
	}
	auth.Use(middleware.JWTAuthWithRedis(jwtSecret, redisClient))
	// 添加 CSRF 防护（在认证之后，只对需要认证的路由启用）
	csrfConfig := middleware.DefaultCSRFConfig()
	csrfConfig.CookieSecure = cfg.IsProduction() // 生产环境使用 HTTPS Cookie
	auth.Use(middleware.CSRF(csrfConfig))

	// 用户相关路由
	auth.GET("/users", adminUserHandler.GetUsers)
	auth.GET("/users/:id", adminUserHandler.GetUser)
	auth.POST("/users", adminUserHandler.CreateUser)
	auth.PUT("/users/:id", adminUserHandler.UpdateUser)
	auth.PUT("/users/:id/status", adminUserHandler.UpdateUserStatus)
	auth.PUT("/users/:id/reset-password", adminUserHandler.ResetPassword)
	auth.DELETE("/users/:id", adminUserHandler.DeleteUser)
	auth.POST("/logout", adminAuthHandler.Logout)
	auth.GET("/profile", adminAuthHandler.GetProfile)
	auth.PUT("/profile", adminAuthHandler.UpdateProfile)
	auth.PUT("/change-password", adminAuthHandler.ChangePassword)

	// 角色相关路由
	auth.GET("/roles", adminRoleHandler.GetRoles)
	auth.GET("/roles/:id", adminRoleHandler.GetRole)
	auth.POST("/roles", adminRoleHandler.CreateRole)
	auth.PUT("/roles/:id", adminRoleHandler.UpdateRole)
	auth.DELETE("/roles/:id", adminRoleHandler.DeleteRole)

	// 权限相关路由
	auth.GET("/permissions", adminPermissionHandler.GetPermissions)
	auth.GET("/permissions/:id", adminPermissionHandler.GetPermission)
	auth.POST("/permissions", adminPermissionHandler.CreatePermission)
	auth.PUT("/permissions/:id", adminPermissionHandler.UpdatePermission)
	auth.PUT("/permissions/:id/status", adminPermissionHandler.UpdatePermissionStatus)
	auth.DELETE("/permissions/:id", adminPermissionHandler.DeletePermission)

	// 菜单相关路由
	auth.GET("/menus", adminMenuHandler.GetMenus)
	auth.GET("/menus/:id", adminMenuHandler.GetMenu)
	auth.POST("/menus", adminMenuHandler.CreateMenu)
	auth.PUT("/menus/:id", adminMenuHandler.UpdateMenu)
	auth.DELETE("/menus/:id", adminMenuHandler.DeleteMenu)

	// Articles管理 路由
	auth.GET("/articleses", articlesHandler.List)
	auth.GET("/articleses/:id", articlesHandler.Get)
	auth.POST("/articleses", articlesHandler.Create)
	auth.PUT("/articleses/:id", articlesHandler.Update)
	auth.DELETE("/articleses/:id", articlesHandler.Delete)
	auth.GET("/menus/tree", adminMenuHandler.GetMenuTree)

	// 系统管理路由
	auth.GET("/system/info", adminSystemHandler.GetSystemInfo)
	auth.GET("/system/configs", adminSystemHandler.GetConfigList)
	auth.GET("/system/configs/:key", adminSystemHandler.GetConfig)
	auth.PUT("/system/configs/:key", adminSystemHandler.UpdateConfig)
	auth.POST("/system/email/test", adminSystemHandler.TestEmailConfig)
	auth.POST("/system/cache/test", adminSystemHandler.TestCacheConnection)
	auth.POST("/system/cache/clear", adminSystemHandler.ClearCache)
	auth.GET("/system/logs", adminSystemHandler.GetLogList)
	auth.POST("/system/logs/clear", adminSystemHandler.ClearLogs)
	auth.GET("/system/recent-login-users", adminSystemHandler.GetRecentLoginUsers)
	auth.GET("/system/runtime", adminSystemHandler.GetRuntimeInfo)
	auth.GET("/system/health", adminSystemHandler.HealthCheck)

	// 通知相关路由
	auth.POST("/notifications/system", notificationHandler.SendSystemNotification)
	auth.POST("/notifications/user", notificationHandler.SendUserNotification)
	auth.POST("/notifications/order", notificationHandler.SendOrderNotification)

	return r
}
