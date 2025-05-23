package model

import (
	"time"

	"gorm.io/gorm"
)

// Cluster 集群模型
type Cluster struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	Name       string         `gorm:"size:100;not null;uniqueIndex" json:"name"` // 集群名称
	Kubeconfig string         `gorm:"type:text;not null" json:"kubeconfig"`      // kubeconfig 内容
}

// TableName 设置表名
func (Cluster) TableName() string {
	return "clusters"
}
