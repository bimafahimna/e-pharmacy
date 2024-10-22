package dto

import "github.com/shopspring/decimal"

type Item struct {
	PharmacyID    int             `json:"selected_pharmacy_id"`
	ProductID     int             `json:"product_id"`
	ProductFormID int             `json:"product_form_id"`
	Name          string          `json:"product_name"`
	Price         decimal.Decimal `json:"product_price"`
	Stock         int             `json:"product_stock"`
	SellingUnit   string          `json:"product_selling_unit"`
	ImageURL      string          `json:"image_url"`
}

type ListPopularProductQueries struct {
	Limit     int    `form:"limit" binding:"omitempty"`
	Latitude  string `form:"latitude" binding:"omitempty"`
	Longitude string `form:"longitude" binding:"omitempty"`
}

type ListProductQueries struct {
	Limit     int
	Page      int    `form:"page" binding:"omitempty"`
	Search    string `form:"s" binding:"omitempty"`
	Latitude  string `form:"latitude" binding:"omitempty"`
	Longitude string `form:"longitude" binding:"omitempty"`
}

func (q *ListPopularProductQueries) EnsureDefaults() {
	if q.Limit == 0 {
		q.Limit = 24
	}
	if q.Latitude == "" {
		q.Latitude = "0"
	}
	if q.Longitude == "" {
		q.Longitude = "0"
	}
}

func (q *ListProductQueries) EnsureDefaults() {
	q.Limit = 20
	// if q.Limit == 0 {
	// }
	if q.Page == 0 {
		q.Page = 1
	}
	if q.Latitude == "" {
		q.Latitude = "0"
	}
	if q.Longitude == "" {
		q.Longitude = "0"
	}
	if q.Search == "" {
		q.Search = "Obat"
	}
}
