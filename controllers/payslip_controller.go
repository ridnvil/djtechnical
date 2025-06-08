package controllers

import (
	"DeallsJobsTest/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PayslipController struct {
	DB *gorm.DB
}

func NewPayslipController(db *gorm.DB) *PayslipController {
	return &PayslipController{
		DB: db,
	}
}

func (h *PayslipController) GeneratePayslip(c *fiber.Ctx) error {
	id := c.Locals("userID").(uint)
	requestID := c.Locals("trackingID").(string)
	dataPaySlip, err := services.GeneratePaySlipByEmployeeID(h.DB, id, c.IP(), requestID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to generate payslip",
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Payslip generated successfully",
		"data":    dataPaySlip,
	})
}
