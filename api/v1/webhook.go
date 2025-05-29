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

var webhookService = new(service.WebhookService)

// CreateWebhook 创建webhook
// @Summary 创建webhook
// @Description 创建新的webhook
// @Tags Webhook管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body request.CreateWebhook true "webhook信息"
// @Success 200 {object} response.Response{data=model.Webhook} "创建成功"
// @Router /webhook [post]
func CreateWebhook(c *gin.Context) {
	var req request.CreateWebhook
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 从上下文获取用户ID
	userID := c.GetUint("userId")
	if userID == 0 {
		response.FailWithMessage("创建webhook失败", c)
		return
	}

	webhook := model.Webhook{
		Name:       req.Name,
		URL:        req.URL,
		Secret:     req.Secret,
		Events:     req.Events,
		IsActive:   req.IsActive,
		PipelineID: req.PipelineID,
		CreatedBy:  userID,
	}

	if err := webhookService.CreateWebhook(&webhook); err != nil {
		global.Log.Error("创建webhook失败", zap.Error(err))
		response.FailWithMessage("创建webhook失败", c)
		return
	}

	response.OkWithData(webhook, c)
}

// GetWebhooksByPipelineID 获取流水线的webhooks
// @Summary 获取流水线的webhooks
// @Description 获取指定流水线的所有webhooks
// @Tags Webhook管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param pipelineId path int true "流水线ID"
// @Success 200 {object} response.Response{data=[]model.Webhook} "获取成功"
// @Router /webhook/pipeline/{pipelineId} [get]
func GetWebhooksByPipelineID(c *gin.Context) {
	pipelineIDStr := c.Param("pipelineId")
	pipelineID, err := strconv.ParseUint(pipelineIDStr, 10, 32)
	if err != nil {
		global.Log.Error("无效的流水线ID", zap.Error(err))
		response.FailWithMessage("无效的流水线ID", c)
		return
	}

	webhooks, err := webhookService.GetWebhooksByPipelineID(uint(pipelineID))
	if err != nil {
		global.Log.Error("获取webhooks失败", zap.Error(err))
		response.FailWithMessage("获取webhooks失败", c)
		return
	}

	response.OkWithData(webhooks, c)
}
