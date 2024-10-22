package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type PharmacyProduct struct {
	PharmacyID            int
	ProductID             int
	Name                  string
	GenericName           string
	Manufacturer          string
	ProductClassification string
	ProductForm           string
	SoldAmount            int
	Price                 decimal.Decimal
	Stock                 int
	IsActive              bool
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             *time.Time
}

type Item struct {
	PharmacyID    int
	ProductID     int
	ProductFormID int
	Name          string
	Price         decimal.Decimal
	Stock         int
	SellingUnit   string
	ImageURL      string
}

type AvailablePharmacy struct {
	PharmacyID   int
	PharmacyName string
	PartnerName  string
	PartnerLogo  string
	Address      string
	CityName     string
	ProductID    int
	Price        string
	Stock        int
}
