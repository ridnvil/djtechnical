package models

import "time"

type Payslip struct {
	ID                 uint    `gorm:"primaryKey"`
	UserID             uint    `gorm:"not null"`
	PeriodID           uint    `gorm:"not null"`
	BaseSalary         float64 `gorm:"type:numeric(15,2);not null"`
	WorkingDays        int     `gorm:"not null"`
	AttendedDays       int     `gorm:"not null"`
	ProratedSalary     float64 `gorm:"type:numeric(15,2);not null"`
	OvertimeHours      float64 `gorm:"type:numeric(5,2);not null"`
	OvertimePay        float64 `gorm:"type:numeric(15,2);not null"`
	ReimbursementTotal float64 `gorm:"type:numeric(15,2);not null"`
	TakeHomePay        float64 `gorm:"type:numeric(15,2);not null"`

	GeneratedAt time.Time
	GeneratedBy *uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CreatedBy   *uint
	UpdatedBy   *uint
	RequestIP   string  `gorm:"type:inet"`
	RequestID   *string `gorm:"type:uuid"`

	User   User              `gorm:"foreignKey:UserID"`
	Period AttendancesPeriod `gorm:"foreignKey:PeriodID"`
}
