package services

import (
	"DeallsJobsTest/models"
	"DeallsJobsTest/utils"
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"time"
)

func GeneratePayroll(db *gorm.DB, periodID uint, ctx *fiber.Ctx, client *redis.Client, requestID string) error {
	var users []models.User
	if err := db.Where("role = ?", "employee").Find(&users).Error; err != nil {
		return err
	}

	var dataSending struct {
		PeriodID  uint          `json:"period_id"`
		IP        string        `json:"ip"`
		Users     []models.User `json:"users"`
		RequestID string        `json:"request_id"`
	}

	dataSending.PeriodID = periodID
	dataSending.IP = ctx.IP()
	dataSending.Users = users
	dataSending.RequestID = requestID

	channel := "payroll_channel_employees"
	jsonData, err := json.Marshal(dataSending)
	if err != nil {
		return err
	}
	if errpub := client.Publish(context.Background(), channel, jsonData); errpub != nil {
		return errpub.Err()
	}
	return nil
}

func GetPayrollDataPagination(db *gorm.DB, page int, pageSize int, sort string, periodID uint) (models.PayRollResponse, error) {
	var payroll models.PayRollResponse
	var payslips []models.Payslip
	if err := db.Scopes(utils.Paginate(page, pageSize)).Preload("User").Preload("Period").Where("period_id = ?", periodID).Order("created_at " + sort).Find(&payslips).Error; err != nil {
		return models.PayRollResponse{}, err
	}

	var totalCount int64
	if err := db.Table("users").Select("COUNT(*) AS totalCount").Where("role = ?", "employee").Scan(&totalCount).Error; err != nil {
		return models.PayRollResponse{}, err
	}
	payroll.TotalEmployees = int(totalCount)

	var periodDate time.Time
	if err := db.Table("attendance_periods").Select("end_date").Where("id = ?", periodID).Scan(&periodDate).Error; err != nil {
		return models.PayRollResponse{}, err
	}

	payroll.Period = utils.ConvertMonthToIDString(periodDate)

	for _, slip := range payslips {
		payroll.Payslips = append(payroll.Payslips, models.PayslipResponse{
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
		payroll.TotalTHP += slip.TakeHomePay
	}
	return payroll, nil
}

func GetPayrollAll(db *gorm.DB, sort string, periodID uint) (models.PayRollResponse, error) {
	var payroll models.PayRollResponse
	var payslips []models.Payslip
	if err := db.Preload("User").Preload("Period").Where("period_id = ?", periodID).Order("created_at " + sort).Find(&payslips).Error; err != nil {
		return models.PayRollResponse{}, err
	}

	var totalCount int64
	if err := db.Table("users").Select("COUNT(*) AS totalCount").Where("role = ?", "employee").Scan(&totalCount).Error; err != nil {
		return models.PayRollResponse{}, err
	}
	payroll.TotalEmployees = int(totalCount)

	var periodDate time.Time
	if err := db.Table("attendance_periods").Select("end_date").Where("id = ?", periodID).Scan(&periodDate).Error; err != nil {
		return models.PayRollResponse{}, err
	}

	payroll.Period = utils.ConvertMonthToIDString(periodDate)

	for _, slip := range payslips {
		payroll.Payslips = append(payroll.Payslips, models.PayslipResponse{
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
		payroll.TotalTHP += slip.TakeHomePay
	}
	return payroll, nil
}
