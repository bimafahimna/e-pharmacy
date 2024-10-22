package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Payment struct {
	ID            int
	PaymentMethod string
	ImageURL      *string
	Amount        decimal.Decimal
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
