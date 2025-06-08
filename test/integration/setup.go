package integration

import (
	"DeallsJobsTest/models"
	"DeallsJobsTest/routes"
	"DeallsJobsTest/test/initconfig"
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	app         *fiber.App
	db          *gorm.DB
	redisClient *redis.Client
)

// SetupTest initializes the test environment
func SetupTest() {
	ctx := context.Background()
	db = initconfig.InitDatabase()
	redisClient = initconfig.InitRedis(ctx)
	app = fiber.New(fiber.Config{AppName: "DeallsJobsTest"})

	// Migrate database tables
	db.AutoMigrate(&models.User{},
		&models.OvertimePaid{},
		&models.Attendance{},
		&models.AttendancesPeriod{},
		&models.Attendance{},
		&models.Overtime{},
		&models.Reimbursement{},
		&models.Payslip{},
		&models.AuditLog{})

	// Setup routes
	routes.SetupRoutes(app, db, redisClient)
}

// MakeRequest makes an HTTP request to the test server
func MakeRequest(t *testing.T, method, path string, body string, token string) (*http.Response, string) {
	// Create a new http request
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Add authorization header if token is provided
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	// Perform the request
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	return resp, string(respBody)
}

// GetAuthToken gets an authentication token for testing
func GetAuthToken(t *testing.T, username, password string) string {
	loginBody := `{"username":"` + username + `","password":"` + password + `"}`
	_, respBody := MakeRequest(t, "POST", "/api/auth/login", loginBody, "")

	// Extract token from response
	// This is a simple extraction and might need to be adjusted based on the actual response format
	if !strings.Contains(respBody, "token") {
		t.Fatalf("Failed to get auth token: %s", respBody)
	}

	// Extract token from JSON response
	// This is a simple extraction and assumes the token is in a specific format
	tokenStart := strings.Index(respBody, `"token":"`) + 9
	tokenEnd := strings.Index(respBody[tokenStart:], `"`) + tokenStart
	return respBody[tokenStart:tokenEnd]
}

// CleanupTest cleans up the test environment
func CleanupTest() {
	// Clean up resources if needed
}
