package controller

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/xuri/excelize/v2"
	"io"
	"log"
	"path/filepath"
	"reflect"
	"server/app/models"
	"server/config"
	"server/helpers"
)

// CreateProduct : Create Product
func CreateProduct(ctx *gin.Context) {
	var product models.Product
	// Check data type
	if err := helpers.DataContentType(ctx, &product); err != nil {
		helpers.RespondJSON(ctx, 400, helpers.StatusCodeFromInt(400), err.Error(), nil)
		return
	}
	// Check validate field
	if err := validator.New().Struct(&product); err != nil {
		listErrors := helpers.ValidateErrors(err.(validator.ValidationErrors))
		helpers.RespondJSON(ctx, 400, helpers.StatusCodeFromInt(400), listErrors, nil)
		return
	}
	// Create
	statusCode, Message := product.Create(config.DB)
	helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), Message, nil)
	return
}

func ReadProducts(ctx *gin.Context) {
	var products models.Products
	// Reads
	statusCode, Message, output := products.Reads(config.DB)
	helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), Message, output)
	return
}

func ReadProduct(ctx *gin.Context) {
	var product models.Product
	// Read
	statusCode, Message, output := product.Read(config.DB, ctx.Param("id"))
	helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), Message, output)
	return
}

func UpdateProduct(ctx *gin.Context) {
	// Find product
	var currProduct models.Product
	statusCode, Message, _ := currProduct.Read(config.DB, ctx.Param("id"))
	if statusCode != 200 {
		helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), Message, nil)
		return
	}
	// Get request
	var newProduct models.Product
	if err := helpers.DataContentType(ctx, &newProduct); err != nil {
		helpers.RespondJSON(ctx, 400, helpers.StatusCodeFromInt(400), err.Error(), nil)
		return
	}
	if err := validator.New().Struct(&newProduct); err != nil {
		listErrors := helpers.ValidateErrors(err.(validator.ValidationErrors))
		helpers.RespondJSON(ctx, 400, helpers.StatusCodeFromInt(400), listErrors, nil)
		return
	}
	// Update
	statusCode, Message = currProduct.Update(config.DB, &newProduct)
	helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), Message, nil)
	return
}

func DeleteProduct(ctx *gin.Context) {
	// Find Product
	var currProduct models.Product
	statusCode, Message, _ := currProduct.Read(config.DB, ctx.Param("id"))
	if statusCode != 200 {
		helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), Message, nil)
		return
	}
	//Delete
	statusCode, Message = currProduct.Delete(config.DB)
	helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), Message, nil)
	return
}

func CreateProduct_FromFile(ctx *gin.Context) {
	// Read file
	file, err := ctx.FormFile("file")
	if err != nil {
		helpers.RespondJSON(ctx, 400, helpers.StatusCodeFromInt(400), err.Error(), nil)
		return
	}

	if filepath.Ext(file.Filename) != ".xlsx" {
		helpers.RespondJSON(ctx, 400, helpers.StatusCodeFromInt(400), "Type file is must Excel (xlsx)", nil)
		return
	}
	src, err := file.Open()
	if err != nil {
		helpers.RespondJSON(ctx, 400, helpers.StatusCodeFromInt(400), "Dont read this", nil)
		return
	}
	defer src.Close()

	buf := bytes.Buffer{}
	io.Copy(&buf, src)

	// Read Excel file
	xlsx, err := excelize.OpenReader(&buf)
	if err != nil {
		helpers.RespondJSON(ctx, 500, helpers.StatusCodeFromInt(500), err.Error(), nil)
		return
	}
	// Get all the rows in the first sheet
	rows, err := xlsx.GetRows("Sheet1")
	if err != nil {
		helpers.RespondJSON(ctx, 400, helpers.StatusCodeFromInt(400), err.Error(), nil)
		return
	}

	var listDataErr []helpers.LineError
	for i, row := range rows {
		if i == 0 {
			continue
		}
		// Len of row
		numField := reflect.ValueOf(models.Product{}).NumField()
		log.Println(numField, len(row), row)
		if len(row) != (numField - 4) { // ID, Amount, Discount, CUD time
			listDataErr = append(listDataErr, helpers.LineError{Line: i + 1, Message: "Invalid length or empty field"})
			continue
		}

		brands, fieldErrorBrands := models.String2Brands(config.DB, row[9])
		product := models.Product{
			Name:        row[0],
			Quantity:    helpers.String2Int(row[1]),
			Price:       helpers.String2Float(row[2]),
			Discount:    helpers.String2Float(row[3]),
			Thumbnail:   row[4],
			Description: row[5],
			Year:        row[6],
			Quality:     row[7],
			//Gallery:   (row[8]),
			Brands: brands}

		// Create new Product
		if errdb := config.DB.Create(&product).Error; errdb != nil || len(fieldErrorBrands) != 0 {
			_, ErrorDB := helpers.DBError(errdb)
			lineErr := helpers.LineError{
				Line:    i + 1,
				Message: append(ErrorDB, fieldErrorBrands...)}
			listDataErr = append(listDataErr, lineErr)
		} else {
			// Create Galleries for Product
			var galleries []models.Gallery
			galleries = models.String2Galleries(row[8], product.Id)
			if len(galleries) != 0 {
				if errGaleery := config.DB.Create(&galleries).Error; errGaleery != nil {
					_, ErrorDB := helpers.DBError(errGaleery)
					lineErr := helpers.LineError{
						Line:    i + 1,
						Message: ErrorDB}
					listDataErr = append(listDataErr, lineErr)
				}
			}
		}
	}
	if len(listDataErr) != 0 {
		helpers.RespondJSON(ctx, 207, helpers.StatusCodeFromInt(207), nil, listDataErr)
		return
	} else {
		helpers.RespondJSON(ctx, 201, helpers.StatusCodeFromInt(201), nil, nil)
		return
	}
}
