package handler

import (
	"net/http"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/appconst"
	"github.com/gin-gonic/gin"
)

func (h *AdminHandler) ListUsers(ctx *gin.Context) {
	var params dto.ListUserParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.Error(err)
		return
	}
	params.EnsureDefaults()

	data, pagination, err := h.useCase.ListUsers(ctx, params)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message:    appconst.MsgListUserOk,
		Data:       data,
		Pagination: pagination,
	})
}
