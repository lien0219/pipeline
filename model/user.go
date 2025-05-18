package model

import (
	"time"
)

// User 用户模型
type User struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Username  string    `gorm:"size:32;not null;uniqueIndex" json:"username"`
	Password  string    `gorm:"size:128;not null" json:"-"` // 不序列化密码
	Email     string    `gorm:"size:128;not null;uniqueIndex" json:"email"`
	Name      string    `gorm:"size:50" json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (User) TableName() string {
	return "user"
}
