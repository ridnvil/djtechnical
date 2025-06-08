package integration

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestGeneratePayslip(t *testing.T) {
	SetupTest()
	defer CleanupTest()

	token := GetAuthToken(t, "employee1", "employee123")

	tests := []struct {
		name           string
		id             uint
		expectedStatus int
		withToken      bool
	}{
		{
			name:           "Valid payslip generation",
			id:             1,
			expectedStatus: http.StatusOK,
			withToken:      true,
		},
		{
			name:           "No authentication",
			id:             0,
			expectedStatus: http.StatusUnauthorized,
			withToken:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var authToken string
			if tt.withToken {
				authToken = token
			}

			resp, body := MakeRequest(t, "GET", "/api/employee/payslip", "", authToken)

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Response: %s", tt.expectedStatus, resp.StatusCode, body)
			}

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(body), &response); err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}

				if message, ok := response["message"]; !ok || message != "Payslip generated successfully" {
					t.Errorf("Expected 'Payslip generated successfully' message, got %v", message)
				}

				if data, ok := response["data"]; !ok {
					t.Error("Expected 'data' in response, but none found")
				} else {
					payslipData, ok := data.(map[string]interface{})
					if !ok {
						t.Error("Expected 'data' to be a JSON object")
					} else {
						expectedFields := []string{"ID", "UserID", "PeriodID", "BaseSalary"}
						for _, field := range expectedFields {
							if _, ok := payslipData[field]; !ok {
								t.Errorf("Expected field '%s' in payslip data, but not found", field)
							}
						}
					}
				}
			}
		})
	}
}
