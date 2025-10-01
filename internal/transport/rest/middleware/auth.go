package middleware

import (
	"backend/pkg/auth"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// AuthMiddleware проверяет access токен и сохраняет userID в контексте
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "missing Authorization header")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid Authorization header format")
		}
		token := parts[1]

		claims, err := auth.ParseAccessToken(token)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
		}

		// Сохраняем userID в context
		c.Set("userID", claims.UserID)

		return next(c)
	}
}
