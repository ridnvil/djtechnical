package models

type PayRollResponse struct {
	Payslips []PayslipResponse `json:"payslips"`
	TotalTHP float64           `json:"total_thp"`
}

type PayslipResponse struct {
	ID                 uint    `json:"id"`
	UserID             uint    `json:"user_id"`
	FullName           string  `json:"full_name"`
	PeriodID           uint    `json:"period_id"`
	BaseSalary         float64 `json:"base_salary"`
	WorkingDays        int     `json:"working_days"`
	AttendedDays       int     `json:"attended_days"`
	ProratedSalary     float64 `json:"prorated_salary"`
	OvertimeHours      float64 `json:"overtime_hours"`
	OvertimePay        float64 `json:"overtime_pay"`
	ReimbursementTotal float64 `json:"reimbursement_total"`
	TakeHomePay        float64 `json:"take_home_pay"`
}
