package request

// CreateArtifact 创建制品请求参数
type CreateArtifact struct {
	Name          string `json:"name" binding:"required"`
	Type          string `json:"type" binding:"required"`
	Version       string `json:"version"`
	Description   string `json:"description"`
	PipelineID    uint   `json:"pipeline_id"`
	PipelineRunID uint   `json:"pipeline_run_id"`
}

// UpdateArtifact 更新制品请求参数
type UpdateArtifact struct {
	Name        string `json:"name" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Version     string `json:"version"`
	Description string `json:"description"`
}
