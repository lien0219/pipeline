package request

// CreateWebhook 创建webhook请求参数
type CreateWebhook struct {
	Name       string `json:"name" binding:"required"`
	URL        string `json:"url" binding:"required"`
	Secret     string `json:"secret"`
	Events     string `json:"events" binding:"required"`
	IsActive   bool   `json:"is_active"`
	PipelineID uint   `json:"pipeline_id" binding:"required"`
}

// UpdateWebhook 更新webhook请求参数
type UpdateWebhook struct {
	Name     string `json:"name" binding:"required"`
	URL      string `json:"url" binding:"required"`
	Secret   string `json:"secret"`
	Events   string `json:"events" binding:"required"`
	IsActive bool   `json:"is_active"`
}
