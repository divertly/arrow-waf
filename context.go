package main

import "github.com/labstack/echo/v4"

type ServerContext struct {
	echo.Context
	serverName string
}

func (c *ServerContext) ServerName() string {
	return c.serverName
}

func newEchoContext(e echo.Context, serverName string) echo.Context {
	return &ServerContext{e, serverName}
}
