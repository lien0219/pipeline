package model

import (
	"gorm.io/gorm"
	"time"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Username  string         `gorm:"size:50;not null;unique" json:"username"`
	Password  string         `gorm:"size:100;not null" json:"-"`
	Name      string         `gorm:"size:50" json:"name"`
	Email     string         `gorm:"size:100" json:"email"`
	Phone     string         `gorm:"size:20" json:"phone"`
	Avatar    string         `gorm:"size:255" json:"avatar"`
	Role      string         `gorm:"size:20;default:user" json:"role"` // admin, user
	LastLogin time.Time      `json:"last_login"`
}

// TableName 设置表名
func (User) TableName() string {
	return "users"
}
