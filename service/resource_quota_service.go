package service

import (
	"context"
	"encoding/json"
	"gin_pipeline/global"
	"gin_pipeline/model"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CreateResourceQuota 创建资源配额
func CreateResourceQuota(quota model.ResourceQuota) error {
	return global.DB.Create(&quota).Error
}

// GetResourceQuotaByTenantID 根据租户ID获取资源配额
func GetResourceQuotaByTenantID(tenantID string) (model.ResourceQuota, error) {
	var quota model.ResourceQuota
	// 尝试从 Redis 中获取缓存
	cacheKey := "resource_quota:" + tenantID
	// 添加 context.Context 参数
	cachedData, err := global.Redis.Get(context.Background(), cacheKey).Bytes()
	if err == nil {
		err = json.Unmarshal(cachedData, &quota)
		if err == nil {
			return quota, nil
		}
		global.Log.Warn("解析 Redis 缓存数据失败", zap.String("cacheKey", cacheKey), zap.Error(err))
	}

	// Redis 中没有缓存，从数据库中获取
	err = global.DB.Where("tenant_id = ?", tenantID).First(&quota).Error
	if err != nil {
		global.Log.Error("从数据库获取资源配额失败", zap.String("tenantID", tenantID), zap.Error(err))
		return quota, err
	}

	// 将结果存入 Redis 缓存
	cachedData, err = json.Marshal(quota)
	if err != nil {
		global.Log.Warn("序列化资源配额数据失败", zap.String("tenantID", tenantID), zap.Error(err))
		return quota, nil
	}
	// 添加 context.Context 参数
	err = global.Redis.Set(context.Background(), cacheKey, cachedData, 5*time.Minute).Err()
	if err != nil {
		global.Log.Warn("将资源配额数据存入 Redis 缓存失败", zap.String("cacheKey", cacheKey), zap.Error(err))
	}

	return quota, nil
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
			if err == gorm.ErrRecordNotFound {
				// 记录日志并跳过该请求
				global.Log.Warn("资源请求记录不存在，跳过处理", zap.Uint("requestID", requestID))
				return nil
			}
			global.Log.Error("查询资源请求失败", zap.Uint("requestID", requestID), zap.Error(err))
			return err
		}

		// 更新资源配额
		var quota model.ResourceQuota
		if err := tx.Where("tenant_id = ?", request.TenantID).First(&quota).Error; err != nil {
			global.Log.Error("查询资源配额失败", zap.String("tenantID", request.TenantID), zap.Error(err))
			return err
		}

		quota.CPUQuota += request.CPURequest
		quota.MemoryQuota += request.MemoryRequest
		quota.StorageQuota += request.StorageRequest

		if err := tx.Save(&quota).Error; err != nil {
			global.Log.Error("更新资源配额失败", zap.String("tenantID", request.TenantID), zap.Error(err))
			return err
		}

		// 更新请求状态
		request.Status = "approved"
		if err := tx.Save(&request).Error; err != nil {
			global.Log.Error("更新资源请求状态失败", zap.Uint("requestID", requestID), zap.Error(err))
			return err
		}
		return nil
	})
}

// RejectResourceRequest 拒绝资源请求
func RejectResourceRequest(requestID uint) error {
	return global.DB.Model(&model.TenantResourceRequest{}).Where("id = ?", requestID).Update("status", "rejected").Error
}

func StartResourceMonitor() {
	go func() {
		for {
			var quotas []model.ResourceQuota
			err := global.DB.Find(&quotas).Error
			if err != nil {
				global.Log.Error("监控资源使用情况时查询数据库失败", zap.Error(err))
				time.Sleep(1 * time.Hour)
				continue
			}

			for _, quota := range quotas {
				// 假设阈值为 80%
				if quota.CPUUsage >= quota.CPUQuota*0.8 ||
					float64(quota.MemoryUsage) >= float64(quota.MemoryQuota)*0.8 ||
					float64(quota.StorageUsage) >= float64(quota.StorageQuota)*0.8 {
					SendAlert(quota.TenantID)
				}
			}

			time.Sleep(1 * time.Hour) // 每小时检查一次
		}
	}()
}

func SendAlert(tenantID string) {
	// 实现告警逻辑，例如发送邮件、短信等
	global.Log.Warn("租户资源使用超过阈值", zap.String("tenant_id", tenantID))
}
func StartResourceReclamation() {
	go func() {
		for {
			// 回收 30 天未使用的资源
			thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
			var requests []model.TenantResourceRequest
			err := global.DB.Where("status = ? AND updated_at < ?", "approved", thirtyDaysAgo).Find(&requests).Error
			if err != nil {
				global.Log.Error("资源回收时查询数据库失败", zap.Error(err))
				time.Sleep(1 * time.Hour)
				continue
			}

			for _, request := range requests {
				err := ReclaimResources(request.TenantID)
				if err != nil {
					global.Log.Error("回收资源失败", zap.String("tenant_id", request.TenantID), zap.Error(err))
				}
			}

			time.Sleep(1 * time.Hour) // 每小时检查一次
		}
	}()
}

