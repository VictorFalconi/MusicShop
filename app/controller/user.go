package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"server/app/model"
	"server/app/service"
	"server/config"
	"server/helpers"
)

type UserController struct {
	service service.UserServiceInterface
}

func NewUserController(service service.UserServiceInterface) *UserController {
	return &UserController{service}
}

func (uc *UserController) RegisterHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user model.User
		// Check data type
		if err := helpers.DataContentType(ctx, &user); err != nil {
			helpers.RespondJSON(ctx, 400, helpers.StatusCodeFromInt(400), err.Error(), nil)
			return
		}
		// Check validate field
		if err := validator.New().Struct(&user); err != nil {
			listErrors := helpers.ValidateErrors(err.(validator.ValidationErrors))
			helpers.RespondJSON(ctx, 400, helpers.StatusCodeFromInt(400), listErrors, nil)
			return
		}
		// Create
		if err := uc.service.Register(&user); err != nil {
			statusCode, message := helpers.DBError(err)
			helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), message, nil)
			return
		}
		helpers.RespondJSON(ctx, 201, helpers.StatusCodeFromInt(201), nil, nil)
		return
	}
}

//Login
func (uc *UserController) LoginHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var loginUser model.LoginUser
		// Check data type
		if err := helpers.DataContentType(ctx, &loginUser); err != nil {
			helpers.RespondJSON(ctx, 400, helpers.StatusCodeFromInt(400), err.Error(), nil)
			return
		}
		// Check validate field
		if err := validator.New().Struct(&loginUser); err != nil {
			listErrors := helpers.ValidateErrors(err.(validator.ValidationErrors))
			helpers.RespondJSON(ctx, 400, helpers.StatusCodeFromInt(400), listErrors, nil)
			return
		}
		// Login
		err, token := uc.service.Login(&loginUser)
		if err != nil {
			statusCode, message := helpers.DBError(err)
			helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), message, nil)
			return
		}
		//Response token
		ctx.SetSameSite(http.SameSiteLaxMode)
		ctx.SetCookie("Authorization", token, 3600*12, "/", "", false, false)
		//mapToken := map[string]string{"Authorization": token}
		helpers.RespondJSON(ctx, 201, helpers.StatusCodeFromInt(201), nil, nil)
		return
	}
}

//OAuth2

func (uc *UserController) OAuth2Home() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.WriteHeader(http.StatusOK)
		fmt.Fprintf(ctx.Writer, `<html><body><a href="/auth/login">Google LogIn</a></body></html>`)
	}
}

func (uc *UserController) OAuth2LoginHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		url := config.GoogleOauthConfig.AuthCodeURL("state")
		ctx.Redirect(http.StatusTemporaryRedirect, url)
	}
}

func (uc *UserController) OAuth2CallbackHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get token
		code := ctx.Query("code")
		token, err := config.GoogleOauthConfig.Exchange(context.Background(), code)
		if err != nil {
			fieldErr := helpers.FieldError{Field: "OAuth2", Message: "Failed to exchange token"}
			helpers.RespondJSON(ctx, 500, helpers.StatusCodeFromInt(500), fieldErr, nil)
			return
		}
		// Get userinfo
		client := config.GoogleOauthConfig.Client(context.Background(), token)
		resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
		if err != nil {
			fieldErr := helpers.FieldError{Field: "OAuth2", Message: "Failed to get user info"}
			helpers.RespondJSON(ctx, 500, helpers.StatusCodeFromInt(500), fieldErr, nil)
			return
		}
		defer resp.Body.Close()
		// Decode userinfo
		userInfo := make(map[string]interface{})
		if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
			fieldErr := helpers.FieldError{Field: "OAuth2", Message: "Failed to decode user info"}
			helpers.RespondJSON(ctx, 500, helpers.StatusCodeFromInt(500), fieldErr, nil)
			return
		}
		// Login || Register
		err, jwtToken := uc.service.GoogleLogin(userInfo)
		if err != nil {
			statusCode, message := helpers.DBError(err)
			helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), message, nil)
			return
		}
		//Response token
		ctx.SetSameSite(http.SameSiteLaxMode)
		ctx.SetCookie("Authorization", jwtToken, 3600*12, "/", "", false, false)
		helpers.RespondJSON(ctx, 201, helpers.StatusCodeFromInt(201), nil, nil)
		return
	}
}

// Read
func (uc *UserController) ReadUserHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(model.User)
		helpers.RespondJSON(ctx, 200, helpers.StatusCodeFromInt(200), nil, user.ReadUser())
		return
	}
}

// Update
func (uc *UserController) UpdateUserHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Current User
		currUser := ctx.MustGet("user").(model.User)
		// Get request
		var newUser model.User
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
		if err := uc.service.Update(&currUser, &newUser); err != nil {
			statusCode, message := helpers.DBError(err)
			helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), message, nil)
			return
		}
		helpers.RespondJSON(ctx, 200, helpers.StatusCodeFromInt(200), nil, nil)
		return
	}
}
