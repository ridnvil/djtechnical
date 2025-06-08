package controllers

import (
	"DeallsJobsTest/models"
	"DeallsJobsTest/services"
	"DeallsJobsTest/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"time"
)

type AttendanceController struct {
	DB *gorm.DB
}

func NewAttendanceController(db *gorm.DB) *AttendanceController {
	return &AttendanceController{
		DB: db,
	}
}

func (h *AttendanceController) SubmitAttendance(c *fiber.Ctx) error {
	employeeID := c.Locals("userID").(uint)
	var attendanceData struct {
		Date string `json:"date"`
	}
	if err := c.BodyParser(&attendanceData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request data",
			"error":   err.Error(),
		})
	}

	if attendanceData.Date == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Date is required",
			"error":   "Please provide a valid date",
		})
	}

	checkInDate, errFormatDate := time.Parse("2006-01-02 15:04:05", attendanceData.Date)

	if errFormatDate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid date format",
			"error":   "Date must be in 'YYYY-MM-DD HH:MM:SS' format",
		})
	}

	if checkInDate.Before(time.Now()) && checkInDate.After(time.Now()) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid date",
			"error":   "Check-in date must be today",
		})
	}

	period, errper := services.GetPeriodData(checkInDate, h.DB)
	if errper != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve attendance period",
			"error":   errper.Error(),
		})
	}

	if period.IsLocked {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Attendance period is locked",
			"error":   "Cannot submit attendance for locked period",
		})
	}

	weekDay := utils.IsWeekend(checkInDate)
	if weekDay {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot submit attendance on weekends",
			"error":   "Weekend attendance submission not allowed",
		})
	}

	var checkedIn struct {
		UsedID uint      `json:"used_id"`
		Date   time.Time `json:"date"`
	}

	if err := h.DB.Table("attendances").Where("user_id = ? AND date = ?", employeeID, checkInDate).Scan(&checkedIn).Error; err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Error check attendance",
			"error":   err,
		})
	}

	var attendance models.Attendance
	attendance.UserID = employeeID
	attendance.Date = checkInDate
	attendance.PeriodID = period.ID
	attendance.CreatedBy = &employeeID
	attendance.UpdatedBy = &employeeID
	attendance.RequestIP = c.IP()
	attendance.RequestID = c.Locals("trackingID").(string)

	if checkedIn.UsedID != 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Attendance already submitted for this date",
			"error":   "Attendance already exists",
		})
	}

	if err := h.DB.Create(&attendance).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to submit attendance",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Attendance submitted successfully",
	})
}
