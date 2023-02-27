package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"server/app/model"
	"server/app/service"
	"server/helpers"
)

type UserController struct {
	service service.UserServiceInterface
}

func NewUserController(service service.UserServiceInterface) *UserController {
	return &UserController{service}
}

func (c *UserController) RegisterHandler() gin.HandlerFunc {
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
		if err := c.service.Register(&user); err != nil {
			statusCode, message := helpers.DBError(err)
			helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), message, nil)
			return
		}
		helpers.RespondJSON(ctx, 201, helpers.StatusCodeFromInt(201), nil, nil)
		return
	}
}

//Login
func (c *UserController) LoginHandler() gin.HandlerFunc {
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
		err, token := c.service.Login(&loginUser)
		if err != nil {
			statusCode, message := helpers.DBError(err)
			helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), message, nil)
			return
		}
		//Response token
		ctx.SetSameSite(http.SameSiteLaxMode)
		ctx.SetCookie("Authorization", token, 3600*12, "/", "", false, false)
		//mapToken := map[string]string{"Authorization": token}
		helpers.RespondJSON(ctx, 201, helpers.StatusCodeFromInt(201), nil, nil) //mapToken
		return
	}
}

// Read
func (c *UserController) ReadUserHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(model.User)
		helpers.RespondJSON(ctx, 200, helpers.StatusCodeFromInt(200), nil, user.ReadUser())
		return
	}
}

// Update
func (c *UserController) UpdateUserHandler() gin.HandlerFunc {
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
		if err := c.service.Update(&currUser, &newUser); err != nil {
			statusCode, message := helpers.DBError(err)
			helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), message, nil)
			return
		}
		helpers.RespondJSON(ctx, 200, helpers.StatusCodeFromInt(200), nil, nil)
		return
	}
}
