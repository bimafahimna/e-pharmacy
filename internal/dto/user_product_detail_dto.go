package dto

import "github.com/shopspring/decimal"

type GetProductDetailParams struct {
	Limit     int    `form:"limit" binding:"omitempty"`
	Latitude  string `form:"latitude" binding:"omitempty"`
	Longitude string `form:"longitude" binding:"omitempty"`
}

type GetProductDetailURI struct {
	PharmacyID int `uri:"id" binding:"required"`
	ProductID  int `uri:"product_id" binding:"required"`
}

type AvailablePharmacyURI struct {
	ProductID int `uri:"product_id" binding:"required"`
}

type GetProductDetail struct {
	ProductID                     int             `json:"product_id"`
	ProductName                   string          `json:"product_name"`
	ProductGenericName            string          `json:"product_generic_name"`
	ProductImageUrl               string          `json:"product_image_url"`
	ProductForm                   string          `json:"product_form"`
	ProductManufacturer           string          `json:"product_manufacturer"`
	ProductClassification         string          `json:"product_classification"`
	ProductCategory               string          `json:"product_category"`
	ProductDescription            string          `json:"product_description"`
	ProductSellingUnit            string          `json:"product_selling_unit"`
	ProductUnitInPack             *int            `json:"product_unit_in_pack"`
	ProductWeight                 decimal.Decimal `json:"product_weight"`
	SelectedProductStock          int             `json:"selected_product_stock"`
	SelectedProductPrice          string          `json:"selected_product_price"`
	SelectedProductSoldAmount     int             `json:"selected_product_sold_amount"`
	SelectedPharmacyID            int             `json:"selected_pharmacy_id"`
	SelectedPharmacyName          string          `json:"selected_pharmacy_name"`
	SelectedPharmacistName        string          `json:"selected_pharmacist_name"`
	SelectedPharmacistPhoneNumber string          `json:"selected_pharmacist_phone_number"`
	SelectedPharmacistSIPANumber  string          `json:"selected_pharmacist_SIPA_number"`
	SelectedPharmacyAddress       string          `json:"selected_pharmacist_address"`
	SelectedPharmacyCityName      *string         `json:"selected_pharmacist_city_name"`
}

type AvailablePharmacy struct {
	PharmacyID   int    `json:"pharmacy_id"`
	PharmacyName string `json:"pharmacy_name"`
	PartnerName  string `json:"partner_name"`
	PartnerLogo  string `json:"partner_logo"`
	Address      string `json:"address"`
	CityName     string `json:"city_name"`
	ProductID    int    `json:"product_id"`
	ProductPrice string `json:"product_price"`
	Stock        int    `json:"stock"`
}

func (q *GetProductDetailParams) EnsureDefaults() {
	if q.Limit == 0 {
		q.Limit = 5
	}
	if q.Latitude == "" {
		q.Latitude = "0"
	}
	if q.Longitude == "" {
		q.Longitude = "0"
	}
}
