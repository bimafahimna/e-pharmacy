package dto

import (
	"strconv"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/model"
	"github.com/shopspring/decimal"
)

type Product struct {
	ID                    int    `json:"id,omitempty"`
	Name                  string `json:"name,omitempty"`
	GenericName           string `json:"generic_name,omitempty"`
	Manufacturer          string `json:"manufacturer,omitempty"`
	ProductClassification string `json:"product_classification,omitempty"`
	ProductForm           string `json:"product_form,omitempty"`
	IsActive              bool   `json:"is_active"`
}

type ProductDetail struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type AddProductRequest struct {
	ManufacturerID          int             `binding:"required" form:"manufacturer_id"`
	ProductClassificationID int             `binding:"required" form:"product_classification_id"`
	ProductFormID           int             `binding:"required" form:"product_form_id"`
	Name                    string          `binding:"required,max=75" form:"name"`
	GenericName             string          `binding:"required" form:"generic_name"`
	Categories              string          `binding:"required" form:"categories"`
	Description             string          `binding:"required" form:"description"`
	UnitInPack              *int            `binding:"required" form:"unit_in_pack"`
	SellingUnit             string          `binding:"required" form:"selling_unit"`
	Weight                  decimal.Decimal `binding:"required" form:"weight"`
	Height                  decimal.Decimal `binding:"required" form:"height"`
	Length                  decimal.Decimal `binding:"required" form:"length"`
	Width                   decimal.Decimal `binding:"required" form:"width"`
	IsActive                string          `binding:"required,boolean" form:"is_active"`
	ImageURL                string
}

func ConvertToProductModel(product AddProductRequest) *model.Product {
	isActive, _ := strconv.ParseBool(product.IsActive)
	agg := product.Name + " " + product.GenericName + " " + product.Categories + " " + product.Description + " " + product.SellingUnit
	return &model.Product{
		ManufacturerID:          product.ManufacturerID,
		ProductClassificationID: product.ProductClassificationID,
		ProductFormID:           product.ProductFormID,
		Name:                    product.Name,
		GenericName:             product.GenericName,
		Categories:              product.Categories,
		Description:             product.Description,
		UnitInPack:              product.UnitInPack,
		SellingUnit:             product.SellingUnit,
		Weight:                  product.Weight,
		Height:                  product.Height,
		Length:                  product.Length,
		Width:                   product.Width,
		ImageURL:                product.ImageURL,
		IsActive:                isActive,
		Agg:                     agg,
	}
}

type ListProductParams struct {
	Limit                 int    `form:"limit" binding:"omitempty,gte=1"`
	Page                  int    `form:"page" binding:"omitempty,gte=1"`
	SortBy                string `form:"sort_by" binding:"omitempty,oneof=name created_at stock"`
	Sort                  string `form:"sort" binding:"omitempty,oneof=asc desc"`
	Name                  string `form:"name" binding:"omitempty"`
	GenericName           string `form:"generic_name" binding:"omitempty"`
	Manufacturer          string `form:"manufacturer" binding:"omitempty"`
	ProductClassification string `form:"product_classification" binding:"omitempty"`
	ProductForm           string `form:"product_form" binding:"omitempty"`
	IsActive              string `form:"is_active" binding:"omitempty,boolean"`
	Usage                 string `form:"usage" binding:"omitempty,number"`
}

type ListProductDetailParams struct {
	Name string `form:"name" binding:"omitempty"`
}

func (p *ListProductParams) EnsureDefaults() {
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
