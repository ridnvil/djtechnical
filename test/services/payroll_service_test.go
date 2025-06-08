package services_test

import (
	"DeallsJobsTest/models"
	"reflect"
	"testing"
)

func TestPayrollResponseConstruction(t *testing.T) {
	payslips := []models.Payslip{
		{
			ID:                 1,
			UserID:             101,
			PeriodID:           201,
			BaseSalary:         5000000,
			WorkingDays:        22,
			AttendedDays:       20,
			ProratedSalary:     4545454.55,
			OvertimeHours:      5,
			OvertimePay:        250000,
			ReimbursementTotal: 100000,
			TakeHomePay:        4895454.55,
			User: models.User{
				ID:       101,
				FullName: "John Doe",
			},
			Period: models.AttendancesPeriod{
				ID: 201,
			},
		},
		{
			ID:                 2,
			UserID:             102,
			PeriodID:           201,
			BaseSalary:         6000000,
			WorkingDays:        22,
			AttendedDays:       22,
			ProratedSalary:     6000000,
			OvertimeHours:      0,
			OvertimePay:        0,
			ReimbursementTotal: 200000,
			TakeHomePay:        6200000,
			User: models.User{
				ID:       102,
				FullName: "Jane Smith",
			},
			Period: models.AttendancesPeriod{
				ID: 201,
			},
		},
	}

	expected := models.PayRollResponse{
		Payslips: []models.PayslipResponse{
			{
				ID:                 1,
				UserID:             101,
				FullName:           "John Doe",
				PeriodID:           201,
				BaseSalary:         5000000,
				WorkingDays:        22,
				AttendedDays:       20,
				ProratedSalary:     4545454.55,
				OvertimeHours:      5,
				OvertimePay:        250000,
				ReimbursementTotal: 100000,
				TakeHomePay:        4895454.55,
			},
			{
				ID:                 2,
				UserID:             102,
				FullName:           "Jane Smith",
				PeriodID:           201,
				BaseSalary:         6000000,
				WorkingDays:        22,
				AttendedDays:       22,
				ProratedSalary:     6000000,
				OvertimeHours:      0,
				OvertimePay:        0,
				ReimbursementTotal: 200000,
				TakeHomePay:        6200000,
			},
		},
		TotalTHP: 11095454.55,
	}

	actual := models.PayRollResponse{}
	for _, slip := range payslips {
		actual.Payslips = append(actual.Payslips, models.PayslipResponse{
			ID:                 slip.ID,
			UserID:             slip.UserID,
			FullName:           slip.User.FullName,
			PeriodID:           slip.PeriodID,
			BaseSalary:         slip.BaseSalary,
			WorkingDays:        slip.WorkingDays,
			AttendedDays:       slip.AttendedDays,
			ProratedSalary:     slip.ProratedSalary,
			OvertimeHours:      slip.OvertimeHours,
			OvertimePay:        slip.OvertimePay,
			ReimbursementTotal: slip.ReimbursementTotal,
			TakeHomePay:        slip.TakeHomePay,
		})
		actual.TotalTHP += slip.TakeHomePay
	}

	if !reflect.DeepEqual(expected.Payslips, actual.Payslips) {
		t.Errorf("Payslips do not match.\nExpected: %+v\nActual: %+v", expected.Payslips, actual.Payslips)
	}

	if diff := expected.TotalTHP - actual.TotalTHP; diff > 0.01 || diff < -0.01 {
		t.Errorf("TotalTHP does not match. Expected: %.2f, Actual: %.2f", expected.TotalTHP, actual.TotalTHP)
	}
}

func TestPayrollPagination(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		pageSize int
		expected struct {
			page     int
			pageSize int
			offset   int
		}
	}{
		{
			name:     "Normal pagination",
			page:     2,
			pageSize: 10,
			expected: struct {
				page     int
				pageSize int
				offset   int
			}{
				page:     2,
				pageSize: 10,
				offset:   10,
			},
		},
		{
			name:     "Page less than 1",
			page:     0,
			pageSize: 10,
			expected: struct {
				page     int
				pageSize int
				offset   int
			}{
				page:     1,
				pageSize: 10,
				offset:   0,
			},
		},
		{
			name:     "PageSize greater than 100",
			page:     1,
			pageSize: 150,
			expected: struct {
				page     int
				pageSize int
				offset   int
			}{
				page:     1,
				pageSize: 100,
				offset:   0,
			},
		},
		{
			name:     "PageSize less than or equal to 0",
			page:     1,
			pageSize: 0,
			expected: struct {
				page     int
				pageSize int
				offset   int
			}{
				page:     1,
				pageSize: 10,
				offset:   0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page := tt.page
			pageSize := tt.pageSize

			if page < 1 {
				page = 1
			}

			switch {
			case pageSize > 100:
				pageSize = 100
			case pageSize <= 0:
				pageSize = 10
			}

			offset := (page - 1) * pageSize

			if page != tt.expected.page {
				t.Errorf("Page: expected %d, got %d", tt.expected.page, page)
			}

			if pageSize != tt.expected.pageSize {
				t.Errorf("PageSize: expected %d, got %d", tt.expected.pageSize, pageSize)
			}

			if offset != tt.expected.offset {
				t.Errorf("Offset: expected %d, got %d", tt.expected.offset, offset)
			}
		})
	}
}
