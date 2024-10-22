package handler

import (
	"net/http"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/cache"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/gin-gonic/gin"
)

type CacheHandler struct {
	cache cache.Provider
}

func NewCacheHandler(cache cache.Provider) *CacheHandler {
	return &CacheHandler{cache: cache}
}

func (h *CacheHandler) FlushAll(ctx *gin.Context) {
	err := h.cache.FlushAll()
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dto.Response{
		Message: "Successfully flush all cache",
	})
}

func (h *CacheHandler) DeleteAll(ctx *gin.Context) {
	err := h.cache.DeleteAll()
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dto.Response{
		Message: "Successfully delete all cache",
	})
}
