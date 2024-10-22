package dto

import (
	"strconv"
	"time"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/model"
	"github.com/shopspring/decimal"
)

type Logistic struct {
	ID            int64           `json:"id"`
	Name          string          `json:"name"`
	LogoUrl       string          `json:"image_url"`
	Service       string          `json:"service"`
	Estimation    string          `json:"estimation,omitempty"`
	Price         decimal.Decimal `json:"price,omitempty"`
	IsRecommended bool            `json:"is_recommended,omitempty"`
}

type ListPharmacyLogistics struct {
	PharmacyID   int        `json:"pharmacy_id"`
	PharmacyName string     `json:"pharmacy_name"`
	Logistics    []Logistic `json:"logistics"`
}

type ListPharmacyParams struct {
	Limit          int    `form:"limit" binding:"omitempty,gte=1"`
	Page           int    `form:"page" binding:"omitempty,gte=1"`
	SortBy         string `form:"sort_by" binding:"omitempty,oneof=name created_at is_assigned"`
	Sort           string `form:"sort" binding:"omitempty,oneof=asc desc"`
	Name           string `form:"name" binding:"omitempty"`
	PharmacistName string `form:"pharmacist_name" binding:"omitempty"`
	PartnerName    string `form:"partner_name" binding:"omitempty"`
	CityName       string `form:"city_name" binding:"omitempty"`
	Address        string `form:"address" binding:"omitempty"`
	IsActive       string `form:"is_active" binding:"omitempty,boolean"`
}

type PharmacyResponse struct {
	ID         int64              `json:"id"`
	Pharmacist PharmacyPharmacist `json:"pharmacist"`
	Partner    PharmacyPartner    `json:"partner"`
	Name       string             `json:"name"`
	Address    string             `json:"address"`
	City       PharmacyCity       `json:"city"`
	Latitude   decimal.Decimal    `json:"latitude"`
	Longitude  decimal.Decimal    `json:"longitude"`
	IsActive   string             `json:"is_active"`
	CreatedAt  time.Time          `json:"created_at"`
}

type PharmacyPharmacist struct {
	ID   *int    `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

type PharmacyPartner struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type PharmacyCity struct {
	ID   int     `json:"id"`
	Name *string `json:"name,omitempty"`
}

type AddPharmacyRequest struct {
	PharmacistID int    `json:"pharmacist_id" binding:"required"`
	PartnerID    int    `json:"partner_id" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Address      string `json:"address" binding:"required"`
	CityId       int    `json:"city_id" binding:"required"`
	Logistics    []int  `json:"logistics" binding:"required"`
	Latitude     string `json:"latitude" binding:"required,numeric"`
	Longitude    string `json:"longitude" binding:"required,numeric"`
	IsActive     string `json:"is_active" binding:"required,boolean"`
}

func (p *ListPharmacyParams) EnsureDefaults() {
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

func (p *ListPharmacyParams) Filters() (filters map[string]interface{}) {
	filters = map[string]interface{}{
		"name":            p.Name,
		"pharmacist_name": p.PharmacistName,
		"partner_name":    p.PartnerName,
		"city_name":       p.CityName,
		"address":         p.Address,
	}

	if p.IsActive != "" {
		isActive, _ := strconv.ParseBool(p.IsActive)
		filters["is_active"] = isActive
	}
	return
}

func ConvertToPharmacyModel(pharmacy AddPharmacyRequest) *model.Pharmacy {
	isActive, _ := strconv.ParseBool(pharmacy.IsActive)
	latitude, _ := decimal.NewFromString(pharmacy.Latitude)
	longitude, _ := decimal.NewFromString(pharmacy.Longitude)
	return &model.Pharmacy{
		PharmacistId: &pharmacy.PharmacistID,
		PartnerId:    pharmacy.PartnerID,
		Name:         pharmacy.Name,
		Address:      pharmacy.Address,
		CityId:       pharmacy.CityId,
		Latitude:     latitude,
		Longitude:    longitude,
		IsActive:     isActive,
	}
}

func ConvertToPharmacyDto(pharmacy model.Pharmacy) PharmacyResponse {
	return PharmacyResponse{
		ID:         int64(pharmacy.ID),
		Pharmacist: PharmacyPharmacist{ID: pharmacy.PharmacistId, Name: pharmacy.PharmacistName},
		Partner:    PharmacyPartner{ID: pharmacy.PartnerId, Name: pharmacy.PartnerName},
		Name:       pharmacy.Name,
		Address:    pharmacy.Address,
		City:       PharmacyCity{ID: pharmacy.CityId, Name: pharmacy.CityName},
		Latitude:   pharmacy.Latitude,
		Longitude:  pharmacy.Longitude,
		IsActive:   strconv.FormatBool(pharmacy.IsActive),
		CreatedAt:  pharmacy.CreatedAt,
	}
}

func ConvertToListPharmacies(pharmacyModel []model.Pharmacy) []PharmacyResponse {
	pharmacies := []PharmacyResponse{}
	for _, pharmacyModel := range pharmacyModel {
		pharmacy := ConvertToPharmacyDto(pharmacyModel)
		pharmacies = append(pharmacies, pharmacy)
	}
	return pharmacies
}
