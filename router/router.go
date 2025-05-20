package router

import (
	v1 "gin_pipeline/api/v1"
	"gin_pipeline/middleware"
	"github.com/gin-gonic/gin"
)

// InitPublicRouter 初始化公共路由
func InitPublicRouter(Router *gin.RouterGroup) {
	PublicRouter := Router.Group("")
	{
		// 健康检查
		PublicRouter.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
			})
		})

		// 用户注册登录
		PublicRouter.POST("/user/register", v1.Register)
		PublicRouter.POST("/user/login", v1.Login)
	}
}

// InitUserRouter 初始化用户路由
func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("/user").Use(middleware.JWTAuth())
	{
		UserRouter.GET("/info", v1.GetUserInfo)
		UserRouter.PUT("/info", v1.UpdateUserInfo)
		UserRouter.PUT("/password", v1.ChangePassword)
	}
}

// InitPipelineRouter 初始化流水线路由
func InitPipelineRouter(Router *gin.RouterGroup) {
	PipelineRouter := Router.Group("/pipeline").Use(middleware.JWTAuth())
	{
		PipelineRouter.POST("", v1.CreatePipeline)
		PipelineRouter.GET("", v1.GetPipelines)
		PipelineRouter.GET("/:id", v1.GetPipelineByID)
		PipelineRouter.PUT("/:id", v1.UpdatePipeline)
		PipelineRouter.DELETE("/:id", v1.DeletePipeline)
		PipelineRouter.POST("/:id/trigger", v1.TriggerPipeline)
		PipelineRouter.GET("/:id/runs", v1.GetPipelineRuns)
		PipelineRouter.GET("/:id/runs/:runId", v1.GetPipelineRunByID)
		PipelineRouter.GET("/:id/runs/:runId/logs", v1.GetPipelineRunLogs)
		PipelineRouter.POST("/:id/runs/:runId/cancel", v1.CancelPipelineRun)
	}
}

// InitArtifactRouter 初始化制品路由
func InitArtifactRouter(Router *gin.RouterGroup) {
	ArtifactRouter := Router.Group("/artifact").Use(middleware.JWTAuth())
	{
		ArtifactRouter.POST("", v1.CreateArtifact)
		ArtifactRouter.GET("", v1.GetArtifacts)
		ArtifactRouter.GET("/:id", v1.GetArtifactByID)
		ArtifactRouter.DELETE("/:id", v1.DeleteArtifact)
		ArtifactRouter.GET("/:id/download", v1.DownloadArtifact)
	}
}

// InitEnvironmentRouter 初始化环境路由
func InitEnvironmentRouter(Router *gin.RouterGroup) {
	EnvironmentRouter := Router.Group("/environment").Use(middleware.JWTAuth())
	{
		EnvironmentRouter.POST("", v1.CreateEnvironment)
		EnvironmentRouter.GET("", v1.GetEnvironments)
		EnvironmentRouter.GET("/:id", v1.GetEnvironmentByID)
		EnvironmentRouter.PUT("/:id", v1.UpdateEnvironment)
		EnvironmentRouter.DELETE("/:id", v1.DeleteEnvironment)
	}
}

// InitReleaseRouter 初始化发布路由
func InitReleaseRouter(Router *gin.RouterGroup) {
	ReleaseRouter := Router.Group("/release").Use(middleware.JWTAuth())
	{
		ReleaseRouter.POST("", v1.CreateRelease)
		ReleaseRouter.GET("", v1.GetReleases)
		ReleaseRouter.GET("/:id", v1.GetReleaseByID)
		ReleaseRouter.DELETE("/:id", v1.DeleteRelease)
		ReleaseRouter.POST("/:id/rollback", v1.RollbackRelease)
	}
}

// InitBuildTemplateRouter 初始化构建模板路由
func InitBuildTemplateRouter(Router *gin.RouterGroup) {
	BuildTemplateRouter := Router.Group("/build-template").Use(middleware.JWTAuth())
	{
		BuildTemplateRouter.POST("", v1.CreateBuildTemplate)
		BuildTemplateRouter.GET("", v1.GetBuildTemplates)
		BuildTemplateRouter.GET("/:id", v1.GetBuildTemplateByID)
		BuildTemplateRouter.PUT("/:id", v1.UpdateBuildTemplate)
		BuildTemplateRouter.DELETE("/:id", v1.DeleteBuildTemplate)
		BuildTemplateRouter.POST("/:id/apply", v1.ApplyBuildTemplate)
	}
}
