package service

import (
	"gin_pipeline/global"
	"gin_pipeline/model"

	"go.uber.org/zap"
)

// WebhookService webhook服务
type WebhookService struct{}

// CreateWebhook 创建webhook
func (s *WebhookService) CreateWebhook(webhook *model.Webhook) error {
	return global.DB.Create(webhook).Error
}

// GetWebhooksByPipelineID 获取流水线的webhooks
func (s *WebhookService) GetWebhooksByPipelineID(pipelineID uint) ([]model.Webhook, error) {
	var webhooks []model.Webhook
	err := global.DB.Where("pipeline_id = ?", pipelineID).Find(&webhooks).Error
	return webhooks, err
}

// TriggerWebhook 触发webhook
func (s *WebhookService) TriggerWebhook(webhookID uint, payload interface{}) error {
	var webhook model.Webhook
	if err := global.DB.First(&webhook, webhookID).Error; err != nil {
		return err
	}

	// 这里实现实际的webhook触发逻辑
	// 通常使用HTTP客户端发送POST请求到webhook.URL
	// 包含payload和签名头

	global.Log.Info("触发webhook",
		zap.Uint("webhookID", webhookID),
		zap.String("url", webhook.URL),
	)

	return nil
}
