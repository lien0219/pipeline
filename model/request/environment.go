package request

// EnvironmentVariable 环境变量
type EnvironmentVariable struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value"`
}

// CreateEnvironment 创建环境请求参数
type CreateEnvironment struct {
	Name        string                `json:"name" binding:"required"`
	Type        string                `json:"type" binding:"required"`
	URL         string                `json:"url"`
	Description string                `json:"description"`
	Variables   []EnvironmentVariable `json:"variables"`
}

// UpdateEnvironment 更新环境请求参数
type UpdateEnvironment struct {
	Name        string                `json:"name" binding:"required"`
	Type        string                `json:"type" binding:"required"`
	URL         string                `json:"url"`
	Description string                `json:"description"`
	Status      string                `json:"status"`
	Variables   []EnvironmentVariable `json:"variables"`
}
