package helpers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ReponseData struct {
	Status  int
	Message string
	Error   interface{}
	Data    interface{}
}

func RespondJSON(w *gin.Context, status int, message string, errors interface{}, payload interface{}) {
	var res ReponseData
	res.Status = status
	res.Message = message
	res.Error = errors
	res.Data = payload

	w.JSON(200, res)
}

// String to list errors from validator.Struct.Error
type FieldError struct {
	Field   string
	Message string
}

func MessageForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	//Validate
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	case "min":
		return "This field is short"
	case "max":
		return "This field is long"
	}
	return fe.Error()
}

func StringToListErrors(errs validator.ValidationErrors) interface{} {
	dictErrors := make([]FieldError, len(errs))
	for index, e := range errs {
		dictErrors[index] = FieldError{Field: e.Field(), Message: MessageForTag(e)}
	}
	return dictErrors
}
