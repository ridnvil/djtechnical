package services

import (
	"DeallsJobsTest/models"
	"errors"
	"gorm.io/gorm"
	"time"
)

func GetPeriodData(date time.Time, db *gorm.DB) (models.AttendancesPeriod, error) {
	var pariod models.AttendancesPeriod
	startDate := time.Date(date.Year(), date.Month(), 1, 6, 0, 0, 0, date.Location())
	firstOfNextMonth := time.Date(date.Year(), date.Month()+1, 1, 6, 0, 0, 0, date.Location())
	endDate := firstOfNextMonth.AddDate(0, 0, -1)
	if err := db.Where("start_date <= ? AND end_date >= ?", endDate, startDate).First(&pariod).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.AttendancesPeriod{}, nil // No overlap found
		}
		return models.AttendancesPeriod{}, err // Some other error occurred
	}
	return pariod, nil
}
