package models

import "time"

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"size:50;unique;not null"`
	Password string `gorm:"not null" json:"-"` // should be hashed
	FullName string `gorm:"size:100"`
	Role     string `gorm:"type:varchar(10);not null;check:role IN ('admin','employee')"`

	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy *uint
	UpdatedBy *uint
	RequestIP string `gorm:"type:inet"`
	RequestID string `gorm:"type:uuid"`
}
