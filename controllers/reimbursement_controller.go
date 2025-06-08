package controllers

import (
	"DeallsJobsTest/models"
	"DeallsJobsTest/services"
	"DeallsJobsTest/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"log"
	"path/filepath"
	"time"
)

type ReimbursementController struct {
	DB *gorm.DB
}

func NewReimbursementController(db *gorm.DB) *ReimbursementController {
	return &ReimbursementController{
		DB: db,
	}
}

func (h *ReimbursementController) SubmitReimbursement(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	requestID := c.Locals("trackingID").(string)

	log.Println(requestID)
	var reimbusementData struct {
		Date        string  `json:"date"`
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
	}

	if err := c.BodyParser(&reimbusementData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	timeNow := time.Now()
	// create Time now only hours, minutes, seconds
	reimbusementData.Date = reimbusementData.Date + " " + timeNow.Format("15:04:05")

	inputReimDate, err := time.Parse("2006-01-02 15:04:05", reimbusementData.Date)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid date format",
			"error":   err.Error(),
		})
	}

	if inputReimDate.After(timeNow) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Reimbursement date cannot be in the future",
			"error":   "Please provide a valid date",
		})
	}

	period, errgetper := services.GetPeriodData(time.Now(), h.DB)
	if errgetper != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve period data",
			"error":   errgetper.Error(),
		})
	}

	var reimbursement models.Reimbursement
	reimbursement.Amount = reimbusementData.Amount
	reimbursement.Description = reimbusementData.Description
	reimbursement.Date = inputReimDate
	reimbursement.UserID = userID
	reimbursement.PeriodID = period.ID
	reimbursement.CreatedBy = &userID
	reimbursement.RequestIP = c.IP()
	reimbursement.RequestID = requestID

	if err := h.DB.Create(&reimbursement).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to submit reimbursement",
			"error":   err.Error(),
		})
	}

	if errget := h.DB.Preload("User").Preload("Period").First(&reimbursement, reimbursement.ID).Error; errget != nil {
		if errors.Is(errget, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Reimbursement not found",
				"error":   "No reimbursement record found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving reimbursement",
			"error":   errget.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Reimbursement submitted successfully",
		"ID":      reimbursement.ID,
	})
}

func (h *ReimbursementController) UploadAttcahments(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	reimbursementID := c.Params("id")
	userIDString := fmt.Sprint(userID)

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed parsing form multipart",
		})
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No file diupload",
		})
	}

	var reimbursement models.Reimbursement
	if errget := h.DB.First(&reimbursement, reimbursementID).Error; errget != nil {
		if errors.Is(errget, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Reimbursement not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed get data reimbursement: " + errget.Error(),
		})
	}

	period, errgetper := services.GetPeriodData(reimbursement.Date, h.DB)
	if errgetper != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve period data",
			"error":   errgetper.Error(),
		})
	}

	if period.IsLocked {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Attendance period is locked",
			"error":   "Cannot submit attendance for locked period",
		})
	}

	uploadDir := "uploads/reimbursements/" + fmt.Sprint(period.ID) + "/" + userIDString
	if errcreatedir := utils.AutoCreateFolder(uploadDir); errcreatedir != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": errcreatedir.Error(),
		})
	}

	if errdelfile := utils.DeleteAllFilesInFolder(uploadDir); errdelfile != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": errdelfile.Error(),
		})
	}

	var uploadedFiles []string

	for _, file := range files {
		ext := filepath.Ext(file.Filename)
		timestamp := time.Now().Format("20060102150405")
		newFileName := fmt.Sprintf("reimbursements_%d%s%s", reimbursement.ID, timestamp, ext)
		savePath := filepath.Join(uploadDir, newFileName)

		if err := c.SaveFile(file, savePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failde save file: " + file.Filename,
			})
		}

		uploadedFiles = append(uploadedFiles, newFileName)
	}

	uploadedFilesJSON, errprsing := json.Marshal(uploadedFiles)
	if errprsing != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed convertion file to JSON: " + errprsing.Error(),
		})
	}

	reimbursement.PathFile = string(uploadedFilesJSON)
	reimbursement.RequestIP = c.IP()
	reimbursement.UpdatedBy = &userID

	if errupdate := h.DB.Model(&models.Reimbursement{}).Where("id = ?", reimbursement.ID).Updates(&reimbursement).Error; errupdate != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed update data reimbursement: " + errupdate.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Sucessfully uploaded files Reimbursement",
		"files":   uploadedFiles,
	})
}
