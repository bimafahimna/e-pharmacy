package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type UnpaidOrder struct {
	PaymentID     int             `json:"payment_id"`
	Orders        []Order         `json:"orders"`
	PaymentAmount decimal.Decimal `json:"payment_amount"`
}

type Order struct {
	Status          string      `json:"status,omitempty" binding:"omitempty"`
	Address         string      `json:"address,omitempty" binding:"required"`
	ContactName     string      `json:"contact_name,omitempty" binding:"required"`
	ContactPhone    string      `json:"contact_phone,omitempty" binding:"required"`
	PharmacyID      int         `json:"pharmacy_id,omitempty" binding:"required"`
	PharmacyName    string      `json:"pharmacy_name,omitempty" binding:"required"`
	LogisticID      int         `json:"logistic_id,omitempty" binding:"required"`
	LogisticName    string      `json:"logistic_name,omitempty" binding:"omitempty"`
	LogisticService string      `json:"logistic_service,omitempty" binding:"omitempty"`
	LogisticCost    string      `json:"logistic_cost,omitempty" binding:"required"`
	OrderID         int         `json:"order_id" binding:"required"`
	OrderItems      []OrderItem `json:"order_items" binding:"required"`
	OrderAmount     string      `json:"order_amount" binding:"required,numeric"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

type OrderItem struct {
	PharmacyID int    `json:"pharmacy_id" binding:"required"`
	ProductID  int    `json:"product_id" binding:"required"`
	ImageURL   string `json:"image_url" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required"`
	Price      string `json:"price" binding:"required,numeric"`
}

type PaymentUri struct {
	PaymentID int `uri:"payment_id" binding:"required"`
}

type ListOrderQuery struct {
	Status string `form:"status" binding:"required,oneof='Processed' 'Sent' 'Order Confirmed' 'Canceled'"`
}

type ListPharmacyOrderQuery struct {
	Status string `form:"status" binding:"required,oneof='Waiting for payment' 'Processed' 'Sent' 'Order Confirmed' 'Canceled'"`
	Limit  int    `form:"limit" binding:"omitempty,gte=1"`
	Page   int    `form:"page" binding:"omitempty,gte=1"`
	SortBy string `form:"sort_by" binding:"omitempty,oneof=contact_name created_at updated_at"`
	Sort   string `form:"sort" binding:"omitempty,oneof=asc desc"`
}

type CreateOrderRequest struct {
	Orders        []Order `json:"orders" binding:"required"`
	PaymentMethod string  `json:"payment_method" binding:"required,oneof='Manual transfer'"`
	PaymentAmount string  `json:"payment_amount" binding:"required,numeric"`
}

func (p *ListPharmacyOrderQuery) EnsureDefaults() {
	if p.Limit == 0 {
		p.Limit = 15
	}
	if p.Page == 0 {
		p.Page = 1
	}
	if p.SortBy == "" {
		p.SortBy = "created_at"
	}
	if p.Sort == "" {
		p.Sort = "desc"
	}
}
