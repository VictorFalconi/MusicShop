package models

import (
	"gorm.io/gorm"
	"server/helpers"
	"time"
)

type Order struct {
	Id          uint    `      json:"ID"               form:"ID"                  gorm:"primary_key"`
	TotalAmount float32 `      json:"totalamount"      form:"totalamount"         gorm:"not null" validate:""`
	//PaymentMethod   string `json:"paymentmethod"    form:"paymentmethod"       gorm:"" validate:""`
	//PaymentStatus   string `json:"paymentstatus"    form:"paymentstatus"       gorm:"" validate:""`
	ShippingAddress string `   json:"shippingaddress"  form:"shippingaddress"     gorm:"" validate:""`
	//BillingAddress  string `json:"billingaddress"  form:"billingaddress"     gorm:"" validate:""`
	Status string `            json:"status"           form:"status"              gorm:"not null;default:'Pending'" validate:""` //Pending - Accept - Completed - Canceled

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	UserId uint

	Products []Product `json:"products" form:"products"   gorm:"many2many:order_products"` //Order n-n Products
}

type Orders []Order

type OrderProducts struct {
	Id       uint    `json:"ID"              form:"ID"             gorm:"primary_key"`
	Quantity int     `json:"quantity"        form:"quantity"       gorm:"not null;default:0"     validate:"required"`
	Price    float32 `json:"price"           form:"price"          gorm:"not null;default:0.0"   validate:"required"`
	Discount float32 `json:"discount"        form:"discount"       gorm:"not null;default:0.0"   validate:"required"`

	CreatedAt time.Time
	UpdatedAt time.Time
	//DeletedAt  gorm.DeletedAt `gorm:"index"`

	OrderID   uint
	ProductID uint `  json:"product_id"      form:"product_id"     gorm:""                       validate:""`
}

type InputOrder struct {
	TotalAmount     float32 `json:"totalamount"      form:"totalamount"         gorm:"not null;default:0" validate:"required"`
	ShippingAddress string  `json:"shippingaddress"  form:"shippingaddress"     gorm:""                   validate:""`
	//Status          string  `json:"status"           form:"status"              gorm:"not null"           validate:"required"`
	Products []struct {
		ProductID uint    `json:"product_id"      form:"product_id"     gorm:"not null"               validate:"required"`
		Quantity  int     `json:"quantity"        form:"quantity"       gorm:"not null;default:0"     validate:"required"`
		Price     float32 `json:"price"           form:"price"          gorm:"not null;default:0.0"   validate:"required"`
		Discount  float32 `json:"discount"        form:"discount"       gorm:"not null;default:0.0"   validate:"required"`
	} `                     json:"products"          form:"products"            gorm:"not null"           validate:"required"`
}

func (o *Order) SetUserID(UserId uint) {
	o.UserId = UserId
}

func (o *Order) MapOrder(inputOrder *InputOrder) {
	o.TotalAmount = inputOrder.TotalAmount
	//o.PaymentMethod = inputOrder.PaymentMethod
	//o.PaymentStatus = inputOrder.PaymentStatus
	o.ShippingAddress = inputOrder.ShippingAddress
	//o.BillingAddress = inputOrder.BillingAddress
	//o.Status = inputOrder.Status
}

func OutputOrder(o *Order, ops *[]OrderProducts) interface{} {
	products := make([]map[string]interface{}, 0)
	for _, orderProduct := range *ops {
		product := map[string]interface{}{
			"product_id": orderProduct.ProductID,
			"quantity":   orderProduct.Quantity,
			"price":      orderProduct.Price,
			"discount":   orderProduct.Discount,
		}
		products = append(products, product)
	}
	output := map[string]interface{}{
		"id":              o.Id,
		"totalamount":     o.TotalAmount,
		"shippingaddress": o.ShippingAddress,
		"status":          o.Status,
		"products":        products,
	}
	return output
}

