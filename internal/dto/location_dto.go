package dto

type ProvinceResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"province"`
}

type CityResponse struct {
	ID                   int64  `json:"id"`
	ProvinceId           int64  `json:"province_id"`
	CityName             string `json:"city_name"`
	CityType             string `json:"city_type"`
	CityUnofficialId     *int64 `json:"city_unofficial_id"`
	ProvinceUnofficialId int64  `json:"province_unofficial_id"`
}

type DistrictResponse struct {
	ID     int64  `json:"id"`
	CityId int64  `json:"city_id"`
	Name   string `json:"district"`
}

type SubDistrictResponse struct {
	ID         int64  `json:"id"`
	DistrictId int64  `json:"district_id"`
	Name       string `json:"sub_district"`
}

type ListCityParams struct {
	ProvinceId *string `form:"province_id" binding:"omitempty"`
	Name       string  `form:"name" binding:"omitempty"`
}

func (p *ListCityParams) EnsureDefaults() {
	if p.ProvinceId != nil && *p.ProvinceId == "" {
		p.ProvinceId = nil
	}
}

type ListDistrictParams struct {
	CityId *string `form:"city_id" binding:"omitempty"`
	Name   string  `form:"name" binding:"omitempty"`
}

func (p *ListDistrictParams) EnsureDefaults() {
	if p.CityId != nil && *p.CityId == "" {
		p.CityId = nil
	}
}

type ListSubDistrictParams struct {
	DistrictId *string `form:"district_id" binding:"omitempty"`
	Name       string  `form:"name" binding:"omitempty"`
}

func (p *ListSubDistrictParams) EnsureDefaults() {
	if p.DistrictId != nil && *p.DistrictId == "" {
		p.DistrictId = nil
	}
}
