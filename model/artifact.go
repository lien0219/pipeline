package model

import (
	"gorm.io/gorm"
	"time"
)

// Artifact 制品模型
type Artifact struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Name          string         `gorm:"size:255;not null" json:"name"`
	Type          string         `gorm:"size:50;not null" json:"type"` // zip, jar, war, docker, etc.
	Path          string         `gorm:"size:500;not null" json:"path"`
	Size          int64          `json:"size"` // 文件大小(字节)
	Version       string         `gorm:"size:50" json:"version"`
	Description   string         `gorm:"size:500" json:"description"`
	PipelineID    uint           `json:"pipeline_id"`
	Pipeline      Pipeline       `gorm:"foreignKey:PipelineID" json:"pipeline"`
	PipelineRunID uint           `json:"pipeline_run_id"`
	DownloadCount int            `gorm:"default:0" json:"download_count"`
	CreatedBy     uint           `json:"created_by"`
	User          User           `gorm:"foreignKey:CreatedBy" json:"user"`
}

// TableName 设置表名
func (Artifact) TableName() string {
	return "artifacts"
}
