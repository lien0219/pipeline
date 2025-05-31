package request

// CreatePipeline 创建流水线请求参数
type CreatePipeline struct {
	Name        string  `json:"name" binding:"required,min=2,max=100"`
	Description string  `json:"description"`
	GitRepo     string  `json:"git_repo" binding:"required"`
	GitBranch   string  `json:"git_branch" default:"main"`
	Stages      []Stage `json:"stages" binding:"required,min=1"`
}

// Stage 阶段请求参数
type Stage struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description"`
	Order       int    `json:"order"`
	Jobs        []Job  `json:"jobs" binding:"required,min=1"`
}

// Job 作业请求参数
type Job struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description"`
	Command     string `json:"command" binding:"required"`
	Image       string `json:"image"`
	Timeout     int    `json:"timeout" default:"3600"`
}

// UpdatePipeline 更新流水线请求参数
type UpdatePipeline struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description"`
	GitRepo     string `json:"git_repo" binding:"required"`
	GitBranch   string `json:"git_branch" default:"main"`
	Status      string `json:"status"`
}

// TriggerPipeline 触发流水线请求参数
type TriggerPipeline struct {
	GitBranch string `json:"git_branch"`
}

// GetDashboardStatsRequest 获取仪表盘统计数据请求参数
type GetDashboardStatsRequest struct {
	// 根据实际需求添加参数
}

// GetRecentActivitiesRequest 获取最近活动数据请求参数
type GetRecentActivitiesRequest struct {
	Limit int `json:"limit" binding:"required,min=1,max=100"` // 限制返回的活动数量
}
