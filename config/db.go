package config

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"server/app/models"
)

var DB *gorm.DB

func ConnectDB() {
	//dsn := os.Getenv("DB_URL") "postgres://thanhliem:liem1234@localhost:41234/musicshop"
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("DATABASE_HOST"), os.Getenv("DATABASE_PORT"), os.Getenv("POSTGRES_DB"))
	fmt.Println(dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.User{})
	fmt.Println("Migration complete")
	DB = db
}
