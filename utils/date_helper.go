package utils

import "time"

func GetWorkingDaysInMonth(year int, month time.Month) int {
	startDate := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, -1) // hari terakhir bulan itu

	workingDays := 0
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		if d.Weekday() >= time.Monday && d.Weekday() <= time.Friday {
			workingDays++
		}
	}
	return workingDays
}
