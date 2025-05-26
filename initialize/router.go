package initialize

import (
	"gin_pipeline/global"
	"gin_pipeline/middleware"
	"gin_pipeline/router"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "gin_pipeline/docs" // 导入swagger文档
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	// 设置gin模式
	if global.Config.System.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// 使用中间件
	r.Use(middleware.GinLogger())
	r.Use(middleware.GinRecovery(true))
	r.Use(ginzap.RecoveryWithZap(global.Log.Desugar(), true))

	// 跨域配置
	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	// 根据配置设置CORS
	if global.Config.CORS.Mode == "whitelist" && len(global.Config.CORS.Whitelist) > 0 {
		corsConfig.AllowOrigins = global.Config.CORS.Whitelist
	}

	r.Use(cors.New(corsConfig))

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Swagger文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 注册路由
	apiGroup := r.Group("/api/v1")
	router.InitPublicRouter(apiGroup)         // 公共路由
	router.InitUserRouter(apiGroup)           // 用户路由
	router.InitPipelineRouter(apiGroup)       // 流水线路由
	router.InitArtifactRouter(apiGroup)       // 制品路由
	router.InitEnvironmentRouter(apiGroup)    // 环境路由
	router.InitReleaseRouter(apiGroup)        // 发布路由
	router.InitBuildTemplateRouter(apiGroup)  // 构建模板路由
	router.InitDAGRouter(apiGroup)            // DAG路由
	router.InitYAMLValidatorRouter(apiGroup)  // YAML验证路由
	router.InitTemplateMarketRouter(apiGroup) // 模板市场路由
	router.InitK8sRouter(apiGroup)            // Kubernetes路由
	router.InitClusterRouter(apiGroup)        // 集群路由
	router.InitHPARouter(apiGroup)            // HPA路由
	router.InitResourceQuotaRouter(apiGroup)  // 资源配额路由
	router.InitResourceReportRouter(apiGroup) // 资源报告路由

	global.Log.Info("路由注册成功")
	return r
}
