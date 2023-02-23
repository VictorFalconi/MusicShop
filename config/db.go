package config

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"server/app/model"
)

var DB *gorm.DB

//Ceate admin & user role
func Init_Role() error {
	//admin
	var admin_role model.Role
	admin_role.Name = "admin"
	if err := DB.Create(&admin_role).Error; err != nil {
		//ErrorDB := helpers.DBError(err)
		fmt.Println("Error Database: Dont create admin role")
	}
	// user
	var user_role model.Role
	user_role.Name = "user"
	if err := DB.Create(&user_role).Error; err != nil {
		//ErrorDB := helpers.DBError(err)
		fmt.Println("Error Database: Dont create user role")
	}
	// employee
	var employee_role model.Role
	employee_role.Name = "employee"
	if err := DB.Create(&employee_role).Error; err != nil {
		//ErrorDB := helpers.DBError(err)
		fmt.Println("Error Database: Dont create employee role")
	}
	return nil
}

func ConnectDB() {
	//dsn := os.Getenv("DB_URL") "postgres://thanhliem:liem1234@localhost:41234/musicshop"
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("DATABASE_HOST"), os.Getenv("DATABASE_PORT"), os.Getenv("POSTGRES_DB"))
	fmt.Println(dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&model.Role{}, &model.User{})
	db.AutoMigrate(&model.Brand{})
	db.AutoMigrate(&model.Product{}, &model.Gallery{})
	db.SetupJoinTable(&model.Order{}, "Products", &model.OrderProducts{})
	db.AutoMigrate(&model.Order{})

	fmt.Println("Migration complete")
	DB = db

	// Create "admin, user, employee" role
	if errRole := Init_Role(); errRole != nil {
		fmt.Println("Dont create role")
	}
}