func (o *Order) IsPending() bool {
	if o.Status == "Pending" {
		return true
	}
	return false
}

func (op *OrderProducts) IsStocking(db *gorm.DB) bool {
	var product Product
	db.Where("id = ?", op.ProductID).First(&product)

	QuantityOrder := op.Quantity
	QuantityProduct := product.Quantity

	if QuantityProduct >= QuantityOrder {
		return true
	}
	return false
}

func (op *OrderProducts) GetQuantity(db *gorm.DB) (int, int) {
	var product Product
	db.Where("id = ?", op.ProductID).First(&product)

	QuantityOrder := op.Quantity
	QuantityProduct := product.Quantity

	return QuantityOrder, QuantityProduct
}

func (op *OrderProducts) GetNameProduct(db *gorm.DB) string {
	var product Product
	db.Where("id = ?", op.ProductID).First(&product)
	return product.Name
}

// CRUD
// User

func (order *Order) User_Create(db *gorm.DB, input *InputOrder, userId uint) (int, interface{}) {
	//Map
	order.MapOrder(input)
	// Set UserID for Order
	order.SetUserID(userId)
	// Create
	tx := db.Begin()
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		statusCode, ErrorDB := helpers.DBError(err)
		return statusCode, ErrorDB
	} else {
		// OrderProducts
		var orderProducts []OrderProducts
		for _, product := range input.Products {
			orderProduct := OrderProducts{
				OrderID:   order.Id,
				ProductID: product.ProductID,
				Quantity:  product.Quantity,
				Price:     product.Price,
				Discount:  product.Discount,
			}
			// Check quantity of order with product
			if !orderProduct.IsStocking(db) {
				tx.Rollback()
				fError := helpers.FieldError{Field: "quantity", Message: "'" + orderProduct.GetNameProduct(db) + "' is not enough or out of stock!"}
				return 400, fError
			}
			orderProducts = append(orderProducts, orderProduct)
		}
		if errOrderProducts := tx.Create(&orderProducts).Error; err != nil {
			tx.Rollback()
			statusCode, ErrorDB := helpers.DBError(errOrderProducts)
			return statusCode, ErrorDB
		}
		tx.Commit()
		return 201, nil
	}
}

func (order *Order) User_Read(db *gorm.DB, orderId string, userId uint) (int, interface{}, interface{}) {
	if err := db.Where("id = ? AND user_id = ? ", orderId, userId).First(&order).Error; err != nil {
		return 404, helpers.FieldError{Field: "", Message: "URL not found"}, nil
	} else {
		//OrderProduct -> OrderID
		var orderProducts []OrderProducts
		if errOrderProducts := db.Where("order_id = ? ", order.Id).Find(&orderProducts).Error; errOrderProducts != nil {
			return 404, helpers.FieldError{Field: "", Message: "URL not found"}, nil
		}
		// order + orderProducts
		output := OutputOrder(order, &orderProducts)
		return 200, nil, output
	}
}

func (orders *Orders) User_ReadsOfUser(db *gorm.DB, userId uint) (int, interface{}, interface{}) {
	if err := db.Where("user_id = ?", userId).Find(&orders).Error; err != nil {
		return 404, helpers.FieldError{Field: "", Message: "URL not found"}, nil
	} else {
		var outputs []interface{}
		for _, order := range *orders {
			//OrderProduct -> OrderID
			var orderProducts []OrderProducts
			if errOrderProducts := db.Where("order_id = ? ", order.Id).Find(&orderProducts).Error; errOrderProducts != nil {
				return 404, helpers.FieldError{Field: "", Message: "URL not found"}, nil
			}
			// order + orderProducts
			output := OutputOrder(&order, &orderProducts)
			outputs = append(outputs, output)
		}
		return 200, nil, outputs
	}
}

