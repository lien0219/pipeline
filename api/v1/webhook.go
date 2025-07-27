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

// GetWebhooks 获取webhooks（支持分页和按流水线ID筛选）
// @Summary 获取webhooks
// @Description 获取所有webhooks，支持分页和按流水线ID筛选
// @Tags Webhook管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param pipelineId query int false "流水线ID（可选）"
// @Param page query int false "页码，默认1"
// @Param pageSize query int false "每页条数，默认10"
// @Param name query string false "Webhook名称（可选）"
// @Param status query bool false "状态（可选）"
// @Success 200 {object} response.Response{data=response.PageResult{list=[]model.Webhook}}
// @Router /webhook [get]
func GetWebhooks(c *gin.Context) {
	pipelineIDStr := c.Query("pipelineId")
	var pipelineID uint = 0
	if pipelineIDStr != "" {
		pid, err := strconv.ParseUint(pipelineIDStr, 10, 32)
		if err != nil {
			global.Log.Error("无效的流水线ID", zap.Error(err))
			response.FailWithMessage("无效的流水线ID", c)
			return
		}
		pipelineID = uint(pid)
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	name := c.Query("name")
	statusStr := c.Query("status")
	var status *bool = nil
	if statusStr != "" {
		statusVal, _ := strconv.ParseBool(statusStr)
		status = &statusVal
	}

	webhooks, total, err := webhookService.GetWebhooksWithPage(pipelineID, name, status, page, pageSize)
	if err != nil {
		global.Log.Error("获取webhooks失败", zap.Error(err))
		response.FailWithMessage("获取webhooks失败", c)
		return
	}

	response.OkWithData(response.PageResult{
		List:     webhooks,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, c)
}
