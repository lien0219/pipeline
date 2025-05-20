package model

import (
	"gorm.io/gorm"
	"time"
)

// Job 作业模型
type Job struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Description string         `gorm:"size:500" json:"description"`
	Command     string         `gorm:"type:text;not null" json:"command"`
	Image       string         `gorm:"size:255" json:"image"`
	Timeout     int            `gorm:"default:3600" json:"timeout"` // 超时时间(秒)
	StageID     uint           `json:"stage_id"`
}

// TableName 设置表名
func (Job) TableName() string {
	return "jobs"
}
