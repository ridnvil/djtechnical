package services_test

import (
	"testing"
)

func TestProratedSalaryCalculation(t *testing.T) {
	tests := []struct {
		name           string
		attendanceDays int64
		workingDays    int
		baseSalary     float64
		expected       float64
	}{
		{
			name:           "Full attendance",
			attendanceDays: 22,
			workingDays:    22,
			baseSalary:     5000000,
			expected:       5000000,
		},
		{
			name:           "Partial attendance",
			attendanceDays: 15,
			workingDays:    22,
			baseSalary:     5000000,
			expected:       3409090.91,
		},
		{
			name:           "No attendance",
			attendanceDays: 0,
			workingDays:    22,
			baseSalary:     5000000,
			expected:       0,
		},
		{
			name:           "More attendance than working days (edge case)",
			attendanceDays: 25,
			workingDays:    22,
			baseSalary:     5000000,
			expected:       5681818.18,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prorateSalary := float64(tt.attendanceDays) / float64(tt.workingDays) * tt.baseSalary

			if diff := prorateSalary - tt.expected; diff > 0.01 || diff < -0.01 {
				t.Errorf("Expected prorated salary %.2f, got %.2f", tt.expected, prorateSalary)
			}
		})
	}
}

func TestOvertimePayCalculation(t *testing.T) {
	tests := []struct {
		name          string
		overtimeHours float64
		overtimeRate  float64
		expected      float64
	}{
		{
			name:          "No overtime",
			overtimeHours: 0,
			overtimeRate:  50000,
			expected:      0,
		},
		{
			name:          "Standard overtime",
			overtimeHours: 10,
			overtimeRate:  50000,
			expected:      500000,
		},
		{
			name:          "Fractional overtime",
			overtimeHours: 5.5,
			overtimeRate:  50000,
			expected:      275000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			overtimePay := tt.overtimeHours * tt.overtimeRate

			if overtimePay != tt.expected {
				t.Errorf("Expected overtime pay %.2f, got %.2f", tt.expected, overtimePay)
			}
		})
	}
}

func TestTakeHomePayCalculation(t *testing.T) {
	tests := []struct {
		name               string
		proratedSalary     float64
		overtimePay        float64
		reimbursementTotal float64
		expected           float64
	}{
		{
			name:               "Salary only",
			proratedSalary:     5000000,
			overtimePay:        0,
			reimbursementTotal: 0,
			expected:           5000000,
		},
		{
			name:               "Salary and overtime",
			proratedSalary:     5000000,
			overtimePay:        500000,
			reimbursementTotal: 0,
			expected:           5500000,
		},
		{
			name:               "Salary and reimbursement",
			proratedSalary:     5000000,
			overtimePay:        0,
			reimbursementTotal: 250000,
			expected:           5250000,
		},
		{
			name:               "All components",
			proratedSalary:     5000000,
			overtimePay:        500000,
			reimbursementTotal: 250000,
			expected:           5750000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			takeHomePay := tt.proratedSalary + tt.overtimePay + tt.reimbursementTotal

			if takeHomePay != tt.expected {
				t.Errorf("Expected take-home pay %.2f, got %.2f", tt.expected, takeHomePay)
			}
		})
	}
}
