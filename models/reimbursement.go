package models

import "time"

type Reimbursement struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      uint      `gorm:"not null"`
	PeriodID    uint      `gorm:"not null"`
	Date        time.Time `gorm:"not null"`
	Amount      float64   `gorm:"type:numeric(15,2);not null"`
	PathFile    string    `gorm:"type:text;not null"`
	Description string    `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy *uint
	UpdatedBy *uint
	RequestIP string `gorm:"type:inet"`

	User   User              `gorm:"foreignKey:UserID"`
	Period AttendancesPeriod `gorm:"foreignKey:PeriodID"`
}
