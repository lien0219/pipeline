package router

import (
	v1 "gin_pipeline/api/v1"
	"gin_pipeline/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine) {
	// API v1 路由组
	apiV1 := r.Group("/api/v1")

	// 公开路由
	{
		// 健康检查
		apiV1.GET("/health", v1.HealthCheck)

		// 用户认证
		auth := apiV1.Group("/auth")
		{
			auth.POST("/login", v1.Login)
			auth.POST("/register", v1.Register)
		}
	}

	// 需要认证的路由
	apiV1.Use(middleware.JWTAuth())
	{
		// 用户相关
		user := apiV1.Group("/users")
		{
			user.GET("/me", v1.GetUserInfo)
			user.PUT("/me", v1.UpdateUserInfo)
			user.PUT("/password", v1.ChangePassword)
		}

		// 系统配置
		system := apiV1.Group("/system")
		{
			system.GET("/config", v1.GetSystemConfig)
		}
	}
}
