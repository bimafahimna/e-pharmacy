package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/appconst"
	"github.com/gin-gonic/gin"
)

func (h *UserHandler) Login(ctx *gin.Context) {
	loginReq := new(dto.LoginRequest)
	if err := ctx.ShouldBindJSON(loginReq); err != nil {
		ctx.Error(err)
		return
	}

	data, err := h.useCase.Login(ctx, *loginReq)
	if err != nil {
		ctx.Error(err)
		return
	}

	returnedData := dto.ReturnedURL{
		ReturnedURL: fmt.Sprintf("?is_verified=%s&role=%s&user_id=%d", data.IsVerified, data.Role, data.UserID),
	}

	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie("token", data.Token, 3600, "/", h.config.URL.DomainName, true, true)
	ctx.SetCookie("is_verified", data.IsVerified, 3600, "/", h.config.URL.DomainName, true, true)
	ctx.SetCookie("role", data.Role, 3600, "/", h.config.URL.DomainName, true, true)
	ctx.JSON(http.StatusOK, dto.Response{Message: "Successfully Logged in",
		Data: returnedData})
}

func (h *UserHandler) GoogleLogin(ctx *gin.Context) {
	var req dto.GoogleLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	data, err := h.useCase.GoogleLogin(ctx, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message: appconst.MsgLoginOk,
		Data:    data,
	})
}

func (h *UserHandler) RedirectGoogleLogin(ctx *gin.Context) {
	ctx.Redirect(http.StatusPermanentRedirect, h.config.Google.AuthCodeURL("login"))
}

func (h *UserHandler) GoogleLoginCallback(ctx *gin.Context) {
	state := ctx.Query("state")

	if !(state == "login" || state == "register") {
		ctx.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("%s/auth/register?toast=Something+went+wrong", h.config.URL.Frontend))
		return
	}

	code := ctx.Query("code")
	token, err := h.config.Google.Exchange(context.Background(), code)
	if err != nil {
		ctx.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("%s/auth/register?toast=Something+went+wrong", h.config.URL.Frontend))
		return
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		ctx.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("%s/auth/register?toast=Something+went+wrong", h.config.URL.Frontend))
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		ctx.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("%s/auth/register?toast=Something+went+wrong", h.config.URL.Frontend))
		return
	}

	if state == "register" {
		var req dto.GoogleRegisterCallback
		json.NewDecoder(response.Body).Decode(&req)

		err := h.useCase.GoogleRegisterCallback(ctx, req)
		if err != nil {
			ctx.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("%s/auth/login?toast=Email+has+been+registered", h.config.URL.Frontend))
			return
		}

		ctx.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("%s/auth/login", h.config.URL.Frontend))
		return
	}

	var req dto.GoogleLoginCallback
	json.NewDecoder(response.Body).Decode(&req)

	data, err := h.useCase.GoogleLoginCallback(ctx, req)
	if err != nil {
		ctx.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("%s/auth/login?toast=Email+is+not+registered", h.config.URL.Frontend))
		return
	}

	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie("token", data.Token, 3600, "/", h.config.URL.DomainName, true, true)
	ctx.SetCookie("is_verified", data.IsVerified, 3600, "/", h.config.URL.DomainName, true, true)
	ctx.SetCookie("role", data.Role, 3600, "/", h.config.URL.DomainName, true, true)
	ctx.Redirect(http.StatusPermanentRedirect, h.config.Client.AuthRedirectURL+fmt.Sprintf("?is_verified=%s&role=%s&user_id=%d", data.IsVerified, data.Role, data.UserID))
}
