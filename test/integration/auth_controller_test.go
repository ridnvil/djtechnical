package integration

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestLogin(t *testing.T) {
	SetupTest()
	defer CleanupTest()

	tests := []struct {
		name           string
		username       string
		password       string
		expectedStatus int
		expectToken    bool
	}{
		{
			name:           "Valid login",
			username:       "employee1",
			password:       "employee123",
			expectedStatus: http.StatusOK,
			expectToken:    true,
		},
		{
			name:           "Invalid username",
			username:       "nonexistent",
			password:       "password",
			expectedStatus: http.StatusUnauthorized,
			expectToken:    false,
		},
		{
			name:           "Invalid password",
			username:       "admin",
			password:       "wrongpassword",
			expectedStatus: http.StatusUnauthorized,
			expectToken:    false,
		},
		{
			name:           "Empty credentials",
			username:       "",
			password:       "",
			expectedStatus: http.StatusUnauthorized,
			expectToken:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loginBody := `{"username":"` + tt.username + `","password":"` + tt.password + `"}`

			resp, body := MakeRequest(t, "POST", "/api/auth/login", loginBody, "")

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			var response map[string]interface{}
			if err := json.Unmarshal([]byte(body), &response); err != nil {
				t.Fatalf("Failed to parse response: %v", err)
			}

			_, hasToken := response["token"]
			if tt.expectToken && !hasToken {
				t.Error("Expected token in response, but none found")
			} else if !tt.expectToken && hasToken {
				t.Error("Did not expect token in response, but one was found")
			}

			if tt.expectedStatus == http.StatusOK {
				if message, ok := response["message"]; !ok || message != "Login successful" {
					t.Errorf("Expected 'Login successful' message, got %v", message)
				}
			}
		})
	}
}
