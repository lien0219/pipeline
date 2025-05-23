package model

import (
	"time"

	"gorm.io/gorm"
)

// HPAPolicy HPA 策略模型
type HPAPolicy struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	Name            string         `gorm:"size:255;not null" json:"name"`
	Namespace       string         `gorm:"size:255;not null" json:"namespace"`
	Deployment      string         `gorm:"size:255;not null" json:"deployment"`
	MinReplicas     int            `json:"min_replicas"`
	MaxReplicas     int            `json:"max_replicas"`
	CPUThreshold    int            `json:"cpu_threshold"`    // CPU 使用率阈值百分比
	MemoryThreshold int            `json:"memory_threshold"` // 内存使用率阈值百分比
}

// TableName 设置表名
func (HPAPolicy) TableName() string {
	return "hpa_policies"
}
