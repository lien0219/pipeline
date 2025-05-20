package model

import (
	"gorm.io/gorm"
	"time"
)

// BuildTemplate 构建模板
type BuildTemplate struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Type        string         `gorm:"size:50;not null" json:"type"` // docker, nodejs, golang, java, etc.
	Description string         `gorm:"size:500" json:"description"`
	Content     string         `gorm:"type:text;not null" json:"content"` // YAML格式的模板内容
	IsPublic    bool           `gorm:"default:true" json:"is_public"`
	UsageCount  int            `gorm:"default:0" json:"usage_count"`
	CreatedBy   uint           `json:"created_by"`
	Creator     User           `gorm:"foreignKey:CreatedBy" json:"creator"`
}

// TableName 设置表名
func (BuildTemplate) TableName() string {
	return "build_templates"
}
