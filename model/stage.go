package model

import (
	"gorm.io/gorm"
	"time"
)

// Stage 阶段模型
type Stage struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Description string         `gorm:"size:500" json:"description"`
	Order       int            `gorm:"not null" json:"order"`
	PipelineID  uint           `json:"pipeline_id"`
	Jobs        []Job          `gorm:"foreignKey:StageID" json:"jobs"`
}

// TableName 设置表名
func (Stage) TableName() string {
	return "stages"
}
