package main

import (
	"fmt"
	"gin_pipeline/global"
	"gin_pipeline/initialize"
	"gin_pipeline/model"
	"gin_pipeline/service"
	"gin_pipeline/utils"
	"time"
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

	// 启动监控任务
	service.StartResourceMonitor()
	// 启动资源回收任务
	service.StartResourceReclamation()
	// 启动队列处理函数
	service.ProcessResourceRequests()

	// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<测试
	// 创建一个资源请求实例
	request := model.TenantResourceRequest{
		TenantID:       "6",
		CPURequest:     2.0,
		MemoryRequest:  4096,
		StorageRequest: 10240,
	}
	// 将资源请求加入队列
	err := service.QueueResourceRequest(request)
	if err != nil {
		utils.Error("将资源请求加入队列失败: %v", err)
		// 可以选择继续执行或者返回，根据实际情况决定
		// return
	}

	tenantID := "60"
	// 动态调整资源配额
	err = service.AdjustResourceQuota(tenantID)
	if err != nil {
		utils.Error("动态调整资源配额失败: %v", err)
		// 可以选择继续执行或者返回，根据实际情况决定
		// return
	}

	// 模拟资源使用
	// 生成资源使用报告
	report, err := service.GenerateResourceReport(tenantID)
	if err != nil {
		utils.Error("生成资源使用报告失败: %v", err)
		// 可以选择继续执行或者返回，根据实际情况决定
		// return
	} else {
		// 打印报告信息
		fmt.Printf("租户 ID: %s\n", report.TenantID)
		fmt.Printf("CPU 使用量: %f\n", report.CPUUsage)
		fmt.Printf("CPU 配额: %f\n", report.CPUQuota)
		fmt.Printf("内存使用量: %d\n", report.MemoryUsage)
		fmt.Printf("内存配额: %d\n", report.MemoryQuota)
		fmt.Printf("存储使用量: %d\n", report.StorageUsage)
		fmt.Printf("存储配额: %d\n", report.StorageQuota)
	}
	// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>测试

	global.StartTime = time.Now()

	// 启动服务器
	port := global.Config.System.Port
	utils.Info("服务器启动成功，监听端口: %s", port)
	if err := r.Run(":" + port); err != nil {
		utils.Error("服务器启动失败: %v", err)
	}
}
