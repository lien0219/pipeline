package request

// CreateRelease 创建发布请求参数
type CreateRelease struct {
	Version      string `json:"version" binding:"required"`
	Description  string `json:"description"`
	ReleaseNotes string `json:"release_notes"`
	Environment  string `json:"environment" binding:"required"`
	ArtifactID   uint   `json:"artifact_id" binding:"required"`
}

// UpdateRelease 更新发布请求参数
type UpdateRelease struct {
	Version      string `json:"version" binding:"required"`
	Description  string `json:"description"`
	ReleaseNotes string `json:"release_notes"`
	Status       string `json:"status"`
}
