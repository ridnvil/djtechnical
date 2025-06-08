package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestSubmitOvertime(t *testing.T) {
	SetupTest()
	defer CleanupTest()

	token := GetAuthToken(t, "employee1", "employee123")

	tests := []struct {
		name           string
		date           string
		hours          float64
		expectedStatus int
		withToken      bool
	}{
		{
			name:           "Valid overtime submission",
			date:           time.Now().Format("2006-01-02 15:04:05"),
			hours:          2.0,
			expectedStatus: http.StatusCreated,
			withToken:      true,
		},
		{
			name:           "No authentication",
			date:           time.Now().Format("2006-01-02 15:04:05"),
			hours:          2.0,
			expectedStatus: http.StatusUnauthorized,
			withToken:      false,
		},
		{
			name:           "Invalid date format",
			date:           "invalid-date",
			hours:          2.0,
			expectedStatus: http.StatusBadRequest,
			withToken:      true,
		},
		{
			name:           "Hours too low",
			date:           time.Now().Format("2006-01-02 15:04:05"),
			hours:          0.5,
			expectedStatus: http.StatusBadRequest,
			withToken:      true,
		},
		{
			name:           "Hours too high",
			date:           time.Now().Format("2006-01-02 15:04:05"),
			hours:          4.0,
			expectedStatus: http.StatusBadRequest,
			withToken:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBody := `{"date":"` + tt.date + `","hours":` + formatFloat(tt.hours) + `}`

			var authToken string
			if tt.withToken {
				authToken = token
			}

			resp, body := MakeRequest(t, "POST", "/api/employee/overtime", requestBody, authToken)

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Response: %s", tt.expectedStatus, resp.StatusCode, body)
			}

			if tt.expectedStatus == http.StatusCreated {
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(body), &response); err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}

				if message, ok := response["message"]; !ok || message != "Overtime submitted successfully" {
					t.Errorf("Expected 'Overtime submitted successfully' message, got %v", message)
				}
			}
		})
	}
}

func formatFloat(f float64) string {
	return fmt.Sprintf("%.1f", f)
}
