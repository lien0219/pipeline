package model

import (
	"time"
)

// PipelineActivity 流水线活动记录
type PipelineActivity struct {
	ID        uint      `gorm:"primaryKey"`
	Type      string    `gorm:"column:type"`
	Content   string    `gorm:"column:content"`
	Timestamp time.Time `gorm:"column:timestamp"`
	Hollow    bool      `gorm:"column:hollow"`
}

// TableName 指定表名
func (PipelineActivity) TableName() string {
	return "pipeline_activities"
}
