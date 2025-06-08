package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestPayrollController(t *testing.T) {
	SetupTest()
	defer CleanupTest()

	adminToken := GetAuthToken(t, "admin", "admin123")
	employeeToken := GetAuthToken(t, "employee1", "employee123")

	t.Run("CreatePeriod", func(t *testing.T) {
		tests := []struct {
			name           string
			startDate      string
			endDate        string
			isLocked       bool
			token          string
			expectedStatus int
		}{
			{
				name:           "Valid period creation by admin",
				startDate:      time.Now().AddDate(0, 0, -15).Format("2006-01-02 15:04:05"),
				endDate:        time.Now().Format("2006-01-02 15:04:05"),
				isLocked:       false,
				token:          adminToken,
				expectedStatus: http.StatusOK,
			},
			{
				name:           "Unauthorized access by employee",
				startDate:      time.Now().AddDate(0, 0, -15).Format("2006-01-02 15:04:05"),
				endDate:        time.Now().Format("2006-01-02 15:04:05"),
				isLocked:       false,
				token:          employeeToken,
				expectedStatus: http.StatusForbidden,
			},
			{
				name:           "No authentication",
				startDate:      time.Now().AddDate(0, 0, -15).Format("2006-01-02 15:04:05"),
				endDate:        time.Now().Format("2006-01-02 15:04:05"),
				isLocked:       false,
				token:          "",
				expectedStatus: http.StatusUnauthorized,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				requestBody := fmt.Sprintf(`{
					"start_date": "%s",
					"end_date": "%s",
					"is_locked": %t
				}`, tt.startDate, tt.endDate, tt.isLocked)

				resp, body := MakeRequest(t, "POST", "/api/admin/period", requestBody, tt.token)

				if resp.StatusCode != tt.expectedStatus {
					t.Errorf("Expected status %d, got %d. Response: %s", tt.expectedStatus, resp.StatusCode, body)
				}

				if tt.expectedStatus == http.StatusOK {
					var response map[string]interface{}
					if err := json.Unmarshal([]byte(body), &response); err != nil {
						t.Fatalf("Failed to parse response: %v", err)
					}

					if message, ok := response["message"]; !ok || message != "Payroll retrieved successfully" {
						t.Errorf("Expected 'Payroll retrieved successfully' message, got %v", message)
					}
				}
			})
		}
	})

	t.Run("RunPayroll", func(t *testing.T) {
		tests := []struct {
			name           string
			payrollDate    string
			token          string
			expectedStatus int
		}{
			{
				name:           "Valid payroll run by admin",
				payrollDate:    time.Now().Format("2006-01-02 15:04:05"),
				token:          adminToken,
				expectedStatus: http.StatusOK,
			},
			{
				name:           "Unauthorized access by employee",
				payrollDate:    time.Now().Format("2006-01-02 15:04:05"),
				token:          employeeToken,
				expectedStatus: http.StatusForbidden,
			},
			{
				name:           "No authentication",
				payrollDate:    time.Now().Format("2006-01-02 15:04:05"),
				token:          "",
				expectedStatus: http.StatusUnauthorized,
			},
			{
				name:           "Invalid date format",
				payrollDate:    "invalid-date",
				token:          adminToken,
				expectedStatus: http.StatusBadRequest,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				requestBody := fmt.Sprintf(`{
					"payroll_date": "%s"
				}`, tt.payrollDate)

				resp, body := MakeRequest(t, "POST", "/api/admin/payroll", requestBody, tt.token)

				if resp.StatusCode != tt.expectedStatus {
					t.Errorf("Expected status %d, got %d. Response: %s", tt.expectedStatus, resp.StatusCode, body)
				}

				if tt.expectedStatus == http.StatusOK {
					var response map[string]interface{}
					if err := json.Unmarshal([]byte(body), &response); err != nil {
						t.Fatalf("Failed to parse response: %v", err)
					}

					if message, ok := response["message"]; !ok || message != "Payroll running.." {
						t.Errorf("Expected 'Payroll running..' message, got %v", message)
					}
				}
			})
		}
	})

	t.Run("GetPayrollSummary", func(t *testing.T) {
		tests := []struct {
			name           string
			queryParams    string
			token          string
			expectedStatus int
		}{
			{
				name:           "Get payroll summary by admin",
				queryParams:    "?page=1&limit=10",
				token:          adminToken,
				expectedStatus: http.StatusOK,
			},
			{
				name:           "Get all payroll data by admin",
				queryParams:    "?all=true",
				token:          adminToken,
				expectedStatus: http.StatusOK,
			},
			{
				name:           "Unauthorized access by employee",
				queryParams:    "?page=1&limit=10",
				token:          employeeToken,
				expectedStatus: http.StatusForbidden,
			},
			{
				name:           "No authentication",
				queryParams:    "?page=1&limit=10",
				token:          "",
				expectedStatus: http.StatusUnauthorized,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				resp, body := MakeRequest(t, "GET", "/api/admin/summary"+tt.queryParams, "", tt.token)

				if resp.StatusCode != tt.expectedStatus {
					t.Errorf("Expected status %d, got %d. Response: %s", tt.expectedStatus, resp.StatusCode, body)
				}

				if tt.expectedStatus == http.StatusOK {
					var response map[string]interface{}
					if err := json.Unmarshal([]byte(body), &response); err != nil {
						t.Fatalf("Failed to parse response: %v", err)
					}

					if message, ok := response["message"]; !ok || message != "Payroll retrieved successfully" {
						t.Errorf("Expected 'Payroll retrieved successfully' message, got %v", message)
					}

					if _, ok := response["data"]; !ok {
						t.Error("Expected 'data' in response, but none found")
					}
				}
			})
		}
	})
}
