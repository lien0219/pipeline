package request

// CreateCanaryRelease 创建金丝雀发布请求参数
type CreateCanaryRelease struct {
	Name            string `json:"name" binding:"required"`
	PipelineID      uint   `json:"pipeline_id" binding:"required"`
	PipelineRunID   uint   `json:"pipeline_run_id" binding:"required"`
	TargetNamespace string `json:"target_namespace" binding:"required"`
	TargetService   string `json:"target_service" binding:"required"`
	TrafficPercent  int    `json:"traffic_percent" binding:"required,min=1,max=100"`
}

// UpdateCanaryRelease 更新金丝雀发布请求参数
type UpdateCanaryRelease struct {
	TrafficPercent int    `json:"traffic_percent" binding:"required,min=1,max=100"`
	Status         string `json:"status" binding:"required,oneof=pending running completed rolledback"`
}
