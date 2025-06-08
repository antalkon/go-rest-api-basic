package handlers

import (
	"backend/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type PingHandler struct {
	svc *service.PingService
}

func NewPingHandler(svc *service.PingService) *PingHandler {
	return &PingHandler{svc: svc}
}

func (h *PingHandler) Ping(c echo.Context) error {
	ip := c.RealIP()
	if err := h.svc.SavePing(c.Request().Context(), ip); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not save ping"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "pong"})
}

func (h *PingHandler) GetAll(c echo.Context) error {
	list, err := h.svc.GetPings(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get pings"})
	}
	return c.JSON(http.StatusOK, list)
}
