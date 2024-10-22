package handler

import (
	"net/http"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/appconst"
	"github.com/gin-gonic/gin"
)

func (h *UserHandler) Register(ctx *gin.Context) {
	var req dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	if err := h.useCase.Register(ctx, req); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, dto.Response{
		Message: appconst.MsgRegisterCreated,
	})
}

func (h *UserHandler) Verify(ctx *gin.Context) {
	var req dto.VerifyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	if err := h.useCase.Verify(ctx, req); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message: appconst.MsgVerifiedOk,
	})
}

func (h *UserHandler) GoogleRegister(ctx *gin.Context) {
	var req dto.GoogleRegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	err := h.useCase.GoogleRegister(ctx, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, dto.Response{
		Message: appconst.MsgRegisterCreated,
	})
}

func (h *UserHandler) RedirectGoogleRegister(ctx *gin.Context) {
	ctx.Redirect(http.StatusPermanentRedirect, h.config.Google.AuthCodeURL("register"))
}
