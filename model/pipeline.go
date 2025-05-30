package model

import (
	"gin_pipeline/global"
	"gin_pipeline/utils"
	"time"

	"gorm.io/gorm"
)

// Pipeline 流水线模型
type Pipeline struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Description string         `gorm:"size:500" json:"description"`
	GitRepo     string         `gorm:"size:255;not null" json:"git_repo"`
	GitBranch   string         `gorm:"size:100;default:main" json:"git_branch"`
	Status      string         `gorm:"size:20;default:inactive" json:"status"` // inactive, active, running, success, failed
	LastRunAt   *time.Time     `json:"last_run_at"`
	CreatorID   uint           `json:"creator_id"`
	Creator     User           `gorm:"foreignKey:CreatorID" json:"creator"`
	Stages      []Stage        `gorm:"foreignKey:PipelineID" json:"stages"`
	Config      string         `gorm:"type:text"`  // 加密存储
	RawConfig   string         `gorm:"-" json:"-"` // 不序列化到
}

// TableName 设置表名
func (Pipeline) TableName() string {
	return "pipelines"
}

func (p *Pipeline) BeforeSave(tx *gorm.DB) error {
	if p.RawConfig != "" {
		encrypted, err := utils.AESEncrypt(p.RawConfig, global.Config.System.EncryptKey)
		if err != nil {
			return err
		}
		p.Config = encrypted
	}
	return nil
}

func (p *Pipeline) AfterFind(tx *gorm.DB) error {
	if p.Config != "" {
		decrypted, err := utils.AESDecrypt(p.Config, global.Config.System.EncryptKey)
		if err != nil {
			return err
		}
		p.RawConfig = decrypted
	}
	return nil
}
