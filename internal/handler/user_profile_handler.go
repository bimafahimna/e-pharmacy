package handler

import (
	"net/http"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/appconst"
	"github.com/gin-gonic/gin"
)

func (h *UserHandler) AddAddress(ctx *gin.Context) {
	req := new(dto.CustomerAddressRequest)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	userId := ctx.GetInt64(appconst.KeyUserID)
	if err := h.useCase.AddAddress(ctx, *req, userId); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, dto.Response{Message: appconst.MsgAddAddressCreated})
}

func (h *UserHandler) GetAddresses(ctx *gin.Context) {
	userId := ctx.GetInt64(appconst.KeyUserID)
	data, err := h.useCase.GetAddresses(ctx, userId)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, dto.Response{Message: appconst.MsgListAddressesOk, Data: data})
}
