package controllers

import (
	"DeallsJobsTest/models"
	"DeallsJobsTest/services"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type PayrollController struct {
	DB  *gorm.DB
	RDB *redis.Client
}

func NewPayrollController(db *gorm.DB, rdb *redis.Client) *PayrollController {
	return &PayrollController{
		DB:  db,
		RDB: rdb,
	}
}

func (h *PayrollController) CreatePeriod(c *fiber.Ctx) error {
	userCreatedID := c.Locals("userID").(uint)
	role := c.Locals("isAdmin").(bool)

	if !role {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "You do not have permission to create a payroll period",
		})
	}

	var period struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		IsLocked  bool   `json:"is_locked"`
	}
	if err := c.BodyParser(&period); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	var attendancePeriod models.AttendancesPeriod
	fromDate, _ := time.Parse("2006-01-02 15:04:05", period.StartDate)
	toDate, _ := time.Parse("2006-01-02 15:04:05", period.EndDate)

	attendancePeriod.StartDate = fromDate
	attendancePeriod.EndDate = toDate
	attendancePeriod.IsLocked = period.IsLocked
	attendancePeriod.CreatedBy = &userCreatedID
	attendancePeriod.RequestIP = c.IP()

	if err := h.DB.Create(&attendancePeriod).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create attendance period",
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Payroll retrieved successfully",
	})
}

func (h *PayrollController) RunPayroll(c *fiber.Ctx) error {
	role := c.Locals("isAdmin").(bool)
	requestID := c.Locals("trackingID").(string)
	var payRollData struct {
		PayrollDate string `json:"payroll_date"`
	}

	if err := c.BodyParser(&payRollData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	payrollDate, err := time.Parse("2006-01-02 15:04:05", payRollData.PayrollDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid payroll date format",
			"error":   err.Error(),
		})
	}

	if !role {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "You do not have permission to create a payroll period",
		})
	}

	period, errper := services.GetPeriodData(payrollDate, h.DB)
	if errper != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve attendance period",
			"error":   errper.Error(),
		})
	}

	errpublish := services.GeneratePayroll(h.DB, period.ID, c, h.RDB, requestID)

	if errpublish != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to run payroll",
			"error":   errpublish.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Payroll running..",
	})
}

func (h *PayrollController) GetPayrollSummary(c *fiber.Ctx) error {
	role := c.Locals("isAdmin").(bool)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sort := c.Query("sort")
	all := c.Query("all") == "true"

	if !role {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "You do not have permission to create a payroll period",
		})
	}

	if all {
		datapayroll, errget := services.GetPayrollAll(h.DB, sort)
		if errget != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to retrieve payroll data",
				"error":   errget.Error(),
			})
		}

		if len(datapayroll.Payslips) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "No payroll data found",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Payroll retrieved successfully",
			"data":    datapayroll,
		})
	} else {
		datapayroll, errget := services.GetPayrollDataPagination(h.DB, page, limit, sort)
		if errget != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to retrieve payroll data",
				"error":   errget.Error(),
			})
		}

		if len(datapayroll.Payslips) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "No payroll data found",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Payroll retrieved successfully",
			"data":    datapayroll,
		})
	}
}
