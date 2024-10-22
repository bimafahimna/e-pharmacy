package middleware

import (
	"net/http"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/appconst"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwt jwt.Provider
}

func NewAuthMiddleware(jwt jwt.Provider) *AuthMiddleware {
	return &AuthMiddleware{
		jwt: jwt,
	}
}

func (m *AuthMiddleware) RequireToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.SetSameSite(http.SameSiteNoneMode)

		token, err := ctx.Cookie("token")
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
				Message: apperror.ErrUnauthorized.Message,
			})
			return
		}

		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
				Message: apperror.ErrUnauthorized.Message,
			})
			return
		}

		claims, err := m.jwt.Parse(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
				Message: apperror.ErrUnauthorized.Message,
			})
			return
		}

		ctx.Set(appconst.KeyUserID, claims.UserID)
		ctx.Set(appconst.KeyRole, claims.Role)
		ctx.Set(appconst.KeyIsVerified, claims.IsVerified)
		ctx.Next()
	}
}
