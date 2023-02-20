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

// GetNameRoleUser : Get Name_Role of User
func (u *User) GetNameRoleUser(db *gorm.DB) (error, string) {
	var role Role
	if err := db.Where("id = ?", u.RoleId).First(&role).Error; err != nil {
		return err, ""
	}
	return nil, role.Name
}

// SetUserRole : Set Role of User form Name_Role:
func (u *User) SetUserRole(db *gorm.DB, roleName string) error {
	var role Role
	if err := db.Where("name = ?", roleName).First(&role).Error; err != nil {
		return err
	}
	u.RoleId = role.Id
	return nil
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

// ComparePassword : Compare between password and HashPassword
func (u *User) ComparePassword(password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return false
	}
	return true
}

// CRUD

func (u *User) Register(db *gorm.DB, roleName string) (int, interface{}) {
	// Set role for user
	if err := u.SetUserRole(db, roleName); err != nil {
		statusCode, ErrorDB := helpers.DBError(err)
		return statusCode, ErrorDB
	}
	// Hash password
	if err := u.HashPassword(); err != nil {
		fError := helpers.FieldError{Field: "password", Message: "Cant hash password"}
		return 500, fError
	}
	// Create new User (Check validate Database)
	if err := db.Create(&u).Error; err != nil {
		statusCode, ErrorDB := helpers.DBError(err)
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

func (u *User) UpdateStruct(newUser *User) {
	//currUser.Name = newUser.Name
	u.Email = newUser.Email
	u.PhoneNumber = newUser.PhoneNumber
	u.Password = newUser.Password
	u.Address = newUser.Address
}

func (u *User) Update(db *gorm.DB, newUser *User) (int, interface{}) {
	u.UpdateStruct(newUser)
	if err := u.HashPassword(); err != nil {
		fError := helpers.FieldError{Field: "password", Message: "Cant hash password"}
		return 500, fError
	}
	if err := db.Save(&u).Error; err != nil {
		statusCode, ErrorDB := helpers.DBError(err)
		return statusCode, ErrorDB
	} else {
		return 200, nil
	}
}
