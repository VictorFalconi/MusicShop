package model

import (
	"gorm.io/gorm"
	"server/helpers"
	"time"
)

//							                             1-n	  n-n     n-n     n-n     n-n    n-n
//#Name, Price, Thumbnail, Description, Year, Quality, Category, Brand, Country, Format, Genre, Style

type Product struct {
	Id          uint    `json:"ID"          gorm:"primary_key"`
	Name        string  `json:"name"        gorm:"unique;not null"          validate:"required,min=4,max=128"`
	Quantity    int     `json:"quantity"    gorm:"not null;default:0"       validate:""`
	Price       float32 `json:"price"       gorm:"not null"                 validate:"required"`
	Discount    float32 `json:"discount"    gorm:"not null;default:0.0"     validate:""`
	Thumbnail   string  `json:"thumbnail"   gorm:""                         validate:""`
	Description string  `json:"description" gorm:""                         validate:""`
	Year        string  `json:"year"        gorm:""                         validate:"len=4"`
	Quality     string  `json:"quality"     gorm:""                         validate:""`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `            gorm:"index"`

	Galleries []Gallery `json:"galleries"       gorm:"foreignKey:ProductId;references:Id"` //Product 1-n Gallery
	Brands    []Brand   `json:"brands"          gorm:"many2many:product_brands"`           //Product n-n Brand
}

type Products []Product

type Gallery struct {
	Id        uint   `json:"ID"          gorm:"primary_key"`
	Thumbnail string `json:"thumbnail"   gorm:"not null"             validate:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	ProductId uint
}

// CRUD

func (product *Product) Create(db *gorm.DB) (int, interface{}) {
	if err := db.Create(&product).Error; err != nil {
		statusCode, ErrorDB := helpers.DBError(err)
		return statusCode, ErrorDB
	}
	return 201, nil
}

func (products *Products) Reads(db *gorm.DB) (int, interface{}, interface{}) {
	if err := db.Preload("Galleries").Preload("Brands").Find(&products).Error; err != nil {
		statusCode, ErrorDB := helpers.DBError(err)
		return statusCode, ErrorDB, nil
	}
	return 200, nil, products
}

func (product *Product) Read(db *gorm.DB, id string) (int, interface{}, interface{}) {
	if err := db.Preload("Galleries").Preload("Brands").Where("id = ?", id).First(&product).Error; err != nil {
		return 404, helpers.FieldError{Field: "", Message: "URL not found"}, nil
	}
	return 200, nil, product
}

func (product *Product) UpdateStruct(newProduct *Product) {
	product.Name = newProduct.Name
	product.Quantity = newProduct.Quantity
	product.Price = newProduct.Price
	product.Discount = newProduct.Discount
	product.Thumbnail = newProduct.Thumbnail
	product.Description = newProduct.Description
	product.Year = newProduct.Year
	product.Quality = newProduct.Quality
}

func (product *Product) Update(db *gorm.DB, newProduct *Product) (int, interface{}) {
	// Map
	product.UpdateStruct(newProduct)
	// Update
	if err := db.Model(&product).Association("Galleries").Replace(newProduct.Galleries); err != nil {
		statusCode, ErrorDB := helpers.DBError(err)
		return statusCode, ErrorDB
	}
	if err := db.Model(&product).Association("Brands").Replace(newProduct.Brands); err != nil {
		statusCode, ErrorDB := helpers.DBError(err)
		return statusCode, ErrorDB
	}
	if err := db.Save(&product).Error; err != nil {
		statusCode, ErrorDB := helpers.DBError(err)
		return statusCode, ErrorDB
	}
	return 200, nil
}

func (product *Product) Delete(db *gorm.DB) (int, interface{}) {
	// Delete
	if err := db.Delete(&product).Error; err != nil {
		statusCode, ErrorDB := helpers.DBError(err)
		return statusCode, ErrorDB
	} else {
		return 200, nil
	}
}

// Create products from Excel file

func String2Galleries(str string, productId uint) []Gallery {
	slice := helpers.String2Slice(str)
	var galleries []Gallery
	if len(slice) == 0 {
		return galleries
	}

	for _, name := range slice {
		if name == "NULL" || name == "" || name == " " {
			continue
		}
		gallery := Gallery{
			Thumbnail: name,
			ProductId: productId}
		galleries = append(galleries, gallery)
	}
	return galleries
}

func String2Brands(db *gorm.DB, str string) ([]Brand, []helpers.FieldError) {
	slice := helpers.String2Slice(str)
	var brands []Brand
	var fielderrors []helpers.FieldError
	for _, name := range slice {
		if name == "NULL" || name == "" || name == " " {
			continue
		}
		var brand Brand
		if err := db.Where("name = ?", name).First(&brand).Error; err == nil {
			brands = append(brands, brand)
		} else {
			fielderrors = append(fielderrors, helpers.FieldError{Field: "Brand", Message: "Dont add '" + name + "' into Brand field !"})
		}
	}
	return brands, fielderrors
}
