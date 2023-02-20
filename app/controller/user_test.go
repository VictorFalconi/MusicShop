package controller

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"server/app/models"
	"testing"
)

func TestRegister(t *testing.T) {
	r := gin.Default()

	user := models.User{
		Name:        "Liem",
		Email:       "liem@gmail.com",
		PhoneNumber: "0987123123",
		Password:    "12345",
	}
	// Define a mock request with the user data as JSON
	body, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/user/register", bytes.NewReader(body))

	// Set the request Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Define a mock response recorder
	w := httptest.NewRecorder()

	// Call the API endpoint using the Gin router and the mock request and response
	r.ServeHTTP(w, req)

	// Define a subtest for successful user registration
	t.Run("Success", func(t *testing.T) {
		// Check the response status code
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d but got %d", http.StatusOK, w.Code)
		}

		// Check the response body
		respBody, err := ioutil.ReadAll(w.Body)
		if err != nil {
			t.Errorf("Error reading response body: %v", err)
		}

		expectedResp := "User registered successfully"
		if string(respBody) != expectedResp {
			t.Errorf("Expected response body '%s' but got '%s'", expectedResp, string(respBody))
		}
	})

	// Define a subtest for invalid request data
	t.Run("InvalidRequest", func(t *testing.T) {
		// Define a mock request with invalid data
		invalidReq, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer([]byte{}))
		invalidReq.Header.Set("Content-Type", "application/json")
		invalidRes := httptest.NewRecorder()

		// Call the API endpoint with invalid request data
		r.ServeHTTP(invalidRes, invalidReq)

		// Check the response status code
		if invalidRes.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d but got %d", http.StatusBadRequest, invalidRes.Code)
		}

		// Check the response body
		respBody, err := ioutil.ReadAll(invalidRes.Body)
		if err != nil {
			t.Errorf("Error reading response body: %v", err)
		}

		expectedResp := "Invalid request data"
		if string(respBody) != expectedResp {
			t.Errorf("Expected response body '%s' but got '%s'", expectedResp, string(respBody))
		}
	})

}
