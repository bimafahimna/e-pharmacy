package handler

import (
	"net/http"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/appconst"
	"github.com/gin-gonic/gin"
)

func (h *AdminHandler) ProcessOrder(ctx *gin.Context) {
	var uri dto.PaymentUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.Error(err)
		return
	}

	var req dto.ProcessOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	if err := h.useCase.UpdateOrderStatus(ctx, uri.PaymentID, req.Status); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message: appconst.MsgUpdateOrderStatusOk,
	})
}

func (h *AdminHandler) ConfirmOrder(ctx *gin.Context) {
	var uri dto.OrderUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.Error(err)
		return
	}

	var req dto.ConfirmOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	if err := h.useCase.UpdateOrderStatus(ctx, uri.OrderID, req.Status); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message: appconst.MsgUpdateOrderStatusOk,
	})
}
