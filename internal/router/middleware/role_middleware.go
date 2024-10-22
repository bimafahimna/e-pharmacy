package middleware

import (
	"net/http"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/appconst"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
	"github.com/gin-gonic/gin"
)

func Authorize(requiredRole string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role := ctx.GetString(appconst.KeyRole)
		if role != requiredRole {
			ctx.AbortWithStatusJSON(http.StatusForbidden, dto.Response{
				Message: apperror.ErrForbiddenRole.Message,
			})
		}

		ctx.Next()
	}
}
