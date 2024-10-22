package handler

import (
	"net/http"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/appconst"
	"github.com/gin-gonic/gin"
)

func (h *PharmacistHandler) ListMasterProducts(ctx *gin.Context) {
	var params dto.ListProductParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.Error(err)
		return
	}
	params.EnsureDefaults()

	data, pagination, err := h.useCase.ListMasterProducts(ctx, ctx.GetInt64(appconst.KeyUserID), params)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message:    appconst.MsgListPharmacyProductsOk,
		Data:       data,
		Pagination: pagination,
	})
}

func (h *PharmacistHandler) ListPharmacyProducts(ctx *gin.Context) {
	var params dto.ListPharmacyProductParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.Error(err)
		return
	}
	params.EnsureDefaults()

	data, pagination, err := h.useCase.ListPharmacyProducts(ctx, ctx.GetInt64(appconst.KeyUserID), params)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message:    appconst.MsgListPharmacyProductsOk,
		Data:       data,
		Pagination: pagination,
	})
}

func (h *PharmacistHandler) GetPharmacyProduct(ctx *gin.Context) {
	var uri dto.ProductUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.Error(err)
		return
	}

	data, err := h.useCase.GetPharmacyProduct(ctx, ctx.GetInt64(appconst.KeyUserID), uri.ProductID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message: appconst.MsgGetPharmacyProductOk,
		Data:    data,
	})
}

func (h *PharmacistHandler) AddPharmacyProduct(ctx *gin.Context) {
	var req dto.AddPharmacyProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	if err := h.useCase.AddPharmacyProduct(ctx, ctx.GetInt64(appconst.KeyUserID), req); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, dto.Response{
		Message: appconst.MsgAddPharmacyProductCreated,
	})
}

func (h *PharmacistHandler) UpdatePharmacyProduct(ctx *gin.Context) {
	var uri dto.ProductUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.Error(err)
		return
	}

	var req dto.UpdatePharmacyProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	if err := h.useCase.UpdatePharmacyProduct(ctx, ctx.GetInt64(appconst.KeyUserID), uri.ProductID, req); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message: appconst.MsgUpdatePharmacyProductOk,
	})
}

func (h *PharmacistHandler) DeletePharmacyProduct(ctx *gin.Context) {
	var uri dto.ProductUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.Error(err)
		return
	}

	if err := h.useCase.DeletePharmacyProduct(ctx, ctx.GetInt64(appconst.KeyUserID), uri.ProductID); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message: appconst.MsgDeletePharmacyProductOk,
	})
}
