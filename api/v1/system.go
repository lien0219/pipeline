package v1

import (
	"gin_pipeline/model/response"
	"gin_pipeline/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

var systemService = service.SystemService{}
var logService = service.LogService{}

// GetSystemStatus 获取系统状态
// @Tags 系统维护
// @Summary 获取系统状态
// @Security ApiKeyAuth
// @Produce application/json
// @Success 200 {object} response.Response{data=service.SystemStatus}
// @Router /system/status [get]
func GetSystemStatus(c *gin.Context) {
	status, err := systemService.GetSystemStatus()
	if err != nil {
		response.FailWithMessage("获取系统状态失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(status, "获取系统状态成功", c)
}

// GetLogFiles 获取日志文件列表
// @Tags 系统维护
// @Summary 获取日志文件列表
// @Security ApiKeyAuth
// @Produce application/json
// @Success 200 {object} response.Response{data=[]service.LogFile}
// @Router /system/logs [get]
func GetLogFiles(c *gin.Context) {
	files, err := logService.GetLogFiles()
	if err != nil {
		response.FailWithMessage("获取日志文件失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(files, "获取日志文件成功", c)
}

// GetLogContent 获取日志文件内容
// @Tags System
// @Summary 获取日志内容
// @Security ApiKeyAuth
// @Produce  application/json
// @Param name path string true "日志文件名"
// @Param startLine query int false "起始行号，默认0"
// @Param limit query int false "获取行数，默认100"
// @Success 200 {object} response.Response{data=map[string]interface{}}
// @Router /system/logs/{name} [get]
// 修复：移除结构体接收者`s *SystemApi`
func GetLogContent(c *gin.Context) {
	fileName := c.Param("name")
	startLine, _ := strconv.Atoi(c.DefaultQuery("startLine", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))

	if limit > 1000 {
		limit = 1000 // 限制最大行数
	}

	logService := service.LogService{}
	content, totalLines, err := logService.GetLogContent(fileName, startLine, limit)
	if err != nil {
		response.FailWithMessage("获取日志失败"+err.Error(), c)
		return
	}

	response.OkWithData(map[string]interface{}{
		"content":    content,
		"totalLines": totalLines,
		"startLine":  startLine,
		"limit":      limit,
	}, c)
}
