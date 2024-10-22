package handler

import (
	"net/http"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/appconst"
	"github.com/gin-gonic/gin"
)

func (h *UserHandler) GetProductDetails(ctx *gin.Context) {
	var queries dto.GetProductDetailParams
	if err := ctx.ShouldBindQuery(&queries); err != nil {
		ctx.Error(err)
		return
	}
	queries.EnsureDefaults()

	var uri dto.GetProductDetailURI
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.Error(err)
		return
	}

	data, err := h.useCase.GetProductDetails(ctx, uri, queries)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dto.Response{
		Message: appconst.MsgListPopularProductOk,
		Data:    data,
	})
}
func (h *UserHandler) AvailablePharmacy(ctx *gin.Context) {
	var uri dto.AvailablePharmacyURI
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.Error(err)
		return
	}

	data, err := h.useCase.AvailablePharmacy(ctx, uri.ProductID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message: appconst.MsgListAvailPharmaciesOk,
		Data:    data,
	})
}
