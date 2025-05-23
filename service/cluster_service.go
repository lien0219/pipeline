package service

import (
	"gin_pipeline/global"
	"gin_pipeline/model"
)

// AddClusterFromFrontend 从前端接收集群配置，保存到数据库
func AddClusterFromFrontend(cluster model.Cluster) error {
	result := global.DB.Create(&cluster)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetAllClusters 获取所有集群配置
func GetAllClusters() ([]model.Cluster, error) {
	var clusters []model.Cluster
	result := global.DB.Find(&clusters)
	if result.Error != nil {
		return nil, result.Error
	}
	return clusters, nil
}
