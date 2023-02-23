package repository

import (
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"server/app/model"
)

type UserRepo struct {
	db *gorm.DB
}

type UserRepoInterface interface {
	FindRoleByName(roleName string) (*model.Role, error)
	Create(user *model.User) error
	FindUserByName(name string) (*model.User, error)
	Update(user *model.User) error

	// Middleware

	GetUserFromToken(token *jwt.Token) (interface{}, error)
	GetRoleOfUser(user *model.User) (*model.Role, error)
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db}
}

func (ur *UserRepo) FindRoleByName(roleName string) (*model.Role, error) {
	var role *model.Role
	if err := ur.db.Where("name = ?", roleName).First(&role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func (ur *UserRepo) Create(user *model.User) error {
	return ur.db.Create(user).Error
}

func (ur *UserRepo) FindUserByName(name string) (*model.User, error) {
	var user *model.User
	if err := ur.db.Where("name = ?", name).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRepo) Update(user *model.User) error {
	return ur.db.Save(user).Error
}

// Middleware

func (ur *UserRepo) GetUserFromToken(token *jwt.Token) (interface{}, error) {
	claims := token.Claims.(jwt.MapClaims)
	// Query User from id
	var user model.User
	if err := ur.db.Where("id = ?", claims["id"]).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRepo) GetRoleOfUser(user *model.User) (*model.Role, error) {
	var role *model.Role
	if err := ur.db.Where("id = ?", user.RoleId).First(&role).Error; err != nil {
		return nil, err
	}
	return role, nil
}
