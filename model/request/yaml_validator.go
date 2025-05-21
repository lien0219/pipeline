package request

// ValidateYAML YAML验证请求参数
type ValidateYAML struct {
	Name       string `json:"name" binding:"required,min=2,max=100"`
	Content    string `json:"content" binding:"required"`
	SchemaType string `json:"schema_type" binding:"required"` // kubernetes, custom, etc.
	SchemaName string `json:"schema_name"`                    // 例如：Deployment, Service, etc.
}

// CreateYAMLSchema 创建YAML Schema请求参数
type CreateYAMLSchema struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Type        string `json:"type" binding:"required"` // kubernetes, custom, etc.
	Version     string `json:"version" binding:"required"`
	Schema      string `json:"schema" binding:"required"` // JSON Schema内容
	Description string `json:"description"`
}

// UpdateYAMLSchema 更新YAML Schema请求参数
type UpdateYAMLSchema struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Version     string `json:"version" binding:"required"`
	Schema      string `json:"schema" binding:"required"` // JSON Schema内容
	Description string `json:"description"`
}
