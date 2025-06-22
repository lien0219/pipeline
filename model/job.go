package model

import (
	"time"

	"gorm.io/gorm"
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
	Timeout     int            `gorm:"default:3600" json:"timeout"`
	StageID     uint           `json:"stage_id"`
	Status      string         `json:"status" gorm:"column:status"`
	Logs        string         `json:"logs"`
}

// TableName 设置表名
func (Job) TableName() string {
	return "jobs"
}
