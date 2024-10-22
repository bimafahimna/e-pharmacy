package dto

type Query struct {
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	Weight      int    `json:"weight"`
	Courier     string `json:"courier"`
}

type Status struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

type OriginDetails struct {
	CityID     string `json:"city_id"`
	ProvinceID string `json:"province_id"`
	Province   string `json:"province"`
	Type       string `json:"type"`
	CityName   string `json:"city_name"`
	PostalCode string `json:"postal_code"`
}

type DestinationDetails struct {
	CityID     string `json:"city_id"`
	ProvinceID string `json:"province_id"`
	Province   string `json:"province"`
	Type       string `json:"type"`
	CityName   string `json:"city_name"`
	PostalCode string `json:"postal_code"`
}

type Cost struct {
	Service     string `json:"service"`
	Description string `json:"description"`
	Cost        []struct {
		Value int    `json:"value"`
		Etd   string `json:"etd"`
		Note  string `json:"note"`
	} `json:"cost"`
}

type Result struct {
	Code  string `json:"code"`
	Name  string `json:"name"`
	Costs []Cost `json:"costs"`
}

type RajaOngkirResponse struct {
	Query              Query              `json:"query"`
	Status             Status             `json:"status"`
	OriginDetails      OriginDetails      `json:"origin_details"`
	DestinationDetails DestinationDetails `json:"destination_details"`
	Results            []Result           `json:"results"`
}
