package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/appconst"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
	"github.com/gin-gonic/gin"
)

func (h *AdminHandler) AddPartner(ctx *gin.Context) {
	var reqForm dto.PartnerForm
	if err := ctx.ShouldBind(&reqForm); err != nil {
		ctx.Error(err)
		return
	}
	tags := reqForm.Tags

	var req dto.AddPartnerRequest
	err := json.Unmarshal([]byte(tags), &req)
	if err != nil {
		ctx.Error(err)
		return
	}
	req.LogoUrl = ctx.GetString("image_url")

	if err := h.useCase.AddPartner(ctx, req); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, dto.Response{
		Message: appconst.MsgAddPartnerCreated,
	})
}

func (h *AdminHandler) ListPartners(ctx *gin.Context) {
	var params dto.ListPartnerParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.Error(err)
		return
	}
	params.EnsureDefaults()

	data, pagination, err := h.useCase.ListPartners(ctx, params)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message:    appconst.MsgListPartnersOk,
		Data:       data,
		Pagination: pagination,
	})
}

func (h *AdminHandler) GetPartnerByID(ctx *gin.Context) {
	var uri dto.PartnerUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.Error(apperror.ErrIDMustBeNumeric)
		return
	}

	data, err := h.useCase.GetPartnerByID(ctx, uri.ID)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dto.Response{
		Message: appconst.MsgGetPartnerOk,
		Data:    data,
	})
}

func (h *AdminHandler) EditPartner(ctx *gin.Context) {
	req := new(dto.EditPartnerRequest)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	var uri dto.PartnerUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.Error(apperror.ErrIDMustBeNumeric)
		return
	}

	req.ID = uri.ID
	if err := h.useCase.EditPartner(ctx, *req); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message: appconst.MsgUpdatePartnerOk,
	})
}

func (h *AdminHandler) EditPartnerDaysAndHours(ctx *gin.Context) {
	req := new(dto.EditPartnerDaysAndHoursRequest)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	var uri dto.PartnerUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.Error(apperror.ErrIDMustBeNumeric)
		return
	}

	req.ID = uri.ID
	if err := h.useCase.EditPartnerDaysAndHours(ctx, *req); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message: appconst.MsgUpdatePartnerOk,
	})
}
