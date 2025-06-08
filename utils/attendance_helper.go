package utils

import "time"

func IsWeekend(dateStr time.Time) bool {
	weekday := dateStr.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}
