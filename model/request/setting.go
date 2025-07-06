package request

import "github.com/go-playground/validator/v10"

// SaveSettingsRequest 保存配置请求
type SaveSettingsRequest struct {
	Type  string                 `json:"type" validate:"required,oneof=basic email integration system"` // 配置类型
	Items map[string]interface{} `json:"items" validate:"required"`                                     // 配置项
}

// Validate 验证请求
func (s *SaveSettingsRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}

// GetSettingsRequest 获取配置请求
type GetSettingsRequest struct {
	Type string `form:"type" validate:"omitempty,oneof=basic email integration system all"` // 配置类型，all表示所有
}

// Validate 验证请求
func (g *GetSettingsRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(g)
}
