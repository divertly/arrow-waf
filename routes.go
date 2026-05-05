package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RouteSysGetPing(c echo.Context) error {
	return c.JSON(http.StatusOK, &GenericAPIResponse{
		Status: http.StatusOK,
		Body:   "OK",
	})
}

func RouteSysGetConfig(c echo.Context) error {
	return c.JSON(http.StatusOK, &GenericAPIResponse{
		Status: http.StatusOK,
		Body:   G.Config,
	})
}

func RouteSysGetMetrics(c echo.Context) error {
	return c.String(200, "# OK")
}
