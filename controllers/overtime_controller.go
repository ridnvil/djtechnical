package controllers

import (
	"DeallsJobsTest/models"
	"DeallsJobsTest/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"time"
)

type OvertimeController struct {
	DB *gorm.DB
}

func NewOvertimeController(db *gorm.DB) *OvertimeController {
	return &OvertimeController{
		DB: db,
	}
}

func (h *OvertimeController) SubmitOvertime(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var overtimeData struct {
		Date  string  `json:"date"`
		Hours float64 `json:"hours"`
	}
	if err := c.BodyParser(&overtimeData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request data",
			"error":   err.Error(),
		})
	}

	inputDate, err := time.Parse("2006-01-02 15:04:05", overtimeData.Date)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid date format",
			"error":   err.Error(),
		})
	}

	thisDay := time.Now()

	if !(overtimeData.Hours >= 1) && !(overtimeData.Hours <= 3) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid overtime hours",
			"error":   "Overtime hours must be between 1 and 3",
		})
	}

	thisHours := thisDay.Hour()

	if !(thisHours >= 8 && thisHours <= 23) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "You can propose overtime after working hours",
			"error":   "Overtime can only be submitted between 17.00 and 23.00",
		})
	}

	period, errgetper := services.GetPeriodData(inputDate, h.DB)
	if errgetper != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve attendance period",
			"error":   errgetper.Error(),
		})
	}

	if period.IsLocked {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Attendance period is locked",
			"error":   "Cannot submit attendance for locked period",
		})
	}

	var overtime models.Overtime
	overtime.Date = inputDate
	overtime.Hours = overtimeData.Hours
	overtime.UserID = userID
	overtime.RequestIP = c.IP()
	overtime.CreatedBy = &userID
	overtime.PeriodID = period.ID

	if err := h.DB.Create(&overtime).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to submit overtime",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Overtime submitted successfully",
	})
}
