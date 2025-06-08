package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSubmitReimbursement(t *testing.T) {
	SetupTest()
	defer CleanupTest()

	token := GetAuthToken(t, "employee1", "employee123")

	tests := []struct {
		name           string
		date           string
		amount         float64
		description    string
		expectedStatus int
		withToken      bool
	}{
		{
			name:           "Valid reimbursement submission",
			date:           time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
			amount:         100000,
			description:    "Office supplies",
			expectedStatus: http.StatusOK,
			withToken:      true,
		},
		{
			name:           "No authentication",
			date:           time.Now().Format("2006-01-02 15:04:05"),
			amount:         100000,
			description:    "Office supplies",
			expectedStatus: http.StatusUnauthorized,
			withToken:      false,
		},
		{
			name:           "Future date",
			date:           time.Now().AddDate(0, 0, 10).Format("2006-01-02"),
			amount:         100000,
			description:    "Office supplies",
			expectedStatus: http.StatusBadRequest,
			withToken:      true,
		},
		{
			name:           "Empty description",
			date:           time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
			amount:         100000,
			description:    "",
			expectedStatus: http.StatusOK,
			withToken:      true,
		},
	}

	var reimbursementID uint

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBody := fmt.Sprintf(`{
				"date": "%s",
				"amount": %.2f,
				"description": "%s"
			}`, tt.date, tt.amount, tt.description)

			var authToken string
			if tt.withToken {
				authToken = token
			}

			resp, body := MakeRequest(t, "POST", "/api/employee/reimbursement", requestBody, authToken)

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Response: %s", tt.expectedStatus, resp.StatusCode, body)
			}

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(body), &response); err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}

				if message, ok := response["message"]; !ok || message != "Reimbursement submitted successfully" {
					t.Errorf("Expected 'Reimbursement submitted successfully' message, got %v", message)
				}

				if id, ok := response["ID"]; ok {
					reimbursementID = uint(id.(float64))
				}
			}
		})
	}

	if reimbursementID > 0 {
		t.Run("Upload attachments", func(t *testing.T) {
			tempFile, err := os.CreateTemp("", "test-*.txt")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tempFile.Name())

			if _, err := tempFile.WriteString("Test file content"); err != nil {
				t.Fatalf("Failed to write to temp file: %v", err)
			}
			tempFile.Close()

			var b bytes.Buffer
			w := multipart.NewWriter(&b)

			f, err := w.CreateFormFile("files", filepath.Base(tempFile.Name()))
			if err != nil {
				t.Fatalf("Failed to create form file: %v", err)
			}
			fileContent, err := os.ReadFile(tempFile.Name())
			if err != nil {
				t.Fatalf("Failed to read temp file: %v", err)
			}
			if _, err := f.Write(fileContent); err != nil {
				t.Fatalf("Failed to write to form file: %v", err)
			}

			w.Close()

			req, err := http.NewRequest("POST", fmt.Sprintf("/api/employee/reimbursement/uploads/%d", reimbursementID), &b)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", w.FormDataContentType())
			req.Header.Set("Authorization", "Bearer "+token)

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}

			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status %d, got %d. Response: %s", http.StatusOK, resp.StatusCode, string(respBody))
			}

			var response map[string]interface{}
			if err := json.Unmarshal(respBody, &response); err != nil {
				t.Fatalf("Failed to parse response: %v", err)
			}

			if message, ok := response["message"]; !ok || message != "Sucessfully uploaded files Reimbursement" {
				t.Errorf("Expected 'Sucessfully uploaded files Reimbursement' message, got %v", message)
			}

			if _, ok := response["files"]; !ok {
				t.Error("Expected 'files' array in response, but none found")
			}
		})
	}
}
