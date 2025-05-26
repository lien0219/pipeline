package model

import (
	"time"

	"gorm.io/gorm"
)

// ResourceQuota 资源配额模型
type ResourceQuota struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	TenantID     string         `gorm:"size:255;not null;unique" json:"tenant_id"` // 租户ID
	CPUQuota     float64        `gorm:"not null" json:"cpu_quota"`                 // CPU配额
	MemoryQuota  int64          `gorm:"not null" json:"memory_quota"`              // 内存配额
	StorageQuota int64          `gorm:"not null" json:"storage_quota"`             // 存储配额
	CPUUsage     float64        `gorm:"not null" json:"cpu_usage"`                 // CPU使用量
	MemoryUsage  int64          `gorm:"not null" json:"memory_usage"`              // 内存使用量
	StorageUsage int64          `gorm:"not null" json:"storage_usage"`             // 存储使用量
}

// TableName 设置表名
func (ResourceQuota) TableName() string {
	return "resource_quotas"
}

// TenantResourceRequest 租户资源请求模型
type TenantResourceRequest struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	TenantID       string         `json:"tenant_id"`
	CPURequest     float64        `json:"cpu_request"`
	MemoryRequest  int64          `json:"memory_request"`
	StorageRequest int64          `json:"storage_request"`
	Status         string         `gorm:"size:50;default:'pending'" json:"status"` // pending, approved, rejected
}

// TableName 设置表名
func (TenantResourceRequest) TableName() string {
	return "tenant_resource_requests"
}

// ResourceReport 资源使用报告模型
type ResourceReport struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	TenantID     string    `json:"tenant_id"`
	CPUUsage     float64   `json:"cpu_usage"`
	CPUQuota     float64   `json:"cpu_quota"`
	MemoryUsage  int64     `json:"memory_usage"`
	MemoryQuota  int64     `json:"memory_quota"`
	StorageUsage int64     `json:"storage_usage"`
	StorageQuota int64     `json:"storage_quota"`
}

// TableName 设置表名
func (ResourceReport) TableName() string {
	return "resource_reports"
}
