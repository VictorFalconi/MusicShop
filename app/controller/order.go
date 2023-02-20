package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"server/app/models"
	"server/config"
	"server/helpers"
)

func User_CreateOrder(ctx *gin.Context) {
	var input models.InputOrder
	// Check data type
	if err := helpers.DataContentType(ctx, &input); err != nil {
		helpers.RespondJSON(ctx, 400, helpers.StatusCodeFromInt(400), err.Error(), nil)
		return
	}
	// Check validate field
	if err := validator.New().Struct(&input); err != nil {
		listErrors := helpers.ValidateErrors(err.(validator.ValidationErrors))
		helpers.RespondJSON(ctx, 400, helpers.StatusCodeFromInt(400), listErrors, nil)
		return
	}
	// Order
	var order models.Order
	currUser := ctx.MustGet("user").(models.User)
	// Create
	statusCode, Message := order.User_Create(config.DB, &input, currUser.Id)
	helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), Message, nil)
	return

}

func User_ReadOrder(ctx *gin.Context) {
	currUser := ctx.MustGet("user").(models.User)
	var order models.Order
	// Read
	statusCode, Message, output := order.User_Read(config.DB, ctx.Param("id"), currUser.Id)
	helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), Message, output)
	return

}

func User_ReadOrders(ctx *gin.Context) {
	currUser := ctx.MustGet("user").(models.User)
	var orders models.Orders
	//Reads
	statusCode, Message, output := orders.User_ReadsOfUser(config.DB, currUser.Id)
	helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), Message, output)

}

func User_CancelOrder(ctx *gin.Context) {
	currUser := ctx.MustGet("user").(models.User)
	var order models.Order
	// Find
	statusCode, Message, output := order.User_Read(config.DB, ctx.Param("id"), currUser.Id)
	if statusCode != 200 {
		helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), Message, output)
		return
	}
	// Cancel
	statusCode, Message = order.User_Cancel(config.DB)
	helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), Message, nil)
	return
}

// Admin

func Admin_ReadOrders(ctx *gin.Context) {
	var orders models.Orders
	statusCode, Message, output := orders.Admin_Reads(config.DB)
	helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), Message, output)
	return
}

func Admin_ReadOrder(ctx *gin.Context) {
	var order models.Order
	statusCode, Message, output := order.Admin_Read(config.DB, ctx.Param("id"))
	helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), Message, output)
	return
}

func Admin_AcceptOrder(ctx *gin.Context) {
	var order models.Order
	// Find
	statusCode, Message, _ := order.Admin_Read(config.DB, ctx.Param("id"))
	if statusCode != 200 {
		helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), Message, nil)
		return
	}
	// Update
	statusCode, Message = order.Admin_AcceptOrder(config.DB)
	helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), Message, nil)
	return
}

func Admin_CancelOrder(ctx *gin.Context) {
	var order models.Order
	// Find
	statusCode, Message, _ := order.Admin_Read(config.DB, ctx.Param("id"))
	if statusCode != 200 {
		helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), Message, nil)
		return
	}
	// Update
	statusCode, Message = order.Admin_CancelOrder(config.DB)
	helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), Message, nil)
	return
}
