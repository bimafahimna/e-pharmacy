package handler

import (
	"net/http"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/appconst"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/usecase"
	"github.com/gin-gonic/gin"
)

type LogisticHandler struct {
	useCase usecase.LogisticUseCase
}

func NewLogisticHandler(useCase usecase.LogisticUseCase) *LogisticHandler {
	return &LogisticHandler{
		useCase: useCase,
	}
}

func (h *LogisticHandler) ListLogistics(ctx *gin.Context) {
	data, err := h.useCase.ListLogistics(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dto.Response{
		Message: appconst.MsgListLogisticOk,
		Data:    data,
	})
}
