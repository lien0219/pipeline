package request

// CreateTemplateCategory 创建模板分类请求参数
type CreateTemplateCategory struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Order       int    `json:"order"`
}

// UpdateTemplateCategory 更新模板分类请求参数
type UpdateTemplateCategory struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Order       int    `json:"order"`
}

// CreateTemplate 创建模板请求参数
type CreateTemplate struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description"`
	CategoryID  uint   `json:"category_id" binding:"required"`
	Icon        string `json:"icon"`
	Tags        string `json:"tags"`
	IsPublic    bool   `json:"is_public"`
	Version     string `json:"version" binding:"required"` // 初始版本号
	Content     string `json:"content" binding:"required"` // 模板内容
	Changelog   string `json:"changelog"`
}

// UpdateTemplate 更新模板请求参数
type UpdateTemplate struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description"`
	CategoryID  uint   `json:"category_id" binding:"required"`
	Icon        string `json:"icon"`
	Tags        string `json:"tags"`
	IsPublic    bool   `json:"is_public"`
}

// CreateTemplateVersion 创建模板版本请求参数
type CreateTemplateVersion struct {
	Version   string `json:"version" binding:"required"`
	Content   string `json:"content" binding:"required"`
	Changelog string `json:"changelog"`
}

// SearchTemplate 搜索模板请求参数
type SearchTemplate struct {
	Keyword    string `json:"keyword" form:"keyword"`
	CategoryID uint   `json:"category_id" form:"category_id"`
	Tags       string `json:"tags" form:"tags"`
}
