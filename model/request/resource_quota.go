package request

import (
	"reflect"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// CreateResourceQuota 创建资源配额请求参数
type CreateResourceQuota struct {
	TenantID     string  `json:"tenant_id" binding:"required"`
	CPUQuota     float64 `json:"cpu_quota" binding:"required,min=0"`
	MemoryQuota  int64   `json:"memory_quota" binding:"required,min=0"`
	StorageQuota int64   `json:"storage_quota" binding:"required,min=0"`
}

// 自定义验证器：允许0值但必须存在
func validateNonNegative(fl validator.FieldLevel) bool {
	switch fl.Field().Kind() {
	case reflect.Int, reflect.Int64:
		return fl.Field().Int() >= 0
	case reflect.Float32, reflect.Float64:
		return fl.Field().Float() >= 0
	}
	return false
}

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("non_negative", validateNonNegative)
	}
}

// UpdateResourceQuota 更新资源配额请求参数
type UpdateResourceQuota struct {
	CPUQuota     float64 `json:"cpu_quota" binding:"omitempty,min=0"`
	MemoryQuota  int64   `json:"memory_quota" binding:"omitempty,min=0"`
	StorageQuota int64   `json:"storage_quota" binding:"omitempty,min=0"`
}

// CreateResourceRequest 创建资源请求请求参数
type CreateResourceRequest struct {
	TenantID       string  `json:"tenant_id" binding:"required"`
	CPURequest     float64 `json:"cpu_request" binding:"required"`
	MemoryRequest  int64   `json:"memory_request" binding:"required"`
	StorageRequest int64   `json:"storage_request" binding:"required"`
}
