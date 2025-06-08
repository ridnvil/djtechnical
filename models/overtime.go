package models

import "time"

type Overtime struct {
	ID       uint      `gorm:"primaryKey"`
	UserID   uint      `gorm:"not null"`
	Date     time.Time `gorm:"not null"`
	Hours    float64   `gorm:"type:numeric(4,2);not null"`
	PeriodID uint      `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy *uint
	UpdatedBy *uint
	RequestIP string `gorm:"type:inet"`

	User   User              `gorm:"foreignKey:UserID"`
	Period AttendancesPeriod `gorm:"foreignKey:PeriodID"`
}

type OvertimePaid struct {
	ID          uint    `gorm:"primaryKey"`
	Amount      float64 `gorm:"type:numeric(10,2);not null"`
	Description string  `gorm:"type:text;not null"`
	PeriodID    uint    `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Period AttendancesPeriod `gorm:"foreignKey:PeriodID"`
}
