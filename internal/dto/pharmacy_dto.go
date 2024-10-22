package dto

type ListLogisticsParams struct {
	AddressID int    `form:"address_id" binding:"required,numeric,gte=1"`
	Weight    string `form:"weight" binding:"required"`
}

type ListLogisticsURI struct {
	PharmacyID int `uri:"id" binding:"required"`
}
