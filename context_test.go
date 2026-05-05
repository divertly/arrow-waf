package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestServerContext_ServerName(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	serverCtx := &ServerContext{
		Context:    c,
		serverName: "test-server",
	}

	assert.Equal(t, "test-server", serverCtx.ServerName())
}

func TestNewEchoContext(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	serverName := "public-server"
	newCtx := newEchoContext(c, serverName)

	assert.IsType(t, &ServerContext{}, newCtx)
	serverCtx := newCtx.(*ServerContext)
	assert.Equal(t, serverName, serverCtx.ServerName())
	assert.Equal(t, c, serverCtx.Context)
}
