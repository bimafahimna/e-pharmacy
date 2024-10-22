package middleware

import (
	"net/http"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/appconst"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
	"github.com/gin-gonic/gin"
)

func VerifiedOnly() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		isVerified := ctx.GetBool(appconst.KeyIsVerified)
		if !isVerified {
			ctx.AbortWithStatusJSON(http.StatusForbidden, dto.Response{
				Message: apperror.ErrForbiddenNotVerified.Error(),
			})
		}
		ctx.Next()
	}
}
