package middleware

import (
	"net/http"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/config"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
	"github.com/gin-gonic/gin"
)

type workerMiddleware struct {
	config *config.WorkerConfig
}

func NewWorkerMiddleware(config *config.WorkerConfig) *workerMiddleware {
	return &workerMiddleware{
		config: config,
	}
}

func (m *workerMiddleware) Bypass() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		secretKey := ctx.Request.Header.Get("X-Secret-Key")
		if secretKey != m.config.SecretKey {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
				Message: apperror.ErrUnauthorized.Error(),
			})
			return
		}

		ctx.Next()
	}
}
