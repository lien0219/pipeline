package main

import (
	"gin_pipeline/global"
	"gin_pipeline/initialize"
	_ "gin_pipeline/service"
	"gin_pipeline/utils"
)

// @title           CI/CD Pipeline Visualization API
// @version         1.0
// @description     CI/CD Pipeline Visualization Platform API
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// 初始化配置
	initialize.InitConfig()

	// 初始化日志
	initialize.InitLogger()
	utils.Success("日志系统初始化成功")

	// 初始化数据库
	initialize.InitDB()
	utils.Success("数据库连接初始化成功")

	// // 初始化K8sSvc
	// kubeconfigPath := global.Config.Kubernetes.Kubeconfig
	// var err error
	// service.K8sSvc, err = service.NewK8sService(kubeconfigPath)
	// if err != nil {
	// 	utils.Error("K8sSvc 初始化失败: %v", err)
	// 	return
	// }
	// utils.Success("K8sSvc 初始化成功")

	// 初始化Redis
	initialize.InitRedis()
	utils.Success("Redis连接初始化成功")

	// 初始化路由
	r := initialize.InitRouter()
	utils.Success("路由初始化成功")

	// 启动服务器
	port := global.Config.System.Port
	utils.Info("服务器启动成功，监听端口: %s", port)
	if err := r.Run(":" + port); err != nil {
		utils.Error("服务器启动失败: %v", err)
	}
}
