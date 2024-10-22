package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type CustomerAddress struct {
	ID                  int64
	UserID              int64
	Name                string
	ReceiverName        string
	ReceiverPhoneNumber string
	Latitude            decimal.Decimal
	Longitude           decimal.Decimal
	Province            string
	CityID              int64
	City                string
	District            string
	SubDistrict         string
	AddressDetails      string
	IsActive            bool
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
