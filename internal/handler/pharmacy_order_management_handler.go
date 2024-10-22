package handler

import (
	"net/http"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/appconst"
	"github.com/gin-gonic/gin"
)

func (h *PharmacistHandler) ListOrder(ctx *gin.Context) {
	var query dto.ListPharmacyOrderQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.Error(err)
		return
	}
	query.EnsureDefaults()

	data, pagination, err := h.useCase.ListOrder(ctx, ctx.GetInt64(appconst.KeyUserID), query)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message:    appconst.MsgListOrderOk,
		Data:       data,
		Pagination: pagination,
	})
}

func (h *PharmacistHandler) SendOrder(ctx *gin.Context) {
	var uri dto.OrderUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.Error(err)
		return
	}

	var req dto.SendOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	if err := h.useCase.SendOrder(ctx, ctx.GetInt64(appconst.KeyUserID), uri.OrderID, req); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message: appconst.MsgUpdateOrderStatusOk,
	})
}
