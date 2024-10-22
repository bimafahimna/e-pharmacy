package handler

import (
	"net/http"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/appconst"
	"github.com/gin-gonic/gin"
)

func (h *AdminHandler) AddPharmacy(ctx *gin.Context) {
	req := new(dto.AddPharmacyRequest)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	if err := h.useCase.AddPharmacy(ctx, *req); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, dto.Response{
		Message: appconst.MsgAddPharmacyCreated,
	})
}

func (h *AdminHandler) ListPharmacies(ctx *gin.Context) {
	var params dto.ListPharmacyParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.Error(err)
		return
	}
	params.EnsureDefaults()

	data, pagination, err := h.useCase.ListPharmacies(ctx, params)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message:    appconst.MsgListPharmaciesOk,
		Data:       data,
		Pagination: pagination,
	})
}
