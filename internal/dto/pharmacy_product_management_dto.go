package dto

type PharmacyProduct struct {
	PharmacyID            int    `json:"pharmacy_id,omitempty"`
	ProductID             int    `json:"product_id,omitempty"`
	Stock                 int    `json:"stock,omitempty"`
	Price                 string `json:"price,omitempty"`
	Name                  string `json:"name,omitempty"`
	GenericName           string `json:"generic_name,omitempty"`
	Manufacturer          string `json:"manufacturer,omitempty"`
	ProductClassification string `json:"product_classification,omitempty"`
	ProductForm           string `json:"product_form,omitempty"`
	IsActive              bool   `json:"is_active"`
}

type ProductUri struct {
	ProductID int `uri:"product_id" binding:"required"`
}

type AddPharmacyProductRequest struct {
	ProductID int    `json:"product_id" binding:"required"`
	Price     string `json:"price" binding:"required,numeric"`
	Stock     int    `json:"stock" binding:"required,gte=1"`
}

type UpdatePharmacyProductRequest struct {
	Stock    int    `json:"stock" binding:"required"`
	IsActive string `json:"status" binding:"required,boolean"`
}

type ListPharmacyProductParams struct {
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
}

func (p *ListPharmacyProductParams) EnsureDefaults() {
	if p.Limit == 0 {
		p.Limit = 15
	}
	if p.Page == 0 {
		p.Page = 1
	}
	if p.SortBy == "" {
		p.SortBy = "name"
	}
	if p.Sort == "" {
		p.Sort = "asc"
	}
}
