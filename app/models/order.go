package models

import (
	"gorm.io/gorm"
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

func (o *Order) SetUserID(UserId uint) {
	o.UserId = UserId
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
