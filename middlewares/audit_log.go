package middlewares

import (
	"DeallsJobsTest/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

type AuditLogHandler struct {
	DB *gorm.DB
}

func NewAuditLog(db *gorm.DB) *AuditLogHandler {
	return &AuditLogHandler{
		DB: db,
	}
}

func (a *AuditLogHandler) LogAction(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uint)
	recordID := uuid.New().String()
	var actionType string

	switch ctx.Method() {
	case fiber.MethodPost:
		actionType = "create"
	case fiber.MethodPut:
		actionType = "update"
	case fiber.MethodDelete:
		actionType = "delete"
	case fiber.MethodGet:
		actionType = "read"
	default:
		return fiber.NewError(fiber.StatusBadRequest, "Unsupported HTTP method")
	}

	actionType = strings.ToUpper(actionType)
	logEntry := &models.AuditLog{
		Method:      ctx.Method(),
		Path:        ctx.Path(),
		ActionType:  actionType,
		PerformedBy: &userID,
		RequestIP:   ctx.IP(),
		RequestID:   recordID,
		CreatedAt:   time.Now(),
	}

	if err := a.DB.Create(logEntry).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to log action",
			"error":   err.Error(),
		})
	}

	ctx.Locals("trackingID", recordID)

	return ctx.Next()
}
