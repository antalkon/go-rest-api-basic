package handlers

import (
	"backend/intenal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type PingHandler struct {
	svc *service.PingService
}

func NewPingHandler(svc *service.PingService) *PingHandler {
	return &PingHandler{svc: svc}
}

// Ping godoc
// @Summary     Ping the server
// @Description Returns pong and saves IP to database
// @Tags        ping
// @Produce     json
// @Success     200 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /ping [get]
func (h *PingHandler) Ping(c echo.Context) error {
	ip := c.RealIP()
	if err := h.svc.SavePing(c.Request().Context(), ip); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not save ping"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "pong"})
}

// GetAll godoc
// @Summary     Get all pings
// @Description Returns list of all pings from DB
// @Tags        ping
// @Produce     json
// @Success     200 {array} string
// @Failure     500 {object} map[string]string
// @Router      /ping/all [get]
func (h *PingHandler) GetAll(c echo.Context) error {
	list, err := h.svc.GetPings(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get pings"})
	}
	return c.JSON(http.StatusOK, list)
}
