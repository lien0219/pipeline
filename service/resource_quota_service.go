package service

import (
	"gin_pipeline/global"
	"gin_pipeline/model"

	"gorm.io/gorm"
)

// CreateResourceQuota 创建资源配额
func CreateResourceQuota(quota model.ResourceQuota) error {
	return global.DB.Create(&quota).Error
}

// GetResourceQuotaByTenantID 根据租户ID获取资源配额
func GetResourceQuotaByTenantID(tenantID string) (model.ResourceQuota, error) {
	var quota model.ResourceQuota
	err := global.DB.Where("tenant_id = ?", tenantID).First(&quota).Error
	return quota, err
}

// UpdateResourceQuota 更新资源配额
func UpdateResourceQuota(tenantID string, quota model.ResourceQuota) error {
	return global.DB.Model(&model.ResourceQuota{}).Where("tenant_id = ?", tenantID).Updates(quota).Error
}

// CreateResourceRequest 创建资源请求
func CreateResourceRequest(request model.TenantResourceRequest) error {
	return global.DB.Create(&request).Error
}

// GetResourceRequests 获取所有资源请求
func GetResourceRequests() ([]model.TenantResourceRequest, error) {
	var requests []model.TenantResourceRequest
	err := global.DB.Find(&requests).Error
	return requests, err
}

// ApproveResourceRequest 批准资源请求
func ApproveResourceRequest(requestID uint) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		var request model.TenantResourceRequest
		if err := tx.First(&request, requestID).Error; err != nil {
			return err
		}

		// 更新资源配额
		var quota model.ResourceQuota
		if err := tx.Where("tenant_id = ?", request.TenantID).First(&quota).Error; err != nil {
			return err
		}

		quota.CPUQuota += request.CPURequest
		quota.MemoryQuota += request.MemoryRequest
		quota.StorageQuota += request.StorageRequest

		if err := tx.Save(&quota).Error; err != nil {
			return err
		}

		// 更新请求状态
		request.Status = "approved"
		return tx.Save(&request).Error
	})
}

// RejectResourceRequest 拒绝资源请求
func RejectResourceRequest(requestID uint) error {
	return global.DB.Model(&model.TenantResourceRequest{}).Where("id = ?", requestID).Update("status", "rejected").Error
}
