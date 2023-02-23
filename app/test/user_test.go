package test

import (
	"bytes"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"server/app/controller"
	"server/app/model"
	"server/app/repository"
	"server/app/service"
	"server/app/test/mocks"
	"server/helpers"
	"testing"
)

func Test_Repo_Create(t *testing.T) {
	mockRepo := new(mocks.MockDBUSer)

	user := model.User{
		Name:        "Liem",
		Email:       "liem@gmail.com",
		PhoneNumber: "0987123123",
		Password:    "12345",
	}
	mockRepo.On("Create", &user).Return(nil)
	err := mockRepo.Create(&user)

	//Test
	assert.Equal(t, nil, err)
	var expected_RoleID uint = 1
	assert.Equal(t, expected_RoleID, user.RoleId)
}

func Test_Repo_GetRoleByName(t *testing.T) {
	mockUser := new(mock.Mock)

	role := model.Role{Name: "user"}

	mockUser.On("GetRoleByName", "user").Return(&role, nil)

	//result := mockUser.MethodCalled("GetRoleByName", "user").Get(0)
	err := mockUser.MethodCalled("GetRoleByName", "user").Error(1)

	if err != nil {
		t.Errorf("GetRoleByName function returned an unexpected error: %v", err)
	}

}

func Test_API_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock database & mocks GORM
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Create repo-service-controller
	userRepo := repository.NewUserRepo(gormDB)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	// Create JSON payload
	user := model.User{
		Name:        "Liem",
		Email:       "liem@gmail.com",
		PhoneNumber: "0987123123",
		Password:    "12345",
		Address:     "HD,KG",
	}
	payload, _ := json.Marshal(user)

	//Set up the mocked database response     - Select ( Find ID from name) -> Create (User -> json)
	role := &model.Role{Name: "user"}

	// Select
	mock.ExpectQuery(`SELECT \* FROM "roles" WHERE name = \$1 AND "roles"."deleted_at" IS NULL ORDER BY "roles"."id" LIMIT 1`).
		WithArgs(role.Name).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, role.Name))
	// Insert
	mock.ExpectBegin()
	var roleID uint = 1
	mock.ExpectExec("INSERT INTO users").WithArgs(user.Name, user.Email, user.PhoneNumber, user.Password, user.Address, nil, nil, nil, roleID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Create a new HTTP request with the test case input
	r := gin.Default()
	r.POST("/user/register", userController.RegisterHandler())

	req, err := http.NewRequest("POST", "/user/register", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}
	// header
	req.Header.Set("Content-Type", "application/json")
	// Create a new HTTP recorder to capture the response
	w := httptest.NewRecorder()
	// Call the RegisterHandler function with the mocked dependencies
	r.ServeHTTP(w, req)

	// Check status
	assert.Equal(t, 201, w.Code)

	// Check response
	expectedResponse := &helpers.ResponseData{
		Status:  201,
		Message: helpers.StatusCodeFromInt(201),
		Error:   nil,
		Data:    nil,
	}
	assert.Equal(t, expectedResponse, w.Body.String())
}
