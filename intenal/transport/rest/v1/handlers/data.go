package handlers

import (
	"backend/intenal/service"
	"backend/intenal/transport/rest/v1/res"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type DataHandler struct {
	svc *service.DataService
}

func NewDataHandler(svc *service.DataService) *DataHandler {
	return &DataHandler{svc: svc}
}

// GetUserData godoc
// @Summary     Get user data
// @Description Returns username and email of the authorized user
// @Tags        data
// @Security    BearerAuth
// @Produce     json
// @Success     200 {object} res.UserData
// @Failure     401 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /data/user [get]
func (h *DataHandler) GetUserData(c echo.Context) error {
	userIDRaw := c.Get("userID")
	userUUID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "userID not found or invalid in context")
	}

	username, email, err := h.svc.GetUserData(c.Request().Context(), userUUID.String())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, res.UserData{
		Username: username,
		Email:    email,
	})
}
