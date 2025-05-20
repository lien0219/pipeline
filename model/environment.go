package model

import (
	"gorm.io/gorm"
	"time"
)

// Environment 环境模型
type Environment struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	Name           string         `gorm:"size:100;not null" json:"name"`
	Type           string         `gorm:"size:50;not null" json:"type"` // development, testing, staging, production
	URL            string         `gorm:"size:255" json:"url"`
	Description    string         `gorm:"size:500" json:"description"`
	Status         string         `gorm:"size:20;default:active" json:"status"` // active, inactive, error
	Variables      string         `gorm:"type:text" json:"variables"`           // JSON格式的环境变量
	LastDeployedAt *time.Time     `json:"last_deployed_at"`
	CreatedBy      uint           `json:"created_by"`
	User           User           `gorm:"foreignKey:CreatedBy" json:"user"`
}

// TableName 设置表名
func (Environment) TableName() string {
	return "environments"
}
