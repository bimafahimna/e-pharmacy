package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Order struct {
	ID              int
	UserID          int64
	PharmacyID      int
	PaymentID       int
	Status          string
	Address         string
	PharmacyName    string
	ContactName     string
	ContactPhone    string
	LogisticID      int
	LogisticName    string
	LogisticService string
	LogisticCost    decimal.Decimal
	Items           []OrderItem
	Amount          decimal.Decimal
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type OrderItem struct {
	ID         int
	OrderID    int
	PharmacyID int
	ProductID  int
	ImageURL   string
	Name       string
	Quantity   int
	Price      decimal.Decimal
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type UnpaidOrder struct {
	ID              int
	PharmacyID      int
	ProductID       int
	PaymentID       int
	Address         string
	PharmacyName    string
	ImageURL        string
	Name            string
	Quantity        int
	Price           decimal.Decimal
	ContactName     string
	ContactPhone    string
	LogisticName    string
	LogisticService string
	LogisticCost    decimal.Decimal
	OrderAmount     decimal.Decimal
	PaymentAmount   decimal.Decimal
}
