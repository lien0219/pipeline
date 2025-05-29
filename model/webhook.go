package model

import (
	"time"

	"gorm.io/gorm"
)

// Webhook webhook模型
type Webhook struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	Name       string         `gorm:"size:255;not null" json:"name"`
	URL        string         `gorm:"size:500;not null" json:"url"`
	Secret     string         `gorm:"size:255" json:"-"`
	Events     string         `gorm:"size:500" json:"events"` // 逗号分隔的事件列表
	IsActive   bool           `gorm:"default:true" json:"is_active"`
	PipelineID uint           `json:"pipeline_id"`
	Pipeline   Pipeline       `gorm:"foreignKey:PipelineID" json:"pipeline"`
	CreatedBy  uint           `json:"created_by"`
	User       User           `gorm:"foreignKey:CreatedBy" json:"user"`
}

// TableName 设置表名
func (Webhook) TableName() string {
	return "webhooks"
}
