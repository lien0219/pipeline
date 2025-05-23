package request

// CreateHPAPolicy 创建 HPA 策略请求参数
type CreateHPAPolicy struct {
	Name            string `json:"name" binding:"required"`
	Namespace       string `json:"namespace" binding:"required"`
	Deployment      string `json:"deployment" binding:"required"`
	MinReplicas     int    `json:"min_replicas" binding:"required"`
	MaxReplicas     int    `json:"max_replicas" binding:"required"`
	CPUThreshold    int    `json:"cpu_threshold" binding:"required"`
	MemoryThreshold int    `json:"memory_threshold" binding:"required"`
}

// UpdateHPAPolicy 更新 HPA 策略请求参数
type UpdateHPAPolicy struct {
	Name            string `json:"name" binding:"required"`
	Namespace       string `json:"namespace" binding:"required"`
	Deployment      string `json:"deployment" binding:"required"`
	MinReplicas     int    `json:"min_replicas" binding:"required"`
	MaxReplicas     int    `json:"max_replicas" binding:"required"`
	CPUThreshold    int    `json:"cpu_threshold" binding:"required"`
	MemoryThreshold int    `json:"memory_threshold" binding:"required"`
}