func (order *Order) User_Cancel(db *gorm.DB) (int, interface{}) {
	// Pending -> Canceled
	if order.IsPending() {
		order.Status = "Canceled"
		if errUpdate := db.Save(&order).Error; errUpdate != nil {
			statusCode, ErrorDB := helpers.DBError(errUpdate)
			return statusCode, ErrorDB
		} else {
			return 200, nil
		}
	} else {
		return 400, helpers.FieldError{Field: "status", Message: "Cant cancel this order"}
	}
}

//Admin

func (orders *Orders) Admin_Reads(db *gorm.DB) (int, interface{}, interface{}) {
	if err := db.Find(&orders).Error; err != nil {
		return 404, helpers.FieldError{Field: "", Message: "URL not found"}, nil
	} else {
		var outputs []interface{}
		for _, order := range *orders {
			//OrderProduct -> OrderID
			var orderProducts []OrderProducts
			if errOrderProducts := db.Where("order_id = ? ", order.Id).Find(&orderProducts).Error; errOrderProducts != nil {
				statusCode, ErrorDB := helpers.DBError(errOrderProducts)
				return statusCode, ErrorDB, nil
			}
			// order + orderProducts
			output := OutputOrder(&order, &orderProducts)
			outputs = append(outputs, output)
		}
		return 201, nil, outputs
	}
}

func (order *Order) Admin_Read(db *gorm.DB, orderId string) (int, interface{}, interface{}) {
	if err := db.Where("id = ?", orderId).First(&order).Error; err != nil {
		return 404, helpers.FieldError{Field: "", Message: "URL not found"}, nil
	} else {
		//OrderProduct -> OrderID
		var orderProducts []OrderProducts
		if errOrderProducts := db.Where("order_id = ? ", order.Id).Find(&orderProducts).Error; errOrderProducts != nil {
			return 404, helpers.FieldError{Field: "", Message: "URL not found"}, nil
		}
		// order + orderProducts
		output := OutputOrder(order, &orderProducts)
		return 200, nil, output
	}
}

func (order *Order) Admin_AcceptOrder(db *gorm.DB) (int, interface{}) {
	// Pending -> Accept (Quantity > 0)
	if order.IsPending() {
		//OrderProduct -> OrderID
		var orderProducts []OrderProducts
		if errOrderProducts := db.Where("order_id = ? ", order.Id).Find(&orderProducts).Error; errOrderProducts != nil {
			statusCode, ErrorDB := helpers.DBError(errOrderProducts)
			return statusCode, ErrorDB
		}
		// Check Quantity
		tx := db.Begin()
		for _, op := range orderProducts {
			if !op.IsStocking(db) {
				return 400, helpers.FieldError{Field: "quantity", Message: "Quantity of '" + op.GetNameProduct(db) + "' is not enough"}
			} else {
				// (product - order) Quantity
				var product Product
				db.Where("id = ?", op.ProductID).First(&product)
				product.Quantity = product.Quantity - op.Quantity
				if errProduct := tx.Save(&product).Error; errProduct != nil {
					tx.Rollback()
					statusCode, ErrorDB := helpers.DBError(errProduct)
					return statusCode, ErrorDB
				}
			}
		}
		order.Status = "Accept"
		if errUpdate := tx.Save(&order).Error; errUpdate != nil {
			tx.Rollback()
			statusCode, ErrorDB := helpers.DBError(errUpdate)
			return statusCode, ErrorDB
		} else {
			tx.Commit()
			return 200, nil
		}
	} else {
		return 400, helpers.FieldError{Field: "status", Message: "Cant accept this order"}
	}
}

func (order *Order) Admin_CancelOrder(db *gorm.DB) (int, interface{}) {
	// all status -> Canceled
	order.Status = "Canceled"
	if errUpdate := db.Save(&order).Error; errUpdate != nil {
		statusCode, ErrorDB := helpers.DBError(errUpdate)
		return statusCode, ErrorDB
	} else {
		return 200, nil
	}
}
