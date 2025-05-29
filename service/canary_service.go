package service

import (
	"gin_pipeline/global"
	"gin_pipeline/model"

	"go.uber.org/zap"
)

// CanaryService 金丝雀发布服务
type CanaryService struct{}

// CreateCanaryRelease 创建金丝雀发布
func (s *CanaryService) CreateCanaryRelease(canary *model.CanaryRelease) error {
	return global.DB.Create(canary).Error
}

// UpdateCanaryRelease 更新金丝雀发布
func (s *CanaryService) UpdateCanaryRelease(id uint, updates map[string]interface{}) error {
	return global.DB.Model(&model.CanaryRelease{}).Where("id = ?", id).Updates(updates).Error
}

// DeployCanaryRelease 部署金丝雀发布
func (s *CanaryService) DeployCanaryRelease(canaryID uint) error {
	var canary model.CanaryRelease
	if err := global.DB.First(&canary, canaryID).Error; err != nil {
		return err
	}

	// 这里实现实际的金丝雀部署逻辑
	// 通常调用Kubernetes API修改流量分配

	global.Log.Info("部署金丝雀发布",
		zap.Uint("canaryID", canaryID),
		zap.String("service", canary.TargetService),
		zap.Int("trafficPercent", canary.TrafficPercent),
	)

	return nil
}
