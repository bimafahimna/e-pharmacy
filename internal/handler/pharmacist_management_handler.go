package handler

import (
	"net/http"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/appconst"
	"github.com/gin-gonic/gin"
)

func (h *AdminHandler) AddPharmacist(ctx *gin.Context) {
	var req dto.AddPharmacistRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	if err := h.useCase.AddPharmacist(ctx, req); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, dto.Response{
		Message: appconst.MsgAddPharmacistCreated,
	})
}

func (h *AdminHandler) ListPharmacists(ctx *gin.Context) {
	var params dto.ListPharmacistParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.Error(err)
		return
	}
	params.EnsureDefaults()

	data, pagination, err := h.useCase.ListPharmacists(ctx, params)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message:    appconst.MsgListPharmacistsOk,
		Data:       data,
		Pagination: pagination,
	})
}

func (h *AdminHandler) GetPharmacist(ctx *gin.Context) {
	var uri dto.PharmacistUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.Error(err)
		return
	}

	data, err := h.useCase.GetPharmacist(ctx, uri.UserID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message: appconst.MsgGetPharmacistOk,
		Data:    data,
	})
}

func (h *AdminHandler) EditPharmacist(ctx *gin.Context) {
	var uri dto.PharmacistUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.Error(err)
		return
	}

	var req dto.UpdatePharmacistRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	if err := h.useCase.EditPharmacist(ctx, uri.UserID, req); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message: appconst.MsgUpdatePharmacistOk,
	})
}

func (h *AdminHandler) RemovePharmacist(ctx *gin.Context) {
	var uri dto.PharmacistUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.Error(err)
		return
	}

	if err := h.useCase.RemovePharmacist(ctx, uri.UserID); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message: appconst.MsgRemovePharmacistOk,
	})
}
