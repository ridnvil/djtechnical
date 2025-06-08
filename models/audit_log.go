package models

import "time"

type AuditLog struct {
	ID          uint   `gorm:"primaryKey"`
	Method      string `gorm:"type:varchar(10);not null;check:method IN ('GET', 'POST', 'PUT', 'DELETE')"`
	ActionType  string `gorm:"type:varchar(50);not null"`
	PerformedBy *uint
	RequestIP   string `gorm:"type:inet"`
	RequestID   string `gorm:"type:uuid"`
	CreatedAt   time.Time

	User *User `gorm:"foreignKey:PerformedBy"`
}
