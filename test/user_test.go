package test

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"server/app/models"
	"server/app/routes"
	"testing"
)

func TestRegister(t *testing.T) {
	r := gin.Default()
	routes.UserRouter(r)

	// create a new test server using the test router
	ts := httptest.NewServer(r)
	defer ts.Close()

	user := models.User{
		Name:        "Liem",
		Email:       "liem@gmail.com",
		PhoneNumber: "0987123123",
		Password:    "12345",
	}
	// Define a mock request with the user data as JSON
	body, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", ts.URL+"/user/register", bytes.NewBuffer(body))

	// Set the request Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Define a mock response recorder
	w := httptest.NewRecorder()

	// Call the API endpoint using the Gin router and the mock request and response
	r.ServeHTTP(w, req)

	// Check the response status code and body
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "User registered successfully", w.Body.String())
}
