package middleware

import (
	"time"

	"backend/pkg/logger"

	"github.com/labstack/echo/v4"
)

func RequestLogger() echo.MiddlewareFunc {
	log := logger.L()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			stop := time.Now()

			req := c.Request()
			res := c.Response()

			log.Info("HTTP",
				"method", req.Method,
				"path", req.URL.Path,
				"status", res.Status,
				"duration", stop.Sub(start).String(),
				"ip", c.RealIP(),
			)

			return err
		}
	}
}
