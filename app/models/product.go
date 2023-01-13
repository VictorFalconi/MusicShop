package models

import (
	"gorm.io/gorm"
	"time"
)

//							                             1-n	  n-n     n-n     n-n     n-n    n-n
//#Name, Price, Thumbnail, Description, Year, Quality, Category, Brand, Country, Format, Genre, Style

type Product struct {
	Id          uint    `json:"ID"          gorm:"primary_key"`
	Name        string  `json:"name"        gorm:"unique;not null"          validate:"required,min=4,max=128"`
	Price       float32 `json:"price"       gorm:"not null"                 validate:"required"`
	Discount    float32 `json:"discount"    gorm:"not null;default:0.0"     validate:""`
	Thumbnail   string  `json:"thumbnail"   gorm:""                         validate:""`
	Description string  `json:"description" gorm:""                         validate:""`
	Year        string  `json:"year"        gorm:""                         validate:"len=4"`
	Quality     string  `json:"quality"     gorm:""                         validate:""`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	Gallery Gallery `gorm:"foreignKey:ProductId;references:Id"` //Product 1-n Gallery
	Brands  []Brand `gorm:"many2many:product_brands;"`          //Product n-n Brand
}

type Gallery struct {
	Id        uint   `json:"ID"          gorm:"primary_key"`
	Thumbnail string `json:"thumbnail"   gorm:"not null"             validate:""`

	ProductId uint
}

type Brand struct {
	Id   uint   `json:"ID"   form:"ID"     gorm:"primary_key"`
	Name string `json:"name" form:"name"   gorm:"unique;not null"  validate:"required,max=128"`
}
