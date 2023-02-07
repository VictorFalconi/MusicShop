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

func CreateBrand(ctx *gin.Context) {
	var brand models.Brand
	// Check data type
	if err := helpers.DataContentType(ctx, &brand); err != nil {
		helpers.RespondJSON(ctx, 400, "Error data type!", err.Error(), nil)
		return
	}
	// Check validate field
	if err := validator.New().Struct(&brand); err != nil {
		listErrors := helpers.ValidateErrors(err.(validator.ValidationErrors))
		helpers.RespondJSON(ctx, 400, "Errors validate!", listErrors, nil)
		return
	}
	// Create new Brand (Check validate Database)
	if err := config.DB.Create(&brand).Error; err != nil {
		statusCode, ErrorDB := helpers.DBError(err)
		helpers.RespondJSON(ctx, statusCode, "Error Database", ErrorDB, nil)
		return
	} else {
		helpers.RespondJSON(ctx, 201, "Created brand successful!", nil, nil)
		return
	}
}

func ReadBrands(ctx *gin.Context) {
	var brands []models.Brand
	if err := config.DB.Find(&brands).Error; err != nil {
		statusCode, ErrorDB := helpers.DBError(err)
		helpers.RespondJSON(ctx, statusCode, "Error Database", ErrorDB, nil)
		return
	} else {
		helpers.RespondJSON(ctx, 200, "Read brands successful!", nil, &brands)
		return
	}
}

func ReadBrand(ctx *gin.Context) {
	var brand models.Brand
	if err := config.DB.Where("id = ?", ctx.Param("id")).First(&brand).Error; err != nil {
		//ErrorDB := helpers.DBError(err)
		helpers.RespondJSON(ctx, 404, "Error URL", "URL not found", nil)
		return
	} else {
		helpers.RespondJSON(ctx, 200, "Read brand successful!", nil, &brand)
		return
	}
}

func UpdateBrand(ctx *gin.Context) {
	// Find brand
	var currBrand models.Brand
	config.DB.Where("id = ?", ctx.Param("id")).First((&currBrand))
	// Get request
	var newBrand models.Brand
	if err := helpers.DataContentType(ctx, &newBrand); err != nil {
		helpers.RespondJSON(ctx, 400, "Error data type!", err.Error(), nil)
		return
	}
	if err := validator.New().Struct(&newBrand); err != nil {
		listErrors := helpers.ValidateErrors(err.(validator.ValidationErrors))
		helpers.RespondJSON(ctx, 400, "Errors validate!", listErrors, nil)
		return
	}
	// Update
	currBrand.Name = newBrand.Name
	if err := config.DB.Save(&currBrand).Error; err != nil {
		statusCode, ErrorDB := helpers.DBError(err)
		helpers.RespondJSON(ctx, statusCode, "Error Database", ErrorDB, nil)
		return
	} else {
		helpers.RespondJSON(ctx, 200, "Updated brand successful!", nil, nil)
		return
	}
}

func DeleteBrand(ctx *gin.Context) {
	// Find Brand
	var currBrand models.Brand
	if err := config.DB.Where("id = ?", ctx.Param("id")).First(&currBrand).Error; err != nil {
		helpers.RespondJSON(ctx, 404, "Error URL", "URL not found", nil)
		return
	}
	// Delete
	if err := config.DB.Delete(&currBrand).Error; err != nil {
		statusCode, ErrorDB := helpers.DBError(err)
		helpers.RespondJSON(ctx, statusCode, "Error Database", ErrorDB, nil)
		return
	} else {
		helpers.RespondJSON(ctx, 204, "Deleted brand successful!", nil, nil)
		return
	}
}

func CreateBrand_FromFile(ctx *gin.Context) {
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
	var brands []models.Brand
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			helpers.RespondJSON(ctx, 400, "Error data!", err.Error(), nil)
			return
		}
		brand := models.Brand{Name: record[0]}
		brands = append(brands, brand)
	}
	// Create new Brands
	if err := config.DB.Create(&brands).Error; err != nil {
		statusCode, ErrorDB := helpers.DBError(err)
		helpers.RespondJSON(ctx, statusCode, "Error Database", ErrorDB, nil)
		return
	} else {
		helpers.RespondJSON(ctx, 201, "Created brand successful!", nil, nil)
		return
	}
}
