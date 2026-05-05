package main

import (
	"github.com/labstack/echo/v4"
	"github.com/oklog/ulid/v2"
)

func ULIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ulid := ulid.Make()
		c.Set("RequestID", ulid.String())
		c.Response().Header().Set("X-Request-ID", ulid.String())
		return next(c)
	}
}
