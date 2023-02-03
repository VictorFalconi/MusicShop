package controller

import (
	"encoding/csv"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"io"
	"path/filepath"
	"server/app/models"
	"server/config"
	"server/helpers"
)

// CreateProduct : Create Product
func CreateProduct(ctx *gin.Context) {
	var product models.Product
	// Check data type
	if err := helpers.DataContentType(ctx, &product); err != nil {
		helpers.RespondJSON(ctx, 400, "Error data type!", err.Error(), nil)
		return
	}
	// Check validate field
	if err := validator.New().Struct(&product); err != nil {
		listErrors := helpers.ValidateErrors(err.(validator.ValidationErrors))
		helpers.RespondJSON(ctx, 400, "Errors validate!", listErrors, nil)
		return
	}
	// Create new Product (Check validate Database)
	if err := config.DB.Create(&product).Error; err != nil {
		ErrorDB := helpers.DBError(err)
		helpers.RespondJSON(ctx, 400, "Error Database", ErrorDB, nil)
		return
	} else {
		helpers.RespondJSON(ctx, 201, "Created product successful!", nil, nil)
		return
	}
}

func ReadProducts(ctx *gin.Context) {
	var products []models.Product
	if err := config.DB.Preload("Brands").Find(&products).Error; err != nil {
		ErrorDB := helpers.DBError(err)
		helpers.RespondJSON(ctx, 400, "Error Database", ErrorDB, nil)
		return
	} else {
		helpers.RespondJSON(ctx, 200, "Read products successful!", nil, &products)
		return
	}
}

func ReadProduct(ctx *gin.Context) {
	var product models.Product
	if err := config.DB.Preload("Brands").Where("id = ?", ctx.Param("id")).First(&product).Error; err != nil {
		//ErrorDB := helpers.DBError(err)
		helpers.RespondJSON(ctx, 404, "Error URL", "URL not found", nil)
		return
	} else {
		helpers.RespondJSON(ctx, 200, "Read product successful!", nil, &product)
		return
	}
}

func UpdateProduct(ctx *gin.Context) {
	// Find product
	var currProduct models.Product
	config.DB.Preload("Brands").Where("id = ?", ctx.Param("id")).First(&currProduct)
	// Get request
	var newProduct models.Product
	if err := helpers.DataContentType(ctx, &newProduct); err != nil {
		helpers.RespondJSON(ctx, 400, "Error data type!", err.Error(), nil)
		return
	}
	if err := validator.New().Struct(&newProduct); err != nil {
		listErrors := helpers.ValidateErrors(err.(validator.ValidationErrors))
		helpers.RespondJSON(ctx, 400, "Errors validate!", listErrors, nil)
		return
	}
	// Update
	currProduct.UpdateStruct(&newProduct)
	config.DB.Model(&currProduct).Association("Brands").Replace(newProduct.Brands)
	if err := config.DB.Save(&currProduct).Error; err != nil {
		ErrorDB := helpers.DBError(err)
		helpers.RespondJSON(ctx, 400, "Error Database", ErrorDB, nil)
		return
	} else {
		helpers.RespondJSON(ctx, 200, "Updated product successful!", nil, nil)
		return
	}
}

func DeleteProduct(ctx *gin.Context) {
	// Find Product
	var currProduct models.Product
	if err := config.DB.Preload("Brands").Where("id = ?", ctx.Param("id")).First(&currProduct).Error; err != nil {
		helpers.RespondJSON(ctx, 404, "Error URL", "URL not found", nil)
		return
	}
	// Delete
	if err := config.DB.Delete(&currProduct).Error; err != nil {
		ErrorDB := helpers.DBError(err)
		helpers.RespondJSON(ctx, 400, "Error Database", ErrorDB, nil)
		return
	} else {
		helpers.RespondJSON(ctx, 204, "Deleted product successful!", nil, nil)
		return
	}
}

func CreateProduct_FromFile(ctx *gin.Context) {
	// Read file
	file, err := ctx.FormFile("file")
	if err != nil {
		helpers.RespondJSON(ctx, 400, "Error file type!", err.Error(), nil)
		return
	}
	//log.Println(file.Filename)
	//log.Println(filepath.Ext(file.Filename))
	if filepath.Ext(file.Filename) != ".csv" {
		helpers.RespondJSON(ctx, 400, "Error type!", "Type file is must CSV", nil)
		return
	}
	// Read the contents of the file into a variable
	csvFile, err := file.Open()
	if err != nil {
		helpers.RespondJSON(ctx, 400, "Error file!", err.Error(), nil)
		return
	}
	defer csvFile.Close()

	// read csv
	reader := csv.NewReader(csvFile)

	// Struct file csv products
	var products []models.Product
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			helpers.RespondJSON(ctx, 400, "Error data!", err.Error(), nil)
			return
		}
		product := models.Product{Name: record[0]} //????  ---- ????
		products = append(products, product)
	}
	// Create new Products
	if err := config.DB.Create(&products).Error; err != nil {
		ErrorDB := helpers.DBError(err)
		helpers.RespondJSON(ctx, 400, "Error Database", ErrorDB, nil)
		return
	} else {
		helpers.RespondJSON(ctx, 201, "Created brand successful!", nil, nil)
		return
	}
}
