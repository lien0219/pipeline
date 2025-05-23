package service

import (
	"gin_pipeline/global"
	"gin_pipeline/model"
)

// CreateHPAPolicy 创建 HPA 策略
func CreateHPAPolicy(hpaPolicy model.HPAPolicy) error {
	return global.DB.Create(&hpaPolicy).Error
}

// GetHPAPolicies 获取所有 HPA 策略
func GetHPAPolicies() ([]model.HPAPolicy, error) {
	var hpaPolicies []model.HPAPolicy
	err := global.DB.Find(&hpaPolicies).Error
	return hpaPolicies, err
}

// GetHPAPolicyByID 根据 ID 获取 HPA 策略
func GetHPAPolicyByID(id string) (model.HPAPolicy, error) {
	var hpaPolicy model.HPAPolicy
	err := global.DB.Where("id = ?", id).First(&hpaPolicy).Error
	return hpaPolicy, err
}

// UpdateHPAPolicy 更新 HPA 策略
func UpdateHPAPolicy(id string, hpaPolicy model.HPAPolicy) error {
	return global.DB.Where("id = ?", id).Updates(&hpaPolicy).Error
}

// DeleteHPAPolicy 删除 HPA 策略
func DeleteHPAPolicy(id string) error {
	return global.DB.Where("id = ?", id).Delete(&model.HPAPolicy{}).Error
}
