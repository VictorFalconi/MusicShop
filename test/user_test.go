package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"server/app/controller"
	"server/app/models"
	"testing"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestRegister(t *testing.T) {
	r := SetUpRouter()
	r.POST("/user/register", controller.Register)
	user := models.User{
		Name:        "Liem",
		Email:       "liem@gmail.com",
		PhoneNumber: "0987123123",
		Password:    "12345",
	}
	user.HashPassword()
	//github.com/fatih/structs
	jsonUser, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/user/register", bytes.NewBuffer(jsonUser))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	//https://github.com/stretchr/testify
	assert.Equal(t, http.StatusCreated, w.Code)
}
