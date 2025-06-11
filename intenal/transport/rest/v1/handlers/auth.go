// handlers/auth.go
package handlers

import (
	"backend/intenal/service"
	"backend/intenal/transport/rest/v1/req"
	"context"
	"fmt"
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
		Path:     "localhost:8080",
		HttpOnly: true,
		Secure:   false, // обязательно для HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   7 * 24 * 60 * 60, // 7 дней
	})

	// Можно вернуть что-то для фронта (например, user info), но токены — в куке и заголовке
	return c.NoContent(http.StatusOK)
}

func (h *AuthHandler) Refresh(c echo.Context) error {

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "refresh token not found")
	}
	fmt.Println("Refresh token:", refreshToken.Value)
	access, rerefreshToken, err := h.svc.Refresh(c.Request().Context(), refreshToken.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	return setAuthTokens(c, access, rerefreshToken)
}

func (h *AuthHandler) Logout(c echo.Context) error {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "refresh token not found")
	}

	if err := h.svc.Logout(c.Request().Context(), refreshToken.Value); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Удаляем куку с refresh токеном
	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "localhost:8080", // путь, где кука будет удалена
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1, // удалить куку
	})
	return c.NoContent(http.StatusOK)
}
