package model

import (
	"gorm.io/gorm"
	"time"
)

// TemplateCategory 模板分类
type TemplateCategory struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Description string         `gorm:"size:500" json:"description"`
	Icon        string         `gorm:"size:255" json:"icon"`
	Order       int            `gorm:"default:0" json:"order"` // 排序顺序
	CreatorID   uint           `json:"creator_id"`
	Creator     User           `gorm:"foreignKey:CreatorID" json:"creator"`
}

// TableName 设置表名
func (TemplateCategory) TableName() string {
	return "template_categories"
}

// Template 模板
type Template struct {
	ID            uint              `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	DeletedAt     gorm.DeletedAt    `gorm:"index" json:"-"`
	Name          string            `gorm:"size:100;not null" json:"name"`
	Description   string            `gorm:"size:500" json:"description"`
	CategoryID    uint              `json:"category_id"`
	Category      TemplateCategory  `gorm:"foreignKey:CategoryID" json:"category"`
	Icon          string            `gorm:"size:255" json:"icon"`
	Tags          string            `gorm:"size:255" json:"tags"` // 逗号分隔的标签
	IsPublic      bool              `gorm:"default:true" json:"is_public"`
	DownloadCount int               `gorm:"default:0" json:"download_count"`
	CreatorID     uint              `json:"creator_id"`
	Creator       User              `gorm:"foreignKey:CreatorID" json:"creator"`
	Versions      []TemplateVersion `gorm:"foreignKey:TemplateID" json:"versions"`
}

// TableName 设置表名
func (Template) TableName() string {
	return "templates"
}

// TemplateVersion 模板版本
type TemplateVersion struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	TemplateID    uint           `json:"template_id"`
	Version       string         `gorm:"size:50;not null" json:"version"`   // 语义化版本号
	Content       string         `gorm:"type:text;not null" json:"content"` // 模板内容
	Changelog     string         `gorm:"type:text" json:"changelog"`        // 变更日志
	IsLatest      bool           `gorm:"default:false" json:"is_latest"`    // 是否为最新版本
	DownloadCount int            `gorm:"default:0" json:"download_count"`
	CreatorID     uint           `json:"creator_id"`
	Creator       User           `gorm:"foreignKey:CreatorID" json:"creator"`
}

// TableName 设置表名
func (TemplateVersion) TableName() string {
	return "template_versions"
}
