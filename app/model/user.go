package model

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
	PhoneNumber string `json:"phonenumber"  form:"phonenumber"  gorm:"unique"          validate:"required,len=10"`
	Password    string `json:"password"     form:"password"     gorm:"not null"        validate:"required,min=4,max=32"`
	Address     string `json:"address"      form:"address"      gorm:""                validate:""`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	RoleId uint
}

type LoginUser struct {
	Input    string `json:"input"     form:"input"     validate:"required,min=4,max=32"`
	Password string `json:"password"  form:"password"  validate:"required,min=4,max=32"`
}

type ReadUser struct {
	Name        string
	Email       string
	PhoneNumber string
	Address     string
}

func (u *User) ReadUser() interface{} {
	return ReadUser{Name: u.Name, Email: u.Email, PhoneNumber: u.PhoneNumber, Address: u.Address}
}

// HashPassword :
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// Compare between password and HashPassword
func (u *User) ComparePassword(password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return false
	}
	return true
}

func (u *User) UpdateStruct(newUser *User) {
	//currUser.Name = newUser.Name
	u.Email = newUser.Email
	u.PhoneNumber = newUser.PhoneNumber
	u.Password = newUser.Password
	u.Address = newUser.Address
}
