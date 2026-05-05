package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRouteSysGetPing(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/system/ping", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := RouteSysGetPing(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response GenericAPIResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Status)
	assert.Equal(t, "OK", response.Body)
}

func TestRouteSysGetConfig(t *testing.T) {
	// Setup global config for testing
	G = &Global{
		Config: &Config{},
	}
	G.Config.Core.TestMode = true

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/system/config", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := RouteSysGetConfig(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response GenericAPIResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Status)
	assert.NotNil(t, response.Body)
}

func TestRouteSysGetMetrics(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/system/metrics", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := RouteSysGetMetrics(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "# OK", rec.Body.String())
}
