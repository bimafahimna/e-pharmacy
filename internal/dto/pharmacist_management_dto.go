package dto

import "time"

type Pharmacist struct {
	UserID            int64     `json:"id"`
	Name              string    `json:"name"`
	Email             string    `json:"email"`
	SipaNumber        string    `json:"sipa_number"`
	WhatsappNumber    string    `json:"whatsapp_number"`
	YearsOfExperience string    `json:"years_of_experience"`
	IsAssigned        bool      `json:"is_assigned"`
	CreatedAt         time.Time `json:"created_at"`
}

type PharmacistUri struct {
	UserID int64 `uri:"id" binding:"required"`
}

type AddPharmacistRequest struct {
	Name              string `json:"name" binding:"required"`
	SipaNumber        string `json:"sipa_number" binding:"required"`
	WhatsappNumber    string `json:"whatsapp_number" binding:"required,numeric"`
	YearsOfExperience string `json:"years_of_experience" binding:"required,numeric"`
	Email             string `json:"email" binding:"required,email"`
	Password          string `json:"password" binding:"required,min=8"`
}

type UpdatePharmacistRequest struct {
	WhatsappNumber    string `json:"whatsapp_number" binding:"required,numeric"`
	YearsOfExperience string `json:"years_of_experience" binding:"required"`
}

type ListPharmacistParams struct {
	Limit             int    `form:"limit" binding:"omitempty,gte=1"`
	Page              int    `form:"page" binding:"omitempty,gte=1"`
	SortBy            string `form:"sort_by" binding:"omitempty,oneof=name created_at is_assigned"`
	Sort              string `form:"sort" binding:"omitempty,oneof=asc desc"`
	Name              string `form:"name" binding:"omitempty"`
	Email             string `form:"email" binding:"omitempty"`
	SipaNumber        string `form:"sipa_number" binding:"omitempty"`
	WhatsappNumber    string `form:"whatsapp_number" binding:"omitempty"`
	YearsOfExperience string `form:"years_of_experience" binding:"omitempty,number"`
	IsAssigned        string `form:"is_assigned" binding:"omitempty,boolean"`
}

func (p *ListPharmacistParams) EnsureDefaults() {
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
