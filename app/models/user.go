package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type Role struct {
	Id        uint   `json:"ID"           form:"ID"           gorm:"primary_key" `
	Name      string `json:"name"         form:"name"         gorm:"unique;not null" validate:"required,min=4,max=32"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	User User `gorm:"foreignKey:RoleId;references:Id"` //role 1-n user
}

type User struct {
	Id          uint   `json:"ID"           form:"ID"           gorm:"primary_key" `
	Name        string `json:"name"         form:"name"         gorm:"unique;not null" validate:"required,min=4,max=32"`
	Email       string `json:"email"        form:"email"        gorm:"unique"          validate:"required,email,min=4,max=32"`
	PhoneNumber string `json:"phonenumber" form:"phonenumber" gorm:"unique"          validate:"required,len=10"`
	Password    string `json:"password"     form:"password"     gorm:"not null"        validate:"required,min=4,max=32"`
	Address     string `json:"address"      form:"address"      gorm:""                validate:""`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	RoleId uint
}

func (currUser *User) UpdateStruct(newUser *User) {
	//currUser.Name = newUser.Name
	currUser.Email = newUser.Email
	currUser.PhoneNumber = newUser.PhoneNumber
	currUser.Password = newUser.Password
	currUser.Address = newUser.Address
}

type LoginUser struct {
	Name     string `json:"name"     form:"name"     validate:"required,min=4,max=32"`
	Password string `json:"password" form:"password" validate:"required,min=4,max=32"`
}

// Get Name_Role of User
func (currUser *User) GetNameRoleUser(db *gorm.DB) (error, string) {
	var role Role
	if err := db.Where("id = ?", currUser.RoleId).First(&role).Error; err != nil {
		return err, ""
	}
	return nil, role.Name
}

// Set Role of User form Name_Role:
func (currUser *User) SetUserRole(db *gorm.DB, roleName string) error {
	var role Role
	if err := db.Where("name = ?", roleName).First(&role).Error; err != nil {
		return err
	}
	currUser.RoleId = role.Id
	return nil
}

// HashPassword :
func (currUser *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(currUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	currUser.Password = string(hashedPassword)
	return nil
}

// ComparePassword : Compare between password and HashPassword
func (currUser *User) ComparePassword(password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(currUser.Password), []byte(password)); err != nil {
		return false
	}
	return true
}
