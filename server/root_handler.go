package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RootHandler(c echo.Context) error {
	return c.JSON(
		http.StatusBadRequest,
		map[string]interface{}{
			"message": "Method Not Allowed",
		},
	)
}

func RobotsTXTHandler(c echo.Context) error {
	return c.String(http.StatusBadRequest, "User-agent: *\nDisallow: /")
}
