package dto

type ProcessOrderRequest struct {
	Status string `json:"status" binding:"required,eq=Processed"`
}

type ConfirmOrderRequest struct {
	Status string `json:"status" binding:"required,eq=Order Confirmed"`
}
