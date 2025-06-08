package services

import (
	"DeallsJobsTest/models"
	"DeallsJobsTest/utils"
	"errors"
	"gorm.io/gorm"
	"log"
	"time"
)

func GeneratePaySlipEmployees(db *gorm.DB) error {
	var employees []models.User
	if erremm := db.Find(&employees).Error; erremm != nil {
		return erremm
	}

	period, errgetper := GetPeriodData(time.Now(), db)
	if errgetper != nil {
		return errgetper
	}

	workingDays := utils.GetWorkingDaysInMonth(time.Now().Year(), time.Now().Month())
	for _, user := range employees {
		salary := utils.RandomSalary(5000000, 10000000)
		var attendanceDays int64

		if errattday := db.Model(&models.Attendance{}).Where("user_id = ? AND period_id = ?", user.ID, period.ID).Count(&attendanceDays).Error; errattday != nil {
			return errattday
		}

		prorateSalary := float64(attendanceDays) / float64(workingDays) * salary
		paySlip := models.Payslip{
			UserID:             user.ID,
			PeriodID:           period.ID,
			BaseSalary:         salary,
			WorkingDays:        workingDays,
			AttendedDays:       int(attendanceDays),
			ProratedSalary:     prorateSalary,
			OvertimeHours:      0.0,
			OvertimePay:        0.0,
			ReimbursementTotal: 0.0,
			TakeHomePay:        0.0,
			CreatedBy:          &user.ID,
			UpdatedBy:          &user.ID,
			RequestIP:          "127.0.0.1",
			GeneratedAt:        time.Now(),
			GeneratedBy:        &user.ID,
		}

		if err := db.FirstOrCreate(&paySlip).Error; err != nil {
			return err
		}
	}
	return nil
}

func GeneratePaySlipByEmployeeID(db *gorm.DB, employeeID uint, ip string, requestID string) (models.Payslip, error) {
	var user models.User
	if err := db.First(&user, employeeID).Error; err != nil {
		return models.Payslip{}, err
	}

	period, errgetper := GetPeriodData(time.Now(), db)
	if errgetper != nil {
		return models.Payslip{}, errgetper
	}

	var tempPaySlip models.Payslip
	if err := db.Preload("User").Preload("Period").Table("payslips").Where("user_id = ? AND period_id = ?", user.ID, period.ID).Scan(&tempPaySlip).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("No payslip found for user:", user.ID, "and period:", period.ID)
			return models.Payslip{}, err // No payslip found for this user and period
		}
		return models.Payslip{}, err
	}

	workingDays := utils.GetWorkingDaysInMonth(time.Now().Year(), time.Now().Month())
	var attendanceDays int64

	if errattday := db.Model(&models.Attendance{}).Where("user_id = ? AND period_id = ?", user.ID, period.ID).Count(&attendanceDays).Error; errattday != nil {
		return models.Payslip{}, errattday
	}

	prorateSalary := float64(attendanceDays) / float64(workingDays) * tempPaySlip.BaseSalary

	var overtimeHours float64
	if errot := db.Table("overtimes").Select("SUM(hours) AS overtimeHours").Where("user_id = ? AND period_id = ?", user.ID, period.ID).Group("user_id").Scan(&overtimeHours).Error; errot != nil {
		return models.Payslip{}, errot
	}

	overtimeAmount, errgetamount := GetOvertimeAmount(db, period.ID)
	if errgetamount != nil {
		return models.Payslip{}, errgetamount
	}

	overtimePay := float64(overtimeHours) * overtimeAmount

	var reimbursementList []models.Reimbursement
	if errreim := db.Where("user_id = ? AND period_id = ?", user.ID, period.ID).Find(&reimbursementList).Error; errreim != nil {
		return models.Payslip{}, errreim
	}

	var reimbursementTotal float64
	for _, reimbursement := range reimbursementList {
		reimbursementTotal += reimbursement.Amount
	}

	takeHomePay := prorateSalary + overtimePay + reimbursementTotal

	paySlip := models.Payslip{
		ID:                 tempPaySlip.ID,
		UserID:             user.ID,
		PeriodID:           period.ID,
		BaseSalary:         tempPaySlip.BaseSalary,
		WorkingDays:        workingDays,
		AttendedDays:       int(attendanceDays),
		ProratedSalary:     prorateSalary,
		OvertimeHours:      float64(overtimeHours),
		OvertimePay:        overtimePay,
		ReimbursementTotal: reimbursementTotal,
		TakeHomePay:        takeHomePay,
		CreatedBy:          &user.ID,
		UpdatedBy:          &user.ID,
		RequestIP:          ip,
		GeneratedAt:        time.Now(),
		GeneratedBy:        &user.ID,
		RequestID:          &requestID,
		User:               user,
		Period:             period,
	}

	log.Println(paySlip)

	if err := db.Table("payslips").Where("id = ? AND user_id = ? AND period_id = ?", tempPaySlip.ID, user.ID, period.ID).Updates(&paySlip).Error; err != nil {
		return models.Payslip{}, err
	}

	return paySlip, nil
}
