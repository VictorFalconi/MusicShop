package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"server/app/models"
	"server/config"
	"server/helpers"
	"server/middleware"
)

// Register new User
func Register(ctx *gin.Context) {
	var user models.User
	// Check data type
	if err := helpers.DataContentType(ctx, &user); err != nil {
		helpers.RespondJSON(ctx, 400, helpers.StatusCodeFromInt(400), err.Error(), nil)
		return
	}
	// Check validate field     //Thiếu 2 cái: nhập dư field k có trong user; nhập trùng 2 field giống nhau
	if err := validator.New().Struct(&user); err != nil {
		listErrors := helpers.ValidateErrors(err.(validator.ValidationErrors))
		helpers.RespondJSON(ctx, 400, helpers.StatusCodeFromInt(400), listErrors, nil)
		return
	}
	// Create
	statusCode, Message := user.Register(config.DB, "user")
	helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), Message, nil)
	return
}

// Login user
func Login(ctx *gin.Context) {
	var currUser models.LoginUser
	if err := helpers.DataContentType(ctx, &currUser); err != nil {
		helpers.RespondJSON(ctx, 400, helpers.StatusCodeFromInt(400), err.Error(), nil)
		return
	}
	validate := validator.New()
	if err := validate.Struct(&currUser); err != nil {
		dictErrors := helpers.ValidateErrors(err.(validator.ValidationErrors))
		helpers.RespondJSON(ctx, 400, helpers.StatusCodeFromInt(400), dictErrors, nil)
		return
	}
	// Login
	statusCode, Message, UserID := currUser.Login(config.DB)
	if statusCode != 201 {
		helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), Message, nil)
		return
	} else {
		//Create token
		token, errCreate := middleware.CreateToken(*UserID)
		if errCreate != nil {
			fError := helpers.FieldError{Field: "token", Message: errCreate.Error()}
			helpers.RespondJSON(ctx, 500, helpers.StatusCodeFromInt(500), fError, nil)
			return
		}
		//Response token
		ctx.SetSameSite(http.SameSiteLaxMode)
		ctx.SetCookie("Authorization", token, 3600*12, "/", "", false, true)
		helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), Message, nil)
		return
	}
}

func ReadUser(ctx *gin.Context) {
	// Get User from Token
	user := ctx.MustGet("user").(models.User)
	helpers.RespondJSON(ctx, 200, helpers.StatusCodeFromInt(200), nil, user.ReadUser())
	return
}

func UpdateUser(ctx *gin.Context) {
	// Current User
	currUser := ctx.MustGet("user").(models.User)
	// Get request
	var newUser models.User
	if err := helpers.DataContentType(ctx, &newUser); err != nil {
		helpers.RespondJSON(ctx, 400, helpers.StatusCodeFromInt(400), err.Error(), nil)
		return
	}
	if err := validator.New().Struct(&newUser); err != nil {
		listErrors := helpers.ValidateErrors(err.(validator.ValidationErrors))
		helpers.RespondJSON(ctx, 400, helpers.StatusCodeFromInt(400), listErrors, nil)
		return
	}
	// Update
	statusCode, Message := currUser.Update(config.DB, &newUser)
	helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), Message, nil)
}
