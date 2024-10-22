package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Product struct {
	ID                      int
	ManufacturerID          int
	Manufacturer            string
	ProductClassificationID int
	ProductClassification   string
	ProductFormID           int
	ProductForm             string
	Name                    string
	GenericName             string
	Categories              string
	Description             string
	UnitInPack              *int
	SellingUnit             string
	Weight                  decimal.Decimal
	Height                  decimal.Decimal
	Length                  decimal.Decimal
	Width                   decimal.Decimal
	ImageURL                string
	IsActive                bool
	Usage                   int
	Agg                     string
	CreatedAt               time.Time
	UpdatedAt               time.Time
	DeletedAt               *time.Time
}

type Manufacturer struct {
	ID   int
	Name string
}

type ProductClassification struct {
	ID   int
	Name string
}

type ProductForm struct {
	ID   int
	Name string
}
