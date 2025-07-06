package model

import (
	"time"

	"gorm.io/gorm"
)

// Setting 系统配置表
type Setting struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Key       string         `gorm:"uniqueIndex;size:100" json:"key"` // 配置键
	Value     string         `gorm:"type:text" json:"value"`          // 配置值
	Type      string         `gorm:"size:50" json:"type"`             // 配置类型(basic,email,integration,system)
	Desc      string         `gorm:"size:255" json:"desc"`            // 配置描述
}

// TableName 设置表名
func (Setting) TableName() string {
	return "settings"
}
