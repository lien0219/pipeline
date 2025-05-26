package request

// CreateResourceQuota 创建资源配额请求参数
type CreateResourceQuota struct {
	TenantID     string  `json:"tenant_id" binding:"required"`
	CPUQuota     float64 `json:"cpu_quota" binding:"required"`
	MemoryQuota  int64   `json:"memory_quota" binding:"required"`
	StorageQuota int64   `json:"storage_quota" binding:"required"`
}

// UpdateResourceQuota 更新资源配额请求参数
type UpdateResourceQuota struct {
	CPUQuota     float64 `json:"cpu_quota" binding:"required"`
	MemoryQuota  int64   `json:"memory_quota" binding:"required"`
	StorageQuota int64   `json:"storage_quota" binding:"required"`
}

// CreateResourceRequest 创建资源请求请求参数
type CreateResourceRequest struct {
	TenantID       string  `json:"tenant_id" binding:"required"`
	CPURequest     float64 `json:"cpu_request" binding:"required"`
	MemoryRequest  int64   `json:"memory_request" binding:"required"`
	StorageRequest int64   `json:"storage_request" binding:"required"`
}
