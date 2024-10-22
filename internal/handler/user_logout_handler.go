package handler

import (
	"net/http"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/appconst"
	"github.com/gin-gonic/gin"
)

func (u *UserHandler) Logout(ctx *gin.Context) {
	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie("token", "", -1, "/", u.config.URL.DomainName, true, true)
	ctx.SetCookie("is_verified", "", -1, "/", u.config.URL.DomainName, true, true)
	ctx.SetCookie("role", "", -1, "/", u.config.URL.DomainName, true, true)

	ctx.JSON(http.StatusOK, dto.Response{
		Message: appconst.MsgLogoutOk,
	})
}
