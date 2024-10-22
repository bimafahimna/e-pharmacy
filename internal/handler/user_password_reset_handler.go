package handler

import (
	"net/http"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/appconst"
	"github.com/gin-gonic/gin"
)

func (h *UserHandler) ResetPassword(ctx *gin.Context) {
	req := new(dto.ResetPasswordRequest)

	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.Error(err)
		return
	}

	if err := h.useCase.ResetPassword(ctx, *req); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, dto.Response{
		Message: appconst.MsgResetPasswordOk,
	})
}

func (h *UserHandler) ConfirmReset(ctx *gin.Context) {
	req := new(dto.ConfirmResetRequest)

	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.Error(err)
		return
	}

	if err := h.useCase.ConfirmResetPassword(ctx, *req); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, dto.Response{
		Message: appconst.MsgConfirmResetOk,
	})
}
