package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"server/app/models"
	"server/config"
	"server/helpers"
)

func CreateOrder(ctx *gin.Context) {
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
	order.MapOrder(&input)
	// Set UserID for Order
	currUser := ctx.MustGet("user").(models.User)
	order.SetUserID(currUser.Id)
	// Create (Check validate Database)
	tx := config.DB.Begin()
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		statusCode, ErrorDB := helpers.DBError(err)
		helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), ErrorDB, nil)
		return
	} else {
		// OrderProducts
		var orderProducts []models.OrderProducts
		for _, product := range input.Products {
			orderProduct := models.OrderProducts{
				OrderID:   order.Id,
				ProductID: product.ProductID,
				Quantity:  product.Quantity,
				Price:     product.Price,
				Discount:  product.Discount,
			}
			// Check quantity of order with product
			if !orderProduct.IsStocking(config.DB) {
				tx.Rollback()
				fError := helpers.FieldError{Field: "quantity", Message: "'" + orderProduct.GetNameProduct(config.DB) + "' is not enough or out of stock!"}
				helpers.RespondJSON(ctx, 400, helpers.StatusCodeFromInt(400), fError, nil)
				return
			}
			orderProducts = append(orderProducts, orderProduct)
		}
		if errOrderProducts := tx.Create(&orderProducts).Error; err != nil {
			tx.Rollback()
			statusCode, ErrorDB := helpers.DBError(errOrderProducts)
			helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), ErrorDB, nil)
			return
		}
		tx.Commit()
		helpers.RespondJSON(ctx, 201, helpers.StatusCodeFromInt(201), nil, nil)
		return
	}
}

func ReadOrder(ctx *gin.Context) {
	currUser := ctx.MustGet("user").(models.User)
	var order models.Order
	if err := config.DB.Where("id = ? AND user_id = ? ", ctx.Param("id"), currUser.Id).First(&order).Error; err != nil {
		helpers.RespondJSON(ctx, 404, helpers.StatusCodeFromInt(404), "URL not found", nil)
		return
	} else {
		//OrderProduct -> OrderID
		var orderProducts []models.OrderProducts
		if errorderProducts := config.DB.Where("order_id = ? ", order.Id).Find(&orderProducts).Error; errorderProducts != nil {
			statusCode, ErrorDB := helpers.DBError(errorderProducts)
			helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), ErrorDB, nil)
			return
		}
		// order + orderProducts
		output := models.OutputOrder(&order, &orderProducts)
		helpers.RespondJSON(ctx, 200, helpers.StatusCodeFromInt(200), nil, output)
		return
	}
}

func ReadOrdersOfUser(ctx *gin.Context) {
	currUser := ctx.MustGet("user").(models.User)
	var orders []models.Order
	if err := config.DB.Where("user_id = ?", currUser.Id).Find(&orders).Error; err != nil {
		helpers.RespondJSON(ctx, 404, helpers.StatusCodeFromInt(404), "URL not found", nil)
		return
	} else {
		var outputs []interface{}
		for _, order := range orders {
			//OrderProduct -> OrderID
			var orderProducts []models.OrderProducts
			if errorderProducts := config.DB.Where("order_id = ? ", order.Id).Find(&orderProducts).Error; errorderProducts != nil {
				statusCode, ErrorDB := helpers.DBError(errorderProducts)
				helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), ErrorDB, nil)
				return
			}
			// order + orderProducts
			output := models.OutputOrder(&order, &orderProducts)
			outputs = append(outputs, output)
		}
		helpers.RespondJSON(ctx, 200, helpers.StatusCodeFromInt(200), nil, outputs)
		return
	}
}

func UserCancelOrder(ctx *gin.Context) {
	currUser := ctx.MustGet("user").(models.User)
	var order models.Order
	if err := config.DB.Where("id = ? AND user_id = ? ", ctx.Param("id"), currUser.Id).First(&order).Error; err != nil {
		helpers.RespondJSON(ctx, 404, helpers.StatusCodeFromInt(404), "URL not found", nil)
		return
	} else {
		// Pending -> Canceled
		if order.IsPending() {
			order.Status = "Canceled"
			if errUpdate := config.DB.Save(&order).Error; errUpdate != nil {
				statusCode, ErrorDB := helpers.DBError(errUpdate)
				helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), ErrorDB, nil)
				return
			} else {
				helpers.RespondJSON(ctx, 200, helpers.StatusCodeFromInt(200), nil, nil)
				return
			}
		} else {
			fError := helpers.FieldError{Field: "status", Message: "Cant cancel this order"}
			helpers.RespondJSON(ctx, 400, helpers.StatusCodeFromInt(400), fError, nil)
			return
		}
	}
}

