package utils_test

import (
	"DeallsJobsTest/utils"
	"testing"
)

func TestRandomSalary(t *testing.T) {
	min := 5000000
	max := 10000000

	for i := 0; i < 100; i++ {
		salary := utils.RandomSalary(min, max)

		if salary < float64(min) || salary > float64(max) {
			t.Errorf("RandomSalary(%d, %d) = %f; want value between %d and %d",
				min, max, salary, min, max)
		}
	}
}

func TestPaginate(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		pageSize int
	}{
		{
			name:     "Normal pagination",
			page:     2,
			pageSize: 10,
		},
		{
			name:     "Page less than 1",
			page:     0,
			pageSize: 10,
		},
		{
			name:     "PageSize greater than 100",
			page:     1,
			pageSize: 150,
		},
		{
			name:     "PageSize less than or equal to 0",
			page:     1,
			pageSize: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paginateFunc := utils.Paginate(tt.page, tt.pageSize)

			if paginateFunc == nil {
				t.Errorf("Paginate(%d, %d) returned nil", tt.page, tt.pageSize)
			}
		})
	}
}
