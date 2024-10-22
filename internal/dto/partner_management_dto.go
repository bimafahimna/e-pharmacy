package dto

import (
	"strconv"
	"time"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/model"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/util"
)

type PartnerResponse struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	LogoUrl          string `json:"logo_url"`
	YearFounded      int    `json:"year_founded"`
	ActiveDays       string `json:"active_days"`
	OperationalStart string `json:"operational_start"`
	OperationalStop  string `json:"operational_stop"`
	IsActive         string `json:"is_active"`
}

type Partner struct {
	ID               int64     `json:"id"`
	Name             string    `json:"name"`
	LogoUrl          string    `json:"logo_url"`
	YearFounded      int       `json:"year_founded"`
	ActiveDays       string    `json:"active_days"`
	OperationalStart string    `json:"operational_start"`
	OperationalStop  string    `json:"operational_stop"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type PartnerForm struct {
	Tags string `form:"tags" binding:"required"`
}

type AddPartnerRequest struct {
	Name             string `json:"name" binding:"required"`
	YearFounded      int    `json:"year_founded" binding:"required"`
	ActiveDays       []int  `json:"active_days" binding:"required"`
	OperationalStart string `json:"operational_start" binding:"required"`
	OperationalStop  string `json:"operational_stop" binding:"required"`
	IsActive         string `json:"is_active" binding:"required,boolean"`
	LogoUrl          string
}

type PartnerUri struct {
	ID int `uri:"id" binding:"required"`
}

type EditPartnerRequest struct {
	ID               int
	LogoUrl          string `json:"logo_url" binding:"required"`
	Name             string `json:"name" binding:"required"`
	YearFounded      int    `json:"year_founded" binding:"required"`
	ActiveDays       []int  `json:"active_days" binding:"required"`
	OperationalStart string `json:"operational_start" binding:"required"`
	OperationalStop  string `json:"operational_stop" binding:"required"`
	IsActive         string `json:"is_active" binding:"required"`
}

type EditPartnerDaysAndHoursRequest struct {
	ID               int
	ActiveDays       []int  `json:"active_days" binding:"required"`
	OperationalStart string `json:"operational_start" binding:"required"`
	OperationalStop  string `json:"operational_stop" binding:"required"`
}

type ListPartnerParams struct {
	Limit            int     `form:"limit" binding:"omitempty,gte=1"`
	Page             int     `form:"page" binding:"omitempty,gte=1"`
	SortBy           string  `form:"sort_by" binding:"omitempty,oneof=id name year_founded operational_start operational_stop is_active created_at"`
	Sort             string  `form:"sort" binding:"omitempty,oneof=asc desc"`
	Id               *string `form:"id" binding:"omitempty"`
	Name             string  `form:"name" binding:"omitempty"`
	YearFounded      *string `form:"year_founded" binding:"omitempty"`
	ActiveDays       string  `form:"active_days" binding:"omitempty"`
	OperationalStart *string `form:"operational_start" binding:"omitempty"`
	OperationalStop  *string `form:"operational_stop" binding:"omitempty"`
	IsActive         *string `form:"is_active" binding:"omitempty"`
}

func (p *ListPartnerParams) EnsureDefaults() {
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
	if p.Id != nil && *p.Id == "" {
		p.Id = nil
	}
	if p.YearFounded != nil && *p.YearFounded == "" {
		p.YearFounded = nil
	}
	if p.IsActive != nil && *p.IsActive == "" {
		p.IsActive = nil
	}
	if p.OperationalStart != nil && *p.OperationalStart == "" {
		p.OperationalStart = nil
	}
	if p.OperationalStop != nil && *p.OperationalStop == "" {
		p.OperationalStop = nil
	}
	if p.IsActive != nil && *p.IsActive == "" {
		p.IsActive = nil
	}
}

func ConvertAddPartnerDtoToModel(dto AddPartnerRequest) *model.Partner {
	isActive, _ := strconv.ParseBool(dto.IsActive)
	return &model.Partner{
		Name:             dto.Name,
		LogoUrl:          dto.LogoUrl,
		YearFounded:      dto.YearFounded,
		ActiveDays:       util.IntArrayToString(dto.ActiveDays),
		OperationalStart: dto.OperationalStart,
		OperationalStop:  dto.OperationalStop,
		IsActive:         isActive,
	}
}

func ConvertEditPartnerDtoToModel(dto EditPartnerRequest) *model.Partner {
	isActive, _ := strconv.ParseBool(dto.IsActive)
	return &model.Partner{
		ID:               int64(dto.ID),
		Name:             dto.Name,
		LogoUrl:          dto.LogoUrl,
		YearFounded:      dto.YearFounded,
		ActiveDays:       util.IntArrayToString(dto.ActiveDays),
		OperationalStart: dto.OperationalStart,
		OperationalStop:  dto.OperationalStop,
		IsActive:         isActive,
	}
}
