package model

import (
	"time"

	"gorm.io/gorm"
)

// CanaryRelease 金丝雀发布模型
type CanaryRelease struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	Name            string         `gorm:"size:255;not null" json:"name"`
	PipelineID      uint           `json:"pipeline_id"`
	Pipeline        Pipeline       `gorm:"foreignKey:PipelineID" json:"pipeline"`
	PipelineRunID   uint           `json:"pipeline_run_id"`
	PipelineRun     PipelineRun    `gorm:"foreignKey:PipelineRunID" json:"pipeline_run"`
	TargetNamespace string         `gorm:"size:255;not null" json:"target_namespace"`
	TargetService   string         `gorm:"size:255;not null" json:"target_service"`
	TrafficPercent  int            `gorm:"default:10" json:"traffic_percent"`
	Status          string         `gorm:"size:50;default:pending" json:"status"`
	CreatedBy       uint           `json:"created_by"`
	User            User           `gorm:"foreignKey:CreatedBy" json:"user"`
}

// TableName 设置表名
func (CanaryRelease) TableName() string {
	return "canary_releases"
}
