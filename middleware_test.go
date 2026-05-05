package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestULIDMiddleware(t *testing.T) {
	e := echo.New()
	
	// Create a handler that checks if RequestID is set
	handler := func(c echo.Context) error {
		requestID := c.Get("RequestID")
		assert.NotNil(t, requestID)
		assert.NotEmpty(t, requestID)
		
		// Check if header is set
		headerRequestID := c.Response().Header().Get("X-Request-ID")
		assert.Equal(t, requestID, headerRequestID)
		
		return c.NoContent(http.StatusOK)
	}

	// Apply middleware
	middleware := ULIDMiddleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := middleware(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestULIDMiddleware_Uniqueness(t *testing.T) {
	e := echo.New()
	
	var requestIDs []string
	
	// Create a handler that collects RequestIDs
	handler := func(c echo.Context) error {
		requestID := c.Get("RequestID").(string)
		requestIDs = append(requestIDs, requestID)
		return c.NoContent(http.StatusOK)
	}

	middleware := ULIDMiddleware(handler)

	// Make multiple requests to ensure uniqueness
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := middleware(c)
		assert.NoError(t, err)
	}

	// Check that all RequestIDs are unique
	uniqueIDs := make(map[string]bool)
	for _, id := range requestIDs {
		assert.False(t, uniqueIDs[id], "RequestID should be unique: %s", id)
		uniqueIDs[id] = true
	}
	
	assert.Len(t, requestIDs, 10)
	assert.Len(t, uniqueIDs, 10)
}
