package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"goweb/pkg/config"
	"goweb/pkg/db"
	"goweb/pkg/logger"
	"goweb/pkg/middleware"
	pkgRedis "goweb/pkg/redis"
	_ "goweb/services/admin-api/docs" // 导入生成的docs包
	"goweb/services/admin-api/internal/handler"
	"goweb/services/admin-api/internal/service"
)

func NewRouter(cfg *config.Config, log logger.Logger) *gin.Engine {
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.Recovery(log))
	r.Use(middleware.RequestID())
	r.Use(middleware.AccessLogger(log))
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-Id"},
		ExposeHeaders:    []string{"X-Request-Id"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 健康检查
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "admin-api"})
	})

	// Swagger文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 初始化数据库
	dbManager := db.NewManager(cfg, log)
	if err := dbManager.Connect(); err != nil {
		log.Fatal("Failed to connect database", err)
	}
	db := dbManager.GetDB()

	// 初始化Redis客户端（用于token黑名单）
	redisConfig := cfg.GetRedisConfig()
	redisClient := pkgRedis.NewClient(&redisConfig, log)
	if redisConfig.Enabled {
		log.Info("Redis client initialized for token blacklist")
	} else {
		log.Warn("Redis is disabled, token blacklist feature will not work")
	}

	// 创建服务实例
	userService := service.NewUserService(db, cfg, redisClient)
	roleService := service.NewRoleService(db, cfg, redisClient)
	menuService := service.NewMenuService(db, cfg, redisClient)
	permissionService := service.NewPermissionService(db, cfg, redisClient)

	// 设置服务层日志
	userService.SetLogger(log)
	roleService.SetLogger(log)
	menuService.SetLogger(log)
	permissionService.SetLogger(log)

	// 创建处理器实例
	authHandler := handler.NewAdminAuthHandler(userService)
	userHandler := handler.NewAdminUserHandler(userService)
	roleHandler := handler.NewAdminRoleHandler(roleService)
	menuHandler := handler.NewAdminMenuHandler(menuService)
	permissionHandler := handler.NewAdminPermissionHandler(permissionService)

	// 设置处理器日志
	authHandler.SetLogger(log)
	userHandler.SetLogger(log)
	roleHandler.SetLogger(log)
	menuHandler.SetLogger(log)
	permissionHandler.SetLogger(log)

	// API路由
	api := r.Group("/api/v1")
	{
		// 认证相关路由（不需要JWT验证）
		authGroup := api.Group("/admin/auth")
		{
			authGroup.POST("/login", authHandler.Login)
		}

		// 需要JWT验证的路由（支持token黑名单）
		adminGroup := api.Group("/admin")
		adminGroup.Use(middleware.JWTAuthWithRedis(cfg.GetString("jwt.secret"), redisClient))
		{
			// 认证相关（需要JWT验证）
			adminGroup.POST("/auth/logout", authHandler.Logout)
			adminGroup.GET("/auth/profile", authHandler.GetProfile)
			adminGroup.PUT("/auth/profile", authHandler.UpdateProfile)
			adminGroup.PUT("/auth/change-password", authHandler.ChangePassword)

			// 用户管理
			adminGroup.GET("/users", userHandler.GetUsers)
			adminGroup.POST("/users", userHandler.CreateUser)
			adminGroup.GET("/users/:id", userHandler.GetUser)
			adminGroup.PUT("/users/:id", userHandler.UpdateUser)
			adminGroup.PUT("/users/:id/status", userHandler.UpdateUserStatus)
			adminGroup.DELETE("/users/:id", userHandler.DeleteUser)

			// 角色管理
			adminGroup.GET("/roles", roleHandler.GetRoles)
			adminGroup.POST("/roles", roleHandler.CreateRole)
			adminGroup.GET("/roles/:id", roleHandler.GetRole)
			adminGroup.PUT("/roles/:id", roleHandler.UpdateRole)
			adminGroup.DELETE("/roles/:id", roleHandler.DeleteRole)

			// 菜单管理
			adminGroup.GET("/menus", menuHandler.GetMenus)
			adminGroup.GET("/menus/tree", menuHandler.GetMenuTree)
			adminGroup.POST("/menus", menuHandler.CreateMenu)
			adminGroup.GET("/menus/:id", menuHandler.GetMenu)
			adminGroup.PUT("/menus/:id", menuHandler.UpdateMenu)
			adminGroup.DELETE("/menus/:id", menuHandler.DeleteMenu)

			// 权限管理
			adminGroup.GET("/permissions", permissionHandler.GetPermissions)
			adminGroup.POST("/permissions", permissionHandler.CreatePermission)
			adminGroup.GET("/permissions/:id", permissionHandler.GetPermission)
			adminGroup.PUT("/permissions/:id", permissionHandler.UpdatePermission)
			adminGroup.DELETE("/permissions/:id", permissionHandler.DeletePermission)
		}
	}

	return r
}
