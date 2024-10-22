package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Pharmacy struct {
	ID             int
	PharmacistId   *int
	PharmacistName *string
	PartnerId      int
	PartnerName    string
	Name           string
	Address        string
	CityId         int
	CityName       *string
	Latitude       decimal.Decimal
	Longitude      decimal.Decimal
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}

type Logistic struct {
	ID                 int64
	Name               string
	LogoUrl            string
	Service            string
	PricePerKilometers decimal.NullDecimal
	EDA                int
}

type PharmacyLogistic struct {
	Pharmacies     Pharmacy
	Logistics      Logistic
	DistanceKM     decimal.Decimal
	CustomerCityID int
	PharmacyCityID int
}
