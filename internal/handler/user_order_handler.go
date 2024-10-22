package handler

import (
	"net/http"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/appconst"
	"github.com/gin-gonic/gin"
)

func (h *UserHandler) CreateOrder(ctx *gin.Context) {
	var req dto.CreateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	if err := h.useCase.CreateOrder(ctx, ctx.GetInt64(appconst.KeyUserID), req); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, dto.Response{
		Message: appconst.MsgCreateOrderCreated,
	})
}

func (h *UserHandler) ListUnpaidOrder(ctx *gin.Context) {
	data, err := h.useCase.ListUnpaidOrder(ctx, ctx.GetInt64(appconst.KeyUserID))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message: appconst.MsgListUnpaidOrderOk,
		Data:    data,
	})
}

func (h *UserHandler) ListOrder(ctx *gin.Context) {
	var query dto.ListOrderQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.Error(err)
		return
	}

	data, err := h.useCase.ListOrder(ctx, ctx.GetInt64(appconst.KeyUserID), query.Status)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message: appconst.MsgListOrderOk,
		Data:    data,
	})
}

func (h *UserHandler) UploadPaymentProof(ctx *gin.Context) {
	var uri dto.PaymentUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.Error(err)
		return
	}

	if err := h.useCase.UploadPaymentProof(ctx, uri.PaymentID, ctx.GetString("image_url")); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message: appconst.MsgUploadPaymentProofOk,
	})
}

func (h *UserHandler) ConfirmOrder(ctx *gin.Context) {
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

	userID := ctx.GetInt64(appconst.KeyUserID)

	if err := h.useCase.UpdateOrderStatus(ctx, uri.OrderID, userID, req.Status); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message: appconst.MsgUpdateOrderStatusOk,
	})
}
