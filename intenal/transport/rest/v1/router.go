package v1

import (
	"backend/intenal/transport/rest/v1/handlers"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group, h *handlers.PingHandler) {
	g.GET("/ping", h.Ping)
	g.GET("/ping/all", h.GetAll)
}
