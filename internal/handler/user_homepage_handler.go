package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/appconst"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/logger"
	"github.com/gin-gonic/gin"
)

func (h *UserHandler) ListPopularProduct(ctx *gin.Context) {
	var queries dto.ListPopularProductQueries
	if err := ctx.ShouldBindQuery(&queries); err != nil {
		ctx.Error(err)
		return
	}
	queries.EnsureDefaults()
	userId := ctx.GetString(appconst.KeyUserID)

	var (
		data []dto.Item
		err  error
	)
	if queries.Latitude == "0" || queries.Longitude == "0" {
		if value, err := h.cache.Get(appconst.CacheBestSeller); err == nil {
			err := json.Unmarshal(value, &data)
			if err != nil {
				ctx.Error(err)
				return
			}
			ctx.JSON(http.StatusOK, dto.Response{
				Message: appconst.MsgListPopularProductOk,
				Data:    data,
			})
			return
		}
		data, err = h.useCase.ListBestsellerProduct(ctx, queries.Limit)
		if err != nil {
			ctx.Error(err)
			return
		}
		go func() {
			value, err := json.Marshal(data)
			if err != nil {
				logger.Log.Errorf("failed to marshal data: %v", err)
				return
			}
			err = h.cache.Set(appconst.CacheBestSeller, value, 60)
			if err != nil {
				logger.Log.Errorf("failed to set cache: %v", err)
				return
			}
		}()
	} else {
		if value, err := h.cache.Get(appconst.CacheBestSeller + userId); err == nil {
			err := json.Unmarshal(value, &data)
			if err != nil {
				ctx.Error(err)
				return
			}
			ctx.JSON(http.StatusOK, dto.Response{
				Message: appconst.MsgListPopularProductOk,
				Data:    data,
			})
			return
		}
		data, err = h.useCase.ListRecommendedProduct(ctx, queries)
		if err != nil {
			ctx.Error(err)
			return
		}
		go func() {
			value, err := json.Marshal(data)
			if err != nil {
				logger.Log.Errorf("failed to marshal data: %v", err)
				return
			}
			err = h.cache.Set(appconst.CacheBestSeller+userId, value, 60)
			if err != nil {
				logger.Log.Errorf("failed to set cache: %v", err)
				return
			}
		}()
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message: appconst.MsgListPopularProductOk,
		Data:    data,
	})
}

func (h *UserHandler) ListProduct(ctx *gin.Context) {
	var queries dto.ListProductQueries
	if err := ctx.ShouldBindQuery(&queries); err != nil {
		ctx.Error(err)
		return
	}
	queries.EnsureDefaults()

	data, pagination, err := h.useCase.ListProduct(ctx, queries)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message:    appconst.MsgListProductOk,
		Data:       data,
		Pagination: pagination,
	})
}
