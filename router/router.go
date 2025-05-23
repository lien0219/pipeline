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

// InitDAGRouter 初始化DAG路由
func InitDAGRouter(Router *gin.RouterGroup) {
	DAGRouter := Router.Group("/dag").Use(middleware.JWTAuth())
	{
		DAGRouter.POST("", v1.CreateDAG)
		DAGRouter.GET("/:id", v1.GetDAGByID)
		DAGRouter.GET("/pipeline/:pipelineId", v1.GetDAGsByPipelineID)
		DAGRouter.GET("/pipeline/:pipelineId/active", v1.GetActiveDAG)
		DAGRouter.PUT("/:id", v1.UpdateDAG)
		DAGRouter.DELETE("/:id", v1.DeleteDAG)
		DAGRouter.POST("/validate", v1.ValidateDAG)
		DAGRouter.POST("/:id/version", v1.CreateDAGVersion)
		DAGRouter.GET("/pipeline/:pipelineId/history", v1.GetDAGHistory)
		DAGRouter.POST("/:id/activate", v1.ActivateDAG)
	}
}

// InitYAMLValidatorRouter 初始化YAML验证路由
func InitYAMLValidatorRouter(Router *gin.RouterGroup) {
	YAMLRouter := Router.Group("/yaml").Use(middleware.JWTAuth())
	{
		YAMLRouter.POST("/validate", v1.ValidateYAML)
		YAMLRouter.GET("/history", v1.GetValidationHistory)
		YAMLRouter.POST("/schema", v1.CreateYAMLSchema)
		YAMLRouter.GET("/schema", v1.GetYAMLSchemas)
		YAMLRouter.GET("/schema/:id", v1.GetYAMLSchemaByID)
		YAMLRouter.PUT("/schema/:id", v1.UpdateYAMLSchema)
		YAMLRouter.DELETE("/schema/:id", v1.DeleteYAMLSchema)
	}
}

// InitTemplateMarketRouter 初始化模板市场路由
func InitTemplateMarketRouter(Router *gin.RouterGroup) {
	TemplateMarketRouter := Router.Group("/template-market").Use(middleware.JWTAuth())
	{
		// 分类管理
		TemplateMarketRouter.POST("/category", v1.CreateTemplateCategory)
		TemplateMarketRouter.GET("/category", v1.GetTemplateCategories)
		TemplateMarketRouter.PUT("/category/:id", v1.UpdateTemplateCategory)
		TemplateMarketRouter.DELETE("/category/:id", v1.DeleteTemplateCategory)

		// 模板管理
		TemplateMarketRouter.POST("/template", v1.CreateTemplate)
		TemplateMarketRouter.GET("/template", v1.GetTemplates)
		TemplateMarketRouter.GET("/template/:id", v1.GetTemplateByID)
		TemplateMarketRouter.PUT("/template/:id", v1.UpdateTemplate)
		TemplateMarketRouter.DELETE("/template/:id", v1.DeleteTemplate)

		// 版本管理
		TemplateMarketRouter.POST("/template/:id/version", v1.CreateTemplateVersion)
		TemplateMarketRouter.GET("/template/:id/version", v1.GetTemplateVersions)
		TemplateMarketRouter.GET("/template/:id/version/:versionId", v1.GetTemplateVersionByID)
		TemplateMarketRouter.DELETE("/template/:id/version/:versionId", v1.DeleteTemplateVersion)
		TemplateMarketRouter.POST("/template/:id/version/:versionId/latest", v1.SetVersionAsLatest)

		// 搜索和下载
		TemplateMarketRouter.GET("/search", v1.SearchTemplates)
		TemplateMarketRouter.GET("/template/:id/download", v1.DownloadTemplate)
	}
}

// InitK8sRouter 初始化 Kubernetes 路由
func InitK8sRouter(Router *gin.RouterGroup) {
	K8sRouter := Router.Group("/k8s")
	K8sRouter.Use(middleware.JWTAuth())
	{
		// Job 操作路由
		JobRouter := K8sRouter.Group("/jobs")
		JobRouter.Use(middleware.JWTAuth())
		{
			JobRouter.POST("/:namespace", v1.CreateJobHandler)
			JobRouter.GET("/:namespace/:name", v1.GetJobHandler)
			JobRouter.PUT("/:namespace/:name", v1.UpdateJobHandler)
			JobRouter.DELETE("/:namespace/:name", v1.DeleteJobHandler)
		}

		// CronJob 操作路由
		CronJobRouter := K8sRouter.Group("/cronjobs")
		CronJobRouter.Use(middleware.JWTAuth())
		{
			CronJobRouter.POST("/:namespace", v1.CreateCronJobHandler)
			CronJobRouter.GET("/:namespace/:name", v1.GetCronJobHandler)
			CronJobRouter.PUT("/:namespace/:name", v1.UpdateCronJobHandler)
			CronJobRouter.DELETE("/:namespace/:name", v1.DeleteCronJobHandler)
		}

		// Deployment 操作路由
		DeploymentRouter := K8sRouter.Group("/deployments")
		DeploymentRouter.Use(middleware.JWTAuth())
		{
			DeploymentRouter.POST("/:namespace", v1.CreateDeploymentHandler)
			DeploymentRouter.GET("/:namespace/:name", v1.GetDeploymentHandler)
			DeploymentRouter.PUT("/:namespace/:name", v1.UpdateDeploymentHandler)
			DeploymentRouter.DELETE("/:namespace/:name", v1.DeleteDeploymentHandler)
		}
	}
}

// InitClusterRouter 初始化集群路由
func InitClusterRouter(Router *gin.RouterGroup) {
	ClusterRouter := Router.Group("/cluster").Use(middleware.JWTAuth())
	{
		ClusterRouter.POST("", v1.AddCluster)
		ClusterRouter.GET("/all", v1.GetAllClusters)
	}
}
