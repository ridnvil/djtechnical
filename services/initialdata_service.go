package services

import (
	"DeallsJobsTest/models"
	"gorm.io/gorm"
	"time"
)

func GenerateOvertimeAmount(db *gorm.DB) error {
	period, errgetper := GetPeriodData(time.Now(), db)
	if errgetper != nil {
		return errgetper
	}

	overtimePaid := &models.OvertimePaid{
		Amount:      50000,
		Description: "Overtime payment each hours for the current period",
		PeriodID:    period.ID,
	}

	if err := db.Create(overtimePaid).Error; err != nil {
		return err
	}

	return nil
}

func GetOvertimeAmount(db *gorm.DB, periodId uint) (float64, error) {
	var overtimePaid models.OvertimePaid
	if err := db.Where("period_id = ?", periodId).First(&overtimePaid).Error; err != nil {
		return 0, err
	}
	return overtimePaid.Amount, nil
}
