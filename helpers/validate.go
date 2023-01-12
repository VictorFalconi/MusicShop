package helpers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgconn"
)

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

func ValidateErrors(errs validator.ValidationErrors) interface{} {
	listErrors := make([]FieldError, len(errs))
	for index, e := range errs {
		listErrors[index] = FieldError{Field: e.Field(), Message: MessageForTag(e)}
	}
	return listErrors
}

//Validate Database: Error Handling
func MessageForTagDB(pgErr *pgconn.PgError) string {
	switch pgErr.Code {
	//Validate DB
	case "23505":
		return "This field is duplicate"
	}
	return pgErr.Error()
}

func DBError(err error) interface{} {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		//dictError := FieldError{Field: pgErr.ConstraintName, Message: MessageForTagDB(pgErr)}
		//dictError := FieldError{Field: pgErr.ColumnName, Message: MessageForTagDB(pgErr)} //Return ''
		listError := make([]FieldError, 1)
		listError[0] = FieldError{Field: pgErr.ConstraintName, Message: MessageForTagDB(pgErr)}
		return listError
	}
	return nil
}
