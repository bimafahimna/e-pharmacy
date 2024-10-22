package handler

import (
	"net/http"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/appconst"
	"github.com/gin-gonic/gin"
)

func (h *AdminHandler) AddProduct(ctx *gin.Context) {
	var req dto.AddProductRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(err)
		return
	}
	req.ImageURL = ctx.GetString("image_url")

	if err := h.useCase.AddProduct(ctx, req); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, dto.Response{Message: appconst.MsgAddProductOk})
}

func (h *AdminHandler) ListProduct(ctx *gin.Context) {
	var params dto.ListProductParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.Error(err)
		return
	}
	params.EnsureDefaults()

	data, pagination, err := h.useCase.ListProduct(ctx, params)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message:    appconst.MsgListProductOk,
		Data:       data,
		Pagination: pagination,
	})
}

func (h *AdminHandler) ListManufacturer(ctx *gin.Context) {
	var params dto.ListProductDetailParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.Error(err)
		return
	}

	data, err := h.useCase.ListManufacturer(ctx, params)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message: appconst.MsgListManufacturerOk,
		Data:    data,
	})
}

func (h *AdminHandler) ListProductClassification(ctx *gin.Context) {
	var params dto.ListProductDetailParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.Error(err)
		return
	}

	data, err := h.useCase.ListProductClassification(ctx, params)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message: appconst.MsgListProductClassificationOk,
		Data:    data,
	})
}

func (h *AdminHandler) ListProductForm(ctx *gin.Context) {
	var params dto.ListProductDetailParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.Error(err)
		return
	}

	data, err := h.useCase.ListProductForm(ctx, params)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message: appconst.MsgListProductFormOk,
		Data:    data,
	})
}
