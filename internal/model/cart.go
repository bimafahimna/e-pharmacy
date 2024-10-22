package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Cart struct {
	UserID     int64
	PharmacyID int
	ProductID  int
	Quantity   int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type CartItem struct {
	UserID       int64
	PharmacyID   int
	ProductID    int
	PharmacyName string
	ProductName  string
	ImageURL     string
	Description  string
	SellingUnit  string
	Price        decimal.Decimal
	Stock        int
	Quantity     int
	Weight       decimal.Decimal
}
