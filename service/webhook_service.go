package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gin_pipeline/global"
	"gin_pipeline/model"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// WebhookService webhook服务
type WebhookService struct{}

type WebhookTriggerHistory struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	WebhookID uint      `json:"webhook_id"`
	Status    string    `json:"status"`
	Response  string    `json:"response"`
	Duration  int64     `json:"duration"`
}

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

	startTime := time.Now()

	// 记录触发历史
	history := WebhookTriggerHistory{
		WebhookID: webhookID,
		Status:    "pending",
	}
	if err := global.DB.Create(&history).Error; err != nil {
		global.Log.Error("创建webhook触发历史失败", zap.Error(err))
	}

	// 这里实现实际的webhook触发逻辑
	err := s.actualTriggerWebhook(webhook, payload)
	duration := time.Since(startTime).Milliseconds()

	// 更新触发历史状态
	update := map[string]interface{}{
		"status":   "success",
		"duration": duration,
	}
	if err != nil {
		update["status"] = "failed"
		update["response"] = err.Error()
	}
	global.DB.Model(&WebhookTriggerHistory{}).Where("id = ?", history.ID).Updates(update)

	global.Log.Info("触发webhook",
		zap.Uint("webhookID", webhookID),
		zap.String("url", webhook.URL),
		zap.Int64("duration", duration),
		zap.Error(err),
	)

	return err
}
func (s *WebhookService) actualTriggerWebhook(webhook model.Webhook, payload interface{}) error {
	// 1. 准备HTTP请求
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// 2. 序列化payload
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("序列化payload失败: %v", err)
	}

	// 3. 创建HTTP请求
	req, err := http.NewRequest("POST", webhook.URL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	// 4. 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Pipeline-Webhook-Sender/1.0")

	// 5. 添加签名头(如果配置了secret)
	if webhook.Secret != "" {
		hmac := hmac.New(sha256.New, []byte(webhook.Secret))
		hmac.Write(jsonPayload)
		signature := hex.EncodeToString(hmac.Sum(nil))
		req.Header.Set("X-Hub-Signature-256", "sha256="+signature)
	}

	// 6. 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("请求发送失败: %v", err)
	}
	defer resp.Body.Close()

	// 7. 检查响应状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("webhook返回错误状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	return nil
}
