package model

import (
	"gorm.io/gorm"
	"time"
)

// Release 发布模型
type Release struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Version      string         `gorm:"size:50;not null" json:"version"`
	Description  string         `gorm:"size:500" json:"description"`
	ReleaseNotes string         `gorm:"type:text" json:"release_notes"`
	Status       string         `gorm:"size:20;default:pending" json:"status"` // pending, in_progress, success, failed, rolled_back
	Environment  string         `gorm:"size:50;not null" json:"environment"`   // development, testing, staging, production
	ArtifactID   uint           `json:"artifact_id"`
	Artifact     Artifact       `gorm:"foreignKey:ArtifactID" json:"artifact"`
	DeployedAt   time.Time      `json:"deployed_at"`
	DeployedBy   uint           `json:"deployed_by"`
	User         User           `gorm:"foreignKey:DeployedBy" json:"user"`
	IsRollback   bool           `gorm:"default:false" json:"is_rollback"`
}

// TableName 设置表名
func (Release) TableName() string {
	return "releases"
}
