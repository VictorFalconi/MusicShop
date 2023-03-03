package service

import (
	"errors"
	"server/app/model"
	"server/app/repository"
	"server/middleware"
)

type UserService struct {
	repo repository.UserRepoInterface
}

type UserServiceInterface interface {
	Register(user *model.User) error
	Login(loginUser *model.LoginUser) (error, string)
	Update(currUser *model.User, newUser *model.User) error
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
