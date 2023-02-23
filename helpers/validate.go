package helpers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
	"regexp"
	"strings"
)

// StatusCode
func StatusCodeFromInt(value int) string {
	switch value {
	case 200:
		return "The request was successful"
	case 201:
		return "The new resource has been created"
	case 204:
		return "This resource has been deleted"
	case 207:
		return "The request was successful, but some resource is error"
	case 400:
		return "The request was invalid"
	case 401:
		return "The request requires authentication"
	case 403:
		return "The request not authorized"
	case 404:
		return "URL not found"
	case 500:
		return "Internal Server Error"
	case 502:
		return "Bad Gateway"
	case 503:
		return "The server is currently unavailable"
	}
	return "Status code: " + string(value)
}

//Data Type: Get multiple data type form client (form-data, XML, JSON)
func DataContentType(ctx *gin.Context, entity interface{}) error {
	switch ctx.ContentType() {
	case "application/json":
		return ctx.ShouldBindJSON(entity)
	case "application/xml":
		return ctx.ShouldBindXML(entity)
	case "multipart/form-data":
		return ctx.ShouldBind(entity)
	}
	return nil
}

//Validate Field: Error Handling
type FieldError struct {
	Field   string
	Message string
}

func messageForTag(fe validator.FieldError) string {
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
	case "len":
		return "Invalid length"
	}
	return fe.Error()
}

func ValidateErrors(errs validator.ValidationErrors) interface{} {
	listErrors := make([]FieldError, len(errs))
	for index, e := range errs {
		listErrors[index] = FieldError{Field: strings.ToLower(e.Field()), Message: messageForTag(e)}
	}
	return listErrors
}

//Validate Database: Error Handling
func messageFormCodeDB(pgErr *pgconn.PgError) (int, string) {
	switch pgErr.Code {
	//Validate DB
	case "23505":
		return 400, "This field is duplicate"
	case "23503":
		return 400, "Item doesn't exist"
	case "42P01":
		return 500, "Table doesn't exist"
	}
	return 400, pgErr.Error()
}

func fieldfromError(err error) (int, FieldError) {
	switch err.Error() {
	// Register
	case "dont find role name":
		return 400, FieldError{Field: "role", Message: err.Error()}
	case "cant hash password":
		return 500, FieldError{Field: "password", Message: err.Error()}

	// Login
	case "name isn't already exist":
		return 400, FieldError{Field: "name", Message: err.Error()}
	case "incorrect password":
		return 400, FieldError{Field: "password", Message: err.Error()}

	}
	return 0, FieldError{Field: "", Message: ""}
}

func DBError(err error) (int, []FieldError) {
	var fieldErrors []FieldError
	if err == nil {
		return 400, fieldErrors
	}
	// Error have code
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		StatusCode, MessageDB := messageFormCodeDB(pgErr)
		fieldErrors = append(fieldErrors, FieldError{Field: detail2ColumnName(pgErr.Detail), Message: MessageDB})
		return StatusCode, fieldErrors
	}
	// gorm.ErrRecordNotFound (Query "Where" don't find)
	statusCode, fieldError := fieldfromError(err)
	if statusCode > 0 { //New error
		fieldErrors = append(fieldErrors, fieldError)
		return statusCode, fieldErrors
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fieldErrors = append(fieldErrors, FieldError{Field: "", Message: "URL not found"})
		return 404, fieldErrors
	}
	// else
	fieldErrors = append(fieldErrors, FieldError{Field: "Unknown", Message: err.Error()})
	return 500, fieldErrors

}

func detail2ColumnName(str string) string {
	if str == "" {
		return str
	}
	re := regexp.MustCompile(`Key \(([^\)]+)\)=\(([^\)]+)\)`)
	match := re.FindStringSubmatch(str)
	if len(match) == 3 {
		return match[1]
	} else {
		fmt.Println(str)
		return "Unknown"
	}
}
