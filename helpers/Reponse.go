package helpers

import (
	"github.com/gin-gonic/gin"
	"server/app/models"
	"server/config"
	"strconv"
	"strings"
)

// Response: Response JSON data to client
type ResponseData struct {
	Status  int
	Message string
	Error   interface{}
	Data    interface{}
}

func RespondJSON(w *gin.Context, status int, message string, errors interface{}, payload interface{}) {
	var res ResponseData
	res.Status = status
	res.Message = message
	res.Error = errors
	res.Data = payload

	w.JSON(status, res)
}

//Respond data : return positon & fields when dont create elements from file (excel, csv)
type LineError struct {
	Line    int
	Message interface{}
}

// String to float
func String2Float(str string) float32 {
	fPrice, err := strconv.ParseFloat(str, 64)
	if err != nil {
		fPrice = 0.0
	}
	return float32(fPrice)
}

// String to entity "LB CLUB, Wanner Music VN" -> []Brand{}
func String2Slice(str string) []string {
	slice := strings.Split(str, ",")
	for i, name := range slice {
		slice[i] = strings.TrimSpace(name)
	}
	return slice
}

func String2Galleries(str string, product_id uint) []models.Gallery {
	slice := String2Slice(str)
	var galleries []models.Gallery
	if len(slice) == 0 {
		return galleries
	}

	for _, name := range slice {
		if name == "NULL" || name == "" || name == " " {
			continue
		}
		gallery := models.Gallery{
			Thumbnail: name,
			ProductId: product_id}
		galleries = append(galleries, gallery)
	}
	return galleries
}

func String2Brands(str string) ([]models.Brand, []FieldError) {
	slice := String2Slice(str)
	var brands []models.Brand
	var fielderrors []FieldError
	for _, name := range slice {
		if name == "NULL" || name == "" || name == " " {
			continue
		}
		var brand models.Brand
		if err := config.DB.Where("name = ?", name).First(&brand).Error; err == nil {
			brands = append(brands, brand)
		} else {
			fielderrors = append(fielderrors, FieldError{Field: "Brand", Message: "Dont add " + name + " into Brand field !"})
		}
	}
	return brands, fielderrors
}