// Admin function

func ReadOrders(ctx *gin.Context) {
	var orders []models.Order
	if err := config.DB.Find(&orders).Error; err != nil {
		statusCode, ErrorDB := helpers.DBError(err)
		helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), ErrorDB, nil)
		return
	} else {
		var outputs []interface{}
		for _, order := range orders {
			//OrderProduct -> OrderID
			var orderProducts []models.OrderProducts
			if errorderProducts := config.DB.Where("order_id = ? ", order.Id).Find(&orderProducts).Error; errorderProducts != nil {
				statusCode, ErrorDB := helpers.DBError(errorderProducts)
				helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), ErrorDB, nil)
				return
			}
			// order + orderProducts
			output := models.OutputOrder(&order, &orderProducts)
			outputs = append(outputs, output)
		}
		helpers.RespondJSON(ctx, 200, helpers.StatusCodeFromInt(200), nil, &outputs)
		return
	}
}

func AcceptOrder(ctx *gin.Context) {
	var order models.Order
	if err := config.DB.Where("id = ?", ctx.Param("id")).First(&order).Error; err != nil {
		helpers.RespondJSON(ctx, 404, helpers.StatusCodeFromInt(404), "URL not found", nil)
		return
	} else {
		// Pending -> Accept (Quantity > 0)
		if order.IsPending() {
			//OrderProduct -> OrderID
			var orderProducts []models.OrderProducts
			if errorderProducts := config.DB.Where("order_id = ? ", order.Id).Find(&orderProducts).Error; errorderProducts != nil {
				statusCode, ErrorDB := helpers.DBError(errorderProducts)
				helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), ErrorDB, nil)
				return
			}
			// Check Quantity
			tx := config.DB.Begin()
			for _, op := range orderProducts {
				if !op.IsStocking(config.DB) {
					fError := helpers.FieldError{Field: "quantity", Message: "Quantity of '" + op.GetNameProduct(config.DB) + "' is not enough"}
					helpers.RespondJSON(ctx, 400, helpers.StatusCodeFromInt(400), fError, nil)
					return
				} else {
					// (product - order) Quantity
					var product models.Product
					config.DB.Where("id = ?", op.ProductID).First(&product)
					product.Quantity = product.Quantity - op.Quantity
					if errProduct := tx.Save(&product).Error; err != nil {
						tx.Rollback()
						statusCode, ErrorDB := helpers.DBError(errProduct)
						helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), ErrorDB, nil)
						return
					}
				}
			}
			order.Status = "Accept"
			if errUpdate := tx.Save(&order).Error; errUpdate != nil {
				tx.Rollback()
				statusCode, ErrorDB := helpers.DBError(errUpdate)
				helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), ErrorDB, nil)
				return
			} else {
				tx.Commit()
				helpers.RespondJSON(ctx, 200, helpers.StatusCodeFromInt(200), nil, nil)
				return
			}
		} else {
			fError := helpers.FieldError{Field: "status", Message: "Cant accept this order"}
			helpers.RespondJSON(ctx, 400, helpers.StatusCodeFromInt(400), fError, nil)
			return
		}
	}
}

func CancelOrder(ctx *gin.Context) {
	var order models.Order
	if err := config.DB.Where("id = ?", ctx.Param("id")).First(&order).Error; err != nil {
		helpers.RespondJSON(ctx, 404, helpers.StatusCodeFromInt(404), "URL not found", nil)
		return
	} else {
		// all status -> Canceled
		order.Status = "Canceled"
		if errUpdate := config.DB.Save(&order).Error; errUpdate != nil {
			statusCode, ErrorDB := helpers.DBError(errUpdate)
			helpers.RespondJSON(ctx, statusCode, helpers.StatusCodeFromInt(statusCode), ErrorDB, nil)
			return
		} else {
			helpers.RespondJSON(ctx, 200, helpers.StatusCodeFromInt(200), nil, nil)
			return
		}
	}
}
