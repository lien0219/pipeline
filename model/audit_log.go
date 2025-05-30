package model

import "time"

type AuditLog struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UserID     uint      `json:"user_id"`
	Action     string    `json:"action" gorm:"size:100"`
	ObjectType string    `json:"object_type" gorm:"size:50"`
	ObjectID   uint      `json:"object_id"`
	RequestIP  string    `json:"request_ip" gorm:"size:50"`
	Details    string    `json:"details" gorm:"type:text"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
