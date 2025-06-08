package models

import "time"

type AttendancesPeriod struct {
	ID        uint      `gorm:"primaryKey"`
	StartDate time.Time `gorm:"not null"`
	EndDate   time.Time `gorm:"not null"`
	IsLocked  bool      `gorm:"default:false"`

	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy *uint
	UpdatedBy *uint
	RequestIP string `gorm:"type:inet"`
	RequestID string `gorm:"type:uuid"`

	CreatedByUser *User `gorm:"foreignKey:CreatedBy"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedBy"`
}

func (att *AttendancesPeriod) TableName() string {
	return "attendance_periods"
}

type Attendance struct {
	ID       uint      `gorm:"primaryKey"`
	UserID   uint      `gorm:"not null"`
	Date     time.Time `gorm:"not null"`
	PeriodID uint      `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy *uint
	UpdatedBy *uint
	RequestIP string `gorm:"type:inet"`
	RequestID string `gorm:"type:uuid"`

	User   User              `gorm:"foreignKey:UserID"`
	Period AttendancesPeriod `gorm:"foreignKey:PeriodID"`
}
