package dto

type OrderUri struct {
	OrderID int `uri:"order_id" binding:"required"`
}

type SendOrderRequest struct {
	Status string `json:"status" binding:"required,eq=Sent"`
}
