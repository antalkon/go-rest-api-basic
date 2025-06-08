package v1

import (
	"backend/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	ping *service.PingService
}

func NewHandler(pingService *service.PingService) *Handler {
	return &Handler{ping: pingService}
}

func (h *Handler) Ping(c echo.Context) error {
	ip := c.RealIP()
	if err := h.ping.SavePing(c.Request().Context(), ip); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not save ping"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "pong"})
}

func (h *Handler) GetAllPings(c echo.Context) error {
	data, err := h.ping.GetPings(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get pings"})
	}
	return c.JSON(http.StatusOK, data)
}
