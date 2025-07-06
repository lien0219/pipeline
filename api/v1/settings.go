package v1

import (
	"gin_pipeline/model/request"
	"gin_pipeline/model/response"
	"gin_pipeline/service"

	"github.com/gin-gonic/gin"
)

var settingService = service.SettingService{}

// SaveSettings 保存配置
// @Tags 系统设置
// @Summary 保存系统配置
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SaveSettingsRequest true "配置类型和配置项"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"保存成功"}"
// @Router /settings [post]
func SaveSettings(c *gin.Context) {
	var req request.SaveSettingsRequest
	_ = c.ShouldBindJSON(&req)

	// 验证请求
	if err := req.Validate(); err != nil {
		response.FailWithMessage("请求参数错误: "+err.Error(), c)
		return
	}

	// 保存配置
	if err := settingService.SaveSettings(req); err != nil {
		response.FailWithMessage("保存配置失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("保存配置成功", c)
}

// GetSettings 获取配置
// @Tags 系统设置
// @Summary 获取系统配置
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param type query string false "配置类型(basic,email,integration,system,all)"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /settings [get]
func GetSettings(c *gin.Context) {
	var req request.GetSettingsRequest
	_ = c.ShouldBindQuery(&req)

	// 验证请求
	if err := req.Validate(); err != nil {
		response.FailWithMessage("请求参数错误: "+err.Error(), c)
		return
	}

	// 获取配置
	settings, err := settingService.GetSettings(req)
	if err != nil {
		response.FailWithMessage("获取配置失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(settings, "获取配置成功", c)
}
