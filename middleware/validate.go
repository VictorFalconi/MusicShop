package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"server/helpers"
)

func CheckTypeRequest(ctx *gin.Context, entity interface{}) {
	if err := helpers.DataContentType(ctx, &entity); err != nil {
		helpers.RespondJSON(ctx, 400, "Error data type!", err.Error(), nil)
		return
	}
}

func CheckValidateFields(ctx *gin.Context, entity interface{}) {
	if err := validator.New().Struct(&entity); err != nil {
		listErrors := helpers.ValidateErrors(err.(validator.ValidationErrors))
		helpers.RespondJSON(ctx, 400, "Errors validate!", listErrors, nil)
		return
	}
}
