package utils_test

import (
	"DeallsJobsTest/utils"
	"testing"
	"time"
)

func TestGetWorkingDaysInMonth(t *testing.T) {
	tests := []struct {
		name     string
		year     int
		month    time.Month
		expected int
	}{
		{
			name:     "January 2023",
			year:     2023,
			month:    time.January,
			expected: 22,
		},
		{
			name:     "February 2023",
			year:     2023,
			month:    time.February,
			expected: 20,
		},
		{
			name:     "February 2024 (Leap Year)",
			year:     2024,
			month:    time.February,
			expected: 21,
		},
		{
			name:     "December 2023",
			year:     2023,
			month:    time.December,
			expected: 21,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.GetWorkingDaysInMonth(tt.year, tt.month)
			if result != tt.expected {
				t.Errorf("GetWorkingDaysInMonth(%d, %v) = %d; want %d",
					tt.year, tt.month, result, tt.expected)
			}
		})
	}
}
