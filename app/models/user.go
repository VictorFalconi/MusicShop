package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"server/helpers"
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

	Orders []Order `gorm:"foreignKey:UserId;references:Id"` //user 1-n order

}

type LoginUser struct {
	Name     string `json:"name"     form:"name"     validate:"required,min=4,max=32"`
	Password string `json:"password" form:"password" validate:"required,min=4,max=32"`
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

func (currUser *User) UpdateStruct(newUser *User) {
	//currUser.Name = newUser.Name
	currUser.Email = newUser.Email
	currUser.PhoneNumber = newUser.PhoneNumber
	currUser.Password = newUser.Password
	currUser.Address = newUser.Address
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

// CRUD in gorm
func (newUser *User) Register(db *gorm.DB, roleName string) (int, interface{}) {
	// Set role for user
	if err := newUser.SetUserRole(db, roleName); err != nil {
		statusCode, ErrorDB := helpers.DBError(err)
		return statusCode, ErrorDB
	}
	// Hash password
	if err := newUser.HashPassword(); err != nil {
		fError := helpers.FieldError{Field: "password", Message: "Cant hash password"}
		return 500, fError
	}
	// Create new User (Check validate Database)
	if err := db.Create(&newUser).Error; err != nil {
		statusCode, ErrorDB := helpers.DBError(err) // Hiện 3 lỗi khi nhập trùng cả 3 fields
		return statusCode, ErrorDB
	}
	return 201, nil
}

func (uLogin *LoginUser) Login(db *gorm.DB) (int, interface{}, *uint) { //StatusCode, Message, user.id
	// Check Field "name" in db
	var user User
	if err := db.Where("name = ?", uLogin.Name).First(&user).Error; err != nil {
		fError := helpers.FieldError{Field: "name", Message: "Name isn't already exist"}
		return 400, fError, nil
	} else {
		// Compare password
		if user.ComparePassword(uLogin.Password) == false {
			fError := helpers.FieldError{Field: "password", Message: "Incorrect Password"}
			return 400, fError, nil
		}
		return 201, nil, &user.Id
	}
}

func (currUser *User) Update(db *gorm.DB, newUser *User) (int, interface{}) {
	currUser.UpdateStruct(newUser)
	if err := currUser.HashPassword(); err != nil {
		fError := helpers.FieldError{Field: "password", Message: "Cant hash password"}
		return 500, fError
	}
	if err := db.Save(&currUser).Error; err != nil {
		statusCode, ErrorDB := helpers.DBError(err)
		return statusCode, ErrorDB
	} else {
		return 200, nil
	}
}
