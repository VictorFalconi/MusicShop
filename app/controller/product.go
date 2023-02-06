package controller

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/xuri/excelize/v2"
	"io"
	"log"
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

	if filepath.Ext(file.Filename) != ".xlsx" {
		helpers.RespondJSON(ctx, 400, "Error type!", "Type file is must Excel (xlsx)", nil)
		return
	}
	src, err := file.Open()
	if err != nil {
		helpers.RespondJSON(ctx, 400, "Error type!", "Dont read this", nil)
		return
	}
	defer src.Close()

	buf := bytes.Buffer{}
	io.Copy(&buf, src)

	// Read Excel file
	xlsx, err := excelize.OpenReader(&buf)
	if err != nil {
		helpers.RespondJSON(ctx, 500, "Error file!", "Internal Server Error: "+err.Error(), nil)
		return
	}
	// Get all the rows in the first sheet
	rows, err := xlsx.GetRows("Sheet1")
	if err != nil {
		helpers.RespondJSON(ctx, 500, "Error file!", "Internal Server Error: "+err.Error(), nil)
		return
	}

	var listDataErr []helpers.LineError
	for i, row := range rows {
		if i == 0 {
			continue
		}

		brands, fieldErrorBrands := helpers.String2Brands(row[7])
		product := models.Product{
			Name:        row[0],
			Price:       helpers.String2Float(row[1]),
			Thumbnail:   row[2],
			Description: row[3],
			Year:        row[4],
			Quality:     row[5],
			//Gallery:   (row[6]),
			Brands: brands}

		// Create new Product
		if errdb := config.DB.Create(&product).Error; errdb != nil || len(fieldErrorBrands) != 0 {
			ErrorDB := helpers.DBError(errdb)
			lineErr := helpers.LineError{
				Line:    i + 1,
				Message: append(ErrorDB, fieldErrorBrands...)}
			listDataErr = append(listDataErr, lineErr)
		} else {
			// Create Galleries for Product
			var galleries []models.Gallery
			galleries = helpers.String2Galleries(row[6], product.Id)
			if len(galleries) != 0 {
				if errGaleery := config.DB.Create(&galleries).Error; errGaleery != nil {
					ErrorDB := helpers.DBError(errGaleery)
					lineErr := helpers.LineError{
						Line:    i + 1,
						Message: ErrorDB}
					listDataErr = append(listDataErr, lineErr)
				}
			}
		}
	}
	if len(listDataErr) != 0 {
		helpers.RespondJSON(ctx, 207, "Created some products successful!", nil, listDataErr)
		return
	} else {
		helpers.RespondJSON(ctx, 201, "Created all products successful!", nil, nil)
		return
	}
}

//// Read many excel files
//form, err := ctx.MultipartForm()
//if err != nil {
//	helpers.RespondJSON(ctx, 400, "Error file type!", err.Error(), nil)
//	return
//}
//files := form.File["file"]
//if len(files) != 2 {
//	helpers.RespondJSON(ctx, 400, "Error type!", "Must 2 Excel file", nil)
//	return
//}
//if filepath.Ext(files[0].Filename) != ".xlsx" || filepath.Ext(files[1].Filename) != ".xlsx" {
//	helpers.RespondJSON(ctx, 400, "Error type!", "Type file is must Excel (xlsx)", nil)
//	return
//}
//var bufs []*bytes.Buffer
//for _, file := range files {
//	fmt.Println(file.Filename)
//	src, err := file.Open()
//	if err != nil {
//		helpers.RespondJSON(ctx, 400, "Error type!", "Dont read this", nil)
//		return
//	}
//	defer src.Close()
//
//	buf := bytes.Buffer{}
//	io.Copy(&buf, src)
//	bufs = append(bufs, &buf)
//}

func CreateProduct_FromFile1(ctx *gin.Context) {
	// Read file
	file, err := ctx.FormFile("file")
	if err != nil {
		helpers.RespondJSON(ctx, 400, "Error file type!", err.Error(), nil)
		return
	}
	log.Println(file.Filename)
	log.Println(filepath.Ext(file.Filename))
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

	// Create a new buffered reader for the file
	reader := bufio.NewReader(csvFile)

	// Read all lines in the file
	for {
		line, err := reader.ReadString('\n') // '\r\n'
		if err != nil {
			break
		}
		fmt.Println(line)
	}

	//// read csv
	//reader := csv.NewReader(csvFile)
	//reader.Comma = ','
	//reader.LazyQuotes = true
	//
	//// Struct file csv products
	//var products []models.Product
	//for {
	//	record, err := reader.Read()
	//	if err == io.EOF {
	//		break
	//	} else if err != nil {
	//		helpers.RespondJSON(ctx, 400, "Error data!", err.Error(), nil)
	//		return
	//	}
	//	fmt.Println(record, len(record))
	//	product := models.Product{Name: record[0]} //????  ---- ????
	//	products = append(products, product)
	//}
	//// Create new Products
	//if err := config.DB.Create(&products).Error; err != nil {
	//	ErrorDB := helpers.DBError(err)
	//	helpers.RespondJSON(ctx, 400, "Error Database", ErrorDB, nil)
	//	return
	//} else {
	//	helpers.RespondJSON(ctx, 201, "Created brand successful!", nil, nil)
	//	return
	//}
}
