package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"server/app/models"
	"server/config"
	"server/helpers"
	"server/middleware"
	"strings"
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
	// Set role for user
	if err := user.SetUserRole(config.DB, "user"); err != nil {
		statusCode, ErrorDB := helpers.DBError(err)
		helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), ErrorDB, nil)
		return
	}
	// Hash password
	if err := user.HashPassword(); err != nil {
		helpers.RespondJSON(ctx, 500, helpers.StatusCodeFromInt(500), "Cant Hash Password", nil)
		return
	}
	// Create new User (Check validate Database)
	if err := config.DB.Create(&user).Error; err != nil {
		statusCode, ErrorDB := helpers.DBError(err) // Hiện 3 lỗi khi nhập trùng cả 3 fields
		helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), ErrorDB, nil)
		return
	} else {
		helpers.RespondJSON(ctx, 201, helpers.StatusCodeFromInt(201), nil, nil)
		return
	}
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
	// Check Field "name" in db
	user := &models.User{}
	if err := config.DB.Where("name = ?", currUser.Name).First(&user).Error; err != nil {
		var fieldErrors []helpers.FieldError
		fieldError := helpers.FieldError{Field: "name", Message: "Name isn't already exist"}
		fieldErrors = append(fieldErrors, fieldError)
		helpers.RespondJSON(ctx, 400, helpers.StatusCodeFromInt(400), fieldErrors, nil)
		return
	} else {
		// Compare password
		if user.ComparePassword(currUser.Password) == false {
			var fieldErrors []helpers.FieldError
			fieldError := helpers.FieldError{Field: "password", Message: "Incorrect Password"}
			fieldErrors = append(fieldErrors, fieldError)
			helpers.RespondJSON(ctx, 400, helpers.StatusCodeFromInt(400), fieldErrors, nil)
			return
		} else {
			//Create token
			token, errCreate := middleware.CreateToken(user.Id)
			if errCreate != nil {
				helpers.RespondJSON(ctx, 500, helpers.StatusCodeFromInt(500), errCreate.Error(), nil)
				return
			}
			//Repose token
			ctx.SetSameSite(http.SameSiteLaxMode)
			ctx.SetCookie("Authorization", token, 3600*12, "", "", false, true)
			helpers.RespondJSON(ctx, 201, helpers.StatusCodeFromInt(201), nil, nil)
			return
		}
	}
}

// Update user
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
	// Update & Hash password
	currUser.UpdateStruct(&newUser)
	if err := currUser.HashPassword(); err != nil {
		helpers.RespondJSON(ctx, 500, helpers.StatusCodeFromInt(500), "Cant Hash Password", nil)
		return
	}
	if err := config.DB.Save(&currUser).Error; err != nil {
		statusCode, ErrorDB := helpers.DBError(err)
		helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), ErrorDB, nil)
		return
	} else {
		helpers.RespondJSON(ctx, 200, helpers.StatusCodeFromInt(200), nil, nil)
		return
	}
}

// Compare Token
func AuthenticToken(ctx *gin.Context) {
	//Get token
	bearerToken := ctx.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		//Validate token
		token, err := middleware.ValidateToken(strings.Split(bearerToken, " ")[1])
		if token.Valid && err == nil {
			claims := token.Claims.(jwt.MapClaims)
			// Query User from id
			var user models.User
			if errDB := config.DB.Where("id = ?", claims["id"]).First(&user).Error; errDB == nil {
				helpers.RespondJSON(ctx, 200, helpers.StatusCodeFromInt(200), nil, &user)
				return
			} else {
				helpers.RespondJSON(ctx, 400, helpers.StatusCodeFromInt(400), errDB.Error(), &user)
				return
			}
		} else {
			helpers.RespondJSON(ctx, 401, helpers.StatusCodeFromInt(401), err.Error(), nil)
			return
		}
	} else {
		helpers.RespondJSON(ctx, 400, helpers.StatusCodeFromInt(400), "No token found", nil)
		return
	}
}
