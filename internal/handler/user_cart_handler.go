package handler

import (
	"net/http"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/appconst"
	"github.com/gin-gonic/gin"
)

func (h *UserHandler) GetCartItems(ctx *gin.Context) {
	items, err := h.useCase.GetItems(ctx, ctx.GetInt64(appconst.KeyUserID))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message: appconst.MsgGetCartItemsOk,
		Data:    items,
	})
}

func (h *UserHandler) AddCartItem(ctx *gin.Context) {
	var req dto.AddCartItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	if err := h.useCase.AddItem(ctx, ctx.GetInt64(appconst.KeyUserID), req); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, dto.Response{
		Message: appconst.MsgAddCartItemCreated,
	})
}

func (h *UserHandler) UpdateCartItem(ctx *gin.Context) {
	var uri dto.CartUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.Error(err)
		return
	}

	var req dto.UpdateCartItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	if err := h.useCase.UpdateItem(ctx, ctx.GetInt64(appconst.KeyUserID), uri.PharmacyID, uri.ProductID, req); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message: appconst.MsgUpdateCartItemOk,
	})
}

func (h *UserHandler) RemoveCartItem(ctx *gin.Context) {
	var uri dto.CartUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.Error(err)
		return
	}

	if err := h.useCase.RemoveItem(ctx, ctx.GetInt64(appconst.KeyUserID), uri.PharmacyID, uri.ProductID); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message: appconst.MsgRemoveCartItemOk,
	})
}
