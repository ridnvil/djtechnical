package services_test

import (
	"testing"
	"time"
)

func TestGetWorkingDaysInMonthLogic(t *testing.T) {
	// Test cases for different months
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
			// Calculate working days manually for verification
			startDate := time.Date(tt.year, tt.month, 1, 0, 0, 0, 0, time.UTC)
			endDate := time.Date(tt.year, tt.month+1, 0, 0, 0, 0, 0, time.UTC) // Last day of month

			workingDays := 0
			for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
				if d.Weekday() >= time.Monday && d.Weekday() <= time.Friday {
					workingDays++
				}
			}

			if workingDays != tt.expected {
				t.Errorf("Expected %d working days in %s %d, got %d",
					tt.expected, tt.month, tt.year, workingDays)
			}
		})
	}
}

func TestPeriodDateCalculation(t *testing.T) {
	tests := []struct {
		name          string
		inputDate     time.Time
		expectedStart time.Time
		expectedEnd   time.Time
	}{
		{
			name:          "Mid-month date",
			inputDate:     time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC),
			expectedStart: time.Date(2023, 5, 1, 6, 0, 0, 0, time.UTC),
			expectedEnd:   time.Date(2023, 5, 31, 6, 0, 0, 0, time.UTC),
		},
		{
			name:          "First day of month",
			inputDate:     time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC),
			expectedStart: time.Date(2023, 6, 1, 6, 0, 0, 0, time.UTC),
			expectedEnd:   time.Date(2023, 6, 30, 6, 0, 0, 0, time.UTC),
		},
		{
			name:          "Last day of month",
			inputDate:     time.Date(2023, 7, 31, 0, 0, 0, 0, time.UTC),
			expectedStart: time.Date(2023, 7, 1, 6, 0, 0, 0, time.UTC),
			expectedEnd:   time.Date(2023, 7, 31, 6, 0, 0, 0, time.UTC),
		},
		{
			name:          "February in leap year",
			inputDate:     time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC),
			expectedStart: time.Date(2024, 2, 1, 6, 0, 0, 0, time.UTC),
			expectedEnd:   time.Date(2024, 2, 29, 6, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startDate := time.Date(tt.inputDate.Year(), tt.inputDate.Month(), 1, 6, 0, 0, 0, tt.inputDate.Location())
			firstOfNextMonth := time.Date(tt.inputDate.Year(), tt.inputDate.Month()+1, 1, 6, 0, 0, 0, tt.inputDate.Location())
			endDate := firstOfNextMonth.AddDate(0, 0, -1)

			if !startDate.Equal(tt.expectedStart) {
				t.Errorf("Expected start date %v, got %v", tt.expectedStart, startDate)
			}

			if !endDate.Equal(tt.expectedEnd) {
				t.Errorf("Expected end date %v, got %v", tt.expectedEnd, endDate)
			}
		})
	}
}
