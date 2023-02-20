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
		helpers.RespondJSON(ctx, 400, helpers.StatusCodeFromInt(400), err.Error(), nil)
		return
	}
	// Check validate field
	if err := validator.New().Struct(&brand); err != nil {
		listErrors := helpers.ValidateErrors(err.(validator.ValidationErrors))
		helpers.RespondJSON(ctx, 400, helpers.StatusCodeFromInt(400), listErrors, nil)
		return
	}
	// Create
	statusCode, Message := brand.Create(config.DB)
	helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), Message, nil)
	return
}

func ReadBrands(ctx *gin.Context) {
	var brands models.Brands
	// Reads
	statusCode, Message, output := brands.Reads(config.DB)
	helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), Message, output)
	return
}

func ReadBrand(ctx *gin.Context) {
	var brand models.Brand
	//Read
	statusCode, Message, output := brand.Read(config.DB, ctx.Param("id"))
	helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), Message, output)
	return
}

func UpdateBrand(ctx *gin.Context) {
	// Find brand
	var currBrand models.Brand
	statusCode, Message, _ := currBrand.Read(config.DB, ctx.Param("id"))
	if statusCode != 200 {
		helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), Message, nil)
		return
	}
	// Get request
	var newBrand models.Brand
	if err := helpers.DataContentType(ctx, &newBrand); err != nil {
		helpers.RespondJSON(ctx, 400, helpers.StatusCodeFromInt(400), err.Error(), nil)
		return
	}
	if err := validator.New().Struct(&newBrand); err != nil {
		listErrors := helpers.ValidateErrors(err.(validator.ValidationErrors))
		helpers.RespondJSON(ctx, 400, helpers.StatusCodeFromInt(400), listErrors, nil)
		return
	}
	// Update
	statusCode, Message = currBrand.Update(config.DB, &newBrand)
	helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), Message, nil)
	return
}

func DeleteBrand(ctx *gin.Context) {
	// Find Brand
	var currBrand models.Brand
	statusCode, Message, _ := currBrand.Read(config.DB, ctx.Param("id"))
	if statusCode != 200 {
		helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), Message, nil)
		return
	}
	//Delete
	statusCode, Message = currBrand.Delete(config.DB)
	helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), Message, nil)
	return
}

func CreateBrand_FromFile(ctx *gin.Context) {
	// Read file
	file, err := ctx.FormFile("file")
	if err != nil {
		helpers.RespondJSON(ctx, 400, helpers.StatusCodeFromInt(400), err.Error(), nil)
		return
	}
	//log.Println(file.Filename)
	//log.Println(filepath.Ext(file.Filename))
	if filepath.Ext(file.Filename) != ".csv" {
		helpers.RespondJSON(ctx, 400, helpers.StatusCodeFromInt(400), "Type file is must CSV", nil)
		return
	}
	// Read the contents of the file into a variable
	csvFile, err := file.Open()
	if err != nil {
		helpers.RespondJSON(ctx, 400, helpers.StatusCodeFromInt(400), err.Error(), nil)
		return
	}
	defer csvFile.Close()

	// read csv
	reader := csv.NewReader(csvFile)
	var brands models.Brands
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			helpers.RespondJSON(ctx, 400, helpers.StatusCodeFromInt(400), err.Error(), nil)
			return
		}
		brand := models.Brand{Name: record[0]}
		brands = append(brands, brand)
	}
	// Creates
	statusCode, Message := brands.Creates(config.DB)
	helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), Message, nil)
	return
}
