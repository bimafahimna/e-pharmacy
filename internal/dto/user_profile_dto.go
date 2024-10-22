package dto

import (
	"strconv"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/model"
	"github.com/shopspring/decimal"
)

type CustomerAddressResponse struct {
	ID                  int64           `json:"id"`
	UserID              int64           `json:"user_id"`
	Name                string          `json:"name"`
	ReceiverName        string          `json:"receiver_name"`
	ReceiverPhoneNumber string          `json:"receiver_phone_number"`
	Latitude            decimal.Decimal `json:"latitude"`
	Longitude           decimal.Decimal `json:"longitude"`
	Province            string          `json:"province"`
	City                string          `json:"city"`
	District            string          `json:"district"`
	SubDistrict         string          `json:"sub_district"`
	AddressDetails      string          `json:"address_details"`
	IsActive            string          `json:"is_active"`
}

type CustomerAddressRequest struct {
	Name                string `json:"name" binding:"required"`
	ReceiverName        string `json:"receiver_name" binding:"required"`
	ReceiverPhoneNumber string `json:"receiver_phone_number" binding:"required"`
	Latitude            string `json:"latitude" binding:"required,numeric"`
	Longitude           string `json:"longitude" binding:"required,numeric"`
	Province            string `json:"province" binding:"required"`
	CityID              int64  `json:"city_id" binding:"required"`
	City                string `json:"city" binding:"required"`
	District            string `json:"district" binding:"required"`
	SubDistrict         string `json:"sub_district" binding:"required"`
	AddressDetails      string `json:"address_details" binding:"required"`
	IsActive            string `json:"is_active" binding:"required,boolean"`
}

func ConvertToCustomerAddressModel(address CustomerAddressRequest) *model.CustomerAddress {
	isActive, _ := strconv.ParseBool(address.IsActive)
	latitude, _ := decimal.NewFromString(address.Latitude)
	longitude, _ := decimal.NewFromString(address.Longitude)
	return &model.CustomerAddress{
		Name:                address.Name,
		ReceiverName:        address.ReceiverName,
		ReceiverPhoneNumber: address.ReceiverPhoneNumber,
		Latitude:            latitude,
		Longitude:           longitude,
		Province:            address.Province,
		CityID:              address.CityID,
		City:                address.City,
		District:            address.District,
		SubDistrict:         address.SubDistrict,
		AddressDetails:      address.AddressDetails,
		IsActive:            isActive,
	}
}

func ConvertToCustomerAddressResponseDto(address model.CustomerAddress) *CustomerAddressResponse {
	isActive := strconv.FormatBool(address.IsActive)
	return &CustomerAddressResponse{
		ID:                  address.ID,
		UserID:              address.UserID,
		Name:                address.Name,
		ReceiverName:        address.ReceiverName,
		ReceiverPhoneNumber: address.ReceiverPhoneNumber,
		Latitude:            address.Latitude,
		Longitude:           address.Longitude,
		Province:            address.Province,
		City:                address.City,
		District:            address.District,
		SubDistrict:         address.SubDistrict,
		AddressDetails:      address.AddressDetails,
		IsActive:            isActive,
	}
}
