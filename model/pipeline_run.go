package model

import (
	"gorm.io/gorm"
	"time"
)

// PipelineRun 流水线运行记录
type PipelineRun struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	PipelineID uint           `json:"pipeline_id"`
	Pipeline   Pipeline       `gorm:"foreignKey:PipelineID" json:"pipeline"`
	Status     string         `gorm:"size:20;default:pending" json:"status"` // pending, running, success, failed, canceled
	StartTime  *time.Time     `json:"start_time"`
	EndTime    *time.Time     `json:"end_time"`
	Duration   int            `json:"duration"` // 持续时间(秒)
	GitBranch  string         `gorm:"size:100" json:"git_branch"`
	GitCommit  string         `gorm:"size:100" json:"git_commit"`
	TriggerBy  uint           `json:"trigger_by"`
	User       User           `gorm:"foreignKey:TriggerBy" json:"user"`
	Logs       string         `gorm:"type:longtext" json:"logs"`
}

// TableName 设置表名
func (PipelineRun) TableName() string {
	return "pipeline_runs"
}