func ReclaimResources(tenantID string) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		var quota model.ResourceQuota
		if err := tx.Where("tenant_id = ?", tenantID).First(&quota).Error; err != nil {
			return err
		}

		// 回收资源
		quota.CPUQuota = 0
		quota.MemoryQuota = 0
		quota.StorageQuota = 0

		if err := tx.Save(&quota).Error; err != nil {
			return err
		}

		return nil
	})
}

// GenerateResourceReport 生成资源使用报告
func GenerateResourceReport(tenantID string) (model.ResourceReport, error) {
	var quota model.ResourceQuota
	err := global.DB.Where("tenant_id = ?", tenantID).First(&quota).Error
	if err != nil {
		return model.ResourceReport{}, err
	}

	report := model.ResourceReport{
		TenantID:     tenantID,
		CPUUsage:     quota.CPUUsage,
		CPUQuota:     quota.CPUQuota,
		MemoryUsage:  quota.MemoryUsage,
		MemoryQuota:  quota.MemoryQuota,
		StorageUsage: quota.StorageUsage,
		StorageQuota: quota.StorageQuota,
		CreatedAt:    time.Now(),
	}

	err = global.DB.Create(&report).Error
	return report, err
}

// AdjustResourceQuota 动态调整资源配额
func AdjustResourceQuota(tenantID string) error {
	var quota model.ResourceQuota
	err := global.DB.Where("tenant_id = ?", tenantID).First(&quota).Error
	if err != nil {
		return err
	}

	// 示例策略：如果使用率超过 90%，增加 10% 的配额；如果使用率低于 30%，减少 10% 的配额
	if quota.CPUUsage >= quota.CPUQuota*0.9 {
		quota.CPUQuota *= 1.1
	} else if quota.CPUUsage <= quota.CPUQuota*0.3 {
		quota.CPUQuota *= 0.9
	}

	if float64(quota.MemoryUsage) >= float64(quota.MemoryQuota)*0.9 {
		quota.MemoryQuota = int64(float64(quota.MemoryQuota) * 1.1)
	} else if float64(quota.MemoryUsage) <= float64(quota.MemoryQuota)*0.3 {
		quota.MemoryQuota = int64(float64(quota.MemoryQuota) * 0.9)
	}

	if float64(quota.StorageUsage) >= float64(quota.StorageQuota)*0.9 {
		quota.StorageQuota = int64(float64(quota.StorageQuota) * 1.1)
	} else if float64(quota.StorageUsage) <= float64(quota.StorageQuota)*0.3 {
		quota.StorageQuota = int64(float64(quota.StorageQuota) * 0.9)
	}

	return global.DB.Save(&quota).Error
}

// QueueResourceRequest 将资源请求加入队列
func QueueResourceRequest(request model.TenantResourceRequest) error {
	requestData, err := json.Marshal(request)
	if err != nil {
		return err
	}

	_, err = global.Redis.RPush(context.Background(), "resource_request_queue", requestData).Result()
	return err
}

// ProcessResourceRequests 处理资源请求队列
func ProcessResourceRequests() {
	go func() {
		for {
			result, err := global.Redis.BRPop(context.Background(), 0, "resource_request_queue").Result()
			if err != nil {
				global.Log.Error("处理资源请求队列时出错", zap.Error(err))
				continue
			}

			var request model.TenantResourceRequest
			err = json.Unmarshal([]byte(result[1]), &request)
			if err != nil {
				global.Log.Error("解析资源请求时出错", zap.Error(err))
				continue
			}

			err = ApproveResourceRequest(request.ID)
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					global.Log.Error("批准资源请求失败，重新加入队列", zap.Uint("requestID", request.ID), zap.Error(err))
					_, err = global.Redis.RPush(context.Background(), "resource_request_queue", result[1]).Result()
					if err != nil {
						global.Log.Error("重新加入资源请求队列时出错", zap.Error(err))
					}
					// 可以添加重试机制，例如等待一段时间后重试
					time.Sleep(5 * time.Second)
				} else {
					global.Log.Warn("资源请求记录不存在，不重新加入队列", zap.Uint("requestID", request.ID))
				}
			}
		}
	}()
}

// GetAllResourceReports 获取所有资源报告
func GetAllResourceReports() ([]model.ResourceReport, error) {
	var reports []model.ResourceReport
	err := global.DB.Find(&reports).Error
	return reports, err
}

// GetResourceReportByID 根据 ID 获取资源报告
func GetResourceReportByID(id uint) (model.ResourceReport, error) {
	var report model.ResourceReport
	err := global.DB.First(&report, id).Error
	return report, err
}

// CreateResourceReport 创建资源报告
func CreateResourceReport(report model.ResourceReport) error {
	return global.DB.Create(&report).Error
}

// UpdateResourceReport 更新资源报告
func UpdateResourceReport(id uint, report model.ResourceReport) error {
	return global.DB.Model(&model.ResourceReport{}).Where("id = ?", id).Updates(report).Error
}

// DeleteResourceReport 删除资源报告
func DeleteResourceReport(id uint) error {
	return global.DB.Delete(&model.ResourceReport{}, id).Error
}
