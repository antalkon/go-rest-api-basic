// handlers/auth.go
package handlers

import (
	"backend/intenal/service"
	"backend/intenal/transport/rest/v1/req"
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var r req.RegisterRequest
	if err := c.Bind(&r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid input")
	}

	if err := r.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	access, refresh, err := h.svc.Register(c.Request().Context(), r)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return setAuthTokens(c, access, refresh)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var r req.LoginRequest
	if err := c.Bind(&r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid input")
	}

	if err := r.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	access, refresh, err := h.svc.Login(context.Background(), r)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	return setAuthTokens(c, access, refresh)
}
func setAuthTokens(c echo.Context, access, refresh string) error {
	// Set Access Token в header — клиент сам может его сохранить
	c.Response().Header().Set("Authorization", "Bearer "+access)

	// Set Refresh Token в HttpOnly secure cookie
	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    refresh,
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // обязательно для HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   7 * 24 * 60 * 60, // 7 дней
	})

	// Можно вернуть что-то для фронта (например, user info), но токены — в куке и заголовке
	return c.NoContent(http.StatusOK)
}
