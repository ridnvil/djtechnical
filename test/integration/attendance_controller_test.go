package integration

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestSubmitAttendance(t *testing.T) {
	SetupTest()
	defer CleanupTest()

	token := GetAuthToken(t, "employee1", "employee123")

	tests := []struct {
		name           string
		date           string
		expectedStatus int
		withToken      bool
	}{
		{
			name:           "Valid attendance submission",
			date:           time.Now().AddDate(0, 0, +2).Format("2006-01-02 15:04:05"),
			expectedStatus: http.StatusOK,
			withToken:      true,
		},
		{
			name:           "No authentication",
			date:           time.Now().Format("2006-01-02 15:04:05"),
			expectedStatus: http.StatusUnauthorized,
			withToken:      false,
		},
		{
			name:           "Invalid date format",
			date:           "invalid-date",
			expectedStatus: http.StatusBadRequest,
			withToken:      true,
		},
		{
			name:           "Weekend date",
			date:           getNextWeekendDate(),
			expectedStatus: http.StatusBadRequest,
			withToken:      true,
		},
		{
			name:           "Empty date",
			date:           "",
			expectedStatus: http.StatusBadRequest,
			withToken:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBody := `{"date":"` + tt.date + `"}`

			var authToken string
			if tt.withToken {
				authToken = token
			}

			resp, body := MakeRequest(t, "POST", "/api/employee/attendance", requestBody, authToken)

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Response: %s", tt.expectedStatus, resp.StatusCode, body)
			}

			var response map[string]interface{}
			if err := json.Unmarshal([]byte(body), &response); err != nil {
				t.Fatalf("Failed to parse response: %v", err)
			}

			if tt.expectedStatus == http.StatusOK {
				if message, ok := response["message"]; !ok || message != "Attendance submitted successfully" {
					t.Errorf("Expected 'Attendance submitted successfully' message, got %v", message)
				}
			}
		})
	}
}

func getNextWeekendDate() string {
	now := time.Now()
	daysUntilSaturday := (6 - int(now.Weekday())) % 7
	if daysUntilSaturday == 0 {
		daysUntilSaturday = 7
	}
	nextWeekend := now.AddDate(0, 0, daysUntilSaturday)
	return nextWeekend.Format("2006-01-02 15:04:05")
}
