package test

import (
	"bytes"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"server/app/controller"
	"server/app/models"
	"testing"
)

//go test -v -cover C:\Users\Liem\Desktop\MusicShop\server\test

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a new instance of the mock database
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	// Create a new instance of the GORM database using the mock DB
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: mockDB,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create GORM database: %v", err)
	}

	r := gin.Default()
	r.POST("/user/register", controller.RegisterHandler(gormDB))

	//role := models.Role{
	//	Name: "user",
	//}

	user := models.User{
		Name:        "liem",
		Email:       "liem@gmail.com",
		PhoneNumber: "0123123321",
		Password:    "12345",
		Address:     "HD,KG",
	}

	//
	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT \* FROM "roles" WHERE name = \$1 AND "roles"."deleted_at" IS NULL ORDER BY "roles"."id" LIMIT 1`).WithArgs("user").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "user"))
	mock.ExpectExec(`INSERT INTO "users" \("name","email","phone_number","password","address","role_id"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6\)`).
		WithArgs("liem", "liem@gmail.com", "0123123321", "12345", "HD,KG", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Define a mock request with the user data as JSON
	body, _ := json.Marshal(user)
	req, err := http.NewRequest("POST", "/user/register", bytes.NewBuffer(body))
	assert.NoError(t, err)

	// Set the request Content-Type header
	req.Header.Set("Access-Control-Allow-Credentials", "true")
	req.Header.Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	req.Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	req.Header.Set("Access-Control-Allow-Origin", "http://localhost:3000")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	// Define a mock response recorder
	w := httptest.NewRecorder()

	// Call the API endpoint using the Gin router and the mock request and response
	r.ServeHTTP(w, req)

	//Check the response status code and body
	assert.Equal(t, http.StatusCreated, w.Code)

	expectedBody := `{
    "Status": 201,
    "Message": "The new resource has been created",
    "Error": null,
    "Data": null
}`
	assert.Equal(t, expectedBody, w.Body.String())

	// Check the mock database expectations
	//err = mock.ExpectationsWereMet()
	//assert.NoError(t, err)
}
