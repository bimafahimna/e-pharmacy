package handler

import (
	"net/http"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/appconst"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/usecase"
	"github.com/gin-gonic/gin"
)

type LocationHandler struct {
	locationUseCase usecase.LocationUseCase
}

func NewLocationHandler(locationUseCase usecase.LocationUseCase) *LocationHandler {
	return &LocationHandler{
		locationUseCase: locationUseCase,
	}
}

func (h *LocationHandler) ListProvinces(ctx *gin.Context) {
	data, err := h.locationUseCase.ListProvinces(ctx, ctx.Query("name"))
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dto.Response{Message: appconst.MsgListProvinceOk, Data: data})
}

func (h *LocationHandler) ListCities(ctx *gin.Context) {
	var params dto.ListCityParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.Error(err)
		return
	}
	params.EnsureDefaults()

	data, err := h.locationUseCase.ListCities(ctx, params)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dto.Response{Message: appconst.MsgListCityOk, Data: data})
}

func (h *LocationHandler) ListDistricts(ctx *gin.Context) {
	var params dto.ListDistrictParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.Error(err)
		return
	}
	params.EnsureDefaults()

	data, err := h.locationUseCase.ListDistricts(ctx, params)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dto.Response{Message: appconst.MsgListCityOk, Data: data})
}

func (h *LocationHandler) ListSubDistricts(ctx *gin.Context) {
	var params dto.ListSubDistrictParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.Error(err)
		return
	}
	params.EnsureDefaults()

	data, err := h.locationUseCase.ListSubDistricts(ctx, params)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dto.Response{Message: appconst.MsgListCityOk, Data: data})
}
