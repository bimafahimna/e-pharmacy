package dto

import (
	"github.com/shopspring/decimal"
)

type CartContent struct {
	PharmacyName string     `json:"pharmacy_name"`
	Items        []CartItem `json:"items"`
}

type CartItem struct {
	PharmacyID  int             `json:"pharmacy_id"`
	ProductID   int             `json:"product_id"`
	Name        string          `json:"name"`
	ImageURL    string          `json:"image_url"`
	Price       decimal.Decimal `json:"price"`
	Stock       int             `json:"stock"`
	Quantity    int             `json:"quantity"`
	SellingUnit string          `json:"selling_unit"`
	Description string          `json:"description"`
	Weight      string          `json:"weight"`
}

type CartUri struct {
	PharmacyID int `uri:"pharmacy_id" binding:"required"`
	ProductID  int `uri:"product_id" binding:"required"`
}

type AddCartItemRequest struct {
	PharmacyID int `json:"pharmacy_id" binding:"required"`
	ProductID  int `json:"product_id" binding:"required"`
	Quantity   int `json:"quantity" binding:"required,gte=1"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,numeric"`
}
