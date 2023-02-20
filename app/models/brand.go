package models

import (
	"gorm.io/gorm"
	"server/helpers"
	"time"
)

type Brand struct {
	Id        uint   `json:"ID"   form:"ID"     gorm:"primary_key"`
	Name      string `json:"name" form:"name"   gorm:"unique;not null"  validate:"required,max=128"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Brands []Brand

// CRUD

func (brand *Brand) Create(db *gorm.DB) (int, interface{}) {
	if err := db.Create(&brand).Error; err != nil {
		statusCode, ErrorDB := helpers.DBError(err)
		return statusCode, ErrorDB
	}
	return 201, nil
}

func (brands *Brands) Reads(db *gorm.DB) (int, interface{}, interface{}) {
	if err := db.Find(&brands).Error; err != nil {
		statusCode, ErrorDB := helpers.DBError(err)
		return statusCode, ErrorDB, nil
	}
	return 200, nil, brands
}

func (brand *Brand) Read(db *gorm.DB, id string) (int, interface{}, interface{}) {
	if err := db.Where("id = ?", id).First(&brand).Error; err != nil {
		//statusCode, ErrorDB := helpers.DBError(err)
		return 404, helpers.FieldError{Field: "", Message: "URL not found"}, nil
	}
	return 200, nil, brand
}

func (brand *Brand) UpdateStruct(newBrand *Brand) {
	brand.Name = newBrand.Name
}

func (brand *Brand) Update(db *gorm.DB, newBrand *Brand) (int, interface{}) {
	brand.UpdateStruct(newBrand)
	if err := db.Save(&brand).Error; err != nil {
		statusCode, ErrorDB := helpers.DBError(err)
		return statusCode, ErrorDB
	}
	return 200, nil
}

func (brand *Brand) Delete(db *gorm.DB) (int, interface{}) {
	// Delete
	if err := db.Delete(&brand).Error; err != nil {
		statusCode, ErrorDB := helpers.DBError(err)
		return statusCode, ErrorDB
	}
	return 200, nil
}

// Creates from csv
func (brands *Brands) Creates(db *gorm.DB) (int, interface{}) {
	if err := db.Create(&brands).Error; err != nil {
		statusCode, ErrorDB := helpers.DBError(err)
		return statusCode, ErrorDB
	}
	return 201, nil
}
