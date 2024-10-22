package handler

import (
	"net/http"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/appconst"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/usecase"
	"github.com/gin-gonic/gin"
)

type PharmacyHandler struct {
	useCase usecase.PharmacyUseCase
}

func NewPharmacyHandler(useCase usecase.PharmacyUseCase) *PharmacyHandler {
	return &PharmacyHandler{
		useCase: useCase,
	}
}

func (h *PharmacyHandler) ListLogistics(ctx *gin.Context) {
	var address dto.ListLogisticsParams
	if err := ctx.ShouldBindQuery(&address); err != nil {
		ctx.Error(err)
		return
	}

	var pharmacy dto.ListLogisticsURI
	if err := ctx.ShouldBindUri(&pharmacy); err != nil {
		ctx.Error(err)
		return
	}

	data, err := h.useCase.ListLogistics(ctx, address.AddressID, pharmacy.PharmacyID, address.Weight)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message: appconst.MsgListPharmacyLogisticsOk,
		Data:    data,
	})
}
