package model

import (
	"gorm.io/gorm"
	"time"
)

// YAMLValidation YAML验证记录
type YAMLValidation struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	Name       string         `gorm:"size:100;not null" json:"name"`
	Content    string         `gorm:"type:text;not null" json:"content"`   // YAML内容
	SchemaType string         `gorm:"size:50;not null" json:"schema_type"` // kubernetes, custom, etc.
	SchemaName string         `gorm:"size:100" json:"schema_name"`         // 例如：Deployment, Service, etc.
	IsValid    bool           `json:"is_valid"`
	Errors     string         `gorm:"type:text" json:"errors"` // 验证错误信息
	CreatorID  uint           `json:"creator_id"`
	Creator    User           `gorm:"foreignKey:CreatorID" json:"creator"`
}

// TableName 设置表名
func (YAMLValidation) TableName() string {
	return "yaml_validations"
}

// YAMLSchema YAML Schema定义
type YAMLSchema struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Type        string         `gorm:"size:50;not null" json:"type"` // kubernetes, custom, etc.
	Version     string         `gorm:"size:50;not null" json:"version"`
	Schema      string         `gorm:"type:text;not null" json:"schema"` // JSON Schema内容
	Description string         `gorm:"size:500" json:"description"`
	CreatorID   uint           `json:"creator_id"`
	Creator     User           `gorm:"foreignKey:CreatorID" json:"creator"`
}

// TableName 设置表名
func (YAMLSchema) TableName() string {
	return "yaml_schemas"
}
