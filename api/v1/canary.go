package v1

import (
	"gin_pipeline/global"
	"gin_pipeline/model"
	"gin_pipeline/model/request"
	"gin_pipeline/model/response"
	"gin_pipeline/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var canaryService = new(service.CanaryService)

// CreateCanaryRelease 创建金丝雀发布
// @Summary 创建金丝雀发布
// @Description 创建新的金丝雀发布
// @Tags 金丝雀发布
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body request.CreateCanaryRelease true "金丝雀发布信息"
// @Success 200 {object} response.Response{data=model.CanaryRelease} "创建成功"
// @Router /canary [post]
func CreateCanaryRelease(c *gin.Context) {
	var req request.CreateCanaryRelease
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 从上下文获取用户ID
	userID := c.GetUint("userId")
	if userID == 0 {
		response.FailWithMessage("创建金丝雀发布失败", c)
		return
	}

	canary := model.CanaryRelease{
		Name:            req.Name,
		PipelineID:      req.PipelineID,
		PipelineRunID:   req.PipelineRunID,
		TargetNamespace: req.TargetNamespace,
		TargetService:   req.TargetService,
		TrafficPercent:  req.TrafficPercent,
		Status:          "pending",
		CreatedBy:       userID,
	}

	if err := canaryService.CreateCanaryRelease(&canary); err != nil {
		global.Log.Error("创建金丝雀发布失败", zap.Error(err))
		response.FailWithMessage("创建金丝雀发布失败", c)
		return
	}

	response.OkWithData(canary, c)
}

// DeployCanaryRelease 部署金丝雀发布
// @Summary 部署金丝雀发布
// @Description 部署指定的金丝雀发布
// @Tags 金丝雀发布
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "金丝雀发布ID"
// @Success 200 {object} response.Response "部署成功"
// @Router /canary/{id}/deploy [post]
func DeployCanaryRelease(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		global.Log.Error("无效的金丝雀发布ID", zap.Error(err))
		response.FailWithMessage("无效的金丝雀发布ID", c)
		return
	}

	if err := canaryService.DeployCanaryRelease(uint(id)); err != nil {
		global.Log.Error("部署金丝雀发布失败", zap.Error(err))
		response.FailWithMessage("部署金丝雀发布失败", c)
		return
	}

	response.OkWithMessage("部署金丝雀发布成功", c)
}
