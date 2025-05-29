package service

import (
	"gin_pipeline/global"
	"gin_pipeline/model"
	"time"

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

	// 更新状态为部署中
	if err := global.DB.Model(&canary).Update("status", "in_progress").Error; err != nil {
		return err
	}

	// 启动渐进式流量切换
	go s.gradualTrafficShift(canaryID, canary.TrafficPercent)

	return nil
}
func (s *CanaryService) gradualTrafficShift(canaryID uint, targetPercent int) {
	var canary model.CanaryRelease
	if err := global.DB.First(&canary, canaryID).Error; err != nil {
		global.Log.Error("获取金丝雀发布信息失败", zap.Error(err))
		return
	}

	// 渐进式增加流量百分比
	for percent := 10; percent <= targetPercent; percent += 10 {
		time.Sleep(30 * time.Second) // 每30秒增加10%流量

		// 更新流量百分比
		updates := map[string]interface{}{
			"traffic_percent": percent,
		}
		if percent == targetPercent {
			updates["status"] = "completed"
		}

		if err := s.UpdateCanaryRelease(canaryID, updates); err != nil {
			global.Log.Error("更新金丝雀发布流量失败",
				zap.Uint("canaryID", canaryID),
				zap.Int("percent", percent),
				zap.Error(err))
			return
		}

		global.Log.Info("更新金丝雀发布流量",
			zap.Uint("canaryID", canaryID),
			zap.Int("percent", percent))
	}
}
