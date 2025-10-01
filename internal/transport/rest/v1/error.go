package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleError(c echo.Context, err error) error {
	return c.JSON(http.StatusInternalServerError, map[string]string{
		"error": err.Error(),
	})
}
