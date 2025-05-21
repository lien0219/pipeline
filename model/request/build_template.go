package request

// CreateBuildTemplate 创建构建模板请求参数
type CreateBuildTemplate struct {
	Name        string `json:"name" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Description string `json:"description"`
	Content     string `json:"content" binding:"required"`
	IsPublic    bool   `json:"is_public"`
}

// UpdateBuildTemplate 更新构建模板请求参数
type UpdateBuildTemplate struct {
	Name        string `json:"name" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Description string `json:"description"`
	Content     string `json:"content" binding:"required"`
	IsPublic    bool   `json:"is_public"`
}

// ApplyBuildTemplate 应用构建模板请求参数
type ApplyBuildTemplate struct {
	PipelineID uint   `json:"pipeline_id"`
	Name       string `json:"name"`
	GitRepo    string `json:"git_repo"`
	GitBranch  string `json:"git_branch"`
}
