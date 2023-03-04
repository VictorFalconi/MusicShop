package service

import (
	"errors"
	"fmt"
	"server/app/model"
	"server/app/repository"
	"server/helpers"
	"server/middleware"
)

type UserService struct {
	repo repository.UserRepoInterface
}

type UserServiceInterface interface {
	Register(user *model.User) error
	Login(loginUser *model.LoginUser) (error, string)
	Update(currUser *model.User, newUser *model.User) error

	//OAuth
	GoogleLogin(userInfo map[string]interface{}) (error, string)
}

func NewUserService(repo repository.UserRepoInterface) *UserService {
	return &UserService{repo}
}

func (s *UserService) Register(user *model.User) error {
	// Set role for User
	roleName := "user"
	// Find role
	role, err := s.repo.FindRoleByName(roleName)
	if err != nil {
		return errors.New("dont find role name")
	}
	// Set role
	user.RoleId = role.Id
	// HashPassword
	if errHP := user.HashPassword(); errHP != nil {
		return errors.New("cant hash password")
	}
	// Create
	return s.repo.Create(user)
}

func (s *UserService) Login(loginUser *model.LoginUser) (error, string) {
	// Find
	user, err := s.repo.FindUser(loginUser.Input)
	if err != nil {
		return errors.New("name or email isn't already exist"), ""
	}
	// Compare password
	if user.ComparePassword(loginUser.Password) == false {
		return errors.New("incorrect password"), ""
	}
	// Create token
	token, errCreate := middleware.CreateToken(user.Id)
	if errCreate != nil {
		return errCreate, ""
	}
	return nil, token
}

func (s *UserService) Update(currUser *model.User, newUser *model.User) error {
	// Update struct
	currUser.UpdateStruct(newUser)
	// HashPassword
	if errHP := currUser.HashPassword(); errHP != nil {
		return errors.New("cant hash password")
	}
	//Save
	return s.repo.Update(currUser)
}

func (s *UserService) GoogleLogin(userInfo map[string]interface{}) (error, string) {
	email := userInfo["email"].(string)
	name := userInfo["name"].(string)
	verified := userInfo["verified_email"].(bool)
	// Verify
	if verified != true {
		return errors.New("unverified email"), ""
	}
	// gmail ( chua dk -> random password -> register -> login , da dk = gmail -> login bth)    |gmail trung voi gmail co san trong db -> ???|

	// Find email in db
	user, err := s.repo.FindUser(email)
	if err != nil {
		// Register ( Phone is null)
		password := helpers.RandomString(10)
		fmt.Println(password)
		newUser := model.User{Name: name, Email: email, Password: password}
		if errRegister := s.Register(&newUser); errRegister != nil {
			return errRegister, ""
		}
		// Login
		loginUser := model.LoginUser{Input: email, Password: password}
		errLogin, token := s.Login(&loginUser)
		return errLogin, token
	}
	// Login -> Create token
	token, errCreate := middleware.CreateToken(user.Id)
	if errCreate != nil {
		return errCreate, ""
	}
	return nil, token
}
