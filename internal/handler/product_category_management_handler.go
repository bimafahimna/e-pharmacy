package handler

import (
	"net/http"
	"strconv"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/appconst"
	"github.com/gin-gonic/gin"
)

func (h *AdminHandler) ListProductCategories(ctx *gin.Context) {
	data, err := h.useCase.ListProductCategories(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message: appconst.MsgListProductCategoriesOk,
		Data:    data,
	})
}

func (h *AdminHandler) AddProductCategory(ctx *gin.Context) {
	var req dto.AddProductCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	if err := h.useCase.AddProductCategory(ctx, req); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, dto.Response{
		Message: appconst.MsgAddProductCategoryCreated,
	})
}

func (h *AdminHandler) UpdateProductCategory(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(err)
		return
	}

	var req dto.UpdateProductCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	if err := h.useCase.UpdateProductCategory(ctx, id, req); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message: appconst.MsgUpdateProductCategoryOk,
	})
}

func (h *AdminHandler) RemoveProductCategory(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(err)
		return
	}

	if err := h.useCase.RemoveProductCategory(ctx, id); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message: appconst.MsgRemoveProductCategoryOk,
	})
}
