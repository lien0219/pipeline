package v1

import (
	"gin_pipeline/global"
	"gin_pipeline/model"
	"gin_pipeline/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AddCluster 处理前端提交的集群配置
func AddCluster(c *gin.Context) {
	var cluster model.Cluster
	if err := c.ShouldBindJSON(&cluster); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := service.AddClusterFromFrontend(cluster); err != nil {
		global.Log.Error("添加集群配置失败", zap.Error(err))
		c.JSON(500, gin.H{"error": "添加集群配置失败"})
		return
	}
	c.JSON(200, gin.H{"message": "集群配置添加成功"})
}

// GetAllClusters 获取所有集群配置
func GetAllClusters(c *gin.Context) {
	clusters, err := service.GetAllClusters()
	if err != nil {
		global.Log.Error("获取集群配置失败", zap.Error(err))
		c.JSON(500, gin.H{"error": "获取集群配置失败"})
		return
	}
	c.JSON(200, clusters)
}
